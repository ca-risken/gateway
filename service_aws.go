package main

import (
	"net/http"

	"github.com/ca-risken/datasource-api/proto/aws"
)

func (g *gatewayService) attachDataSourceHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	req := &aws.AttachDataSourceRequest{}
	if err := bind(req, r); err != nil {
		writeResponse(ctx, w, http.StatusBadRequest, map[string]interface{}{errorJSONKey: err.Error()})
		return
	}
	if err := req.ValidateForUser(); err != nil {
		writeResponse(ctx, w, http.StatusBadRequest, map[string]interface{}{errorJSONKey: err.Error()})
		return
	}
	resp, err := g.awsClient.AttachDataSource(ctx, req)
	if err != nil {
		writeResponse(ctx, w, http.StatusInternalServerError, map[string]interface{}{errorJSONKey: err.Error()})
		return
	}
	writeResponse(ctx, w, http.StatusOK, map[string]interface{}{successJSONKey: resp})
}
