package main

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/CyberAgent/mimosa-osint/proto/osint"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/stretchr/testify/mock"
	"google.golang.org/grpc"
)

func TestListOsintHandler(t *testing.T) {
	osintMock := &mockOsintClient{}
	svc := gatewayService{
		osintClient: osintMock,
	}
	cases := []struct {
		name       string
		input      string
		mockResp   *osint.ListOsintResponse
		mockErr    error
		wantStatus int
	}{
		{
			name:       "OK",
			input:      `project_id=1`,
			mockResp:   &osint.ListOsintResponse{},
			wantStatus: http.StatusOK,
		},
		{
			name:       "NG Invalid parameter",
			input:      `resource_name=123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901`,
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
				osintMock.On("ListOsint").Return(c.mockResp, c.mockErr).Once()
			}
			// Invoke HTTP Request
			rec := httptest.NewRecorder()
			req, _ := http.NewRequest(http.MethodGet, "/api/v1/osint/list-osint/?"+c.input, nil)
			svc.listOsintHandler(rec, req)
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

func TestGetOsintHandler(t *testing.T) {
	osintMock := &mockOsintClient{}
	svc := gatewayService{
		osintClient: osintMock,
	}
	cases := []struct {
		name       string
		input      string
		mockResp   *osint.GetOsintResponse
		mockErr    error
		wantStatus int
	}{
		{
			name:       "OK",
			input:      `project_id=1&osint_id=1`,
			mockResp:   &osint.GetOsintResponse{},
			wantStatus: http.StatusOK,
		},
		{
			name:       "NG Invalid parameter",
			input:      `project_id=1`,
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "NG Invalid parameter",
			input:      `osint_id=1`,
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "NG Backend service error",
			input:      `project_id=1&osint_id=1`,
			wantStatus: http.StatusInternalServerError,
			mockErr:    errors.New("something wrong"),
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if c.mockResp != nil || c.mockErr != nil {
				osintMock.On("GetOsint").Return(c.mockResp, c.mockErr).Once()
			}
			// Invoke HTTP Request
			rec := httptest.NewRecorder()
			req, _ := http.NewRequest(http.MethodGet, "/api/v1/osint/get-osint/?"+c.input, nil)
			svc.getOsintHandler(rec, req)
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

func TestPutOsintHandler(t *testing.T) {
	osintMock := &mockOsintClient{}
	svc := gatewayService{
		osintClient: osintMock,
	}
	cases := []struct {
		name       string
		input      string
		mockResp   *osint.PutOsintResponse
		mockErr    error
		wantStatus int
	}{
		{
			name:       "OK",
			input:      `{"project_id":1, "osint":{"project_id":1,"resource_name":"hoge","resource_type":"fuga"}}`,
			mockResp:   &osint.PutOsintResponse{},
			wantStatus: http.StatusOK,
		},
		{
			name:       "NG Invalid parameter",
			input:      `invalid_param`,
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "NG Backend service error",
			input:      `{"project_id":1, "osint":{"project_id":1,"resource_name":"hoge","resource_type":"fuga"}}`,
			wantStatus: http.StatusInternalServerError,
			mockErr:    errors.New("something wrong"),
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if c.mockResp != nil || c.mockErr != nil {
				osintMock.On("PutOsint").Return(c.mockResp, c.mockErr).Once()
			}
			// Invoke HTTP Request
			rec := httptest.NewRecorder()
			req, _ := http.NewRequest(http.MethodPost, "/api/v1/osint/put-osint/", strings.NewReader(c.input))
			req.Header.Add("Content-Type", "application/json")
			svc.putOsintHandler(rec, req)
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

func TestDeleteOsintHandler(t *testing.T) {
	osintMock := &mockOsintClient{}
	svc := gatewayService{
		osintClient: osintMock,
	}
	cases := []struct {
		name       string
		input      string
		mockResp   *empty.Empty
		mockErr    error
		wantStatus int
	}{
		{
			name: "OK",
			input: `{"project_id":	1, "osint_id":1}`,
			mockResp:   &empty.Empty{},
			wantStatus: http.StatusOK,
		},
		{
			name:       "NG Invalid parameter",
			input:      `invalid_param`,
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "NG Backend service error",
			input: `{"project_id":	1, "osint_id":1}`,
			wantStatus: http.StatusInternalServerError,
			mockErr:    errors.New("something wrong"),
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if c.mockResp != nil || c.mockErr != nil {
				osintMock.On("DeleteOsint").Return(c.mockResp, c.mockErr).Once()
			}
			// Invoke HTTP Request
			rec := httptest.NewRecorder()
			req, _ := http.NewRequest(http.MethodPost, "/api/v1/osint/delete-osint/", strings.NewReader(c.input))
			req.Header.Add("Content-Type", "application/json")
			svc.deleteOsintHandler(rec, req)
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

func TestListOsintDataSourceHandler(t *testing.T) {
	osintMock := &mockOsintClient{}
	svc := gatewayService{
		osintClient: osintMock,
	}
	cases := []struct {
		name       string
		input      string
		mockResp   *osint.ListOsintDataSourceResponse
		mockErr    error
		wantStatus int
	}{
		{
			name:       "OK",
			input:      `project_id=1`,
			mockResp:   &osint.ListOsintDataSourceResponse{},
			wantStatus: http.StatusOK,
		},
		{
			name:       "NG Invalid parameter",
			input:      `name=123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901`,
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
				osintMock.On("ListOsintDataSource").Return(c.mockResp, c.mockErr).Once()
			}
			// Invoke HTTP Request
			rec := httptest.NewRecorder()
			req, _ := http.NewRequest(http.MethodGet, "/api/v1/osint/list-datasource/?"+c.input, nil)
			svc.listOsintDataSourceHandler(rec, req)
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

func TestGetOsintDataSourceHandler(t *testing.T) {
	osintMock := &mockOsintClient{}
	svc := gatewayService{
		osintClient: osintMock,
	}
	cases := []struct {
		name       string
		input      string
		mockResp   *osint.GetOsintDataSourceResponse
		mockErr    error
		wantStatus int
	}{
		{
			name:       "OK",
			input:      `project_id=1&osint_data_source_id=1`,
			mockResp:   &osint.GetOsintDataSourceResponse{},
			wantStatus: http.StatusOK,
		},
		{
			name:       "NG Invalid parameter",
			input:      `project_id=1`,
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "NG Invalid parameter",
			input:      `osint_data_source_id=1`,
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "NG Backend service error",
			input:      `project_id=1&osint_data_source_id=1`,
			wantStatus: http.StatusInternalServerError,
			mockErr:    errors.New("something wrong"),
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if c.mockResp != nil || c.mockErr != nil {
				osintMock.On("GetOsintDataSource").Return(c.mockResp, c.mockErr).Once()
			}
			// Invoke HTTP Request
			rec := httptest.NewRecorder()
			req, _ := http.NewRequest(http.MethodGet, "/api/v1/osint/get-datasource/?"+c.input, nil)
			svc.getOsintDataSourceHandler(rec, req)
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

func TestPutOsintDataSourceHandler(t *testing.T) {
	osintMock := &mockOsintClient{}
	svc := gatewayService{
		osintClient: osintMock,
	}
	cases := []struct {
		name       string
		input      string
		mockResp   *osint.PutOsintDataSourceResponse
		mockErr    error
		wantStatus int
	}{
		{
			name:       "OK",
			input:      `{"project_id":1, "osint_data_source":{"name":"osint_data_source-name","description":"description","max_score":10.0}}`,
			mockResp:   &osint.PutOsintDataSourceResponse{},
			wantStatus: http.StatusOK,
		},
		{
			name:       "NG Invalid parameter",
			input:      `invalid_param`,
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "NG Backend service error",
			input:      `{"project_id":1, "osint_data_source":{"name":"osint_data_source-name","description":"description","max_score":10.0}}`,
			wantStatus: http.StatusInternalServerError,
			mockErr:    errors.New("something wrong"),
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if c.mockResp != nil || c.mockErr != nil {
				osintMock.On("PutOsintDataSource").Return(c.mockResp, c.mockErr).Once()
			}
			// Invoke HTTP Request
			rec := httptest.NewRecorder()
			req, _ := http.NewRequest(http.MethodPost, "/api/v1/osint/put-datasource/", strings.NewReader(c.input))
			req.Header.Add("Content-Type", "application/json")
			svc.putOsintDataSourceHandler(rec, req)
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

func TestDeleteOsintDataSourceHandler(t *testing.T) {
	osintMock := &mockOsintClient{}
	svc := gatewayService{
		osintClient: osintMock,
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
			input:      `{"project_id":1, "osint_data_source_id":1}`,
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
			input:      `{"project_id":1, "osint_data_source_id":1}`,
			wantStatus: http.StatusInternalServerError,
			mockErr:    errors.New("something wrong"),
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if c.mockResp != nil || c.mockErr != nil {
				osintMock.On("DeleteOsintDataSource").Return(c.mockResp, c.mockErr).Once()
			}
			// Invoke HTTP Request
			rec := httptest.NewRecorder()
			req, _ := http.NewRequest(http.MethodPost, "/api/v1/osint/delete-datasource/", strings.NewReader(c.input))
			req.Header.Add("Content-Type", "application/json")
			svc.deleteOsintDataSourceHandler(rec, req)
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

func TestListRelOsintDataSourceHandler(t *testing.T) {
	osintMock := &mockOsintClient{}
	svc := gatewayService{
		osintClient: osintMock,
	}
	cases := []struct {
		name       string
		input      string
		mockResp   *osint.ListRelOsintDataSourceResponse
		mockErr    error
		wantStatus int
	}{
		{
			name:       "OK",
			input:      `project_id=1`,
			mockResp:   &osint.ListRelOsintDataSourceResponse{},
			wantStatus: http.StatusOK,
		},
		{
			name:       "NG Invalid parameter",
			input:      `name=123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901`,
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
				osintMock.On("ListRelOsintDataSource").Return(c.mockResp, c.mockErr).Once()
			}
			// Invoke HTTP Request
			rec := httptest.NewRecorder()
			req, _ := http.NewRequest(http.MethodGet, "/api/v1/osint/list-rel-datasource/?"+c.input, nil)
			svc.listRelOsintDataSourceHandler(rec, req)
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

func TestGetRelOsintDataSourceHandler(t *testing.T) {
	osintMock := &mockOsintClient{}
	svc := gatewayService{
		osintClient: osintMock,
	}
	cases := []struct {
		name       string
		input      string
		mockResp   *osint.GetRelOsintDataSourceResponse
		mockErr    error
		wantStatus int
	}{
		{
			name:       "OK",
			input:      `project_id=1&rel_osint_data_source_id=1`,
			mockResp:   &osint.GetRelOsintDataSourceResponse{},
			wantStatus: http.StatusOK,
		},
		{
			name:       "NG Invalid parameter",
			input:      `project_id=1`,
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "NG Invalid parameter",
			input:      `rel_osint_data_source_id=1`,
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "NG Backend service error",
			input:      `project_id=1&rel_osint_data_source_id=1`,
			wantStatus: http.StatusInternalServerError,
			mockErr:    errors.New("something wrong"),
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if c.mockResp != nil || c.mockErr != nil {
				osintMock.On("GetRelOsintDataSource").Return(c.mockResp, c.mockErr).Once()
			}
			// Invoke HTTP Request
			rec := httptest.NewRecorder()
			req, _ := http.NewRequest(http.MethodGet, "/api/v1/osint/get-rel-datasource/?"+c.input, nil)
			svc.getRelOsintDataSourceHandler(rec, req)
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

func TestPutRelOsintDataSourceHandler(t *testing.T) {
	osintMock := &mockOsintClient{}
	svc := gatewayService{
		osintClient: osintMock,
	}
	cases := []struct {
		name       string
		input      string
		mockResp   *osint.PutRelOsintDataSourceResponse
		mockErr    error
		wantStatus int
	}{
		{
			name:       "OK",
			input:      `{"project_id":1, "rel_osint_data_source":{"project_id":1,"osint_id":1,"osint_data_source_id":1,"status":1}}`,
			mockResp:   &osint.PutRelOsintDataSourceResponse{},
			wantStatus: http.StatusOK,
		},
		{
			name:       "NG Invalid parameter",
			input:      `invalid_param`,
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "NG Backend service error",
			input:      `{"project_id":1, "rel_osint_data_source":{"project_id":1,"osint_id":1,"osint_data_source_id":1,"status":1}}`,
			wantStatus: http.StatusInternalServerError,
			mockErr:    errors.New("something wrong"),
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if c.mockResp != nil || c.mockErr != nil {
				osintMock.On("PutRelOsintDataSource").Return(c.mockResp, c.mockErr).Once()
			}
			// Invoke HTTP Request
			rec := httptest.NewRecorder()
			req, _ := http.NewRequest(http.MethodPost, "/api/v1/osint/put-rel-datasource/", strings.NewReader(c.input))
			req.Header.Add("Content-Type", "application/json")
			svc.putRelOsintDataSourceHandler(rec, req)
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

func TestDeleteRelOsintDataSourceHandler(t *testing.T) {
	osintMock := &mockOsintClient{}
	svc := gatewayService{
		osintClient: osintMock,
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
			input:      `{"project_id": 1, "rel_osint_data_source_id":1}`,
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
			input:      `{"project_id": 1, "rel_osint_data_source_id":1}`,
			wantStatus: http.StatusInternalServerError,
			mockErr:    errors.New("something wrong"),
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if c.mockResp != nil || c.mockErr != nil {
				osintMock.On("DeleteRelOsintDataSource").Return(c.mockResp, c.mockErr).Once()
			}
			// Invoke HTTP Request
			rec := httptest.NewRecorder()
			req, _ := http.NewRequest(http.MethodPost, "/api/v1/osint/delete-rel-datasource/", strings.NewReader(c.input))
			req.Header.Add("Content-Type", "application/json")
			svc.deleteRelOsintDataSourceHandler(rec, req)
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

func TestListOsintDetectWordHandler(t *testing.T) {
	osintMock := &mockOsintClient{}
	svc := gatewayService{
		osintClient: osintMock,
	}
	cases := []struct {
		name       string
		input      string
		mockResp   *osint.ListOsintDetectWordResponse
		mockErr    error
		wantStatus int
	}{
		{
			name:       "OK",
			input:      `project_id=1`,
			mockResp:   &osint.ListOsintDetectWordResponse{},
			wantStatus: http.StatusOK,
		},
		{
			name:       "NG Invalid parameter",
			input:      `hoge=123`,
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
				osintMock.On("ListOsintDetectWord").Return(c.mockResp, c.mockErr).Once()
			}
			// Invoke HTTP Request
			rec := httptest.NewRecorder()
			req, _ := http.NewRequest(http.MethodGet, "/api/v1/osint/list-word/?"+c.input, nil)
			svc.listOsintDetectWordHandler(rec, req)
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

func TestGetOsintDetectWordHandler(t *testing.T) {
	osintMock := &mockOsintClient{}
	svc := gatewayService{
		osintClient: osintMock,
	}
	cases := []struct {
		name       string
		input      string
		mockResp   *osint.GetOsintDetectWordResponse
		mockErr    error
		wantStatus int
	}{
		{
			name:       "OK",
			input:      `project_id=1&osint_detect_word_id=1`,
			mockResp:   &osint.GetOsintDetectWordResponse{},
			wantStatus: http.StatusOK,
		},
		{
			name:       "NG Invalid parameter",
			input:      `project_id=1`,
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "NG Invalid parameter",
			input:      `osint_detect_word_id=1`,
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "NG Backend service error",
			input:      `project_id=1&osint_detect_word_id=1`,
			wantStatus: http.StatusInternalServerError,
			mockErr:    errors.New("something wrong"),
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if c.mockResp != nil || c.mockErr != nil {
				osintMock.On("GetOsintDetectWord").Return(c.mockResp, c.mockErr).Once()
			}
			// Invoke HTTP Request
			rec := httptest.NewRecorder()
			req, _ := http.NewRequest(http.MethodGet, "/api/v1/osint/get-word/?"+c.input, nil)
			svc.getOsintDetectWordHandler(rec, req)
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

func TestPutOsintDetectWordHandler(t *testing.T) {
	osintMock := &mockOsintClient{}
	svc := gatewayService{
		osintClient: osintMock,
	}
	cases := []struct {
		name       string
		input      string
		mockResp   *osint.PutOsintDetectWordResponse
		mockErr    error
		wantStatus int
	}{
		{
			name:       "OK",
			input:      `{"project_id":1, "osint_detect_word":{"project_id":1,"word":"hoge","rel_osint_data_source_id":1}}`,
			mockResp:   &osint.PutOsintDetectWordResponse{},
			wantStatus: http.StatusOK,
		},
		{
			name:       "NG Invalid parameter",
			input:      `invalid_param`,
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "NG Backend service error",
			input:      `{"project_id":1, "osint_detect_word":{"project_id":1,"word":"hoge","rel_osint_data_source_id":1}}`,
			wantStatus: http.StatusInternalServerError,
			mockErr:    errors.New("something wrong"),
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if c.mockResp != nil || c.mockErr != nil {
				osintMock.On("PutOsintDetectWord").Return(c.mockResp, c.mockErr).Once()
			}
			// Invoke HTTP Request
			rec := httptest.NewRecorder()
			req, _ := http.NewRequest(http.MethodPost, "/api/v1/osint/put-word/", strings.NewReader(c.input))
			req.Header.Add("Content-Type", "application/json")
			svc.putOsintDetectWordHandler(rec, req)
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

func TestDeleteOsintDetectWordHandler(t *testing.T) {
	osintMock := &mockOsintClient{}
	svc := gatewayService{
		osintClient: osintMock,
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
			input:      `{"project_id": 1, "osint_detect_word_id":1}`,
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
			input:      `{"project_id": 1, "osint_detect_word_id":1}`,
			wantStatus: http.StatusInternalServerError,
			mockErr:    errors.New("something wrong"),
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if c.mockResp != nil || c.mockErr != nil {
				osintMock.On("DeleteOsintDetectWord").Return(c.mockResp, c.mockErr).Once()
			}
			// Invoke HTTP Request
			rec := httptest.NewRecorder()
			req, _ := http.NewRequest(http.MethodPost, "/api/v1/osint/delete-word/", strings.NewReader(c.input))
			req.Header.Add("Content-Type", "application/json")
			svc.deleteOsintDetectWordHandler(rec, req)
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

func TestOsintInvokeScanHandler(t *testing.T) {
	osintMock := &mockOsintClient{}
	svc := gatewayService{
		osintClient: osintMock,
	}
	cases := []struct {
		name       string
		input      string
		mockResp   *osint.InvokeScanResponse
		mockErr    error
		wantStatus int
	}{
		{
			name:       "OK",
			input:      `{"project_id": 1, "rel_osint_data_source_id":1}`,
			mockResp:   &osint.InvokeScanResponse{},
			wantStatus: http.StatusOK,
		},
		{
			name:       "NG Invalid parameter",
			input:      `invalid_param`,
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "NG Backend service error",
			input:      `{"project_id": 1, "rel_osint_data_source_id":1}`,
			wantStatus: http.StatusInternalServerError,
			mockErr:    errors.New("something wrong"),
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if c.mockResp != nil || c.mockErr != nil {
				osintMock.On("InvokeScan").Return(c.mockResp, c.mockErr).Once()
			}
			// Invoke HTTP Request
			rec := httptest.NewRecorder()
			req, _ := http.NewRequest(http.MethodPost, "/api/v1/osint/start-osint/", strings.NewReader(c.input))
			req.Header.Add("Content-Type", "application/json")
			svc.invokeOsintScanHandler(rec, req)
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
type mockOsintClient struct {
	mock.Mock
}

func (m *mockOsintClient) ListOsint(context.Context, *osint.ListOsintRequest, ...grpc.CallOption) (*osint.ListOsintResponse, error) {
	args := m.Called()
	return args.Get(0).(*osint.ListOsintResponse), args.Error(1)
}
func (m *mockOsintClient) GetOsint(context.Context, *osint.GetOsintRequest, ...grpc.CallOption) (*osint.GetOsintResponse, error) {
	args := m.Called()
	return args.Get(0).(*osint.GetOsintResponse), args.Error(1)
}
func (m *mockOsintClient) PutOsint(context.Context, *osint.PutOsintRequest, ...grpc.CallOption) (*osint.PutOsintResponse, error) {
	args := m.Called()
	return args.Get(0).(*osint.PutOsintResponse), args.Error(1)
}
func (m *mockOsintClient) DeleteOsint(context.Context, *osint.DeleteOsintRequest, ...grpc.CallOption) (*empty.Empty, error) {
	args := m.Called()
	return args.Get(0).(*empty.Empty), args.Error(1)
}
func (m *mockOsintClient) ListOsintDataSource(context.Context, *osint.ListOsintDataSourceRequest, ...grpc.CallOption) (*osint.ListOsintDataSourceResponse, error) {
	args := m.Called()
	return args.Get(0).(*osint.ListOsintDataSourceResponse), args.Error(1)
}
func (m *mockOsintClient) GetOsintDataSource(context.Context, *osint.GetOsintDataSourceRequest, ...grpc.CallOption) (*osint.GetOsintDataSourceResponse, error) {
	args := m.Called()
	return args.Get(0).(*osint.GetOsintDataSourceResponse), args.Error(1)
}
func (m *mockOsintClient) PutOsintDataSource(context.Context, *osint.PutOsintDataSourceRequest, ...grpc.CallOption) (*osint.PutOsintDataSourceResponse, error) {
	args := m.Called()
	return args.Get(0).(*osint.PutOsintDataSourceResponse), args.Error(1)
}
func (m *mockOsintClient) DeleteOsintDataSource(context.Context, *osint.DeleteOsintDataSourceRequest, ...grpc.CallOption) (*empty.Empty, error) {
	args := m.Called()
	return args.Get(0).(*empty.Empty), args.Error(1)
}
func (m *mockOsintClient) ListRelOsintDataSource(context.Context, *osint.ListRelOsintDataSourceRequest, ...grpc.CallOption) (*osint.ListRelOsintDataSourceResponse, error) {
	args := m.Called()
	return args.Get(0).(*osint.ListRelOsintDataSourceResponse), args.Error(1)
}
func (m *mockOsintClient) GetRelOsintDataSource(context.Context, *osint.GetRelOsintDataSourceRequest, ...grpc.CallOption) (*osint.GetRelOsintDataSourceResponse, error) {
	args := m.Called()
	return args.Get(0).(*osint.GetRelOsintDataSourceResponse), args.Error(1)
}
func (m *mockOsintClient) PutRelOsintDataSource(context.Context, *osint.PutRelOsintDataSourceRequest, ...grpc.CallOption) (*osint.PutRelOsintDataSourceResponse, error) {
	args := m.Called()
	return args.Get(0).(*osint.PutRelOsintDataSourceResponse), args.Error(1)
}
func (m *mockOsintClient) DeleteRelOsintDataSource(context.Context, *osint.DeleteRelOsintDataSourceRequest, ...grpc.CallOption) (*empty.Empty, error) {
	args := m.Called()
	return args.Get(0).(*empty.Empty), args.Error(1)
}
func (m *mockOsintClient) InvokeScan(context.Context, *osint.InvokeScanRequest, ...grpc.CallOption) (*osint.InvokeScanResponse, error) {
	args := m.Called()
	return args.Get(0).(*osint.InvokeScanResponse), args.Error(1)
}
func (m *mockOsintClient) ListOsintDetectWord(context.Context, *osint.ListOsintDetectWordRequest, ...grpc.CallOption) (*osint.ListOsintDetectWordResponse, error) {
	args := m.Called()
	return args.Get(0).(*osint.ListOsintDetectWordResponse), args.Error(1)
}
func (m *mockOsintClient) GetOsintDetectWord(context.Context, *osint.GetOsintDetectWordRequest, ...grpc.CallOption) (*osint.GetOsintDetectWordResponse, error) {
	args := m.Called()
	return args.Get(0).(*osint.GetOsintDetectWordResponse), args.Error(1)
}
func (m *mockOsintClient) PutOsintDetectWord(context.Context, *osint.PutOsintDetectWordRequest, ...grpc.CallOption) (*osint.PutOsintDetectWordResponse, error) {
	args := m.Called()
	return args.Get(0).(*osint.PutOsintDetectWordResponse), args.Error(1)
}
func (m *mockOsintClient) DeleteOsintDetectWord(context.Context, *osint.DeleteOsintDetectWordRequest, ...grpc.CallOption) (*empty.Empty, error) {
	args := m.Called()
	return args.Get(0).(*empty.Empty), args.Error(1)
}
func (m *mockOsintClient) InvokeScanAll(context.Context, *empty.Empty, ...grpc.CallOption) (*empty.Empty, error) {
	args := m.Called()
	return args.Get(0).(*empty.Empty), args.Error(1)
}
