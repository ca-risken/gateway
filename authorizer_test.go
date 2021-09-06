package main

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ca-risken/core/proto/iam"
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
	iamMock := &mockIAMClient{}
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
				iamMock.On("IsAuthorized").Return(c.mockResp, c.mockErr).Once()
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
	iamMock := &mockIAMClient{}
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
				iamMock.On("IsAuthorizedToken").Return(c.mockResp, c.mockErr).Once()
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
	iamMock := &mockIAMClient{}
	svc := gatewayService{
		iamClient: iamMock,
	}
	cases := []struct {
		name          string
		inputUser     *requestUser
		inputProject  string
		want          bool
		mockAdminResp *iam.IsAdminResponse
		mockAdminErr  error
		mockAuthzResp *iam.IsAuthorizedResponse
		mockAuthzErr  error
	}{
		{
			name:          "OK",
			inputUser:     &requestUser{sub: "sub", userID: 1},
			inputProject:  "project_id=1",
			mockAdminResp: &iam.IsAdminResponse{Ok: true},
			mockAuthzResp: &iam.IsAuthorizedResponse{Ok: true},
			want:          true,
		},
		{
			name:         "NG Invalid userID",
			inputUser:    &requestUser{sub: "sub", userID: 0},
			inputProject: "project_id=1",
			want:         false,
		},
		{
			name:         "NG Admin API error",
			inputUser:    &requestUser{sub: "sub", userID: 1},
			inputProject: "project_id=1",
			want:         false,
			mockAdminErr: errors.New("something error"),
		},
		{
			name:          "NG Authz API error",
			inputUser:     &requestUser{sub: "sub", userID: 1},
			inputProject:  "project_id=1",
			want:          false,
			mockAdminResp: &iam.IsAdminResponse{Ok: true},
			mockAuthzErr:  errors.New("something error"),
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if c.mockAdminResp != nil || c.mockAdminErr != nil {
				iamMock.On("IsAdmin").Return(c.mockAdminResp, c.mockAdminErr).Once()
			}
			if c.mockAuthzResp != nil || c.mockAuthzErr != nil {
				iamMock.On("IsAuthorized").Return(c.mockAuthzResp, c.mockAuthzErr).Once()
			}
			req, _ := http.NewRequest(http.MethodGet, "/api/v1/admin/api/?"+c.inputProject, nil)
			got := svc.authzAdmin(c.inputUser, req)
			if got != c.want {
				t.Fatalf("Unexpected response. want=%t, got=%t", c.want, got)
			}
			c.mockAdminResp = nil
			c.mockAdminErr = nil
			c.mockAuthzResp = nil
			c.mockAuthzErr = nil

		})
	}
}
