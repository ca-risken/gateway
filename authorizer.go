package main

import (
	"bytes"
	"context"
	"errors"
	"io"
	"net/http"
	"strings"

	"github.com/ca-risken/core/proto/iam"
	"github.com/ca-risken/core/proto/organization_iam"
	"github.com/vikyd/zero"
)

type key int

const (
	userKey key = iota
)

type requestUser struct {
	// human access
	sub    string
	userID uint32
	name   string

	// program access
	accessTokenID        uint32
	accessTokenProjectID uint32
}

func getRequestUser(r *http.Request) (*requestUser, error) {
	if u, ok := r.Context().Value(userKey).(*requestUser); !ok || u == nil || (zero.IsZeroVal(u.userID) && zero.IsZeroVal(u.accessTokenID)) {
		return nil, errors.New("user not found")
	}
	return r.Context().Value(userKey).(*requestUser), nil
}

func getRequestUserSub(r *http.Request) (*requestUser, error) {
	if u, ok := r.Context().Value(userKey).(*requestUser); !ok || u == nil || zero.IsZeroVal(u.sub) {
		return nil, errors.New("user not found")
	}
	return r.Context().Value(userKey).(*requestUser), nil
}

func (g *gatewayService) UpdateUserFromIdp(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		sub := r.Header.Get(g.uidHeader)
		if sub == "" {
			next.ServeHTTP(w, r)
			return
		}
		oidcData := r.Header.Get(g.oidcDataHeader) // r.Header.Get("X-Amzn-Oidc-Data")
		claims, err := g.claimsClient.getClaims(ctx, oidcData)
		appLogger.Debugf(ctx, "claims: %+v", claims)
		if err != nil {
			appLogger.Warnf(ctx, "Failed to validate id token, err=%+v", err)
			http.Error(w, "Failed to validate id token", http.StatusForbidden)
			return
		}
		userIdpKey := g.claimsClient.getUserIdpKey(claims)
		if userIdpKey == "" {
			appLogger.Warnf(ctx, "UserIdpKey is not found in token, err=%+v", err)
			http.Error(w, "UserIdpKey is not found in token", http.StatusForbidden)
			return
		}
		userName := g.claimsClient.getUserName(claims)
		if userName == "" {
			userName = userIdpKey
		}
		putUserRequest := &iam.PutUserRequest{
			User: &iam.UserForUpsert{
				Sub:        sub,
				Name:       userName,
				UserIdpKey: userIdpKey,
				Activated:  true,
			},
		}
		// 既存のユーザーであれば、手動でNameを変更している可能性があるので、contextのユーザー名を使用する
		u, err := getRequestUser(r)
		if err == nil && u != nil && u.name != "" {
			putUserRequest.User.Name = u.name
		}
		putResp, err := g.iamClient.PutUser(ctx, putUserRequest)
		if err != nil {
			appLogger.Warnf(ctx, "Failed to PutUser request, err=%+v", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
		if putResp != nil && putResp.User != nil {
			next.ServeHTTP(w, r.WithContext(
				context.WithValue(ctx, userKey, &requestUser{sub: sub, userID: putResp.User.UserId})))
			return
		}
		next.ServeHTTP(w, r)
	}
	return http.HandlerFunc(fn)
}

// Authentication for human access
func (g *gatewayService) authn(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		sub := r.Header.Get(g.uidHeader)
		if sub == "" {
			next.ServeHTTP(w, r)
			return
		}
		appLogger.Debugf(ctx, "sub: %s", sub)
		resp, err := g.iamClient.GetUser(ctx, &iam.GetUserRequest{Sub: sub})
		if err != nil {
			appLogger.Warnf(ctx, "Failed to GetUser request, err=%+v", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
		if resp != nil && resp.User != nil {
			next.ServeHTTP(w, r.WithContext(
				context.WithValue(ctx, userKey, &requestUser{sub: sub, userID: resp.User.UserId, name: resp.User.Name})))
			return
		}
		next.ServeHTTP(w, r)
	}
	return http.HandlerFunc(fn)
}

// Authentication for programable API access
func (g *gatewayService) authnToken(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		bearer := r.Header.Get("Authorization")
		tokenBody := ""
		if len(bearer) > 7 && strings.ToUpper(bearer[0:7]) == "BEARER " {
			tokenBody = strings.TrimSpace(bearer[7:])
		}
		if tokenBody == "" {
			next.ServeHTTP(w, r)
			return
		}
		projectID, accessTokenID, plainTextToken, err := decodeAccessToken(ctx, tokenBody)
		if err != nil {
			// TODO アクセストークンが不要な後続処理があるかを確認、不要な場合はすぐに403などを返したい
			next.ServeHTTP(w, r)
			return
		}
		resp, err := g.iamClient.AuthenticateAccessToken(ctx, &iam.AuthenticateAccessTokenRequest{
			ProjectId:      projectID,
			AccessTokenId:  accessTokenID,
			PlainTextToken: plainTextToken,
		})
		if err != nil {
			// TODO 認証でエラーになった後に継続する後続の処理があるか確認、できる限りすぐに403などを返したい
			appLogger.Errorf(ctx, "Failed to AuthenticateAccessToken API, err=%+v", err)
			next.ServeHTTP(w, r)
			return
		}
		if resp.AccessToken == nil || resp.AccessToken.AccessTokenId == 0 {
			appLogger.Error(ctx, "Failed to get AccessTokenId")
			next.ServeHTTP(w, r)
			return
		}
		next.ServeHTTP(w, r.WithContext(
			context.WithValue(ctx, userKey, &requestUser{
				accessTokenID:        accessTokenID,
				accessTokenProjectID: projectID,
			})))
	}
	return http.HandlerFunc(fn)
}

func (g *gatewayService) authzWithProject(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		buf, err := io.ReadAll(r.Body)
		if err != nil {
			appLogger.Errorf(ctx, "Failed to read body, err=%+v", err)
			http.Error(w, "Could not read body", http.StatusInternalServerError)
			return
		}
		r.Body = io.NopCloser(bytes.NewBuffer(buf))

		u, err := getRequestUser(r)
		if err != nil {
			appLogger.Infof(ctx, "Unauthenticated: %+v", err)
			http.Error(w, "Unauthenticated", http.StatusUnauthorized)
			return
		}

		if isHumanAccess(u) {
			// Human Access
			if !g.authzProject(u, r) {
				http.Error(w, "Unauthorized the project resource for human access", http.StatusForbidden)
				return
			}
		} else {
			// Program Access
			if !g.authzProjectForToken(u, r) {
				http.Error(w, "Unauthorized the project resource for token access", http.StatusForbidden)
				return
			}
		}
		r.Body = io.NopCloser(bytes.NewBuffer(buf)) // 後続のハンドラでもリクエストボディを読み取れるように上書きしとく
		next.ServeHTTP(w, r)
	}
	return http.HandlerFunc(fn)
}

func (g *gatewayService) authzOnlyAdmin(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		buf, err := io.ReadAll(r.Body)
		if err != nil {
			appLogger.Errorf(ctx, "Failed to read body, err=%+v", err)
			http.Error(w, "Could not read body", http.StatusInternalServerError)
			return
		}
		r.Body = io.NopCloser(bytes.NewBuffer(buf))

		u, err := getRequestUser(r)
		if err != nil {
			appLogger.Infof(ctx, "Unauthenticated: %+v", err)
			http.Error(w, "Unauthenticated", http.StatusUnauthorized)
			return
		}
		if !g.authzAdmin(u, r) {
			http.Error(w, "Unauthorized admin API", http.StatusForbidden)
			return
		}
		r.Body = io.NopCloser(bytes.NewBuffer(buf)) // 後続のハンドラでもリクエストボディを読み取れるように上書きしとく
		next.ServeHTTP(w, r)
	}
	return http.HandlerFunc(fn)
}

func (g *gatewayService) authzWithOrganization(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		buf, err := io.ReadAll(r.Body)
		if err != nil {
			appLogger.Errorf(ctx, "Failed to read body, err=%+v", err)
			http.Error(w, "Could not read body", http.StatusInternalServerError)
			return
		}
		r.Body = io.NopCloser(bytes.NewBuffer(buf))

		u, err := getRequestUser(r)
		if err != nil {
			appLogger.Infof(ctx, "Unauthenticated: %+v", err)
			http.Error(w, "Unauthenticated", http.StatusUnauthorized)
			return
		}

		if !isHumanAccess(u) {
			http.Error(w, "Organization API does not support program access", http.StatusForbidden)
			return
		}

		if !g.authzOrganization(u, r) {
			http.Error(w, "Unauthorized the organization resource", http.StatusForbidden)
			return
		}

		r.Body = io.NopCloser(bytes.NewBuffer(buf)) // 後続のハンドラでもリクエストボディを読み取れるように上書きしとく
		next.ServeHTTP(w, r)
	}
	return http.HandlerFunc(fn)
}

func (g *gatewayService) verifyCSRF(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		u, _ := getRequestUser(r)
		if isHumanAccess(u) &&
			shouldVerifyCSRFTokenURI(r.URL.Path) &&
			!validCSRFToken(r) {
			appLogger.Debugf(ctx, "Invalid CSRF token: request_user=%+v, uri=%s", u, r.RequestURI)
			writeResponse(ctx, w, http.StatusForbidden, map[string]interface{}{errorJSONKey: "Invalid token"})
			return
		}
		next.ServeHTTP(w, r)
	}
	return http.HandlerFunc(fn)
}

func isHumanAccess(u *requestUser) bool {
	if u == nil || zero.IsZeroVal(u.userID) {
		return false
	}
	if !zero.IsZeroVal(u.accessTokenID) {
		return false
	}
	return true
}

var ignoreURI4CSRF = []string{
	"/healthz",
	"/api/v1/signin",
}

func shouldVerifyCSRFTokenURI(uri string) bool {
	trimedURI := strings.TrimSuffix(uri, "/")
	for _, ignoreURI := range ignoreURI4CSRF {
		if trimedURI == ignoreURI {
			return false
		}
	}
	return true
}

func validCSRFToken(r *http.Request) bool {
	headerToken := r.Header.Get("X-XSRF-TOKEN")
	if headerToken == "" {
		return false
	}
	cookieToken, err := r.Cookie(XSRF_TOKEN)
	if err != nil || cookieToken.Value == "" {
		return false
	}
	return cookieToken.Value == headerToken
}

type requestProject struct {
	ProjectID uint32 `json:"project_id"`
}

type requestOrganization struct {
	OrganizationID uint32 `json:"organization_id"`
}

func (g *gatewayService) authzProject(u *requestUser, r *http.Request) bool {
	ctx := r.Context()
	if zero.IsZeroVal(u.userID) {
		return false
	}
	p := &requestProject{}
	err := bind(p, r)
	if err != nil {
		appLogger.Warnf(ctx, "Failed to bind request, err=%+v", err)
	}
	if zero.IsZeroVal(p.ProjectID) {
		return false
	}
	req := &iam.IsAuthorizedRequest{
		UserId:       u.userID,
		ProjectId:    p.ProjectID,
		ActionName:   getActionNameFromURI(r.URL.Path),
		ResourceName: getServiceNameFromURI(r.URL.Path) + "/resource_any",
	}
	resp, err := g.iamClient.IsAuthorized(ctx, req)
	if err != nil {
		appLogger.Errorf(ctx, "Failed to IsAuthorized request, request=%+v, err=%+v", req, err)
		return false
	}
	return resp.Ok
}

func (g *gatewayService) authzProjectForToken(u *requestUser, r *http.Request) bool {
	ctx := r.Context()
	if u.accessTokenID == 0 || u.accessTokenProjectID == 0 {
		return false
	}
	p := &requestProject{}
	err := bind(p, r)
	if err != nil {
		appLogger.Warnf(ctx, "Failed to bind request, err=%+v", err)
	}
	if p.ProjectID != 0 && p.ProjectID != u.accessTokenProjectID {
		return false
	}
	req := &iam.IsAuthorizedTokenRequest{
		AccessTokenId: u.accessTokenID,
		ProjectId:     u.accessTokenProjectID,
		ActionName:    getActionNameFromURI(r.URL.Path),
		ResourceName:  getServiceNameFromURI(r.URL.Path) + "/resource_any",
	}
	resp, err := g.iamClient.IsAuthorizedToken(ctx, req)
	if err != nil {
		appLogger.Errorf(ctx, "Failed to IsAuthorizedToken request, request=%+v, err=%+v", req, err)
		return false
	}
	return resp.Ok
}

const prefixURI = "/api/v1/"

// getActionNameFromURI: `/api/v1/service/path1/path2/...` will return `service/path1`
func getActionNameFromURI(uri string) string {
	if !strings.HasPrefix(uri, prefixURI) {
		return ""
	}
	paths := strings.Split(strings.Replace(uri, prefixURI, "", 1), "/")
	if len(paths) < 2 {
		return ""
	}
	return paths[0] + "/" + paths[1]
}

// getServiceNameFromURI: `/service/path1/path2/...` will return `service`
func getServiceNameFromURI(uri string) string {
	if !strings.HasPrefix(uri, prefixURI) {
		return ""
	}
	paths := strings.Split(strings.Replace(uri, prefixURI, "", 1), "/")
	if len(paths) < 1 {
		return ""
	}
	return paths[0]
}

func (g *gatewayService) authzAdmin(u *requestUser, r *http.Request) bool {
	ctx := r.Context()
	if zero.IsZeroVal(u.userID) {
		return false
	}
	req := &iam.IsAuthorizedAdminRequest{
		UserId:       u.userID,
		ActionName:   getActionNameFromURI(r.URL.Path),
		ResourceName: getServiceNameFromURI(r.URL.Path) + "/resource_any",
	}
	resp, err := g.iamClient.IsAuthorizedAdmin(ctx, req)
	if err != nil {
		appLogger.Errorf(ctx, "Failed to IsAuthorizedAdmin request, request=%+v, err=%+v", req, err)
		return false
	}
	if !resp.Ok {
		appLogger.Debugf(ctx, "user=%d is not Admin, request=%+v", u.userID, req)
		return false
	}
	return resp.Ok
}

func (g *gatewayService) authzOrganization(u *requestUser, r *http.Request) bool {
	ctx := r.Context()
	if u.userID == 0 {
		return false
	}
	o := &requestOrganization{}
	err := bind(o, r)
	if err != nil {
		appLogger.Warnf(ctx, "Failed to bind request, err=%+v", err)
	}
	if o.OrganizationID == 0 {
		return false
	}
	req := &organization_iam.IsAuthorizedOrganizationRequest{
		UserId:         u.userID,
		OrganizationId: o.OrganizationID,
		ActionName:     getActionNameFromURI(r.URL.Path),
	}
	resp, err := g.organization_iamClient.IsAuthorizedOrganization(ctx, req)
	if err != nil {
		appLogger.Errorf(ctx, "Failed to IsAuthorizedOrganization request, request=%+v, err=%+v", req, err)
		return false
	}
	return resp.Ok
}
