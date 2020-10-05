package main

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/CyberAgent/mimosa-core/proto/alert"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/stretchr/testify/mock"
	"google.golang.org/grpc"
)

func TestListAlertHandler(t *testing.T) {
	alertMock := &mockAlertClient{}
	svc := gatewayService{
		alertClient: alertMock,
	}
	cases := []struct {
		name       string
		input      string
		mockResp   *alert.ListAlertResponse
		mockErr    error
		wantStatus int
	}{
		{
			name:       "OK",
			input:      `project_id=1`,
			mockResp:   &alert.ListAlertResponse{},
			wantStatus: http.StatusOK,
		},
		{
			name:       "NG Invalid severity",
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
				alertMock.On("ListAlert").Return(c.mockResp, c.mockErr).Once()
			}
			// Invoke HTTP Request
			rec := httptest.NewRecorder()
			req, _ := http.NewRequest(http.MethodGet, "/api/v1/alert/list-alert/?"+c.input, nil)
			svc.listAlertHandler(rec, req)
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

func TestGetAlertHandler(t *testing.T) {
	alertMock := &mockAlertClient{}
	svc := gatewayService{
		alertClient: alertMock,
	}
	cases := []struct {
		name       string
		input      string
		mockResp   *alert.GetAlertResponse
		mockErr    error
		wantStatus int
	}{
		{
			name:       "OK",
			input:      `project_id=1&alert_id=1`,
			mockResp:   &alert.GetAlertResponse{},
			wantStatus: http.StatusOK,
		},
		{
			name:       "NG Invalid parameter",
			input:      `project_id=1`,
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "NG Invalid parameter",
			input:      `alert_id=1`,
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "NG Backend service error",
			input:      `project_id=1&alert_id=1`,
			wantStatus: http.StatusInternalServerError,
			mockErr:    errors.New("something wrong"),
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if c.mockResp != nil || c.mockErr != nil {
				alertMock.On("GetAlert").Return(c.mockResp, c.mockErr).Once()
			}
			// Invoke HTTP Request
			rec := httptest.NewRecorder()
			req, _ := http.NewRequest(http.MethodGet, "/api/v1/alert/get-alert/?"+c.input, nil)
			svc.getAlertHandler(rec, req)
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

func TestPutAlertHandler(t *testing.T) {
	alertMock := &mockAlertClient{}
	svc := gatewayService{
		alertClient: alertMock,
	}
	cases := []struct {
		name       string
		input      string
		mockResp   *alert.PutAlertResponse
		mockErr    error
		wantStatus int
	}{
		{
			name:       "OK",
			input:      `{"project_id":1001, "alert":{"alert_condition_id":1001,"description":"desc","severity":"high","project_id":1001,"activated":true}}`,
			mockResp:   &alert.PutAlertResponse{},
			wantStatus: http.StatusOK,
		},
		{
			name:       "NG Invalid parameter",
			input:      `invalid_param`,
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "NG Backend service error",
			input:      `{"project_id":1001, "alert":{"alert_condition_id":1001,"description":"desc","severity":"high","project_id":1001,"activated":true}}`,
			wantStatus: http.StatusInternalServerError,
			mockErr:    errors.New("something wrong"),
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if c.mockResp != nil || c.mockErr != nil {
				alertMock.On("PutAlert").Return(c.mockResp, c.mockErr).Once()
			}
			// Invoke HTTP Request
			rec := httptest.NewRecorder()
			req, _ := http.NewRequest(http.MethodPost, "/api/v1/alert/put-alert/", strings.NewReader(c.input))
			req.Header.Add("Content-Type", "application/json")
			svc.putAlertHandler(rec, req)
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

func TestDeleteAlertHandler(t *testing.T) {
	alertMock := &mockAlertClient{}
	svc := gatewayService{
		alertClient: alertMock,
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
			input: `{"project_id":	1, "alert_id":1}`,
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
			input: `{"project_id":	1, "alert_id":1}`,
			wantStatus: http.StatusInternalServerError,
			mockErr:    errors.New("something wrong"),
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if c.mockResp != nil || c.mockErr != nil {
				alertMock.On("DeleteAlert").Return(c.mockResp, c.mockErr).Once()
			}
			// Invoke HTTP Request
			rec := httptest.NewRecorder()
			req, _ := http.NewRequest(http.MethodPost, "/api/v1/alert/delete-alert/", strings.NewReader(c.input))
			req.Header.Add("Content-Type", "application/json")
			svc.deleteAlertHandler(rec, req)
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

func TestListAlertHistoryHandler(t *testing.T) {
	alertMock := &mockAlertClient{}
	svc := gatewayService{
		alertClient: alertMock,
	}
	cases := []struct {
		name       string
		input      string
		mockResp   *alert.ListAlertHistoryResponse
		mockErr    error
		wantStatus int
	}{
		{
			name:       "OK",
			input:      `project_id=1`,
			mockResp:   &alert.ListAlertHistoryResponse{},
			wantStatus: http.StatusOK,
		},
		{
			name:       "NG Invalid severity",
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
				alertMock.On("ListAlertHistory").Return(c.mockResp, c.mockErr).Once()
			}
			// Invoke HTTP Request
			rec := httptest.NewRecorder()
			req, _ := http.NewRequest(http.MethodGet, "/api/v1/alert/list-history/?"+c.input, nil)
			svc.listAlertHistoryHandler(rec, req)
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

func TestGetAlertHistoryHandler(t *testing.T) {
	alertMock := &mockAlertClient{}
	svc := gatewayService{
		alertClient: alertMock,
	}
	cases := []struct {
		name       string
		input      string
		mockResp   *alert.GetAlertHistoryResponse
		mockErr    error
		wantStatus int
	}{
		{
			name:       "OK",
			input:      `project_id=1&alert_history_id=1`,
			mockResp:   &alert.GetAlertHistoryResponse{},
			wantStatus: http.StatusOK,
		},
		{
			name:       "NG Invalid parameter",
			input:      `project_id=1`,
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "NG Invalid parameter",
			input:      `alert_history_id=1`,
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "NG Backend service error",
			input:      `project_id=1&alert_history_id=1`,
			wantStatus: http.StatusInternalServerError,
			mockErr:    errors.New("something wrong"),
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if c.mockResp != nil || c.mockErr != nil {
				alertMock.On("GetAlertHistory").Return(c.mockResp, c.mockErr).Once()
			}
			// Invoke HTTP Request
			rec := httptest.NewRecorder()
			req, _ := http.NewRequest(http.MethodGet, "/api/v1/alert/get-history/?"+c.input, nil)
			svc.getAlertHistoryHandler(rec, req)
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

func TestPutAlertHistoryHandler(t *testing.T) {
	alertMock := &mockAlertClient{}
	svc := gatewayService{
		alertClient: alertMock,
	}
	cases := []struct {
		name       string
		input      string
		mockResp   *alert.PutAlertHistoryResponse
		mockErr    error
		wantStatus int
	}{
		{
			name:       "OK",
			input:      `{"project_id":1001, "alert_history":{"alert_id":1001,"description":"desc","severity":"high","project_id":1001,"history_type":"created"}}`,
			mockResp:   &alert.PutAlertHistoryResponse{},
			wantStatus: http.StatusOK,
		},
		{
			name:       "NG Invalid parameter",
			input:      `invalid_param`,
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "NG Backend service error",
			input:      `{"project_id":1001, "alert_history":{"alert_id":1001,"description":"desc","severity":"high","project_id":1001,"history_type":"created"}}`,
			wantStatus: http.StatusInternalServerError,
			mockErr:    errors.New("something wrong"),
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if c.mockResp != nil || c.mockErr != nil {
				alertMock.On("PutAlertHistory").Return(c.mockResp, c.mockErr).Once()
			}
			// Invoke HTTP Request
			rec := httptest.NewRecorder()
			req, _ := http.NewRequest(http.MethodPost, "/api/v1/alert/put-history/", strings.NewReader(c.input))
			req.Header.Add("Content-Type", "application/json")
			svc.putAlertHistoryHandler(rec, req)
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

func TestDeleteAlertHistoryHandler(t *testing.T) {
	alertMock := &mockAlertClient{}
	svc := gatewayService{
		alertClient: alertMock,
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
			input: `{"project_id":	1, "alert_history_id":1}`,
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
			input: `{"project_id":	1, "alert_history_id":1}`,
			wantStatus: http.StatusInternalServerError,
			mockErr:    errors.New("something wrong"),
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if c.mockResp != nil || c.mockErr != nil {
				alertMock.On("DeleteAlertHistory").Return(c.mockResp, c.mockErr).Once()
			}
			// Invoke HTTP Request
			rec := httptest.NewRecorder()
			req, _ := http.NewRequest(http.MethodPost, "/api/v1/alert/delete-alert/", strings.NewReader(c.input))
			req.Header.Add("Content-Type", "application/json")
			svc.deleteAlertHistoryHandler(rec, req)
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

func TestListRelAlertFindingHandler(t *testing.T) {
	alertMock := &mockAlertClient{}
	svc := gatewayService{
		alertClient: alertMock,
	}
	cases := []struct {
		name       string
		input      string
		mockResp   *alert.ListRelAlertFindingResponse
		mockErr    error
		wantStatus int
	}{
		{
			name:       "OK",
			input:      `project_id=1`,
			mockResp:   &alert.ListRelAlertFindingResponse{},
			wantStatus: http.StatusOK,
		},
		{
			name:       "NG Invalid severity",
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
				alertMock.On("ListRelAlertFinding").Return(c.mockResp, c.mockErr).Once()
			}
			// Invoke HTTP Request
			rec := httptest.NewRecorder()
			req, _ := http.NewRequest(http.MethodGet, "/api/v1/alert/list-rel_alert_finding/?"+c.input, nil)
			svc.listRelAlertFindingHandler(rec, req)
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

func TestGetRelAlertFindingHandler(t *testing.T) {
	alertMock := &mockAlertClient{}
	svc := gatewayService{
		alertClient: alertMock,
	}
	cases := []struct {
		name       string
		input      string
		mockResp   *alert.GetRelAlertFindingResponse
		mockErr    error
		wantStatus int
	}{
		{
			name:       "OK",
			input:      `project_id=1&alert_id=1&finding_id=1`,
			mockResp:   &alert.GetRelAlertFindingResponse{},
			wantStatus: http.StatusOK,
		},
		{
			name:       "NG Invalid parameter",
			input:      `project_id=1`,
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "NG Backend service error",
			input:      `project_id=1&alert_id=1&finding_id=1`,
			wantStatus: http.StatusInternalServerError,
			mockErr:    errors.New("something wrong"),
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if c.mockResp != nil || c.mockErr != nil {
				alertMock.On("GetRelAlertFinding").Return(c.mockResp, c.mockErr).Once()
			}
			// Invoke HTTP Request
			rec := httptest.NewRecorder()
			req, _ := http.NewRequest(http.MethodGet, "/api/v1/alert/get-rel_alert_finding/?"+c.input, nil)
			svc.getRelAlertFindingHandler(rec, req)
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

func TestPutRelAlertFindingHandler(t *testing.T) {
	alertMock := &mockAlertClient{}
	svc := gatewayService{
		alertClient: alertMock,
	}
	cases := []struct {
		name       string
		input      string
		mockResp   *alert.PutRelAlertFindingResponse
		mockErr    error
		wantStatus int
	}{
		{
			name:       "OK",
			input:      `{"project_id":1001, "rel_alert_finding":{"alert_id":1001,"finding_id":1001,"project_id":1001}}`,
			mockResp:   &alert.PutRelAlertFindingResponse{},
			wantStatus: http.StatusOK,
		},
		{
			name:       "NG Invalid parameter",
			input:      `invalid_param`,
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "NG Backend service error",
			input:      `{"project_id":1001, "rel_alert_finding":{"alert_id":1001,"finding_id":1001,"project_id":1001}}`,
			wantStatus: http.StatusInternalServerError,
			mockErr:    errors.New("something wrong"),
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if c.mockResp != nil || c.mockErr != nil {
				alertMock.On("PutRelAlertFinding").Return(c.mockResp, c.mockErr).Once()
			}
			// Invoke HTTP Request
			rec := httptest.NewRecorder()
			req, _ := http.NewRequest(http.MethodPost, "/api/v1/alert/put-rel_alert_finding/", strings.NewReader(c.input))
			req.Header.Add("Content-Type", "application/json")
			svc.putRelAlertFindingHandler(rec, req)
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

func TestDeleteRelAlertFindingHandler(t *testing.T) {
	alertMock := &mockAlertClient{}
	svc := gatewayService{
		alertClient: alertMock,
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
			input: `{"project_id":	1, "alert_id":1,"finding_id":1}`,
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
			input: `{"project_id":	1, "alert_id":1,"finding_id":1}`,
			wantStatus: http.StatusInternalServerError,
			mockErr:    errors.New("something wrong"),
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if c.mockResp != nil || c.mockErr != nil {
				alertMock.On("DeleteRelAlertFinding").Return(c.mockResp, c.mockErr).Once()
			}
			// Invoke HTTP Request
			rec := httptest.NewRecorder()
			req, _ := http.NewRequest(http.MethodPost, "/api/v1/alert/delete-rel_alert_finding/", strings.NewReader(c.input))
			req.Header.Add("Content-Type", "application/json")
			svc.deleteRelAlertFindingHandler(rec, req)
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

func TestListAlertConditionHandler(t *testing.T) {
	alertMock := &mockAlertClient{}
	svc := gatewayService{
		alertClient: alertMock,
	}
	cases := []struct {
		name       string
		input      string
		mockResp   *alert.ListAlertConditionResponse
		mockErr    error
		wantStatus int
	}{
		{
			name:       "OK",
			input:      `project_id=1`,
			mockResp:   &alert.ListAlertConditionResponse{},
			wantStatus: http.StatusOK,
		},
		{
			name:       "NG Invalid severity",
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
				alertMock.On("ListAlertCondition").Return(c.mockResp, c.mockErr).Once()
			}
			// Invoke HTTP Request
			rec := httptest.NewRecorder()
			req, _ := http.NewRequest(http.MethodGet, "/api/v1/alert/list-condition/?"+c.input, nil)
			svc.listAlertConditionHandler(rec, req)
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

func TestGetAlertConditionHandler(t *testing.T) {
	alertMock := &mockAlertClient{}
	svc := gatewayService{
		alertClient: alertMock,
	}
	cases := []struct {
		name       string
		input      string
		mockResp   *alert.GetAlertConditionResponse
		mockErr    error
		wantStatus int
	}{
		{
			name:       "OK",
			input:      `project_id=1&alert_condition_id=1`,
			mockResp:   &alert.GetAlertConditionResponse{},
			wantStatus: http.StatusOK,
		},
		{
			name:       "NG Invalid parameter",
			input:      `project_id=1`,
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "NG Backend service error",
			input:      `project_id=1&alert_condition_id=1`,
			wantStatus: http.StatusInternalServerError,
			mockErr:    errors.New("something wrong"),
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if c.mockResp != nil || c.mockErr != nil {
				alertMock.On("GetAlertCondition").Return(c.mockResp, c.mockErr).Once()
			}
			// Invoke HTTP Request
			rec := httptest.NewRecorder()
			req, _ := http.NewRequest(http.MethodGet, "/api/v1/alert/get-condition/?"+c.input, nil)
			svc.getAlertConditionHandler(rec, req)
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

func TestPutAlertConditionHandler(t *testing.T) {
	alertMock := &mockAlertClient{}
	svc := gatewayService{
		alertClient: alertMock,
	}
	cases := []struct {
		name       string
		input      string
		mockResp   *alert.PutAlertConditionResponse
		mockErr    error
		wantStatus int
	}{
		{
			name:       "OK",
			input:      `{"project_id":1001, "alert_condition":{"description":"test_desc","severity":"high","and_or":"and","enabled":true,"project_id":1001}}`,
			mockResp:   &alert.PutAlertConditionResponse{},
			wantStatus: http.StatusOK,
		},
		{
			name:       "NG Invalid parameter",
			input:      `invalid_param`,
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "NG Backend service error",
			input:      `{"project_id":1001, "alert_condition":{"description":"test_desc","severity":"high","and_or":"and","enabled":true,"project_id":1001}}`,
			wantStatus: http.StatusInternalServerError,
			mockErr:    errors.New("something wrong"),
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if c.mockResp != nil || c.mockErr != nil {
				alertMock.On("PutAlertCondition").Return(c.mockResp, c.mockErr).Once()
			}
			// Invoke HTTP Request
			rec := httptest.NewRecorder()
			req, _ := http.NewRequest(http.MethodPost, "/api/v1/alert/put-condition/", strings.NewReader(c.input))
			req.Header.Add("Content-Type", "application/json")
			svc.putAlertConditionHandler(rec, req)
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

func TestDeleteAlertConditionHandler(t *testing.T) {
	alertMock := &mockAlertClient{}
	svc := gatewayService{
		alertClient: alertMock,
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
			input: `{"project_id":	1, "alert_condition_id":1}`,
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
			input: `{"project_id":	1, "alert_condition_id":1}`,
			wantStatus: http.StatusInternalServerError,
			mockErr:    errors.New("something wrong"),
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if c.mockResp != nil || c.mockErr != nil {
				alertMock.On("DeleteAlertCondition").Return(c.mockResp, c.mockErr).Once()
			}
			// Invoke HTTP Request
			rec := httptest.NewRecorder()
			req, _ := http.NewRequest(http.MethodPost, "/api/v1/alert/delete-condition/", strings.NewReader(c.input))
			req.Header.Add("Content-Type", "application/json")
			svc.deleteAlertConditionHandler(rec, req)
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

func TestListAlertRuleHandler(t *testing.T) {
	alertMock := &mockAlertClient{}
	svc := gatewayService{
		alertClient: alertMock,
	}
	cases := []struct {
		name       string
		input      string
		mockResp   *alert.ListAlertRuleResponse
		mockErr    error
		wantStatus int
	}{
		{
			name:       "OK",
			input:      `project_id=1`,
			mockResp:   &alert.ListAlertRuleResponse{},
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
				alertMock.On("ListAlertRule").Return(c.mockResp, c.mockErr).Once()
			}
			// Invoke HTTP Request
			rec := httptest.NewRecorder()
			req, _ := http.NewRequest(http.MethodGet, "/api/v1/alert/list-condition/?"+c.input, nil)
			svc.listAlertRuleHandler(rec, req)
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

func TestGetAlertRuleHandler(t *testing.T) {
	alertMock := &mockAlertClient{}
	svc := gatewayService{
		alertClient: alertMock,
	}
	cases := []struct {
		name       string
		input      string
		mockResp   *alert.GetAlertRuleResponse
		mockErr    error
		wantStatus int
	}{
		{
			name:       "OK",
			input:      `project_id=1&alert_rule_id=1`,
			mockResp:   &alert.GetAlertRuleResponse{},
			wantStatus: http.StatusOK,
		},
		{
			name:       "NG Invalid parameter",
			input:      `project_id=1`,
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "NG Backend service error",
			input:      `project_id=1&alert_rule_id=1`,
			wantStatus: http.StatusInternalServerError,
			mockErr:    errors.New("something wrong"),
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if c.mockResp != nil || c.mockErr != nil {
				alertMock.On("GetAlertRule").Return(c.mockResp, c.mockErr).Once()
			}
			// Invoke HTTP Request
			rec := httptest.NewRecorder()
			req, _ := http.NewRequest(http.MethodGet, "/api/v1/alert/get-condition/?"+c.input, nil)
			svc.getAlertRuleHandler(rec, req)
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

func TestPutAlertRuleHandler(t *testing.T) {
	alertMock := &mockAlertClient{}
	svc := gatewayService{
		alertClient: alertMock,
	}
	cases := []struct {
		name       string
		input      string
		mockResp   *alert.PutAlertRuleResponse
		mockErr    error
		wantStatus int
	}{
		{
			name:       "OK",
			input:      `{"project_id":1001, "alert_rule": {"name": "test_desc", "score":0.1, "resource_name": "test_rn", "tag": "test_tag", "finding_cnt": 1, "project_id": 1001}}`,
			mockResp:   &alert.PutAlertRuleResponse{},
			wantStatus: http.StatusOK,
		},
		{
			name:       "NG Invalid parameter",
			input:      `invalid_param`,
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "NG Backend service error",
			input:      `{"project_id":1001, "alert_rule":{"name":"test_desc","score":0.1,"resource_name":"test_rn","tag":"test_tag","finding_cnt":1,"project_id":1001}}`,
			wantStatus: http.StatusInternalServerError,
			mockErr:    errors.New("something wrong"),
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if c.mockResp != nil || c.mockErr != nil {
				alertMock.On("PutAlertRule").Return(c.mockResp, c.mockErr).Once()
			}
			// Invoke HTTP Request
			rec := httptest.NewRecorder()
			req, _ := http.NewRequest(http.MethodPost, "/api/v1/alert/put-rule/", strings.NewReader(c.input))
			req.Header.Add("Content-Type", "application/json")
			svc.putAlertRuleHandler(rec, req)
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

func TestDeleteAlertRuleHandler(t *testing.T) {
	alertMock := &mockAlertClient{}
	svc := gatewayService{
		alertClient: alertMock,
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
			input: `{"project_id":	1, "alert_rule_id":1}`,
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
			input: `{"project_id":	1, "alert_rule_id":1}`,
			wantStatus: http.StatusInternalServerError,
			mockErr:    errors.New("something wrong"),
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if c.mockResp != nil || c.mockErr != nil {
				alertMock.On("DeleteAlertRule").Return(c.mockResp, c.mockErr).Once()
			}
			// Invoke HTTP Request
			rec := httptest.NewRecorder()
			req, _ := http.NewRequest(http.MethodPost, "/api/v1/alert/delete-rule/", strings.NewReader(c.input))
			req.Header.Add("Content-Type", "application/json")
			svc.deleteAlertRuleHandler(rec, req)
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

func TestListAlertCondRuleHandler(t *testing.T) {
	alertMock := &mockAlertClient{}
	svc := gatewayService{
		alertClient: alertMock,
	}
	cases := []struct {
		name       string
		input      string
		mockResp   *alert.ListAlertCondRuleResponse
		mockErr    error
		wantStatus int
	}{
		{
			name:       "OK",
			input:      `project_id=1`,
			mockResp:   &alert.ListAlertCondRuleResponse{},
			wantStatus: http.StatusOK,
		},
		{
			name:       "NG Invalid severity",
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
				alertMock.On("ListAlertCondRule").Return(c.mockResp, c.mockErr).Once()
			}
			// Invoke HTTP Request
			rec := httptest.NewRecorder()
			req, _ := http.NewRequest(http.MethodGet, "/api/v1/alert/list-condition_rule/?"+c.input, nil)
			svc.listAlertCondRuleHandler(rec, req)
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

func TestGetAlertCondRuleHandler(t *testing.T) {
	alertMock := &mockAlertClient{}
	svc := gatewayService{
		alertClient: alertMock,
	}
	cases := []struct {
		name       string
		input      string
		mockResp   *alert.GetAlertCondRuleResponse
		mockErr    error
		wantStatus int
	}{
		{
			name:       "OK",
			input:      `project_id=1&alert_condition_id=1&alert_rule_id=1`,
			mockResp:   &alert.GetAlertCondRuleResponse{},
			wantStatus: http.StatusOK,
		},
		{
			name:       "NG Invalid parameter",
			input:      `project_id=1`,
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "NG Backend service error",
			input:      `project_id=1&alert_condition_id=1&alert_rule_id=1`,
			wantStatus: http.StatusInternalServerError,
			mockErr:    errors.New("something wrong"),
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if c.mockResp != nil || c.mockErr != nil {
				alertMock.On("GetAlertCondRule").Return(c.mockResp, c.mockErr).Once()
			}
			// Invoke HTTP Request
			rec := httptest.NewRecorder()
			req, _ := http.NewRequest(http.MethodGet, "/api/v1/alert/get-condition_rule/?"+c.input, nil)
			svc.getAlertCondRuleHandler(rec, req)
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

func TestPutAlertCondRuleHandler(t *testing.T) {
	alertMock := &mockAlertClient{}
	svc := gatewayService{
		alertClient: alertMock,
	}
	cases := []struct {
		name       string
		input      string
		mockResp   *alert.PutAlertCondRuleResponse
		mockErr    error
		wantStatus int
	}{
		{
			name:       "OK",
			input:      `{"project_id":1001, "alert_cond_rule":{"alert_condition_id":1001,"alert_rule_id":1001,"project_id":1001}}`,
			mockResp:   &alert.PutAlertCondRuleResponse{},
			wantStatus: http.StatusOK,
		},
		{
			name:       "NG Invalid parameter",
			input:      `invalid_param`,
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "NG Backend service error",
			input:      `{"project_id":1001, "alert_cond_rule":{"alert_condition_id":1001,"alert_rule_id":1001,"project_id":1001}}`,
			wantStatus: http.StatusInternalServerError,
			mockErr:    errors.New("something wrong"),
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if c.mockResp != nil || c.mockErr != nil {
				alertMock.On("PutAlertCondRule").Return(c.mockResp, c.mockErr).Once()
			}
			// Invoke HTTP Request
			rec := httptest.NewRecorder()
			req, _ := http.NewRequest(http.MethodPost, "/api/v1/alert/put-condition_rule/", strings.NewReader(c.input))
			req.Header.Add("Content-Type", "application/json")
			svc.putAlertCondRuleHandler(rec, req)
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

func TestDeleteAlertCondRuleHandler(t *testing.T) {
	alertMock := &mockAlertClient{}
	svc := gatewayService{
		alertClient: alertMock,
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
			input: `{"project_id":	1, "alert_condition_id":1,"alert_rule_id":1}`,
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
			input: `{"project_id":	1, "alert_condition_id":1,"alert_rule_id":1}`,
			wantStatus: http.StatusInternalServerError,
			mockErr:    errors.New("something wrong"),
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if c.mockResp != nil || c.mockErr != nil {
				alertMock.On("DeleteAlertCondRule").Return(c.mockResp, c.mockErr).Once()
			}
			// Invoke HTTP Request
			rec := httptest.NewRecorder()
			req, _ := http.NewRequest(http.MethodPost, "/api/v1/alert/delete-condition_rule/", strings.NewReader(c.input))
			req.Header.Add("Content-Type", "application/json")
			svc.deleteAlertCondRuleHandler(rec, req)
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

func TestListNotificationHandler(t *testing.T) {
	alertMock := &mockAlertClient{}
	svc := gatewayService{
		alertClient: alertMock,
	}
	cases := []struct {
		name       string
		input      string
		mockResp   *alert.ListNotificationResponse
		mockErr    error
		wantStatus int
	}{
		{
			name:       "OK",
			input:      `project_id=1`,
			mockResp:   &alert.ListNotificationResponse{},
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
				alertMock.On("ListNotification").Return(c.mockResp, c.mockErr).Once()
			}
			// Invoke HTTP Request
			rec := httptest.NewRecorder()
			req, _ := http.NewRequest(http.MethodGet, "/api/v1/alert/list-notification/?"+c.input, nil)
			svc.listNotificationHandler(rec, req)
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

func TestGetNotificationHandler(t *testing.T) {
	alertMock := &mockAlertClient{}
	svc := gatewayService{
		alertClient: alertMock,
	}
	cases := []struct {
		name       string
		input      string
		mockResp   *alert.GetNotificationResponse
		mockErr    error
		wantStatus int
	}{
		{
			name:       "OK",
			input:      `project_id=1&notification_id=1`,
			mockResp:   &alert.GetNotificationResponse{},
			wantStatus: http.StatusOK,
		},
		{
			name:       "NG Invalid parameter",
			input:      `project_id=1`,
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "NG Backend service error",
			input:      `project_id=1&notification_id=1`,
			wantStatus: http.StatusInternalServerError,
			mockErr:    errors.New("something wrong"),
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if c.mockResp != nil || c.mockErr != nil {
				alertMock.On("GetNotification").Return(c.mockResp, c.mockErr).Once()
			}
			// Invoke HTTP Request
			rec := httptest.NewRecorder()
			req, _ := http.NewRequest(http.MethodGet, "/api/v1/alert/get-notification/?"+c.input, nil)
			svc.getNotificationHandler(rec, req)
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

func TestPutNotificationHandler(t *testing.T) {
	alertMock := &mockAlertClient{}
	svc := gatewayService{
		alertClient: alertMock,
	}
	cases := []struct {
		name       string
		input      string
		mockResp   *alert.PutNotificationResponse
		mockErr    error
		wantStatus int
	}{
		{
			name:       "OK",
			input:      `{"project_id":1001, "notification":{"name":"test_name","type":"test_type","notify_setting":"{}","project_id":1001}}`,
			mockResp:   &alert.PutNotificationResponse{},
			wantStatus: http.StatusOK,
		},
		{
			name:       "NG Invalid parameter",
			input:      `invalid_param`,
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "NG Backend service error",
			input:      `{"project_id":1001, "notification":{"name":"test_name","type":"test_type","notify_setting":"{}","project_id":1001}}`,
			wantStatus: http.StatusInternalServerError,
			mockErr:    errors.New("something wrong"),
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if c.mockResp != nil || c.mockErr != nil {
				alertMock.On("PutNotification").Return(c.mockResp, c.mockErr).Once()
			}
			// Invoke HTTP Request
			rec := httptest.NewRecorder()
			req, _ := http.NewRequest(http.MethodPost, "/api/v1/alert/put-notification/", strings.NewReader(c.input))
			req.Header.Add("Content-Type", "application/json")
			svc.putNotificationHandler(rec, req)
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

func TestDeleteNotificationHandler(t *testing.T) {
	alertMock := &mockAlertClient{}
	svc := gatewayService{
		alertClient: alertMock,
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
			input: `{"project_id":	1, "notification_id":1}`,
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
			input: `{"project_id":	1, "notification_id":1}`,
			wantStatus: http.StatusInternalServerError,
			mockErr:    errors.New("something wrong"),
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if c.mockResp != nil || c.mockErr != nil {
				alertMock.On("DeleteNotification").Return(c.mockResp, c.mockErr).Once()
			}
			// Invoke HTTP Request
			rec := httptest.NewRecorder()
			req, _ := http.NewRequest(http.MethodPost, "/api/v1/alert/delete-notification/", strings.NewReader(c.input))
			req.Header.Add("Content-Type", "application/json")
			svc.deleteNotificationHandler(rec, req)
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

func TestListAlertCondNotificationHandler(t *testing.T) {
	alertMock := &mockAlertClient{}
	svc := gatewayService{
		alertClient: alertMock,
	}
	cases := []struct {
		name       string
		input      string
		mockResp   *alert.ListAlertCondNotificationResponse
		mockErr    error
		wantStatus int
	}{
		{
			name:       "OK",
			input:      `project_id=1`,
			mockResp:   &alert.ListAlertCondNotificationResponse{},
			wantStatus: http.StatusOK,
		},
		{
			name:       "NG Invalid severity",
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
				alertMock.On("ListAlertCondNotification").Return(c.mockResp, c.mockErr).Once()
			}
			// Invoke HTTP Request
			rec := httptest.NewRecorder()
			req, _ := http.NewRequest(http.MethodGet, "/api/v1/alert/list-condition_notification/?"+c.input, nil)
			svc.listAlertCondNotificationHandler(rec, req)
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

func TestGetAlertCondNotificationHandler(t *testing.T) {
	alertMock := &mockAlertClient{}
	svc := gatewayService{
		alertClient: alertMock,
	}
	cases := []struct {
		name       string
		input      string
		mockResp   *alert.GetAlertCondNotificationResponse
		mockErr    error
		wantStatus int
	}{
		{
			name:       "OK",
			input:      `project_id=1&alert_condition_id=1&notification_id=1`,
			mockResp:   &alert.GetAlertCondNotificationResponse{},
			wantStatus: http.StatusOK,
		},
		{
			name:       "NG Invalid parameter",
			input:      `project_id=1`,
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "NG Backend service error",
			input:      `project_id=1&alert_condition_id=1&notification_id=1`,
			wantStatus: http.StatusInternalServerError,
			mockErr:    errors.New("something wrong"),
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if c.mockResp != nil || c.mockErr != nil {
				alertMock.On("GetAlertCondNotification").Return(c.mockResp, c.mockErr).Once()
			}
			// Invoke HTTP Request
			rec := httptest.NewRecorder()
			req, _ := http.NewRequest(http.MethodGet, "/api/v1/alert/get-condition_notification/?"+c.input, nil)
			svc.getAlertCondNotificationHandler(rec, req)
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

func TestPutAlertCondNotificationHandler(t *testing.T) {
	alertMock := &mockAlertClient{}
	svc := gatewayService{
		alertClient: alertMock,
	}
	cases := []struct {
		name       string
		input      string
		mockResp   *alert.PutAlertCondNotificationResponse
		mockErr    error
		wantStatus int
	}{
		{
			name:       "OK",
			input:      `{"project_id":1001, "alert_cond_notification":{"alert_condition_id":1001,"notification_id":1001,"project_id":1001}}`,
			mockResp:   &alert.PutAlertCondNotificationResponse{},
			wantStatus: http.StatusOK,
		},
		{
			name:       "NG Invalid parameter",
			input:      `invalid_param`,
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "NG Backend service error",
			input:      `{"project_id":1001, "alert_cond_notification":{"alert_condition_id":1001,"notification_id":1001,"project_id":1001}}`,
			wantStatus: http.StatusInternalServerError,
			mockErr:    errors.New("something wrong"),
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if c.mockResp != nil || c.mockErr != nil {
				alertMock.On("PutAlertCondNotification").Return(c.mockResp, c.mockErr).Once()
			}
			// Invoke HTTP Request
			rec := httptest.NewRecorder()
			req, _ := http.NewRequest(http.MethodPost, "/api/v1/alert/put-condition_notification/", strings.NewReader(c.input))
			req.Header.Add("Content-Type", "application/json")
			svc.putAlertCondNotificationHandler(rec, req)
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

func TestDeleteAlertCondNotificationHandler(t *testing.T) {
	alertMock := &mockAlertClient{}
	svc := gatewayService{
		alertClient: alertMock,
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
			input: `{"project_id":	1, "alert_condition_id":1,"notification_id":1}`,
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
			input: `{"project_id":	1, "alert_condition_id":1,"notification_id":1}`,
			wantStatus: http.StatusInternalServerError,
			mockErr:    errors.New("something wrong"),
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if c.mockResp != nil || c.mockErr != nil {
				alertMock.On("DeleteAlertCondNotification").Return(c.mockResp, c.mockErr).Once()
			}
			// Invoke HTTP Request
			rec := httptest.NewRecorder()
			req, _ := http.NewRequest(http.MethodPost, "/api/v1/alert/delete-condition_notification/", strings.NewReader(c.input))
			req.Header.Add("Content-Type", "application/json")
			svc.deleteAlertCondNotificationHandler(rec, req)
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

func TestAnalyzeAlertHandler(t *testing.T) {
	alertMock := &mockAlertClient{}
	svc := gatewayService{
		alertClient: alertMock,
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
			input:      `{"project_id": 1}`,
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
			input:      `{"project_id": 1}`,
			wantStatus: http.StatusInternalServerError,
			mockErr:    errors.New("something wrong"),
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if c.mockResp != nil || c.mockErr != nil {
				alertMock.On("AnalyzeAlert").Return(c.mockResp, c.mockErr).Once()
			}
			// Invoke HTTP Request
			rec := httptest.NewRecorder()
			req, _ := http.NewRequest(http.MethodPost, "/api/v1/alert/analyze-alert/", strings.NewReader(c.input))
			req.Header.Add("Content-Type", "application/json")
			svc.analyzeAlertHandler(rec, req)
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
type mockAlertClient struct {
	mock.Mock
}

func (m *mockAlertClient) ListAlert(context.Context, *alert.ListAlertRequest, ...grpc.CallOption) (*alert.ListAlertResponse, error) {
	args := m.Called()
	return args.Get(0).(*alert.ListAlertResponse), args.Error(1)
}
func (m *mockAlertClient) GetAlert(context.Context, *alert.GetAlertRequest, ...grpc.CallOption) (*alert.GetAlertResponse, error) {
	args := m.Called()
	return args.Get(0).(*alert.GetAlertResponse), args.Error(1)
}
func (m *mockAlertClient) PutAlert(context.Context, *alert.PutAlertRequest, ...grpc.CallOption) (*alert.PutAlertResponse, error) {
	args := m.Called()
	return args.Get(0).(*alert.PutAlertResponse), args.Error(1)
}
func (m *mockAlertClient) DeleteAlert(context.Context, *alert.DeleteAlertRequest, ...grpc.CallOption) (*empty.Empty, error) {
	args := m.Called()
	return args.Get(0).(*empty.Empty), args.Error(1)
}
func (m *mockAlertClient) ListAlertHistory(context.Context, *alert.ListAlertHistoryRequest, ...grpc.CallOption) (*alert.ListAlertHistoryResponse, error) {
	args := m.Called()
	return args.Get(0).(*alert.ListAlertHistoryResponse), args.Error(1)
}
func (m *mockAlertClient) GetAlertHistory(context.Context, *alert.GetAlertHistoryRequest, ...grpc.CallOption) (*alert.GetAlertHistoryResponse, error) {
	args := m.Called()
	return args.Get(0).(*alert.GetAlertHistoryResponse), args.Error(1)
}
func (m *mockAlertClient) PutAlertHistory(context.Context, *alert.PutAlertHistoryRequest, ...grpc.CallOption) (*alert.PutAlertHistoryResponse, error) {
	args := m.Called()
	return args.Get(0).(*alert.PutAlertHistoryResponse), args.Error(1)
}
func (m *mockAlertClient) DeleteAlertHistory(context.Context, *alert.DeleteAlertHistoryRequest, ...grpc.CallOption) (*empty.Empty, error) {
	args := m.Called()
	return args.Get(0).(*empty.Empty), args.Error(1)
}
func (m *mockAlertClient) ListRelAlertFinding(context.Context, *alert.ListRelAlertFindingRequest, ...grpc.CallOption) (*alert.ListRelAlertFindingResponse, error) {
	args := m.Called()
	return args.Get(0).(*alert.ListRelAlertFindingResponse), args.Error(1)
}
func (m *mockAlertClient) GetRelAlertFinding(context.Context, *alert.GetRelAlertFindingRequest, ...grpc.CallOption) (*alert.GetRelAlertFindingResponse, error) {
	args := m.Called()
	return args.Get(0).(*alert.GetRelAlertFindingResponse), args.Error(1)
}
func (m *mockAlertClient) PutRelAlertFinding(context.Context, *alert.PutRelAlertFindingRequest, ...grpc.CallOption) (*alert.PutRelAlertFindingResponse, error) {
	args := m.Called()
	return args.Get(0).(*alert.PutRelAlertFindingResponse), args.Error(1)
}
func (m *mockAlertClient) DeleteRelAlertFinding(context.Context, *alert.DeleteRelAlertFindingRequest, ...grpc.CallOption) (*empty.Empty, error) {
	args := m.Called()
	return args.Get(0).(*empty.Empty), args.Error(1)
}
func (m *mockAlertClient) ListAlertCondition(context.Context, *alert.ListAlertConditionRequest, ...grpc.CallOption) (*alert.ListAlertConditionResponse, error) {
	args := m.Called()
	return args.Get(0).(*alert.ListAlertConditionResponse), args.Error(1)
}
func (m *mockAlertClient) GetAlertCondition(context.Context, *alert.GetAlertConditionRequest, ...grpc.CallOption) (*alert.GetAlertConditionResponse, error) {
	args := m.Called()
	return args.Get(0).(*alert.GetAlertConditionResponse), args.Error(1)
}
func (m *mockAlertClient) PutAlertCondition(context.Context, *alert.PutAlertConditionRequest, ...grpc.CallOption) (*alert.PutAlertConditionResponse, error) {
	args := m.Called()
	return args.Get(0).(*alert.PutAlertConditionResponse), args.Error(1)
}
func (m *mockAlertClient) DeleteAlertCondition(context.Context, *alert.DeleteAlertConditionRequest, ...grpc.CallOption) (*empty.Empty, error) {
	args := m.Called()
	return args.Get(0).(*empty.Empty), args.Error(1)
}
func (m *mockAlertClient) ListAlertRule(context.Context, *alert.ListAlertRuleRequest, ...grpc.CallOption) (*alert.ListAlertRuleResponse, error) {
	args := m.Called()
	return args.Get(0).(*alert.ListAlertRuleResponse), args.Error(1)
}
func (m *mockAlertClient) GetAlertRule(context.Context, *alert.GetAlertRuleRequest, ...grpc.CallOption) (*alert.GetAlertRuleResponse, error) {
	args := m.Called()
	return args.Get(0).(*alert.GetAlertRuleResponse), args.Error(1)
}
func (m *mockAlertClient) PutAlertRule(context.Context, *alert.PutAlertRuleRequest, ...grpc.CallOption) (*alert.PutAlertRuleResponse, error) {
	args := m.Called()
	return args.Get(0).(*alert.PutAlertRuleResponse), args.Error(1)
}
func (m *mockAlertClient) DeleteAlertRule(context.Context, *alert.DeleteAlertRuleRequest, ...grpc.CallOption) (*empty.Empty, error) {
	args := m.Called()
	return args.Get(0).(*empty.Empty), args.Error(1)
}
func (m *mockAlertClient) ListAlertCondRule(context.Context, *alert.ListAlertCondRuleRequest, ...grpc.CallOption) (*alert.ListAlertCondRuleResponse, error) {
	args := m.Called()
	return args.Get(0).(*alert.ListAlertCondRuleResponse), args.Error(1)
}
func (m *mockAlertClient) GetAlertCondRule(context.Context, *alert.GetAlertCondRuleRequest, ...grpc.CallOption) (*alert.GetAlertCondRuleResponse, error) {
	args := m.Called()
	return args.Get(0).(*alert.GetAlertCondRuleResponse), args.Error(1)
}
func (m *mockAlertClient) PutAlertCondRule(context.Context, *alert.PutAlertCondRuleRequest, ...grpc.CallOption) (*alert.PutAlertCondRuleResponse, error) {
	args := m.Called()
	return args.Get(0).(*alert.PutAlertCondRuleResponse), args.Error(1)
}
func (m *mockAlertClient) DeleteAlertCondRule(context.Context, *alert.DeleteAlertCondRuleRequest, ...grpc.CallOption) (*empty.Empty, error) {
	args := m.Called()
	return args.Get(0).(*empty.Empty), args.Error(1)
}
func (m *mockAlertClient) ListNotification(context.Context, *alert.ListNotificationRequest, ...grpc.CallOption) (*alert.ListNotificationResponse, error) {
	args := m.Called()
	return args.Get(0).(*alert.ListNotificationResponse), args.Error(1)
}
func (m *mockAlertClient) GetNotification(context.Context, *alert.GetNotificationRequest, ...grpc.CallOption) (*alert.GetNotificationResponse, error) {
	args := m.Called()
	return args.Get(0).(*alert.GetNotificationResponse), args.Error(1)
}
func (m *mockAlertClient) PutNotification(context.Context, *alert.PutNotificationRequest, ...grpc.CallOption) (*alert.PutNotificationResponse, error) {
	args := m.Called()
	return args.Get(0).(*alert.PutNotificationResponse), args.Error(1)
}
func (m *mockAlertClient) DeleteNotification(context.Context, *alert.DeleteNotificationRequest, ...grpc.CallOption) (*empty.Empty, error) {
	args := m.Called()
	return args.Get(0).(*empty.Empty), args.Error(1)
}
func (m *mockAlertClient) ListAlertCondNotification(context.Context, *alert.ListAlertCondNotificationRequest, ...grpc.CallOption) (*alert.ListAlertCondNotificationResponse, error) {
	args := m.Called()
	return args.Get(0).(*alert.ListAlertCondNotificationResponse), args.Error(1)
}
func (m *mockAlertClient) GetAlertCondNotification(context.Context, *alert.GetAlertCondNotificationRequest, ...grpc.CallOption) (*alert.GetAlertCondNotificationResponse, error) {
	args := m.Called()
	return args.Get(0).(*alert.GetAlertCondNotificationResponse), args.Error(1)
}
func (m *mockAlertClient) PutAlertCondNotification(context.Context, *alert.PutAlertCondNotificationRequest, ...grpc.CallOption) (*alert.PutAlertCondNotificationResponse, error) {
	args := m.Called()
	return args.Get(0).(*alert.PutAlertCondNotificationResponse), args.Error(1)
}
func (m *mockAlertClient) DeleteAlertCondNotification(context.Context, *alert.DeleteAlertCondNotificationRequest, ...grpc.CallOption) (*empty.Empty, error) {
	args := m.Called()
	return args.Get(0).(*empty.Empty), args.Error(1)
}
func (m *mockAlertClient) AnalyzeAlert(context.Context, *alert.AnalyzeAlertRequest, ...grpc.CallOption) (*empty.Empty, error) {
	args := m.Called()
	return args.Get(0).(*empty.Empty), args.Error(1)
}
