package biz

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/arda-labs/arda/arda-be-go/services/iam-service/internal/conf"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/golang-jwt/jwt/v5"
	"github.com/lestrrat-go/jwx/v2/jwk"
)

type AuthUsecase struct {
	conf          *conf.Zitadel
	jwtConf       *conf.JWT
	jwks          *jwkCache
	perms         *PermissionUsecase
	loginPAT      string
	managementPAT string
	log           *log.Helper
}

// jwkCache giữ bản cache JWKS set với TTL 1 giờ
type jwkCache struct {
	mu       sync.RWMutex
	set      jwk.Set
	cached   time.Time
	ttl      time.Duration
	endpoint string
}

func newJWKCache(endpoint string) *jwkCache {
	return &jwkCache{endpoint: endpoint, ttl: time.Hour}
}

func (c *jwkCache) getSet() (jwk.Set, error) {
	c.mu.RLock()
	if c.set != nil && time.Since(c.cached) < c.ttl {
		s := c.set
		c.mu.RUnlock()
		return s, nil
	}
	c.mu.RUnlock()

	c.mu.Lock()
	defer c.mu.Unlock()
	if c.set != nil && time.Since(c.cached) < c.ttl {
		return c.set, nil
	}
	set, err := jwk.Fetch(context.Background(), c.endpoint)
	if err != nil {
		return nil, fmt.Errorf("fetching JWKS from %s: %w", c.endpoint, err)
	}
	c.set = set
	c.cached = time.Now()
	return set, nil
}

func NewAuthUsecase(zitadelConf *conf.Zitadel, jwtConf *conf.JWT, perms *PermissionUsecase, logger log.Logger) *AuthUsecase {
	loginPAT := firstConfiguredPAT(
		os.Getenv("ZITADEL_LOGIN_CLIENT_PAT"),
		zitadelConf.GetPat(),
	)
	managementPAT := firstConfiguredPAT(
		os.Getenv("ZITADEL_MANAGEMENT_PAT"),
		os.Getenv("ZITADEL_PAT"),
	)

	return &AuthUsecase{
		conf:          zitadelConf,
		jwtConf:       jwtConf,
		jwks:          newJWKCache(jwtConf.JwksEndpoint),
		perms:         perms,
		loginPAT:      loginPAT,
		managementPAT: managementPAT,
		log:           log.NewHelper(logger),
	}
}

func firstConfiguredPAT(values ...string) string {
	for _, value := range values {
		pat := strings.TrimSpace(value)
		if pat != "" && !strings.Contains(pat, "${") {
			return pat
		}
	}
	return ""
}

func (uc *AuthUsecase) HasZitadelManagementPAT() bool {
	return uc.managementPAT != ""
}

type PasswordPolicy struct {
	MinLength        int
	RequireUppercase bool
	RequireLowercase bool
	RequireNumber    bool
	RequireSymbol    bool
}

type LoginPolicy struct {
	PasswordLoginEnabled bool
	ExternalIDPEnabled   bool
	MFARequired          bool
}

type AuthSettings struct {
	TenantID       string
	AuthMode       string
	Provider       string
	PasswordPolicy PasswordPolicy
	LoginPolicy    LoginPolicy
}

type ZitadelAPIError struct {
	StatusCode int
	Status     string
	Code       string
	Message    string
}

func (e *ZitadelAPIError) Error() string {
	if e.Message != "" {
		return e.Message
	}
	if e.Status != "" {
		return "zitadel error: " + e.Status
	}
	return "zitadel error"
}

func (uc *AuthUsecase) GetAuthSettings(ctx context.Context, tenantID, authMode string) (*AuthSettings, error) {
	passwordPolicy, err := uc.getZitadelPasswordPolicy(ctx)
	if err != nil {
		return nil, err
	}

	return &AuthSettings{
		TenantID:       tenantID,
		AuthMode:       authMode,
		Provider:       "ZITADEL",
		PasswordPolicy: passwordPolicy,
		LoginPolicy: LoginPolicy{
			PasswordLoginEnabled: true,
			ExternalIDPEnabled:   true,
			MFARequired:          false,
		},
	}, nil
}

type sessionResponse struct {
	SessionID    string `json:"sessionId"`
	SessionToken string `json:"sessionToken"`
}

type finalizeResponse struct {
	CallbackURL string `json:"callbackUrl"`
}

/**
 * CustomLogin thực hiện flow "Login App" chuẩn Zitadel:
 * 1. Tạo session liên kết với authRequestID
 * 2. Xác thực password cho session dùng PAT của Login Client
 * 3. Finalize OIDC request dùng PAT của Login Client
 */
func (uc *AuthUsecase) CustomLogin(ctx context.Context, loginName, password, authRequestID string) (string, error) {
	uc.log.Infof("Login App: Starting auth for %s, request: %s", loginName, authRequestID)

	// 1. Create Session
	session, err := uc.createSession(loginName, authRequestID)
	if err != nil {
		return "", err
	}

	// 2. Verify Password
	sessionToken, err := uc.verifyPassword(session.SessionID, session.SessionToken, password)
	if err != nil {
		return "", err
	}

	// 3. Finalize OIDC
	callbackURL, err := uc.finalizeAuthRequest(authRequestID, session.SessionID, sessionToken)
	if err != nil {
		return "", err
	}

	uc.log.Infof("Login App: Success. Redirecting to callback")
	return callbackURL, nil
}

/**
 * ForwardAuth thực hiện kiểm tra quyền hạn cho APISIX Gateway:
 * 1. Validate access token từ Zitadel
 * 2. Trích xuất userID và tenantID từ claims
 * 3. (TODO) Kiểm tra quyền RBAC/ABAC trong DB
 */
func (uc *AuthUsecase) ForwardAuth(ctx context.Context, method, path, token, selectedTenantID string) (bool, string, string, error) {
	// Strip "Bearer " prefix nếu có
	if len(token) > 7 && token[:7] == "Bearer " {
		token = token[7:]
	}
	if token == "" {
		return false, "", "", fmt.Errorf("missing token")
	}

	claims, err := uc.claimsFromAccessToken(token)
	if err != nil {
		uc.log.Warnf("ForwardAuth: invalid token for %s %s: %v", method, path, err)
		return false, "", "", err
	}

	userID := stringClaimVal(claims["sub"])
	tenantID := selectedTenantID
	if tenantID == "" {
		tenantID = stringClaimVal(claims["tenant_id"])
	}
	if tenantID == "" {
		tenantID = stringClaimVal(claims["urn:zitadel:iam:org:id"])
	}

	// Map Method+Path sang Resource+Action để kiểm tra quyền
	resource, action := uc.mapPathToAction(method, path)
	uc.log.Infof("ForwardAuth: user=%s tenant=%s checking %s:%s for %s %s", userID, tenantID, resource, action, method, path)

	// Phase 4.1: Thực hiện kiểm tra quyền thực tế
	allowed, source, err := uc.perms.CheckPermission(ctx, userID, tenantID, resource, action, "")
	if err != nil {
		uc.log.Errorf("ForwardAuth: CheckPermission error: %v", err)
		return false, userID, tenantID, err
	}

	if !allowed {
		uc.log.Warnf("ForwardAuth: permission denied for user=%s tenant=%s resource=%s action=%s", userID, tenantID, resource, action)
		return false, userID, tenantID, nil
	}

	uc.log.Infof("ForwardAuth: user=%s tenant=%s ALLOWED by %s", userID, tenantID, source)
	return true, userID, tenantID, nil
}

func (uc *AuthUsecase) claimsFromAccessToken(token string) (jwt.MapClaims, error) {
	if strings.HasPrefix(token, "V2_") || strings.Count(token, ".") != 2 {
		return uc.fetchUserInfoClaims(token)
	}

	set, err := uc.jwks.getSet()
	if err != nil {
		return nil, err
	}

	parseOptions := []jwt.ParserOption{
		jwt.WithExpirationRequired(),
		jwt.WithValidMethods([]string{"RS256"}),
	}
	if uc.jwtConf.Issuer != "" {
		parseOptions = append(parseOptions, jwt.WithIssuer(uc.jwtConf.Issuer))
	}
	if uc.jwtConf.Audience != "" {
		parseOptions = append(parseOptions, jwt.WithAudience(uc.jwtConf.Audience))
	}

	// Verify JWT signature and standard claims with JWKS.
	tok, err := jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
		kid, ok := t.Header["kid"].(string)
		if !ok {
			return nil, fmt.Errorf("missing kid in JWT header")
		}
		key, found := set.LookupKeyID(kid)
		if !found {
			return nil, fmt.Errorf("key %s not found in JWKS", kid)
		}
		var raw interface{}
		if err := key.Raw(&raw); err != nil {
			return nil, err
		}
		return raw, nil
	}, parseOptions...)
	if err != nil || !tok.Valid {
		return nil, fmt.Errorf("invalid token: %w", err)
	}

	claims, ok := tok.Claims.(jwt.MapClaims)
	if !ok {
		return nil, fmt.Errorf("invalid claims format")
	}
	return claims, nil
}

func (uc *AuthUsecase) fetchUserInfoClaims(token string) (jwt.MapClaims, error) {
	issuer := strings.TrimRight(uc.jwtConf.Issuer, "/")
	if issuer == "" {
		issuer = strings.TrimRight(uc.conf.Authority, "/")
	}
	if issuer == "" {
		return nil, fmt.Errorf("missing Zitadel issuer")
	}

	req, err := http.NewRequest(http.MethodGet, issuer+"/oidc/v1/userinfo", nil)
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

func stringClaimVal(v interface{}) string {
	if v == nil {
		return ""
	}
	if s, ok := v.(string); ok {
		return s
	}
	return fmt.Sprintf("%v", v)
}

type routeRule struct {
	method   string
	pattern  *regexp.Regexp
	resource string
	action   string
}

var routeRules = []routeRule{
	// Users
	{http.MethodGet, regexp.MustCompile(`^/v1/users`), "user", "read"},
	{http.MethodPost, regexp.MustCompile(`^/v1/users$`), "user", "create"},
	{http.MethodPut, regexp.MustCompile(`^/v1/users/`), "user", "update"},
	{http.MethodDelete, regexp.MustCompile(`^/v1/users/`), "user", "delete"},
	// Tenants
	{http.MethodGet, regexp.MustCompile(`^/v1/tenants`), "tenant", "read"},
	{http.MethodPost, regexp.MustCompile(`^/v1/tenants$`), "tenant", "create"},
	{http.MethodPut, regexp.MustCompile(`^/v1/tenants/`), "tenant", "update"},
	{http.MethodDelete, regexp.MustCompile(`^/v1/tenants/`), "tenant", "delete"},
	// Members
	{http.MethodGet, regexp.MustCompile(`^/v1/members`), "member", "read"},
	{http.MethodPost, regexp.MustCompile(`^/v1/members`), "member", "invite"},
	{http.MethodDelete, regexp.MustCompile(`^/v1/members/`), "member", "remove"},
	// Roles
	{http.MethodGet, regexp.MustCompile(`^/v1/roles`), "role", "read"},
	{http.MethodPost, regexp.MustCompile(`^/v1/roles$`), "role", "create"},
	{http.MethodPut, regexp.MustCompile(`^/v1/roles/`), "role", "update"},
	{http.MethodDelete, regexp.MustCompile(`^/v1/roles/`), "role", "delete"},
	// Permissions / Approvals
	{http.MethodGet, regexp.MustCompile(`^/v1/permissions`), "permission", "read"},
	{http.MethodPost, regexp.MustCompile(`^/v1/permissions`), "permission", "grant"},
	{http.MethodGet, regexp.MustCompile(`^/v1/approvals`), "approval", "read"},
	{http.MethodPost, regexp.MustCompile(`^/v1/approvals/`), "approval", "manage"},
	// Self-service — luôn cho qua (auth middleware đã xác thực token)
	{http.MethodGet, regexp.MustCompile(`^/v1/me`), "me", "read"},
}

func (uc *AuthUsecase) mapPathToAction(method, path string) (string, string) {
	for _, rule := range routeRules {
		if rule.method == method && rule.pattern.MatchString(path) {
			return rule.resource, rule.action
		}
	}
	// Fallback: public/unknown → dùng resource="public" action="access"
	return "public", "access"
}

func (uc *AuthUsecase) createSession(loginName, authRequestID string) (*sessionResponse, error) {
	url := fmt.Sprintf("%s/v2/sessions", uc.conf.Authority)
	body := map[string]any{
		"checks": map[string]any{
			"user": map[string]string{"loginName": loginName},
		},
		"authRequestId": authRequestID,
	}

	var resp sessionResponse
	if err := uc.callZitadelWithPAT(http.MethodPost, url, body, &resp, uc.loginPAT, "login client"); err != nil {
		return nil, err
	}
	return &resp, nil
}

func (uc *AuthUsecase) verifyPassword(sessionID, sessionToken, password string) (string, error) {
	url := fmt.Sprintf("%s/v2/sessions/%s", uc.conf.Authority, sessionID)
	body := map[string]any{
		"checks": map[string]any{
			"password": map[string]string{"password": password},
		},
		"sessionToken": sessionToken,
	}

	var resp sessionResponse
	if err := uc.callZitadelWithPAT(http.MethodPatch, url, body, &resp, uc.loginPAT, "login client"); err != nil {
		return "", err
	}
	return resp.SessionToken, nil
}

func (uc *AuthUsecase) finalizeAuthRequest(authRequestID, sessionID, sessionToken string) (string, error) {
	url := fmt.Sprintf("%s/v2/oidc/auth_requests/%s", uc.conf.Authority, authRequestID)
	body := map[string]any{
		"session": map[string]string{
			"sessionId":    sessionID,
			"sessionToken": sessionToken,
		},
	}

	var resp finalizeResponse
	if err := uc.callZitadelWithPAT(http.MethodPost, url, body, &resp, uc.loginPAT, "login client"); err != nil {
		return "", err
	}
	return resp.CallbackURL, nil
}

type zitadelCreateUserResponse struct {
	UserID string `json:"userId"`
}

type zitadelPasswordComplexityResponse struct {
	Policy struct {
		MinLength    string `json:"minLength"`
		HasUppercase bool   `json:"hasUppercase"`
		HasLowercase bool   `json:"hasLowercase"`
		HasNumber    bool   `json:"hasNumber"`
		HasSymbol    bool   `json:"hasSymbol"`
	} `json:"policy"`
}

func (uc *AuthUsecase) getZitadelPasswordPolicy(ctx context.Context) (PasswordPolicy, error) {
	url := fmt.Sprintf("%s/management/v1/policies/password/complexity", strings.TrimRight(uc.conf.Authority, "/"))
	var resp zitadelPasswordComplexityResponse
	if err := uc.callZitadelWithPATContext(ctx, http.MethodGet, url, nil, &resp, uc.managementPAT, "management"); err != nil {
		return PasswordPolicy{}, err
	}

	minLength := 8
	if resp.Policy.MinLength != "" {
		if _, err := fmt.Sscanf(resp.Policy.MinLength, "%d", &minLength); err != nil {
			minLength = 8
		}
	}

	return PasswordPolicy{
		MinLength:        minLength,
		RequireUppercase: resp.Policy.HasUppercase,
		RequireLowercase: resp.Policy.HasLowercase,
		RequireNumber:    resp.Policy.HasNumber,
		RequireSymbol:    resp.Policy.HasSymbol,
	}, nil
}

func (uc *AuthUsecase) CreateZitadelUser(ctx context.Context, username, email, displayName, password string) (string, error) {
	if username == "" {
		username = email
	}

	url := fmt.Sprintf("%s/v2/users/human", uc.conf.Authority)
	body := map[string]any{
		"username": username,
		"profile": map[string]string{
			"givenName":  displayName,
			"familyName": "User", // Default
		},
		"email": map[string]any{
			"email":           email,
			"isEmailVerified": true,
		},
		"password": map[string]any{
			"password":       password,
			"changeRequired": false,
		},
	}

	var resp zitadelCreateUserResponse
	if err := uc.callZitadelWithPAT(http.MethodPost, url, body, &resp, uc.managementPAT, "management"); err != nil {
		return "", err
	}
	return resp.UserID, nil
}

func (uc *AuthUsecase) callZitadelWithPAT(method, url string, body any, result any, pat, purpose string) error {
	return uc.callZitadelWithPATContext(context.Background(), method, url, body, result, pat, purpose)
}

func (uc *AuthUsecase) callZitadelWithPATContext(ctx context.Context, method, url string, body any, result any, pat, purpose string) error {
	pat = strings.TrimSpace(pat)
	if pat == "" || strings.Contains(pat, "${") {
		switch purpose {
		case "management":
			return fmt.Errorf("zitadel management PAT is not configured; set ZITADEL_MANAGEMENT_PAT or ZITADEL_PAT")
		default:
			return fmt.Errorf("zitadel login client PAT is not configured; set ZITADEL_LOGIN_CLIENT_PAT")
		}
	}

	var reader io.Reader
	if body != nil {
		jsonBody, _ := json.Marshal(body)
		reader = bytes.NewBuffer(jsonBody)
	}
	req, err := http.NewRequestWithContext(ctx, method, url, reader)
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", pat))

	client := &http.Client{Timeout: 15 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		respBody, _ := io.ReadAll(resp.Body)
		uc.log.Errorf("Zitadel API Error [%s]: %s", url, string(respBody))
		apiErr := parseZitadelAPIError(resp.StatusCode, resp.Status, respBody)
		return apiErr
	}

	if result == nil {
		return nil
	}
	return json.NewDecoder(resp.Body).Decode(result)
}

func parseZitadelAPIError(statusCode int, status string, body []byte) *ZitadelAPIError {
	apiErr := &ZitadelAPIError{StatusCode: statusCode, Status: status}
	var payload struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
		Details []struct {
			ID      string `json:"id"`
			Message string `json:"message"`
		} `json:"details"`
	}
	if err := json.Unmarshal(body, &payload); err == nil {
		apiErr.Code = fmt.Sprintf("%d", payload.Code)
		apiErr.Message = payload.Message
		for _, detail := range payload.Details {
			if detail.ID != "" {
				apiErr.Code = detail.ID
			}
			if detail.Message != "" {
				apiErr.Message = detail.Message
			}
		}
	}
	return apiErr
}
