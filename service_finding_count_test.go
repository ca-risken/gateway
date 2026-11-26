package main

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ca-risken/core/proto/finding"
	findingmocks "github.com/ca-risken/core/proto/finding/mocks"
	"github.com/ca-risken/core/proto/iam"
	iammocks "github.com/ca-risken/core/proto/iam/mocks"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestGetFindingCountHandler(t *testing.T) {
	cases := []struct {
		name, query string
		wantStatus  int
		dbErr       error
		call        bool
	}{
		{name: "count only", query: "project_id=1&from_score=0.8&data_source=aws:&status=1", wantStatus: 200, call: true},
		{name: "missing project", query: "from_score=0.8", wantStatus: 400},
		{name: "invalid score", query: "project_id=1&from_score=2", wantStatus: 400},
		{name: "upstream failure", query: "project_id=1&from_score=0.8&data_source=aws:&status=1", wantStatus: 500, call: true, dbErr: errors.New("db failed")},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			client := findingmocks.NewFindingServiceClient(t)
			if c.call {
				client.On("GetFindingCount", mock.Anything, mock.MatchedBy(func(r *finding.GetFindingCountRequest) bool {
					return r.ProjectId == 1 && r.FromScore == 0.8 && len(r.DataSource) == 1 && r.DataSource[0] == "aws:" && r.Status == finding.FindingStatus_FINDING_ACTIVE
				})).Return(&finding.GetFindingCountResponse{Total: 345}, c.dbErr).Once()
			}
			svc := gatewayService{findingClient: client}
			rec := httptest.NewRecorder()
			svc.getFindingCountFindingHandler(rec, httptest.NewRequest(http.MethodGet, "/api/v1/finding/get-finding-count?"+c.query, nil))
			require.Equal(t, c.wantStatus, rec.Code)
			if c.call && c.dbErr == nil {
				require.JSONEq(t, `{"data":{"total":345}}`, rec.Body.String())
			}
		})
	}
}

func TestGetFindingCountProjectAuthorization(t *testing.T) {
	cases := []struct {
		name              string
		signedIn, allowed bool
		status            int
	}{
		{name: "anonymous", status: 401},
		{name: "denied", signedIn: true, status: 403},
		{name: "authorized", signedIn: true, allowed: true, status: 200},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			iamClient := iammocks.NewIAMServiceClient(t)
			findingClient := findingmocks.NewFindingServiceClient(t)
			svc := &gatewayService{uidHeader: "X-User-ID", iamClient: iamClient, findingClient: findingClient}
			req := httptest.NewRequest(http.MethodGet, "/api/v1/finding/get-finding-count?project_id=1", nil)
			if c.signedIn {
				iamClient.On("GetUser", mock.Anything, mock.Anything).Return(&iam.GetUserResponse{User: &iam.User{UserId: 10}}, nil).Once()
				iamClient.On("IsAuthorized", mock.Anything, mock.MatchedBy(func(r *iam.IsAuthorizedRequest) bool {
					return r.UserId == 10 && r.ProjectId == 1 && r.ActionName == "finding/get-finding-count"
				})).Return(&iam.IsAuthorizedResponse{Ok: c.allowed}, nil).Once()
				req.Header.Set(svc.uidHeader, "sub")
				req.Header.Set("X-XSRF-TOKEN", "csrf")
				req.AddCookie(&http.Cookie{Name: XSRF_TOKEN, Value: "csrf"})
			}
			if c.allowed {
				findingClient.On("GetFindingCount", mock.Anything, mock.Anything).Return(&finding.GetFindingCountResponse{}, nil).Once()
			}
			rec := httptest.NewRecorder()
			newRouter(svc).ServeHTTP(rec, req)
			require.Equal(t, c.status, rec.Code)
		})
	}
}
