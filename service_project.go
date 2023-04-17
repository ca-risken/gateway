package main

import (
	"errors"
	"net/http"

	"github.com/ca-risken/core/proto/project"
)

func (g *gatewayService) createProjectHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	user, err := getRequestUser(r)
	if err != nil {
		writeResponse(ctx, w, http.StatusUnauthorized, map[string]interface{}{errorJSONKey: errors.New("InvalidUser")})
	}
	req := &project.CreateProjectRequest{}
	req.UserId = user.userID // force update by own userID
	if err := bind(req, r); err != nil {
		appLogger.Warnf(ctx, "Failed to bind request, req=%s, err=%+v", "CreateProjectRequest", err)
	}
	if err := req.Validate(); err != nil {
		writeResponse(ctx, w, http.StatusBadRequest, map[string]interface{}{errorJSONKey: err.Error()})
		return
	}
	resp, err := g.projectClient.CreateProject(ctx, req)
	if err != nil {
		if handleErr := handleGRPCError(ctx, w, err); handleErr != nil {
			appLogger.Errorf(ctx, "HandleGRPCError: %+v", handleErr)
			writeResponse(ctx, w, http.StatusInternalServerError, map[string]interface{}{errorJSONKey: "InternalServerError"})
		}
		return
	}
	writeResponse(ctx, w, http.StatusOK, map[string]interface{}{successJSONKey: resp})
}
