package main

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/ca-risken/core/proto/finding"
	findingmocks "github.com/ca-risken/core/proto/finding/mocks"
	"github.com/golang-jwt/jwt/v4"
	"github.com/stretchr/testify/mock"
)

func TestPutPendFindingHandler(t *testing.T) {
	findingMock := findingmocks.NewFindingServiceClient(t)
	svc := gatewayService{
		findingClient: findingMock,
	}
	cases := []struct {
		name       string
		input      string
		claims     *jwt.MapClaims
		userKey    *requestUser
		userIdpKey string
		claimsErr  error
		mockResp   *finding.PutPendFindingResponse
		mockErr    error
		wantStatus int
	}{
		{
			name:  "OK userID is set",
			input: `{"project_id": 1, "pend_finding":{"finding_id":1, "project_id":1}}`,
			claims: &jwt.MapClaims{
				"user_idp_key": "uik",
			},
			userKey:    &requestUser{sub: "uik", userID: 1},
			userIdpKey: "uik",
			mockResp:   &finding.PutPendFindingResponse{},
			wantStatus: http.StatusOK,
		},
		{
			name:  "OK userID is not set",
			input: `{"project_id": 1, "pend_finding":{"finding_id":1, "project_id":1}}`,
			claims: &jwt.MapClaims{
				"user_idp_key": "uik",
			},
			userKey:    &requestUser{sub: "uik", accessTokenID: 1},
			userIdpKey: "uik",
			mockResp:   &finding.PutPendFindingResponse{},
			wantStatus: http.StatusOK,
		},
		{
			name:  "NG Invalid parameter",
			input: `invalid_param`,
			claims: &jwt.MapClaims{
				"user_idp_key": "uik",
			},
			userKey:    &requestUser{sub: "uik", userID: 1},
			userIdpKey: "uik",
			wantStatus: http.StatusBadRequest,
		},
		{
			name:  "NG InvalidUser",
			input: `{"project_id": 1, "pend_finding":{"finding_id":1, "project_id":1}}`,
			claims: &jwt.MapClaims{
				"user_idp_key": "uik",
			},
			userKey:    nil,
			userIdpKey: "uik",
			wantStatus: http.StatusUnauthorized,
		},
		{
			name:  "NG Backend service error",
			input: `{"project_id": 1, "pend_finding":{"finding_id":1, "project_id":1}}`,
			claims: &jwt.MapClaims{
				"user_idp_key": "uik",
			},
			userKey:    &requestUser{sub: "uik", userID: 1},
			userIdpKey: "uik",
			wantStatus: http.StatusInternalServerError,
			mockErr:    errors.New("something wrong"),
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			svc.claimsClient = newMockClaimsClient(c.claims, "", c.userIdpKey, c.claimsErr)
			if c.mockResp != nil || c.mockErr != nil {
				findingMock.On("PutPendFinding", mock.Anything, mock.Anything).Return(c.mockResp, c.mockErr).Once()
			}
			// Invoke HTTP Request
			rec := httptest.NewRecorder()
			req, _ := http.NewRequest(http.MethodPost, "/api/v1/finding/put-pend-finding/", strings.NewReader(c.input))
			req = req.WithContext(context.WithValue(req.Context(), userKey, c.userKey))
			req.Header.Add("Content-Type", "application/json")
			svc.putPendFindingHandler(rec, req)
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

func TestGetAISummaryStreamHandler_DisablesWriteTimeout(t *testing.T) {
	findingMock := findingmocks.NewFindingServiceClient(t)
	streamMock := findingmocks.NewFindingService_GetAISummaryStreamClient(t)
	svc := gatewayService{
		findingClient: findingMock,
	}

	findingMock.On("GetAISummaryStream", mock.Anything, mock.Anything).Return(streamMock, nil).Once()
	streamMock.On("Recv").Return(&finding.GetAISummaryResponse{Answer: "first"}, nil).Once()
	streamMock.On("Recv").Run(func(args mock.Arguments) {
		time.Sleep(80 * time.Millisecond)
	}).Return(&finding.GetAISummaryResponse{Answer: "second"}, nil).Once()
	streamMock.On("Recv").Return((*finding.GetAISummaryResponse)(nil), io.EOF).Once()

	server := httptest.NewUnstartedServer(http.HandlerFunc(svc.getAISummaryStreamHandler))
	server.Config.WriteTimeout = 50 * time.Millisecond
	server.Start()
	defer server.Close()

	resp, err := http.Get(server.URL + "?project_id=1&finding_id=1&lang=ja")
	if err != nil {
		t.Fatalf("failed to request SSE endpoint: err=%+v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("failed to read SSE response body: err=%+v", err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("unexpected HTTP status code: got=%d, body=%q", resp.StatusCode, string(body))
	}
	if got := string(body); got != "firstsecond" {
		t.Fatalf("unexpected SSE body: got=%q", got)
	}
}
