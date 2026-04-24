package middleware

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/arda-labs/arda/arda-be-go/services/iam-service/internal/conf"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware"
	"github.com/go-kratos/kratos/v2/transport"
	"github.com/golang-jwt/jwt/v5"
	"github.com/lestrrat-go/jwx/v2/jwk"
)

type contextKey string

const (
	claimsKey   contextKey = "iam_claims"
	subjectKey  contextKey = "iam_subject"
	tenantIDKey contextKey = "iam_tenant_id"
	emailKey    contextKey = "iam_email"
)

type IAMClaims struct {
	Sub      string
	TenantID string
	Email    string
	Name     string
	Roles    []string
}

func SubjectFromContext(ctx context.Context) string {
	v, _ := ctx.Value(subjectKey).(string)
	return v
}

func TenantIDFromContext(ctx context.Context) string {
	v, _ := ctx.Value(tenantIDKey).(string)
	return v
}

func EmailFromContext(ctx context.Context) string {
	v, _ := ctx.Value(emailKey).(string)
	return v
}

func IAMClaimsFromContext(ctx context.Context) *IAMClaims {
	v, _ := ctx.Value(claimsKey).(*IAMClaims)
	return v
}

func iamContext(claims *IAMClaims) context.Context {
	ctx := context.Background()
	ctx = context.WithValue(ctx, claimsKey, claims)
	ctx = context.WithValue(ctx, subjectKey, claims.Sub)
	ctx = context.WithValue(ctx, tenantIDKey, claims.TenantID)
	ctx = context.WithValue(ctx, emailKey, claims.Email)
	return ctx
}

type jwksFetcher struct {
	mu       sync.RWMutex
	set      jwk.Set
	cached   time.Time
	ttl      time.Duration
	endpoint string
}

func newJWKSFetcher(endpoint string) *jwksFetcher {
	return &jwksFetcher{endpoint: endpoint, ttl: 1 * time.Hour}
}

func (f *jwksFetcher) fetch() (jwk.Set, error) {
	f.mu.RLock()
	if f.set != nil && time.Since(f.cached) < f.ttl {
		s := f.set
		f.mu.RUnlock()
		return s, nil
	}
	f.mu.RUnlock()

	f.mu.Lock()
	defer f.mu.Unlock()
	if f.set != nil && time.Since(f.cached) < f.ttl {
		return f.set, nil
	}

	set, err := jwk.Fetch(context.Background(), f.endpoint)
	if err != nil {
		return nil, fmt.Errorf("fetching JWKS: %w", err)
	}
	f.set = set
	f.cached = time.Now()
	return set, nil
}

func fetchUserInfo(issuer, token string) (*IAMClaims, error) {
	userInfoURL := strings.TrimRight(issuer, "/") + "/oidc/v1/userinfo"
	req, err := http.NewRequest("GET", userInfoURL, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+token)

	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("userinfo request failed with status %d", resp.StatusCode)
	}

	var claims map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&claims); err != nil {
		return nil, err
	}

	iamClaims := &IAMClaims{
		Sub:      stringVal(claims["sub"]),
		TenantID: stringVal(claims["urn:zitadel:iam:org:id"]), // Zitadel standard claim for org/tenant
		Email:    stringVal(claims["email"]),
		Name:     stringVal(claims["name"]),
	}
	
	if tenantID := stringVal(claims["tenant_id"]); tenantID != "" {
		iamClaims.TenantID = tenantID
	}

	if roles, ok := claims["urn:zitadel:iam:org:project:roles"].(map[string]interface{}); ok {
		for r := range roles {
			iamClaims.Roles = append(iamClaims.Roles, r)
		}
	} else if roles, ok := claims["roles"].([]interface{}); ok {
		for _, r := range roles {
			if s, ok := r.(string); ok {
				iamClaims.Roles = append(iamClaims.Roles, s)
			}
		}
	}

	return iamClaims, nil
}

func JWTMiddleware(jwtConf *conf.JWT) middleware.Middleware {
	fetcher := newJWKSFetcher(jwtConf.JwksEndpoint)

	return func(handler middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req interface{}) (interface{}, error) {
			tr, ok := transport.FromServerContext(ctx)
			if !ok {
				log.Warn("JWTMiddleware transport missing")
				return handler(ctx, req)
			}

			tokenStr := tr.RequestHeader().Get("Authorization")
			if len(tokenStr) > 7 && strings.EqualFold(tokenStr[:7], "Bearer ") {
				tokenStr = tokenStr[7:]
			}
			if tokenStr == "" {
				return handler(ctx, req)
			}

			var iamClaims *IAMClaims

			// Check if token is likely an opaque token (doesn't have two dots) or starts with V2_
			if strings.HasPrefix(tokenStr, "V2_") || strings.Count(tokenStr, ".") != 2 {
				claims, err := fetchUserInfo(jwtConf.Issuer, tokenStr)
				if err != nil {
					log.Warnf("JWTMiddleware userinfo fetch error: %v", err)
					return handler(ctx, req)
				}
				iamClaims = claims
			} else {
				set, err := fetcher.fetch()
				if err != nil {
					log.Warnf("JWTMiddleware JWKS fetch error: %v", err)
					return handler(ctx, req)
				}

				tok, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
					kid, ok := t.Header["kid"].(string)
					if !ok {
						return nil, fmt.Errorf("missing kid")
					}
					key, found := set.LookupKeyID(kid)
					if !found {
						return nil, fmt.Errorf("key %s not found", kid)
					}
					var raw interface{}
					if err := key.Raw(&raw); err != nil {
						return nil, err
					}
					return raw, nil
				})
				if err != nil || tok == nil {
					log.Warnf("JWTMiddleware parse error: %v", err)
					return handler(ctx, req)
				}

				claims, ok := tok.Claims.(jwt.MapClaims)
				if !ok {
					log.Warn("JWTMiddleware claims format error")
					return handler(ctx, req)
				}

				iamClaims = &IAMClaims{
					Sub:      stringVal(claims["sub"]),
					TenantID: stringVal(claims["tenant_id"]),
					Email:    stringVal(claims["email"]),
					Name:     stringVal(claims["name"]),
				}
				
				if tenantID := stringVal(claims["urn:zitadel:iam:org:id"]); tenantID != "" && iamClaims.TenantID == "" {
					iamClaims.TenantID = tenantID
				}
				
				if roles, ok := claims["urn:zitadel:iam:org:project:roles"].(map[string]interface{}); ok {
					for r := range roles {
						iamClaims.Roles = append(iamClaims.Roles, r)
					}
				} else if roles, ok := claims["roles"].([]interface{}); ok {
					for _, r := range roles {
						if s, ok := r.(string); ok {
							iamClaims.Roles = append(iamClaims.Roles, s)
						}
					}
				}
			}

			if iamClaims != nil {
				ctx = context.WithValue(ctx, claimsKey, iamClaims)
				ctx = context.WithValue(ctx, subjectKey, iamClaims.Sub)
				ctx = context.WithValue(ctx, tenantIDKey, iamClaims.TenantID)
				ctx = context.WithValue(ctx, emailKey, iamClaims.Email)
			}

			return handler(ctx, req)
		}
	}
}

func stringVal(v interface{}) string {
	if v == nil {
		return ""
	}
	if s, ok := v.(string); ok {
		return s
	}
	return fmt.Sprintf("%v", v)
}
