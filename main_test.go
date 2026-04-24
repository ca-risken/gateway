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
