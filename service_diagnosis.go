package main

import (
	"net/http"

	"github.com/CyberAgent/mimosa-diagnosis/pkg/pb/diagnosis"
)

func (g *gatewayService) listDiagnosisHandler(w http.ResponseWriter, r *http.Request) {
	req := &diagnosis.ListDiagnosisRequest{}
	bind(req, r)
	if err := req.Validate(); err != nil {
		writeResponse(w, http.StatusBadRequest, map[string]interface{}{errorJSONKey: err.Error()})
		return
	}
	resp, err := g.diagnosisClient.ListDiagnosis(r.Context(), req)
	if err != nil {
		writeResponse(w, http.StatusInternalServerError, map[string]interface{}{errorJSONKey: err.Error()})
		return
	}
	writeResponse(w, http.StatusOK, map[string]interface{}{successJSONKey: resp})
}

func (g *gatewayService) getDiagnosisHandler(w http.ResponseWriter, r *http.Request) {
	req := &diagnosis.GetDiagnosisRequest{}
	bind(req, r)
	if err := req.Validate(); err != nil {
		writeResponse(w, http.StatusBadRequest, map[string]interface{}{errorJSONKey: err.Error()})
		return
	}
	resp, err := g.diagnosisClient.GetDiagnosis(r.Context(), req)
	if err != nil {
		writeResponse(w, http.StatusInternalServerError, map[string]interface{}{errorJSONKey: err.Error()})
		return
	}
	writeResponse(w, http.StatusOK, map[string]interface{}{successJSONKey: resp})
}

func (g *gatewayService) putDiagnosisHandler(w http.ResponseWriter, r *http.Request) {
	req := &diagnosis.PutDiagnosisRequest{}
	bind(req, r)
	if err := req.Validate(); err != nil {
		writeResponse(w, http.StatusBadRequest, map[string]interface{}{errorJSONKey: err.Error()})
		return
	}
	resp, err := g.diagnosisClient.PutDiagnosis(r.Context(), req)
	if err != nil {
		writeResponse(w, http.StatusInternalServerError, map[string]interface{}{errorJSONKey: err.Error()})
		return
	}
	writeResponse(w, http.StatusOK, map[string]interface{}{successJSONKey: resp})
}

func (g *gatewayService) deleteDiagnosisHandler(w http.ResponseWriter, r *http.Request) {
	req := &diagnosis.DeleteDiagnosisRequest{}
	bind(req, r)
	if err := req.Validate(); err != nil {
		writeResponse(w, http.StatusBadRequest, map[string]interface{}{errorJSONKey: err.Error()})
		return
	}
	resp, err := g.diagnosisClient.DeleteDiagnosis(r.Context(), req)
	if err != nil {
		writeResponse(w, http.StatusInternalServerError, map[string]interface{}{errorJSONKey: err.Error()})
		return
	}
	writeResponse(w, http.StatusOK, map[string]interface{}{successJSONKey: resp})
}

func (g *gatewayService) listDiagnosisDataSourceHandler(w http.ResponseWriter, r *http.Request) {
	req := &diagnosis.ListDiagnosisDataSourceRequest{}
	bind(req, r)
	if err := req.Validate(); err != nil {
		writeResponse(w, http.StatusBadRequest, map[string]interface{}{errorJSONKey: err.Error()})
		return
	}
	resp, err := g.diagnosisClient.ListDiagnosisDataSource(r.Context(), req)
	if err != nil {
		writeResponse(w, http.StatusInternalServerError, map[string]interface{}{errorJSONKey: err.Error()})
		return
	}
	writeResponse(w, http.StatusOK, map[string]interface{}{successJSONKey: resp})
}

func (g *gatewayService) getDiagnosisDataSourceHandler(w http.ResponseWriter, r *http.Request) {
	req := &diagnosis.GetDiagnosisDataSourceRequest{}
	bind(req, r)
	if err := req.Validate(); err != nil {
		writeResponse(w, http.StatusBadRequest, map[string]interface{}{errorJSONKey: err.Error()})
		return
	}
	resp, err := g.diagnosisClient.GetDiagnosisDataSource(r.Context(), req)
	if err != nil {
		writeResponse(w, http.StatusInternalServerError, map[string]interface{}{errorJSONKey: err.Error()})
		return
	}
	writeResponse(w, http.StatusOK, map[string]interface{}{successJSONKey: resp})
}

func (g *gatewayService) putDiagnosisDataSourceHandler(w http.ResponseWriter, r *http.Request) {
	req := &diagnosis.PutDiagnosisDataSourceRequest{}
	bind(req, r)
	if err := req.Validate(); err != nil {
		writeResponse(w, http.StatusBadRequest, map[string]interface{}{errorJSONKey: err.Error()})
		return
	}
	resp, err := g.diagnosisClient.PutDiagnosisDataSource(r.Context(), req)
	if err != nil {
		writeResponse(w, http.StatusInternalServerError, map[string]interface{}{errorJSONKey: err.Error()})
		return
	}
	writeResponse(w, http.StatusOK, map[string]interface{}{successJSONKey: resp})
}

func (g *gatewayService) deleteDiagnosisDataSourceHandler(w http.ResponseWriter, r *http.Request) {
	req := &diagnosis.DeleteDiagnosisDataSourceRequest{}
	bind(req, r)
	if err := req.Validate(); err != nil {
		writeResponse(w, http.StatusBadRequest, map[string]interface{}{errorJSONKey: err.Error()})
		return
	}
	resp, err := g.diagnosisClient.DeleteDiagnosisDataSource(r.Context(), req)
	if err != nil {
		writeResponse(w, http.StatusInternalServerError, map[string]interface{}{errorJSONKey: err.Error()})
		return
	}
	writeResponse(w, http.StatusOK, map[string]interface{}{successJSONKey: resp})
}

func (g *gatewayService) listRelDiagnosisDataSourceHandler(w http.ResponseWriter, r *http.Request) {
	req := &diagnosis.ListRelDiagnosisDataSourceRequest{}
	bind(req, r)
	if err := req.Validate(); err != nil {
		writeResponse(w, http.StatusBadRequest, map[string]interface{}{errorJSONKey: err.Error()})
		return
	}
	resp, err := g.diagnosisClient.ListRelDiagnosisDataSource(r.Context(), req)
	if err != nil {
		writeResponse(w, http.StatusInternalServerError, map[string]interface{}{errorJSONKey: err.Error()})
		return
	}
	writeResponse(w, http.StatusOK, map[string]interface{}{successJSONKey: resp})
}

func (g *gatewayService) getRelDiagnosisDataSourceHandler(w http.ResponseWriter, r *http.Request) {
	req := &diagnosis.GetRelDiagnosisDataSourceRequest{}
	bind(req, r)
	if err := req.Validate(); err != nil {
		writeResponse(w, http.StatusBadRequest, map[string]interface{}{errorJSONKey: err.Error()})
		return
	}
	resp, err := g.diagnosisClient.GetRelDiagnosisDataSource(r.Context(), req)
	if err != nil {
		writeResponse(w, http.StatusInternalServerError, map[string]interface{}{errorJSONKey: err.Error()})
		return
	}
	writeResponse(w, http.StatusOK, map[string]interface{}{successJSONKey: resp})
}

func (g *gatewayService) putRelDiagnosisDataSourceHandler(w http.ResponseWriter, r *http.Request) {
	req := &diagnosis.PutRelDiagnosisDataSourceRequest{}
	bind(req, r)
	if err := req.Validate(); err != nil {
		writeResponse(w, http.StatusBadRequest, map[string]interface{}{errorJSONKey: err.Error()})
		return
	}
	resp, err := g.diagnosisClient.PutRelDiagnosisDataSource(r.Context(), req)
	if err != nil {
		writeResponse(w, http.StatusInternalServerError, map[string]interface{}{errorJSONKey: err.Error()})
		return
	}
	writeResponse(w, http.StatusOK, map[string]interface{}{successJSONKey: resp})
}

func (g *gatewayService) deleteRelDiagnosisDataSourceHandler(w http.ResponseWriter, r *http.Request) {
	req := &diagnosis.DeleteRelDiagnosisDataSourceRequest{}
	bind(req, r)
	if err := req.Validate(); err != nil {
		writeResponse(w, http.StatusBadRequest, map[string]interface{}{errorJSONKey: err.Error()})
		return
	}
	resp, err := g.diagnosisClient.DeleteRelDiagnosisDataSource(r.Context(), req)
	if err != nil {
		writeResponse(w, http.StatusInternalServerError, map[string]interface{}{errorJSONKey: err.Error()})
		return
	}
	writeResponse(w, http.StatusOK, map[string]interface{}{successJSONKey: resp})
}

func (g *gatewayService) startDiagnosisHandler(w http.ResponseWriter, r *http.Request) {
	req := &diagnosis.StartDiagnosisRequest{}
	bind(req, r)
	if err := req.Validate(); err != nil {
		writeResponse(w, http.StatusBadRequest, map[string]interface{}{errorJSONKey: err.Error()})
		return
	}
	resp, err := g.diagnosisClient.StartDiagnosis(r.Context(), req)
	if err != nil {
		writeResponse(w, http.StatusInternalServerError, map[string]interface{}{errorJSONKey: err.Error()})
		return
	}
	writeResponse(w, http.StatusOK, map[string]interface{}{successJSONKey: resp})
}
