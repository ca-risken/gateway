package main

import (
	"net/http"

	"github.com/ca-risken/core/proto/report"
)

func (g *gatewayService) getReportFindingForOrganizationReportHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	req := &report.GetReportFindingForOrganizationRequest{}
	if err := bind(req, r); err != nil {
		appLogger.Warnf(ctx, "Failed to bind request, req=%s, err=%+v", "GetReportFindingForOrganizationRequest", err)
		writeResponse(ctx, w, http.StatusBadRequest, map[string]interface{}{errorJSONKey: err.Error()})
		return
	}
	if err := req.Validate(); err != nil {
		writeResponse(ctx, w, http.StatusBadRequest, map[string]interface{}{errorJSONKey: err.Error()})
		return
	}
	resp, err := g.reportClient.GetReportFindingForOrganization(ctx, req)
	if err != nil {
		if handleErr := handleGRPCError(ctx, w, err); handleErr != nil {
			appLogger.Errorf(ctx, "HandleGRPCError: %+v", handleErr)
			writeResponse(ctx, w, http.StatusInternalServerError, map[string]interface{}{errorJSONKey: "InternalServerError"})
		}
		return
	}
	writeResponse(ctx, w, http.StatusOK, map[string]interface{}{successJSONKey: resp})
}
