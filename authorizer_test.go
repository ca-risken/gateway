package main

import (
	"context"
	"crypto/ecdsa"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ca-risken/core/proto/iam"
	iammocks "github.com/ca-risken/core/proto/iam/mocks"
	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/mock"
)

func TestSigninHandler(t *testing.T) {
	cases := []struct {
		name  string
		input *requestUser
		want  int
	}{
		{
			name:  "OK",
			input: &requestUser{sub: "sub", userID: 123},
			want:  http.StatusOK,
		},
		{
			name:  "NG No user",
			input: nil,
			want:  http.StatusUnauthorized,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			rec := httptest.NewRecorder()
			req, _ := http.NewRequest(http.MethodGet, "/api/v1/signin", nil)
			req = req.WithContext(context.WithValue(req.Context(), userKey, c.input))
			signinHandler(rec, req)
			got := rec.Result().StatusCode
			if got != c.want {
				t.Fatalf("Unexpected responce. want=%d, got=%d", c.want, got)
			}
		})
	}
}

func TestValidCSRFToken(t *testing.T) {
	cases := []struct {
		name        string
		inputHeader string
		inputCookie string
		want        bool
	}{
		{
			name:        "OK",
			inputHeader: "same_value",
			inputCookie: "same_value",
			want:        true,
		},
		{
			name:        "NG Header blank",
			inputHeader: "",
			inputCookie: "exists",
			want:        false,
		},
		{
			name:        "NG Cookie blank",
			inputHeader: "exists",
			inputCookie: "",
			want:        false,
		},
		{
			name:        "NG Wrong value",
			inputHeader: "wrong",
			inputCookie: "value",
			want:        false,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			req, _ := http.NewRequest(http.MethodPost, "/api/v1/test", nil)
			req.Header.Add("X-XSRF-TOKEN", c.inputHeader)
			req.AddCookie(&http.Cookie{Name: "XSRF-TOKEN", Value: c.inputCookie})
			got := validCSRFToken(req)
			if got != c.want {
				t.Fatalf("Unexpected response. want=%t, got=%t", c.want, got)
			}
		})
	}
}

func TestAuthzProject(t *testing.T) {
	iamMock := iammocks.NewIAMServiceClient(t)
	svc := gatewayService{
		iamClient: iamMock,
	}
	cases := []struct {
		name         string
		inputUser    *requestUser
		inputProject string
		want         bool
		mockResp     *iam.IsAuthorizedResponse
		mockErr      error
	}{
		{
			name:         "OK",
			inputUser:    &requestUser{sub: "sub", userID: 123},
			inputProject: "project_id=1",
			mockResp:     &iam.IsAuthorizedResponse{Ok: true},
			want:         true,
		},
		{
			name:         "NG Invalid user",
			inputUser:    &requestUser{sub: "sub"},
			inputProject: "project_id=1",
			want:         false,
		},
		{
			name:         "NG Invalid project",
			inputUser:    &requestUser{sub: "sub", userID: 123},
			inputProject: "project_id=aaa",
			want:         false,
		},
		{
			name:         "NG IAM error",
			inputUser:    &requestUser{sub: "sub", userID: 123},
			inputProject: "project_id=1",
			want:         false,
			mockErr:      errors.New("something error"),
		}}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if c.mockResp != nil || c.mockErr != nil {
				iamMock.On("IsAuthorized", mock.Anything, mock.Anything).Return(c.mockResp, c.mockErr).Once()
			}
			req, _ := http.NewRequest(http.MethodGet, "/api/v1/service/action?"+c.inputProject, nil)
			got := svc.authzProject(c.inputUser, req)
			if got != c.want {
				t.Fatalf("Unexpected response. want=%t, got=%t", c.want, got)
			}
		})
	}
}

func TestAuthzProjectForToken(t *testing.T) {
	iamMock := iammocks.NewIAMServiceClient(t)
	svc := gatewayService{
		iamClient: iamMock,
	}
	cases := []struct {
		name         string
		inputUser    *requestUser
		inputProject string
		want         bool
		mockResp     *iam.IsAuthorizedTokenResponse
		mockErr      error
	}{
		{
			name:         "OK",
			inputUser:    &requestUser{accessTokenID: 123},
			inputProject: "project_id=1",
			mockResp:     &iam.IsAuthorizedTokenResponse{Ok: true},
			want:         true,
		},
		{
			name:         "NG No token",
			inputUser:    &requestUser{sub: "sub"},
			inputProject: "project_id=1",
			want:         false,
		},
		{
			name:         "NG Invalid project",
			inputUser:    &requestUser{accessTokenID: 123},
			inputProject: "project_id=aaa",
			want:         false,
		},
		{
			name:         "NG IAM error",
			inputUser:    &requestUser{accessTokenID: 123},
			inputProject: "project_id=1",
			want:         false,
			mockErr:      errors.New("something error"),
		}}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if c.mockResp != nil || c.mockErr != nil {
				iamMock.On("IsAuthorizedToken", mock.Anything, mock.Anything).Return(c.mockResp, c.mockErr).Once()
			}
			req, _ := http.NewRequest(http.MethodGet, "/api/v1/service/action?"+c.inputProject, nil)
			got := svc.authzProjectForToken(c.inputUser, req)
			if got != c.want {
				t.Fatalf("Unexpected response. want=%t, got=%t", c.want, got)
			}
		})
	}
}

func TestGetActionNameFromURI(t *testing.T) {
	cases := []struct {
		name  string
		input string
		want  string
	}{
		{
			name:  "OK",
			input: "/api/v1/service/path1/path2",
			want:  "service/path1",
		},
		{
			name:  "OK No sub paths",
			input: "/api/v1/service/",
			want:  "service/",
		},
		{
			name:  "NG blank",
			input: "",
			want:  "",
		},
		{
			name:  "NG No prefix(/)",
			input: "service-action1-action2",
			want:  "",
		},
		{
			name:  "NG No sub slashes",
			input: "/api/v1/service-action1-action2",
			want:  "",
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got := getActionNameFromURI(c.input)
			if got != c.want {
				t.Fatalf("Unexpected response. want=%s, got=%s", c.want, got)
			}
		})
	}
}

func TestGetServiceNameFromURI(t *testing.T) {
	cases := []struct {
		name  string
		input string
		want  string
	}{
		{
			name:  "OK",
			input: "/api/v1/service/path1/path2",
			want:  "service",
		},
		{
			name:  "OK No sub paths",
			input: "/api/v1/service",
			want:  "service",
		},
		{
			name:  "NG blank",
			input: "",
			want:  "",
		},
		{
			name:  "NG No prefix(/)",
			input: "service-action1-action2",
			want:  "",
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got := getServiceNameFromURI(c.input)
			if got != c.want {
				t.Fatalf("Unexpected response. want=%s, got=%s", c.want, got)
			}
		})
	}
}

func TestAuthzAdmin(t *testing.T) {
	iamMock := iammocks.NewIAMServiceClient(t)
	svc := gatewayService{
		iamClient: iamMock,
	}
	cases := []struct {
		name          string
		inputUser     *requestUser
		inputProject  string
		want          bool
		mockAuthzResp *iam.IsAuthorizedAdminResponse
		mockAuthzErr  error
	}{
		{
			name:          "OK",
			inputUser:     &requestUser{sub: "sub", userID: 1},
			inputProject:  "project_id=1",
			mockAuthzResp: &iam.IsAuthorizedAdminResponse{Ok: true},
			want:          true,
		},
		{
			name:         "NG Invalid userID",
			inputUser:    &requestUser{sub: "sub", userID: 0},
			inputProject: "project_id=1",
			want:         false,
		},
		{
			name:         "NG Authz API error",
			inputUser:    &requestUser{sub: "sub", userID: 1},
			inputProject: "project_id=1",
			want:         false,
			mockAuthzErr: errors.New("something error"),
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if c.mockAuthzResp != nil || c.mockAuthzErr != nil {
				iamMock.On("IsAuthorizedAdmin", mock.Anything, mock.Anything).Return(c.mockAuthzResp, c.mockAuthzErr).Once()
			}
			req, _ := http.NewRequest(http.MethodGet, "/api/v1/admin/api/?"+c.inputProject, nil)
			got := svc.authzAdmin(c.inputUser, req)
			if got != c.want {
				t.Fatalf("Unexpected response. want=%t, got=%t", c.want, got)
			}
			c.mockAuthzResp = nil
			c.mockAuthzErr = nil
		})
	}
}

func TestIsHumanAccess(t *testing.T) {
	cases := []struct {
		name  string
		input *requestUser
		want  bool
	}{
		{
			name:  "Human",
			input: &requestUser{sub: "sub", userID: 1},
			want:  true,
		},
		{
			name:  "Not human 1",
			input: &requestUser{sub: "sub", accessTokenID: 1},
			want:  false,
		},
		{
			name:  "Not human 2",
			input: &requestUser{sub: "sub", accessTokenID: 1, userID: 1},
			want:  false,
		},
		{
			name:  "Nil",
			input: nil,
			want:  false,
		},
		{
			name:  "No userID",
			input: &requestUser{sub: "sub"},
			want:  false,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got := isHumanAccess(c.input)
			if got != c.want {
				t.Fatalf("Unexpected response. want=%t, got=%t", c.want, got)
			}
		})
	}
}

func TestShouldVerifyCSRFTokenURI(t *testing.T) {
	cases := []struct {
		name  string
		input string
		want  bool
	}{
		{
			name:  "should verify 1",
			input: "/api/v1/uri",
			want:  true,
		},
		{
			name:  "ignore URI 1",
			input: "/healthz",
			want:  false,
		},
		{
			name:  "ignore URI 2",
			input: "/api/v1/signin",
			want:  false,
		},
		{
			name:  "slash suffix 1",
			input: "/api/v1/uri/",
			want:  true,
		},
		{
			name:  "slash suffix 2",
			input: "/api/v1/signin/",
			want:  false,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got := shouldVerifyCSRFTokenURI(c.input)
			if got != c.want {
				t.Fatalf("Unexpected response. want=%t, got=%t", c.want, got)
			}
		})
	}
}

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
			g := gatewayService{
				Region: "ap-northeast-1",
			}
			err := g.verifyTokenForALB(context.Background(), c.tokenString)
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
			g := gatewayService{
				Region: "ap-northeast-1",
			}
			_, err := g.fetchALBPublicKey(context.Background(), c.keyURL)
			if (c.wantErr && err == nil) || (!c.wantErr && err != nil) {
				t.Fatalf("Unexpected error: wantErr=%t, err=%+v", c.wantErr, err)
			}
		})
	}
}
