package main

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/ca-risken/core/proto/iam"
	iammocks "github.com/ca-risken/core/proto/iam/mocks"
	"github.com/ca-risken/core/proto/org_iam"
	orgiammocks "github.com/ca-risken/core/proto/org_iam/mocks"
	"github.com/ca-risken/datasource-api/proto/remediation"
	remediationmocks "github.com/ca-risken/datasource-api/proto/remediation/mocks"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestGenerateRemediationProposalRoute(t *testing.T) {
	const path = "/api/v1/remediation/generate-remediation-proposal"
	cases := []struct {
		name        string
		user        *requestUser
		body        string
		checkAuth   bool
		authorized  bool
		authErr     error
		callBackend bool
		backendErr  error
		noCSRF      bool
		contentType string
		wantStatus  int
	}{
		{name: "human", user: &requestUser{userID: 1}, checkAuth: true, authorized: true, callBackend: true, wantStatus: http.StatusOK},
		{name: "project token", user: &requestUser{accessTokenID: 2, accessTokenProjectID: 1001}, checkAuth: true, authorized: true, callBackend: true, noCSRF: true, wantStatus: http.StatusOK},
		{name: "organization token", user: &requestUser{orgAccessTokenID: 3, orgAccessTokenOrgID: 4}, checkAuth: true, authorized: true, callBackend: true, noCSRF: true, wantStatus: http.StatusOK},
		{name: "unauthenticated", wantStatus: http.StatusUnauthorized},
		{name: "missing csrf", user: &requestUser{userID: 1}, noCSRF: true, wantStatus: http.StatusForbidden},
		{name: "forbidden", user: &requestUser{userID: 1}, checkAuth: true, wantStatus: http.StatusForbidden},
		{name: "authorization error", user: &requestUser{userID: 1}, checkAuth: true, authErr: errors.New("IAM unavailable"), wantStatus: http.StatusForbidden},
		{name: "different project token", user: &requestUser{accessTokenID: 2, accessTokenProjectID: 2002}, wantStatus: http.StatusForbidden},
		{name: "missing project", user: &requestUser{userID: 1}, body: `{"finding_id":20}`, wantStatus: http.StatusForbidden},
		{name: "missing finding", user: &requestUser{userID: 1}, body: `{"project_id":1001}`, checkAuth: true, authorized: true, wantStatus: http.StatusBadRequest},
		{name: "invalid finding", user: &requestUser{userID: 1}, body: `{"project_id":1001,"finding_id":"invalid"}`, checkAuth: true, authorized: true, wantStatus: http.StatusBadRequest},
		{name: "malformed json", user: &requestUser{userID: 1}, body: `{"project_id":1001,`, wantStatus: http.StatusForbidden},
		{name: "unsupported content type", user: &requestUser{userID: 1}, checkAuth: true, authorized: true, contentType: "text/plain", wantStatus: http.StatusUnsupportedMediaType},
		{name: "unsupported data source", user: &requestUser{userID: 1}, checkAuth: true, authorized: true, callBackend: true, backendErr: status.Error(codes.InvalidArgument, "unsupported data source"), wantStatus: http.StatusBadRequest},
		{name: "finding not found", user: &requestUser{userID: 1}, checkAuth: true, authorized: true, callBackend: true, backendErr: status.Error(codes.NotFound, "finding not found"), wantStatus: http.StatusNotFound},
		{name: "cooldown", user: &requestUser{userID: 1}, checkAuth: true, authorized: true, callBackend: true, backendErr: status.Error(codes.FailedPrecondition, "cooldown"), wantStatus: http.StatusPreconditionFailed},
		{name: "backend unavailable", user: &requestUser{userID: 1}, checkAuth: true, authorized: true, callBackend: true, backendErr: status.Error(codes.Unavailable, "unavailable"), wantStatus: http.StatusServiceUnavailable},
		{name: "backend error", user: &requestUser{userID: 1}, checkAuth: true, authorized: true, callBackend: true, backendErr: errors.New("backend error"), wantStatus: http.StatusInternalServerError},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			iamMock := iammocks.NewIAMServiceClient(t)
			orgIAMMock := orgiammocks.NewOrgIAMServiceClient(t)
			remediationMock := remediationmocks.NewRemediationServiceClient(t)
			svc := gatewayService{iamClient: iamMock, org_iamClient: orgIAMMock, remediationClient: remediationMock}
			if c.checkAuth {
				switch {
				case c.user.userID != 0:
					iamMock.On("IsAuthorized", mock.Anything, mock.MatchedBy(func(req *iam.IsAuthorizedRequest) bool {
						return req.UserId == 1 && req.ProjectId == 1001 && req.ActionName == "remediation/generate-remediation-proposal" && req.ResourceName == "remediation/resource_any"
					})).Return(&iam.IsAuthorizedResponse{Ok: c.authorized}, c.authErr).Once()
				case c.user.accessTokenID != 0:
					iamMock.On("IsAuthorizedToken", mock.Anything, mock.MatchedBy(func(req *iam.IsAuthorizedTokenRequest) bool {
						return req.AccessTokenId == 2 && req.ProjectId == 1001 && req.ActionName == "remediation/generate-remediation-proposal" && req.ResourceName == "remediation/resource_any"
					})).Return(&iam.IsAuthorizedTokenResponse{Ok: c.authorized}, c.authErr).Once()
				default:
					orgIAMMock.On("IsAuthorizedOrgToken", mock.Anything, mock.MatchedBy(func(req *org_iam.IsAuthorizedOrgTokenRequest) bool {
						return req.OrganizationId == 4 && req.AccessTokenId == 3 && req.ProjectId == 1001 && req.ActionName == "remediation/generate-remediation-proposal"
					})).Return(&org_iam.IsAuthorizedOrgTokenResponse{Ok: c.authorized}, c.authErr).Once()
				}
			}
			if c.callBackend {
				remediationMock.On("GenerateRemediationProposal", mock.Anything, mock.MatchedBy(func(req *remediation.GenerateRemediationProposalRequest) bool {
					return req.ProjectId == 1001 && req.FindingId == 20
				})).Return(&remediation.GenerateRemediationProposalResponse{RemediationProposalId: 10}, c.backendErr).Once()
			}
			body := c.body
			if body == "" {
				body = `{"project_id":1001,"finding_id":20}`
			}
			req := httptest.NewRequest(http.MethodPost, path, strings.NewReader(body))
			contentType := c.contentType
			if contentType == "" {
				contentType = contenTypeJSON
			}
			req.Header.Set("Content-Type", contentType)
			if !c.noCSRF {
				req.Header.Set("X-XSRF-TOKEN", "test-token")
				req.AddCookie(&http.Cookie{Name: XSRF_TOKEN, Value: "test-token"})
			}
			if c.user != nil {
				req = req.WithContext(context.WithValue(req.Context(), userKey, c.user))
			}
			rec := httptest.NewRecorder()
			newRouter(&svc).ServeHTTP(rec, req)
			require.Equal(t, c.wantStatus, rec.Code, rec.Body.String())
			if c.wantStatus == http.StatusOK {
				require.JSONEq(t, `{"data":{"remediation_proposal_id":10}}`, rec.Body.String())
			} else if c.backendErr != nil {
				require.Contains(t, rec.Body.String(), `"error"`)
			}
		})
	}
}
