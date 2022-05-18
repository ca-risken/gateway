package main

import (
	"context"
	"encoding/base64"
	"testing"
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

func TestDecodeAccessToken(t *testing.T) {
	cases := []struct {
		name              string
		input             string
		wantProjectID     uint32
		wantAccessTokenID uint32
		wantPlainText     string
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
			name:              "Blank",
			input:             "",
			wantAccessTokenID: 0,
			wantPlainText:     "",
		},
		{
			name:              "Invalid token",
			input:             "xxx",
			wantAccessTokenID: 0,
			wantPlainText:     "",
		},
		{
			name:              "Invalid format",
			input:             base64.RawURLEncoding.EncodeToString([]byte("1001/plain_text")), // `/` char is invalid
			wantAccessTokenID: 0,
			wantPlainText:     "",
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			ctx := context.Background()
			gotProjectID, gotAccessTokenID, gotPlainText := decodeAccessToken(ctx, c.input)
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
