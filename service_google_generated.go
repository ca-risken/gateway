package main

import (
	"net/http"

	"github.com/ca-risken/datasource-api/proto/google"
)

func (g *gatewayService) listGoogleDataSourceHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	req := &google.ListGoogleDataSourceRequest{}
	bind(req, r)
	if err := req.Validate(); err != nil {
		writeResponse(ctx, w, http.StatusBadRequest, map[string]interface{}{errorJSONKey: err.Error()})
		return
	}
	resp, err := g.googleClient.ListGoogleDataSource(ctx, req)
	if err != nil {
		writeResponse(ctx, w, http.StatusInternalServerError, map[string]interface{}{errorJSONKey: err.Error()})
		return
	}
	writeResponse(ctx, w, http.StatusOK, map[string]interface{}{successJSONKey: resp})
}

func (g *gatewayService) listGCPHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	req := &google.ListGCPRequest{}
	bind(req, r)
	if err := req.Validate(); err != nil {
		writeResponse(ctx, w, http.StatusBadRequest, map[string]interface{}{errorJSONKey: err.Error()})
		return
	}
	resp, err := g.googleClient.ListGCP(ctx, req)
	if err != nil {
		writeResponse(ctx, w, http.StatusInternalServerError, map[string]interface{}{errorJSONKey: err.Error()})
		return
	}
	writeResponse(ctx, w, http.StatusOK, map[string]interface{}{successJSONKey: resp})
}

func (g *gatewayService) getGCPHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	req := &google.GetGCPRequest{}
	bind(req, r)
	if err := req.Validate(); err != nil {
		writeResponse(ctx, w, http.StatusBadRequest, map[string]interface{}{errorJSONKey: err.Error()})
		return
	}
	resp, err := g.googleClient.GetGCP(ctx, req)
	if err != nil {
		writeResponse(ctx, w, http.StatusInternalServerError, map[string]interface{}{errorJSONKey: err.Error()})
		return
	}
	writeResponse(ctx, w, http.StatusOK, map[string]interface{}{successJSONKey: resp})
}

func (g *gatewayService) putGCPHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	req := &google.PutGCPRequest{}
	bind(req, r)
	if err := req.Validate(); err != nil {
		writeResponse(ctx, w, http.StatusBadRequest, map[string]interface{}{errorJSONKey: err.Error()})
		return
	}
	resp, err := g.googleClient.PutGCP(ctx, req)
	if err != nil {
		writeResponse(ctx, w, http.StatusInternalServerError, map[string]interface{}{errorJSONKey: err.Error()})
		return
	}
	writeResponse(ctx, w, http.StatusOK, map[string]interface{}{successJSONKey: resp})
}

func (g *gatewayService) deleteGCPHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	req := &google.DeleteGCPRequest{}
	bind(req, r)
	if err := req.Validate(); err != nil {
		writeResponse(ctx, w, http.StatusBadRequest, map[string]interface{}{errorJSONKey: err.Error()})
		return
	}
	resp, err := g.googleClient.DeleteGCP(ctx, req)
	if err != nil {
		writeResponse(ctx, w, http.StatusInternalServerError, map[string]interface{}{errorJSONKey: err.Error()})
		return
	}
	writeResponse(ctx, w, http.StatusOK, map[string]interface{}{successJSONKey: resp})
}

func (g *gatewayService) listGCPDataSourceHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	req := &google.ListGCPDataSourceRequest{}
	bind(req, r)
	if err := req.Validate(); err != nil {
		writeResponse(ctx, w, http.StatusBadRequest, map[string]interface{}{errorJSONKey: err.Error()})
		return
	}
	resp, err := g.googleClient.ListGCPDataSource(ctx, req)
	if err != nil {
		writeResponse(ctx, w, http.StatusInternalServerError, map[string]interface{}{errorJSONKey: err.Error()})
		return
	}
	writeResponse(ctx, w, http.StatusOK, map[string]interface{}{successJSONKey: resp})
}

func (g *gatewayService) getGCPDataSourceHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	req := &google.GetGCPDataSourceRequest{}
	bind(req, r)
	if err := req.Validate(); err != nil {
		writeResponse(ctx, w, http.StatusBadRequest, map[string]interface{}{errorJSONKey: err.Error()})
		return
	}
	resp, err := g.googleClient.GetGCPDataSource(ctx, req)
	if err != nil {
		writeResponse(ctx, w, http.StatusInternalServerError, map[string]interface{}{errorJSONKey: err.Error()})
		return
	}
	writeResponse(ctx, w, http.StatusOK, map[string]interface{}{successJSONKey: resp})
}

func (g *gatewayService) attachGCPDataSourceHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	req := &google.AttachGCPDataSourceRequest{}
	bind(req, r)
	if err := req.Validate(); err != nil {
		writeResponse(ctx, w, http.StatusBadRequest, map[string]interface{}{errorJSONKey: err.Error()})
		return
	}
	resp, err := g.googleClient.AttachGCPDataSource(ctx, req)
	if err != nil {
		writeResponse(ctx, w, http.StatusInternalServerError, map[string]interface{}{errorJSONKey: err.Error()})
		return
	}
	writeResponse(ctx, w, http.StatusOK, map[string]interface{}{successJSONKey: resp})
}

func (g *gatewayService) detachGCPDataSourceHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	req := &google.DetachGCPDataSourceRequest{}
	bind(req, r)
	if err := req.Validate(); err != nil {
		writeResponse(ctx, w, http.StatusBadRequest, map[string]interface{}{errorJSONKey: err.Error()})
		return
	}
	resp, err := g.googleClient.DetachGCPDataSource(ctx, req)
	if err != nil {
		writeResponse(ctx, w, http.StatusInternalServerError, map[string]interface{}{errorJSONKey: err.Error()})
		return
	}
	writeResponse(ctx, w, http.StatusOK, map[string]interface{}{successJSONKey: resp})
}

func (g *gatewayService) invokeScanGCPHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	req := &google.InvokeScanGCPRequest{}
	bind(req, r)
	if err := req.Validate(); err != nil {
		writeResponse(ctx, w, http.StatusBadRequest, map[string]interface{}{errorJSONKey: err.Error()})
		return
	}
	resp, err := g.googleClient.InvokeScanGCP(ctx, req)
	if err != nil {
		writeResponse(ctx, w, http.StatusInternalServerError, map[string]interface{}{errorJSONKey: err.Error()})
		return
	}
	writeResponse(ctx, w, http.StatusOK, map[string]interface{}{successJSONKey: resp})
}
