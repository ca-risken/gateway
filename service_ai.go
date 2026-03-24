package main

import (
	"net/http"

	"github.com/ca-risken/core/proto/ai"
)

func (g *gatewayService) chatAIAiHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	req := &ai.ChatAIRequest{}
	if err := bind(req, r); err != nil {
		appLogger.Warnf(ctx, "Failed to bind request, req=%s, err=%+v", "ChatAIRequest", err)
	}

	u, err := getRequestUser(r)
	if err != nil {
		appLogger.Infof(ctx, "Unauthenticated: %+v", err)
		http.Error(w, "Unauthenticated", http.StatusUnauthorized)
		return
	}

	if isHumanAccess(u) {
		if err := req.Validate(); err != nil {
			writeResponse(ctx, w, http.StatusBadRequest, map[string]interface{}{errorJSONKey: err.Error()})
			return
		}
		if req.ProjectId != 0 && !g.isAuthorizedProject(ctx, u.userID, req.ProjectId, r.URL.Path) {
			http.Error(w, "Unauthorized the project resource for human access", http.StatusForbidden)
			return
		}
	} else {
		if u.accessTokenID == 0 || u.accessTokenProjectID == 0 {
			http.Error(w, "Unauthorized the project resource for token access", http.StatusForbidden)
			return
		}
		req.ProjectId = u.accessTokenProjectID
		if err := req.Validate(); err != nil {
			writeResponse(ctx, w, http.StatusBadRequest, map[string]interface{}{errorJSONKey: err.Error()})
			return
		}
		if !g.isAuthorizedProjectToken(ctx, u.accessTokenID, u.accessTokenProjectID, r.URL.Path) {
			http.Error(w, "Unauthorized the project resource for token access", http.StatusForbidden)
			return
		}
	}

	resp, err := g.aiClient.ChatAI(ctx, req)
	if err != nil {
		if handleErr := handleGRPCError(ctx, w, err); handleErr != nil {
			appLogger.Errorf(ctx, "HandleGRPCError: %+v", handleErr)
			writeResponse(ctx, w, http.StatusInternalServerError, map[string]interface{}{errorJSONKey: "InternalServerError"})
		}
		return
	}
	writeResponse(ctx, w, http.StatusOK, map[string]interface{}{successJSONKey: resp})
}
