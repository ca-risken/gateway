package main

import (
	"context"
	"encoding/base64"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEncodeAccessToken(t *testing.T) {
	cases := []struct {
		name          string
		projectID     uint32
		accessTokenID uint32
		plainText     string
		want          string
	}{
		{
			name:          "OK",
			projectID:     111,
			accessTokenID: 222,
			plainText:     "plain_text",
			want:          base64.RawURLEncoding.EncodeToString([]byte("111@222@plain_text")),
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got := encodeAccessToken(c.projectID, c.accessTokenID, c.plainText)
			if got != c.want {
				t.Fatalf("Unexpected responce. want=%s, got=%s", c.want, got)
			}
		})
	}
}

func TestEncodeOrganizationAccessToken(t *testing.T) {
	got := encodeOrganizationAccessToken(1, 2, "plain")
	wantSuffix := encodeAccessToken(1, 2, "plain")
	if !strings.HasPrefix(got, organizationTokenPrefix) {
		t.Fatalf("organization token must have prefix, got=%s", got)
	}
	if !strings.HasSuffix(got, wantSuffix) {
		t.Fatalf("organization token suffix mismatch, want suffix=%s, got=%s", wantSuffix, got)
	}
}

func TestDecodeAccessToken(t *testing.T) {
	cases := []struct {
		name              string
		input             string
		wantProjectID     uint32
		wantAccessTokenID uint32
		wantPlainText     string
		wantErr           bool
	}{
		{
			name:              "OK 1",
			input:             base64.RawURLEncoding.EncodeToString([]byte("111@222@plain_text")),
			wantProjectID:     111,
			wantAccessTokenID: 222,
			wantPlainText:     "plain_text",
		},
		{
			name:              "OK 2",
			input:             base64.RawURLEncoding.EncodeToString([]byte("111@222@333@plain_text")), // too many `@`
			wantProjectID:     111,
			wantAccessTokenID: 222,
			wantPlainText:     "333@plain_text",
		},
		{
			name:    "Blank",
			input:   "",
			wantErr: true,
		},
		{
			name:    "Invalid token",
			input:   "xxx",
			wantErr: true,
		},
		{
			name:    "Invalid format",
			input:   base64.RawURLEncoding.EncodeToString([]byte("1001/plain_text")), // `/` char is invalid
			wantErr: true,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			ctx := context.Background()
			gotProjectID, gotAccessTokenID, gotPlainText, err := decodeAccessToken(ctx, c.input)
			if c.wantErr {
				assert.Error(t, err)
			}
			if gotProjectID != c.wantProjectID {
				t.Fatalf("Unexpected ProjectID. want=%d, got=%d", c.wantProjectID, gotProjectID)
			}
			if gotAccessTokenID != c.wantAccessTokenID {
				t.Fatalf("Unexpected AccessTokenID. want=%d, got=%d", c.wantAccessTokenID, gotAccessTokenID)
			}
			if gotPlainText != c.wantPlainText {
				t.Fatalf("Unexpected responce. wantPlainText=%s, gotPlainText=%s", c.wantPlainText, gotPlainText)
			}
		})
	}
}

func TestDecodeOrganizationAccessToken(t *testing.T) {
	ctx := context.Background()
	ownerID, accessTokenID, plainText, err := decodeOrganizationAccessToken(ctx, encodeOrganizationAccessToken(1, 2, "plain"))
	if err != nil {
		t.Fatalf("unexpected error: %+v", err)
	}
	if ownerID != 1 || accessTokenID != 2 || plainText != "plain" {
		t.Fatalf("unexpected decoded result: %d %d %s", ownerID, accessTokenID, plainText)
	}

	if _, _, _, err := decodeOrganizationAccessToken(ctx, encodeAccessToken(1, 2, "plain")); err == nil {
		t.Fatalf("expected error for token without prefix")
	}
}
