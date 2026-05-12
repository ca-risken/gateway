package main

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/ca-risken/core/proto/iam"
	iammocks "github.com/ca-risken/core/proto/iam/mocks"
	"github.com/golang-jwt/jwt/v4"
	"github.com/stretchr/testify/mock"
)

type countingClaimsClient struct {
	claims         *jwt.MapClaims
	userIdpKey     string
	err            error
	getClaimsCalls int
}

func (c *countingClaimsClient) getClaims(ctx context.Context, tokenString string) (*jwt.MapClaims, error) {
	c.getClaimsCalls++
	return c.claims, c.err
}

func (c *countingClaimsClient) getUserName(claims *jwt.MapClaims) string {
	return ""
}

func (c *countingClaimsClient) getUserIdpKey(claims *jwt.MapClaims) string {
	return c.userIdpKey
}

func TestPutUserHandler(t *testing.T) {
	cases := []struct {
		name            string
		input           string
		user            *requestUser
		claims          *jwt.MapClaims
		userIdpKey      string
		claimsErr       error
		mockResp        *iam.PutUserResponse
		mockErr         error
		wantStatus      int
		wantClaimsCalls int
		wantPutUser     bool
		wantSub         string
		wantName        string
		wantActivated   bool
		wantUserIdpKey  string
	}{
		{
			name:  "OK ignores attacker controlled sub",
			input: `{"user": {"sub":"attacker-controlled-sub", "name":"nm1", "activated":true}}`,
			user:  &requestUser{sub: "authenticated-sub", userID: 1},
			claims: &jwt.MapClaims{
				"user_idp_key": "uik",
			},
			userIdpKey:      "uik",
			mockResp:        &iam.PutUserResponse{},
			wantStatus:      http.StatusOK,
			wantClaimsCalls: 1,
			wantPutUser:     true,
			wantSub:         "authenticated-sub",
			wantName:        "nm1",
			wantActivated:   true,
			wantUserIdpKey:  "uik",
		},
		{
			name:  "OK keeps existing request format",
			input: `{"user": {"sub":"authenticated-sub", "name":"nm2", "activated":true}}`,
			user:  &requestUser{sub: "authenticated-sub", userID: 1},
			claims: &jwt.MapClaims{
				"user_idp_key": "uik",
			},
			userIdpKey:      "uik",
			mockResp:        &iam.PutUserResponse{},
			wantStatus:      http.StatusOK,
			wantClaimsCalls: 1,
			wantPutUser:     true,
			wantSub:         "authenticated-sub",
			wantName:        "nm2",
			wantActivated:   true,
			wantUserIdpKey:  "uik",
		},
		{
			name:  "OK empty sub in body is ignored",
			input: `{"user": {"sub":"", "name":"nm3", "activated":true}}`,
			user:  &requestUser{sub: "authenticated-sub", userID: 1},
			claims: &jwt.MapClaims{
				"user_idp_key": "uik",
			},
			userIdpKey:      "uik",
			mockResp:        &iam.PutUserResponse{},
			wantStatus:      http.StatusOK,
			wantClaimsCalls: 1,
			wantPutUser:     true,
			wantSub:         "authenticated-sub",
			wantName:        "nm3",
			wantActivated:   true,
			wantUserIdpKey:  "uik",
		},
		{
			name:            "NG user is missing",
			input:           `{}`,
			user:            &requestUser{sub: "authenticated-sub", userID: 1},
			wantStatus:      http.StatusBadRequest,
			wantClaimsCalls: 0,
		},
		{
			name:            "NG user is null",
			input:           `{"user":null}`,
			user:            &requestUser{sub: "authenticated-sub", userID: 1},
			wantStatus:      http.StatusBadRequest,
			wantClaimsCalls: 0,
		},
		{
			name:            "NG Invalid parameter",
			input:           `invalid_param`,
			user:            &requestUser{sub: "authenticated-sub", userID: 1},
			wantStatus:      http.StatusBadRequest,
			wantClaimsCalls: 0,
		},
		{
			name:            "NG Invalid user",
			input:           `{"user": {"sub":"victim-sub", "name":"nm4", "activated":true}}`,
			user:            &requestUser{userID: 1},
			wantStatus:      http.StatusUnauthorized,
			wantClaimsCalls: 0,
		},
		{
			name:            "NG verifying token error",
			input:           `{"user": {"sub":"victim-sub", "name":"nm5", "activated":true}}`,
			user:            &requestUser{sub: "authenticated-sub", userID: 1},
			claimsErr:       errors.New("something error"),
			wantStatus:      http.StatusForbidden,
			wantClaimsCalls: 1,
		},
		{
			name:            "NG userIdpKey is empty",
			input:           `{"user": {"sub":"victim-sub", "name":"nm6", "activated":true}}`,
			user:            &requestUser{sub: "authenticated-sub", userID: 1},
			claims:          &jwt.MapClaims{},
			userIdpKey:      "",
			wantStatus:      http.StatusForbidden,
			wantClaimsCalls: 1,
		},
		{
			name:  "NG Backend service error",
			input: `{"user": {"sub":"victim-sub", "name":"nm7", "activated":true}}`,
			user:  &requestUser{sub: "authenticated-sub", userID: 1},
			claims: &jwt.MapClaims{
				"user_idp_key": "uik",
			},
			userIdpKey:      "uik",
			wantStatus:      http.StatusInternalServerError,
			wantClaimsCalls: 1,
			wantPutUser:     true,
			wantSub:         "authenticated-sub",
			wantName:        "nm7",
			wantActivated:   true,
			wantUserIdpKey:  "uik",
			mockErr:         errors.New("something wrong"),
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			iamMock := iammocks.NewIAMServiceClient(t)
			claimsMock := &countingClaimsClient{
				claims:     c.claims,
				userIdpKey: c.userIdpKey,
				err:        c.claimsErr,
			}
			svc := gatewayService{
				iamClient:    iamMock,
				claimsClient: claimsMock,
			}
			if c.wantPutUser {
				iamMock.On("PutUser", mock.Anything, mock.MatchedBy(func(req *iam.PutUserRequest) bool {
					return req != nil &&
						req.GetUser() != nil &&
						req.GetUser().GetSub() == c.wantSub &&
						req.GetUser().GetName() == c.wantName &&
						req.GetUser().GetActivated() == c.wantActivated &&
						req.GetUser().GetUserIdpKey() == c.wantUserIdpKey
				})).Return(c.mockResp, c.mockErr).Once()
			}
			// Invoke HTTP Request
			rec := httptest.NewRecorder()
			req, _ := http.NewRequest(http.MethodPost, "/api/v1/iam/put-user/", strings.NewReader(c.input))
			req = req.WithContext(context.WithValue(req.Context(), userKey, c.user))
			req.Header.Add("Content-Type", "application/json")
			svc.putUserHandler(rec, req)
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
			if claimsMock.getClaimsCalls != c.wantClaimsCalls {
				t.Fatalf("Unexpected claims call count: want=%d, got=%d", c.wantClaimsCalls, claimsMock.getClaimsCalls)
			}
			if !c.wantPutUser {
				iamMock.AssertNotCalled(t, "PutUser", mock.Anything, mock.Anything)
			}
		})
	}
}

func TestUpdateAccessTokenHandler(t *testing.T) {
	iamMock := iammocks.NewIAMServiceClient(t)
	svc := gatewayService{
		iamClient: iamMock,
	}
	cases := []struct {
		name       string
		input      string
		mockResp   *iam.PutAccessTokenResponse
		mockErr    error
		wantStatus int
	}{
		{
			name:       "OK",
			input:      `{"project_id":1, "access_token": {"project_id":1, "access_token_id":1, "name":"nm", "last_updated_user_id":1}}`,
			mockResp:   &iam.PutAccessTokenResponse{},
			wantStatus: http.StatusOK,
		},
		{
			name:       "NG Invalid parameter",
			input:      `invalid_param`,
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "NG Backend service error",
			input:      `{"project_id":1, "access_token": {"project_id":1, "access_token_id":1, "name":"nm", "expired_at": 1, "last_updated_user_id":1}}`,
			wantStatus: http.StatusInternalServerError,
			mockErr:    errors.New("something wrong"),
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if c.mockResp != nil || c.mockErr != nil {
				iamMock.On("PutAccessToken", mock.Anything, mock.Anything).Return(c.mockResp, c.mockErr).Once()
			}
			rec := httptest.NewRecorder()
			req, _ := http.NewRequest(http.MethodPost, "/api/v1/iam/update-access-token/", strings.NewReader(c.input))
			req = req.WithContext(context.WithValue(req.Context(), userKey, &requestUser{userID: 1}))
			req.Header.Add("Content-Type", "application/json")
			svc.updateAccessTokenHandler(rec, req)
			if c.wantStatus != rec.Code {
				t.Fatalf("Unexpected HTTP status code: want=%+v, got=%+v", c.wantStatus, rec.Code)
			}
			resp := map[string]interface{}{}
			if err := json.NewDecoder(rec.Body).Decode(&resp); err != nil {
				t.Fatalf("Unexpected json decode error to response body: err=%+v", err)
			}
			// sappLogger.Info(resp)
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
