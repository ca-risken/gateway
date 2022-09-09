package main

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/ca-risken/core/proto/iam"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/stretchr/testify/mock"
	"google.golang.org/grpc"
)

func TestPutUserHandler(t *testing.T) {
	iamMock := &mockIAMClient{}
	svc := gatewayService{
		iamClient: iamMock,
	}
	cases := []struct {
		name       string
		input      string
		mockResp   *iam.PutUserResponse
		mockErr    error
		wantStatus int
	}{
		{
			name:       "OK",
			input:      `{"user": {"sub":"xxx", "name":"nm", "activated":"true"}}`,
			mockResp:   &iam.PutUserResponse{},
			wantStatus: http.StatusOK,
		},
		{
			name:       "NG Invalid parameter",
			input:      `invalid_param`,
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "NG Backend service error",
			input:      `{"user": {"sub":"xxx", "name":"nm", "activated":"true"}}`,
			wantStatus: http.StatusInternalServerError,
			mockErr:    errors.New("something wrong"),
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if c.mockResp != nil || c.mockErr != nil {
				iamMock.On("PutUser").Return(c.mockResp, c.mockErr).Once()
			}
			// Invoke HTTP Request
			rec := httptest.NewRecorder()
			req, _ := http.NewRequest(http.MethodPost, "/api/v1/iam/put-user/", strings.NewReader(c.input))
			req = req.WithContext(context.WithValue(req.Context(), userKey, &requestUser{sub: "xxx", userID: 1}))
			req.Header.Add("Content-Type", "application/json")
			svc.putUserHandler(rec, req)
			// Check Response
			if c.wantStatus != rec.Code {
				t.Fatalf("Unexpected HTTP status code: want=%+v, got=%+v", c.wantStatus, rec.Code)
			}
			resp := map[string]interface{}{}
			if err := json.NewDecoder(rec.Body).Decode(&resp); err != nil {
				t.Fatalf("Unexpected json decode error to response body: err=%+v", err)
			}
			jsonKey := successJSONKey
			if c.wantStatus != http.StatusOK {
				jsonKey = errorJSONKey
			}
			if _, ok := resp[jsonKey]; !ok {
				t.Fatalf("Unexpected no response key: want key=%s", jsonKey)
			}
		})
	}
}

func TestListAdminRoleHandler(t *testing.T) {
	iamMock := &mockIAMClient{}
	svc := gatewayService{
		iamClient: iamMock,
	}
	cases := []struct {
		name       string
		input      string
		mockResp   *iam.ListRoleResponse
		mockErr    error
		wantStatus int
	}{
		{
			name:       "OK",
			input:      ``,
			mockResp:   &iam.ListRoleResponse{},
			wantStatus: http.StatusOK,
		},
		{
			name:       "NG Backend service error",
			input:      ``,
			wantStatus: http.StatusInternalServerError,
			mockErr:    errors.New("something wrong"),
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if c.mockResp != nil || c.mockErr != nil {
				iamMock.On("ListRole").Return(c.mockResp, c.mockErr).Once()
			}
			rec := httptest.NewRecorder()
			req, _ := http.NewRequest(http.MethodGet, "/api/v1/iam/list-admin-role?"+c.input, nil)
			svc.listAdminRoleHandler(rec, req)
			if c.wantStatus != rec.Code {
				t.Fatalf("Unexpected HTTP status code: want=%+v, got=%+v", c.wantStatus, rec.Code)
			}
			resp := map[string]interface{}{}
			if err := json.NewDecoder(rec.Body).Decode(&resp); err != nil {
				t.Fatalf("Unexpected json decode error to response body: err=%+v", err)
			}
			jsonKey := successJSONKey
			if c.wantStatus != http.StatusOK {
				jsonKey = errorJSONKey
			}
			if _, ok := resp[jsonKey]; !ok {
				t.Fatalf("Unexpected no response key: want key=%s", jsonKey)
			}
		})
	}
}

func TestGetAdminRoleHandler(t *testing.T) {
	iamMock := &mockIAMClient{}
	svc := gatewayService{
		iamClient: iamMock,
	}
	cases := []struct {
		name       string
		input      string
		mockResp   *iam.GetRoleResponse
		mockErr    error
		wantStatus int
	}{
		{
			name:       "OK",
			input:      `role_id=1`,
			mockResp:   &iam.GetRoleResponse{},
			wantStatus: http.StatusOK,
		},
		{
			name:       "NG Invalid parameter",
			input:      `no_param`,
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "NG Backend service error",
			input:      `role_id=1`,
			wantStatus: http.StatusInternalServerError,
			mockErr:    errors.New("something wrong"),
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if c.mockResp != nil || c.mockErr != nil {
				iamMock.On("GetRole").Return(c.mockResp, c.mockErr).Once()
			}
			rec := httptest.NewRecorder()
			req, _ := http.NewRequest(http.MethodGet, "/api/v1/iam/get-admin-role?"+c.input, nil)
			svc.getAdminRoleHandler(rec, req)
			if c.wantStatus != rec.Code {
				t.Fatalf("Unexpected HTTP status code: want=%+v, got=%+v", c.wantStatus, rec.Code)
			}
			resp := map[string]interface{}{}
			if err := json.NewDecoder(rec.Body).Decode(&resp); err != nil {
				t.Fatalf("Unexpected json decode error to response body: err=%+v", err)
			}
			jsonKey := successJSONKey
			if c.wantStatus != http.StatusOK {
				jsonKey = errorJSONKey
			}
			if _, ok := resp[jsonKey]; !ok {
				t.Fatalf("Unexpected no response key: want key=%s", jsonKey)
			}
		})
	}
}

func TestAttachAdminRoleHandler(t *testing.T) {
	iamMock := &mockIAMClient{}
	svc := gatewayService{
		iamClient: iamMock,
	}
	cases := []struct {
		name       string
		input      string
		mockResp   *iam.AttachRoleResponse
		mockErr    error
		wantStatus int
	}{
		{
			name:       "OK",
			input:      `{"role_id":1, "user_id":1}`,
			mockResp:   &iam.AttachRoleResponse{},
			wantStatus: http.StatusOK,
		},
		{
			name:       "NG Invalid parameter",
			input:      `invalid_param`,
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "NG Backend service error",
			input:      `{"role_id":1, "user_id":1}`,
			wantStatus: http.StatusInternalServerError,
			mockErr:    errors.New("something wrong"),
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if c.mockResp != nil || c.mockErr != nil {
				iamMock.On("AttachRole").Return(c.mockResp, c.mockErr).Once()
			}
			rec := httptest.NewRecorder()
			req, _ := http.NewRequest(http.MethodPost, "/api/v1/iam/attach-admin-role/", strings.NewReader(c.input))
			req.Header.Add("Content-Type", "application/json")
			svc.attachAdminRoleHandler(rec, req)
			if c.wantStatus != rec.Code {
				t.Fatalf("Unexpected HTTP status code: want=%+v, got=%+v", c.wantStatus, rec.Code)
			}
			resp := map[string]interface{}{}
			if err := json.NewDecoder(rec.Body).Decode(&resp); err != nil {
				t.Fatalf("Unexpected json decode error to response body: err=%+v", err)
			}
			jsonKey := successJSONKey
			if c.wantStatus != http.StatusOK {
				jsonKey = errorJSONKey
			}
			if _, ok := resp[jsonKey]; !ok {
				t.Fatalf("Unexpected no response key: want key=%s", jsonKey)
			}
		})
	}
}

func TestDetachAdminRoleHandler(t *testing.T) {
	iamMock := &mockIAMClient{}
	svc := gatewayService{
		iamClient: iamMock,
	}
	cases := []struct {
		name       string
		input      string
		mockResp   *empty.Empty
		mockErr    error
		wantStatus int
	}{
		{
			name:       "OK",
			input:      `{"role_id":1, "user_id":1}`,
			mockResp:   &empty.Empty{},
			wantStatus: http.StatusOK,
		},
		{
			name:       "NG Invalid parameter",
			input:      `invalid_param`,
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "NG Backend service error",
			input:      `{"role_id":1, "user_id":1}`,
			wantStatus: http.StatusInternalServerError,
			mockErr:    errors.New("something wrong"),
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if c.mockResp != nil || c.mockErr != nil {
				iamMock.On("DetachRole").Return(c.mockResp, c.mockErr).Once()
			}
			rec := httptest.NewRecorder()
			req, _ := http.NewRequest(http.MethodPost, "/api/v1/iam/detach-admin-role/", strings.NewReader(c.input))
			req.Header.Add("Content-Type", "application/json")
			svc.detachAdminRoleHandler(rec, req)
			if c.wantStatus != rec.Code {
				t.Fatalf("Unexpected HTTP status code: want=%+v, got=%+v", c.wantStatus, rec.Code)
			}
			resp := map[string]interface{}{}
			if err := json.NewDecoder(rec.Body).Decode(&resp); err != nil {
				t.Fatalf("Unexpected json decode error to response body: err=%+v", err)
			}
			jsonKey := successJSONKey
			if c.wantStatus != http.StatusOK {
				jsonKey = errorJSONKey
			}
			if _, ok := resp[jsonKey]; !ok {
				t.Fatalf("Unexpected no response key: want key=%s", jsonKey)
			}
		})
	}
}

func TestUpdateAccessTokenHandler(t *testing.T) {
	iamMock := &mockIAMClient{}
	svc := gatewayService{
		iamClient: iamMock,
	}
	cases := []struct {
		name       string
		input      string
		mockResp   *iam.PutAccessTokenResponse
		mockErr    error
		wantStatus int
	}{
		{
			name:       "OK",
			input:      `{"project_id":1, "access_token": {"project_id":1, "access_token_id":1, "name":"nm", "last_updated_user_id":1}}`,
			mockResp:   &iam.PutAccessTokenResponse{},
			wantStatus: http.StatusOK,
		},
		{
			name:       "NG Invalid parameter",
			input:      `invalid_param`,
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "NG Backend service error",
			input:      `{"project_id":1, "access_token": {"project_id":1, "access_token_id":1, "name":"nm", "expired_at": 1, "last_updated_user_id":1}}`,
			wantStatus: http.StatusInternalServerError,
			mockErr:    errors.New("something wrong"),
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if c.mockResp != nil || c.mockErr != nil {
				iamMock.On("PutAccessToken").Return(c.mockResp, c.mockErr).Once()
			}
			rec := httptest.NewRecorder()
			req, _ := http.NewRequest(http.MethodPost, "/api/v1/iam/update-access-token/", strings.NewReader(c.input))
			req = req.WithContext(context.WithValue(req.Context(), userKey, &requestUser{userID: 1}))
			req.Header.Add("Content-Type", "application/json")
			svc.updateAccessTokenHandler(rec, req)
			if c.wantStatus != rec.Code {
				t.Fatalf("Unexpected HTTP status code: want=%+v, got=%+v", c.wantStatus, rec.Code)
			}
			resp := map[string]interface{}{}
			if err := json.NewDecoder(rec.Body).Decode(&resp); err != nil {
				t.Fatalf("Unexpected json decode error to response body: err=%+v", err)
			}
			// sappLogger.Info(resp)
			jsonKey := successJSONKey
			if c.wantStatus != http.StatusOK {
				jsonKey = errorJSONKey
			}
			if _, ok := resp[jsonKey]; !ok {
				t.Fatalf("Unexpected no response key: want key=%s", jsonKey)
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
func (m *mockIAMClient) ListAccessToken(context.Context, *iam.ListAccessTokenRequest, ...grpc.CallOption) (*iam.ListAccessTokenResponse, error) {
	args := m.Called()
	return args.Get(0).(*iam.ListAccessTokenResponse), args.Error(1)
}
func (m *mockIAMClient) AuthenticateAccessToken(context.Context, *iam.AuthenticateAccessTokenRequest, ...grpc.CallOption) (*iam.AuthenticateAccessTokenResponse, error) {
	args := m.Called()
	return args.Get(0).(*iam.AuthenticateAccessTokenResponse), args.Error(1)
}
func (m *mockIAMClient) PutAccessToken(context.Context, *iam.PutAccessTokenRequest, ...grpc.CallOption) (*iam.PutAccessTokenResponse, error) {
	args := m.Called()
	return args.Get(0).(*iam.PutAccessTokenResponse), args.Error(1)
}
func (m *mockIAMClient) DeleteAccessToken(context.Context, *iam.DeleteAccessTokenRequest, ...grpc.CallOption) (*empty.Empty, error) {
	args := m.Called()
	return args.Get(0).(*empty.Empty), args.Error(1)
}
func (m *mockIAMClient) AttachAccessTokenRole(context.Context, *iam.AttachAccessTokenRoleRequest, ...grpc.CallOption) (*iam.AttachAccessTokenRoleResponse, error) {
	args := m.Called()
	return args.Get(0).(*iam.AttachAccessTokenRoleResponse), args.Error(1)
}
func (m *mockIAMClient) DetachAccessTokenRole(context.Context, *iam.DetachAccessTokenRoleRequest, ...grpc.CallOption) (*empty.Empty, error) {
	args := m.Called()
	return args.Get(0).(*empty.Empty), args.Error(1)
}
func (m *mockIAMClient) IsAuthorized(context.Context, *iam.IsAuthorizedRequest, ...grpc.CallOption) (*iam.IsAuthorizedResponse, error) {
	args := m.Called()
	return args.Get(0).(*iam.IsAuthorizedResponse), args.Error(1)
}
func (m *mockIAMClient) IsAuthorizedAdmin(context.Context, *iam.IsAuthorizedAdminRequest, ...grpc.CallOption) (*iam.IsAuthorizedAdminResponse, error) {
	args := m.Called()
	return args.Get(0).(*iam.IsAuthorizedAdminResponse), args.Error(1)
}
func (m *mockIAMClient) IsAuthorizedToken(context.Context, *iam.IsAuthorizedTokenRequest, ...grpc.CallOption) (*iam.IsAuthorizedTokenResponse, error) {
	args := m.Called()
	return args.Get(0).(*iam.IsAuthorizedTokenResponse), args.Error(1)
}
func (m *mockIAMClient) IsAdmin(context.Context, *iam.IsAdminRequest, ...grpc.CallOption) (*iam.IsAdminResponse, error) {
	args := m.Called()
	return args.Get(0).(*iam.IsAdminResponse), args.Error(1)
}
func (m *mockIAMClient) AnalyzeTokenExpiration(context.Context, *empty.Empty, ...grpc.CallOption) (*empty.Empty, error) {
	args := m.Called()
	return args.Get(0).(*empty.Empty), args.Error(1)
}
