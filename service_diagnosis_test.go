package main

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/CyberAgent/mimosa-diagnosis/pkg/pb/diagnosis"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/stretchr/testify/mock"
	"google.golang.org/grpc"
)

func TestListDiagnosisHandler(t *testing.T) {
	diagnosisMock := &mockDiagnosisClient{}
	svc := gatewayService{
		diagnosisClient: diagnosisMock,
	}
	cases := []struct {
		name       string
		input      string
		mockResp   *diagnosis.ListDiagnosisResponse
		mockErr    error
		wantStatus int
	}{
		{
			name:       "OK",
			input:      `project_id=1`,
			mockResp:   &diagnosis.ListDiagnosisResponse{},
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
				diagnosisMock.On("ListDiagnosis").Return(c.mockResp, c.mockErr).Once()
			}
			// Invoke HTTP Request
			rec := httptest.NewRecorder()
			req, _ := http.NewRequest(http.MethodGet, "/api/v1/diagnosis/list-diagnosis/?"+c.input, nil)
			svc.listDiagnosisHandler(rec, req)
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

func TestGetDiagnosisHandler(t *testing.T) {
	diagnosisMock := &mockDiagnosisClient{}
	svc := gatewayService{
		diagnosisClient: diagnosisMock,
	}
	cases := []struct {
		name       string
		input      string
		mockResp   *diagnosis.GetDiagnosisResponse
		mockErr    error
		wantStatus int
	}{
		{
			name:       "OK",
			input:      `project_id=1&diagnosis_id=1`,
			mockResp:   &diagnosis.GetDiagnosisResponse{},
			wantStatus: http.StatusOK,
		},
		{
			name:       "NG Invalid parameter",
			input:      `project_id=1`,
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "NG Invalid parameter",
			input:      `diagnosis_id=1`,
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "NG Backend service error",
			input:      `project_id=1&diagnosis_id=1`,
			wantStatus: http.StatusInternalServerError,
			mockErr:    errors.New("something wrong"),
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if c.mockResp != nil || c.mockErr != nil {
				diagnosisMock.On("GetDiagnosis").Return(c.mockResp, c.mockErr).Once()
			}
			// Invoke HTTP Request
			rec := httptest.NewRecorder()
			req, _ := http.NewRequest(http.MethodGet, "/api/v1/diagnosis/get-diagnosis/?"+c.input, nil)
			svc.getDiagnosisHandler(rec, req)
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

func TestPutDiagnosisHandler(t *testing.T) {
	diagnosisMock := &mockDiagnosisClient{}
	svc := gatewayService{
		diagnosisClient: diagnosisMock,
	}
	cases := []struct {
		name       string
		input      string
		mockResp   *diagnosis.PutDiagnosisResponse
		mockErr    error
		wantStatus int
	}{
		{
			name:       "OK",
			input:      `{"project_id":1, "diagnosis":{"name":"diagnosis-name"}}`,
			mockResp:   &diagnosis.PutDiagnosisResponse{},
			wantStatus: http.StatusOK,
		},
		{
			name:       "NG Invalid parameter",
			input:      `invalid_param`,
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "NG Backend service error",
			input:      `{"project_id":1, "diagnosis":{"name":"diagnosis-name"}}`,
			wantStatus: http.StatusInternalServerError,
			mockErr:    errors.New("something wrong"),
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if c.mockResp != nil || c.mockErr != nil {
				diagnosisMock.On("PutDiagnosis").Return(c.mockResp, c.mockErr).Once()
			}
			// Invoke HTTP Request
			rec := httptest.NewRecorder()
			req, _ := http.NewRequest(http.MethodPost, "/api/v1/diagnosis/put-diagnosis/", strings.NewReader(c.input))
			req.Header.Add("Content-Type", "application/json")
			svc.putDiagnosisHandler(rec, req)
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

func TestDeleteDiagnosisHandler(t *testing.T) {
	diagnosisMock := &mockDiagnosisClient{}
	svc := gatewayService{
		diagnosisClient: diagnosisMock,
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
			input: `{"project_id":	1, "diagnosis_id":1}`,
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
			input: `{"project_id":	1, "diagnosis_id":1}`,
			wantStatus: http.StatusInternalServerError,
			mockErr:    errors.New("something wrong"),
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if c.mockResp != nil || c.mockErr != nil {
				diagnosisMock.On("DeleteDiagnosis").Return(c.mockResp, c.mockErr).Once()
			}
			// Invoke HTTP Request
			rec := httptest.NewRecorder()
			req, _ := http.NewRequest(http.MethodPost, "/api/v1/diagnosis/delete-diagnosis/", strings.NewReader(c.input))
			req.Header.Add("Content-Type", "application/json")
			svc.deleteDiagnosisHandler(rec, req)
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

func TestListDiagnosisDataSourceHandler(t *testing.T) {
	diagnosisMock := &mockDiagnosisClient{}
	svc := gatewayService{
		diagnosisClient: diagnosisMock,
	}
	cases := []struct {
		name       string
		input      string
		mockResp   *diagnosis.ListDiagnosisDataSourceResponse
		mockErr    error
		wantStatus int
	}{
		{
			name:       "OK",
			input:      `project_id=1`,
			mockResp:   &diagnosis.ListDiagnosisDataSourceResponse{},
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
				diagnosisMock.On("ListDiagnosisDataSource").Return(c.mockResp, c.mockErr).Once()
			}
			// Invoke HTTP Request
			rec := httptest.NewRecorder()
			req, _ := http.NewRequest(http.MethodGet, "/api/v1/diagnosis/list-datasource/?"+c.input, nil)
			svc.listDiagnosisDataSourceHandler(rec, req)
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

func TestGetDiagnosisDataSourceHandler(t *testing.T) {
	diagnosisMock := &mockDiagnosisClient{}
	svc := gatewayService{
		diagnosisClient: diagnosisMock,
	}
	cases := []struct {
		name       string
		input      string
		mockResp   *diagnosis.GetDiagnosisDataSourceResponse
		mockErr    error
		wantStatus int
	}{
		{
			name:       "OK",
			input:      `project_id=1&diagnosis_data_source_id=1`,
			mockResp:   &diagnosis.GetDiagnosisDataSourceResponse{},
			wantStatus: http.StatusOK,
		},
		{
			name:       "NG Invalid parameter",
			input:      `project_id=1`,
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "NG Invalid parameter",
			input:      `diagnosis_data_source_id=1`,
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "NG Backend service error",
			input:      `project_id=1&diagnosis_data_source_id=1`,
			wantStatus: http.StatusInternalServerError,
			mockErr:    errors.New("something wrong"),
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if c.mockResp != nil || c.mockErr != nil {
				diagnosisMock.On("GetDiagnosisDataSource").Return(c.mockResp, c.mockErr).Once()
			}
			// Invoke HTTP Request
			rec := httptest.NewRecorder()
			req, _ := http.NewRequest(http.MethodGet, "/api/v1/diagnosis/get-datasource/?"+c.input, nil)
			svc.getDiagnosisDataSourceHandler(rec, req)
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

func TestPutDiagnosisDataSourceHandler(t *testing.T) {
	diagnosisMock := &mockDiagnosisClient{}
	svc := gatewayService{
		diagnosisClient: diagnosisMock,
	}
	cases := []struct {
		name       string
		input      string
		mockResp   *diagnosis.PutDiagnosisDataSourceResponse
		mockErr    error
		wantStatus int
	}{
		{
			name:       "OK",
			input:      `{"project_id":1, "diagnosis_data_source":{"name":"diagnosis_data_source-name","description":"description","max_score":10.0}}`,
			mockResp:   &diagnosis.PutDiagnosisDataSourceResponse{},
			wantStatus: http.StatusOK,
		},
		{
			name:       "NG Invalid parameter",
			input:      `invalid_param`,
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "NG Backend service error",
			input:      `{"project_id":1, "diagnosis_data_source":{"name":"diagnosis_data_source-name","description":"description","max_score":10.0}}`,
			wantStatus: http.StatusInternalServerError,
			mockErr:    errors.New("something wrong"),
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if c.mockResp != nil || c.mockErr != nil {
				diagnosisMock.On("PutDiagnosisDataSource").Return(c.mockResp, c.mockErr).Once()
			}
			// Invoke HTTP Request
			rec := httptest.NewRecorder()
			req, _ := http.NewRequest(http.MethodPost, "/api/v1/diagnosis/put-datasource/", strings.NewReader(c.input))
			req.Header.Add("Content-Type", "application/json")
			svc.putDiagnosisDataSourceHandler(rec, req)
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

func TestDeleteDiagnosisDataSourceHandler(t *testing.T) {
	diagnosisMock := &mockDiagnosisClient{}
	svc := gatewayService{
		diagnosisClient: diagnosisMock,
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
			input:      `{"project_id":1, "diagnosis_data_source_id":1}`,
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
			input:      `{"project_id":1, "diagnosis_data_source_id":1}`,
			wantStatus: http.StatusInternalServerError,
			mockErr:    errors.New("something wrong"),
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if c.mockResp != nil || c.mockErr != nil {
				diagnosisMock.On("DeleteDiagnosisDataSource").Return(c.mockResp, c.mockErr).Once()
			}
			// Invoke HTTP Request
			rec := httptest.NewRecorder()
			req, _ := http.NewRequest(http.MethodPost, "/api/v1/diagnosis/delete-datasource/", strings.NewReader(c.input))
			req.Header.Add("Content-Type", "application/json")
			svc.deleteDiagnosisDataSourceHandler(rec, req)
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

func TestListRelDiagnosisDataSourceHandler(t *testing.T) {
	diagnosisMock := &mockDiagnosisClient{}
	svc := gatewayService{
		diagnosisClient: diagnosisMock,
	}
	cases := []struct {
		name       string
		input      string
		mockResp   *diagnosis.ListRelDiagnosisDataSourceResponse
		mockErr    error
		wantStatus int
	}{
		{
			name:       "OK",
			input:      `project_id=1`,
			mockResp:   &diagnosis.ListRelDiagnosisDataSourceResponse{},
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
				diagnosisMock.On("ListRelDiagnosisDataSource").Return(c.mockResp, c.mockErr).Once()
			}
			// Invoke HTTP Request
			rec := httptest.NewRecorder()
			req, _ := http.NewRequest(http.MethodGet, "/api/v1/diagnosis/list-rel-datasource/?"+c.input, nil)
			svc.listRelDiagnosisDataSourceHandler(rec, req)
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

func TestGetRelDiagnosisDataSourceHandler(t *testing.T) {
	diagnosisMock := &mockDiagnosisClient{}
	svc := gatewayService{
		diagnosisClient: diagnosisMock,
	}
	cases := []struct {
		name       string
		input      string
		mockResp   *diagnosis.GetRelDiagnosisDataSourceResponse
		mockErr    error
		wantStatus int
	}{
		{
			name:       "OK",
			input:      `project_id=1&rel_diagnosis_data_source_id=1`,
			mockResp:   &diagnosis.GetRelDiagnosisDataSourceResponse{},
			wantStatus: http.StatusOK,
		},
		{
			name:       "NG Invalid parameter",
			input:      `project_id=1`,
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "NG Invalid parameter",
			input:      `rel_diagnosis_data_source_id=1`,
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "NG Backend service error",
			input:      `project_id=1&rel_diagnosis_data_source_id=1`,
			wantStatus: http.StatusInternalServerError,
			mockErr:    errors.New("something wrong"),
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if c.mockResp != nil || c.mockErr != nil {
				diagnosisMock.On("GetRelDiagnosisDataSource").Return(c.mockResp, c.mockErr).Once()
			}
			// Invoke HTTP Request
			rec := httptest.NewRecorder()
			req, _ := http.NewRequest(http.MethodGet, "/api/v1/diagnosis/get-rel-datasource/?"+c.input, nil)
			svc.getRelDiagnosisDataSourceHandler(rec, req)
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

func TestPutRelDiagnosisDataSourceHandler(t *testing.T) {
	diagnosisMock := &mockDiagnosisClient{}
	svc := gatewayService{
		diagnosisClient: diagnosisMock,
	}
	cases := []struct {
		name       string
		input      string
		mockResp   *diagnosis.PutRelDiagnosisDataSourceResponse
		mockErr    error
		wantStatus int
	}{
		{
			name:       "OK",
			input:      `{"project_id":1, "rel_diagnosis_data_source":{"name":"rel_diagnosis_data_source-name","diagnosis_id":1,"diagnosis_data_source_id":1,"record_id":"test_record","jira_id":"test_jira_id","jira_key":"test_jira_key"}}`,
			mockResp:   &diagnosis.PutRelDiagnosisDataSourceResponse{},
			wantStatus: http.StatusOK,
		},
		{
			name:       "NG Invalid parameter",
			input:      `invalid_param`,
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "NG Backend service error",
			input:      `{"project_id":1, "rel_diagnosis_data_source":{"name":"rel_diagnosis_data_source-name","diagnosis_id":1,"diagnosis_data_source_id":1,"record_id":"test_record","jira_id":"test_jira_id","jira_key":"test_jira_key"}}`,
			wantStatus: http.StatusInternalServerError,
			mockErr:    errors.New("something wrong"),
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if c.mockResp != nil || c.mockErr != nil {
				diagnosisMock.On("PutRelDiagnosisDataSource").Return(c.mockResp, c.mockErr).Once()
			}
			// Invoke HTTP Request
			rec := httptest.NewRecorder()
			req, _ := http.NewRequest(http.MethodPost, "/api/v1/diagnosis/put-rel-datasource/", strings.NewReader(c.input))
			req.Header.Add("Content-Type", "application/json")
			svc.putRelDiagnosisDataSourceHandler(rec, req)
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

func TestDeleteRelDiagnosisDataSourceHandler(t *testing.T) {
	diagnosisMock := &mockDiagnosisClient{}
	svc := gatewayService{
		diagnosisClient: diagnosisMock,
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
			input:      `{"project_id": 1, "rel_diagnosis_data_source_id":1}`,
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
			input:      `{"project_id": 1, "rel_diagnosis_data_source_id":1}`,
			wantStatus: http.StatusInternalServerError,
			mockErr:    errors.New("something wrong"),
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if c.mockResp != nil || c.mockErr != nil {
				diagnosisMock.On("DeleteRelDiagnosisDataSource").Return(c.mockResp, c.mockErr).Once()
			}
			// Invoke HTTP Request
			rec := httptest.NewRecorder()
			req, _ := http.NewRequest(http.MethodPost, "/api/v1/diagnosis/delete-rel-datasource/", strings.NewReader(c.input))
			req.Header.Add("Content-Type", "application/json")
			svc.deleteRelDiagnosisDataSourceHandler(rec, req)
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

func TestStartDiagnosisHandler(t *testing.T) {
	diagnosisMock := &mockDiagnosisClient{}
	svc := gatewayService{
		diagnosisClient: diagnosisMock,
	}
	cases := []struct {
		name       string
		input      string
		mockResp   *diagnosis.StartDiagnosisResponse
		mockErr    error
		wantStatus int
	}{
		{
			name:       "OK",
			input:      `{"project_id": 1, "rel_diagnosis_data_source_id":1}`,
			mockResp:   &diagnosis.StartDiagnosisResponse{},
			wantStatus: http.StatusOK,
		},
		{
			name:       "NG Invalid parameter",
			input:      `invalid_param`,
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "NG Backend service error",
			input:      `{"project_id": 1, "rel_diagnosis_data_source_id":1}`,
			wantStatus: http.StatusInternalServerError,
			mockErr:    errors.New("something wrong"),
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if c.mockResp != nil || c.mockErr != nil {
				diagnosisMock.On("StartDiagnosis").Return(c.mockResp, c.mockErr).Once()
			}
			// Invoke HTTP Request
			rec := httptest.NewRecorder()
			req, _ := http.NewRequest(http.MethodPost, "/api/v1/diagnosis/start-diagnosis/", strings.NewReader(c.input))
			req.Header.Add("Content-Type", "application/json")
			svc.startDiagnosisHandler(rec, req)
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
type mockDiagnosisClient struct {
	mock.Mock
}

func (m *mockDiagnosisClient) ListDiagnosis(context.Context, *diagnosis.ListDiagnosisRequest, ...grpc.CallOption) (*diagnosis.ListDiagnosisResponse, error) {
	args := m.Called()
	return args.Get(0).(*diagnosis.ListDiagnosisResponse), args.Error(1)
}
func (m *mockDiagnosisClient) GetDiagnosis(context.Context, *diagnosis.GetDiagnosisRequest, ...grpc.CallOption) (*diagnosis.GetDiagnosisResponse, error) {
	args := m.Called()
	return args.Get(0).(*diagnosis.GetDiagnosisResponse), args.Error(1)
}
func (m *mockDiagnosisClient) PutDiagnosis(context.Context, *diagnosis.PutDiagnosisRequest, ...grpc.CallOption) (*diagnosis.PutDiagnosisResponse, error) {
	args := m.Called()
	return args.Get(0).(*diagnosis.PutDiagnosisResponse), args.Error(1)
}
func (m *mockDiagnosisClient) DeleteDiagnosis(context.Context, *diagnosis.DeleteDiagnosisRequest, ...grpc.CallOption) (*empty.Empty, error) {
	args := m.Called()
	return args.Get(0).(*empty.Empty), args.Error(1)
}
func (m *mockDiagnosisClient) ListDiagnosisDataSource(context.Context, *diagnosis.ListDiagnosisDataSourceRequest, ...grpc.CallOption) (*diagnosis.ListDiagnosisDataSourceResponse, error) {
	args := m.Called()
	return args.Get(0).(*diagnosis.ListDiagnosisDataSourceResponse), args.Error(1)
}
func (m *mockDiagnosisClient) GetDiagnosisDataSource(context.Context, *diagnosis.GetDiagnosisDataSourceRequest, ...grpc.CallOption) (*diagnosis.GetDiagnosisDataSourceResponse, error) {
	args := m.Called()
	return args.Get(0).(*diagnosis.GetDiagnosisDataSourceResponse), args.Error(1)
}
func (m *mockDiagnosisClient) PutDiagnosisDataSource(context.Context, *diagnosis.PutDiagnosisDataSourceRequest, ...grpc.CallOption) (*diagnosis.PutDiagnosisDataSourceResponse, error) {
	args := m.Called()
	return args.Get(0).(*diagnosis.PutDiagnosisDataSourceResponse), args.Error(1)
}
func (m *mockDiagnosisClient) DeleteDiagnosisDataSource(context.Context, *diagnosis.DeleteDiagnosisDataSourceRequest, ...grpc.CallOption) (*empty.Empty, error) {
	args := m.Called()
	return args.Get(0).(*empty.Empty), args.Error(1)
}
func (m *mockDiagnosisClient) ListRelDiagnosisDataSource(context.Context, *diagnosis.ListRelDiagnosisDataSourceRequest, ...grpc.CallOption) (*diagnosis.ListRelDiagnosisDataSourceResponse, error) {
	args := m.Called()
	return args.Get(0).(*diagnosis.ListRelDiagnosisDataSourceResponse), args.Error(1)
}
func (m *mockDiagnosisClient) GetRelDiagnosisDataSource(context.Context, *diagnosis.GetRelDiagnosisDataSourceRequest, ...grpc.CallOption) (*diagnosis.GetRelDiagnosisDataSourceResponse, error) {
	args := m.Called()
	return args.Get(0).(*diagnosis.GetRelDiagnosisDataSourceResponse), args.Error(1)
}
func (m *mockDiagnosisClient) PutRelDiagnosisDataSource(context.Context, *diagnosis.PutRelDiagnosisDataSourceRequest, ...grpc.CallOption) (*diagnosis.PutRelDiagnosisDataSourceResponse, error) {
	args := m.Called()
	return args.Get(0).(*diagnosis.PutRelDiagnosisDataSourceResponse), args.Error(1)
}
func (m *mockDiagnosisClient) DeleteRelDiagnosisDataSource(context.Context, *diagnosis.DeleteRelDiagnosisDataSourceRequest, ...grpc.CallOption) (*empty.Empty, error) {
	args := m.Called()
	return args.Get(0).(*empty.Empty), args.Error(1)
}
func (m *mockDiagnosisClient) StartDiagnosis(context.Context, *diagnosis.StartDiagnosisRequest, ...grpc.CallOption) (*diagnosis.StartDiagnosisResponse, error) {
	args := m.Called()
	return args.Get(0).(*diagnosis.StartDiagnosisResponse), args.Error(1)
}
