package main

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/ca-risken/core/proto/ai"
	aimocks "github.com/ca-risken/core/proto/ai/mocks"
	"github.com/ca-risken/core/proto/iam"
	iammocks "github.com/ca-risken/core/proto/iam/mocks"
	"github.com/ca-risken/datasource-api/proto/datasource_ai"
	datasource_aimocks "github.com/ca-risken/datasource-api/proto/datasource_ai/mocks"
	"github.com/stretchr/testify/mock"
)

func TestChatAIAiHandler(t *testing.T) {
	cases := []struct {
		name       string
		inputUser  *requestUser
		inputBody  string
		setupMocks func(*aimocks.AIServiceClient, *iammocks.IAMServiceClient)
		wantStatus int
	}{
		{
			name:      "OK human access",
			inputUser: &requestUser{sub: "sub", userID: 1},
			inputBody: `{"question":"test","project_id":1001}`,
			setupMocks: func(aiMock *aimocks.AIServiceClient, iamMock *iammocks.IAMServiceClient) {
				iamMock.On("IsAuthorized", mock.Anything, mock.MatchedBy(func(req *iam.IsAuthorizedRequest) bool {
					return req.UserId == 1 &&
						req.ProjectId == 1001 &&
						req.ActionName == "ai/chat-ai" &&
						req.ResourceName == "ai/resource_any"
				})).Return(&iam.IsAuthorizedResponse{Ok: true}, nil).Once()
				aiMock.On("ChatAI", mock.Anything, mock.MatchedBy(func(req *ai.ChatAIRequest) bool {
					return req.Question == "test" && req.ProjectId == 1001
				})).Return(&ai.ChatAIResponse{Answer: "ok"}, nil).Once()
			},
			wantStatus: http.StatusOK,
		},
		{
			name:      "OK token access forces authorized project",
			inputUser: &requestUser{accessTokenID: 10, accessTokenProjectID: 2004},
			inputBody: `{"question":"test","project_id":9999}`,
			setupMocks: func(aiMock *aimocks.AIServiceClient, iamMock *iammocks.IAMServiceClient) {
				iamMock.On("IsAuthorizedToken", mock.Anything, mock.MatchedBy(func(req *iam.IsAuthorizedTokenRequest) bool {
					return req.AccessTokenId == 10 &&
						req.ProjectId == 2004 &&
						req.ActionName == "ai/chat-ai" &&
						req.ResourceName == "ai/resource_any"
				})).Return(&iam.IsAuthorizedTokenResponse{Ok: true}, nil).Once()
				aiMock.On("ChatAI", mock.Anything, mock.MatchedBy(func(req *ai.ChatAIRequest) bool {
					return req.Question == "test" && req.ProjectId == 2004
				})).Return(&ai.ChatAIResponse{Answer: "ok"}, nil).Once()
			},
			wantStatus: http.StatusOK,
		},
		{
			name:      "NG unauthorized human access",
			inputUser: &requestUser{sub: "sub", userID: 1},
			inputBody: `{"question":"test","project_id":1001}`,
			setupMocks: func(_ *aimocks.AIServiceClient, iamMock *iammocks.IAMServiceClient) {
				iamMock.On("IsAuthorized", mock.Anything, mock.Anything).Return(&iam.IsAuthorizedResponse{Ok: false}, nil).Once()
			},
			wantStatus: http.StatusForbidden,
		},
		{
			name:      "OK human access without project id",
			inputUser: &requestUser{sub: "sub", userID: 1},
			inputBody: `{"question":"test"}`,
			setupMocks: func(aiMock *aimocks.AIServiceClient, _ *iammocks.IAMServiceClient) {
				aiMock.On("ChatAI", mock.Anything, mock.MatchedBy(func(req *ai.ChatAIRequest) bool {
					return req.Question == "test" && req.ProjectId == 0
				})).Return(&ai.ChatAIResponse{Answer: "ok"}, nil).Once()
			},
			wantStatus: http.StatusOK,
		},
		{
			name:       "NG organization token access",
			inputUser:  &requestUser{orgAccessTokenID: 10, orgAccessTokenOrgID: 20},
			inputBody:  `{"question":"test","project_id":1001}`,
			wantStatus: http.StatusForbidden,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			aiMock := aimocks.NewAIServiceClient(t)
			iamMock := iammocks.NewIAMServiceClient(t)
			svc := gatewayService{
				aiClient:  aiMock,
				iamClient: iamMock,
			}
			if c.setupMocks != nil {
				c.setupMocks(aiMock, iamMock)
			}

			rec := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodPost, "/api/v1/ai/chat-ai", strings.NewReader(c.inputBody))
			req.Header.Add("Content-Type", "application/json")
			req = req.WithContext(context.WithValue(req.Context(), userKey, c.inputUser))

			svc.chatAIAiHandler(rec, req)

			if rec.Code != c.wantStatus {
				t.Fatalf("Unexpected HTTP status code: want=%d, got=%d", c.wantStatus, rec.Code)
			}
			if c.wantStatus == http.StatusOK {
				resp := map[string]interface{}{}
				if err := json.NewDecoder(rec.Body).Decode(&resp); err != nil {
					t.Fatalf("Unexpected json decode error to response body: err=%+v", err)
				}
				if _, ok := resp[successJSONKey]; !ok {
					t.Fatalf("Unexpected no response key: want key=%s", successJSONKey)
				}
			}
		})
	}
}

func TestGenerateRemediationProposalDatasourceAIHandler(t *testing.T) {
	cases := []struct {
		name       string
		inputBody  string
		setupMocks func(*datasource_aimocks.AIServiceClient)
		wantStatus int
	}{
		{
			name:      "OK",
			inputBody: `{"project_id":1001,"finding_id":2001}`,
			setupMocks: func(aiMock *datasource_aimocks.AIServiceClient) {
				aiMock.On("GenerateRemediationProposal", mock.Anything, mock.MatchedBy(func(req *datasource_ai.GenerateRemediationProposalRequest) bool {
					return req.ProjectId == 1001 && req.FindingId == 2001
				})).Return(&datasource_ai.GenerateRemediationProposalResponse{RemediationProposalId: 3001}, nil).Once()
			},
			wantStatus: http.StatusOK,
		},
		{
			name:       "NG Invalid parameter",
			inputBody:  `{"project_id":1001}`,
			wantStatus: http.StatusBadRequest,
		},
		{
			name:      "NG Backend service error",
			inputBody: `{"project_id":1001,"finding_id":2001}`,
			setupMocks: func(aiMock *datasource_aimocks.AIServiceClient) {
				aiMock.On("GenerateRemediationProposal", mock.Anything, mock.Anything).Return(nil, errors.New("something wrong")).Once()
			},
			wantStatus: http.StatusInternalServerError,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			aiMock := datasource_aimocks.NewAIServiceClient(t)
			svc := gatewayService{datasource_aiClient: aiMock}
			if c.setupMocks != nil {
				c.setupMocks(aiMock)
			}

			rec := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodPost, "/api/v1/ai/generate-remediation-proposal", strings.NewReader(c.inputBody))

			svc.generateRemediationProposalDatasource_aiHandler(rec, req)

			if rec.Code != c.wantStatus {
				t.Fatalf("Unexpected HTTP status code: want=%d, got=%d", c.wantStatus, rec.Code)
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

func TestRemediationProposalRoutesWithProjectAuthz(t *testing.T) {
	cases := []struct {
		name       string
		method     string
		path       string
		body       string
		actionName string
		setupMocks func(*aimocks.AIServiceClient, *datasource_aimocks.AIServiceClient)
	}{
		{
			name:       "GetRemediationProposal",
			method:     http.MethodGet,
			path:       "/api/v1/ai/get-remediation-proposal?project_id=1001&remediation_proposal_id=3001",
			actionName: "ai/get-remediation-proposal",
			setupMocks: func(aiMock *aimocks.AIServiceClient, _ *datasource_aimocks.AIServiceClient) {
				aiMock.On("GetRemediationProposal", mock.Anything, mock.MatchedBy(func(req *ai.GetRemediationProposalRequest) bool {
					return req.ProjectId == 1001 && req.RemediationProposalId == 3001
				})).Return(&ai.GetRemediationProposalResponse{}, nil).Once()
			},
		},
		{
			name:       "ListRemediationProposal",
			method:     http.MethodGet,
			path:       "/api/v1/ai/list-remediation-proposal?project_id=1001&finding_id=2001",
			actionName: "ai/list-remediation-proposal",
			setupMocks: func(aiMock *aimocks.AIServiceClient, _ *datasource_aimocks.AIServiceClient) {
				aiMock.On("ListRemediationProposal", mock.Anything, mock.MatchedBy(func(req *ai.ListRemediationProposalRequest) bool {
					return req.ProjectId == 1001 && req.FindingId == 2001
				})).Return(&ai.ListRemediationProposalResponse{}, nil).Once()
			},
		},
		{
			name:       "GenerateRemediationProposal",
			method:     http.MethodPost,
			path:       "/api/v1/ai/generate-remediation-proposal",
			body:       `{"project_id":1001,"finding_id":2001}`,
			actionName: "ai/generate-remediation-proposal",
			setupMocks: func(_ *aimocks.AIServiceClient, aiMock *datasource_aimocks.AIServiceClient) {
				aiMock.On("GenerateRemediationProposal", mock.Anything, mock.MatchedBy(func(req *datasource_ai.GenerateRemediationProposalRequest) bool {
					return req.ProjectId == 1001 && req.FindingId == 2001
				})).Return(&datasource_ai.GenerateRemediationProposalResponse{RemediationProposalId: 3001}, nil).Once()
			},
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			coreAIMock := aimocks.NewAIServiceClient(t)
			datasourceAIMock := datasource_aimocks.NewAIServiceClient(t)
			iamMock := iammocks.NewIAMServiceClient(t)
			iamMock.On("IsAuthorizedToken", mock.Anything, mock.MatchedBy(func(req *iam.IsAuthorizedTokenRequest) bool {
				return req.AccessTokenId == 10 &&
					req.ProjectId == 1001 &&
					req.ActionName == c.actionName &&
					req.ResourceName == "ai/resource_any"
			})).Return(&iam.IsAuthorizedTokenResponse{Ok: true}, nil).Once()
			c.setupMocks(coreAIMock, datasourceAIMock)

			svc := gatewayService{
				aiClient:            coreAIMock,
				datasource_aiClient: datasourceAIMock,
				iamClient:           iamMock,
			}
			rec := httptest.NewRecorder()
			req := httptest.NewRequest(c.method, c.path, strings.NewReader(c.body))
			if c.method == http.MethodPost {
				req.Header.Add("Content-Type", "application/json")
			}
			req = req.WithContext(context.WithValue(req.Context(), userKey, &requestUser{
				accessTokenID:        10,
				accessTokenProjectID: 1001,
			}))

			newRouter(&svc).ServeHTTP(rec, req)

			if rec.Code != http.StatusOK {
				t.Fatalf("Unexpected HTTP status code: want=%d, got=%d", http.StatusOK, rec.Code)
			}
		})
	}
}
