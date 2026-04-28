package middleware

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/middleware"
	"github.com/go-kratos/kratos/v2/transport"
	"github.com/golang-jwt/jwt/v5"
	"github.com/lestrrat-go/jwx/v2/jwk"
)

const (
	userIDKey          contextKey = "user_id"
	userEmailKey       contextKey = "user_email"
	userPermissionsKey contextKey = "user_permissions"
)

type AuthOption func(*authConfig)

type authConfig struct {
	jwksEndpoint string
	issuer       string
	audience     string
}

type jwtVerifier struct {
	mu       sync.RWMutex
	set      jwk.Set
	cached   time.Time
	ttl      time.Duration
	endpoint string
}

func WithJWTValidation(jwksEndpoint, issuer, audience string) AuthOption {
	return func(c *authConfig) {
		c.jwksEndpoint = jwksEndpoint
		c.issuer = issuer
		c.audience = audience
	}
}

// Auth extracts user information from APISIX headers or verifies Bearer tokens.
// Without options it keeps the previous local-dev behavior and parses JWTs
// without verification. Services should pass WithJWTValidation in production.
func Auth(options ...AuthOption) middleware.Middleware {
	cfg := authConfig{}
	for _, option := range options {
		option(&cfg)
	}

	var verifier *jwtVerifier
	if cfg.jwksEndpoint != "" {
		verifier = &jwtVerifier{endpoint: cfg.jwksEndpoint, ttl: time.Hour}
	}

	return func(handler middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req interface{}) (interface{}, error) {
			if tr, ok := transport.FromServerContext(ctx); ok {
				userID := tr.RequestHeader().Get("X-User-Id")
				email := tr.RequestHeader().Get("X-User-Email")
				perms := tr.RequestHeader().Get("X-User-Permissions")

				if userID == "" {
					authHeader := tr.RequestHeader().Get("Authorization")
					if strings.HasPrefix(authHeader, "Bearer ") {
						tokenStr := authHeader[7:]
						claims, err := claimsFromToken(tokenStr, cfg, verifier)
						if err != nil {
							return nil, errors.Unauthorized("UNAUTHORIZED", "invalid bearer token")
						}
						if claims != nil {
							userID, email, perms = claimsToIdentity(claims)
							if tenantID := tr.RequestHeader().Get("X-Tenant-ID"); tenantID == "" {
								setTenantHeaderFromClaims(tr, claims)
							}
						}
					}
				}

				if userID != "" {
					ctx = context.WithValue(ctx, userIDKey, userID)
					ctx = context.WithValue(ctx, userEmailKey, email)
					ctx = context.WithValue(ctx, userPermissionsKey, perms)
				}
			}
			return handler(ctx, req)
		}
	}
}

func (v *jwtVerifier) getSet() (jwk.Set, error) {
	v.mu.RLock()
	if v.set != nil && time.Since(v.cached) < v.ttl {
		set := v.set
		v.mu.RUnlock()
		return set, nil
	}
	v.mu.RUnlock()

	v.mu.Lock()
	defer v.mu.Unlock()
	if v.set != nil && time.Since(v.cached) < v.ttl {
		return v.set, nil
	}

	set, err := jwk.Fetch(context.Background(), v.endpoint)
	if err != nil {
		return nil, fmt.Errorf("fetching JWKS: %w", err)
	}
	v.set = set
	v.cached = time.Now()
	return set, nil
}

func claimsFromToken(tokenStr string, cfg authConfig, verifier *jwtVerifier) (jwt.MapClaims, error) {
	if verifier == nil {
		token, _, _ := new(jwt.Parser).ParseUnverified(tokenStr, jwt.MapClaims{})
		if token == nil {
			return nil, nil
		}
		claims, _ := token.Claims.(jwt.MapClaims)
		return claims, nil
	}

	if strings.Count(tokenStr, ".") != 2 {
		return fetchUserInfo(cfg.issuer, tokenStr)
	}

	set, err := verifier.getSet()
	if err != nil {
		return nil, err
	}

	parseOptions := []jwt.ParserOption{
		jwt.WithExpirationRequired(),
		jwt.WithValidMethods([]string{"RS256"}),
	}
	if cfg.issuer != "" {
		parseOptions = append(parseOptions, jwt.WithIssuer(cfg.issuer))
	}
	if cfg.audience != "" {
		parseOptions = append(parseOptions, jwt.WithAudience(cfg.audience))
	}

	token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
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
	}, parseOptions...)
	if err != nil || token == nil || !token.Valid {
		return nil, fmt.Errorf("invalid JWT: %w", err)
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, fmt.Errorf("invalid claims format")
	}
	return claims, nil
}

func fetchUserInfo(issuer, token string) (jwt.MapClaims, error) {
	userInfoURL := strings.TrimRight(issuer, "/") + "/oidc/v1/userinfo"
	req, err := http.NewRequest(http.MethodGet, userInfoURL, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := (&http.Client{Timeout: 5 * time.Second}).Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("userinfo request failed with status %d", resp.StatusCode)
	}

	var claims jwt.MapClaims
	if err := json.NewDecoder(resp.Body).Decode(&claims); err != nil {
		return nil, err
	}
	return claims, nil
}

func claimsToIdentity(claims jwt.MapClaims) (string, string, string) {
	userID, _ := claims["sub"].(string)
	email, _ := claims["email"].(string)

	var permissions []string
	if roles, ok := claims["urn:zitadel:iam:org:project:roles"].(map[string]any); ok {
		for role := range roles {
			permissions = append(permissions, role)
		}
	} else if roles, ok := claims["roles"].([]any); ok {
		for _, role := range roles {
			if roleName, ok := role.(string); ok {
				permissions = append(permissions, roleName)
			}
		}
	}

	return userID, email, strings.Join(permissions, ",")
}

func setTenantHeaderFromClaims(tr transport.Transporter, claims jwt.MapClaims) {
	if tenantID, ok := claims["tenant_id"].(string); ok {
		tr.RequestHeader().Set("X-Tenant-ID", tenantID)
		return
	}
	if tenantID, ok := claims["urn:zitadel:iam:org:id"].(string); ok {
		tr.RequestHeader().Set("X-Tenant-ID", tenantID)
	}
}

func GetUserID(ctx context.Context) string {
	if v, ok := ctx.Value(userIDKey).(string); ok {
		return v
	}
	return ""
}

func GetEmail(ctx context.Context) string {
	if v, ok := ctx.Value(userEmailKey).(string); ok {
		return v
	}
	return ""
}
