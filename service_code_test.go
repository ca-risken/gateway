package main

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/ca-risken/datasource-api/proto/code"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/stretchr/testify/mock"
	"google.golang.org/grpc"
)

const (
	length65 string = "12345678901234567890123456789012345678901234567890123456789012345"
)

func TestListCodeDataSourceHandler(t *testing.T) {
	mock := &mockCodeClient{}
	svc := gatewayService{
		codeClient: mock,
	}
	cases := []struct {
		name       string
		input      string
		mockResp   *code.ListDataSourceResponse
		mockErr    error
		wantStatus int
	}{
		{
			name:       "OK",
			input:      `project_id=1`,
			mockResp:   &code.ListDataSourceResponse{},
			wantStatus: http.StatusOK,
		},
		{
			name:       "NG Invalid parameter",
			input:      `name=` + length65,
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
				mock.On("ListDataSource").Return(c.mockResp, c.mockErr).Once()
			}
			rec := httptest.NewRecorder()
			req, _ := http.NewRequest(http.MethodGet, "/api/v1/code/list-datasource/?"+c.input, nil)
			svc.listCodeDataSourceHandler(rec, req)
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

func TestListGitHubSettingHandler(t *testing.T) {
	mock := &mockCodeClient{}
	svc := gatewayService{
		codeClient: mock,
	}
	cases := []struct {
		name       string
		input      string
		mockResp   *code.ListGitHubSettingResponse
		mockErr    error
		wantStatus int
	}{
		{
			name:       "OK",
			input:      `project_id=1`,
			mockResp:   &code.ListGitHubSettingResponse{},
			wantStatus: http.StatusOK,
		},
		{
			name:       "NG Invalid parameter",
			input:      "",
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
				mock.On("ListGitHubSetting").Return(c.mockResp, c.mockErr).Once()
			}
			rec := httptest.NewRecorder()
			req, _ := http.NewRequest(http.MethodGet, "/api/v1/code/list-github-setting/?"+c.input, nil)
			svc.listGitHubSettingHandler(rec, req)
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

func TestPutGitHubSettingHandler(t *testing.T) {
	mock := &mockCodeClient{}
	svc := gatewayService{
		codeClient: mock,
	}
	cases := []struct {
		name       string
		input      string
		mockResp   *code.PutGitHubSettingResponse
		mockErr    error
		wantStatus int
	}{
		{
			name:       "OK",
			input:      `{"project_id":1, "github_setting": {"github_setting_id":1, "name":"test", "project_id":1, "type":1, "target_resource":"githubsetting"}}`,
			mockResp:   &code.PutGitHubSettingResponse{},
			wantStatus: http.StatusOK,
		},
		{
			name:       "NG Invalid parameter",
			input:      `invalid_param`,
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "NG Backend service error",
			input:      `{"project_id":1, "github_setting": {"github_setting_id":1, "name":"test", "project_id":1, "type":1, "target_resource":"githubsetting"}}`,
			wantStatus: http.StatusInternalServerError,
			mockErr:    errors.New("something wrong"),
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if c.mockResp != nil || c.mockErr != nil {
				mock.On("PutGitHubSetting").Return(c.mockResp, c.mockErr).Once()
			}
			rec := httptest.NewRecorder()
			req, _ := http.NewRequest(http.MethodPost, "/api/v1/code/put-github-setting/", strings.NewReader(c.input))
			req.Header.Add("Content-Type", "application/json")
			svc.putGitHubSettingHandler(rec, req)
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
func TestDeleteGitHubSettingHandler(t *testing.T) {
	mock := &mockCodeClient{}
	svc := gatewayService{
		codeClient: mock,
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
			input:      `{"project_id":1, "github_setting_id":1}`,
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
			input:      `{"project_id":1, "github_setting_id":1}`,
			wantStatus: http.StatusInternalServerError,
			mockErr:    errors.New("something wrong"),
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if c.mockResp != nil || c.mockErr != nil {
				mock.On("DeleteGitHubSetting").Return(c.mockResp, c.mockErr).Once()
			}
			rec := httptest.NewRecorder()
			req, _ := http.NewRequest(http.MethodPost, "/api/v1/code/delete-github-setting/", strings.NewReader(c.input))
			req.Header.Add("Content-Type", "application/json")
			svc.deleteGitHubSettingHandler(rec, req)
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

func TestPutGitleaksSettingHandler(t *testing.T) {
	mock := &mockCodeClient{}
	svc := gatewayService{
		codeClient: mock,
	}
	cases := []struct {
		name       string
		input      string
		mockResp   *code.PutGitleaksSettingResponse
		mockErr    error
		wantStatus int
	}{
		{
			name:       "OK",
			input:      `{"project_id":1, "gitleaks_setting": {"github_setting_id":1, "code_data_source_id":1, "project_id":1, "type":1}}`,
			mockResp:   &code.PutGitleaksSettingResponse{},
			wantStatus: http.StatusOK,
		},
		{
			name:       "NG Invalid parameter",
			input:      `invalid_param`,
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "NG Backend service error",
			input:      `{"project_id":1, "gitleaks_setting": {"github_setting_id":1, "code_data_source_id":1, "project_id":1, "type":1}}`,
			wantStatus: http.StatusInternalServerError,
			mockErr:    errors.New("something wrong"),
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if c.mockResp != nil || c.mockErr != nil {
				mock.On("PutGitleaksSetting").Return(c.mockResp, c.mockErr).Once()
			}
			rec := httptest.NewRecorder()
			req, _ := http.NewRequest(http.MethodPost, "/api/v1/code/put-gitleaks-setting/", strings.NewReader(c.input))
			req.Header.Add("Content-Type", "application/json")
			svc.putGitleaksSettingHandler(rec, req)
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
func TestDeleteGitleaksSettingHandler(t *testing.T) {
	mock := &mockCodeClient{}
	svc := gatewayService{
		codeClient: mock,
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
			input:      `{"project_id":1, "github_setting_id":1}`,
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
			input:      `{"project_id":1, "github_setting_id":1}`,
			wantStatus: http.StatusInternalServerError,
			mockErr:    errors.New("something wrong"),
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if c.mockResp != nil || c.mockErr != nil {
				mock.On("DeleteGitleaksSetting").Return(c.mockResp, c.mockErr).Once()
			}
			rec := httptest.NewRecorder()
			req, _ := http.NewRequest(http.MethodPost, "/api/v1/code/delete-gitleaks-setting/", strings.NewReader(c.input))
			req.Header.Add("Content-Type", "application/json")
			svc.deleteGitleaksSettingHandler(rec, req)
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
func TestInvokeScanGitleaksHandler(t *testing.T) {
	mock := &mockCodeClient{}
	svc := gatewayService{
		codeClient: mock,
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
			input:      `{"project_id":1, "github_setting_id":1}`,
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
			input:      `{"project_id":1, "github_setting_id":1}`,
			wantStatus: http.StatusInternalServerError,
			mockErr:    errors.New("something wrong"),
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if c.mockResp != nil || c.mockErr != nil {
				mock.On("InvokeScanGitleaks").Return(c.mockResp, c.mockErr).Once()
			}
			rec := httptest.NewRecorder()
			req, _ := http.NewRequest(http.MethodPost, "/api/v1/code/invoke-scan-gitleaks/", strings.NewReader(c.input))
			req.Header.Add("Content-Type", "application/json")
			svc.invokeScanGitleaksHandler(rec, req)
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
type mockCodeClient struct {
	mock.Mock
}

func (m *mockCodeClient) ListDataSource(context.Context, *code.ListDataSourceRequest, ...grpc.CallOption) (*code.ListDataSourceResponse, error) {
	args := m.Called()
	return args.Get(0).(*code.ListDataSourceResponse), args.Error(1)
}
func (m *mockCodeClient) ListGitHubSetting(context.Context, *code.ListGitHubSettingRequest, ...grpc.CallOption) (*code.ListGitHubSettingResponse, error) {
	args := m.Called()
	return args.Get(0).(*code.ListGitHubSettingResponse), args.Error(1)
}
func (m *mockCodeClient) GetGitHubSetting(context.Context, *code.GetGitHubSettingRequest, ...grpc.CallOption) (*code.GetGitHubSettingResponse, error) {
	args := m.Called()
	return args.Get(0).(*code.GetGitHubSettingResponse), args.Error(1)
}
func (m *mockCodeClient) PutGitHubSetting(context.Context, *code.PutGitHubSettingRequest, ...grpc.CallOption) (*code.PutGitHubSettingResponse, error) {
	args := m.Called()
	return args.Get(0).(*code.PutGitHubSettingResponse), args.Error(1)
}
func (m *mockCodeClient) DeleteGitHubSetting(context.Context, *code.DeleteGitHubSettingRequest, ...grpc.CallOption) (*empty.Empty, error) {
	args := m.Called()
	return args.Get(0).(*empty.Empty), args.Error(1)
}
func (m *mockCodeClient) PutGitleaksSetting(context.Context, *code.PutGitleaksSettingRequest, ...grpc.CallOption) (*code.PutGitleaksSettingResponse, error) {
	args := m.Called()
	return args.Get(0).(*code.PutGitleaksSettingResponse), args.Error(1)
}
func (m *mockCodeClient) DeleteGitleaksSetting(context.Context, *code.DeleteGitleaksSettingRequest, ...grpc.CallOption) (*empty.Empty, error) {
	args := m.Called()
	return args.Get(0).(*empty.Empty), args.Error(1)
}
func (m *mockCodeClient) InvokeScanGitleaks(context.Context, *code.InvokeScanGitleaksRequest, ...grpc.CallOption) (*empty.Empty, error) {
	args := m.Called()
	return args.Get(0).(*empty.Empty), args.Error(1)
}
func (m *mockCodeClient) InvokeScanAllGitleaks(context.Context, *empty.Empty, ...grpc.CallOption) (*empty.Empty, error) {
	args := m.Called()
	return args.Get(0).(*empty.Empty), args.Error(1)
}
func (m *mockCodeClient) ListGitHubEnterpriseOrg(ctx context.Context, in *code.ListGitHubEnterpriseOrgRequest, opts ...grpc.CallOption) (*code.ListGitHubEnterpriseOrgResponse, error) {
	args := m.Called()
	return args.Get(0).(*code.ListGitHubEnterpriseOrgResponse), args.Error(1)
}
func (m *mockCodeClient) PutGitHubEnterpriseOrg(ctx context.Context, in *code.PutGitHubEnterpriseOrgRequest, opts ...grpc.CallOption) (*code.PutGitHubEnterpriseOrgResponse, error) {
	args := m.Called()
	return args.Get(0).(*code.PutGitHubEnterpriseOrgResponse), args.Error(1)
}
func (m *mockCodeClient) DeleteGitHubEnterpriseOrg(ctx context.Context, in *code.DeleteGitHubEnterpriseOrgRequest, opts ...grpc.CallOption) (*empty.Empty, error) {
	args := m.Called()
	return args.Get(0).(*empty.Empty), args.Error(1)
}
