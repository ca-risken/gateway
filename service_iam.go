package main

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/ca-risken/core/proto/iam"
	"github.com/vikyd/zero"
)

func (g *gatewayService) putUserHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	user, err := getRequestUserSub(r)
	if err != nil {
		writeResponse(ctx, w, http.StatusUnauthorized, map[string]interface{}{errorJSONKey: errors.New("InvalidUser")})
	}
	req := &iam.PutUserRequest{
		User: &iam.UserForUpsert{},
	}
	req.User.Sub = user.sub // force update sub
	if err := bind(req, r); err != nil {
		appLogger.Warnf(ctx, "Failed to bind request, req=%s, err=%+v", "PutUserRequest", err)
	}
	oidcData := r.Header.Get(g.oidcDataHeader)
	claims, err := g.claimsClient.getClaims(ctx, oidcData)
	if err != nil {
		writeResponse(ctx, w, http.StatusForbidden, map[string]interface{}{errorJSONKey: errors.New("invalid token")})
		return
	}
	userIdpKey := g.claimsClient.getUserIdpKey(claims)
	if userIdpKey == "" {
		writeResponse(ctx, w, http.StatusForbidden, map[string]interface{}{errorJSONKey: errors.New("userIdpKey is not found in token")})
		return
	}
	req.User.UserIdpKey = userIdpKey
	if err := req.Validate(); err != nil {
		appLogger.Debugf(ctx, "debug: %v", err)
		writeResponse(ctx, w, http.StatusBadRequest, map[string]interface{}{errorJSONKey: err.Error()})
		return
	}
	resp, err := g.iamClient.PutUser(ctx, req)
	if err != nil {
		if handleErr := handleGRPCError(ctx, w, err); handleErr != nil {
			appLogger.Errorf(ctx, "HandleGRPCError: %+v", handleErr)
			writeResponse(ctx, w, http.StatusInternalServerError, map[string]interface{}{errorJSONKey: "InternalServerError"})
		}
		return
	}
	writeResponse(ctx, w, http.StatusOK, map[string]interface{}{successJSONKey: resp})
}

type generateAccessTokenResponse struct {
	AccessTokenID uint32 `json:"access_token_id"`
	AccessToken   string `json:"access_token"`
}

func (g *gatewayService) generateAccessTokenHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	req := &iam.PutAccessTokenRequest{AccessToken: &iam.AccessTokenForUpsert{}}
	if err := bind(req, r); err != nil {
		appLogger.Warnf(ctx, "Failed to bind request, req=%s, err=%+v", "PutAccessTokenRequest", err)
	}
	u, err := getRequestUser(r)
	if err != nil || zero.IsZeroVal(u.userID) {
		writeResponse(ctx, w, http.StatusBadRequest, map[string]interface{}{errorJSONKey: fmt.Errorf("Failed to get user info, userInfo=%+v, err=%+v", u, err)})
		return
	}
	req.AccessToken.LastUpdatedUserId = u.userID           // Force update
	req.AccessToken.AccessTokenId = 0                      // Force update
	req.AccessToken.PlainTextToken = generateAccessToken() // Force update

	if err := req.Validate(); err != nil {
		writeResponse(ctx, w, http.StatusBadRequest, map[string]interface{}{errorJSONKey: err.Error()})
		return
	}
	list, err := g.iamClient.ListAccessToken(ctx, &iam.ListAccessTokenRequest{
		ProjectId: req.ProjectId,
		Name:      req.AccessToken.Name,
	})
	if err != nil {
		if handleErr := handleGRPCError(ctx, w, err); handleErr != nil {
			appLogger.Errorf(ctx, "HandleGRPCError: %+v", handleErr)
			writeResponse(ctx, w, http.StatusInternalServerError, map[string]interface{}{errorJSONKey: "InternalServerError"})
		}
		return
	}
	if len(list.AccessToken) > 0 {
		writeResponse(ctx, w, http.StatusBadRequest, map[string]interface{}{errorJSONKey: fmt.Errorf("Token already exists, token_name=%s", req.AccessToken.Name)})
		return
	}

	// Call API
	resp, err := g.iamClient.PutAccessToken(ctx, req)
	if err != nil {
		if handleErr := handleGRPCError(ctx, w, err); handleErr != nil {
			appLogger.Errorf(ctx, "HandleGRPCError: %+v", handleErr)
			writeResponse(ctx, w, http.StatusInternalServerError, map[string]interface{}{errorJSONKey: "InternalServerError"})
		}
		return
	}
	writeResponse(ctx, w, http.StatusOK, map[string]interface{}{successJSONKey: &generateAccessTokenResponse{
		AccessTokenID: resp.AccessToken.AccessTokenId,
		AccessToken:   encodeAccessToken(resp.AccessToken.ProjectId, resp.AccessToken.AccessTokenId, req.AccessToken.PlainTextToken),
	}})
}

func (g *gatewayService) updateAccessTokenHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	req := &iam.PutAccessTokenRequest{AccessToken: &iam.AccessTokenForUpsert{}}
	if err := bind(req, r); err != nil {
		appLogger.Warnf(ctx, "Failed to bind request, req=%s, err=%+v", "PutAccessTokenRequest", err)
	}
	u, err := getRequestUser(r)
	if err != nil || zero.IsZeroVal(u.userID) {
		writeResponse(ctx, w, http.StatusBadRequest, map[string]interface{}{errorJSONKey: fmt.Errorf("Failed to get user info, userInfo=%+v, err=%+v", u, err)})
		return
	}
	req.AccessToken.LastUpdatedUserId = u.userID // Force update data

	if err := req.Validate(); err != nil {
		appLogger.Error(ctx, err)
		writeResponse(ctx, w, http.StatusBadRequest, map[string]interface{}{errorJSONKey: err.Error()})
		return
	}
	if zero.IsZeroVal(req.AccessToken.AccessTokenId) {
		writeResponse(ctx, w, http.StatusBadRequest, map[string]interface{}{errorJSONKey: errors.New("Required access_token_id")})
		return
	}

	// Call API
	resp, err := g.iamClient.PutAccessToken(ctx, req)
	if err != nil {
		if handleErr := handleGRPCError(ctx, w, err); handleErr != nil {
			appLogger.Errorf(ctx, "HandleGRPCError: %+v", handleErr)
			writeResponse(ctx, w, http.StatusInternalServerError, map[string]interface{}{errorJSONKey: "InternalServerError"})
		}
		return
	}
	writeResponse(ctx, w, http.StatusOK, map[string]interface{}{successJSONKey: resp})
}
