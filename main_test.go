package main

import (
	"net/http"
	"testing"
	"time"
)

func TestNewHTTPServer(t *testing.T) {
	conf := &AppConfig{
		Port:                 "8080",
		ReadHeaderTimeoutSec: 5,
		ReadTimeoutSec:       30,
		WriteTimeoutSec:      31,
		IdleTimeoutSec:       120,
		MaxHeaderBytes:       4096,
	}

	server := newHTTPServer(conf, http.NewServeMux())

	if server.Addr != ":8080" {
		t.Fatalf("unexpected addr: %s", server.Addr)
	}
	if server.ReadHeaderTimeout != 5*time.Second {
		t.Fatalf("unexpected read header timeout: %s", server.ReadHeaderTimeout)
	}
	if server.ReadTimeout != 30*time.Second {
		t.Fatalf("unexpected read timeout: %s", server.ReadTimeout)
	}
	if server.WriteTimeout != 31*time.Second {
		t.Fatalf("unexpected write timeout: %s", server.WriteTimeout)
	}
	if server.IdleTimeout != 120*time.Second {
		t.Fatalf("unexpected idle timeout: %s", server.IdleTimeout)
	}
	if server.MaxHeaderBytes != 4096 {
		t.Fatalf("unexpected max header bytes: %d", server.MaxHeaderBytes)
	}
}

func TestValidateGatewayConfig(t *testing.T) {
	cases := []struct {
		name       string
		conf       *AppConfig
		wantErrMsg string
	}{
		{
			name: "github app disabled",
			conf: &AppConfig{},
		},
		{
			name: "github app slug configured",
			conf: &AppConfig{
				GithubAppSlug: "codescan-app-test",
			},
		},
		{
			name: "github app slug invalid",
			conf: &AppConfig{
				GithubAppSlug: "owner/app",
			},
			wantErrMsg: "github app slug is invalid",
		},
		{
			name: "github app disabled with short state secret",
			conf: &AppConfig{
				GithubAppStateSecret: "short",
			},
			wantErrMsg: "github app state secret must be empty or at least 32 bytes when github app oauth client id is not configured",
		},
		{
			name: "github app enabled with strong state secret",
			conf: &AppConfig{
				GithubAppOAuthClientID:    "test-github-app-client-id",
				GithubAppOAuthRedirectURL: "https://risken.example/api/v1/code/github-app/oauth/callback",
				GithubAppStateSecret:      "12345678901234567890123456789012",
			},
		},
		{
			name: "github app enabled with short state secret",
			conf: &AppConfig{
				GithubAppOAuthClientID:    "test-github-app-client-id",
				GithubAppOAuthRedirectURL: "https://risken.example/api/v1/code/github-app/oauth/callback",
				GithubAppStateSecret:      "short",
			},
			wantErrMsg: "github app state secret must be at least 32 bytes when github app oauth client id is configured",
		},
		{
			name: "github app enabled without redirect url",
			conf: &AppConfig{
				GithubAppOAuthClientID: "test-github-app-client-id",
				GithubAppStateSecret:   "12345678901234567890123456789012",
			},
			wantErrMsg: "github app oauth redirect url is required when github app oauth client id is configured",
		},
		{
			name: "github app enabled with relative redirect url",
			conf: &AppConfig{
				GithubAppOAuthClientID:    "test-github-app-client-id",
				GithubAppOAuthRedirectURL: "/api/v1/code/github-app/oauth/callback",
				GithubAppStateSecret:      "12345678901234567890123456789012",
			},
			wantErrMsg: "github app oauth redirect url must be an absolute URL",
		},
		{
			name: "github app enabled with local http redirect url",
			conf: &AppConfig{
				EnvName:                   "local",
				GithubAppOAuthClientID:    "test-github-app-client-id",
				GithubAppOAuthRedirectURL: "http://localhost:8080/api/v1/code/github-app/oauth/callback",
				GithubAppStateSecret:      "12345678901234567890123456789012",
			},
		},
		{
			name: "github app enabled with production http redirect url",
			conf: &AppConfig{
				EnvName:                   "prod",
				GithubAppOAuthClientID:    "test-github-app-client-id",
				GithubAppOAuthRedirectURL: "http://localhost:8080/api/v1/code/github-app/oauth/callback",
				GithubAppStateSecret:      "12345678901234567890123456789012",
			},
			wantErrMsg: "github app oauth redirect url must use https except local localhost",
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			err := validateGatewayConfig(c.conf)
			if c.wantErrMsg != "" && err == nil {
				t.Fatal("expected error but got nil")
			}
			if c.wantErrMsg == "" && err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if c.wantErrMsg != "" && err.Error() != c.wantErrMsg {
				t.Fatalf("unexpected error. want=%q, got=%q", c.wantErrMsg, err.Error())
			}
		})
	}
}
