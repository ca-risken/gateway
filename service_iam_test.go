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

func TestListUserHandler(t *testing.T) {
	iamMock := &mockIAMClient{}
	svc := gatewayService{
		iamClient: iamMock,
	}
	cases := []struct {
		name       string
		input      string
		mockResp   *iam.ListUserResponse
		mockErr    error
		wantStatus int
	}{
		{
			name:       "OK",
			input:      `project_id=1`,
			mockResp:   &iam.ListUserResponse{UserId: []uint32{1, 2, 3}},
			wantStatus: http.StatusOK,
		},
		{
			name:       "NG Invalid parameter",
			input:      `name=12345678901234567890123456789012345678901234567890123456789012345`,
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "NG Backend service error",
			input:      `project_id=1`,
			wantStatus: http.StatusInternalServerError,
			mockErr:    errors.New("something wrong"),
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if c.mockResp != nil || c.mockErr != nil {
				iamMock.On("ListUser").Return(c.mockResp, c.mockErr).Once()
			}
			// Invoke HTTP Request
			rec := httptest.NewRecorder()
			req, _ := http.NewRequest(http.MethodGet, "/api/v1/iam/list-user?"+c.input, nil)
			svc.listUserHandler(rec, req)
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

func TestGetUserHandler(t *testing.T) {
	iamMock := &mockIAMClient{}
	svc := gatewayService{
		iamClient: iamMock,
	}
	cases := []struct {
		name       string
		input      string
		mockResp   *iam.GetUserResponse
		mockErr    error
		wantStatus int
	}{
		{
			name:       "OK",
			input:      `user_id=1`,
			mockResp:   &iam.GetUserResponse{},
			wantStatus: http.StatusOK,
		},
		{
			name:       "NG Invalid parameter",
			input:      `no_param`,
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "NG Backend service error",
			input:      `user_id=1`,
			wantStatus: http.StatusInternalServerError,
			mockErr:    errors.New("something wrong"),
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if c.mockResp != nil || c.mockErr != nil {
				iamMock.On("GetUser").Return(c.mockResp, c.mockErr).Once()
			}
			// Invoke HTTP Request
			rec := httptest.NewRecorder()
			req, _ := http.NewRequest(http.MethodGet, "/api/v1/iam/get-user?"+c.input, nil)
			svc.getUserHandler(rec, req)
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

func TestIsAdminHandler(t *testing.T) {
	iamMock := &mockIAMClient{}
	svc := gatewayService{
		iamClient: iamMock,
	}
	cases := []struct {
		name       string
		input      string
		mockResp   *iam.IsAdminResponse
		mockErr    error
		wantStatus int
	}{
		{
			name:       "OK",
			input:      `user_id=1`,
			mockResp:   &iam.IsAdminResponse{},
			wantStatus: http.StatusOK,
		},
		{
			name:       "NG Invalid parameter",
			input:      `no_param`,
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "NG Backend service error",
			input:      `user_id=1`,
			wantStatus: http.StatusInternalServerError,
			mockErr:    errors.New("something wrong"),
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if c.mockResp != nil || c.mockErr != nil {
				iamMock.On("IsAdmin").Return(c.mockResp, c.mockErr).Once()
			}
			rec := httptest.NewRecorder()
			req, _ := http.NewRequest(http.MethodGet, "/api/v1/iam/is-admin?"+c.input, nil)
			svc.isAdminHandler(rec, req)
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

func TestListRoleHandler(t *testing.T) {
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
			input:      `project_id=1`,
			mockResp:   &iam.ListRoleResponse{},
			wantStatus: http.StatusOK,
		},
		{
			name:       "NG Invalid parameter",
			input:      `no_param`,
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "NG Backend service error",
			input:      `project_id=1`,
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
			req, _ := http.NewRequest(http.MethodGet, "/api/v1/iam/list-role?"+c.input, nil)
			svc.listRoleHandler(rec, req)
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

func TestGetRoleHandler(t *testing.T) {
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
			input:      `project_id=1&role_id=1`,
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
			input:      `project_id=1&role_id=1`,
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
			req, _ := http.NewRequest(http.MethodGet, "/api/v1/iam/get-role?"+c.input, nil)
			svc.getRoleHandler(rec, req)
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

func TestPutRoleHandler(t *testing.T) {
	iamMock := &mockIAMClient{}
	svc := gatewayService{
		iamClient: iamMock,
	}
	cases := []struct {
		name       string
		input      string
		mockResp   *iam.PutRoleResponse
		mockErr    error
		wantStatus int
	}{
		{
			name:       "OK",
			input:      `{"project_id":1, "role":{"project_id":1, "name":"nm"}}`,
			mockResp:   &iam.PutRoleResponse{},
			wantStatus: http.StatusOK,
		},
		{
			name:       "NG Invalid parameter",
			input:      `invalid_param`,
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "NG Backend service error",
			input:      `{"project_id":1, "role":{"project_id":1, "name":"nm"}}`,
			wantStatus: http.StatusInternalServerError,
			mockErr:    errors.New("something wrong"),
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if c.mockResp != nil || c.mockErr != nil {
				iamMock.On("PutRole").Return(c.mockResp, c.mockErr).Once()
			}
			// Invoke HTTP Request
			rec := httptest.NewRecorder()
			req, _ := http.NewRequest(http.MethodPost, "/api/v1/iam/put-role/", strings.NewReader(c.input))
			req.Header.Add("Content-Type", "application/json")
			svc.putRoleHandler(rec, req)
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

func TestDeleteRoleHandler(t *testing.T) {
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
			input:      `{"project_id":1, "role_id":1}`,
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
			input:      `{"project_id":1, "role_id":1}`,
			wantStatus: http.StatusInternalServerError,
			mockErr:    errors.New("something wrong"),
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if c.mockResp != nil || c.mockErr != nil {
				iamMock.On("DeleteRole").Return(c.mockResp, c.mockErr).Once()
			}
			// Invoke HTTP Request
			rec := httptest.NewRecorder()
			req, _ := http.NewRequest(http.MethodPost, "/api/v1/iam/delete-role/", strings.NewReader(c.input))
			req.Header.Add("Content-Type", "application/json")
			svc.deleteRoleHandler(rec, req)
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

func TestAttachRoleHandler(t *testing.T) {
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
			input:      `{"project_id":1, "role_id":1, "user_id":1}`,
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
			input:      `{"project_id":1, "role_id":1, "user_id":1}`,
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
			req, _ := http.NewRequest(http.MethodPost, "/api/v1/iam/attach-role/", strings.NewReader(c.input))
			req.Header.Add("Content-Type", "application/json")
			svc.attachRoleHandler(rec, req)
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

func TestDetachRoleHandler(t *testing.T) {
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
			input:      `{"project_id":1, "role_id":1, "user_id":1}`,
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
			input:      `{"project_id":1, "role_id":1, "user_id":1}`,
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
			req, _ := http.NewRequest(http.MethodPost, "/api/v1/iam/detach-role/", strings.NewReader(c.input))
			req.Header.Add("Content-Type", "application/json")
			svc.detachRoleHandler(rec, req)
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

func TestListPolicyHandler(t *testing.T) {
	iamMock := &mockIAMClient{}
	svc := gatewayService{
		iamClient: iamMock,
	}
	cases := []struct {
		name       string
		input      string
		mockResp   *iam.ListPolicyResponse
		mockErr    error
		wantStatus int
	}{
		{
			name:       "OK",
			input:      `project_id=1`,
			mockResp:   &iam.ListPolicyResponse{},
			wantStatus: http.StatusOK,
		},
		{
			name:       "NG Invalid parameter",
			input:      `no_param`,
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "NG Backend service error",
			input:      `project_id=1`,
			wantStatus: http.StatusInternalServerError,
			mockErr:    errors.New("something wrong"),
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if c.mockResp != nil || c.mockErr != nil {
				iamMock.On("ListPolicy").Return(c.mockResp, c.mockErr).Once()
			}
			// Invoke HTTP Request
			rec := httptest.NewRecorder()
			req, _ := http.NewRequest(http.MethodGet, "/api/v1/iam/list-policy?"+c.input, nil)
			svc.listPolicyHandler(rec, req)
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

func TestGetPolicyHandler(t *testing.T) {
	iamMock := &mockIAMClient{}
	svc := gatewayService{
		iamClient: iamMock,
	}
	cases := []struct {
		name       string
		input      string
		mockResp   *iam.GetPolicyResponse
		mockErr    error
		wantStatus int
	}{
		{
			name:       "OK",
			input:      `project_id=1&policy_id=1`,
			mockResp:   &iam.GetPolicyResponse{},
			wantStatus: http.StatusOK,
		},
		{
			name:       "NG Invalid parameter",
			input:      `no_param`,
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "NG Backend service error",
			input:      `project_id=1&policy_id=1`,
			wantStatus: http.StatusInternalServerError,
			mockErr:    errors.New("something wrong"),
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if c.mockResp != nil || c.mockErr != nil {
				iamMock.On("GetPolicy").Return(c.mockResp, c.mockErr).Once()
			}
			// Invoke HTTP Request
			rec := httptest.NewRecorder()
			req, _ := http.NewRequest(http.MethodGet, "/api/v1/iam/get-policy?"+c.input, nil)
			svc.getPolicyHandler(rec, req)
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

func TestPutPolicyHandler(t *testing.T) {
	iamMock := &mockIAMClient{}
	svc := gatewayService{
		iamClient: iamMock,
	}
	cases := []struct {
		name       string
		input      string
		mockResp   *iam.PutPolicyResponse
		mockErr    error
		wantStatus int
	}{
		{
			name:       "OK",
			input:      `{"project_id":1, "policy":{"project_id":1, "name":"nm", "action_ptn":".*", "resource_ptn":".*"}}`,
			mockResp:   &iam.PutPolicyResponse{},
			wantStatus: http.StatusOK,
		},
		{
			name:       "NG Invalid parameter",
			input:      `invalid_param`,
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "NG Backend service error",
			input:      `{"project_id":1, "policy":{"project_id":1, "name":"nm", "action_ptn":".*", "resource_ptn":".*"}}`,
			wantStatus: http.StatusInternalServerError,
			mockErr:    errors.New("something wrong"),
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if c.mockResp != nil || c.mockErr != nil {
				iamMock.On("PutPolicy").Return(c.mockResp, c.mockErr).Once()
			}
			// Invoke HTTP Request
			rec := httptest.NewRecorder()
			req, _ := http.NewRequest(http.MethodPost, "/api/v1/iam/put-policy/", strings.NewReader(c.input))
			req.Header.Add("Content-Type", "application/json")
			svc.putPolicyHandler(rec, req)
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

func TestDeletePolicyHandler(t *testing.T) {
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
			input:      `{"project_id":1, "policy_id":1}`,
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
			input:      `{"project_id":1, "policy_id":1}`,
			wantStatus: http.StatusInternalServerError,
			mockErr:    errors.New("something wrong"),
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if c.mockResp != nil || c.mockErr != nil {
				iamMock.On("DeletePolicy").Return(c.mockResp, c.mockErr).Once()
			}
			// Invoke HTTP Request
			rec := httptest.NewRecorder()
			req, _ := http.NewRequest(http.MethodPost, "/api/v1/iam/delete-policy/", strings.NewReader(c.input))
			req.Header.Add("Content-Type", "application/json")
			svc.deletePolicyHandler(rec, req)
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

func TestAttachPolicyHandler(t *testing.T) {
	iamMock := &mockIAMClient{}
	svc := gatewayService{
		iamClient: iamMock,
	}
	cases := []struct {
		name       string
		input      string
		mockResp   *iam.AttachPolicyResponse
		mockErr    error
		wantStatus int
	}{
		{
			name:       "OK",
			input:      `{"project_id":1, "role_id":1, "policy_id":1}`,
			mockResp:   &iam.AttachPolicyResponse{},
			wantStatus: http.StatusOK,
		},
		{
			name:       "NG Invalid parameter",
			input:      `invalid_param`,
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "NG Backend service error",
			input:      `{"project_id":1, "role_id":1, "policy_id":1}`,
			wantStatus: http.StatusInternalServerError,
			mockErr:    errors.New("something wrong"),
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if c.mockResp != nil || c.mockErr != nil {
				iamMock.On("AttachPolicy").Return(c.mockResp, c.mockErr).Once()
			}
			// Invoke HTTP Request
			rec := httptest.NewRecorder()
			req, _ := http.NewRequest(http.MethodPost, "/api/v1/iam/attach-policy/", strings.NewReader(c.input))
			req.Header.Add("Content-Type", "application/json")
			svc.attachPolicyHandler(rec, req)
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

func TestDetachPolicyHandler(t *testing.T) {
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
			input:      `{"project_id":1, "role_id":1, "policy_id":1}`,
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
			input:      `{"project_id":1, "role_id":1, "policy_id":1}`,
			wantStatus: http.StatusInternalServerError,
			mockErr:    errors.New("something wrong"),
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if c.mockResp != nil || c.mockErr != nil {
				iamMock.On("DetachPolicy").Return(c.mockResp, c.mockErr).Once()
			}
			// Invoke HTTP Request
			rec := httptest.NewRecorder()
			req, _ := http.NewRequest(http.MethodPost, "/api/v1/iam/detach-policy/", strings.NewReader(c.input))
			req.Header.Add("Content-Type", "application/json")
			svc.detachPolicyHandler(rec, req)
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

func TestListAccessTokenHandler(t *testing.T) {
	iamMock := &mockIAMClient{}
	svc := gatewayService{
		iamClient: iamMock,
	}
	cases := []struct {
		name       string
		input      string
		mockResp   *iam.ListAccessTokenResponse
		mockErr    error
		wantStatus int
	}{
		{
			name:       "OK",
			input:      `project_id=1`,
			mockResp:   &iam.ListAccessTokenResponse{},
			wantStatus: http.StatusOK,
		},
		{
			name:       "NG Invalid parameter",
			input:      `no_param`,
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "NG Backend service error",
			input:      `project_id=1`,
			wantStatus: http.StatusInternalServerError,
			mockErr:    errors.New("something wrong"),
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if c.mockResp != nil || c.mockErr != nil {
				iamMock.On("ListAccessToken").Return(c.mockResp, c.mockErr).Once()
			}
			rec := httptest.NewRecorder()
			req, _ := http.NewRequest(http.MethodGet, "/api/v1/iam/list-access-token?"+c.input, nil)
			svc.listAccessTokenHandler(rec, req)
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

func TestDeleteAccessTokenHandler(t *testing.T) {
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
			input:      `{"project_id":1, "access_token_id":1}`,
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
			input:      `{"project_id":1, "access_token_id":1}`,
			wantStatus: http.StatusInternalServerError,
			mockErr:    errors.New("something wrong"),
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if c.mockResp != nil || c.mockErr != nil {
				iamMock.On("DeleteAccessToken").Return(c.mockResp, c.mockErr).Once()
			}
			// Invoke HTTP Request
			rec := httptest.NewRecorder()
			req, _ := http.NewRequest(http.MethodPost, "/api/v1/iam/delete-access-token/", strings.NewReader(c.input))
			req.Header.Add("Content-Type", "application/json")
			svc.deleteAccessTokenHandler(rec, req)
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

func TestAttachAccessTokenHandler(t *testing.T) {
	iamMock := &mockIAMClient{}
	svc := gatewayService{
		iamClient: iamMock,
	}
	cases := []struct {
		name       string
		input      string
		mockResp   *iam.AttachAccessTokenRoleResponse
		mockErr    error
		wantStatus int
	}{
		{
			name:       "OK",
			input:      `{"project_id":1, "role_id":1, "access_token_id":1}`,
			mockResp:   &iam.AttachAccessTokenRoleResponse{},
			wantStatus: http.StatusOK,
		},
		{
			name:       "NG Invalid parameter",
			input:      `invalid_param`,
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "NG Backend service error",
			input:      `{"project_id":1, "role_id":1, "access_token_id":1}`,
			wantStatus: http.StatusInternalServerError,
			mockErr:    errors.New("something wrong"),
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if c.mockResp != nil || c.mockErr != nil {
				iamMock.On("AttachAccessTokenRole").Return(c.mockResp, c.mockErr).Once()
			}
			// Invoke HTTP Request
			rec := httptest.NewRecorder()
			req, _ := http.NewRequest(http.MethodPost, "/api/v1/iam/attach-access-token/", strings.NewReader(c.input))
			req.Header.Add("Content-Type", "application/json")
			svc.attachAccessTokenRoleHandler(rec, req)
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

func TestDetachAccessTokenRoleHandler(t *testing.T) {
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
			input:      `{"project_id":1, "role_id":1, "access_token_id":1}`,
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
			input:      `{"project_id":1, "role_id":1, "access_token_id":1}`,
			wantStatus: http.StatusInternalServerError,
			mockErr:    errors.New("something wrong"),
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if c.mockResp != nil || c.mockErr != nil {
				iamMock.On("DetachAccessTokenRole").Return(c.mockResp, c.mockErr).Once()
			}
			// Invoke HTTP Request
			rec := httptest.NewRecorder()
			req, _ := http.NewRequest(http.MethodPost, "/api/v1/iam/detach-access-token/", strings.NewReader(c.input))
			req.Header.Add("Content-Type", "application/json")
			svc.detachAccessTokenRoleHandler(rec, req)
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
