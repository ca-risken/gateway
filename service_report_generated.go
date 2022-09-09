package main

import (
	"net/http"

	"github.com/ca-risken/core/proto/report"
)

func (g *gatewayService) getReportFindingHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	req := &report.GetReportFindingRequest{}
	bind(req, r)
	if err := req.Validate(); err != nil {
		writeResponse(ctx, w, http.StatusBadRequest, map[string]interface{}{errorJSONKey: err.Error()})
		return
	}
	resp, err := g.reportClient.GetReportFinding(ctx, req)
	if err != nil {
		writeResponse(ctx, w, http.StatusInternalServerError, map[string]interface{}{errorJSONKey: err.Error()})
		return
	}
	writeResponse(ctx, w, http.StatusOK, map[string]interface{}{successJSONKey: resp})
}

func (g *gatewayService) getReportFindingAllHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	req := &report.GetReportFindingAllRequest{}
	bind(req, r)
	if err := req.Validate(); err != nil {
		writeResponse(ctx, w, http.StatusBadRequest, map[string]interface{}{errorJSONKey: err.Error()})
		return
	}
	resp, err := g.reportClient.GetReportFindingAll(ctx, req)
	if err != nil {
		writeResponse(ctx, w, http.StatusInternalServerError, map[string]interface{}{errorJSONKey: err.Error()})
		return
	}
	writeResponse(ctx, w, http.StatusOK, map[string]interface{}{successJSONKey: resp})
}
