package main

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/ca-risken/core/proto/iam"
	iammocks "github.com/ca-risken/core/proto/iam/mocks"
	"github.com/ca-risken/core/proto/org_iam"
	orgiammocks "github.com/ca-risken/core/proto/org_iam/mocks"
	"github.com/golang-jwt/jwt/v4"
	"github.com/stretchr/testify/mock"
)

type readCounterBody struct {
	readCount int
}

func (b *readCounterBody) Read(_ []byte) (int, error) {
	b.readCount++
	return 0, io.EOF
}

func (b *readCounterBody) Close() error {
	return nil
}

func TestSigninHandler(t *testing.T) {
	svc := gatewayService{
		sessionTimeoutSec: 1,
	}
	cases := []struct {
		name  string
		input *requestUser
		want  int
	}{
		{
			name:  "OK",
			input: &requestUser{sub: "sub", userID: 123},
			want:  http.StatusOK,
		},
		{
			name:  "NG No user",
			input: nil,
			want:  http.StatusUnauthorized,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			rec := httptest.NewRecorder()
			req, _ := http.NewRequest(http.MethodGet, "/api/v1/signin", nil)
			req = req.WithContext(context.WithValue(req.Context(), userKey, c.input))
			svc.signinHandler(rec, req)
			got := rec.Result().StatusCode
			if got != c.want {
				t.Fatalf("Unexpected responce. want=%d, got=%d", c.want, got)
			}
		})
	}
}

func TestAuthnToken(t *testing.T) {
	projectID := uint32(10)
	orgID := uint32(20)
	accessTokenID := uint32(30)
	plainToken := "plain-text"
	projectToken := "Bearer " + encodeAccessToken(projectID, accessTokenID, plainToken)
	orgToken := "Bearer " + encodeOrgAccessToken(orgID, accessTokenID, plainToken)

	cases := []struct {
		name          string
		authorization string
		expectUser    bool
		assert        func(t *testing.T, u *requestUser)
		setupMocks    func(iamMock *iammocks.IAMServiceClient, orgMock *orgiammocks.OrgIAMServiceClient)
	}{
		{
			name:          "Project token",
			authorization: projectToken,
			expectUser:    true,
			setupMocks: func(iamMock *iammocks.IAMServiceClient, _ *orgiammocks.OrgIAMServiceClient) {
				iamMock.On("AuthenticateAccessToken", mock.Anything, mock.MatchedBy(func(req *iam.AuthenticateAccessTokenRequest) bool {
					return req.ProjectId == projectID &&
						req.AccessTokenId == accessTokenID &&
						req.PlainTextToken == plainToken
				})).Return(&iam.AuthenticateAccessTokenResponse{
					AccessToken: &iam.AccessToken{
						AccessTokenId: accessTokenID,
					},
				}, nil).Once()
			},
			assert: func(t *testing.T, u *requestUser) {
				if u.accessTokenID != accessTokenID || u.accessTokenProjectID != projectID {
					t.Fatalf("unexpected project token info: %+v", u)
				}
			},
		},
		{
			name:          "Organization token",
			authorization: orgToken,
			expectUser:    true,
			setupMocks: func(_ *iammocks.IAMServiceClient, orgMock *orgiammocks.OrgIAMServiceClient) {
				orgMock.On("AuthenticateOrgAccessToken", mock.Anything, mock.MatchedBy(func(req *org_iam.AuthenticateOrgAccessTokenRequest) bool {
					return req.OrganizationId == orgID &&
						req.AccessTokenId == accessTokenID &&
						req.PlainTextToken == plainToken
				})).Return(&org_iam.AuthenticateOrgAccessTokenResponse{
					AccessToken: &org_iam.OrgAccessToken{
						AccessTokenId: accessTokenID,
					},
				}, nil).Once()
			},
			assert: func(t *testing.T, u *requestUser) {
				if u.orgAccessTokenID != accessTokenID || u.orgAccessTokenOrgID != orgID {
					t.Fatalf("unexpected org token info: %+v", u)
				}
			},
		},
		{
			name:          "Skip invalid token",
			authorization: "Bearer something-else",
		},
		{
			name: "No header",
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			iamMock := iammocks.NewIAMServiceClient(t)
			orgIAMMock := orgiammocks.NewOrgIAMServiceClient(t)
			svc := gatewayService{
				iamClient:     iamMock,
				org_iamClient: orgIAMMock,
			}
			if c.setupMocks != nil {
				c.setupMocks(iamMock, orgIAMMock)
			}

			rec := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodGet, "/api/v1/organization", nil)
			if c.authorization != "" {
				req.Header.Set("Authorization", c.authorization)
			}

			nextCalled := false
			handler := svc.authnToken(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				nextCalled = true
				u, err := getRequestUser(r)
				if c.expectUser {
					if err != nil {
						t.Fatalf("want user but got error: %+v", err)
					}
					c.assert(t, u)
				} else if err == nil {
					t.Fatalf("unexpected user found: %+v", u)
				}
			}))

			handler.ServeHTTP(rec, req)
			if !nextCalled {
				t.Fatalf("next handler was not called")
			}
		})
	}
}

func TestValidCSRFToken(t *testing.T) {
	cases := []struct {
		name        string
		inputHeader string
		inputCookie string
		want        bool
	}{
		{
			name:        "OK",
			inputHeader: "same_value",
			inputCookie: "same_value",
			want:        true,
		},
		{
			name:        "NG Header blank",
			inputHeader: "",
			inputCookie: "exists",
			want:        false,
		},
		{
			name:        "NG Cookie blank",
			inputHeader: "exists",
			inputCookie: "",
			want:        false,
		},
		{
			name:        "NG Wrong value",
			inputHeader: "wrong",
			inputCookie: "value",
			want:        false,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			req, _ := http.NewRequest(http.MethodPost, "/api/v1/test", nil)
			req.Header.Add("X-XSRF-TOKEN", c.inputHeader)
			req.AddCookie(&http.Cookie{Name: XSRF_TOKEN, Value: c.inputCookie})
			got := validCSRFToken(req)
			if got != c.want {
				t.Fatalf("Unexpected response. want=%t, got=%t", c.want, got)
			}
		})
	}
}

func TestAuthzProject(t *testing.T) {
	iamMock := iammocks.NewIAMServiceClient(t)
	svc := gatewayService{
		iamClient: iamMock,
	}
	cases := []struct {
		name         string
		inputUser    *requestUser
		inputProject string
		want         bool
		mockResp     *iam.IsAuthorizedResponse
		mockErr      error
	}{
		{
			name:         "OK",
			inputUser:    &requestUser{sub: "sub", userID: 123},
			inputProject: "project_id=1",
			mockResp:     &iam.IsAuthorizedResponse{Ok: true},
			want:         true,
		},
		{
			name:         "NG Invalid user",
			inputUser:    &requestUser{sub: "sub"},
			inputProject: "project_id=1",
			want:         false,
		},
		{
			name:         "NG Invalid project",
			inputUser:    &requestUser{sub: "sub", userID: 123},
			inputProject: "project_id=aaa",
			want:         false,
		},
		{
			name:         "NG IAM error",
			inputUser:    &requestUser{sub: "sub", userID: 123},
			inputProject: "project_id=1",
			want:         false,
			mockErr:      errors.New("something error"),
		}}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if c.mockResp != nil || c.mockErr != nil {
				iamMock.On("IsAuthorized", mock.Anything, mock.Anything).Return(c.mockResp, c.mockErr).Once()
			}
			req, _ := http.NewRequest(http.MethodGet, "/api/v1/service/action?"+c.inputProject, nil)
			got := svc.authzProject(c.inputUser, req)
			if got != c.want {
				t.Fatalf("Unexpected response. want=%t, got=%t", c.want, got)
			}
		})
	}
}

func TestAuthzProjectForToken(t *testing.T) {
	iamMock := iammocks.NewIAMServiceClient(t)
	orgIAMMock := orgiammocks.NewOrgIAMServiceClient(t)
	svc := gatewayService{
		iamClient:     iamMock,
		org_iamClient: orgIAMMock,
	}

	const projectID = uint32(1)
	const orgID = uint32(100)
	const tokenID = uint32(123)

	cases := []struct {
		name      string
		inputUser *requestUser
		query     string
		want      bool
		setup     func()
	}{
		{
			name:      "OK project token with project id",
			inputUser: &requestUser{accessTokenID: tokenID, accessTokenProjectID: projectID},
			query:     fmt.Sprintf("project_id=%d", projectID),
			want:      true,
			setup: func() {
				iamMock.On("IsAuthorizedToken", mock.Anything, mock.MatchedBy(func(req *iam.IsAuthorizedTokenRequest) bool {
					return req.ProjectId == projectID &&
						req.AccessTokenId == tokenID &&
						req.ActionName == "service/action"
				})).Return(&iam.IsAuthorizedTokenResponse{Ok: true}, nil).Once()
			},
		},
		{
			name:      "OK project token without project id",
			inputUser: &requestUser{accessTokenID: tokenID, accessTokenProjectID: projectID},
			want:      true,
			setup: func() {
				iamMock.On("IsAuthorizedToken", mock.Anything, mock.Anything).Return(&iam.IsAuthorizedTokenResponse{Ok: true}, nil).Once()
			},
		},
		{
			name:      "NG project token not provided",
			inputUser: &requestUser{sub: "sub", accessTokenProjectID: projectID},
			query:     fmt.Sprintf("project_id=%d", projectID),
			want:      false,
		},
		{
			name:      "NG project token missing project",
			inputUser: &requestUser{accessTokenID: tokenID, accessTokenProjectID: 0},
			query:     "project_id=0",
			want:      false,
		},
		{
			name:      "NG project token mismatched project id",
			inputUser: &requestUser{accessTokenID: tokenID, accessTokenProjectID: projectID},
			query:     "project_id=999",
			want:      false,
		},
		{
			name:      "NG project token IAM error",
			inputUser: &requestUser{accessTokenID: tokenID, accessTokenProjectID: projectID},
			query:     fmt.Sprintf("project_id=%d", projectID),
			want:      false,
			setup: func() {
				iamMock.On("IsAuthorizedToken", mock.Anything, mock.Anything).
					Return((*iam.IsAuthorizedTokenResponse)(nil), errors.New("something error")).Once()
			},
		},
		{
			name:      "OK org token with project id",
			inputUser: &requestUser{orgAccessTokenID: tokenID, orgAccessTokenOrgID: orgID},
			query:     fmt.Sprintf("project_id=%d", projectID),
			want:      true,
			setup: func() {
				orgIAMMock.On("IsAuthorizedOrgToken", mock.Anything, mock.MatchedBy(func(req *org_iam.IsAuthorizedOrgTokenRequest) bool {
					return req.OrganizationId == orgID &&
						req.AccessTokenId == tokenID &&
						req.ProjectId == projectID &&
						req.ActionName == "service/action"
				})).Return(&org_iam.IsAuthorizedOrgTokenResponse{Ok: true}, nil).Once()
			},
		},
		{
			name:      "NG org token missing project id",
			inputUser: &requestUser{orgAccessTokenID: tokenID, orgAccessTokenOrgID: orgID},
			want:      false,
		},
		{
			name:      "NG org token auth error",
			inputUser: &requestUser{orgAccessTokenID: tokenID, orgAccessTokenOrgID: orgID},
			query:     fmt.Sprintf("project_id=%d", projectID),
			want:      false,
			setup: func() {
				orgIAMMock.On("IsAuthorizedOrgToken", mock.Anything, mock.Anything).
					Return((*org_iam.IsAuthorizedOrgTokenResponse)(nil), errors.New("auth error")).Once()
			},
		},
		{
			name:      "NG org token unauthorized",
			inputUser: &requestUser{orgAccessTokenID: tokenID, orgAccessTokenOrgID: orgID},
			query:     fmt.Sprintf("project_id=%d", projectID),
			want:      false,
			setup: func() {
				orgIAMMock.On("IsAuthorizedOrgToken", mock.Anything, mock.Anything).
					Return(&org_iam.IsAuthorizedOrgTokenResponse{Ok: false}, nil).Once()
			},
		},
	}
	for _, c := range cases {
		c := c
		t.Run(c.name, func(t *testing.T) {
			if c.setup != nil {
				c.setup()
			}
			uri := "/api/v1/service/action"
			if c.query != "" {
				uri = uri + "?" + c.query
			}
			req, _ := http.NewRequest(http.MethodGet, uri, nil)
			got := svc.authzProjectForToken(c.inputUser, req)
			if got != c.want {
				t.Fatalf("Unexpected response. want=%t, got=%t", c.want, got)
			}
		})
	}
}

func TestAuthzMiddlewareSkipsBodyReadWhenUnauthenticated(t *testing.T) {
	svc := gatewayService{}
	cases := []struct {
		name   string
		wrap   func(http.Handler) http.Handler
		method string
		target string
	}{
		{
			name:   "project",
			wrap:   svc.authzWithProject,
			method: http.MethodPost,
			target: "/api/v1/finding/put-finding",
		},
		{
			name:   "admin",
			wrap:   svc.authzOnlyAdmin,
			method: http.MethodGet,
			target: "/api/v1/report/get-report-finding-all",
		},
		{
			name:   "organization",
			wrap:   svc.authzWithOrg,
			method: http.MethodPost,
			target: "/api/v1/organization/update-organization",
		},
		{
			name:   "project-member",
			wrap:   svc.authzWithProjectMember,
			method: http.MethodPost,
			target: "/api/v1/alert/put-alert-first-viewed-at",
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			body := &readCounterBody{}
			req, err := http.NewRequest(c.method, c.target, body)
			if err != nil {
				t.Fatalf("failed to create request: %v", err)
			}
			rec := httptest.NewRecorder()
			called := false
			handler := c.wrap(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				called = true
				w.WriteHeader(http.StatusOK)
			}))

			handler.ServeHTTP(rec, req)

			if rec.Code != http.StatusUnauthorized {
				t.Fatalf("unexpected status: want=%d got=%d", http.StatusUnauthorized, rec.Code)
			}
			if called {
				t.Fatal("next handler should not be called")
			}
			if body.readCount != 0 {
				t.Fatalf("request body was read before authentication: readCount=%d", body.readCount)
			}
		})
	}
}

func TestAuthzWithProjectRejectsOversizedBody(t *testing.T) {
	originalLimit := maxRequestBodyBytes
	maxRequestBodyBytes = 8
	t.Cleanup(func() {
		maxRequestBodyBytes = originalLimit
	})

	svc := gatewayService{}
	req, err := http.NewRequest(http.MethodPost, "/api/v1/finding/put-finding", bytes.NewBufferString(`{"project_id":1,"name":"too-large"}`))
	if err != nil {
		t.Fatalf("failed to create request: %v", err)
	}
	req = req.WithContext(context.WithValue(req.Context(), userKey, &requestUser{sub: "sub", userID: 1}))
	rec := httptest.NewRecorder()
	called := false

	svc.authzWithProject(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called = true
		w.WriteHeader(http.StatusOK)
	})).ServeHTTP(rec, req)

	if rec.Code != http.StatusRequestEntityTooLarge {
		t.Fatalf("unexpected status: want=%d got=%d", http.StatusRequestEntityTooLarge, rec.Code)
	}
	if called {
		t.Fatal("next handler should not be called")
	}
}

func TestAuthzOrgForToken(t *testing.T) {
	orgIAMMock := orgiammocks.NewOrgIAMServiceClient(t)
	svc := gatewayService{
		org_iamClient: orgIAMMock,
	}
	const orgID = uint32(100)
	const tokenID = uint32(200)

	cases := []struct {
		name      string
		inputUser *requestUser
		query     string
		want      bool
		mockResp  *org_iam.IsAuthorizedOrgTokenResponse
		mockErr   error
	}{
		{
			name:      "OK with organization id",
			inputUser: &requestUser{orgAccessTokenOrgID: orgID, orgAccessTokenID: tokenID},
			query:     fmt.Sprintf("organization_id=%d", orgID),
			want:      true,
			mockResp:  &org_iam.IsAuthorizedOrgTokenResponse{Ok: true},
		},
		{
			name:      "OK without organization id",
			inputUser: &requestUser{orgAccessTokenOrgID: orgID, orgAccessTokenID: tokenID},
			want:      true,
			mockResp:  &org_iam.IsAuthorizedOrgTokenResponse{Ok: true},
		},
		{
			name:      "NG missing token info",
			inputUser: &requestUser{orgAccessTokenOrgID: orgID},
			query:     fmt.Sprintf("organization_id=%d", orgID),
			want:      false,
		},
		{
			name:      "NG mismatched organization id",
			inputUser: &requestUser{orgAccessTokenOrgID: orgID, orgAccessTokenID: tokenID},
			query:     "organization_id=999",
			want:      false,
		},
		{
			name:      "NG authz error",
			inputUser: &requestUser{orgAccessTokenOrgID: orgID, orgAccessTokenID: tokenID},
			query:     fmt.Sprintf("organization_id=%d", orgID),
			want:      false,
			mockErr:   errors.New("something error"),
		},
		{
			name:      "NG unauthorized response",
			inputUser: &requestUser{orgAccessTokenOrgID: orgID, orgAccessTokenID: tokenID},
			query:     fmt.Sprintf("organization_id=%d", orgID),
			want:      false,
			mockResp:  &org_iam.IsAuthorizedOrgTokenResponse{Ok: false},
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if c.mockResp != nil || c.mockErr != nil {
				orgIAMMock.On(
					"IsAuthorizedOrgToken",
					mock.Anything,
					mock.MatchedBy(func(req *org_iam.IsAuthorizedOrgTokenRequest) bool {
						return req.OrganizationId == c.inputUser.orgAccessTokenOrgID &&
							req.AccessTokenId == c.inputUser.orgAccessTokenID &&
							req.ActionName == "organization/action"
					}),
				).Return(c.mockResp, c.mockErr).Once()
			}
			uri := "/api/v1/organization/action"
			if c.query != "" {
				uri = uri + "?" + c.query
			}
			req, _ := http.NewRequest(http.MethodGet, uri, nil)
			got := svc.authzOrgForToken(c.inputUser, req)
			if got != c.want {
				t.Fatalf("Unexpected response. want=%t, got=%t", c.want, got)
			}
		})
	}
}

func TestGetActionNameFromURI(t *testing.T) {
	cases := []struct {
		name  string
		input string
		want  string
	}{
		{
			name:  "OK",
			input: "/api/v1/service/path1/path2",
			want:  "service/path1",
		},
		{
			name:  "OK No sub paths",
			input: "/api/v1/service/",
			want:  "service/",
		},
		{
			name:  "NG blank",
			input: "",
			want:  "",
		},
		{
			name:  "NG No prefix(/)",
			input: "service-action1-action2",
			want:  "",
		},
		{
			name:  "NG No sub slashes",
			input: "/api/v1/service-action1-action2",
			want:  "",
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got := getActionNameFromURI(c.input)
			if got != c.want {
				t.Fatalf("Unexpected response. want=%s, got=%s", c.want, got)
			}
		})
	}
}

func TestGetServiceNameFromURI(t *testing.T) {
	cases := []struct {
		name  string
		input string
		want  string
	}{
		{
			name:  "OK",
			input: "/api/v1/service/path1/path2",
			want:  "service",
		},
		{
			name:  "OK No sub paths",
			input: "/api/v1/service",
			want:  "service",
		},
		{
			name:  "NG blank",
			input: "",
			want:  "",
		},
		{
			name:  "NG No prefix(/)",
			input: "service-action1-action2",
			want:  "",
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got := getServiceNameFromURI(c.input)
			if got != c.want {
				t.Fatalf("Unexpected response. want=%s, got=%s", c.want, got)
			}
		})
	}
}

func TestAuthzAdmin(t *testing.T) {
	iamMock := iammocks.NewIAMServiceClient(t)
	svc := gatewayService{
		iamClient: iamMock,
	}
	cases := []struct {
		name          string
		inputUser     *requestUser
		inputProject  string
		want          bool
		mockAuthzResp *iam.IsAuthorizedAdminResponse
		mockAuthzErr  error
	}{
		{
			name:          "OK",
			inputUser:     &requestUser{sub: "sub", userID: 1},
			inputProject:  "project_id=1",
			mockAuthzResp: &iam.IsAuthorizedAdminResponse{Ok: true},
			want:          true,
		},
		{
			name:         "NG Invalid userID",
			inputUser:    &requestUser{sub: "sub", userID: 0},
			inputProject: "project_id=1",
			want:         false,
		},
		{
			name:         "NG Authz API error",
			inputUser:    &requestUser{sub: "sub", userID: 1},
			inputProject: "project_id=1",
			want:         false,
			mockAuthzErr: errors.New("something error"),
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if c.mockAuthzResp != nil || c.mockAuthzErr != nil {
				iamMock.On("IsAuthorizedAdmin", mock.Anything, mock.Anything).Return(c.mockAuthzResp, c.mockAuthzErr).Once()
			}
			req, _ := http.NewRequest(http.MethodGet, "/api/v1/admin/api/?"+c.inputProject, nil)
			got := svc.authzAdmin(c.inputUser, req)
			if got != c.want {
				t.Fatalf("Unexpected response. want=%t, got=%t", c.want, got)
			}
			c.mockAuthzResp = nil
			c.mockAuthzErr = nil
		})
	}
}

func TestIsHumanAccess(t *testing.T) {
	cases := []struct {
		name  string
		input *requestUser
		want  bool
	}{
		{
			name:  "Human",
			input: &requestUser{sub: "sub", userID: 1},
			want:  true,
		},
		{
			name:  "Not human 1",
			input: &requestUser{sub: "sub", accessTokenID: 1},
			want:  false,
		},
		{
			name:  "Not human 2",
			input: &requestUser{sub: "sub", accessTokenID: 1, userID: 1},
			want:  false,
		},
		{
			name:  "Organization token",
			input: &requestUser{sub: "sub", orgAccessTokenID: 1, orgAccessTokenOrgID: 2},
			want:  false,
		},
		{
			name:  "Nil",
			input: nil,
			want:  false,
		},
		{
			name:  "No userID",
			input: &requestUser{sub: "sub"},
			want:  false,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got := isHumanAccess(c.input)
			if got != c.want {
				t.Fatalf("Unexpected response. want=%t, got=%t", c.want, got)
			}
		})
	}
}

func TestShouldVerifyCSRFTokenURI(t *testing.T) {
	cases := []struct {
		name  string
		input string
		want  bool
	}{
		{
			name:  "should verify 1",
			input: "/api/v1/uri",
			want:  true,
		},
		{
			name:  "ignore URI 1",
			input: "/healthz",
			want:  false,
		},
		{
			name:  "ignore URI 2",
			input: "/api/v1/signin",
			want:  false,
		},
		{
			name:  "slash suffix 1",
			input: "/api/v1/uri/",
			want:  true,
		},
		{
			name:  "slash suffix 2",
			input: "/api/v1/signin/",
			want:  false,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got := shouldVerifyCSRFTokenURI(c.input)
			if got != c.want {
				t.Fatalf("Unexpected response. want=%t, got=%t", c.want, got)
			}
		})
	}
}

func GetTestHandler() http.HandlerFunc {
	fn := func(w http.ResponseWriter, r *http.Request) {
		_, err := w.Write([]byte("OK"))
		if err != nil {
			log.Fatal(err)
		}
	}
	return http.HandlerFunc(fn)
}

func setContextValue(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		name := r.Header.Get("test-requestUser-userName")
		strUserID := r.Header.Get("test-requestUser-userID")
		userID, _ := strconv.Atoi(strUserID)
		requestUser := &requestUser{
			userID: uint32(userID),
			name:   name,
		}
		r = r.WithContext(context.WithValue(r.Context(), userKey, requestUser))
		next.ServeHTTP(w, r)
	}
	return http.HandlerFunc(fn)
}

type mockClaimsClient struct {
	claims     *jwt.MapClaims
	userName   string
	userIdpKey string
	err        error
}

func (c *mockClaimsClient) getClaims(ctx context.Context, tokenString string) (*jwt.MapClaims, error) {
	return c.claims, c.err
}

func (c *mockClaimsClient) getUserName(claims *jwt.MapClaims) string {
	return c.userName
}

func (c *mockClaimsClient) getUserIdpKey(claims *jwt.MapClaims) string {
	return c.userIdpKey
}

func newMockClaimsClient(claims *jwt.MapClaims, userName, userIdpKey string, err error) *mockClaimsClient {
	return &mockClaimsClient{
		claims:     claims,
		userName:   userName,
		userIdpKey: userIdpKey,
		err:        err,
	}
}

func TestUpdateUserFromIdp(t *testing.T) {
	iamMock := iammocks.NewIAMServiceClient(t)
	g := gatewayService{
		iamClient:      iamMock,
		uidHeader:      "uid",
		oidcDataHeader: "oidc",
	}
	handler := setContextValue(g.UpdateUserFromIdp(GetTestHandler()))

	cases := []struct {
		name            string
		userID          string
		oidcData        string
		claims          *jwt.MapClaims
		userName        string
		userIdpKey      string
		mockClaimsErr   error
		mockResponse    string
		mockErr         error
		requestUser     *requestUser
		mockPutUserResp *iam.PutUserResponse
		mockPutUserErr  error
		wantStatusCode  int
	}{
		{
			name:     "OK Update",
			userID:   "sub",
			oidcData: "",
			claims: &jwt.MapClaims{
				"username":     "username",
				"user_idp_key": "uid",
			},
			requestUser: &requestUser{
				sub:    "sub",
				userID: 1,
				name:   "name",
			},
			mockPutUserResp: &iam.PutUserResponse{},
			userName:        "username",
			userIdpKey:      "uid",
			wantStatusCode:  http.StatusOK,
		},
		{
			name:           "OK userID is not set",
			userID:         "",
			wantStatusCode: http.StatusOK,
		},
		{
			name:           "NG get claims error",
			userID:         "sub",
			mockClaimsErr:  errors.New("something error"),
			wantStatusCode: http.StatusForbidden,
		},
		{
			name:   "NG userIdpKey is empty",
			userID: "sub",
			claims: &jwt.MapClaims{
				"username":     "username",
				"user_idp_key": "",
			},
			wantStatusCode: http.StatusForbidden,
		},
		{
			name:       "OK getRequestUser error",
			userID:     "sub",
			userIdpKey: "uid",
			claims: &jwt.MapClaims{
				"username":     "username",
				"user_idp_key": "uid",
			},
			mockPutUserResp: &iam.PutUserResponse{},
			requestUser:     nil,
			wantStatusCode:  http.StatusOK,
		},
		{
			name:       "NG PutUserError",
			userID:     "sub",
			userIdpKey: "uid",
			claims: &jwt.MapClaims{
				"username":     "username",
				"user_idp_key": "uid",
			},
			requestUser: &requestUser{
				sub:    "sub",
				userID: 1,
				name:   "name",
			},
			mockPutUserErr: errors.New("something error"),
			wantStatusCode: http.StatusInternalServerError,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if c.mockPutUserResp != nil || c.mockPutUserErr != nil {
				iamMock.On("PutUser", mock.Anything, mock.Anything).Return(c.mockPutUserResp, c.mockPutUserErr).Once()
			}
			req := httptest.NewRequest(http.MethodGet, "/api/v1/signin", nil)
			req.Header.Set("uid", c.userID)
			req.Header.Set("oidc", c.oidcData)
			// テスト時にrequestにcontextを注入しても反映されないため、
			// テスト用に作成したmiddlewareでヘッダに挿入した値をcontextに注入します
			if c.requestUser != nil {
				req.Header.Set("test-requestUser-userID", fmt.Sprint(c.requestUser.userID))
				req.Header.Set("test-requestUser-userName", c.requestUser.name)
			}
			g.claimsClient = newMockClaimsClient(c.claims, c.userName, c.userIdpKey, c.mockClaimsErr)
			rec := httptest.NewRecorder()
			handler.ServeHTTP(rec, req)
			res := rec.Result()
			defer res.Body.Close()
			if c.wantStatusCode != res.StatusCode {
				t.Fatalf("Unexpected statusCode: want=%+v, got=%+v", c.wantStatusCode, res.StatusCode)
			}
		})
	}
}
