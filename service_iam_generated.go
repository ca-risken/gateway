// Code generated by protoc-gen-service. DO NOT EDIT.
// source: iam/service.proto

package main

import (
	"net/http"

	"github.com/ca-risken/core/proto/iam"
)

func (g *gatewayService) listUserIamHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	req := &iam.ListUserRequest{}
	if err := bind(req, r); err != nil {
		writeResponse(ctx, w, http.StatusBadRequest, map[string]interface{}{errorJSONKey: err.Error()})
		return
	}
	if err := req.Validate(); err != nil {
		writeResponse(ctx, w, http.StatusBadRequest, map[string]interface{}{errorJSONKey: err.Error()})
		return
	}
	resp, err := g.iamClient.ListUser(ctx, req)
	if err != nil {
		writeResponse(ctx, w, http.StatusInternalServerError, map[string]interface{}{errorJSONKey: err.Error()})
		return
	}
	writeResponse(ctx, w, http.StatusOK, map[string]interface{}{successJSONKey: resp})
}

func (g *gatewayService) getUserIamHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	req := &iam.GetUserRequest{}
	if err := bind(req, r); err != nil {
		writeResponse(ctx, w, http.StatusBadRequest, map[string]interface{}{errorJSONKey: err.Error()})
		return
	}
	if err := req.Validate(); err != nil {
		writeResponse(ctx, w, http.StatusBadRequest, map[string]interface{}{errorJSONKey: err.Error()})
		return
	}
	resp, err := g.iamClient.GetUser(ctx, req)
	if err != nil {
		writeResponse(ctx, w, http.StatusInternalServerError, map[string]interface{}{errorJSONKey: err.Error()})
		return
	}
	writeResponse(ctx, w, http.StatusOK, map[string]interface{}{successJSONKey: resp})
}

func (g *gatewayService) listRoleIamHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	req := &iam.ListRoleRequest{}
	if err := bind(req, r); err != nil {
		writeResponse(ctx, w, http.StatusBadRequest, map[string]interface{}{errorJSONKey: err.Error()})
		return
	}
	if err := req.Validate(); err != nil {
		writeResponse(ctx, w, http.StatusBadRequest, map[string]interface{}{errorJSONKey: err.Error()})
		return
	}
	resp, err := g.iamClient.ListRole(ctx, req)
	if err != nil {
		writeResponse(ctx, w, http.StatusInternalServerError, map[string]interface{}{errorJSONKey: err.Error()})
		return
	}
	writeResponse(ctx, w, http.StatusOK, map[string]interface{}{successJSONKey: resp})
}

func (g *gatewayService) getRoleIamHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	req := &iam.GetRoleRequest{}
	if err := bind(req, r); err != nil {
		writeResponse(ctx, w, http.StatusBadRequest, map[string]interface{}{errorJSONKey: err.Error()})
		return
	}
	if err := req.Validate(); err != nil {
		writeResponse(ctx, w, http.StatusBadRequest, map[string]interface{}{errorJSONKey: err.Error()})
		return
	}
	resp, err := g.iamClient.GetRole(ctx, req)
	if err != nil {
		writeResponse(ctx, w, http.StatusInternalServerError, map[string]interface{}{errorJSONKey: err.Error()})
		return
	}
	writeResponse(ctx, w, http.StatusOK, map[string]interface{}{successJSONKey: resp})
}

func (g *gatewayService) putRoleIamHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	req := &iam.PutRoleRequest{}
	if err := bind(req, r); err != nil {
		writeResponse(ctx, w, http.StatusBadRequest, map[string]interface{}{errorJSONKey: err.Error()})
		return
	}
	if err := req.Validate(); err != nil {
		writeResponse(ctx, w, http.StatusBadRequest, map[string]interface{}{errorJSONKey: err.Error()})
		return
	}
	resp, err := g.iamClient.PutRole(ctx, req)
	if err != nil {
		writeResponse(ctx, w, http.StatusInternalServerError, map[string]interface{}{errorJSONKey: err.Error()})
		return
	}
	writeResponse(ctx, w, http.StatusOK, map[string]interface{}{successJSONKey: resp})
}

func (g *gatewayService) deleteRoleIamHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	req := &iam.DeleteRoleRequest{}
	if err := bind(req, r); err != nil {
		writeResponse(ctx, w, http.StatusBadRequest, map[string]interface{}{errorJSONKey: err.Error()})
		return
	}
	if err := req.Validate(); err != nil {
		writeResponse(ctx, w, http.StatusBadRequest, map[string]interface{}{errorJSONKey: err.Error()})
		return
	}
	resp, err := g.iamClient.DeleteRole(ctx, req)
	if err != nil {
		writeResponse(ctx, w, http.StatusInternalServerError, map[string]interface{}{errorJSONKey: err.Error()})
		return
	}
	writeResponse(ctx, w, http.StatusOK, map[string]interface{}{successJSONKey: resp})
}

func (g *gatewayService) attachRoleIamHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	req := &iam.AttachRoleRequest{}
	if err := bind(req, r); err != nil {
		writeResponse(ctx, w, http.StatusBadRequest, map[string]interface{}{errorJSONKey: err.Error()})
		return
	}
	if err := req.Validate(); err != nil {
		writeResponse(ctx, w, http.StatusBadRequest, map[string]interface{}{errorJSONKey: err.Error()})
		return
	}
	resp, err := g.iamClient.AttachRole(ctx, req)
	if err != nil {
		writeResponse(ctx, w, http.StatusInternalServerError, map[string]interface{}{errorJSONKey: err.Error()})
		return
	}
	writeResponse(ctx, w, http.StatusOK, map[string]interface{}{successJSONKey: resp})
}

func (g *gatewayService) detachRoleIamHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	req := &iam.DetachRoleRequest{}
	if err := bind(req, r); err != nil {
		writeResponse(ctx, w, http.StatusBadRequest, map[string]interface{}{errorJSONKey: err.Error()})
		return
	}
	if err := req.Validate(); err != nil {
		writeResponse(ctx, w, http.StatusBadRequest, map[string]interface{}{errorJSONKey: err.Error()})
		return
	}
	resp, err := g.iamClient.DetachRole(ctx, req)
	if err != nil {
		writeResponse(ctx, w, http.StatusInternalServerError, map[string]interface{}{errorJSONKey: err.Error()})
		return
	}
	writeResponse(ctx, w, http.StatusOK, map[string]interface{}{successJSONKey: resp})
}

func (g *gatewayService) listPolicyIamHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	req := &iam.ListPolicyRequest{}
	if err := bind(req, r); err != nil {
		writeResponse(ctx, w, http.StatusBadRequest, map[string]interface{}{errorJSONKey: err.Error()})
		return
	}
	if err := req.Validate(); err != nil {
		writeResponse(ctx, w, http.StatusBadRequest, map[string]interface{}{errorJSONKey: err.Error()})
		return
	}
	resp, err := g.iamClient.ListPolicy(ctx, req)
	if err != nil {
		writeResponse(ctx, w, http.StatusInternalServerError, map[string]interface{}{errorJSONKey: err.Error()})
		return
	}
	writeResponse(ctx, w, http.StatusOK, map[string]interface{}{successJSONKey: resp})
}

func (g *gatewayService) getPolicyIamHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	req := &iam.GetPolicyRequest{}
	if err := bind(req, r); err != nil {
		writeResponse(ctx, w, http.StatusBadRequest, map[string]interface{}{errorJSONKey: err.Error()})
		return
	}
	if err := req.Validate(); err != nil {
		writeResponse(ctx, w, http.StatusBadRequest, map[string]interface{}{errorJSONKey: err.Error()})
		return
	}
	resp, err := g.iamClient.GetPolicy(ctx, req)
	if err != nil {
		writeResponse(ctx, w, http.StatusInternalServerError, map[string]interface{}{errorJSONKey: err.Error()})
		return
	}
	writeResponse(ctx, w, http.StatusOK, map[string]interface{}{successJSONKey: resp})
}

func (g *gatewayService) putPolicyIamHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	req := &iam.PutPolicyRequest{}
	if err := bind(req, r); err != nil {
		writeResponse(ctx, w, http.StatusBadRequest, map[string]interface{}{errorJSONKey: err.Error()})
		return
	}
	if err := req.Validate(); err != nil {
		writeResponse(ctx, w, http.StatusBadRequest, map[string]interface{}{errorJSONKey: err.Error()})
		return
	}
	resp, err := g.iamClient.PutPolicy(ctx, req)
	if err != nil {
		writeResponse(ctx, w, http.StatusInternalServerError, map[string]interface{}{errorJSONKey: err.Error()})
		return
	}
	writeResponse(ctx, w, http.StatusOK, map[string]interface{}{successJSONKey: resp})
}

func (g *gatewayService) deletePolicyIamHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	req := &iam.DeletePolicyRequest{}
	if err := bind(req, r); err != nil {
		writeResponse(ctx, w, http.StatusBadRequest, map[string]interface{}{errorJSONKey: err.Error()})
		return
	}
	if err := req.Validate(); err != nil {
		writeResponse(ctx, w, http.StatusBadRequest, map[string]interface{}{errorJSONKey: err.Error()})
		return
	}
	resp, err := g.iamClient.DeletePolicy(ctx, req)
	if err != nil {
		writeResponse(ctx, w, http.StatusInternalServerError, map[string]interface{}{errorJSONKey: err.Error()})
		return
	}
	writeResponse(ctx, w, http.StatusOK, map[string]interface{}{successJSONKey: resp})
}

func (g *gatewayService) attachPolicyIamHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	req := &iam.AttachPolicyRequest{}
	if err := bind(req, r); err != nil {
		writeResponse(ctx, w, http.StatusBadRequest, map[string]interface{}{errorJSONKey: err.Error()})
		return
	}
	if err := req.Validate(); err != nil {
		writeResponse(ctx, w, http.StatusBadRequest, map[string]interface{}{errorJSONKey: err.Error()})
		return
	}
	resp, err := g.iamClient.AttachPolicy(ctx, req)
	if err != nil {
		writeResponse(ctx, w, http.StatusInternalServerError, map[string]interface{}{errorJSONKey: err.Error()})
		return
	}
	writeResponse(ctx, w, http.StatusOK, map[string]interface{}{successJSONKey: resp})
}

func (g *gatewayService) detachPolicyIamHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	req := &iam.DetachPolicyRequest{}
	if err := bind(req, r); err != nil {
		writeResponse(ctx, w, http.StatusBadRequest, map[string]interface{}{errorJSONKey: err.Error()})
		return
	}
	if err := req.Validate(); err != nil {
		writeResponse(ctx, w, http.StatusBadRequest, map[string]interface{}{errorJSONKey: err.Error()})
		return
	}
	resp, err := g.iamClient.DetachPolicy(ctx, req)
	if err != nil {
		writeResponse(ctx, w, http.StatusInternalServerError, map[string]interface{}{errorJSONKey: err.Error()})
		return
	}
	writeResponse(ctx, w, http.StatusOK, map[string]interface{}{successJSONKey: resp})
}

func (g *gatewayService) listAccessTokenIamHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	req := &iam.ListAccessTokenRequest{}
	if err := bind(req, r); err != nil {
		writeResponse(ctx, w, http.StatusBadRequest, map[string]interface{}{errorJSONKey: err.Error()})
		return
	}
	if err := req.Validate(); err != nil {
		writeResponse(ctx, w, http.StatusBadRequest, map[string]interface{}{errorJSONKey: err.Error()})
		return
	}
	resp, err := g.iamClient.ListAccessToken(ctx, req)
	if err != nil {
		writeResponse(ctx, w, http.StatusInternalServerError, map[string]interface{}{errorJSONKey: err.Error()})
		return
	}
	writeResponse(ctx, w, http.StatusOK, map[string]interface{}{successJSONKey: resp})
}

func (g *gatewayService) deleteAccessTokenIamHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	req := &iam.DeleteAccessTokenRequest{}
	if err := bind(req, r); err != nil {
		writeResponse(ctx, w, http.StatusBadRequest, map[string]interface{}{errorJSONKey: err.Error()})
		return
	}
	if err := req.Validate(); err != nil {
		writeResponse(ctx, w, http.StatusBadRequest, map[string]interface{}{errorJSONKey: err.Error()})
		return
	}
	resp, err := g.iamClient.DeleteAccessToken(ctx, req)
	if err != nil {
		writeResponse(ctx, w, http.StatusInternalServerError, map[string]interface{}{errorJSONKey: err.Error()})
		return
	}
	writeResponse(ctx, w, http.StatusOK, map[string]interface{}{successJSONKey: resp})
}

func (g *gatewayService) attachAccessTokenRoleIamHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	req := &iam.AttachAccessTokenRoleRequest{}
	if err := bind(req, r); err != nil {
		writeResponse(ctx, w, http.StatusBadRequest, map[string]interface{}{errorJSONKey: err.Error()})
		return
	}
	if err := req.Validate(); err != nil {
		writeResponse(ctx, w, http.StatusBadRequest, map[string]interface{}{errorJSONKey: err.Error()})
		return
	}
	resp, err := g.iamClient.AttachAccessTokenRole(ctx, req)
	if err != nil {
		writeResponse(ctx, w, http.StatusInternalServerError, map[string]interface{}{errorJSONKey: err.Error()})
		return
	}
	writeResponse(ctx, w, http.StatusOK, map[string]interface{}{successJSONKey: resp})
}

func (g *gatewayService) detachAccessTokenRoleIamHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	req := &iam.DetachAccessTokenRoleRequest{}
	if err := bind(req, r); err != nil {
		writeResponse(ctx, w, http.StatusBadRequest, map[string]interface{}{errorJSONKey: err.Error()})
		return
	}
	if err := req.Validate(); err != nil {
		writeResponse(ctx, w, http.StatusBadRequest, map[string]interface{}{errorJSONKey: err.Error()})
		return
	}
	resp, err := g.iamClient.DetachAccessTokenRole(ctx, req)
	if err != nil {
		writeResponse(ctx, w, http.StatusInternalServerError, map[string]interface{}{errorJSONKey: err.Error()})
		return
	}
	writeResponse(ctx, w, http.StatusOK, map[string]interface{}{successJSONKey: resp})
}

func (g *gatewayService) isAdminIamHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	req := &iam.IsAdminRequest{}
	if err := bind(req, r); err != nil {
		writeResponse(ctx, w, http.StatusBadRequest, map[string]interface{}{errorJSONKey: err.Error()})
		return
	}
	if err := req.Validate(); err != nil {
		writeResponse(ctx, w, http.StatusBadRequest, map[string]interface{}{errorJSONKey: err.Error()})
		return
	}
	resp, err := g.iamClient.IsAdmin(ctx, req)
	if err != nil {
		writeResponse(ctx, w, http.StatusInternalServerError, map[string]interface{}{errorJSONKey: err.Error()})
		return
	}
	writeResponse(ctx, w, http.StatusOK, map[string]interface{}{successJSONKey: resp})
}