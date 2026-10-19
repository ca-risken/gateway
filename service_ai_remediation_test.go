package main

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ca-risken/core/proto/ai"
	aimocks "github.com/ca-risken/core/proto/ai/mocks"
	"github.com/ca-risken/core/proto/iam"
	iammocks "github.com/ca-risken/core/proto/iam/mocks"
	"github.com/ca-risken/core/proto/org_iam"
	orgiammocks "github.com/ca-risken/core/proto/org_iam/mocks"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestRemediationProposalRoutes(t *testing.T) {
	proposal := &ai.RemediationProposal{
		RemediationProposalId: 10,
		FindingId:             20,
		ProjectId:             1001,
		Status:                "SUCCEEDED",
		StatusDetail:          "completed",
		RemediationPlan:       `{"summary":"repair"}`,
		GeneratedAt:           100,
		CreatedAt:             90,
		UpdatedAt:             100,
	}
	endpoints := []struct {
		path     string
		method   string
		query    string
		response interface{}
		empty    interface{}
		wantJSON string
	}{
		{
			path: "/api/v1/ai/get-remediation-proposal", method: "GetRemediationProposal",
			query:    "remediation_proposal_id=10",
			response: &ai.GetRemediationProposalResponse{RemediationProposal: proposal},
			empty:    &ai.GetRemediationProposalResponse{},
			wantJSON: `{"data":{"remediation_proposal":{"remediation_proposal_id":10,"finding_id":20,"project_id":1001,"status":"SUCCEEDED","status_detail":"completed","remediation_plan":"{\"summary\":\"repair\"}","generated_at":100,"created_at":90,"updated_at":100}}}`,
		},
		{
			path: "/api/v1/ai/list-remediation-proposal", method: "ListRemediationProposal",
			query:    "finding_id=20",
			response: &ai.ListRemediationProposalResponse{RemediationProposals: []*ai.RemediationProposal{proposal}},
			empty:    &ai.ListRemediationProposalResponse{},
			wantJSON: `{"data":{"remediation_proposals":[{"remediation_proposal_id":10,"finding_id":20,"project_id":1001,"status":"SUCCEEDED","status_detail":"completed","remediation_plan":"{\"summary\":\"repair\"}","generated_at":100,"created_at":90,"updated_at":100}]}}`,
		},
	}
	for _, endpoint := range endpoints {
		t.Run(endpoint.method, func(t *testing.T) {
			type testCase struct {
				name       string
				user       *requestUser
				project    string
				query      string
				checkAuth  bool
				authorized bool
				authErr    error
				callCore   bool
				empty      bool
				coreErr    error
				wantStatus int
			}
			cases := []testCase{
				{name: "human", user: &requestUser{userID: 1}, checkAuth: true, authorized: true, callCore: true, wantStatus: http.StatusOK},
				{name: "project token", user: &requestUser{accessTokenID: 2, accessTokenProjectID: 1001}, checkAuth: true, authorized: true, callCore: true, wantStatus: http.StatusOK},
				{name: "organization token", user: &requestUser{orgAccessTokenID: 3, orgAccessTokenOrgID: 4}, checkAuth: true, authorized: true, callCore: true, wantStatus: http.StatusOK},
				{name: "unauthenticated", wantStatus: http.StatusUnauthorized},
				{name: "forbidden", user: &requestUser{userID: 1}, checkAuth: true, wantStatus: http.StatusForbidden},
				{name: "authorization error", user: &requestUser{userID: 1}, checkAuth: true, authErr: errors.New("IAM unavailable"), wantStatus: http.StatusForbidden},
				{name: "different project token", user: &requestUser{accessTokenID: 2, accessTokenProjectID: 2002}, wantStatus: http.StatusForbidden},
				{name: "missing project", user: &requestUser{userID: 1}, project: "missing", wantStatus: http.StatusForbidden},
				{name: "missing resource id", user: &requestUser{userID: 1}, query: "unused=1", checkAuth: true, authorized: true, wantStatus: http.StatusBadRequest},
				{name: "malformed resource id", user: &requestUser{userID: 1}, query: endpoint.query + "x", checkAuth: true, authorized: true, wantStatus: http.StatusBadRequest},
				{name: "empty result", user: &requestUser{userID: 1}, checkAuth: true, authorized: true, callCore: true, empty: true, wantStatus: http.StatusOK},
				{name: "database error", user: &requestUser{userID: 1}, checkAuth: true, authorized: true, callCore: true, coreErr: errors.New("database unavailable"), wantStatus: http.StatusInternalServerError},
				{name: "grpc unavailable", user: &requestUser{userID: 1}, checkAuth: true, authorized: true, callCore: true, coreErr: status.Error(codes.Unavailable, "unavailable"), wantStatus: http.StatusServiceUnavailable},
			}
			if endpoint.method == "ListRemediationProposal" {
				for _, value := range []string{"PENDING", "SUCCEEDED", "FAILED"} {
					cases = append(cases, testCase{name: value, user: &requestUser{userID: 1}, query: endpoint.query + "&status=" + value, checkAuth: true, authorized: true, callCore: true, wantStatus: http.StatusOK})
				}
				cases = append(cases, testCase{name: "invalid status", user: &requestUser{userID: 1}, query: endpoint.query + "&status=SUCCEEDED,PENDING", checkAuth: true, authorized: true, wantStatus: http.StatusBadRequest})
			}
			for _, c := range cases {
				t.Run(c.name, func(t *testing.T) {
					aiMock := aimocks.NewAIServiceClient(t)
					iamMock := iammocks.NewIAMServiceClient(t)
					orgIAMMock := orgiammocks.NewOrgIAMServiceClient(t)
					svc := gatewayService{aiClient: aiMock, iamClient: iamMock, org_iamClient: orgIAMMock}
					if c.checkAuth {
						switch {
						case c.user.userID != 0:
							iamMock.On("IsAuthorized", mock.Anything, mock.MatchedBy(func(req *iam.IsAuthorizedRequest) bool {
								return req.UserId == 1 && req.ProjectId == 1001 && req.ActionName == getActionNameFromURI(endpoint.path) && req.ResourceName == "ai/resource_any"
							})).Return(&iam.IsAuthorizedResponse{Ok: c.authorized}, c.authErr).Once()
						case c.user.accessTokenID != 0:
							iamMock.On("IsAuthorizedToken", mock.Anything, mock.MatchedBy(func(req *iam.IsAuthorizedTokenRequest) bool {
								return req.AccessTokenId == 2 && req.ProjectId == 1001 && req.ActionName == getActionNameFromURI(endpoint.path) && req.ResourceName == "ai/resource_any"
							})).Return(&iam.IsAuthorizedTokenResponse{Ok: c.authorized}, c.authErr).Once()
						default:
							orgIAMMock.On("IsAuthorizedOrgToken", mock.Anything, mock.MatchedBy(func(req *org_iam.IsAuthorizedOrgTokenRequest) bool {
								return req.OrganizationId == 4 && req.AccessTokenId == 3 && req.ProjectId == 1001 && req.ActionName == getActionNameFromURI(endpoint.path)
							})).Return(&org_iam.IsAuthorizedOrgTokenResponse{Ok: c.authorized}, c.authErr).Once()
						}
					}
					query := c.query
					if query == "" {
						query = endpoint.query
					}
					req := httptest.NewRequest(http.MethodGet, endpoint.path+"?"+query, nil)
					if c.project != "missing" {
						req.URL.RawQuery += "&project_id=1001"
					}
					if c.user != nil {
						req = req.WithContext(context.WithValue(req.Context(), userKey, c.user))
					}
					req.Header.Set("X-XSRF-TOKEN", "test-token")
					req.AddCookie(&http.Cookie{Name: XSRF_TOKEN, Value: "test-token"})
					if c.callCore {
						response := endpoint.response
						if c.empty {
							response = endpoint.empty
						}
						if endpoint.method == "GetRemediationProposal" {
							aiMock.On(endpoint.method, mock.Anything, mock.MatchedBy(func(input *ai.GetRemediationProposalRequest) bool {
								return input.ProjectId == 1001 && input.RemediationProposalId == 10
							})).Return(response, c.coreErr).Once()
						} else {
							aiMock.On(endpoint.method, mock.Anything, mock.MatchedBy(func(input *ai.ListRemediationProposalRequest) bool {
								return input.ProjectId == 1001 && input.FindingId == 20 && input.Status == req.URL.Query().Get("status")
							})).Return(response, c.coreErr).Once()
						}
					}
					rec := httptest.NewRecorder()
					newRouter(&svc).ServeHTTP(rec, req)
					require.Equal(t, c.wantStatus, rec.Code, rec.Body.String())
					if c.wantStatus == http.StatusOK {
						want := endpoint.wantJSON
						if c.empty {
							want = `{"data":{}}`
						}
						require.JSONEq(t, want, rec.Body.String())
					}
				})
			}
		})
	}
}
