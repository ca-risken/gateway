package main

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/ca-risken/datasource-api/proto/aws"
	awsmocks "github.com/ca-risken/datasource-api/proto/aws/mocks"
	"github.com/stretchr/testify/mock"
)

func TestAttachDataSourceHandler(t *testing.T) {
	awsMock := awsmocks.NewAWSServiceClient(t)
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
				awsMock.On("AttachDataSource", mock.Anything, mock.Anything).Return(c.mockResp, c.mockErr).Once()
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
