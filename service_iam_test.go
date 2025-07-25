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

func TestPutUserHandler(t *testing.T) {
	iamMock := iammocks.NewIAMServiceClient(t)
	svc := gatewayService{
		iamClient: iamMock,
	}
	cases := []struct {
		name       string
		input      string
		claims     *jwt.MapClaims
		userIdpKey string
		claimsErr  error
		mockResp   *iam.PutUserResponse
		mockErr    error
		wantStatus int
	}{
		{
			name:  "OK",
			input: `{"user": {"sub":"xxx", "name":"nm1", "activated":true}}`,
			claims: &jwt.MapClaims{
				"user_idp_key": "uik",
			},
			userIdpKey: "uik",
			mockResp:   &iam.PutUserResponse{},
			wantStatus: http.StatusOK,
		},
		{
			name:  "NG Invalid parameter",
			input: `invalid_param`,
			claims: &jwt.MapClaims{
				"user_idp_key": "uik",
			},
			userIdpKey: "uik",
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "NG verifying token error",
			input:      `{"user": {"sub":"xxx", "name":"nm2", "activated":true}}`,
			claimsErr:  errors.New("something error"),
			wantStatus: http.StatusForbidden,
		},
		{
			name:       "NG userIdpKey is empty",
			input:      `{"user": {"sub":"xxx", "name":"nm3", "activated":true}}`,
			claims:     &jwt.MapClaims{},
			userIdpKey: "",
			wantStatus: http.StatusForbidden,
		},
		{
			name:  "NG Backend service error",
			input: `{"user": {"sub":"xxx", "name":"nm4", "activated":true}}`,
			claims: &jwt.MapClaims{
				"user_idp_key": "uik",
			},
			userIdpKey: "uik",
			wantStatus: http.StatusInternalServerError,
			mockErr:    errors.New("something wrong"),
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			svc.claimsClient = newMockClaimsClient(c.claims, "", c.userIdpKey, c.claimsErr)
			if c.mockResp != nil || c.mockErr != nil {
				iamMock.On("PutUser", mock.Anything, mock.Anything).Return(c.mockResp, c.mockErr).Once()
			}
			// Invoke HTTP Request
			rec := httptest.NewRecorder()
			req, _ := http.NewRequest(http.MethodPost, "/api/v1/iam/put-user/", strings.NewReader(c.input))
			req = req.WithContext(context.WithValue(req.Context(), userKey, &requestUser{sub: "xxx", userID: 1}))
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
