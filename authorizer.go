package main

import (
	"bytes"
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/CyberAgent/mimosa-core/proto/iam"
	"github.com/vikyd/zero"
)

type key int

const (
	userKey key = iota
)

type requestUser struct {
	sub    string
	userID uint32
}

// signinHandler: OIDC proxy backend signin process.
func signinHandler(w http.ResponseWriter, r *http.Request) {
	if user, ok := r.Context().Value(userKey).(*requestUser); !ok {
		appLogger.Infof("Unauthenticated: Invalid requestUser type")
		http.Error(w, "Unauthenticated", http.StatusUnauthorized)
		return
	} else if user == nil || zero.IsZeroVal(user.userID) {
		appLogger.Infof("Unauthenticated: No mimosa-user")
		http.Error(w, "Unauthenticated", http.StatusUnauthorized)
		return
	}
	signinUser := r.Context().Value(userKey).(*requestUser)

	token := make([]byte, 24)
	rand.Read(token)
	http.SetCookie(w, &http.Cookie{
		Name:   "XSRF-TOKEN",
		Value:  base64.RawURLEncoding.EncodeToString(token),
		Path:   "/",
		Secure: r.Header.Get("X-Forwarded-Proto") == "https",
	})
	writeResponse(w, http.StatusOK, map[string]interface{}{
		"user_id": signinUser.userID,
	})
}

func (g *gatewayService) authn(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		sub := r.Header.Get(g.uidHeader)
		if sub == "" {
			next.ServeHTTP(w, r)
			return
		}
		appLogger.Debugf("sub: %s", sub)
		resp, err := g.iamClient.GetUser(r.Context(), &iam.GetUserRequest{Sub: sub})
		if err != nil {
			appLogger.Warnf("Failed to GetUser request, err=%+v", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
		if resp != nil && resp.User != nil {
			next.ServeHTTP(w, r.WithContext(
				context.WithValue(r.Context(), userKey, &requestUser{sub: sub, userID: resp.User.UserId})))
			return
		}
		// Try user auto provisioning
		oidcData := r.Header.Get(g.oidcDataHeader) // r.Header.Get("X-Amzn-Oidc-Data")
		userName, err := g.getUserName(oidcData)
		if err != nil {
			appLogger.Warnf("Failed to get username from oidc data, err=%+v", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
		putResp, err := g.iamClient.PutUser(r.Context(), &iam.PutUserRequest{
			User: &iam.UserForUpsert{
				Sub:       sub,
				Name:      userName,
				Activated: true,
			},
		})
		if err != nil {
			appLogger.Warnf("Failed to PutUser request, err=%+v", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
		if putResp != nil && putResp.User != nil {
			next.ServeHTTP(w, r.WithContext(
				context.WithValue(r.Context(), userKey, &requestUser{sub: sub, userID: putResp.User.UserId})))
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

func (g *gatewayService) authzWithProject(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		buf, err := ioutil.ReadAll(r.Body)
		r.Body = ioutil.NopCloser(bytes.NewBuffer(buf))
		if err != nil {
			http.Error(w, "Could not read body", http.StatusInternalServerError)
			return
		}
		u, ok := r.Context().Value(userKey).(*requestUser)
		if !ok {
			appLogger.Infof("Unauthenticated: Invalid requestUser type.")
			http.Error(w, "Unauthenticated", http.StatusUnauthorized)
			return
		}
		if !validCSRFToken(r) {
			http.Error(w, "Invalid token", http.StatusForbidden)
			return
		}
		if !g.authzProject(u, r) {
			http.Error(w, "Unauthorized xxx", http.StatusForbidden)
			return
		}
		r.Body = ioutil.NopCloser(bytes.NewBuffer(buf)) // 後続のハンドラでもリクエストボディを読み取れるように上書きしとく
		next.ServeHTTP(w, r)
	}
	return http.HandlerFunc(fn)
}

func (g *gatewayService) authzOnlyAdmin(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		u, ok := r.Context().Value(userKey).(*requestUser)
		if !ok {
			appLogger.Infof("Unauthenticated: Invalid requestUser type.")
			http.Error(w, "Unauthenticated", http.StatusUnauthorized)
			return
		}
		if !validCSRFToken(r) {
			http.Error(w, "Invalid token", http.StatusForbidden)
			return
		}
		if !g.authzAdmin(u.userID, r) {
			http.Error(w, "Unauthorized admin API", http.StatusForbidden)
			return
		}
		next.ServeHTTP(w, r)
	}
	return http.HandlerFunc(fn)
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
	resp, err := g.iamClient.IsAuthorized(r.Context(), req)
	if err != nil {
		appLogger.Errorf("Failed to IsAuthorized requuest, request=%+v, err=%+v", req, err)
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

func (g *gatewayService) authzAdmin(userID uint32, r *http.Request) bool {
	if zero.IsZeroVal(userID) {
		return false
	}
	req := &iam.IsAdminRequest{UserId: userID}
	resp, err := g.iamClient.IsAdmin(r.Context(), req)
	if err != nil {
		appLogger.Errorf("Failed to IsAdmin requuest, request=%+v, err=%+v", req, err)
		return false
	}
	return resp.Ok
}
