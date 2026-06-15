package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestGitHubAppOAuthState(t *testing.T) {
	svc := &gatewayService{githubAppStateSecret: "state-secret"}
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

func TestVerifyGitHubAppOAuthStateRejectsInvalidState(t *testing.T) {
	svc := &gatewayService{githubAppStateSecret: "state-secret"}
	now := time.Unix(1700000000, 0)
	rawState, err := svc.newGitHubAppOAuthState(1001, 10, 20, "/code/github?project_id=1001", now)
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

func TestBuildGitHubAppInstallURL(t *testing.T) {
	svc := &gatewayService{githubAppInstallURL: "https://github.com/apps/risken-codescan/installations/new?existing=1"}
	got, err := svc.buildGitHubAppInstallURL("state-value")
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	if got != "https://github.com/apps/risken-codescan/installations/new?existing=1&state=state-value" {
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

func TestBuildGitHubAppInstallURLRejectsInvalidConfig(t *testing.T) {
	cases := []struct {
		name       string
		installURL string
	}{
		{name: "empty"},
		{name: "http", installURL: "http://github.com/apps/risken/installations/new"},
		{name: "relative", installURL: "/apps/risken/installations/new"},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			svc := &gatewayService{githubAppInstallURL: c.installURL}
			_, err := svc.buildGitHubAppInstallURL("state")
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
