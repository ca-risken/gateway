package main

import (
	"context"
	"crypto/ecdsa"
	"errors"
	"net/http"
	"reflect"
	"testing"

	"github.com/golang-jwt/jwt/v4"
	"github.com/jarcoal/httpmock"
)

func TestVerifyTokenForALB(t *testing.T) {
	cases := []struct {
		name         string
		tokenString  string
		keyURL       string
		mockStatus   int
		mockResponse string
		mockErr      error
		wantErr      bool
	}{
		{
			name:        "NG alg is not ES256",
			tokenString: "YWxnOmludmFsaWQgdHlwOkpXVAo.eyJmb28iOiJiYXIifQ.MEQCIHoSJnmGlPaVQDqacx_2XlXEhhqtWceVopjomc2PJLtdAiAUTeGPoNYxZw0z8mgOnnIcjoxRuNDVZvybRZF3wR1l8W",
			wantErr:     true,
		},
		{
			name:        "NG kid doesn't exist in JWT header",
			tokenString: "eyJhbGciOiJFUzI1NiIsInR5cCI6IkpXVCJ9.eyJmb28iOiJiYXIifQ.MEQCIHoSJnmGlPaVQDqacx_2XlXEhhqtWceVopjomc2PJLtdAiAUTeGPoNYxZw0z8mgOnnIcjoxRuNDVZvybRZF3wR1l8W",
			wantErr:     true,
		},
		{
			name:        "NG fetch Public key error",
			tokenString: "eyJhbGciOiJFUzI1NiIsInR5cCI6IkpXVCIsImtpZCI6ImhvZ2UifQo.eyJmb28iOiJiYXIifQ.MEQCIHoSJnmGlPaVQDqacx_2XlXEhhqtWceVopjomc2PJLtdAiAUTeGPoNYxZw0z8mgOnnIcjoxRuNDVZvybRZF3wR1l8W",
			keyURL:      "https://public-keys.auth.elb.ap-northeast-1.amazonaws.com/hoge",
			mockErr:     errors.New("something error"),
			wantErr:     true,
		},
		{
			name:        "NG verify Error",
			tokenString: "eyJhbGciOiJFUzI1NiIsInR5cCI6IkpXVCIsImtpZCI6ImhvZ2UifQo.eyJmb28iOiJiYXIifQ.MEQCIHoSJnmGlPaVQDqacx_2XlXEhhqtWceVopjomc2PJLtdAiAUTeGPoNYxZw0z8mgOnnIcjoxRuNDVZvybRZF3wR1l8W",
			keyURL:      "https://public-keys.auth.elb.ap-northeast-1.amazonaws.com/hoge",
			mockStatus:  200,
			mockResponse: `-----BEGIN PUBLIC KEY-----
MFkwEwYHKoZIzj0CAQYIKoZIzj0DAQcDQgAEYD54V/vp+54P9DXarYqx4MPcm+HK
RIQzNasYSoRQHQ/6S6Ps8tpMcT+KvIIC8W/e9k0W7Cm72M1P9jU7SLf/vg==
-----END PUBLIC KEY-----
`,
			wantErr: true,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if c.keyURL != "" {
				httpmock.Activate()
				defer httpmock.DeactivateAndReset()
				if c.mockErr == nil {
					httpmock.RegisterResponder("GET", c.keyURL,
						httpmock.NewStringResponder(c.mockStatus, c.mockResponse))
				} else {
					httpmock.RegisterResponder("GET", c.keyURL,
						func(req *http.Request) (*http.Response, error) {
							return nil, c.mockErr
						})
				}
			}
			cli := claimsClient{
				region: "ap-northeast-1",
			}
			_, err := cli.getClaimsForALB(context.Background(), c.tokenString)
			if (c.wantErr && err == nil) || (!c.wantErr && err != nil) {
				t.Fatalf("Unexpected error: wantErr=%t, err=%+v", c.wantErr, err)
			}
		})
	}
}
func TestFetchPublicKey(t *testing.T) {
	cases := []struct {
		name         string
		keyURL       string
		mockStatus   int
		mockResponse string
		mockErr      error
		want         *ecdsa.PublicKey
		wantErr      bool
	}{
		{
			name:       "OK",
			keyURL:     "http://example.com/valid",
			mockStatus: 200,
			mockResponse: `-----BEGIN PUBLIC KEY-----
MFkwEwYHKoZIzj0CAQYIKoZIzj0DAQcDQgAEYD54V/vp+54P9DXarYqx4MPcm+HK
RIQzNasYSoRQHQ/6S6Ps8tpMcT+KvIIC8W/e9k0W7Cm72M1P9jU7SLf/vg==
-----END PUBLIC KEY-----
`,
			wantErr: false,
		},
		{
			name:         "NG Parse Error",
			keyURL:       "http://example.com/invalid",
			mockStatus:   200,
			mockResponse: `Invalid Key`,
			wantErr:      true,
		},
		{
			name:    "NG HTTP Error",
			keyURL:  "http://example.com/error",
			mockErr: errors.New("something error"),
			wantErr: true,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			httpmock.Activate()
			defer httpmock.DeactivateAndReset()
			if c.mockErr == nil {
				httpmock.RegisterResponder("GET", c.keyURL,
					httpmock.NewStringResponder(c.mockStatus, c.mockResponse))
			} else {
				httpmock.RegisterResponder("GET", c.keyURL,
					func(req *http.Request) (*http.Response, error) {
						return nil, c.mockErr
					})
			}
			_, err := fetchALBPublicKey(context.Background(), c.keyURL)
			if (c.wantErr && err == nil) || (!c.wantErr && err != nil) {
				t.Fatalf("Unexpected error: wantErr=%t, err=%+v", c.wantErr, err)
			}
		})
	}
}

func TestGetClaimsWithoutVerify(t *testing.T) {
	cases := []struct {
		name    string
		token   string
		want    *jwt.MapClaims
		wantErr bool
	}{
		{
			name:  "OK Get Claims",
			token: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiJ0ZXN0In0K.1EmGW2C02O8vcS_w4-OjOjHnegrYuBiBg3okk-WY_qI",
			want: &jwt.MapClaims{
				"sub": "test",
			},
		},
		{
			name:    "NG invalid string pattern",
			token:   "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiJzdWIiLCJpYXQiOjE1MTYyMzkwMjJ9",
			wantErr: true,
		},
		{
			name:    "NG base64 decode error",
			token:   "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiJzdWIiLCJpYXQiOjE1MTYyMzkwMjJ9###.3nBr2r_c2ukFrbhh_aomGSP4B9CCdtnwgcIlp2LN5GE",
			wantErr: true,
		},
		{
			name:    "NG jwt unmarshal error",
			token:   "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.aG9nZWZ1Z2EK.3nBr2r_c2ukFrbhh_aomGSP4B9CCdtnwgcIlp2LN5GE",
			wantErr: true,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			cli := claimsClient{}
			got, err := cli.getClaimsWithoutVerify(c.token)
			if (c.wantErr && err == nil) || (!c.wantErr && err != nil) {
				t.Fatalf("Unexpected error: wantErr=%t, err=%+v", c.wantErr, err)
			}
			if !reflect.DeepEqual(got, c.want) {
				t.Fatalf("Unexpected response. want=%+v, got=%+v", c.want, got)

			}
		})
	}
}

func TestGetUserName(t *testing.T) {
	cases := []struct {
		name            string
		token           *jwt.MapClaims
		idpProviderName []string
		want            string
	}{
		{
			name: "OK Get UserName (Not Trim)",
			token: &jwt.MapClaims{
				"username": "value",
			},
			idpProviderName: []string{"IDP1", "IDP2"},
			want:            "value",
		},
		{
			name: "OK Get Attribute (Trim)",
			token: &jwt.MapClaims{
				"username": "IDP1_value",
			},
			idpProviderName: []string{"IDP1", "IDP2"},
			want:            "value",
		},
		{
			name: "OK Not Found username",
			token: &jwt.MapClaims{
				"field": "value",
			},
			idpProviderName: []string{"IDP1", "IDP2"},
			want:            "",
		},
		{
			name: "OK failed to convert username",
			token: &jwt.MapClaims{
				"username": nil,
			},
			idpProviderName: []string{"IDP1", "IDP2"},
			want:            "",
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			cli := claimsClient{
				idpProviderName: c.idpProviderName,
			}
			got := cli.getUserName(c.token)
			if got != c.want {
				t.Fatalf("Unexpected response. want=%v, got=%v", c.want, got)
			}
		})
	}
}

func TestGetUserIdpKey(t *testing.T) {
	cases := []struct {
		name   string
		token  *jwt.MapClaims
		idpKey string
		want   string
	}{
		{
			name: "OK Get value",
			token: &jwt.MapClaims{
				"key": "value",
			},
			idpKey: "key",
			want:   "value",
		},
		{
			name: "OK failed to convert user_idp_key",
			token: &jwt.MapClaims{
				"key": nil,
			},
			idpKey: "key",
			want:   "",
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			cli := claimsClient{
				userIdpKey: c.idpKey,
			}
			got := cli.getUserIdpKey(c.token)
			if got != c.want {
				t.Fatalf("Unexpected response. want=%v, got=%v", c.want, got)
			}
		})
	}
}
