// Code generated by protoc-gen-service. DO NOT EDIT.
// source: google/service.proto

package main

import (
	"net/http"

	"github.com/ca-risken/datasource-api/proto/google"
)

func (g *gatewayService) listGoogleDataSourceGoogleHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	req := &google.ListGoogleDataSourceRequest{}
	if err := bind(req, r); err != nil {
		writeResponse(ctx, w, http.StatusBadRequest, map[string]interface{}{errorJSONKey: err.Error()})
		return
	}
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

func (g *gatewayService) listGCPGoogleHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	req := &google.ListGCPRequest{}
	if err := bind(req, r); err != nil {
		writeResponse(ctx, w, http.StatusBadRequest, map[string]interface{}{errorJSONKey: err.Error()})
		return
	}
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

func (g *gatewayService) getGCPGoogleHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	req := &google.GetGCPRequest{}
	if err := bind(req, r); err != nil {
		writeResponse(ctx, w, http.StatusBadRequest, map[string]interface{}{errorJSONKey: err.Error()})
		return
	}
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

func (g *gatewayService) putGCPGoogleHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	req := &google.PutGCPRequest{}
	if err := bind(req, r); err != nil {
		writeResponse(ctx, w, http.StatusBadRequest, map[string]interface{}{errorJSONKey: err.Error()})
		return
	}
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

func (g *gatewayService) deleteGCPGoogleHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	req := &google.DeleteGCPRequest{}
	if err := bind(req, r); err != nil {
		writeResponse(ctx, w, http.StatusBadRequest, map[string]interface{}{errorJSONKey: err.Error()})
		return
	}
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

func (g *gatewayService) listGCPDataSourceGoogleHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	req := &google.ListGCPDataSourceRequest{}
	if err := bind(req, r); err != nil {
		writeResponse(ctx, w, http.StatusBadRequest, map[string]interface{}{errorJSONKey: err.Error()})
		return
	}
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

func (g *gatewayService) getGCPDataSourceGoogleHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	req := &google.GetGCPDataSourceRequest{}
	if err := bind(req, r); err != nil {
		writeResponse(ctx, w, http.StatusBadRequest, map[string]interface{}{errorJSONKey: err.Error()})
		return
	}
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

func (g *gatewayService) attachGCPDataSourceGoogleHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	req := &google.AttachGCPDataSourceRequest{}
	if err := bind(req, r); err != nil {
		writeResponse(ctx, w, http.StatusBadRequest, map[string]interface{}{errorJSONKey: err.Error()})
		return
	}
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

func (g *gatewayService) detachGCPDataSourceGoogleHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	req := &google.DetachGCPDataSourceRequest{}
	if err := bind(req, r); err != nil {
		writeResponse(ctx, w, http.StatusBadRequest, map[string]interface{}{errorJSONKey: err.Error()})
		return
	}
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

func (g *gatewayService) invokeScanGCPGoogleHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	req := &google.InvokeScanGCPRequest{}
	if err := bind(req, r); err != nil {
		writeResponse(ctx, w, http.StatusBadRequest, map[string]interface{}{errorJSONKey: err.Error()})
		return
	}
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