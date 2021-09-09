package main

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/ca-risken/google/proto/google"
	"github.com/stretchr/testify/mock"
	"google.golang.org/grpc"
)

func TestListGoogleDataSourceHandler(t *testing.T) {
	mock := &mockGoogleClient{}
	svc := gatewayService{
		googleClient: mock,
	}
	cases := []struct {
		name       string
		input      string
		mockResp   *google.ListGoogleDataSourceResponse
		mockErr    error
		wantStatus int
	}{
		{
			name:       "OK",
			input:      `google_data_source_id=1&name=asset`,
			mockResp:   &google.ListGoogleDataSourceResponse{},
			wantStatus: http.StatusOK,
		},
		{
			name:       "NG Invalid parameter",
			input:      `name=` + length65,
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "NG Backend service error",
			input:      `google_data_source_id=1`,
			wantStatus: http.StatusInternalServerError,
			mockErr:    errors.New("something wrong"),
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if c.mockResp != nil || c.mockErr != nil {
				mock.On("ListGoogleDataSource").Return(c.mockResp, c.mockErr).Once()
			}
			rec := httptest.NewRecorder()
			req, _ := http.NewRequest(http.MethodGet, "/api/v1/google/list-google-datasource/?"+c.input, nil)
			svc.listGoogleDataSourceHandler(rec, req)
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

func TestListGCPHandler(t *testing.T) {
	mock := &mockGoogleClient{}
	svc := gatewayService{
		googleClient: mock,
	}
	cases := []struct {
		name       string
		input      string
		mockResp   *google.ListGCPResponse
		mockErr    error
		wantStatus int
	}{
		{
			name:       "OK",
			input:      `project_id=1&gcp_id=1&gcp_project_id=my-pj`,
			mockResp:   &google.ListGCPResponse{},
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
				mock.On("ListGCP").Return(c.mockResp, c.mockErr).Once()
			}
			rec := httptest.NewRecorder()
			req, _ := http.NewRequest(http.MethodGet, "/api/v1/google/list-gcp/?"+c.input, nil)
			svc.listGCPHandler(rec, req)
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

func TestGetGCPHandler(t *testing.T) {
	mock := &mockGoogleClient{}
	svc := gatewayService{
		googleClient: mock,
	}
	cases := []struct {
		name       string
		input      string
		mockResp   *google.GetGCPResponse
		mockErr    error
		wantStatus int
	}{
		{
			name:       "OK",
			input:      `project_id=1&gcp_id=1`,
			mockResp:   &google.GetGCPResponse{},
			wantStatus: http.StatusOK,
		},
		{
			name:       "NG Invalid parameter",
			input:      "",
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "NG Backend service error",
			input:      `project_id=1&gcp_id=1`,
			wantStatus: http.StatusInternalServerError,
			mockErr:    errors.New("something wrong"),
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if c.mockResp != nil || c.mockErr != nil {
				mock.On("GetGCP").Return(c.mockResp, c.mockErr).Once()
			}
			rec := httptest.NewRecorder()
			req, _ := http.NewRequest(http.MethodGet, "/api/v1/google/get-gcp/?"+c.input, nil)
			svc.getGCPHandler(rec, req)
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

func TestPutGCPHandler(t *testing.T) {
	mock := &mockGoogleClient{}
	svc := gatewayService{
		googleClient: mock,
	}
	cases := []struct {
		name       string
		input      string
		mockResp   *google.PutGCPResponse
		mockErr    error
		wantStatus int
	}{
		{
			name:       "OK",
			input:      `{"project_id":1, "gcp": {"gcp_id":1, "name":"test", "project_id":1, "gcp_project_id":"my-pj", "verification_code":"valid code"}}`,
			mockResp:   &google.PutGCPResponse{},
			wantStatus: http.StatusOK,
		},
		{
			name:       "NG Invalid parameter",
			input:      `invalid_param`,
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "NG Backend service error",
			input:      `{"project_id":1, "gcp": {"gcp_id":1, "name":"test", "project_id":1, "gcp_project_id":"my-pj", "verification_code":"valid code"}}`,
			wantStatus: http.StatusInternalServerError,
			mockErr:    errors.New("something wrong"),
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if c.mockResp != nil || c.mockErr != nil {
				mock.On("PutGCP").Return(c.mockResp, c.mockErr).Once()
			}
			rec := httptest.NewRecorder()
			req, _ := http.NewRequest(http.MethodPost, "/api/v1/google/put-gcp/", strings.NewReader(c.input))
			req.Header.Add("Content-Type", "application/json")
			svc.putGCPHandler(rec, req)
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

func TestDeleteGCPHandler(t *testing.T) {
	mock := &mockGoogleClient{}
	svc := gatewayService{
		googleClient: mock,
	}
	cases := []struct {
		name       string
		input      string
		mockResp   *google.Empty
		mockErr    error
		wantStatus int
	}{
		{
			name:       "OK",
			input:      `{"project_id":1, "gcp_id":1}`,
			mockResp:   &google.Empty{},
			wantStatus: http.StatusOK,
		},
		{
			name:       "NG Invalid parameter",
			input:      `invalid_param`,
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "NG Backend service error",
			input:      `{"project_id":1, "gcp_id":1}`,
			wantStatus: http.StatusInternalServerError,
			mockErr:    errors.New("something wrong"),
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if c.mockResp != nil || c.mockErr != nil {
				mock.On("DeleteGCP").Return(c.mockResp, c.mockErr).Once()
			}
			rec := httptest.NewRecorder()
			req, _ := http.NewRequest(http.MethodPost, "/api/v1/google/delete-gcp/", strings.NewReader(c.input))
			req.Header.Add("Content-Type", "application/json")
			svc.deleteGCPHandler(rec, req)
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

func TestListGCPDataSourceHandler(t *testing.T) {
	mock := &mockGoogleClient{}
	svc := gatewayService{
		googleClient: mock,
	}
	cases := []struct {
		name       string
		input      string
		mockResp   *google.ListGCPDataSourceResponse
		mockErr    error
		wantStatus int
	}{
		{
			name:       "OK",
			input:      `project_id=1&gcp_id=1`,
			mockResp:   &google.ListGCPDataSourceResponse{},
			wantStatus: http.StatusOK,
		},
		{
			name:       "NG Invalid parameter",
			input:      "",
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "NG Backend service error",
			input:      `project_id=1&gcp_id=1`,
			wantStatus: http.StatusInternalServerError,
			mockErr:    errors.New("something wrong"),
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if c.mockResp != nil || c.mockErr != nil {
				mock.On("ListGCPDataSource").Return(c.mockResp, c.mockErr).Once()
			}
			rec := httptest.NewRecorder()
			req, _ := http.NewRequest(http.MethodGet, "/api/v1/google/list-gcp-datasource/?"+c.input, nil)
			svc.listGCPDataSourceHandler(rec, req)
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

func TestGetGCPDataSourceHandler(t *testing.T) {
	mock := &mockGoogleClient{}
	svc := gatewayService{
		googleClient: mock,
	}
	cases := []struct {
		name       string
		input      string
		mockResp   *google.GetGCPDataSourceResponse
		mockErr    error
		wantStatus int
	}{
		{
			name:       "OK",
			input:      `project_id=1&gcp_id=1&google_data_source_id=1`,
			mockResp:   &google.GetGCPDataSourceResponse{},
			wantStatus: http.StatusOK,
		},
		{
			name:       "NG Invalid parameter",
			input:      "",
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "NG Backend service error",
			input:      `project_id=1&gcp_id=1&google_data_source_id=1`,
			wantStatus: http.StatusInternalServerError,
			mockErr:    errors.New("something wrong"),
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if c.mockResp != nil || c.mockErr != nil {
				mock.On("GetGCPDataSource").Return(c.mockResp, c.mockErr).Once()
			}
			rec := httptest.NewRecorder()
			req, _ := http.NewRequest(http.MethodGet, "/api/v1/google/get-gcp-datasource/?"+c.input, nil)
			svc.getGCPDataSourceHandler(rec, req)
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

func TestAttachGCPDataSourceHandler(t *testing.T) {
	mock := &mockGoogleClient{}
	svc := gatewayService{
		googleClient: mock,
	}
	cases := []struct {
		name       string
		input      string
		mockResp   *google.AttachGCPDataSourceResponse
		mockErr    error
		wantStatus int
	}{
		{
			name:       "OK",
			input:      `{"project_id":1, "gcp_data_source": {"gcp_id":1, "google_data_source_id":1, "project_id":1}}`,
			mockResp:   &google.AttachGCPDataSourceResponse{},
			wantStatus: http.StatusOK,
		},
		{
			name:       "NG Invalid parameter",
			input:      `invalid_param`,
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "NG Backend service error",
			input:      `{"project_id":1, "gcp_data_source": {"gcp_id":1, "google_data_source_id":1, "project_id":1}}`,
			wantStatus: http.StatusInternalServerError,
			mockErr:    errors.New("something wrong"),
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if c.mockResp != nil || c.mockErr != nil {
				mock.On("AttachGCPDataSource").Return(c.mockResp, c.mockErr).Once()
			}
			rec := httptest.NewRecorder()
			req, _ := http.NewRequest(http.MethodPost, "/api/v1/google/attach-gcp-datasource/", strings.NewReader(c.input))
			req.Header.Add("Content-Type", "application/json")
			svc.attachGCPDataSourceHandler(rec, req)
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

func TestDetachGCPDataSourceHandler(t *testing.T) {
	mock := &mockGoogleClient{}
	svc := gatewayService{
		googleClient: mock,
	}
	cases := []struct {
		name       string
		input      string
		mockResp   *google.Empty
		mockErr    error
		wantStatus int
	}{
		{
			name:       "OK",
			input:      `{"project_id":1, "gcp_id":1, "google_data_source_id":1}`,
			mockResp:   &google.Empty{},
			wantStatus: http.StatusOK,
		},
		{
			name:       "NG Invalid parameter",
			input:      `invalid_param`,
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "NG Backend service error",
			input:      `{"project_id":1, "gcp_id":1, "google_data_source_id":1}`,
			wantStatus: http.StatusInternalServerError,
			mockErr:    errors.New("something wrong"),
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if c.mockResp != nil || c.mockErr != nil {
				mock.On("DetachGCPDataSource").Return(c.mockResp, c.mockErr).Once()
			}
			rec := httptest.NewRecorder()
			req, _ := http.NewRequest(http.MethodPost, "/api/v1/google/detach-gcp-datasource/", strings.NewReader(c.input))
			req.Header.Add("Content-Type", "application/json")
			svc.detachGCPDataSourceHandler(rec, req)
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
func TestInvokeScanGCPHandler(t *testing.T) {
	mock := &mockGoogleClient{}
	svc := gatewayService{
		googleClient: mock,
	}
	cases := []struct {
		name       string
		input      string
		mockResp   *google.Empty
		mockErr    error
		wantStatus int
	}{
		{
			name:       "OK",
			input:      `{"project_id":1, "gcp_id":1, "google_data_source_id":1}`,
			mockResp:   &google.Empty{},
			wantStatus: http.StatusOK,
		},
		{
			name:       "NG Invalid parameter",
			input:      `invalid_param`,
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "NG Backend service error",
			input:      `{"project_id":1, "gcp_id":1, "google_data_source_id":1}`,
			wantStatus: http.StatusInternalServerError,
			mockErr:    errors.New("something wrong"),
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if c.mockResp != nil || c.mockErr != nil {
				mock.On("InvokeScanGCP").Return(c.mockResp, c.mockErr).Once()
			}
			rec := httptest.NewRecorder()
			req, _ := http.NewRequest(http.MethodPost, "/api/v1/google/invoke-scan-gcp/", strings.NewReader(c.input))
			req.Header.Add("Content-Type", "application/json")
			svc.invokeScanGCPHandler(rec, req)
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
type mockGoogleClient struct {
	mock.Mock
}

func (m *mockGoogleClient) ListGoogleDataSource(context.Context, *google.ListGoogleDataSourceRequest, ...grpc.CallOption) (*google.ListGoogleDataSourceResponse, error) {
	args := m.Called()
	return args.Get(0).(*google.ListGoogleDataSourceResponse), args.Error(1)
}
func (m *mockGoogleClient) ListGCP(context.Context, *google.ListGCPRequest, ...grpc.CallOption) (*google.ListGCPResponse, error) {
	args := m.Called()
	return args.Get(0).(*google.ListGCPResponse), args.Error(1)
}
func (m *mockGoogleClient) GetGCP(context.Context, *google.GetGCPRequest, ...grpc.CallOption) (*google.GetGCPResponse, error) {
	args := m.Called()
	return args.Get(0).(*google.GetGCPResponse), args.Error(1)
}
func (m *mockGoogleClient) PutGCP(context.Context, *google.PutGCPRequest, ...grpc.CallOption) (*google.PutGCPResponse, error) {
	args := m.Called()
	return args.Get(0).(*google.PutGCPResponse), args.Error(1)
}
func (m *mockGoogleClient) DeleteGCP(context.Context, *google.DeleteGCPRequest, ...grpc.CallOption) (*google.Empty, error) {
	args := m.Called()
	return args.Get(0).(*google.Empty), args.Error(1)
}
func (m *mockGoogleClient) ListGCPDataSource(context.Context, *google.ListGCPDataSourceRequest, ...grpc.CallOption) (*google.ListGCPDataSourceResponse, error) {
	args := m.Called()
	return args.Get(0).(*google.ListGCPDataSourceResponse), args.Error(1)
}
func (m *mockGoogleClient) GetGCPDataSource(context.Context, *google.GetGCPDataSourceRequest, ...grpc.CallOption) (*google.GetGCPDataSourceResponse, error) {
	args := m.Called()
	return args.Get(0).(*google.GetGCPDataSourceResponse), args.Error(1)
}
func (m *mockGoogleClient) AttachGCPDataSource(context.Context, *google.AttachGCPDataSourceRequest, ...grpc.CallOption) (*google.AttachGCPDataSourceResponse, error) {
	args := m.Called()
	return args.Get(0).(*google.AttachGCPDataSourceResponse), args.Error(1)
}
func (m *mockGoogleClient) DetachGCPDataSource(context.Context, *google.DetachGCPDataSourceRequest, ...grpc.CallOption) (*google.Empty, error) {
	args := m.Called()
	return args.Get(0).(*google.Empty), args.Error(1)
}

func (m *mockGoogleClient) InvokeScanGCP(context.Context, *google.InvokeScanGCPRequest, ...grpc.CallOption) (*google.Empty, error) {
	args := m.Called()
	return args.Get(0).(*google.Empty), args.Error(1)
}
func (m *mockGoogleClient) InvokeScanAll(context.Context, *google.Empty, ...grpc.CallOption) (*google.Empty, error) {
	args := m.Called()
	return args.Get(0).(*google.Empty), args.Error(1)
}
