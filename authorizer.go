package main

import (
	"bytes"
	"context"
	"crypto/ecdsa"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/ca-risken/core/proto/iam"
	"github.com/golang-jwt/jwt/v4"
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

	// program access
	accessTokenID uint32
}

func getRequestUser(r *http.Request) (*requestUser, error) {
	if u, ok := r.Context().Value(userKey).(*requestUser); !ok || u == nil || (zero.IsZeroVal(u.userID) && zero.IsZeroVal(u.accessTokenID)) {
		return nil, errors.New("user not found")
	}
	appLogger.Infof(context.Background(), "requestUser: %+v", r.Context().Value(userKey).(*requestUser))
	return r.Context().Value(userKey).(*requestUser), nil
}

func getRequestUserSub(r *http.Request) (*requestUser, error) {
	if u, ok := r.Context().Value(userKey).(*requestUser); !ok || u == nil || zero.IsZeroVal(u.sub) {
		return nil, errors.New("user not found")
	}
	return r.Context().Value(userKey).(*requestUser), nil
}

// signinHandler: OIDC proxy backend signin process.
func signinHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	signinUser, err := getRequestUser(r)
	if err != nil {
		appLogger.Infof(ctx, "Unauthenticated: %+v", err)
		http.Error(w, "Unauthenticated", http.StatusUnauthorized)
		return
	}
	token := make([]byte, 24)
	_, _ = rand.Read(token)
	http.SetCookie(w, &http.Cookie{
		Name:   "XSRF-TOKEN",
		Value:  base64.RawURLEncoding.EncodeToString(token),
		Path:   "/",
		Secure: r.Header.Get("X-Forwarded-Proto") == "https",
	})
	writeResponse(ctx, w, http.StatusOK, map[string]interface{}{
		"user_id": signinUser.userID,
	})
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
		appLogger.Debugf(ctx, "oidcData: %+v", r.Header.Get(g.oidcDataHeader))
		resp, err := g.iamClient.GetUser(ctx, &iam.GetUserRequest{Sub: sub})
		if err != nil {
			appLogger.Warnf(ctx, "Failed to GetUser request, err=%+v", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
		if resp != nil && resp.User != nil {
			next.ServeHTTP(w, r.WithContext(
				context.WithValue(ctx, userKey, &requestUser{sub: sub, userID: resp.User.UserId})))
			return
		}
		// Try AUTO PROVISIONING
		oidcData := r.Header.Get(g.oidcDataHeader) // r.Header.Get("X-Amzn-Oidc-Data")
		if g.VerifyIDToken {
			err = g.verifyTokenForALB(oidcData)
			if err != nil {
				appLogger.Warnf(ctx, "Failed to validate id token, err=%+v", err)
				http.Error(w, "Internal server error", http.StatusForbidden)
				return
			}
		}
		userName, err := g.getUserName(oidcData)
		if err != nil || zero.IsZeroVal(userName) {
			appLogger.Warnf(ctx, "Failed to get username from oidc data, err=%+v", err)
			next.ServeHTTP(w, r.WithContext(
				context.WithValue(ctx, userKey, &requestUser{sub: sub})))
			return
		}
		putResp, err := g.iamClient.PutUser(ctx, &iam.PutUserRequest{
			User: &iam.UserForUpsert{
				Sub:       sub,
				Name:      userName,
				Activated: true,
			},
		})
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

type userCalim struct {
	Sub      string `json:"sub"`
	Username string `json:"username"`
}

func (g *gatewayService) getUserName(jwt string) (string, error) {
	parts := strings.Split(jwt, ".")
	if len(parts) != 3 {
		return "", errors.New("Invalid JWT string pattern")
	}
	// Decode JWT
	claimBytes, err := base64.StdEncoding.DecodeString(parts[1])
	if err != nil {
		return "", err
	}
	var user userCalim
	if err := json.NewDecoder(bytes.NewBuffer(claimBytes)).Decode(&user); err != nil {
		return "", err
	}
	username := user.Username
	for _, idp := range g.idpProviderName {
		if strings.HasPrefix(strings.ToLower(username), strings.ToLower(idp)+"_") {
			username = strings.Replace(username, idp+"_", "", 1)
			break
		}
	}
	return username, nil
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
			context.WithValue(ctx, userKey, &requestUser{accessTokenID: accessTokenID})))
	}
	return http.HandlerFunc(fn)
}

func (g *gatewayService) authzWithProject(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		buf, err := ioutil.ReadAll(r.Body)
		if err != nil {
			appLogger.Errorf(ctx, "Failed to read body, err=%+v", err)
			http.Error(w, "Could not read body", http.StatusInternalServerError)
			return
		}
		r.Body = ioutil.NopCloser(bytes.NewBuffer(buf))

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
		r.Body = ioutil.NopCloser(bytes.NewBuffer(buf)) // 後続のハンドラでもリクエストボディを読み取れるように上書きしとく
		next.ServeHTTP(w, r)
	}
	return http.HandlerFunc(fn)
}

func (g *gatewayService) authzOnlyAdmin(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		buf, err := ioutil.ReadAll(r.Body)
		if err != nil {
			appLogger.Errorf(ctx, "Failed to read body, err=%+v", err)
			http.Error(w, "Could not read body", http.StatusInternalServerError)
			return
		}
		r.Body = ioutil.NopCloser(bytes.NewBuffer(buf))

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
		r.Body = ioutil.NopCloser(bytes.NewBuffer(buf)) // 後続のハンドラでもリクエストボディを読み取れるように上書きしとく
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

const (
	PublicKeyURL = "https://public-keys.auth.elb.%s.amazonaws.com/%s"
)

func (g *gatewayService) verifyTokenForALB(tokenString string) error {
	jwt.DecodePaddingAllowed = true
	_, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodECDSA); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		kid, ok := token.Header["kid"].(string)
		if !ok {
			return nil, errors.New("kid is not found in jwt header")
		}
		keyURL := fmt.Sprintf(PublicKeyURL, g.Region, kid)
		key, err := g.fetchALBPublicKey(keyURL)
		if err != nil {
			return nil, fmt.Errorf("failed to get public key, err: %w", err)
		}
		return key, nil
	})
	if err != nil {
		return err
	}
	return nil
}

func (g *gatewayService) fetchALBPublicKey(keyURL string) (*ecdsa.PublicKey, error) {
	ctx := context.Background()
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, keyURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to new GET request for %s, err: %w", keyURL, err)
	}
	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to get public key from %s, err: %w", keyURL, err)
	}
	defer resp.Body.Close()
	pem, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to get public key from %s, err: %w", keyURL, err)
	}
	publicKey, err := jwt.ParseECPublicKeyFromPEM(pem)
	if err != nil {
		return nil, fmt.Errorf("failed to parse public key from %s, err: %w", keyURL, err)
	}

	return publicKey, nil
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
	cookieToken, err := r.Cookie("XSRF-TOKEN")
	if err != nil || cookieToken.Value == "" {
		return false
	}
	return cookieToken.Value == headerToken
}

type requestProject struct {
	ProjectID uint32 `json:"project_id"`
}

func (g *gatewayService) authzProject(u *requestUser, r *http.Request) bool {
	ctx := r.Context()
	if zero.IsZeroVal(u.userID) {
		return false
	}
	p := &requestProject{}
	bind(p, r)
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
	if zero.IsZeroVal(u.accessTokenID) {
		return false
	}
	p := &requestProject{}
	bind(p, r)
	if zero.IsZeroVal(p.ProjectID) {
		return false
	}
	req := &iam.IsAuthorizedTokenRequest{
		AccessTokenId: u.accessTokenID,
		ProjectId:     p.ProjectID,
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
