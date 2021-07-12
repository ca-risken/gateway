package main

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/CyberAgent/mimosa-diagnosis/proto/diagnosis"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/stretchr/testify/mock"
	"google.golang.org/grpc"
)

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

func TestListJiraSettingHandler(t *testing.T) {
	diagnosisMock := &mockDiagnosisClient{}
	svc := gatewayService{
		diagnosisClient: diagnosisMock,
	}
	cases := []struct {
		name       string
		input      string
		mockResp   *diagnosis.ListJiraSettingResponse
		mockErr    error
		wantStatus int
	}{
		{
			name:       "OK",
			input:      `project_id=1`,
			mockResp:   &diagnosis.ListJiraSettingResponse{},
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
				diagnosisMock.On("ListJiraSetting").Return(c.mockResp, c.mockErr).Once()
			}
			// Invoke HTTP Request
			rec := httptest.NewRecorder()
			req, _ := http.NewRequest(http.MethodGet, "/api/v1/diagnosis/list-jira-setting/?"+c.input, nil)
			svc.listJiraSettingHandler(rec, req)
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

func TestGetJiraSettingHandler(t *testing.T) {
	diagnosisMock := &mockDiagnosisClient{}
	svc := gatewayService{
		diagnosisClient: diagnosisMock,
	}
	cases := []struct {
		name       string
		input      string
		mockResp   *diagnosis.GetJiraSettingResponse
		mockErr    error
		wantStatus int
	}{
		{
			name:       "OK",
			input:      `project_id=1&jira_setting_id=1`,
			mockResp:   &diagnosis.GetJiraSettingResponse{},
			wantStatus: http.StatusOK,
		},
		{
			name:       "NG Invalid parameter",
			input:      `project_id=1`,
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "NG Invalid parameter",
			input:      `jira_setting_id=1`,
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "NG Backend service error",
			input:      `project_id=1&jira_setting_id=1`,
			wantStatus: http.StatusInternalServerError,
			mockErr:    errors.New("something wrong"),
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if c.mockResp != nil || c.mockErr != nil {
				diagnosisMock.On("GetJiraSetting").Return(c.mockResp, c.mockErr).Once()
			}
			// Invoke HTTP Request
			rec := httptest.NewRecorder()
			req, _ := http.NewRequest(http.MethodGet, "/api/v1/diagnosis/get-jira-setting/?"+c.input, nil)
			svc.getJiraSettingHandler(rec, req)
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

func TestPutJiraSettingHandler(t *testing.T) {
	diagnosisMock := &mockDiagnosisClient{}
	svc := gatewayService{
		diagnosisClient: diagnosisMock,
	}
	cases := []struct {
		name       string
		input      string
		mockResp   *diagnosis.PutJiraSettingResponse
		mockErr    error
		wantStatus int
	}{
		{
			name:       "OK",
			input:      `{"project_id":1, "jira_setting":{"project_id":1,"name":"jira_setting-name","diagnosis_data_source_id":1,"identity_field":"test_field","identity_value":"test_value","jira_id":"test_jira_id","jira_key":"test_jira_key"}}`,
			mockResp:   &diagnosis.PutJiraSettingResponse{},
			wantStatus: http.StatusOK,
		},
		{
			name:       "NG Invalid parameter",
			input:      `invalid_param`,
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "NG Backend service error",
			input:      `{"project_id":1, "jira_setting":{"project_id":1,"name":"jira_setting-name","diagnosis_data_source_id":1,"identity_field":"test_field","identity_value":"test_value","jira_id":"test_jira_id","jira_key":"test_jira_key"}}`,
			wantStatus: http.StatusInternalServerError,
			mockErr:    errors.New("something wrong"),
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if c.mockResp != nil || c.mockErr != nil {
				diagnosisMock.On("PutJiraSetting").Return(c.mockResp, c.mockErr).Once()
			}
			// Invoke HTTP Request
			rec := httptest.NewRecorder()
			req, _ := http.NewRequest(http.MethodPost, "/api/v1/diagnosis/put-jira-setting/", strings.NewReader(c.input))
			req.Header.Add("Content-Type", "application/json")
			svc.putJiraSettingHandler(rec, req)
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

func TestDeleteJiraSettingHandler(t *testing.T) {
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
			input:      `{"project_id": 1, "jira_setting_id":1}`,
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
			input:      `{"project_id": 1, "jira_setting_id":1}`,
			wantStatus: http.StatusInternalServerError,
			mockErr:    errors.New("something wrong"),
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if c.mockResp != nil || c.mockErr != nil {
				diagnosisMock.On("DeleteJiraSetting").Return(c.mockResp, c.mockErr).Once()
			}
			// Invoke HTTP Request
			rec := httptest.NewRecorder()
			req, _ := http.NewRequest(http.MethodPost, "/api/v1/diagnosis/delete-jira-setting/", strings.NewReader(c.input))
			req.Header.Add("Content-Type", "application/json")
			svc.deleteJiraSettingHandler(rec, req)
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

func TestListWpscanSettingHandler(t *testing.T) {
	diagnosisMock := &mockDiagnosisClient{}
	svc := gatewayService{
		diagnosisClient: diagnosisMock,
	}
	cases := []struct {
		name       string
		input      string
		mockResp   *diagnosis.ListWpscanSettingResponse
		mockErr    error
		wantStatus int
	}{
		{
			name:       "OK",
			input:      `project_id=1&diagnosis_data_source_id=1`,
			mockResp:   &diagnosis.ListWpscanSettingResponse{},
			wantStatus: http.StatusOK,
		},
		{
			name:       "NG Invalid parameter",
			input:      ``,
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
				diagnosisMock.On("ListWpscanSetting").Return(c.mockResp, c.mockErr).Once()
			}
			// Invoke HTTP Request
			rec := httptest.NewRecorder()
			req, _ := http.NewRequest(http.MethodGet, "/api/v1/diagnosis/list-wpscan-setting/?"+c.input, nil)
			svc.listWpscanSettingHandler(rec, req)
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

func TestGetWpscanSettingHandler(t *testing.T) {
	diagnosisMock := &mockDiagnosisClient{}
	svc := gatewayService{
		diagnosisClient: diagnosisMock,
	}
	cases := []struct {
		name       string
		input      string
		mockResp   *diagnosis.GetWpscanSettingResponse
		mockErr    error
		wantStatus int
	}{
		{
			name:       "OK",
			input:      `project_id=1&wpscan_setting_id=1`,
			mockResp:   &diagnosis.GetWpscanSettingResponse{},
			wantStatus: http.StatusOK,
		},
		{
			name:       "NG Invalid parameter",
			input:      `project_id=1`,
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "NG Invalid parameter",
			input:      `wpscan_setting_id=1`,
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "NG Backend service error",
			input:      `project_id=1&wpscan_setting_id=1`,
			wantStatus: http.StatusInternalServerError,
			mockErr:    errors.New("something wrong"),
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if c.mockResp != nil || c.mockErr != nil {
				diagnosisMock.On("GetWpscanSetting").Return(c.mockResp, c.mockErr).Once()
			}
			// Invoke HTTP Request
			rec := httptest.NewRecorder()
			req, _ := http.NewRequest(http.MethodGet, "/api/v1/diagnosis/get-wpscan-setting/?"+c.input, nil)
			svc.getWpscanSettingHandler(rec, req)
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

func TestPutWpscanSettingHandler(t *testing.T) {
	diagnosisMock := &mockDiagnosisClient{}
	svc := gatewayService{
		diagnosisClient: diagnosisMock,
	}
	cases := []struct {
		name       string
		input      string
		mockResp   *diagnosis.PutWpscanSettingResponse
		mockErr    error
		wantStatus int
	}{
		{
			name:       "OK",
			input:      `{"project_id":1, "wpscan_setting":{"project_id":1,"diagnosis_data_source_id":1,"target_url":"http://example.com"}}`,
			mockResp:   &diagnosis.PutWpscanSettingResponse{},
			wantStatus: http.StatusOK,
		},
		{
			name:       "NG Invalid parameter",
			input:      `invalid_param`,
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "NG Backend service error",
			input:      `{"project_id":1, "wpscan_setting":{"project_id":1,"diagnosis_data_source_id":1,"target_url":"http://example.com"}}`,
			wantStatus: http.StatusInternalServerError,
			mockErr:    errors.New("something wrong"),
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if c.mockResp != nil || c.mockErr != nil {
				diagnosisMock.On("PutWpscanSetting").Return(c.mockResp, c.mockErr).Once()
			}
			// Invoke HTTP Request
			rec := httptest.NewRecorder()
			req, _ := http.NewRequest(http.MethodPost, "/api/v1/diagnosis/put-wpscan-setting/", strings.NewReader(c.input))
			req.Header.Add("Content-Type", "application/json")
			svc.putWpscanSettingHandler(rec, req)
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

func TestDeleteWpscanSettingHandler(t *testing.T) {
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
			input:      `{"project_id": 1, "wpscan_setting_id":1}`,
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
			input:      `{"project_id": 1, "wpscan_setting_id":1}`,
			wantStatus: http.StatusInternalServerError,
			mockErr:    errors.New("something wrong"),
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if c.mockResp != nil || c.mockErr != nil {
				diagnosisMock.On("DeleteWpscanSetting").Return(c.mockResp, c.mockErr).Once()
			}
			// Invoke HTTP Request
			rec := httptest.NewRecorder()
			req, _ := http.NewRequest(http.MethodPost, "/api/v1/diagnosis/delete-wpscan-setting/", strings.NewReader(c.input))
			req.Header.Add("Content-Type", "application/json")
			svc.deleteWpscanSettingHandler(rec, req)
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

func TestListPortscanSettingHandler(t *testing.T) {
	diagnosisMock := &mockDiagnosisClient{}
	svc := gatewayService{
		diagnosisClient: diagnosisMock,
	}
	cases := []struct {
		name       string
		input      string
		mockResp   *diagnosis.ListPortscanSettingResponse
		mockErr    error
		wantStatus int
	}{
		{
			name:       "OK",
			input:      `project_id=1&diagnosis_data_source_id=1`,
			mockResp:   &diagnosis.ListPortscanSettingResponse{},
			wantStatus: http.StatusOK,
		},
		{
			name:       "NG Invalid parameter",
			input:      ``,
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
				diagnosisMock.On("ListPortscanSetting").Return(c.mockResp, c.mockErr).Once()
			}
			// Invoke HTTP Request
			rec := httptest.NewRecorder()
			req, _ := http.NewRequest(http.MethodGet, "/api/v1/diagnosis/list-portscan-setting/?"+c.input, nil)
			svc.listPortscanSettingHandler(rec, req)
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

func TestGetPortscanSettingHandler(t *testing.T) {
	diagnosisMock := &mockDiagnosisClient{}
	svc := gatewayService{
		diagnosisClient: diagnosisMock,
	}
	cases := []struct {
		name       string
		input      string
		mockResp   *diagnosis.GetPortscanSettingResponse
		mockErr    error
		wantStatus int
	}{
		{
			name:       "OK",
			input:      `project_id=1&portscan_setting_id=1`,
			mockResp:   &diagnosis.GetPortscanSettingResponse{},
			wantStatus: http.StatusOK,
		},
		{
			name:       "NG Invalid parameter",
			input:      `project_id=1`,
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "NG Invalid parameter",
			input:      `portscan_setting_id=1`,
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "NG Backend service error",
			input:      `project_id=1&portscan_setting_id=1`,
			wantStatus: http.StatusInternalServerError,
			mockErr:    errors.New("something wrong"),
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if c.mockResp != nil || c.mockErr != nil {
				diagnosisMock.On("GetPortscanSetting").Return(c.mockResp, c.mockErr).Once()
			}
			// Invoke HTTP Request
			rec := httptest.NewRecorder()
			req, _ := http.NewRequest(http.MethodGet, "/api/v1/diagnosis/get-portscan-setting/?"+c.input, nil)
			svc.getPortscanSettingHandler(rec, req)
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

func TestPutPortscanSettingHandler(t *testing.T) {
	diagnosisMock := &mockDiagnosisClient{}
	svc := gatewayService{
		diagnosisClient: diagnosisMock,
	}
	cases := []struct {
		name       string
		input      string
		mockResp   *diagnosis.PutPortscanSettingResponse
		mockErr    error
		wantStatus int
	}{
		{
			name:       "OK",
			input:      `{"project_id":1, "portscan_setting":{"project_id":1,"diagnosis_data_source_id":1,"name":"test_portscan"}}`,
			mockResp:   &diagnosis.PutPortscanSettingResponse{},
			wantStatus: http.StatusOK,
		},
		{
			name:       "NG Invalid parameter",
			input:      `invalid_param`,
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "NG Backend service error",
			input:      `{"project_id":1, "portscan_setting":{"project_id":1,"diagnosis_data_source_id":1,"name":"test_portscan"}}`,
			wantStatus: http.StatusInternalServerError,
			mockErr:    errors.New("something wrong"),
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if c.mockResp != nil || c.mockErr != nil {
				diagnosisMock.On("PutPortscanSetting").Return(c.mockResp, c.mockErr).Once()
			}
			// Invoke HTTP Request
			rec := httptest.NewRecorder()
			req, _ := http.NewRequest(http.MethodPost, "/api/v1/diagnosis/put-portscan-setting/", strings.NewReader(c.input))
			req.Header.Add("Content-Type", "application/json")
			svc.putPortscanSettingHandler(rec, req)
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

func TestDeletePortscanSettingHandler(t *testing.T) {
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
			input:      `{"project_id": 1, "portscan_setting_id":1}`,
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
			input:      `{"project_id": 1, "portscan_setting_id":1}`,
			wantStatus: http.StatusInternalServerError,
			mockErr:    errors.New("something wrong"),
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if c.mockResp != nil || c.mockErr != nil {
				diagnosisMock.On("DeletePortscanSetting").Return(c.mockResp, c.mockErr).Once()
			}
			// Invoke HTTP Request
			rec := httptest.NewRecorder()
			req, _ := http.NewRequest(http.MethodPost, "/api/v1/diagnosis/delete-portscan-setting/", strings.NewReader(c.input))
			req.Header.Add("Content-Type", "application/json")
			svc.deletePortscanSettingHandler(rec, req)
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

func TestListPortscanTargetHandler(t *testing.T) {
	diagnosisMock := &mockDiagnosisClient{}
	svc := gatewayService{
		diagnosisClient: diagnosisMock,
	}
	cases := []struct {
		name       string
		input      string
		mockResp   *diagnosis.ListPortscanTargetResponse
		mockErr    error
		wantStatus int
	}{
		{
			name:       "OK",
			input:      `project_id=1&portscan_setting_id=1`,
			mockResp:   &diagnosis.ListPortscanTargetResponse{},
			wantStatus: http.StatusOK,
		},
		{
			name:       "NG Invalid parameter",
			input:      ``,
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
				diagnosisMock.On("ListPortscanTarget").Return(c.mockResp, c.mockErr).Once()
			}
			// Invoke HTTP Request
			rec := httptest.NewRecorder()
			req, _ := http.NewRequest(http.MethodGet, "/api/v1/diagnosis/list-portscan-target/?"+c.input, nil)
			svc.listPortscanTargetHandler(rec, req)
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

func TestGetPortscanTargetHandler(t *testing.T) {
	diagnosisMock := &mockDiagnosisClient{}
	svc := gatewayService{
		diagnosisClient: diagnosisMock,
	}
	cases := []struct {
		name       string
		input      string
		mockResp   *diagnosis.GetPortscanTargetResponse
		mockErr    error
		wantStatus int
	}{
		{
			name:       "OK",
			input:      `project_id=1&portscan_target_id=1`,
			mockResp:   &diagnosis.GetPortscanTargetResponse{},
			wantStatus: http.StatusOK,
		},
		{
			name:       "NG Invalid parameter",
			input:      `project_id=1`,
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "NG Invalid parameter",
			input:      `portscan_target_id=1`,
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "NG Backend service error",
			input:      `project_id=1&portscan_target_id=1`,
			wantStatus: http.StatusInternalServerError,
			mockErr:    errors.New("something wrong"),
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if c.mockResp != nil || c.mockErr != nil {
				diagnosisMock.On("GetPortscanTarget").Return(c.mockResp, c.mockErr).Once()
			}
			// Invoke HTTP Request
			rec := httptest.NewRecorder()
			req, _ := http.NewRequest(http.MethodGet, "/api/v1/diagnosis/get-portscan-target/?"+c.input, nil)
			svc.getPortscanTargetHandler(rec, req)
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

func TestPutPortscanTargetHandler(t *testing.T) {
	diagnosisMock := &mockDiagnosisClient{}
	svc := gatewayService{
		diagnosisClient: diagnosisMock,
	}
	cases := []struct {
		name       string
		input      string
		mockResp   *diagnosis.PutPortscanTargetResponse
		mockErr    error
		wantStatus int
	}{
		{
			name:       "OK",
			input:      `{"project_id":1, "portscan_target":{"project_id":1,"portscan_setting_id":1,"target":"test_portscan"}}`,
			mockResp:   &diagnosis.PutPortscanTargetResponse{},
			wantStatus: http.StatusOK,
		},
		{
			name:       "NG Invalid parameter",
			input:      `invalid_param`,
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "NG Backend service error",
			input:      `{"project_id":1, "portscan_target":{"project_id":1,"portscan_setting_id":1,"target":"test_portscan"}}`,
			wantStatus: http.StatusInternalServerError,
			mockErr:    errors.New("something wrong"),
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if c.mockResp != nil || c.mockErr != nil {
				diagnosisMock.On("PutPortscanTarget").Return(c.mockResp, c.mockErr).Once()
			}
			// Invoke HTTP Request
			rec := httptest.NewRecorder()
			req, _ := http.NewRequest(http.MethodPost, "/api/v1/diagnosis/put-portscan-target/", strings.NewReader(c.input))
			req.Header.Add("Content-Type", "application/json")
			svc.putPortscanTargetHandler(rec, req)
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

func TestDeletePortscanTargetHandler(t *testing.T) {
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
			input:      `{"project_id": 1, "portscan_target_id":1}`,
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
			input:      `{"project_id": 1, "portscan_target_id":1}`,
			wantStatus: http.StatusInternalServerError,
			mockErr:    errors.New("something wrong"),
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if c.mockResp != nil || c.mockErr != nil {
				diagnosisMock.On("DeletePortscanTarget").Return(c.mockResp, c.mockErr).Once()
			}
			// Invoke HTTP Request
			rec := httptest.NewRecorder()
			req, _ := http.NewRequest(http.MethodPost, "/api/v1/diagnosis/delete-portscan-target/", strings.NewReader(c.input))
			req.Header.Add("Content-Type", "application/json")
			svc.deletePortscanTargetHandler(rec, req)
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

func TestInvokeDiagnosisScanHandler(t *testing.T) {
	diagnosisMock := &mockDiagnosisClient{}
	svc := gatewayService{
		diagnosisClient: diagnosisMock,
	}
	cases := []struct {
		name       string
		input      string
		mockResp   *diagnosis.InvokeScanResponse
		mockErr    error
		wantStatus int
	}{
		{
			name:       "OK",
			input:      `{"project_id": 1, "setting_id":1,"diagnosis_data_source_id":1}`,
			mockResp:   &diagnosis.InvokeScanResponse{},
			wantStatus: http.StatusOK,
		},
		{
			name:       "NG Invalid parameter",
			input:      `invalid_param`,
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "NG Backend service error",
			input:      `{"project_id": 1, "setting_id":1, "diagnosis_data_source_id":1}`,
			wantStatus: http.StatusInternalServerError,
			mockErr:    errors.New("something wrong"),
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if c.mockResp != nil || c.mockErr != nil {
				diagnosisMock.On("InvokeScan").Return(c.mockResp, c.mockErr).Once()
			}
			// Invoke HTTP Request
			rec := httptest.NewRecorder()
			req, _ := http.NewRequest(http.MethodPost, "/api/v1/diagnosis/invoke-scan/", strings.NewReader(c.input))
			req.Header.Add("Content-Type", "application/json")
			svc.invokeDiagnosisScanHandler(rec, req)
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
func (m *mockDiagnosisClient) ListJiraSetting(context.Context, *diagnosis.ListJiraSettingRequest, ...grpc.CallOption) (*diagnosis.ListJiraSettingResponse, error) {
	args := m.Called()
	return args.Get(0).(*diagnosis.ListJiraSettingResponse), args.Error(1)
}
func (m *mockDiagnosisClient) GetJiraSetting(context.Context, *diagnosis.GetJiraSettingRequest, ...grpc.CallOption) (*diagnosis.GetJiraSettingResponse, error) {
	args := m.Called()
	return args.Get(0).(*diagnosis.GetJiraSettingResponse), args.Error(1)
}
func (m *mockDiagnosisClient) PutJiraSetting(context.Context, *diagnosis.PutJiraSettingRequest, ...grpc.CallOption) (*diagnosis.PutJiraSettingResponse, error) {
	args := m.Called()
	return args.Get(0).(*diagnosis.PutJiraSettingResponse), args.Error(1)
}
func (m *mockDiagnosisClient) DeleteJiraSetting(context.Context, *diagnosis.DeleteJiraSettingRequest, ...grpc.CallOption) (*empty.Empty, error) {
	args := m.Called()
	return args.Get(0).(*empty.Empty), args.Error(1)
}
func (m *mockDiagnosisClient) ListWpscanSetting(context.Context, *diagnosis.ListWpscanSettingRequest, ...grpc.CallOption) (*diagnosis.ListWpscanSettingResponse, error) {
	args := m.Called()
	return args.Get(0).(*diagnosis.ListWpscanSettingResponse), args.Error(1)
}
func (m *mockDiagnosisClient) GetWpscanSetting(context.Context, *diagnosis.GetWpscanSettingRequest, ...grpc.CallOption) (*diagnosis.GetWpscanSettingResponse, error) {
	args := m.Called()
	return args.Get(0).(*diagnosis.GetWpscanSettingResponse), args.Error(1)
}
func (m *mockDiagnosisClient) PutWpscanSetting(context.Context, *diagnosis.PutWpscanSettingRequest, ...grpc.CallOption) (*diagnosis.PutWpscanSettingResponse, error) {
	args := m.Called()
	return args.Get(0).(*diagnosis.PutWpscanSettingResponse), args.Error(1)
}
func (m *mockDiagnosisClient) DeleteWpscanSetting(context.Context, *diagnosis.DeleteWpscanSettingRequest, ...grpc.CallOption) (*empty.Empty, error) {
	args := m.Called()
	return args.Get(0).(*empty.Empty), args.Error(1)
}

func (m *mockDiagnosisClient) InvokeScan(context.Context, *diagnosis.InvokeScanRequest, ...grpc.CallOption) (*diagnosis.InvokeScanResponse, error) {
	args := m.Called()
	return args.Get(0).(*diagnosis.InvokeScanResponse), args.Error(1)
}
func (m *mockDiagnosisClient) InvokeScanAll(context.Context, *empty.Empty, ...grpc.CallOption) (*empty.Empty, error) {
	args := m.Called()
	return args.Get(0).(*empty.Empty), args.Error(1)
}
func (m *mockDiagnosisClient) ListPortscanSetting(context.Context, *diagnosis.ListPortscanSettingRequest, ...grpc.CallOption) (*diagnosis.ListPortscanSettingResponse, error) {
	args := m.Called()
	return args.Get(0).(*diagnosis.ListPortscanSettingResponse), args.Error(1)
}
func (m *mockDiagnosisClient) GetPortscanSetting(context.Context, *diagnosis.GetPortscanSettingRequest, ...grpc.CallOption) (*diagnosis.GetPortscanSettingResponse, error) {
	args := m.Called()
	return args.Get(0).(*diagnosis.GetPortscanSettingResponse), args.Error(1)
}
func (m *mockDiagnosisClient) PutPortscanSetting(context.Context, *diagnosis.PutPortscanSettingRequest, ...grpc.CallOption) (*diagnosis.PutPortscanSettingResponse, error) {
	args := m.Called()
	return args.Get(0).(*diagnosis.PutPortscanSettingResponse), args.Error(1)
}
func (m *mockDiagnosisClient) DeletePortscanSetting(context.Context, *diagnosis.DeletePortscanSettingRequest, ...grpc.CallOption) (*empty.Empty, error) {
	args := m.Called()
	return args.Get(0).(*empty.Empty), args.Error(1)
}
func (m *mockDiagnosisClient) ListPortscanTarget(context.Context, *diagnosis.ListPortscanTargetRequest, ...grpc.CallOption) (*diagnosis.ListPortscanTargetResponse, error) {
	args := m.Called()
	return args.Get(0).(*diagnosis.ListPortscanTargetResponse), args.Error(1)
}
func (m *mockDiagnosisClient) GetPortscanTarget(context.Context, *diagnosis.GetPortscanTargetRequest, ...grpc.CallOption) (*diagnosis.GetPortscanTargetResponse, error) {
	args := m.Called()
	return args.Get(0).(*diagnosis.GetPortscanTargetResponse), args.Error(1)
}
func (m *mockDiagnosisClient) PutPortscanTarget(context.Context, *diagnosis.PutPortscanTargetRequest, ...grpc.CallOption) (*diagnosis.PutPortscanTargetResponse, error) {
	args := m.Called()
	return args.Get(0).(*diagnosis.PutPortscanTargetResponse), args.Error(1)
}
func (m *mockDiagnosisClient) DeletePortscanTarget(context.Context, *diagnosis.DeletePortscanTargetRequest, ...grpc.CallOption) (*empty.Empty, error) {
	args := m.Called()
	return args.Get(0).(*empty.Empty), args.Error(1)
}
