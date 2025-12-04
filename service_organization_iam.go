package main

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/ca-risken/core/proto/organization_iam"
)

func (g *gatewayService) generateOrganizationAccessTokenOrganization_iamHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	req := &organization_iam.PutOrganizationAccessTokenRequest{}
	if err := bind(req, r); err != nil {
		appLogger.Warnf(ctx, "Failed to bind request, req=%s, err=%+v", "PutOrganizationAccessTokenRequest", err)
	}
	u, err := getRequestUser(r)
	if err != nil || u.userID == 0 {
		writeResponse(ctx, w, http.StatusBadRequest, map[string]interface{}{errorJSONKey: fmt.Errorf("Failed to get user info, userInfo=%+v, err=%+v", u, err)})
		return
	}
	req.LastUpdatedUserId = u.userID
	req.AccessTokenId = 0
	req.PlainTextToken = generateAccessToken()
	if err := req.Validate(); err != nil {
		writeResponse(ctx, w, http.StatusBadRequest, map[string]interface{}{errorJSONKey: err.Error()})
		return
	}

	list, err := g.organization_iamClient.ListOrganizationAccessToken(ctx, &organization_iam.ListOrganizationAccessTokenRequest{
		OrganizationId: req.OrganizationId,
		Name:           req.Name,
	})
	if err != nil {
		if handleErr := handleGRPCError(ctx, w, err); handleErr != nil {
			appLogger.Errorf(ctx, "HandleGRPCError: %+v", handleErr)
			writeResponse(ctx, w, http.StatusInternalServerError, map[string]interface{}{errorJSONKey: "InternalServerError"})
		}
		return
	}
	if list != nil && len(list.AccessToken) > 0 {
		writeResponse(ctx, w, http.StatusBadRequest, map[string]interface{}{errorJSONKey: fmt.Errorf("Token already exists, token_name=%s", req.Name)})
		return
	}

	resp, err := g.organization_iamClient.PutOrganizationAccessToken(ctx, req)
	if err != nil {
		if handleErr := handleGRPCError(ctx, w, err); handleErr != nil {
			appLogger.Errorf(ctx, "HandleGRPCError: %+v", handleErr)
			writeResponse(ctx, w, http.StatusInternalServerError, map[string]interface{}{errorJSONKey: "InternalServerError"})
		}
		return
	}
	if resp == nil || resp.AccessToken == nil {
		writeResponse(ctx, w, http.StatusInternalServerError, map[string]interface{}{errorJSONKey: "Invalid response from OrganizationIAMService"})
		return
	}
	writeResponse(ctx, w, http.StatusOK, map[string]interface{}{successJSONKey: &generateAccessTokenResponse{
		AccessTokenID: resp.AccessToken.AccessTokenId,
		AccessToken:   encodeOrgAccessToken(resp.AccessToken.OrganizationId, resp.AccessToken.AccessTokenId, req.PlainTextToken),
	}})
}

func (g *gatewayService) updateOrganizationAccessTokenOrganization_iamHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	req := &organization_iam.PutOrganizationAccessTokenRequest{}
	if err := bind(req, r); err != nil {
		appLogger.Warnf(ctx, "Failed to bind request, req=%s, err=%+v", "PutOrganizationAccessTokenRequest", err)
	}
	u, err := getRequestUser(r)
	if err != nil || u.userID == 0 {
		writeResponse(ctx, w, http.StatusBadRequest, map[string]interface{}{errorJSONKey: fmt.Errorf("Failed to get user info, userInfo=%+v, err=%+v", u, err)})
		return
	}
	req.LastUpdatedUserId = u.userID
	if err := req.Validate(); err != nil {
		writeResponse(ctx, w, http.StatusBadRequest, map[string]interface{}{errorJSONKey: err.Error()})
		return
	}
	if req.AccessTokenId == 0 {
		writeResponse(ctx, w, http.StatusBadRequest, map[string]interface{}{errorJSONKey: errors.New("Required access_token_id")})
		return
	}

	resp, err := g.organization_iamClient.PutOrganizationAccessToken(ctx, req)
	if err != nil {
		if handleErr := handleGRPCError(ctx, w, err); handleErr != nil {
			appLogger.Errorf(ctx, "HandleGRPCError: %+v", handleErr)
			writeResponse(ctx, w, http.StatusInternalServerError, map[string]interface{}{errorJSONKey: "InternalServerError"})
		}
		return
	}
	writeResponse(ctx, w, http.StatusOK, map[string]interface{}{successJSONKey: resp})
}
