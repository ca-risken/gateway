package main

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/CyberAgent/mimosa-core/proto/iam"
	"github.com/vikyd/zero"
)

func (g *gatewayService) listUserHandler(w http.ResponseWriter, r *http.Request) {
	req := &iam.ListUserRequest{}
	bind(req, r)
	if err := req.Validate(); err != nil {
		writeResponse(w, http.StatusBadRequest, map[string]interface{}{errorJSONKey: err.Error()})
		return
	}
	resp, err := g.iamClient.ListUser(r.Context(), req)
	if err != nil {
		writeResponse(w, http.StatusInternalServerError, map[string]interface{}{errorJSONKey: err.Error()})
		return
	}
	writeResponse(w, http.StatusOK, map[string]interface{}{successJSONKey: resp})
}

func (g *gatewayService) getUserHandler(w http.ResponseWriter, r *http.Request) {
	req := &iam.GetUserRequest{}
	bind(req, r)
	if err := req.Validate(); err != nil {
		writeResponse(w, http.StatusBadRequest, map[string]interface{}{errorJSONKey: err.Error()})
		return
	}
	resp, err := g.iamClient.GetUser(r.Context(), req)
	if err != nil {
		writeResponse(w, http.StatusInternalServerError, map[string]interface{}{errorJSONKey: err.Error()})
		return
	}
	writeResponse(w, http.StatusOK, map[string]interface{}{successJSONKey: resp})
}

func (g *gatewayService) isAdminHandler(w http.ResponseWriter, r *http.Request) {
	req := &iam.IsAdminRequest{}
	bind(req, r)
	if err := req.Validate(); err != nil {
		writeResponse(w, http.StatusBadRequest, map[string]interface{}{errorJSONKey: err.Error()})
		return
	}
	resp, err := g.iamClient.IsAdmin(r.Context(), req)
	if err != nil {
		writeResponse(w, http.StatusInternalServerError, map[string]interface{}{errorJSONKey: err.Error()})
		return
	}
	writeResponse(w, http.StatusOK, map[string]interface{}{successJSONKey: resp})
}

func (g *gatewayService) putUserHandler(w http.ResponseWriter, r *http.Request) {
	user, err := getRequestUser(r)
	if err != nil {
		writeResponse(w, http.StatusUnauthorized, map[string]interface{}{errorJSONKey: errors.New("InvalidUser")})
	}
	req := &iam.PutUserRequest{
		User: &iam.UserForUpsert{},
	}
	req.User.Sub = user.sub // force update sub
	bind(req, r)
	if err := req.Validate(); err != nil {
		appLogger.Debugf("debug: %v", err)
		writeResponse(w, http.StatusBadRequest, map[string]interface{}{errorJSONKey: err.Error()})
		return
	}
	resp, err := g.iamClient.PutUser(r.Context(), req)
	if err != nil {
		writeResponse(w, http.StatusInternalServerError, map[string]interface{}{errorJSONKey: err.Error()})
		return
	}
	writeResponse(w, http.StatusOK, map[string]interface{}{successJSONKey: resp})
}

func (g *gatewayService) listRoleHandler(w http.ResponseWriter, r *http.Request) {
	req := &iam.ListRoleRequest{}
	bind(req, r)
	if err := req.Validate(); err != nil {
		writeResponse(w, http.StatusBadRequest, map[string]interface{}{errorJSONKey: err.Error()})
		return
	}
	resp, err := g.iamClient.ListRole(r.Context(), req)
	if err != nil {
		writeResponse(w, http.StatusInternalServerError, map[string]interface{}{errorJSONKey: err.Error()})
		return
	}
	writeResponse(w, http.StatusOK, map[string]interface{}{successJSONKey: resp})
}

func (g *gatewayService) getRoleHandler(w http.ResponseWriter, r *http.Request) {
	req := &iam.GetRoleRequest{}
	bind(req, r)
	if err := req.Validate(); err != nil {
		writeResponse(w, http.StatusBadRequest, map[string]interface{}{errorJSONKey: err.Error()})
		return
	}
	resp, err := g.iamClient.GetRole(r.Context(), req)
	if err != nil {
		writeResponse(w, http.StatusInternalServerError, map[string]interface{}{errorJSONKey: err.Error()})
		return
	}
	writeResponse(w, http.StatusOK, map[string]interface{}{successJSONKey: resp})
}

func (g *gatewayService) putRoleHandler(w http.ResponseWriter, r *http.Request) {
	req := &iam.PutRoleRequest{}
	bind(req, r)
	if err := req.Validate(); err != nil {
		writeResponse(w, http.StatusBadRequest, map[string]interface{}{errorJSONKey: err.Error()})
		return
	}
	resp, err := g.iamClient.PutRole(r.Context(), req)
	if err != nil {
		writeResponse(w, http.StatusInternalServerError, map[string]interface{}{errorJSONKey: err.Error()})
		return
	}
	writeResponse(w, http.StatusOK, map[string]interface{}{successJSONKey: resp})
}

func (g *gatewayService) deleteRoleHandler(w http.ResponseWriter, r *http.Request) {
	req := &iam.DeleteRoleRequest{}
	bind(req, r)
	if err := req.Validate(); err != nil {
		writeResponse(w, http.StatusBadRequest, map[string]interface{}{errorJSONKey: err.Error()})
		return
	}
	resp, err := g.iamClient.DeleteRole(r.Context(), req)
	if err != nil {
		writeResponse(w, http.StatusInternalServerError, map[string]interface{}{errorJSONKey: err.Error()})
		return
	}
	writeResponse(w, http.StatusOK, map[string]interface{}{successJSONKey: resp})
}

func (g *gatewayService) attachRoleHandler(w http.ResponseWriter, r *http.Request) {
	req := &iam.AttachRoleRequest{}
	bind(req, r)
	if err := req.Validate(); err != nil {
		writeResponse(w, http.StatusBadRequest, map[string]interface{}{errorJSONKey: err.Error()})
		return
	}
	resp, err := g.iamClient.AttachRole(r.Context(), req)
	if err != nil {
		writeResponse(w, http.StatusInternalServerError, map[string]interface{}{errorJSONKey: err.Error()})
		return
	}
	writeResponse(w, http.StatusOK, map[string]interface{}{successJSONKey: resp})
}

func (g *gatewayService) detachRoleHandler(w http.ResponseWriter, r *http.Request) {
	req := &iam.DetachRoleRequest{}
	bind(req, r)
	if err := req.Validate(); err != nil {
		writeResponse(w, http.StatusBadRequest, map[string]interface{}{errorJSONKey: err.Error()})
		return
	}
	resp, err := g.iamClient.DetachRole(r.Context(), req)
	if err != nil {
		writeResponse(w, http.StatusInternalServerError, map[string]interface{}{errorJSONKey: err.Error()})
		return
	}
	writeResponse(w, http.StatusOK, map[string]interface{}{successJSONKey: resp})
}

func (g *gatewayService) listPolicyHandler(w http.ResponseWriter, r *http.Request) {
	req := &iam.ListPolicyRequest{}
	bind(req, r)
	if err := req.Validate(); err != nil {
		writeResponse(w, http.StatusBadRequest, map[string]interface{}{errorJSONKey: err.Error()})
		return
	}
	resp, err := g.iamClient.ListPolicy(r.Context(), req)
	if err != nil {
		writeResponse(w, http.StatusInternalServerError, map[string]interface{}{errorJSONKey: err.Error()})
		return
	}
	writeResponse(w, http.StatusOK, map[string]interface{}{successJSONKey: resp})
}

func (g *gatewayService) getPolicyHandler(w http.ResponseWriter, r *http.Request) {
	req := &iam.GetPolicyRequest{}
	bind(req, r)
	if err := req.Validate(); err != nil {
		writeResponse(w, http.StatusBadRequest, map[string]interface{}{errorJSONKey: err.Error()})
		return
	}
	resp, err := g.iamClient.GetPolicy(r.Context(), req)
	if err != nil {
		writeResponse(w, http.StatusInternalServerError, map[string]interface{}{errorJSONKey: err.Error()})
		return
	}
	writeResponse(w, http.StatusOK, map[string]interface{}{successJSONKey: resp})
}

func (g *gatewayService) putPolicyHandler(w http.ResponseWriter, r *http.Request) {
	req := &iam.PutPolicyRequest{}
	bind(req, r)
	if err := req.Validate(); err != nil {
		writeResponse(w, http.StatusBadRequest, map[string]interface{}{errorJSONKey: err.Error()})
		return
	}
	resp, err := g.iamClient.PutPolicy(r.Context(), req)
	if err != nil {
		writeResponse(w, http.StatusInternalServerError, map[string]interface{}{errorJSONKey: err.Error()})
		return
	}
	writeResponse(w, http.StatusOK, map[string]interface{}{successJSONKey: resp})
}

func (g *gatewayService) deletePolicyHandler(w http.ResponseWriter, r *http.Request) {
	req := &iam.DeletePolicyRequest{}
	bind(req, r)
	if err := req.Validate(); err != nil {
		writeResponse(w, http.StatusBadRequest, map[string]interface{}{errorJSONKey: err.Error()})
		return
	}
	resp, err := g.iamClient.DeletePolicy(r.Context(), req)
	if err != nil {
		writeResponse(w, http.StatusInternalServerError, map[string]interface{}{errorJSONKey: err.Error()})
		return
	}
	writeResponse(w, http.StatusOK, map[string]interface{}{successJSONKey: resp})
}

func (g *gatewayService) attachPolicyHandler(w http.ResponseWriter, r *http.Request) {
	req := &iam.AttachPolicyRequest{}
	bind(req, r)
	if err := req.Validate(); err != nil {
		writeResponse(w, http.StatusBadRequest, map[string]interface{}{errorJSONKey: err.Error()})
		return
	}
	resp, err := g.iamClient.AttachPolicy(r.Context(), req)
	if err != nil {
		writeResponse(w, http.StatusInternalServerError, map[string]interface{}{errorJSONKey: err.Error()})
		return
	}
	writeResponse(w, http.StatusOK, map[string]interface{}{successJSONKey: resp})
}

func (g *gatewayService) detachPolicyHandler(w http.ResponseWriter, r *http.Request) {
	req := &iam.DetachPolicyRequest{}
	bind(req, r)
	if err := req.Validate(); err != nil {
		writeResponse(w, http.StatusBadRequest, map[string]interface{}{errorJSONKey: err.Error()})
		return
	}
	resp, err := g.iamClient.DetachPolicy(r.Context(), req)
	if err != nil {
		writeResponse(w, http.StatusInternalServerError, map[string]interface{}{errorJSONKey: err.Error()})
		return
	}
	writeResponse(w, http.StatusOK, map[string]interface{}{successJSONKey: resp})
}

func (g *gatewayService) listAccessTokenHandler(w http.ResponseWriter, r *http.Request) {
	req := &iam.ListAccessTokenRequest{}
	bind(req, r)
	if err := req.Validate(); err != nil {
		writeResponse(w, http.StatusBadRequest, map[string]interface{}{errorJSONKey: err.Error()})
		return
	}
	resp, err := g.iamClient.ListAccessToken(r.Context(), req)
	if err != nil {
		writeResponse(w, http.StatusInternalServerError, map[string]interface{}{errorJSONKey: err.Error()})
		return
	}
	writeResponse(w, http.StatusOK, map[string]interface{}{successJSONKey: resp})
}

type generateAccessTokenResponse struct {
	AccessToken string `json:"access_token"`
}

func (g *gatewayService) generateAccessTokenHandler(w http.ResponseWriter, r *http.Request) {
	req := &iam.PutAccessTokenRequest{AccessToken: &iam.AccessTokenForUpsert{}}
	bind(req, r)
	u, err := getRequestUser(r)
	if err != nil || zero.IsZeroVal(u.userID) {
		writeResponse(w, http.StatusBadRequest, map[string]interface{}{errorJSONKey: fmt.Errorf("Failed to get user info, userInfo=%+v, err=%+v", u, err)})
		return
	}
	req.AccessToken.LastUpdatedUserId = u.userID           // Force update
	req.AccessToken.AccessTokenId = 0                      // Force update
	req.AccessToken.PlainTextToken = generateAccessToken() // Force update

	if err := req.Validate(); err != nil {
		writeResponse(w, http.StatusBadRequest, map[string]interface{}{errorJSONKey: err.Error()})
		return
	}
	list, err := g.iamClient.ListAccessToken(r.Context(), &iam.ListAccessTokenRequest{
		ProjectId: req.ProjectId,
		Name:      req.AccessToken.Name,
	})
	if err != nil {
		writeResponse(w, http.StatusInternalServerError, map[string]interface{}{errorJSONKey: err.Error()})
		return
	}
	if len(list.AccessToken) > 0 {
		writeResponse(w, http.StatusBadRequest, map[string]interface{}{errorJSONKey: fmt.Errorf("Token already exists, token_name=%s", req.AccessToken.Name)})
		return
	}

	// Call API
	resp, err := g.iamClient.PutAccessToken(r.Context(), req)
	if err != nil {
		writeResponse(w, http.StatusInternalServerError, map[string]interface{}{errorJSONKey: err.Error()})
		return
	}
	writeResponse(w, http.StatusOK, map[string]interface{}{successJSONKey: &generateAccessTokenResponse{
		AccessToken: encodeAccessToken(resp.AccessToken.ProjectId, resp.AccessToken.AccessTokenId, req.AccessToken.PlainTextToken),
	}})
}

func (g *gatewayService) updateAccessTokenHandler(w http.ResponseWriter, r *http.Request) {
	req := &iam.PutAccessTokenRequest{AccessToken: &iam.AccessTokenForUpsert{}}
	bind(req, r)
	u, err := getRequestUser(r)
	if err != nil || zero.IsZeroVal(u.userID) {
		writeResponse(w, http.StatusBadRequest, map[string]interface{}{errorJSONKey: fmt.Errorf("Failed to get user info, userInfo=%+v, err=%+v", u, err)})
		return
	}
	req.AccessToken.LastUpdatedUserId = u.userID // Force update data

	if err := req.Validate(); err != nil {
		appLogger.Error(err)
		writeResponse(w, http.StatusBadRequest, map[string]interface{}{errorJSONKey: err.Error()})
		return
	}
	if zero.IsZeroVal(req.AccessToken.AccessTokenId) {
		writeResponse(w, http.StatusBadRequest, map[string]interface{}{errorJSONKey: errors.New("Required access_token_id")})
		return
	}

	// Call API
	resp, err := g.iamClient.PutAccessToken(r.Context(), req)
	if err != nil {
		writeResponse(w, http.StatusInternalServerError, map[string]interface{}{errorJSONKey: err.Error()})
		return
	}
	writeResponse(w, http.StatusOK, map[string]interface{}{successJSONKey: resp})
}

func (g *gatewayService) deleteAccessTokenHandler(w http.ResponseWriter, r *http.Request) {
	req := &iam.DeleteAccessTokenRequest{}
	bind(req, r)
	if err := req.Validate(); err != nil {
		writeResponse(w, http.StatusBadRequest, map[string]interface{}{errorJSONKey: err.Error()})
		return
	}
	resp, err := g.iamClient.DeleteAccessToken(r.Context(), req)
	if err != nil {
		writeResponse(w, http.StatusInternalServerError, map[string]interface{}{errorJSONKey: err.Error()})
		return
	}
	writeResponse(w, http.StatusOK, map[string]interface{}{successJSONKey: resp})
}

func (g *gatewayService) attachAccessTokenRoleHandler(w http.ResponseWriter, r *http.Request) {
	req := &iam.AttachAccessTokenRoleRequest{}
	bind(req, r)
	if err := req.Validate(); err != nil {
		appLogger.Info(err)
		writeResponse(w, http.StatusBadRequest, map[string]interface{}{errorJSONKey: err.Error()})
		return
	}
	resp, err := g.iamClient.AttachAccessTokenRole(r.Context(), req)
	if err != nil {
		writeResponse(w, http.StatusInternalServerError, map[string]interface{}{errorJSONKey: err.Error()})
		return
	}
	writeResponse(w, http.StatusOK, map[string]interface{}{successJSONKey: resp})
}

func (g *gatewayService) detachAccessTokenRoleHandler(w http.ResponseWriter, r *http.Request) {
	req := &iam.DetachAccessTokenRoleRequest{}
	bind(req, r)
	if err := req.Validate(); err != nil {
		writeResponse(w, http.StatusBadRequest, map[string]interface{}{errorJSONKey: err.Error()})
		return
	}
	resp, err := g.iamClient.DetachAccessTokenRole(r.Context(), req)
	if err != nil {
		writeResponse(w, http.StatusInternalServerError, map[string]interface{}{errorJSONKey: err.Error()})
		return
	}
	writeResponse(w, http.StatusOK, map[string]interface{}{successJSONKey: resp})
}
