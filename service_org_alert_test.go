package main

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/ca-risken/core/proto/org_alert"
	orgalertmocks "github.com/ca-risken/core/proto/org_alert/mocks"
	"github.com/ca-risken/core/proto/org_iam"
	orgiammocks "github.com/ca-risken/core/proto/org_iam/mocks"
	"github.com/stretchr/testify/mock"
)

func TestListOrgAlertCondNotificationHandler(t *testing.T) {
	cases := []struct {
		name       string
		query      string
		mockErr    error
		wantCall   bool
		wantStatus int
	}{
		{
			name:       "OK four key filter",
			query:      "organization_id=1&project_id=2&alert_condition_id=3&notification_id=4",
			wantCall:   true,
			wantStatus: http.StatusOK,
		},
		{
			name:       "NG invalid organization ID",
			query:      "organization_id=0",
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "NG backend error",
			query:      "organization_id=1&project_id=2&alert_condition_id=3&notification_id=4",
			mockErr:    errors.New("something wrong"),
			wantCall:   true,
			wantStatus: http.StatusInternalServerError,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			orgAlertMock := orgalertmocks.NewOrgAlertServiceClient(t)
			if c.wantCall {
				orgAlertMock.On("ListOrgAlertCondNotification", mock.Anything, mock.MatchedBy(func(req *org_alert.ListOrgAlertCondNotificationRequest) bool {
					return req.GetOrganizationId() == 1 && req.GetProjectId() == 2 &&
						req.GetAlertConditionId() == 3 && req.GetNotificationId() == 4
				})).Return(&org_alert.ListOrgAlertCondNotificationResponse{}, c.mockErr).Once()
			}
			svc := gatewayService{org_alertClient: orgAlertMock}
			rec := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodGet, "/api/v1/organization-alert/list-alert-cond-notification?"+c.query, nil)

			svc.listOrgAlertCondNotificationOrg_alertHandler(rec, req)

			if rec.Code != c.wantStatus {
				t.Fatalf("Unexpected HTTP status code: want=%d, got=%d", c.wantStatus, rec.Code)
			}
		})
	}
}

func TestUpdateOrgAlertCondNotificationCacheHandler(t *testing.T) {
	cases := []struct {
		name       string
		body       string
		mockErr    error
		wantCall   bool
		wantStatus int
	}{
		{
			name:       "OK four key update",
			body:       `{"organization_id":1,"project_id":2,"alert_condition_id":3,"notification_id":4,"cache_second":900}`,
			wantCall:   true,
			wantStatus: http.StatusOK,
		},
		{
			name:       "NG missing notification ID",
			body:       `{"organization_id":1,"project_id":2,"alert_condition_id":3,"cache_second":900}`,
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "NG backend error",
			body:       `{"organization_id":1,"project_id":2,"alert_condition_id":3,"notification_id":4,"cache_second":900}`,
			mockErr:    errors.New("something wrong"),
			wantCall:   true,
			wantStatus: http.StatusInternalServerError,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			orgAlertMock := orgalertmocks.NewOrgAlertServiceClient(t)
			if c.wantCall {
				orgAlertMock.On("UpdateOrgAlertCondNotificationCache", mock.Anything, mock.MatchedBy(func(req *org_alert.UpdateOrgAlertCondNotificationCacheRequest) bool {
					return req.GetOrganizationId() == 1 && req.GetProjectId() == 2 &&
						req.GetAlertConditionId() == 3 && req.GetNotificationId() == 4 && req.GetCacheSecond() == 900
				})).Return(&org_alert.UpdateOrgAlertCondNotificationCacheResponse{}, c.mockErr).Once()
			}
			svc := gatewayService{org_alertClient: orgAlertMock}
			rec := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodPost, "/api/v1/organization-alert/update-alert-cond-notification-cache", strings.NewReader(c.body))
			req.Header.Set("Content-Type", contenTypeJSON)

			svc.updateOrgAlertCondNotificationCacheOrg_alertHandler(rec, req)

			if rec.Code != c.wantStatus {
				t.Fatalf("Unexpected HTTP status code: want=%d, got=%d", c.wantStatus, rec.Code)
			}
		})
	}
}

func TestOrgAlertCondNotificationRoutes(t *testing.T) {
	cases := []struct {
		name       string
		method     string
		path       string
		body       string
		authorized bool
		wantStatus int
	}{
		{
			name:       "list is protected by organization authorization",
			method:     http.MethodGet,
			path:       "/api/v1/organization-alert/list-alert-cond-notification?organization_id=1",
			wantStatus: http.StatusUnauthorized,
		},
		{
			name:       "cache update is protected by organization authorization",
			method:     http.MethodPost,
			path:       "/api/v1/organization-alert/update-alert-cond-notification-cache",
			body:       `{"organization_id":1,"project_id":2,"alert_condition_id":3,"notification_id":4,"cache_second":900}`,
			wantStatus: http.StatusUnauthorized,
		},
		{
			name:       "generated get handler is not exposed",
			method:     http.MethodGet,
			path:       "/api/v1/organization-alert/get-alert-cond-notification?organization_id=1&project_id=2&alert_condition_id=3&notification_id=4",
			authorized: true,
			wantStatus: http.StatusNotFound,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			svc := &gatewayService{}
			rec := httptest.NewRecorder()
			req := httptest.NewRequest(c.method, c.path, strings.NewReader(c.body))
			if c.body != "" {
				req.Header.Set("Content-Type", contenTypeJSON)
			}
			if c.authorized {
				orgIAMMock := orgiammocks.NewOrgIAMServiceClient(t)
				orgIAMMock.On("IsAuthorizedOrgToken", mock.Anything, mock.Anything).Return(&org_iam.IsAuthorizedOrgTokenResponse{Ok: true}, nil).Once()
				svc.org_iamClient = orgIAMMock
				req = req.WithContext(context.WithValue(req.Context(), userKey, &requestUser{orgAccessTokenID: 1, orgAccessTokenOrgID: 1}))
			}

			newRouter(svc).ServeHTTP(rec, req)

			if rec.Code != c.wantStatus {
				t.Fatalf("Unexpected HTTP status code: want=%d, got=%d", c.wantStatus, rec.Code)
			}
		})
	}
}
