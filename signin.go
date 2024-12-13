package main

import (
	"crypto/rand"
	"encoding/base64"
	"net/http"
	"time"

	"github.com/ca-risken/common/pkg/logging"
)

const (
	XSRF_TOKEN = "XSRF-TOKEN"
)

// signinHandler: OIDC proxy backend signin process.
func (g *gatewayService) signinHandler(w http.ResponseWriter, r *http.Request) {
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
		Name:   XSRF_TOKEN,
		Value:  base64.RawURLEncoding.EncodeToString(token),
		Path:   "/",
		MaxAge: g.sessionTimeoutSec,
		Secure: r.Header.Get("X-Forwarded-Proto") == "https",
	})
	appLogger.WithItems(ctx, logging.InfoLevel, map[string]interface{}{"user_id": signinUser.userID}, "Signin")

	resp := map[string]interface{}{}
	if signinUser.userID != 0 {
		resp["user_id"] = signinUser.userID
	} else if signinUser.accessTokenID != 0 && signinUser.accessTokenProjectID != 0 {
		resp["access_token_id"] = signinUser.accessTokenID
		resp["project_id"] = signinUser.accessTokenProjectID
	}
	writeResponse(ctx, w, http.StatusOK, resp)
}

// signoutHandler
func (g *gatewayService) signoutHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	signinUser, err := getRequestUser(r)
	if err != nil {
		appLogger.Infof(ctx, "Unauthenticated: %+v", err)
		http.Error(w, "Unauthenticated", http.StatusUnauthorized)
		return
	}

	// Remove cookies
	for _, cookie := range g.sessionCookieName {
		http.SetCookie(w, &http.Cookie{
			Name:     cookie,
			Value:    "",
			Path:     "/",
			MaxAge:   -1,
			Expires:  time.Unix(0, 0),
			HttpOnly: true,
			Secure:   r.Header.Get("X-Forwarded-Proto") == "https",
		})
	}
	http.SetCookie(w, &http.Cookie{
		Name:    XSRF_TOKEN,
		Value:   "",
		Path:    "/",
		MaxAge:  -1,
		Expires: time.Unix(0, 0),
		Secure:  r.Header.Get("X-Forwarded-Proto") == "https",
	})
	appLogger.WithItems(ctx, logging.InfoLevel,
		map[string]interface{}{"user_id": signinUser.userID}, "Signout")
	writeResponse(ctx, w, http.StatusOK, nil)
}
