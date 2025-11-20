package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/ca-risken/core/proto/organization_iam"
	organization_iammocks "github.com/ca-risken/core/proto/organization_iam/mocks"
	"github.com/stretchr/testify/mock"
)

func TestGenerateOrganizationAccessTokenHandler(t *testing.T) {
	orgIAMMock := organization_iammocks.NewOrganizationIAMServiceClient(t)
	svc := gatewayService{organization_iamClient: orgIAMMock}

	cases := []struct {
		name       string
		input      string
		mockSetup  func()
		wantStatus int
	}{
		{
			name:  "OK",
			input: `{"organization_id":1, "name":"token-A"}`,
			mockSetup: func() {
				orgIAMMock.On("ListOrganizationAccessToken", mock.Anything, mock.MatchedBy(func(req *organization_iam.ListOrganizationAccessTokenRequest) bool {
					return req.OrganizationId == 1 && req.Name == "token-A"
				})).Return(&organization_iam.ListOrganizationAccessTokenResponse{}, nil).Once()
				orgIAMMock.On("PutOrganizationAccessToken", mock.Anything, mock.MatchedBy(func(req *organization_iam.PutOrganizationAccessTokenRequest) bool {
					return req.OrganizationId == 1 && req.LastUpdatedUserId == 1 && req.AccessTokenId == 0 && req.PlainTextToken != ""
				})).Return(&organization_iam.PutOrganizationAccessTokenResponse{
					AccessToken: &organization_iam.OrganizationAccessToken{
						OrganizationId: 1,
						AccessTokenId:  10,
					},
				}, nil).Once()
			},
			wantStatus: http.StatusOK,
		},
		{
			name:  "NG duplicate name",
			input: `{"organization_id":1, "name":"token-A"}`,
			mockSetup: func() {
				orgIAMMock.On("ListOrganizationAccessToken", mock.Anything, mock.Anything).Return(&organization_iam.ListOrganizationAccessTokenResponse{
					AccessToken: []*organization_iam.OrganizationAccessToken{{}},
				}, nil).Once()
			},
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "NG invalid parameter",
			input:      `invalid_json`,
			wantStatus: http.StatusBadRequest,
		},
		{
			name:  "NG backend error",
			input: `{"organization_id":1, "name":"token-A"}`,
			mockSetup: func() {
				orgIAMMock.On("ListOrganizationAccessToken", mock.Anything, mock.Anything).Return(&organization_iam.ListOrganizationAccessTokenResponse{}, nil).Once()
				orgIAMMock.On("PutOrganizationAccessToken", mock.Anything, mock.Anything).Return(nil, errors.New("something wrong")).Once()
			},
			wantStatus: http.StatusInternalServerError,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if c.mockSetup != nil {
				c.mockSetup()
			}
			rec := httptest.NewRecorder()
			req, _ := http.NewRequest(http.MethodPost, "/api/v1/organization/generate-organization-access-token/", strings.NewReader(c.input))
			req = req.WithContext(context.WithValue(req.Context(), userKey, &requestUser{userID: 1}))
			req.Header.Add("Content-Type", "application/json")
			svc.generateOrganizationAccessTokenOrganization_iamHandler(rec, req)

			if rec.Code != c.wantStatus {
				t.Fatalf("Unexpected HTTP status code: want=%d, got=%d", c.wantStatus, rec.Code)
			}
			resp := map[string]interface{}{}
			if err := json.NewDecoder(rec.Body).Decode(&resp); err != nil {
				t.Fatalf("Unexpected json decode error to response body: err=%+v", err)
			}
			key := successJSONKey
			if c.wantStatus != http.StatusOK {
				key = errorJSONKey
			}
			if _, ok := resp[key]; !ok {
				t.Fatalf("Unexpected no response key: want key=%s", key)
			}
		})
	}
}

func TestUpdateOrganizationAccessTokenHandler(t *testing.T) {
	orgIAMMock := organization_iammocks.NewOrganizationIAMServiceClient(t)
	svc := gatewayService{organization_iamClient: orgIAMMock}

	cases := []struct {
		name       string
		input      string
		mockSetup  func()
		wantStatus int
	}{
		{
			name:  "OK",
			input: `{"organization_id":1, "access_token_id":1, "name":"token-A"}`,
			mockSetup: func() {
				orgIAMMock.On("PutOrganizationAccessToken", mock.Anything, mock.MatchedBy(func(req *organization_iam.PutOrganizationAccessTokenRequest) bool {
					return req.OrganizationId == 1 && req.AccessTokenId == 1 && req.LastUpdatedUserId == 1
				})).Return(&organization_iam.PutOrganizationAccessTokenResponse{}, nil).Once()
			},
			wantStatus: http.StatusOK,
		},
		{
			name:       "NG invalid parameter",
			input:      `invalid_json`,
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "NG missing access_token_id",
			input:      `{"organization_id":1, "name":"token-A"}`,
			wantStatus: http.StatusBadRequest,
		},
		{
			name:  "NG backend error",
			input: `{"organization_id":1, "access_token_id":1, "name":"token-A"}`,
			mockSetup: func() {
				orgIAMMock.On("PutOrganizationAccessToken", mock.Anything, mock.Anything).Return(nil, errors.New("something wrong")).Once()
			},
			wantStatus: http.StatusInternalServerError,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if c.mockSetup != nil {
				c.mockSetup()
			}
			rec := httptest.NewRecorder()
			req, _ := http.NewRequest(http.MethodPost, "/api/v1/organization/update-organization-access-token/", strings.NewReader(c.input))
			req = req.WithContext(context.WithValue(req.Context(), userKey, &requestUser{userID: 1}))
			req.Header.Add("Content-Type", "application/json")
			svc.updateOrganizationAccessTokenOrganization_iamHandler(rec, req)

			if rec.Code != c.wantStatus {
				t.Fatalf("Unexpected HTTP status code: want=%d, got=%d", c.wantStatus, rec.Code)
			}
			resp := map[string]interface{}{}
			if err := json.NewDecoder(rec.Body).Decode(&resp); err != nil {
				t.Fatalf("Unexpected json decode error to response body: err=%+v", err)
			}
			key := successJSONKey
			if c.wantStatus != http.StatusOK {
				key = errorJSONKey
			}
			if _, ok := resp[key]; !ok {
				t.Fatalf("Unexpected no response key: want key=%s", key)
			}
		})
	}
}

func TestAuthenticateOrganizationAccessTokenHandler(t *testing.T) {
	orgIAMMock := organization_iammocks.NewOrganizationIAMServiceClient(t)
	svc := gatewayService{organization_iamClient: orgIAMMock}

	validToken := encodeOrganizationAccessToken(1, 2, "plain")

	cases := []struct {
		name       string
		input      string
		mockSetup  func()
		wantStatus int
	}{
		{
			name:  "OK",
			input: fmt.Sprintf(`{"access_token":%q}`, validToken),
			mockSetup: func() {
				orgIAMMock.On("AuthenticateOrganizationAccessToken", mock.Anything, mock.MatchedBy(func(req *organization_iam.AuthenticateOrganizationAccessTokenRequest) bool {
					return req.OrganizationId == 1 && req.AccessTokenId == 2 && req.PlainTextToken == "plain"
				})).Return(&organization_iam.AuthenticateOrganizationAccessTokenResponse{}, nil).Once()
			},
			wantStatus: http.StatusOK,
		},
		{
			name:       "NG invalid parameter",
			input:      `{}`,
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "NG invalid token format",
			input:      `{"access_token":"invalid_token"}`,
			wantStatus: http.StatusBadRequest,
		},
		{
			name:  "NG backend error",
			input: fmt.Sprintf(`{"access_token":%q}`, validToken),
			mockSetup: func() {
				orgIAMMock.On("AuthenticateOrganizationAccessToken", mock.Anything, mock.Anything).Return(nil, errors.New("something wrong")).Once()
			},
			wantStatus: http.StatusInternalServerError,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if c.mockSetup != nil {
				c.mockSetup()
			}
			rec := httptest.NewRecorder()
			req, _ := http.NewRequest(http.MethodPost, "/api/v1/organization/authenticate-organization-access-token/", strings.NewReader(c.input))
			req.Header.Add("Content-Type", "application/json")
			svc.authenticateOrganizationAccessTokenOrganization_iamHandler(rec, req)

			if rec.Code != c.wantStatus {
				t.Fatalf("Unexpected HTTP status code: want=%d, got=%d", c.wantStatus, rec.Code)
			}
			resp := map[string]interface{}{}
			if err := json.NewDecoder(rec.Body).Decode(&resp); err != nil {
				t.Fatalf("Unexpected json decode error to response body: err=%+v", err)
			}
			key := successJSONKey
			if c.wantStatus != http.StatusOK {
				key = errorJSONKey
			}
			if _, ok := resp[key]; !ok {
				t.Fatalf("Unexpected no response key: want key=%s", key)
			}
		})
	}
}
