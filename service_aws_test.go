package main

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/ca-risken/aws/proto/activity"
	"github.com/ca-risken/aws/proto/aws"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/stretchr/testify/mock"
	"google.golang.org/grpc"
)

func TestListAWSHandler(t *testing.T) {
	awsMock := &mockAWSClient{}
	svc := gatewayService{
		awsClient: awsMock,
	}
	cases := []struct {
		name       string
		input      string
		mockResp   *aws.ListAWSResponse
		mockErr    error
		wantStatus int
	}{
		{
			name:       "OK",
			input:      `project_id=1`,
			mockResp:   &aws.ListAWSResponse{},
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
				awsMock.On("ListAWS").Return(c.mockResp, c.mockErr).Once()
			}
			// Invoke HTTP Request
			rec := httptest.NewRecorder()
			req, _ := http.NewRequest(http.MethodGet, "/api/v1/aws/list-aws/?"+c.input, nil)
			svc.listAWSHandler(rec, req)
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

func TestPutAWSHandler(t *testing.T) {
	awsMock := &mockAWSClient{}
	svc := gatewayService{
		awsClient: awsMock,
	}
	cases := []struct {
		name       string
		input      string
		mockResp   *aws.PutAWSResponse
		mockErr    error
		wantStatus int
	}{
		{
			name:       "OK",
			input:      `{"project_id":1, "aws":{"project_id":1, "name":"aws-name", "aws_account_id":"123456789012"}}`,
			mockResp:   &aws.PutAWSResponse{},
			wantStatus: http.StatusOK,
		},
		{
			name:       "NG Invalid parameter",
			input:      `invalid_param`,
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "NG Backend service error",
			input:      `{"project_id":1, "aws":{"project_id":1, "name":"aws-name", "aws_account_id":"123456789012"}}`,
			wantStatus: http.StatusInternalServerError,
			mockErr:    errors.New("something wrong"),
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if c.mockResp != nil || c.mockErr != nil {
				awsMock.On("PutAWS").Return(c.mockResp, c.mockErr).Once()
			}
			// Invoke HTTP Request
			rec := httptest.NewRecorder()
			req, _ := http.NewRequest(http.MethodPost, "/api/v1/aws/put-aws/", strings.NewReader(c.input))
			req.Header.Add("Content-Type", "application/json")
			svc.putAWSHandler(rec, req)
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

func TestDeleteAWSHandler(t *testing.T) {
	awsMock := &mockAWSClient{}
	svc := gatewayService{
		awsClient: awsMock,
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
			input: `{"project_id":	1, "aws_id":1}`,
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
			input: `{"project_id":	1, "aws_id":1}`,
			wantStatus: http.StatusInternalServerError,
			mockErr:    errors.New("something wrong"),
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if c.mockResp != nil || c.mockErr != nil {
				awsMock.On("DeleteAWS").Return(c.mockResp, c.mockErr).Once()
			}
			// Invoke HTTP Request
			rec := httptest.NewRecorder()
			req, _ := http.NewRequest(http.MethodPost, "/api/v1/aws/delete-aws/", strings.NewReader(c.input))
			req.Header.Add("Content-Type", "application/json")
			svc.deleteAWSHandler(rec, req)
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

func TestListDataSourceHandler(t *testing.T) {
	awsMock := &mockAWSClient{}
	svc := gatewayService{
		awsClient: awsMock,
	}
	cases := []struct {
		name       string
		input      string
		mockResp   *aws.ListDataSourceResponse
		mockErr    error
		wantStatus int
	}{
		{
			name:       "OK",
			input:      `project_id=1&aws_id=1`,
			mockResp:   &aws.ListDataSourceResponse{},
			wantStatus: http.StatusOK,
		},
		{
			name:       "NG Invalid parameter",
			input:      `project_id=1`,
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "NG Backend service error",
			input:      `project_id=1&aws_id=1`,
			wantStatus: http.StatusInternalServerError,
			mockErr:    errors.New("something wrong"),
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if c.mockResp != nil || c.mockErr != nil {
				awsMock.On("ListDataSource").Return(c.mockResp, c.mockErr).Once()
			}
			// Invoke HTTP Request
			rec := httptest.NewRecorder()
			req, _ := http.NewRequest(http.MethodGet, "/api/v1/aws/list-datasource/?"+c.input, nil)
			svc.listDataSourceHandler(rec, req)
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

func TestAttachDataSourceHandler(t *testing.T) {
	awsMock := &mockAWSClient{}
	svc := gatewayService{
		awsClient: awsMock,
	}
	cases := []struct {
		name       string
		input      string
		mockResp   *aws.AttachDataSourceResponse
		mockErr    error
		wantStatus int
	}{
		{
			name:       "OK",
			input:      `{"project_id":1, "attach_data_source":{"aws_id":1, "aws_data_source_id":1, "project_id":1, "assume_role_arn":"arn", "external_id":"12345678"}}`,
			mockResp:   &aws.AttachDataSourceResponse{},
			wantStatus: http.StatusOK,
		},
		{
			name:       "NG Invalid parameter",
			input:      `invalid_param`,
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "NG Backend service error",
			input:      `{"project_id":1, "attach_data_source":{"aws_id":1, "aws_data_source_id":1, "project_id":1, "assume_role_arn":"arn", "external_id":"12345678"}}`,
			wantStatus: http.StatusInternalServerError,
			mockErr:    errors.New("something wrong"),
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if c.mockResp != nil || c.mockErr != nil {
				awsMock.On("AttachDataSource").Return(c.mockResp, c.mockErr).Once()
			}
			// Invoke HTTP Request
			rec := httptest.NewRecorder()
			req, _ := http.NewRequest(http.MethodPost, "/api/v1/aws/attach-datasource/", strings.NewReader(c.input))
			req.Header.Add("Content-Type", "application/json")
			svc.attachDataSourceHandler(rec, req)
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

func TestDetachDataSourceHandler(t *testing.T) {
	awsMock := &mockAWSClient{}
	svc := gatewayService{
		awsClient: awsMock,
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
			input:      `{"project_id":1, "aws_id":1, "aws_data_source_id":1}`,
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
			input:      `{"project_id":1, "aws_id":1, "aws_data_source_id":1}`,
			wantStatus: http.StatusInternalServerError,
			mockErr:    errors.New("something wrong"),
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if c.mockResp != nil || c.mockErr != nil {
				awsMock.On("DetachDataSource").Return(c.mockResp, c.mockErr).Once()
			}
			// Invoke HTTP Request
			rec := httptest.NewRecorder()
			req, _ := http.NewRequest(http.MethodPost, "/api/v1/aws/detach-datasource/", strings.NewReader(c.input))
			req.Header.Add("Content-Type", "application/json")
			svc.detachDataSourceHandler(rec, req)
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

func TestInvokeScanHandler(t *testing.T) {
	awsMock := &mockAWSClient{}
	svc := gatewayService{
		awsClient: awsMock,
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
			input:      `{"project_id":1, "aws_id":1, "aws_data_source_id":1}`,
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
			input:      `{"project_id":1, "aws_id":1, "aws_data_source_id":1}`,
			wantStatus: http.StatusInternalServerError,
			mockErr:    errors.New("something wrong"),
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if c.mockResp != nil || c.mockErr != nil {
				awsMock.On("InvokeScan").Return(c.mockResp, c.mockErr).Once()
			}
			// Invoke HTTP Request
			rec := httptest.NewRecorder()
			req, _ := http.NewRequest(http.MethodPost, "/api/v1/aws/invoke-scan/", strings.NewReader(c.input))
			req.Header.Add("Content-Type", "application/json")
			svc.invokeScanHandler(rec, req)
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

func TestDescribeARNHandler(t *testing.T) {
	m := &mockAWSActivityClient{}
	svc := gatewayService{
		awsActivityClient: m,
	}
	cases := []struct {
		name       string
		input      string
		mockResp   *activity.DescribeARNResponse
		mockErr    error
		wantStatus int
	}{
		{
			name:       "OK",
			input:      `arn=arn:aws:ec2:region:account-id:customer-gateway/cgw-id`,
			mockResp:   &activity.DescribeARNResponse{},
			wantStatus: http.StatusOK,
		},
		{
			name:       "NG Invalid parameter",
			input:      `arn=`,
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "NG Backend service error",
			input:      `arn=arn:aws:ec2:region:account-id:customer-gateway/cgw-id`,
			wantStatus: http.StatusInternalServerError,
			mockErr:    errors.New("something wrong"),
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if c.mockResp != nil || c.mockErr != nil {
				m.On("DescribeARN").Return(c.mockResp, c.mockErr).Once()
			}
			rec := httptest.NewRecorder()
			req, _ := http.NewRequest(http.MethodGet, "/api/v1/aws/describe-arn/?"+c.input, nil)
			svc.describeARNHandler(rec, req)
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
func TestListCloudTrailHandler(t *testing.T) {
	m := &mockAWSActivityClient{}
	svc := gatewayService{
		awsActivityClient: m,
	}
	cases := []struct {
		name       string
		input      string
		mockResp   *activity.ListCloudTrailResponse
		mockErr    error
		wantStatus int
	}{
		{
			name:       "OK",
			input:      `project_id=1&aws_id=1&region=ap-northeast-1`,
			mockResp:   &activity.ListCloudTrailResponse{},
			wantStatus: http.StatusOK,
		},
		{
			name:       "NG Invalid parameter",
			input:      `project_id=`,
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "NG Backend service error",
			input:      `project_id=1&aws_id=1&region=ap-northeast-1`,
			wantStatus: http.StatusInternalServerError,
			mockErr:    errors.New("something wrong"),
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if c.mockResp != nil || c.mockErr != nil {
				m.On("ListCloudTrail").Return(c.mockResp, c.mockErr).Once()
			}
			// Invoke HTTP Request
			rec := httptest.NewRecorder()
			req, _ := http.NewRequest(http.MethodGet, "/api/v1/aws/list-cloudtrail/?"+c.input, nil)
			svc.listCloudTrailHandler(rec, req)
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
func TestListConfigHistoryHandler(t *testing.T) {
	mock := &mockAWSActivityClient{}
	svc := gatewayService{
		awsActivityClient: mock,
	}
	cases := []struct {
		name       string
		input      string
		mockResp   *activity.ListConfigHistoryResponse
		mockErr    error
		wantStatus int
	}{
		{
			name:       "OK",
			input:      `project_id=1&aws_id=1&region=ap-northeast-1&resource_type=AWS::S3::Bucket&resource_id=bucket_name`,
			mockResp:   &activity.ListConfigHistoryResponse{},
			wantStatus: http.StatusOK,
		},
		{
			name:       "NG Invalid parameter",
			input:      `project_id=`,
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "NG Backend service error",
			input:      `project_id=1&aws_id=1&region=ap-northeast-1&resource_type=AWS::S3::Bucket&resource_id=bucket_name`,
			wantStatus: http.StatusInternalServerError,
			mockErr:    errors.New("something wrong"),
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if c.mockResp != nil || c.mockErr != nil {
				mock.On("ListConfigHistory").Return(c.mockResp, c.mockErr).Once()
			}
			// Invoke HTTP Request
			rec := httptest.NewRecorder()
			req, _ := http.NewRequest(http.MethodGet, "/api/v1/aws/list-config-history/?"+c.input, nil)
			svc.listConfigHistoryHandler(rec, req)
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
type mockAWSClient struct {
	mock.Mock
}

func (m *mockAWSClient) ListAWS(context.Context, *aws.ListAWSRequest, ...grpc.CallOption) (*aws.ListAWSResponse, error) {
	args := m.Called()
	return args.Get(0).(*aws.ListAWSResponse), args.Error(1)
}
func (m *mockAWSClient) PutAWS(context.Context, *aws.PutAWSRequest, ...grpc.CallOption) (*aws.PutAWSResponse, error) {
	args := m.Called()
	return args.Get(0).(*aws.PutAWSResponse), args.Error(1)
}
func (m *mockAWSClient) DeleteAWS(context.Context, *aws.DeleteAWSRequest, ...grpc.CallOption) (*empty.Empty, error) {
	args := m.Called()
	return args.Get(0).(*empty.Empty), args.Error(1)
}
func (m *mockAWSClient) ListDataSource(context.Context, *aws.ListDataSourceRequest, ...grpc.CallOption) (*aws.ListDataSourceResponse, error) {
	args := m.Called()
	return args.Get(0).(*aws.ListDataSourceResponse), args.Error(1)
}
func (m *mockAWSClient) AttachDataSource(context.Context, *aws.AttachDataSourceRequest, ...grpc.CallOption) (*aws.AttachDataSourceResponse, error) {
	args := m.Called()
	return args.Get(0).(*aws.AttachDataSourceResponse), args.Error(1)
}
func (m *mockAWSClient) DetachDataSource(context.Context, *aws.DetachDataSourceRequest, ...grpc.CallOption) (*empty.Empty, error) {
	args := m.Called()
	return args.Get(0).(*empty.Empty), args.Error(1)
}
func (m *mockAWSClient) InvokeScan(context.Context, *aws.InvokeScanRequest, ...grpc.CallOption) (*empty.Empty, error) {
	args := m.Called()
	return args.Get(0).(*empty.Empty), args.Error(1)
}
func (m *mockAWSClient) InvokeScanAll(context.Context, *empty.Empty, ...grpc.CallOption) (*empty.Empty, error) {
	args := m.Called()
	return args.Get(0).(*empty.Empty), args.Error(1)
}

type mockAWSActivityClient struct {
	mock.Mock
}

func (m *mockAWSActivityClient) DescribeARN(context.Context, *activity.DescribeARNRequest, ...grpc.CallOption) (*activity.DescribeARNResponse, error) {
	args := m.Called()
	return args.Get(0).(*activity.DescribeARNResponse), args.Error(1)
}
func (m *mockAWSActivityClient) ListCloudTrail(context.Context, *activity.ListCloudTrailRequest, ...grpc.CallOption) (*activity.ListCloudTrailResponse, error) {
	args := m.Called()
	return args.Get(0).(*activity.ListCloudTrailResponse), args.Error(1)
}
func (m *mockAWSActivityClient) ListConfigHistory(context.Context, *activity.ListConfigHistoryRequest, ...grpc.CallOption) (*activity.ListConfigHistoryResponse, error) {
	args := m.Called()
	return args.Get(0).(*activity.ListConfigHistoryResponse), args.Error(1)
}
