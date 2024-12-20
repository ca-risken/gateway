// Code generated by protoc-gen-service. DO NOT EDIT.
// source: osint/services.proto

package main

import (
	"net/http"

	"github.com/ca-risken/datasource-api/proto/osint"
)

func (g *gatewayService) listOsintOsintHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	req := &osint.ListOsintRequest{}
	if err := bind(req, r); err != nil {
		appLogger.Warnf(ctx, "Failed to bind request, req=%s, err=%+v", "ListOsintRequest", err)
	}
	if err := req.Validate(); err != nil {
		writeResponse(ctx, w, http.StatusBadRequest, map[string]interface{}{errorJSONKey: err.Error()})
		return
	}
	resp, err := g.osintClient.ListOsint(ctx, req)
	if err != nil {
		if handleErr := handleGRPCError(ctx, w, err); handleErr != nil {
			appLogger.Errorf(ctx, "HandleGRPCError: %+v", handleErr)
			writeResponse(ctx, w, http.StatusInternalServerError, map[string]interface{}{errorJSONKey: "InternalServerError"})
		}
		return
	}
	writeResponse(ctx, w, http.StatusOK, map[string]interface{}{successJSONKey: resp})
}

func (g *gatewayService) getOsintOsintHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	req := &osint.GetOsintRequest{}
	if err := bind(req, r); err != nil {
		appLogger.Warnf(ctx, "Failed to bind request, req=%s, err=%+v", "GetOsintRequest", err)
	}
	if err := req.Validate(); err != nil {
		writeResponse(ctx, w, http.StatusBadRequest, map[string]interface{}{errorJSONKey: err.Error()})
		return
	}
	resp, err := g.osintClient.GetOsint(ctx, req)
	if err != nil {
		if handleErr := handleGRPCError(ctx, w, err); handleErr != nil {
			appLogger.Errorf(ctx, "HandleGRPCError: %+v", handleErr)
			writeResponse(ctx, w, http.StatusInternalServerError, map[string]interface{}{errorJSONKey: "InternalServerError"})
		}
		return
	}
	writeResponse(ctx, w, http.StatusOK, map[string]interface{}{successJSONKey: resp})
}

func (g *gatewayService) putOsintOsintHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	req := &osint.PutOsintRequest{}
	if err := bind(req, r); err != nil {
		appLogger.Warnf(ctx, "Failed to bind request, req=%s, err=%+v", "PutOsintRequest", err)
	}
	if err := req.Validate(); err != nil {
		writeResponse(ctx, w, http.StatusBadRequest, map[string]interface{}{errorJSONKey: err.Error()})
		return
	}
	resp, err := g.osintClient.PutOsint(ctx, req)
	if err != nil {
		if handleErr := handleGRPCError(ctx, w, err); handleErr != nil {
			appLogger.Errorf(ctx, "HandleGRPCError: %+v", handleErr)
			writeResponse(ctx, w, http.StatusInternalServerError, map[string]interface{}{errorJSONKey: "InternalServerError"})
		}
		return
	}
	writeResponse(ctx, w, http.StatusOK, map[string]interface{}{successJSONKey: resp})
}

func (g *gatewayService) deleteOsintOsintHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	req := &osint.DeleteOsintRequest{}
	if err := bind(req, r); err != nil {
		appLogger.Warnf(ctx, "Failed to bind request, req=%s, err=%+v", "DeleteOsintRequest", err)
	}
	if err := req.Validate(); err != nil {
		writeResponse(ctx, w, http.StatusBadRequest, map[string]interface{}{errorJSONKey: err.Error()})
		return
	}
	resp, err := g.osintClient.DeleteOsint(ctx, req)
	if err != nil {
		if handleErr := handleGRPCError(ctx, w, err); handleErr != nil {
			appLogger.Errorf(ctx, "HandleGRPCError: %+v", handleErr)
			writeResponse(ctx, w, http.StatusInternalServerError, map[string]interface{}{errorJSONKey: "InternalServerError"})
		}
		return
	}
	writeResponse(ctx, w, http.StatusOK, map[string]interface{}{successJSONKey: resp})
}

func (g *gatewayService) listRelOsintDataSourceOsintHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	req := &osint.ListRelOsintDataSourceRequest{}
	if err := bind(req, r); err != nil {
		appLogger.Warnf(ctx, "Failed to bind request, req=%s, err=%+v", "ListRelOsintDataSourceRequest", err)
	}
	if err := req.Validate(); err != nil {
		writeResponse(ctx, w, http.StatusBadRequest, map[string]interface{}{errorJSONKey: err.Error()})
		return
	}
	resp, err := g.osintClient.ListRelOsintDataSource(ctx, req)
	if err != nil {
		if handleErr := handleGRPCError(ctx, w, err); handleErr != nil {
			appLogger.Errorf(ctx, "HandleGRPCError: %+v", handleErr)
			writeResponse(ctx, w, http.StatusInternalServerError, map[string]interface{}{errorJSONKey: "InternalServerError"})
		}
		return
	}
	writeResponse(ctx, w, http.StatusOK, map[string]interface{}{successJSONKey: resp})
}

func (g *gatewayService) getRelOsintDataSourceOsintHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	req := &osint.GetRelOsintDataSourceRequest{}
	if err := bind(req, r); err != nil {
		appLogger.Warnf(ctx, "Failed to bind request, req=%s, err=%+v", "GetRelOsintDataSourceRequest", err)
	}
	if err := req.Validate(); err != nil {
		writeResponse(ctx, w, http.StatusBadRequest, map[string]interface{}{errorJSONKey: err.Error()})
		return
	}
	resp, err := g.osintClient.GetRelOsintDataSource(ctx, req)
	if err != nil {
		if handleErr := handleGRPCError(ctx, w, err); handleErr != nil {
			appLogger.Errorf(ctx, "HandleGRPCError: %+v", handleErr)
			writeResponse(ctx, w, http.StatusInternalServerError, map[string]interface{}{errorJSONKey: "InternalServerError"})
		}
		return
	}
	writeResponse(ctx, w, http.StatusOK, map[string]interface{}{successJSONKey: resp})
}

func (g *gatewayService) putRelOsintDataSourceOsintHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	req := &osint.PutRelOsintDataSourceRequest{}
	if err := bind(req, r); err != nil {
		appLogger.Warnf(ctx, "Failed to bind request, req=%s, err=%+v", "PutRelOsintDataSourceRequest", err)
	}
	if err := req.Validate(); err != nil {
		writeResponse(ctx, w, http.StatusBadRequest, map[string]interface{}{errorJSONKey: err.Error()})
		return
	}
	resp, err := g.osintClient.PutRelOsintDataSource(ctx, req)
	if err != nil {
		if handleErr := handleGRPCError(ctx, w, err); handleErr != nil {
			appLogger.Errorf(ctx, "HandleGRPCError: %+v", handleErr)
			writeResponse(ctx, w, http.StatusInternalServerError, map[string]interface{}{errorJSONKey: "InternalServerError"})
		}
		return
	}
	writeResponse(ctx, w, http.StatusOK, map[string]interface{}{successJSONKey: resp})
}

func (g *gatewayService) deleteRelOsintDataSourceOsintHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	req := &osint.DeleteRelOsintDataSourceRequest{}
	if err := bind(req, r); err != nil {
		appLogger.Warnf(ctx, "Failed to bind request, req=%s, err=%+v", "DeleteRelOsintDataSourceRequest", err)
	}
	if err := req.Validate(); err != nil {
		writeResponse(ctx, w, http.StatusBadRequest, map[string]interface{}{errorJSONKey: err.Error()})
		return
	}
	resp, err := g.osintClient.DeleteRelOsintDataSource(ctx, req)
	if err != nil {
		if handleErr := handleGRPCError(ctx, w, err); handleErr != nil {
			appLogger.Errorf(ctx, "HandleGRPCError: %+v", handleErr)
			writeResponse(ctx, w, http.StatusInternalServerError, map[string]interface{}{errorJSONKey: "InternalServerError"})
		}
		return
	}
	writeResponse(ctx, w, http.StatusOK, map[string]interface{}{successJSONKey: resp})
}

func (g *gatewayService) listOsintDataSourceOsintHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	req := &osint.ListOsintDataSourceRequest{}
	if err := bind(req, r); err != nil {
		appLogger.Warnf(ctx, "Failed to bind request, req=%s, err=%+v", "ListOsintDataSourceRequest", err)
	}
	if err := req.Validate(); err != nil {
		writeResponse(ctx, w, http.StatusBadRequest, map[string]interface{}{errorJSONKey: err.Error()})
		return
	}
	resp, err := g.osintClient.ListOsintDataSource(ctx, req)
	if err != nil {
		if handleErr := handleGRPCError(ctx, w, err); handleErr != nil {
			appLogger.Errorf(ctx, "HandleGRPCError: %+v", handleErr)
			writeResponse(ctx, w, http.StatusInternalServerError, map[string]interface{}{errorJSONKey: "InternalServerError"})
		}
		return
	}
	writeResponse(ctx, w, http.StatusOK, map[string]interface{}{successJSONKey: resp})
}

func (g *gatewayService) getOsintDataSourceOsintHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	req := &osint.GetOsintDataSourceRequest{}
	if err := bind(req, r); err != nil {
		appLogger.Warnf(ctx, "Failed to bind request, req=%s, err=%+v", "GetOsintDataSourceRequest", err)
	}
	if err := req.Validate(); err != nil {
		writeResponse(ctx, w, http.StatusBadRequest, map[string]interface{}{errorJSONKey: err.Error()})
		return
	}
	resp, err := g.osintClient.GetOsintDataSource(ctx, req)
	if err != nil {
		if handleErr := handleGRPCError(ctx, w, err); handleErr != nil {
			appLogger.Errorf(ctx, "HandleGRPCError: %+v", handleErr)
			writeResponse(ctx, w, http.StatusInternalServerError, map[string]interface{}{errorJSONKey: "InternalServerError"})
		}
		return
	}
	writeResponse(ctx, w, http.StatusOK, map[string]interface{}{successJSONKey: resp})
}

func (g *gatewayService) putOsintDataSourceOsintHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	req := &osint.PutOsintDataSourceRequest{}
	if err := bind(req, r); err != nil {
		appLogger.Warnf(ctx, "Failed to bind request, req=%s, err=%+v", "PutOsintDataSourceRequest", err)
	}
	if err := req.Validate(); err != nil {
		writeResponse(ctx, w, http.StatusBadRequest, map[string]interface{}{errorJSONKey: err.Error()})
		return
	}
	resp, err := g.osintClient.PutOsintDataSource(ctx, req)
	if err != nil {
		if handleErr := handleGRPCError(ctx, w, err); handleErr != nil {
			appLogger.Errorf(ctx, "HandleGRPCError: %+v", handleErr)
			writeResponse(ctx, w, http.StatusInternalServerError, map[string]interface{}{errorJSONKey: "InternalServerError"})
		}
		return
	}
	writeResponse(ctx, w, http.StatusOK, map[string]interface{}{successJSONKey: resp})
}

func (g *gatewayService) deleteOsintDataSourceOsintHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	req := &osint.DeleteOsintDataSourceRequest{}
	if err := bind(req, r); err != nil {
		appLogger.Warnf(ctx, "Failed to bind request, req=%s, err=%+v", "DeleteOsintDataSourceRequest", err)
	}
	if err := req.Validate(); err != nil {
		writeResponse(ctx, w, http.StatusBadRequest, map[string]interface{}{errorJSONKey: err.Error()})
		return
	}
	resp, err := g.osintClient.DeleteOsintDataSource(ctx, req)
	if err != nil {
		if handleErr := handleGRPCError(ctx, w, err); handleErr != nil {
			appLogger.Errorf(ctx, "HandleGRPCError: %+v", handleErr)
			writeResponse(ctx, w, http.StatusInternalServerError, map[string]interface{}{errorJSONKey: "InternalServerError"})
		}
		return
	}
	writeResponse(ctx, w, http.StatusOK, map[string]interface{}{successJSONKey: resp})
}

func (g *gatewayService) listOsintDetectWordOsintHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	req := &osint.ListOsintDetectWordRequest{}
	if err := bind(req, r); err != nil {
		appLogger.Warnf(ctx, "Failed to bind request, req=%s, err=%+v", "ListOsintDetectWordRequest", err)
	}
	if err := req.Validate(); err != nil {
		writeResponse(ctx, w, http.StatusBadRequest, map[string]interface{}{errorJSONKey: err.Error()})
		return
	}
	resp, err := g.osintClient.ListOsintDetectWord(ctx, req)
	if err != nil {
		if handleErr := handleGRPCError(ctx, w, err); handleErr != nil {
			appLogger.Errorf(ctx, "HandleGRPCError: %+v", handleErr)
			writeResponse(ctx, w, http.StatusInternalServerError, map[string]interface{}{errorJSONKey: "InternalServerError"})
		}
		return
	}
	writeResponse(ctx, w, http.StatusOK, map[string]interface{}{successJSONKey: resp})
}

func (g *gatewayService) getOsintDetectWordOsintHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	req := &osint.GetOsintDetectWordRequest{}
	if err := bind(req, r); err != nil {
		appLogger.Warnf(ctx, "Failed to bind request, req=%s, err=%+v", "GetOsintDetectWordRequest", err)
	}
	if err := req.Validate(); err != nil {
		writeResponse(ctx, w, http.StatusBadRequest, map[string]interface{}{errorJSONKey: err.Error()})
		return
	}
	resp, err := g.osintClient.GetOsintDetectWord(ctx, req)
	if err != nil {
		if handleErr := handleGRPCError(ctx, w, err); handleErr != nil {
			appLogger.Errorf(ctx, "HandleGRPCError: %+v", handleErr)
			writeResponse(ctx, w, http.StatusInternalServerError, map[string]interface{}{errorJSONKey: "InternalServerError"})
		}
		return
	}
	writeResponse(ctx, w, http.StatusOK, map[string]interface{}{successJSONKey: resp})
}

func (g *gatewayService) putOsintDetectWordOsintHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	req := &osint.PutOsintDetectWordRequest{}
	if err := bind(req, r); err != nil {
		appLogger.Warnf(ctx, "Failed to bind request, req=%s, err=%+v", "PutOsintDetectWordRequest", err)
	}
	if err := req.Validate(); err != nil {
		writeResponse(ctx, w, http.StatusBadRequest, map[string]interface{}{errorJSONKey: err.Error()})
		return
	}
	resp, err := g.osintClient.PutOsintDetectWord(ctx, req)
	if err != nil {
		if handleErr := handleGRPCError(ctx, w, err); handleErr != nil {
			appLogger.Errorf(ctx, "HandleGRPCError: %+v", handleErr)
			writeResponse(ctx, w, http.StatusInternalServerError, map[string]interface{}{errorJSONKey: "InternalServerError"})
		}
		return
	}
	writeResponse(ctx, w, http.StatusOK, map[string]interface{}{successJSONKey: resp})
}

func (g *gatewayService) deleteOsintDetectWordOsintHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	req := &osint.DeleteOsintDetectWordRequest{}
	if err := bind(req, r); err != nil {
		appLogger.Warnf(ctx, "Failed to bind request, req=%s, err=%+v", "DeleteOsintDetectWordRequest", err)
	}
	if err := req.Validate(); err != nil {
		writeResponse(ctx, w, http.StatusBadRequest, map[string]interface{}{errorJSONKey: err.Error()})
		return
	}
	resp, err := g.osintClient.DeleteOsintDetectWord(ctx, req)
	if err != nil {
		if handleErr := handleGRPCError(ctx, w, err); handleErr != nil {
			appLogger.Errorf(ctx, "HandleGRPCError: %+v", handleErr)
			writeResponse(ctx, w, http.StatusInternalServerError, map[string]interface{}{errorJSONKey: "InternalServerError"})
		}
		return
	}
	writeResponse(ctx, w, http.StatusOK, map[string]interface{}{successJSONKey: resp})
}

func (g *gatewayService) invokeScanOsintHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	req := &osint.InvokeScanRequest{}
	if err := bind(req, r); err != nil {
		appLogger.Warnf(ctx, "Failed to bind request, req=%s, err=%+v", "InvokeScanRequest", err)
	}
	if err := req.Validate(); err != nil {
		writeResponse(ctx, w, http.StatusBadRequest, map[string]interface{}{errorJSONKey: err.Error()})
		return
	}
	resp, err := g.osintClient.InvokeScan(ctx, req)
	if err != nil {
		if handleErr := handleGRPCError(ctx, w, err); handleErr != nil {
			appLogger.Errorf(ctx, "HandleGRPCError: %+v", handleErr)
			writeResponse(ctx, w, http.StatusInternalServerError, map[string]interface{}{errorJSONKey: "InternalServerError"})
		}
		return
	}
	writeResponse(ctx, w, http.StatusOK, map[string]interface{}{successJSONKey: resp})
}
