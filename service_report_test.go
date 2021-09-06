package main

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ca-risken/core/proto/report"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/stretchr/testify/mock"
	"google.golang.org/grpc"
)

func TestGetReportFindingHandler(t *testing.T) {
	reportMock := &mockReportClient{}
	svc := gatewayService{
		reportClient: reportMock,
	}
	cases := []struct {
		name       string
		input      string
		mockResp   *report.GetReportFindingResponse
		mockErr    error
		wantStatus int
	}{
		{
			name:       "OK",
			input:      `project_id=1`,
			mockResp:   &report.GetReportFindingResponse{},
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
				reportMock.On("GetReportFinding").Return(c.mockResp, c.mockErr).Once()
			}
			// Invoke HTTP Request
			rec := httptest.NewRecorder()
			req, _ := http.NewRequest(http.MethodGet, "/api/v1/report/get-report/?"+c.input, nil)
			svc.getReportFindingHandler(rec, req)
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

func TestGetReportFindingAllHandler(t *testing.T) {
	reportMock := &mockReportClient{}
	svc := gatewayService{
		reportClient: reportMock,
	}
	cases := []struct {
		name       string
		input      string
		mockResp   *report.GetReportFindingAllResponse
		mockErr    error
		wantStatus int
	}{
		{
			name:       "OK",
			input:      ``,
			mockResp:   &report.GetReportFindingAllResponse{},
			wantStatus: http.StatusOK,
		},
		{
			name:       "NG Invalid parameter",
			input:      `from_date=hogehoge`,
			wantStatus: http.StatusBadRequest,
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
				reportMock.On("GetReportFindingAll").Return(c.mockResp, c.mockErr).Once()
			}
			// Invoke HTTP Request
			rec := httptest.NewRecorder()
			req, _ := http.NewRequest(http.MethodGet, "/api/v1/report/get-report-all/?"+c.input, nil)
			svc.getReportFindingAllHandler(rec, req)
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
type mockReportClient struct {
	mock.Mock
}

func (m *mockReportClient) GetReportFinding(context.Context, *report.GetReportFindingRequest, ...grpc.CallOption) (*report.GetReportFindingResponse, error) {
	args := m.Called()
	return args.Get(0).(*report.GetReportFindingResponse), args.Error(1)
}
func (m *mockReportClient) GetReportFindingAll(context.Context, *report.GetReportFindingAllRequest, ...grpc.CallOption) (*report.GetReportFindingAllResponse, error) {
	args := m.Called()
	return args.Get(0).(*report.GetReportFindingAllResponse), args.Error(1)
}
func (m *mockReportClient) CollectReportFinding(context.Context, *empty.Empty, ...grpc.CallOption) (*empty.Empty, error) {
	args := m.Called()
	return args.Get(0).(*empty.Empty), args.Error(1)
}
