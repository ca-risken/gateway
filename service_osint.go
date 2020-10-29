package main

import (
	"net/http"

	"github.com/CyberAgent/mimosa-osint/proto/osint"
)

func (g *gatewayService) listOsintHandler(w http.ResponseWriter, r *http.Request) {
	req := &osint.ListOsintRequest{}
	bind(req, r)
	if err := req.Validate(); err != nil {
		writeResponse(w, http.StatusBadRequest, map[string]interface{}{errorJSONKey: err.Error()})
		return
	}
	resp, err := g.osintClient.ListOsint(r.Context(), req)
	if err != nil {
		writeResponse(w, http.StatusInternalServerError, map[string]interface{}{errorJSONKey: err.Error()})
		return
	}
	writeResponse(w, http.StatusOK, map[string]interface{}{successJSONKey: resp})
}

func (g *gatewayService) getOsintHandler(w http.ResponseWriter, r *http.Request) {
	req := &osint.GetOsintRequest{}
	bind(req, r)
	if err := req.Validate(); err != nil {
		writeResponse(w, http.StatusBadRequest, map[string]interface{}{errorJSONKey: err.Error()})
		return
	}
	resp, err := g.osintClient.GetOsint(r.Context(), req)
	if err != nil {
		writeResponse(w, http.StatusInternalServerError, map[string]interface{}{errorJSONKey: err.Error()})
		return
	}
	writeResponse(w, http.StatusOK, map[string]interface{}{successJSONKey: resp})
}

func (g *gatewayService) putOsintHandler(w http.ResponseWriter, r *http.Request) {
	req := &osint.PutOsintRequest{}
	bind(req, r)
	if err := req.Validate(); err != nil {
		writeResponse(w, http.StatusBadRequest, map[string]interface{}{errorJSONKey: err.Error()})
		return
	}
	resp, err := g.osintClient.PutOsint(r.Context(), req)
	if err != nil {
		writeResponse(w, http.StatusInternalServerError, map[string]interface{}{errorJSONKey: err.Error()})
		return
	}
	writeResponse(w, http.StatusOK, map[string]interface{}{successJSONKey: resp})
}

func (g *gatewayService) deleteOsintHandler(w http.ResponseWriter, r *http.Request) {
	req := &osint.DeleteOsintRequest{}
	bind(req, r)
	if err := req.Validate(); err != nil {
		writeResponse(w, http.StatusBadRequest, map[string]interface{}{errorJSONKey: err.Error()})
		return
	}
	resp, err := g.osintClient.DeleteOsint(r.Context(), req)
	if err != nil {
		writeResponse(w, http.StatusInternalServerError, map[string]interface{}{errorJSONKey: err.Error()})
		return
	}
	writeResponse(w, http.StatusOK, map[string]interface{}{successJSONKey: resp})
}

func (g *gatewayService) listOsintDataSourceHandler(w http.ResponseWriter, r *http.Request) {
	req := &osint.ListOsintDataSourceRequest{}
	bind(req, r)
	if err := req.Validate(); err != nil {
		writeResponse(w, http.StatusBadRequest, map[string]interface{}{errorJSONKey: err.Error()})
		return
	}
	resp, err := g.osintClient.ListOsintDataSource(r.Context(), req)
	if err != nil {
		writeResponse(w, http.StatusInternalServerError, map[string]interface{}{errorJSONKey: err.Error()})
		return
	}
	writeResponse(w, http.StatusOK, map[string]interface{}{successJSONKey: resp})
}

func (g *gatewayService) getOsintDataSourceHandler(w http.ResponseWriter, r *http.Request) {
	req := &osint.GetOsintDataSourceRequest{}
	bind(req, r)
	if err := req.Validate(); err != nil {
		writeResponse(w, http.StatusBadRequest, map[string]interface{}{errorJSONKey: err.Error()})
		return
	}
	resp, err := g.osintClient.GetOsintDataSource(r.Context(), req)
	if err != nil {
		writeResponse(w, http.StatusInternalServerError, map[string]interface{}{errorJSONKey: err.Error()})
		return
	}
	writeResponse(w, http.StatusOK, map[string]interface{}{successJSONKey: resp})
}

func (g *gatewayService) putOsintDataSourceHandler(w http.ResponseWriter, r *http.Request) {
	req := &osint.PutOsintDataSourceRequest{}
	bind(req, r)
	if err := req.Validate(); err != nil {
		writeResponse(w, http.StatusBadRequest, map[string]interface{}{errorJSONKey: err.Error()})
		return
	}
	resp, err := g.osintClient.PutOsintDataSource(r.Context(), req)
	if err != nil {
		writeResponse(w, http.StatusInternalServerError, map[string]interface{}{errorJSONKey: err.Error()})
		return
	}
	writeResponse(w, http.StatusOK, map[string]interface{}{successJSONKey: resp})
}

func (g *gatewayService) deleteOsintDataSourceHandler(w http.ResponseWriter, r *http.Request) {
	req := &osint.DeleteOsintDataSourceRequest{}
	bind(req, r)
	if err := req.Validate(); err != nil {
		writeResponse(w, http.StatusBadRequest, map[string]interface{}{errorJSONKey: err.Error()})
		return
	}
	resp, err := g.osintClient.DeleteOsintDataSource(r.Context(), req)
	if err != nil {
		writeResponse(w, http.StatusInternalServerError, map[string]interface{}{errorJSONKey: err.Error()})
		return
	}
	writeResponse(w, http.StatusOK, map[string]interface{}{successJSONKey: resp})
}

func (g *gatewayService) listRelOsintDataSourceHandler(w http.ResponseWriter, r *http.Request) {
	req := &osint.ListRelOsintDataSourceRequest{}
	bind(req, r)
	if err := req.Validate(); err != nil {
		writeResponse(w, http.StatusBadRequest, map[string]interface{}{errorJSONKey: err.Error()})
		return
	}
	resp, err := g.osintClient.ListRelOsintDataSource(r.Context(), req)
	if err != nil {
		writeResponse(w, http.StatusInternalServerError, map[string]interface{}{errorJSONKey: err.Error()})
		return
	}
	writeResponse(w, http.StatusOK, map[string]interface{}{successJSONKey: resp})
}

func (g *gatewayService) getRelOsintDataSourceHandler(w http.ResponseWriter, r *http.Request) {
	req := &osint.GetRelOsintDataSourceRequest{}
	bind(req, r)
	if err := req.Validate(); err != nil {
		writeResponse(w, http.StatusBadRequest, map[string]interface{}{errorJSONKey: err.Error()})
		return
	}
	resp, err := g.osintClient.GetRelOsintDataSource(r.Context(), req)
	if err != nil {
		writeResponse(w, http.StatusInternalServerError, map[string]interface{}{errorJSONKey: err.Error()})
		return
	}
	writeResponse(w, http.StatusOK, map[string]interface{}{successJSONKey: resp})
}

func (g *gatewayService) putRelOsintDataSourceHandler(w http.ResponseWriter, r *http.Request) {
	req := &osint.PutRelOsintDataSourceRequest{}
	bind(req, r)
	if err := req.Validate(); err != nil {
		writeResponse(w, http.StatusBadRequest, map[string]interface{}{errorJSONKey: err.Error()})
		return
	}
	resp, err := g.osintClient.PutRelOsintDataSource(r.Context(), req)
	if err != nil {
		writeResponse(w, http.StatusInternalServerError, map[string]interface{}{errorJSONKey: err.Error()})
		return
	}
	writeResponse(w, http.StatusOK, map[string]interface{}{successJSONKey: resp})
}

func (g *gatewayService) deleteRelOsintDataSourceHandler(w http.ResponseWriter, r *http.Request) {
	req := &osint.DeleteRelOsintDataSourceRequest{}
	bind(req, r)
	if err := req.Validate(); err != nil {
		writeResponse(w, http.StatusBadRequest, map[string]interface{}{errorJSONKey: err.Error()})
		return
	}
	resp, err := g.osintClient.DeleteRelOsintDataSource(r.Context(), req)
	if err != nil {
		writeResponse(w, http.StatusInternalServerError, map[string]interface{}{errorJSONKey: err.Error()})
		return
	}
	writeResponse(w, http.StatusOK, map[string]interface{}{successJSONKey: resp})
}

func (g *gatewayService) listOsintDetectWordHandler(w http.ResponseWriter, r *http.Request) {
	req := &osint.ListOsintDetectWordRequest{}
	bind(req, r)
	if err := req.Validate(); err != nil {
		writeResponse(w, http.StatusBadRequest, map[string]interface{}{errorJSONKey: err.Error()})
		return
	}
	resp, err := g.osintClient.ListOsintDetectWord(r.Context(), req)
	if err != nil {
		writeResponse(w, http.StatusInternalServerError, map[string]interface{}{errorJSONKey: err.Error()})
		return
	}
	writeResponse(w, http.StatusOK, map[string]interface{}{successJSONKey: resp})
}

func (g *gatewayService) getOsintDetectWordHandler(w http.ResponseWriter, r *http.Request) {
	req := &osint.GetOsintDetectWordRequest{}
	bind(req, r)
	if err := req.Validate(); err != nil {
		writeResponse(w, http.StatusBadRequest, map[string]interface{}{errorJSONKey: err.Error()})
		return
	}
	resp, err := g.osintClient.GetOsintDetectWord(r.Context(), req)
	if err != nil {
		writeResponse(w, http.StatusInternalServerError, map[string]interface{}{errorJSONKey: err.Error()})
		return
	}
	writeResponse(w, http.StatusOK, map[string]interface{}{successJSONKey: resp})
}

func (g *gatewayService) putOsintDetectWordHandler(w http.ResponseWriter, r *http.Request) {
	req := &osint.PutOsintDetectWordRequest{}
	bind(req, r)
	if err := req.Validate(); err != nil {
		writeResponse(w, http.StatusBadRequest, map[string]interface{}{errorJSONKey: err.Error()})
		return
	}
	resp, err := g.osintClient.PutOsintDetectWord(r.Context(), req)
	if err != nil {
		writeResponse(w, http.StatusInternalServerError, map[string]interface{}{errorJSONKey: err.Error()})
		return
	}
	writeResponse(w, http.StatusOK, map[string]interface{}{successJSONKey: resp})
}

func (g *gatewayService) deleteOsintDetectWordHandler(w http.ResponseWriter, r *http.Request) {
	req := &osint.DeleteOsintDetectWordRequest{}
	bind(req, r)
	if err := req.Validate(); err != nil {
		writeResponse(w, http.StatusBadRequest, map[string]interface{}{errorJSONKey: err.Error()})
		return
	}
	resp, err := g.osintClient.DeleteOsintDetectWord(r.Context(), req)
	if err != nil {
		writeResponse(w, http.StatusInternalServerError, map[string]interface{}{errorJSONKey: err.Error()})
		return
	}
	writeResponse(w, http.StatusOK, map[string]interface{}{successJSONKey: resp})
}

func (g *gatewayService) listRelOsintDetectWordHandler(w http.ResponseWriter, r *http.Request) {
	req := &osint.ListRelOsintDetectWordRequest{}
	bind(req, r)
	if err := req.Validate(); err != nil {
		writeResponse(w, http.StatusBadRequest, map[string]interface{}{errorJSONKey: err.Error()})
		return
	}
	resp, err := g.osintClient.ListRelOsintDetectWord(r.Context(), req)
	if err != nil {
		writeResponse(w, http.StatusInternalServerError, map[string]interface{}{errorJSONKey: err.Error()})
		return
	}
	writeResponse(w, http.StatusOK, map[string]interface{}{successJSONKey: resp})
}

func (g *gatewayService) getRelOsintDetectWordHandler(w http.ResponseWriter, r *http.Request) {
	req := &osint.GetRelOsintDetectWordRequest{}
	bind(req, r)
	if err := req.Validate(); err != nil {
		writeResponse(w, http.StatusBadRequest, map[string]interface{}{errorJSONKey: err.Error()})
		return
	}
	resp, err := g.osintClient.GetRelOsintDetectWord(r.Context(), req)
	if err != nil {
		writeResponse(w, http.StatusInternalServerError, map[string]interface{}{errorJSONKey: err.Error()})
		return
	}
	writeResponse(w, http.StatusOK, map[string]interface{}{successJSONKey: resp})
}

func (g *gatewayService) putRelOsintDetectWordHandler(w http.ResponseWriter, r *http.Request) {
	req := &osint.PutRelOsintDetectWordRequest{}
	bind(req, r)
	if err := req.Validate(); err != nil {
		writeResponse(w, http.StatusBadRequest, map[string]interface{}{errorJSONKey: err.Error()})
		return
	}
	resp, err := g.osintClient.PutRelOsintDetectWord(r.Context(), req)
	if err != nil {
		writeResponse(w, http.StatusInternalServerError, map[string]interface{}{errorJSONKey: err.Error()})
		return
	}
	writeResponse(w, http.StatusOK, map[string]interface{}{successJSONKey: resp})
}

func (g *gatewayService) deleteRelOsintDetectWordHandler(w http.ResponseWriter, r *http.Request) {
	req := &osint.DeleteRelOsintDetectWordRequest{}
	bind(req, r)
	if err := req.Validate(); err != nil {
		writeResponse(w, http.StatusBadRequest, map[string]interface{}{errorJSONKey: err.Error()})
		return
	}
	resp, err := g.osintClient.DeleteRelOsintDetectWord(r.Context(), req)
	if err != nil {
		writeResponse(w, http.StatusInternalServerError, map[string]interface{}{errorJSONKey: err.Error()})
		return
	}
	writeResponse(w, http.StatusOK, map[string]interface{}{successJSONKey: resp})
}

func (g *gatewayService) invokeOsintScanHandler(w http.ResponseWriter, r *http.Request) {
	req := &osint.InvokeScanRequest{}
	bind(req, r)
	if err := req.Validate(); err != nil {
		writeResponse(w, http.StatusBadRequest, map[string]interface{}{errorJSONKey: err.Error()})
		return
	}
	resp, err := g.osintClient.InvokeScan(r.Context(), req)
	if err != nil {
		writeResponse(w, http.StatusInternalServerError, map[string]interface{}{errorJSONKey: err.Error()})
		return
	}
	writeResponse(w, http.StatusOK, map[string]interface{}{successJSONKey: resp})
}
