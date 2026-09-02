package main

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ca-risken/core/proto/iam"
	iammocks "github.com/ca-risken/core/proto/iam/mocks"
	"github.com/ca-risken/core/proto/org_iam"
	orgiammocks "github.com/ca-risken/core/proto/org_iam/mocks"
	"github.com/ca-risken/core/proto/report"
	reportmocks "github.com/ca-risken/core/proto/report/mocks"
	"github.com/stretchr/testify/mock"
)

func TestGetReportFindingForOrganizationReportHandler(t *testing.T) {
	cases := []struct {
		name       string
		query      string
		setup      func(*reportmocks.ReportServiceClient)
		wantStatus int
	}{
		{
			name:  "OK with organization projects and filters",
			query: "organization_id=1&project_id=2&project_id=3&from_date=2026-08-01&to_date=2026-08-31&score=0.8&data_source=aws,google",
			setup: func(reportMock *reportmocks.ReportServiceClient) {
				reportMock.On("GetReportFindingForOrganization", mock.Anything, mock.MatchedBy(func(req *report.GetReportFindingForOrganizationRequest) bool {
					return req.OrganizationId == 1 &&
						len(req.ProjectId) == 2 && req.ProjectId[0] == 2 && req.ProjectId[1] == 3 &&
						req.FromDate == "2026-08-01" && req.ToDate == "2026-08-31" &&
						req.Score == 0.8 && len(req.DataSource) == 2
				})).Return(&report.GetReportFindingForOrganizationResponse{}, nil).Once()
			},
			wantStatus: http.StatusOK,
		},
		{
			name:       "NG missing organization id",
			query:      "from_date=2026-08-01&to_date=2026-08-31",
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "NG invalid project id",
			query:      "organization_id=1&project_id=invalid",
			wantStatus: http.StatusBadRequest,
		},
		{
			name:  "NG backend error",
			query: "organization_id=1",
			setup: func(reportMock *reportmocks.ReportServiceClient) {
				reportMock.On("GetReportFindingForOrganization", mock.Anything, mock.Anything).Return(nil, errors.New("something wrong")).Once()
			},
			wantStatus: http.StatusInternalServerError,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			reportMock := reportmocks.NewReportServiceClient(t)
			if c.setup != nil {
				c.setup(reportMock)
			}
			svc := gatewayService{reportClient: reportMock}
			rec := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodGet, "/api/v1/report/get-report-finding-for-organization?"+c.query, nil)

			svc.getReportFindingForOrganizationReportHandler(rec, req)

			if rec.Code != c.wantStatus {
				t.Fatalf("Unexpected HTTP status code: want=%d, got=%d", c.wantStatus, rec.Code)
			}
		})
	}
}

func TestGetReportFindingForOrganizationRoute(t *testing.T) {
	cases := []struct {
		name       string
		inputUser  *requestUser
		setup      func(*orgiammocks.OrgIAMServiceClient, *reportmocks.ReportServiceClient)
		wantStatus int
	}{
		{
			name:       "NG unauthenticated",
			wantStatus: http.StatusUnauthorized,
		},
		{
			name:      "NG unauthorized organization",
			inputUser: &requestUser{sub: "sub", userID: 10},
			setup: func(orgIAMMock *orgiammocks.OrgIAMServiceClient, _ *reportmocks.ReportServiceClient) {
				orgIAMMock.On("IsAuthorizedOrg", mock.Anything, mock.MatchedBy(func(req *org_iam.IsAuthorizedOrgRequest) bool {
					return req.UserId == 10 && req.OrganizationId == 1 && req.ActionName == "report/get-report-finding-for-organization"
				})).Return(&org_iam.IsAuthorizedOrgResponse{Ok: false}, nil).Once()
			},
			wantStatus: http.StatusForbidden,
		},
		{
			name:      "OK authorized organization",
			inputUser: &requestUser{sub: "sub", userID: 10},
			setup: func(orgIAMMock *orgiammocks.OrgIAMServiceClient, reportMock *reportmocks.ReportServiceClient) {
				orgIAMMock.On("IsAuthorizedOrg", mock.Anything, mock.Anything).Return(&org_iam.IsAuthorizedOrgResponse{Ok: true}, nil).Once()
				reportMock.On("GetReportFindingForOrganization", mock.Anything, mock.Anything).Return(&report.GetReportFindingForOrganizationResponse{}, nil).Once()
			},
			wantStatus: http.StatusOK,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			iamMock := iammocks.NewIAMServiceClient(t)
			orgIAMMock := orgiammocks.NewOrgIAMServiceClient(t)
			reportMock := reportmocks.NewReportServiceClient(t)
			if c.setup != nil {
				c.setup(orgIAMMock, reportMock)
			}
			svc := &gatewayService{
				uidHeader:     "X-User-ID",
				iamClient:     iamMock,
				org_iamClient: orgIAMMock,
				reportClient:  reportMock,
			}
			rec := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodGet, "/api/v1/report/get-report-finding-for-organization?organization_id=1", nil)
			if c.inputUser != nil {
				iamMock.On("GetUser", mock.Anything, mock.MatchedBy(func(req *iam.GetUserRequest) bool {
					return req.Sub == c.inputUser.sub
				})).Return(&iam.GetUserResponse{User: &iam.User{UserId: c.inputUser.userID}}, nil).Once()
				req.Header.Set(svc.uidHeader, c.inputUser.sub)
				req.Header.Set("X-XSRF-TOKEN", "csrf-token")
				req.AddCookie(&http.Cookie{Name: XSRF_TOKEN, Value: "csrf-token"})
			}

			newRouter(svc).ServeHTTP(rec, req)

			if rec.Code != c.wantStatus {
				t.Fatalf("Unexpected HTTP status code: want=%d, got=%d", c.wantStatus, rec.Code)
			}
		})
	}
}
