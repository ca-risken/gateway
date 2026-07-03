package main

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
	"time"

	"github.com/ca-risken/core/proto/iam"
	iammocks "github.com/ca-risken/core/proto/iam/mocks"
	"github.com/ca-risken/datasource-api/proto/code"
	codemocks "github.com/ca-risken/datasource-api/proto/code/mocks"
	"github.com/stretchr/testify/mock"
)

const testGitHubAppStateSecret = "12345678901234567890123456789012"

func newGitHubAppOAuthCallbackRequest(t *testing.T, svc *gatewayService, userID uint32) *http.Request {
	t.Helper()
	rawState, err := svc.newGitHubAppOAuthState(1001, 10, 20, "/code/github?project_id=1001", time.Now())
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	req := httptest.NewRequest(http.MethodGet, "/api/v1/code/github-app/oauth/callback?state="+url.QueryEscape(rawState)+"&code=oauth-code", nil)
	if userID != 0 {
		req = req.WithContext(context.WithValue(req.Context(), userKey, &requestUser{userID: userID}))
	}
	return req
}

func TestGitHubAppOAuthState(t *testing.T) {
	svc := &gatewayService{githubAppStateSecret: testGitHubAppStateSecret}
	now := time.Unix(1700000000, 0)

	rawState, err := svc.newGitHubAppOAuthState(1001, 10, 20, "/code/github?project_id=1001", now)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	state, err := svc.verifyGitHubAppOAuthState(rawState, now.Add(time.Minute))
	if err != nil {
		t.Fatalf("Unexpected verify error: %v", err)
	}
	if state.ProjectID != 1001 || state.GithubSettingID != 10 || state.UserID != 20 || state.ReturnTo != "/code/github?project_id=1001" {
		t.Fatalf("Unexpected state: %+v", state)
	}
}

func TestNewGitHubAppOAuthStateGeneratesUniqueToken(t *testing.T) {
	svc := &gatewayService{githubAppStateSecret: testGitHubAppStateSecret}
	now := time.Unix(1700000000, 0)

	state1, err := svc.newGitHubAppOAuthState(1001, 10, 20, "/code/github?project_id=1001", now)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	state2, err := svc.newGitHubAppOAuthState(1001, 10, 20, "/code/github?project_id=1001", now)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	if state1 == state2 {
		t.Fatal("Expected unique state tokens")
	}
}

func TestVerifyGitHubAppOAuthStateRejectsInvalidState(t *testing.T) {
	svc := &gatewayService{githubAppStateSecret: testGitHubAppStateSecret}
	now := time.Unix(1700000000, 0)
	rawState, err := svc.newGitHubAppOAuthState(1001, 10, 20, "/code/github?project_id=1001", now)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	invalidReturnToState, err := svc.signGitHubAppOAuthState(&githubAppOAuthState{
		ProjectID:       1001,
		GithubSettingID: 10,
		UserID:          20,
		ReturnTo:        `/\attacker.example/path`,
		Random:          "test-random",
		ExpiresAt:       now.Add(githubAppStateTTL).Unix(),
	})
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	cases := []struct {
		name  string
		state string
		now   time.Time
	}{
		{name: "invalid format", state: "invalid", now: now},
		{name: "invalid signature", state: rawState + "x", now: now},
		{name: "invalid return_to", state: invalidReturnToState, now: now},
		{name: "expired", state: rawState, now: now.Add(githubAppStateTTL + time.Second)},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if _, err := svc.verifyGitHubAppOAuthState(c.state, c.now); err == nil {
				t.Fatal("Expected error but got none")
			}
		})
	}
}

func TestBuildGitHubAppOAuthStartURL(t *testing.T) {
	svc := &gatewayService{
		githubAppClientID:    "test-github-app-client-id",
		githubAppRedirectURL: "https://risken.example/api/v1/code/github-app/oauth/callback",
	}
	got, err := svc.buildGitHubAppOAuthStartURL("state-value")
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	if got != "https://github.com/login/oauth/authorize?client_id=test-github-app-client-id&redirect_uri=https%3A%2F%2Frisken.example%2Fapi%2Fv1%2Fcode%2Fgithub-app%2Foauth%2Fcallback&state=state-value" {
		t.Fatalf("Unexpected URL: %s", got)
	}
}

func TestBuildGitHubAppInstallURL(t *testing.T) {
	svc := &gatewayService{githubAppSlug: "codescan-app-test"}
	got, err := svc.buildGitHubAppInstallURL()
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	if got != "https://github.com/apps/codescan-app-test/installations/select_target" {
		t.Fatalf("Unexpected URL: %s", got)
	}
}

func TestGetGitHubAppInstallationStatusCodeHandler(t *testing.T) {
	codeMock := codemocks.NewCodeServiceClient(t)
	svc := &gatewayService{codeClient: codeMock}
	codeMock.On("GetGitHubAppInstallationStatus", mock.Anything, mock.MatchedBy(func(req *code.GetGitHubAppInstallationStatusRequest) bool {
		return req.ProjectId == 1001 &&
			req.Type == code.Type_ORGANIZATION &&
			req.BaseUrl == "https://api.github.com/" &&
			req.TargetResource == "ca-risken"
	})).Return(&code.GetGitHubAppInstallationStatusResponse{
		GithubAppInstallationStatus: &code.GitHubAppInstallationStatus{
			TargetResource:      "ca-risken",
			Installed:           true,
			RepositorySelection: "selected",
			RepositoryCount:     3,
			Reason:              code.GitHubAppInstallationReasonInstalled,
		},
	}, nil).Once()

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/api/v1/code/github-app/installation-status?project_id=1001&type=1&base_url=https%3A%2F%2Fapi.github.com%2F&target_resource=ca-risken", nil)

	svc.getGitHubAppInstallationStatusCodeHandler(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("Unexpected status. want=%d, got=%d", http.StatusOK, rec.Code)
	}
	resp := map[string]*code.GetGitHubAppInstallationStatusResponse{}
	if err := json.NewDecoder(rec.Body).Decode(&resp); err != nil {
		t.Fatalf("Unexpected json decode error: %v", err)
	}
	status := resp[successJSONKey].GetGithubAppInstallationStatus()
	if !status.GetInstalled() || status.GetTargetResource() != "ca-risken" || status.GetRepositoryCount() != 3 {
		t.Fatalf("Unexpected installation status: %+v", status)
	}
}

func TestBuildGitHubAppInstallURLRejectsInvalidSlug(t *testing.T) {
	cases := []struct {
		name string
		slug string
	}{
		{name: "empty"},
		{name: "blank", slug: "   "},
		{name: "path", slug: "owner/app"},
		{name: "query", slug: "app?state=value"},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			svc := &gatewayService{githubAppSlug: c.slug}
			_, err := svc.buildGitHubAppInstallURL()
			if err == nil {
				t.Fatal("Expected error but got none")
			}
			if strings.TrimSpace(err.Error()) == "" {
				t.Fatal("Expected non-empty error")
			}
		})
	}
}

func TestIsValidGitHubAppSlug(t *testing.T) {
	cases := []struct {
		name string
		slug string
		want bool
	}{
		{name: "lowercase", slug: "codescan-app-test", want: true},
		{name: "number", slug: "codescan-app-1", want: true},
		{name: "empty", want: false},
		{name: "uppercase", slug: "CodeScan-App-Test", want: false},
		{name: "path", slug: "owner/app", want: false},
		{name: "query", slug: "app?state=value", want: false},
		{name: "at sign", slug: "app@example", want: false},
		{name: "plus", slug: "app+test", want: false},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if got := isValidGitHubAppSlug(c.slug); got != c.want {
				t.Fatalf("Unexpected result. want=%t, got=%t", c.want, got)
			}
		})
	}
}

func TestBuildGitHubAppOAuthStartURLAllowsLocalHTTPRedirectURL(t *testing.T) {
	svc := &gatewayService{
		envName:              "local",
		githubAppClientID:    "test-github-app-client-id",
		githubAppRedirectURL: "http://localhost:8080/api/v1/code/github-app/oauth/callback",
	}
	got, err := svc.buildGitHubAppOAuthStartURL("state-value")
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	if got != "https://github.com/login/oauth/authorize?client_id=test-github-app-client-id&redirect_uri=http%3A%2F%2Flocalhost%3A8080%2Fapi%2Fv1%2Fcode%2Fgithub-app%2Foauth%2Fcallback&state=state-value" {
		t.Fatalf("Unexpected URL: %s", got)
	}
}

func TestNormalizeGitHubAppReturnTo(t *testing.T) {
	got, err := normalizeGitHubAppReturnTo("", 1001)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	if got != "/code/github?project_id=1001" {
		t.Fatalf("Unexpected return_to: %s", got)
	}

	if _, err := normalizeGitHubAppReturnTo("https://attacker.example", 1001); err == nil {
		t.Fatal("Expected absolute URL error but got none")
	}
	if _, err := normalizeGitHubAppReturnTo("//attacker.example/path", 1001); err == nil {
		t.Fatal("Expected protocol-relative URL error but got none")
	}
	if _, err := normalizeGitHubAppReturnTo(`/\attacker.example/path`, 1001); err == nil {
		t.Fatal("Expected backslash URL error but got none")
	}
	if _, err := normalizeGitHubAppReturnTo("code/github", 1001); err == nil {
		t.Fatal("Expected non-relative path error but got none")
	}
}

func TestBuildGitHubAppOAuthStartURLRejectsMissingClientID(t *testing.T) {
	cases := []struct {
		name     string
		clientID string
	}{
		{name: "empty"},
		{name: "blank", clientID: "   "},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			svc := &gatewayService{githubAppClientID: c.clientID}
			_, err := svc.buildGitHubAppOAuthStartURL("state")
			if err == nil {
				t.Fatal("Expected error but got none")
			}
			if strings.TrimSpace(err.Error()) == "" {
				t.Fatal("Expected non-empty error")
			}
		})
	}
}

func TestBuildGitHubAppOAuthStartURLRejectsInvalidRedirectURL(t *testing.T) {
	cases := []struct {
		name        string
		redirectURL string
	}{
		{name: "empty"},
		{name: "blank", redirectURL: "   "},
		{name: "relative", redirectURL: "/api/v1/code/github-app/oauth/callback"},
		{name: "host missing", redirectURL: "https:///api/v1/code/github-app/oauth/callback"},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			svc := &gatewayService{
				githubAppClientID:    "test-github-app-client-id",
				githubAppRedirectURL: c.redirectURL,
			}
			_, err := svc.buildGitHubAppOAuthStartURL("state")
			if err == nil {
				t.Fatal("Expected error but got none")
			}
			if strings.TrimSpace(err.Error()) == "" {
				t.Fatal("Expected non-empty error")
			}
		})
	}
}

func TestRedirectGitHubAppOAuthResultRejectsInvalidReturnTo(t *testing.T) {
	cases := []struct {
		name     string
		returnTo string
	}{
		{name: "absolute", returnTo: "https://attacker.example"},
		{name: "backslash", returnTo: `/\attacker.example/path`},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			svc := &gatewayService{}
			req := httptest.NewRequest(http.MethodGet, "/api/v1/code/github-app/oauth/callback", nil)
			rec := httptest.NewRecorder()

			svc.redirectGitHubAppOAuthResult(rec, req, c.returnTo, "success")

			if rec.Code != http.StatusFound {
				t.Fatalf("Unexpected status. want=%d, got=%d", http.StatusFound, rec.Code)
			}
			if got := rec.Header().Get("Location"); got != "/code/github?github_app_oauth=failed" {
				t.Fatalf("Unexpected Location: %s", got)
			}
		})
	}
}

func TestGitHubAppOAuthCallbackHandlerRedirectsWhenSessionExpired(t *testing.T) {
	svc := &gatewayService{githubAppStateSecret: testGitHubAppStateSecret}
	req := newGitHubAppOAuthCallbackRequest(t, svc, 0)
	rec := httptest.NewRecorder()

	svc.githubAppOAuthCallbackHandler(rec, req)

	if rec.Code != http.StatusFound {
		t.Fatalf("Unexpected status. want=%d, got=%d", http.StatusFound, rec.Code)
	}
	if got := rec.Header().Get("Location"); got != "/code/github?github_app_oauth=session_expired&project_id=1001" {
		t.Fatalf("Unexpected Location: %s", got)
	}
}

func TestGitHubAppOAuthCallbackHandlerRedirectsWhenUserMismatch(t *testing.T) {
	svc := &gatewayService{githubAppStateSecret: testGitHubAppStateSecret}
	req := newGitHubAppOAuthCallbackRequest(t, svc, 21)
	rec := httptest.NewRecorder()

	svc.githubAppOAuthCallbackHandler(rec, req)

	if rec.Code != http.StatusFound {
		t.Fatalf("Unexpected status. want=%d, got=%d", http.StatusFound, rec.Code)
	}
	if got := rec.Header().Get("Location"); got != "/code/github?github_app_oauth=unauthorized&project_id=1001" {
		t.Fatalf("Unexpected Location: %s", got)
	}
}

func TestGitHubAppOAuthCallbackHandlerRedirectsWhenProjectUnauthorized(t *testing.T) {
	iamMock := iammocks.NewIAMServiceClient(t)
	iamMock.On("IsAuthorized", mock.Anything, mock.Anything).Return(&iam.IsAuthorizedResponse{Ok: false}, nil).Once()
	svc := &gatewayService{
		githubAppStateSecret: testGitHubAppStateSecret,
		iamClient:            iamMock,
	}
	req := newGitHubAppOAuthCallbackRequest(t, svc, 20)
	rec := httptest.NewRecorder()

	svc.githubAppOAuthCallbackHandler(rec, req)

	if rec.Code != http.StatusFound {
		t.Fatalf("Unexpected status. want=%d, got=%d", http.StatusFound, rec.Code)
	}
	if got := rec.Header().Get("Location"); got != "/code/github?github_app_oauth=unauthorized&project_id=1001" {
		t.Fatalf("Unexpected Location: %s", got)
	}
}
