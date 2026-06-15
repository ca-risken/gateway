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
		name      string
		conf      *AppConfig
		wantError bool
	}{
		{
			name: "github app disabled",
			conf: &AppConfig{},
		},
		{
			name: "github app enabled with strong state secret",
			conf: &AppConfig{
				GithubAppInstallURL:  "https://github.com/apps/risken/installations/new",
				GithubAppStateSecret: "12345678901234567890123456789012",
			},
		},
		{
			name: "github app enabled with short state secret",
			conf: &AppConfig{
				GithubAppInstallURL:  "https://github.com/apps/risken/installations/new",
				GithubAppStateSecret: "short",
			},
			wantError: true,
		},
		{
			name: "github app enabled with http install url",
			conf: &AppConfig{
				GithubAppInstallURL:  "http://github.com/apps/risken/installations/new",
				GithubAppStateSecret: "12345678901234567890123456789012",
			},
			wantError: true,
		},
		{
			name: "github app enabled with relative install url",
			conf: &AppConfig{
				GithubAppInstallURL:  "/apps/risken/installations/new",
				GithubAppStateSecret: "12345678901234567890123456789012",
			},
			wantError: true,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			err := validateGatewayConfig(c.conf)
			if c.wantError && err == nil {
				t.Fatal("expected error but got nil")
			}
			if !c.wantError && err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
		})
	}
}
