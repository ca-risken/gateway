package main

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/CyberAgent/mimosa-code/proto/code"
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

func TestListGitleaksHandler(t *testing.T) {
	mock := &mockCodeClient{}
	svc := gatewayService{
		codeClient: mock,
	}
	cases := []struct {
		name       string
		input      string
		mockResp   *code.ListGitleaksResponse
		mockErr    error
		wantStatus int
	}{
		{
			name:       "OK",
			input:      `project_id=1`,
			mockResp:   &code.ListGitleaksResponse{},
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
				mock.On("ListGitleaks").Return(c.mockResp, c.mockErr).Once()
			}
			rec := httptest.NewRecorder()
			req, _ := http.NewRequest(http.MethodGet, "/api/v1/code/list-gitleaks/?"+c.input, nil)
			svc.listGitleaksHandler(rec, req)
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

func TestPutGitleaksHandler(t *testing.T) {
	mock := &mockCodeClient{}
	svc := gatewayService{
		codeClient: mock,
	}
	cases := []struct {
		name       string
		input      string
		mockResp   *code.PutGitleaksResponse
		mockErr    error
		wantStatus int
	}{
		{
			name:       "OK",
			input:      `{"project_id":1, "gitleaks": {"gitleaks_id":1, "code_data_source_id":1, "name":"test", "project_id":1, "type":1, "target_resource":"gitleakstest"}}`,
			mockResp:   &code.PutGitleaksResponse{},
			wantStatus: http.StatusOK,
		},
		{
			name:       "NG Invalid parameter",
			input:      `invalid_param`,
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "NG Backend service error",
			input:      `{"project_id":1, "gitleaks": {"gitleaks_id":1, "code_data_source_id":1, "name":"test", "project_id":1, "type":1, "target_resource":"gitleakstest"}}`,
			wantStatus: http.StatusInternalServerError,
			mockErr:    errors.New("something wrong"),
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if c.mockResp != nil || c.mockErr != nil {
				mock.On("PutGitleaks").Return(c.mockResp, c.mockErr).Once()
			}
			rec := httptest.NewRecorder()
			req, _ := http.NewRequest(http.MethodPost, "/api/v1/code/put-gitleaks/", strings.NewReader(c.input))
			req.Header.Add("Content-Type", "application/json")
			svc.putGitleaksHandler(rec, req)
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
func TestDeleteGitleaksHandler(t *testing.T) {
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
			input:      `{"project_id":1, "gitleaks_id":1}`,
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
			input:      `{"project_id":1, "gitleaks_id":1}`,
			wantStatus: http.StatusInternalServerError,
			mockErr:    errors.New("something wrong"),
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if c.mockResp != nil || c.mockErr != nil {
				mock.On("DeleteGitleaks").Return(c.mockResp, c.mockErr).Once()
			}
			rec := httptest.NewRecorder()
			req, _ := http.NewRequest(http.MethodPost, "/api/v1/code/delete-gitleaks/", strings.NewReader(c.input))
			req.Header.Add("Content-Type", "application/json")
			svc.deleteGitleaksHandler(rec, req)
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
			input:      `{"project_id":1, "gitleaks_id":1}`,
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
			input:      `{"project_id":1, "gitleaks_id":1}`,
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
func (m *mockCodeClient) ListGitleaks(context.Context, *code.ListGitleaksRequest, ...grpc.CallOption) (*code.ListGitleaksResponse, error) {
	args := m.Called()
	return args.Get(0).(*code.ListGitleaksResponse), args.Error(1)
}
func (m *mockCodeClient) PutGitleaks(context.Context, *code.PutGitleaksRequest, ...grpc.CallOption) (*code.PutGitleaksResponse, error) {
	args := m.Called()
	return args.Get(0).(*code.PutGitleaksResponse), args.Error(1)
}
func (m *mockCodeClient) DeleteGitleaks(context.Context, *code.DeleteGitleaksRequest, ...grpc.CallOption) (*empty.Empty, error) {
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
