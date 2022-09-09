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
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/stretchr/testify/mock"
	"google.golang.org/grpc"
)

func TestCreateProjectHandler(t *testing.T) {
	pjMock := &mockProjectClient{}
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
				pjMock.On("CreateProject").Return(c.mockResp, c.mockErr).Once()
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

/**
 * Mock Client
**/
type mockProjectClient struct {
	mock.Mock
}

func (m *mockProjectClient) ListProject(context.Context, *project.ListProjectRequest, ...grpc.CallOption) (*project.ListProjectResponse, error) {
	args := m.Called()
	return args.Get(0).(*project.ListProjectResponse), args.Error(1)
}
func (m *mockProjectClient) CreateProject(context.Context, *project.CreateProjectRequest, ...grpc.CallOption) (*project.CreateProjectResponse, error) {
	args := m.Called()
	return args.Get(0).(*project.CreateProjectResponse), args.Error(1)
}
func (m *mockProjectClient) UpdateProject(context.Context, *project.UpdateProjectRequest, ...grpc.CallOption) (*project.UpdateProjectResponse, error) {
	args := m.Called()
	return args.Get(0).(*project.UpdateProjectResponse), args.Error(1)
}
func (m *mockProjectClient) DeleteProject(context.Context, *project.DeleteProjectRequest, ...grpc.CallOption) (*empty.Empty, error) {
	args := m.Called()
	return args.Get(0).(*empty.Empty), args.Error(1)
}
func (m *mockProjectClient) TagProject(context.Context, *project.TagProjectRequest, ...grpc.CallOption) (*project.TagProjectResponse, error) {
	args := m.Called()
	return args.Get(0).(*project.TagProjectResponse), args.Error(1)
}
func (m *mockProjectClient) UntagProject(context.Context, *project.UntagProjectRequest, ...grpc.CallOption) (*empty.Empty, error) {
	args := m.Called()
	return args.Get(0).(*empty.Empty), args.Error(1)
}
func (m *mockProjectClient) IsActive(context.Context, *project.IsActiveRequest, ...grpc.CallOption) (*project.IsActiveResponse, error) {
	args := m.Called()
	return args.Get(0).(*project.IsActiveResponse), args.Error(1)
}
