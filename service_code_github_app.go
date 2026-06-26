package main

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/ca-risken/datasource-api/proto/code"
)

const (
	githubAppOAuthCallbackPath = "/api/v1/code/github-app/oauth/callback"
	githubAppStateTTL          = 10 * time.Minute
	githubAppOAuthAuthorizeURL = "https://github.com/login/oauth/authorize"
)

type githubAppOAuthState struct {
	ProjectID       uint32 `json:"project_id"`
	GithubSettingID uint32 `json:"github_setting_id"`
	UserID          uint32 `json:"user_id"`
	ReturnTo        string `json:"return_to"`
	Random          string `json:"random"`
	ExpiresAt       int64  `json:"expires_at"`
}

func (g *gatewayService) githubAppOAuthStartHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	u, err := getRequestUser(r)
	if err != nil || !isHumanAccess(u) {
		writeResponse(ctx, w, http.StatusUnauthorized, map[string]any{errorJSONKey: "Unauthenticated"})
		return
	}
	projectID, err := parseRequiredUint32(r, "project_id")
	if err != nil {
		writeResponse(ctx, w, http.StatusBadRequest, map[string]any{errorJSONKey: err.Error()})
		return
	}
	githubSettingID, err := parseRequiredUint32(r, "github_setting_id")
	if err != nil {
		writeResponse(ctx, w, http.StatusBadRequest, map[string]any{errorJSONKey: err.Error()})
		return
	}
	returnTo, err := normalizeGitHubAppReturnTo(r.URL.Query().Get("return_to"), projectID)
	if err != nil {
		writeResponse(ctx, w, http.StatusBadRequest, map[string]any{errorJSONKey: err.Error()})
		return
	}
	state, err := g.newGitHubAppOAuthState(projectID, githubSettingID, u.userID, returnTo, time.Now())
	if err != nil {
		appLogger.Errorf(ctx, "Failed to create github app oauth state: err=%+v", err)
		writeResponse(ctx, w, http.StatusServiceUnavailable, map[string]any{errorJSONKey: "GitHub App OAuth is not configured"})
		return
	}
	oauthURL, err := g.buildGitHubAppOAuthStartURL(state)
	if err != nil {
		appLogger.Errorf(ctx, "Failed to build github app oauth start url: err=%+v", err)
		writeResponse(ctx, w, http.StatusServiceUnavailable, map[string]any{errorJSONKey: "GitHub App OAuth start URL is not configured"})
		return
	}
	writeResponse(ctx, w, http.StatusOK, map[string]any{successJSONKey: map[string]string{"url": oauthURL}})
}

func (g *gatewayService) githubAppOAuthCallbackHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	state, err := g.verifyGitHubAppOAuthState(r.URL.Query().Get("state"), time.Now())
	if err != nil {
		appLogger.Warnf(ctx, "Invalid github app oauth state: err=%+v", err)
		writeResponse(ctx, w, http.StatusBadRequest, map[string]any{errorJSONKey: "Invalid GitHub App OAuth state"})
		return
	}
	u, err := getRequestUser(r)
	if err != nil || !isHumanAccess(u) {
		g.redirectGitHubAppOAuthResult(w, r, state.ReturnTo, "session_expired")
		return
	}
	if state.UserID != u.userID {
		appLogger.Warnf(ctx, "GitHub App OAuth state user mismatch: state_user_id=%d, request_user_id=%d", state.UserID, u.userID)
		g.redirectGitHubAppOAuthResult(w, r, state.ReturnTo, "unauthorized")
		return
	}
	if !g.isAuthorizedProject(ctx, u.userID, state.ProjectID, r.URL.Path) {
		g.redirectGitHubAppOAuthResult(w, r, state.ReturnTo, "unauthorized")
		return
	}
	oauthCode := r.URL.Query().Get("code")
	if strings.TrimSpace(oauthCode) == "" {
		g.redirectGitHubAppOAuthResult(w, r, state.ReturnTo, "failed")
		return
	}
	if _, err := g.codeClient.VerifyGitHubAppUser(ctx, &code.VerifyGitHubAppUserRequest{
		ProjectId:       state.ProjectID,
		GithubSettingId: state.GithubSettingID,
		Code:            oauthCode,
	}); err != nil {
		appLogger.Warnf(ctx, "Failed to verify github app oauth user: project_id=%d, github_setting_id=%d, err=%+v", state.ProjectID, state.GithubSettingID, err)
		g.redirectGitHubAppOAuthResult(w, r, state.ReturnTo, "failed")
		return
	}
	g.redirectGitHubAppOAuthResult(w, r, state.ReturnTo, "success")
}

func parseRequiredUint32(r *http.Request, name string) (uint32, error) {
	value := r.URL.Query().Get(name)
	if value == "" {
		return 0, fmt.Errorf("required %s", name)
	}
	parsed, err := strconv.ParseUint(value, 10, 32)
	if err != nil || parsed == 0 {
		return 0, fmt.Errorf("invalid %s", name)
	}
	return uint32(parsed), nil
}

func normalizeGitHubAppReturnTo(returnTo string, projectID uint32) (string, error) {
	if returnTo == "" {
		return fmt.Sprintf("/code/github?project_id=%d", projectID), nil
	}
	if err := validateGitHubAppReturnTo(returnTo); err != nil {
		return "", err
	}
	return returnTo, nil
}

func validateGitHubAppReturnTo(returnTo string) error {
	u, err := url.Parse(returnTo)
	if err != nil {
		return err
	}
	if u.IsAbs() || u.Host != "" || !strings.HasPrefix(returnTo, "/") || strings.HasPrefix(returnTo, "//") || strings.ContainsRune(returnTo, '\\') {
		return errors.New("return_to must be a relative path")
	}
	return nil
}

func (g *gatewayService) newGitHubAppOAuthState(projectID, githubSettingID, userID uint32, returnTo string, now time.Time) (string, error) {
	randomBytes := make([]byte, 16)
	if _, err := rand.Read(randomBytes); err != nil {
		return "", err
	}
	state := &githubAppOAuthState{
		ProjectID:       projectID,
		GithubSettingID: githubSettingID,
		UserID:          userID,
		ReturnTo:        returnTo,
		Random:          base64.RawURLEncoding.EncodeToString(randomBytes),
		ExpiresAt:       now.Add(githubAppStateTTL).Unix(),
	}
	return g.signGitHubAppOAuthState(state)
}

func (g *gatewayService) signGitHubAppOAuthState(state *githubAppOAuthState) (string, error) {
	if g.githubAppStateSecret == "" {
		return "", errors.New("github app state secret is required")
	}
	payload, err := json.Marshal(state)
	if err != nil {
		return "", err
	}
	encodedPayload := base64.RawURLEncoding.EncodeToString(payload)
	signature := signGitHubAppState(encodedPayload, g.githubAppStateSecret)
	return encodedPayload + "." + signature, nil
}

func (g *gatewayService) verifyGitHubAppOAuthState(rawState string, now time.Time) (*githubAppOAuthState, error) {
	if g.githubAppStateSecret == "" {
		return nil, errors.New("github app state secret is required")
	}
	parts := strings.Split(rawState, ".")
	if len(parts) != 2 {
		return nil, errors.New("invalid state format")
	}
	expectedSignature := signGitHubAppState(parts[0], g.githubAppStateSecret)
	if !hmac.Equal([]byte(expectedSignature), []byte(parts[1])) {
		return nil, errors.New("invalid state signature")
	}
	payload, err := base64.RawURLEncoding.DecodeString(parts[0])
	if err != nil {
		return nil, err
	}
	state := &githubAppOAuthState{}
	if err := json.Unmarshal(payload, state); err != nil {
		return nil, err
	}
	if state.ProjectID == 0 || state.GithubSettingID == 0 || state.UserID == 0 || state.ReturnTo == "" || state.Random == "" {
		return nil, errors.New("invalid state payload")
	}
	if err := validateGitHubAppReturnTo(state.ReturnTo); err != nil {
		return nil, err
	}
	if now.Unix() > state.ExpiresAt {
		return nil, errors.New("expired state")
	}
	return state, nil
}

func signGitHubAppState(payload, secret string) string {
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write([]byte(payload))
	return base64.RawURLEncoding.EncodeToString(mac.Sum(nil))
}

func (g *gatewayService) buildGitHubAppOAuthStartURL(state string) (string, error) {
	clientID := strings.TrimSpace(g.githubAppClientID)
	if clientID == "" {
		return "", errors.New("github app oauth client id is required")
	}
	u, err := url.Parse(githubAppOAuthAuthorizeURL)
	if err != nil {
		return "", err
	}
	q := u.Query()
	q.Set("client_id", clientID)
	q.Set("state", state)
	u.RawQuery = q.Encode()
	return u.String(), nil
}

func (g *gatewayService) redirectGitHubAppOAuthResult(w http.ResponseWriter, r *http.Request, returnTo, result string) {
	if err := validateGitHubAppReturnTo(returnTo); err != nil {
		http.Redirect(w, r, "/code/github?github_app_oauth=failed", http.StatusFound)
		return
	}
	u, err := url.Parse(returnTo)
	if err != nil {
		http.Redirect(w, r, "/code/github?github_app_oauth=failed", http.StatusFound)
		return
	}
	q := u.Query()
	q.Set("github_app_oauth", result)
	u.RawQuery = q.Encode()
	http.Redirect(w, r, u.String(), http.StatusFound)
}
