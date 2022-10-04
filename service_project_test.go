package main

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/ca-risken/core/proto/project"
	projectmocks "github.com/ca-risken/core/proto/project/mocks"
	"github.com/stretchr/testify/mock"
)

func TestCreateProjectHandler(t *testing.T) {
	pjMock := projectmocks.NewProjectServiceClient(t)
	svc := gatewayService{
		projectClient: pjMock,
	}
	cases := []struct {
		name       string
		input      string
		mockResp   *project.CreateProjectResponse
		mockErr    error
		wantStatus int
	}{
		{
			name:       "OK",
			input:      `{"user_id":1, "name":"nm"}`,
			mockResp:   &project.CreateProjectResponse{},
			wantStatus: http.StatusOK,
		},
		{
			name:       "NG Invalid parameter",
			input:      `invalid_param`,
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "NG Backend service error",
			input:      `{"user_id":1, "name":"nm"}`,
			wantStatus: http.StatusInternalServerError,
			mockErr:    errors.New("something wrong"),
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if c.mockResp != nil || c.mockErr != nil {
				pjMock.On("CreateProject", mock.Anything, mock.Anything).Return(c.mockResp, c.mockErr).Once()
			}
			// Invoke HTTP Request
			rec := httptest.NewRecorder()
			req, _ := http.NewRequest(http.MethodPost, "/api/v1/project/create-project/", strings.NewReader(c.input))
			req = req.WithContext(context.WithValue(req.Context(), userKey, &requestUser{sub: "sub", userID: 1}))
			req.Header.Add("Content-Type", "application/json")
			svc.createProjectHandler(rec, req)
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
