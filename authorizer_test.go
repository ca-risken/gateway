package main

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/ca-risken/core/proto/iam"
	iammocks "github.com/ca-risken/core/proto/iam/mocks"
	"github.com/ca-risken/core/proto/organization_iam"
	organizationiammocks "github.com/ca-risken/core/proto/organization_iam/mocks"
	"github.com/golang-jwt/jwt/v4"
	"github.com/stretchr/testify/mock"
)

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
		setupMocks    func(iamMock *iammocks.IAMServiceClient, orgMock *organizationiammocks.OrganizationIAMServiceClient)
	}{
		{
			name:          "Project token",
			authorization: projectToken,
			expectUser:    true,
			setupMocks: func(iamMock *iammocks.IAMServiceClient, _ *organizationiammocks.OrganizationIAMServiceClient) {
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
			setupMocks: func(_ *iammocks.IAMServiceClient, orgMock *organizationiammocks.OrganizationIAMServiceClient) {
				orgMock.On("AuthenticateOrganizationAccessToken", mock.Anything, mock.MatchedBy(func(req *organization_iam.AuthenticateOrganizationAccessTokenRequest) bool {
					return req.OrganizationId == orgID &&
						req.AccessTokenId == accessTokenID &&
						req.PlainTextToken == plainToken
				})).Return(&organization_iam.AuthenticateOrganizationAccessTokenResponse{
					AccessToken: &organization_iam.OrganizationAccessToken{
						AccessTokenId: accessTokenID,
					},
				}, nil).Once()
			},
			assert: func(t *testing.T, u *requestUser) {
				if u.orgAccessTokenID != accessTokenID || u.orgID != orgID {
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
			orgIAMMock := organizationiammocks.NewOrganizationIAMServiceClient(t)
			svc := gatewayService{
				iamClient:              iamMock,
				organization_iamClient: orgIAMMock,
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
	svc := gatewayService{
		iamClient: iamMock,
	}
	cases := []struct {
		name         string
		inputUser    *requestUser
		inputProject string
		want         bool
		mockResp     *iam.IsAuthorizedTokenResponse
		mockErr      error
	}{
		{
			name:         "OK1",
			inputUser:    &requestUser{accessTokenID: 123, accessTokenProjectID: 1},
			inputProject: "project_id=1",
			mockResp:     &iam.IsAuthorizedTokenResponse{Ok: true},
			want:         true,
		},
		{
			name:      "OK2",
			inputUser: &requestUser{accessTokenID: 123, accessTokenProjectID: 1},
			mockResp:  &iam.IsAuthorizedTokenResponse{Ok: true},
			want:      true,
		},
		{
			name:         "NG No token",
			inputUser:    &requestUser{sub: "sub", accessTokenProjectID: 1},
			inputProject: "project_id=1",
			want:         false,
		},
		{
			name:         "NG Invalid project",
			inputUser:    &requestUser{accessTokenID: 123, accessTokenProjectID: 0},
			inputProject: "project_id=0",
			want:         false,
		},
		{
			name:         "NG Invalid project(unmatch)",
			inputUser:    &requestUser{accessTokenID: 123, accessTokenProjectID: 1},
			inputProject: "project_id=2",
			want:         false,
		},
		{
			name:         "NG IAM error",
			inputUser:    &requestUser{accessTokenID: 123, accessTokenProjectID: 1},
			inputProject: "project_id=1",
			want:         false,
			mockErr:      errors.New("something error"),
		}}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if c.mockResp != nil || c.mockErr != nil {
				iamMock.On("IsAuthorizedToken", mock.Anything, mock.Anything).Return(c.mockResp, c.mockErr).Once()
			}
			req, _ := http.NewRequest(http.MethodGet, "/api/v1/service/action/?"+c.inputProject, nil)
			got := svc.authzProjectForToken(c.inputUser, req)
			if got != c.want {
				t.Fatalf("Unexpected response. want=%t, got=%t", c.want, got)
			}
		})
	}
}

func TestAuthzOrgForToken(t *testing.T) {
	orgIAMMock := organizationiammocks.NewOrganizationIAMServiceClient(t)
	svc := gatewayService{
		organization_iamClient: orgIAMMock,
	}
	const orgID = uint32(100)
	const tokenID = uint32(200)

	cases := []struct {
		name      string
		inputUser *requestUser
		query     string
		want      bool
		mockResp  *organization_iam.IsAuthorizedOrganizationTokenResponse
		mockErr   error
	}{
		{
			name:      "OK with organization id",
			inputUser: &requestUser{orgID: orgID, orgAccessTokenID: tokenID},
			query:     fmt.Sprintf("organization_id=%d", orgID),
			want:      true,
			mockResp:  &organization_iam.IsAuthorizedOrganizationTokenResponse{Ok: true},
		},
		{
			name:      "OK without organization id",
			inputUser: &requestUser{orgID: orgID, orgAccessTokenID: tokenID},
			want:      true,
			mockResp:  &organization_iam.IsAuthorizedOrganizationTokenResponse{Ok: true},
		},
		{
			name:      "NG missing token info",
			inputUser: &requestUser{orgID: orgID},
			query:     fmt.Sprintf("organization_id=%d", orgID),
			want:      false,
		},
		{
			name:      "NG mismatched organization id",
			inputUser: &requestUser{orgID: orgID, orgAccessTokenID: tokenID},
			query:     "organization_id=999",
			want:      false,
		},
		{
			name:      "NG authz error",
			inputUser: &requestUser{orgID: orgID, orgAccessTokenID: tokenID},
			query:     fmt.Sprintf("organization_id=%d", orgID),
			want:      false,
			mockErr:   errors.New("something error"),
		},
		{
			name:      "NG unauthorized response",
			inputUser: &requestUser{orgID: orgID, orgAccessTokenID: tokenID},
			query:     fmt.Sprintf("organization_id=%d", orgID),
			want:      false,
			mockResp:  &organization_iam.IsAuthorizedOrganizationTokenResponse{Ok: false},
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if c.mockResp != nil || c.mockErr != nil {
				orgIAMMock.On(
					"IsAuthorizedOrganizationToken",
					mock.Anything,
					mock.MatchedBy(func(req *organization_iam.IsAuthorizedOrganizationTokenRequest) bool {
						return req.OrganizationId == c.inputUser.orgID &&
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
			input: &requestUser{sub: "sub", orgAccessTokenID: 1, orgID: 2},
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
	ts := httptest.NewServer(setContextValue(g.UpdateUserFromIdp(GetTestHandler())))
	defer ts.Close()

	var u bytes.Buffer
	u.WriteString(string(ts.URL))
	u.WriteString("/api/v1/signin")

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
		wantErr         bool
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
			client := http.Client{}
			req, _ := http.NewRequest(http.MethodGet, u.String(), nil)
			req.Header.Set("uid", c.userID)
			req.Header.Set("oidc", c.oidcData)
			// テスト時にrequestにcontextを注入しても反映されないため、
			// テスト用に作成したmiddlewareでヘッダに挿入した値をcontextに注入します
			if c.requestUser != nil {
				req.Header.Set("test-requestUser-userID", fmt.Sprint(c.requestUser.userID))
				req.Header.Set("test-requestUser-userName", c.requestUser.name)
			}
			g.claimsClient = newMockClaimsClient(c.claims, c.userName, c.userIdpKey, c.mockClaimsErr)
			res, err := client.Do(req)
			if (c.wantErr && err == nil) || (!c.wantErr && err != nil) {
				t.Fatalf("Unexpected error: wantErr=%t, err=%+v", c.wantErr, err)
			}
			if c.wantStatusCode != res.StatusCode {
				t.Fatalf("Unexpected statusCode: want=%+v, got=%+v", c.wantStatusCode, res.StatusCode)
			}
		})
	}
}
