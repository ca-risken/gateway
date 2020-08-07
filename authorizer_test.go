package main

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/CyberAgent/mimosa-core/proto/iam"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/stretchr/testify/mock"
	"google.golang.org/grpc"
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
			want:  http.StatusSeeOther,
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
			req, _ := http.NewRequest(http.MethodGet, "/signin", nil)
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
			req, _ := http.NewRequest(http.MethodPost, "/test", nil)
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
			req, _ := http.NewRequest(http.MethodGet, "/service/action?"+c.inputProject, nil)
			got := svc.authzProject(c.inputUser, req)
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
			input: "/service/path1/path2",
			want:  "service/path1",
		},
		{
			name:  "OK No sub paths",
			input: "/service/",
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
			input: "/service-action1-action2",
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
			input: "/service/path1/path2",
			want:  "service",
		},
		{
			name:  "OK No sub paths",
			input: "/service",
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

/**
 * Mock Client
**/
type mockIAMClient struct {
	mock.Mock
}

func (m *mockIAMClient) ListUser(context.Context, *iam.ListUserRequest, ...grpc.CallOption) (*iam.ListUserResponse, error) {
	args := m.Called()
	return args.Get(0).(*iam.ListUserResponse), args.Error(1)
}
func (m *mockIAMClient) GetUser(context.Context, *iam.GetUserRequest, ...grpc.CallOption) (*iam.GetUserResponse, error) {
	args := m.Called()
	return args.Get(0).(*iam.GetUserResponse), args.Error(1)
}
func (m *mockIAMClient) PutUser(context.Context, *iam.PutUserRequest, ...grpc.CallOption) (*iam.PutUserResponse, error) {
	args := m.Called()
	return args.Get(0).(*iam.PutUserResponse), args.Error(1)
}
func (m *mockIAMClient) ListRole(context.Context, *iam.ListRoleRequest, ...grpc.CallOption) (*iam.ListRoleResponse, error) {
	args := m.Called()
	return args.Get(0).(*iam.ListRoleResponse), args.Error(1)
}
func (m *mockIAMClient) GetRole(context.Context, *iam.GetRoleRequest, ...grpc.CallOption) (*iam.GetRoleResponse, error) {
	args := m.Called()
	return args.Get(0).(*iam.GetRoleResponse), args.Error(1)
}
func (m *mockIAMClient) PutRole(context.Context, *iam.PutRoleRequest, ...grpc.CallOption) (*iam.PutRoleResponse, error) {
	args := m.Called()
	return args.Get(0).(*iam.PutRoleResponse), args.Error(1)
}
func (m *mockIAMClient) DeleteRole(context.Context, *iam.DeleteRoleRequest, ...grpc.CallOption) (*empty.Empty, error) {
	args := m.Called()
	return args.Get(0).(*empty.Empty), args.Error(1)
}
func (m *mockIAMClient) AttachRole(context.Context, *iam.AttachRoleRequest, ...grpc.CallOption) (*iam.AttachRoleResponse, error) {
	args := m.Called()
	return args.Get(0).(*iam.AttachRoleResponse), args.Error(1)
}
func (m *mockIAMClient) DetachRole(context.Context, *iam.DetachRoleRequest, ...grpc.CallOption) (*empty.Empty, error) {
	args := m.Called()
	return args.Get(0).(*empty.Empty), args.Error(1)
}
func (m *mockIAMClient) ListPolicy(context.Context, *iam.ListPolicyRequest, ...grpc.CallOption) (*iam.ListPolicyResponse, error) {
	args := m.Called()
	return args.Get(0).(*iam.ListPolicyResponse), args.Error(1)
}
func (m *mockIAMClient) GetPolicy(context.Context, *iam.GetPolicyRequest, ...grpc.CallOption) (*iam.GetPolicyResponse, error) {
	args := m.Called()
	return args.Get(0).(*iam.GetPolicyResponse), args.Error(1)
}
func (m *mockIAMClient) PutPolicy(context.Context, *iam.PutPolicyRequest, ...grpc.CallOption) (*iam.PutPolicyResponse, error) {
	args := m.Called()
	return args.Get(0).(*iam.PutPolicyResponse), args.Error(1)
}
func (m *mockIAMClient) DeletePolicy(context.Context, *iam.DeletePolicyRequest, ...grpc.CallOption) (*empty.Empty, error) {
	args := m.Called()
	return args.Get(0).(*empty.Empty), args.Error(1)
}
func (m *mockIAMClient) AttachPolicy(context.Context, *iam.AttachPolicyRequest, ...grpc.CallOption) (*iam.AttachPolicyResponse, error) {
	args := m.Called()
	return args.Get(0).(*iam.AttachPolicyResponse), args.Error(1)
}
func (m *mockIAMClient) DetachPolicy(context.Context, *iam.DetachPolicyRequest, ...grpc.CallOption) (*empty.Empty, error) {
	args := m.Called()
	return args.Get(0).(*empty.Empty), args.Error(1)
}

func (m *mockIAMClient) IsAuthorized(context.Context, *iam.IsAuthorizedRequest, ...grpc.CallOption) (*iam.IsAuthorizedResponse, error) {
	args := m.Called()
	return args.Get(0).(*iam.IsAuthorizedResponse), args.Error(1)
}
