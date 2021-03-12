package main

import (
	"net/http"

	"github.com/CyberAgent/mimosa-core/proto/finding"
)

func (g *gatewayService) listFindingHandler(w http.ResponseWriter, r *http.Request) {
	// bind
	req := &finding.ListFindingRequest{}
	bind(req, r)
	// validate
	if err := req.Validate(); err != nil {
		writeResponse(w, http.StatusBadRequest, map[string]interface{}{errorJSONKey: err.Error()})
		return
	}
	// call backend service
	resp, err := g.findingClient.ListFinding(r.Context(), req)
	if err != nil {
		writeResponse(w, http.StatusInternalServerError, map[string]interface{}{errorJSONKey: err.Error()})
		return
	}
	writeResponse(w, http.StatusOK, map[string]interface{}{successJSONKey: resp})
}

func (g *gatewayService) getFindingHandler(w http.ResponseWriter, r *http.Request) {
	// bind
	req := &finding.GetFindingRequest{}
	bind(req, r)
	// validate
	if err := req.Validate(); err != nil {
		writeResponse(w, http.StatusBadRequest, map[string]interface{}{errorJSONKey: err.Error()})
		return
	}
	// call backend service
	resp, err := g.findingClient.GetFinding(r.Context(), req)
	if err != nil {
		writeResponse(w, http.StatusInternalServerError, map[string]interface{}{errorJSONKey: err.Error()})
		return
	}
	writeResponse(w, http.StatusOK, map[string]interface{}{successJSONKey: resp})
}

func (g *gatewayService) putFindingHandler(w http.ResponseWriter, r *http.Request) {
	// bind
	req := &finding.PutFindingRequest{}
	bind(req, r)
	// validate
	if err := req.Validate(); err != nil {
		writeResponse(w, http.StatusBadRequest, map[string]interface{}{errorJSONKey: err.Error()})
		return
	}
	// call backend service
	resp, err := g.findingClient.PutFinding(r.Context(), req)
	if err != nil {
		writeResponse(w, http.StatusInternalServerError, map[string]interface{}{errorJSONKey: err.Error()})
		return
	}
	writeResponse(w, http.StatusOK, map[string]interface{}{successJSONKey: resp})
}

func (g *gatewayService) deleteFindingHandler(w http.ResponseWriter, r *http.Request) {
	// bind
	req := &finding.DeleteFindingRequest{}
	bind(req, r)
	// validate
	if err := req.Validate(); err != nil {
		writeResponse(w, http.StatusBadRequest, map[string]interface{}{errorJSONKey: err.Error()})
		return
	}
	// call backend service
	resp, err := g.findingClient.DeleteFinding(r.Context(), req)
	if err != nil {
		writeResponse(w, http.StatusInternalServerError, map[string]interface{}{errorJSONKey: err.Error()})
		return
	}
	writeResponse(w, http.StatusOK, map[string]interface{}{successJSONKey: resp})
}

func (g *gatewayService) listFindingTagHandler(w http.ResponseWriter, r *http.Request) {
	// bind
	req := &finding.ListFindingTagRequest{}
	bind(req, r)
	// validate
	if err := req.Validate(); err != nil {
		writeResponse(w, http.StatusBadRequest, map[string]interface{}{errorJSONKey: err.Error()})
		return
	}
	// call backend service
	resp, err := g.findingClient.ListFindingTag(r.Context(), req)
	if err != nil {
		writeResponse(w, http.StatusInternalServerError, map[string]interface{}{errorJSONKey: err.Error()})
		return
	}
	writeResponse(w, http.StatusOK, map[string]interface{}{successJSONKey: resp})
}

func (g *gatewayService) listFindingTagNameHandler(w http.ResponseWriter, r *http.Request) {
	// bind
	req := &finding.ListFindingTagNameRequest{}
	bind(req, r)
	// validate
	if err := req.Validate(); err != nil {
		writeResponse(w, http.StatusBadRequest, map[string]interface{}{errorJSONKey: err.Error()})
		return
	}
	// call backend service
	resp, err := g.findingClient.ListFindingTagName(r.Context(), req)
	if err != nil {
		writeResponse(w, http.StatusInternalServerError, map[string]interface{}{errorJSONKey: err.Error()})
		return
	}
	writeResponse(w, http.StatusOK, map[string]interface{}{successJSONKey: resp})
}

func (g *gatewayService) tagFindingHandler(w http.ResponseWriter, r *http.Request) {
	// bind
	req := &finding.TagFindingRequest{}
	bind(req, r)
	// validate
	if err := req.Validate(); err != nil {
		writeResponse(w, http.StatusBadRequest, map[string]interface{}{errorJSONKey: err.Error()})
		return
	}
	// call backend service
	resp, err := g.findingClient.TagFinding(r.Context(), req)
	if err != nil {
		writeResponse(w, http.StatusInternalServerError, map[string]interface{}{errorJSONKey: err.Error()})
		return
	}
	writeResponse(w, http.StatusOK, map[string]interface{}{successJSONKey: resp})
}

func (g *gatewayService) untagFindingHandler(w http.ResponseWriter, r *http.Request) {
	// bind
	req := &finding.UntagFindingRequest{}
	bind(req, r)
	// validate
	if err := req.Validate(); err != nil {
		writeResponse(w, http.StatusBadRequest, map[string]interface{}{errorJSONKey: err.Error()})
		return
	}
	// call backend service
	resp, err := g.findingClient.UntagFinding(r.Context(), req)
	if err != nil {
		writeResponse(w, http.StatusInternalServerError, map[string]interface{}{errorJSONKey: err.Error()})
		return
	}
	writeResponse(w, http.StatusOK, map[string]interface{}{successJSONKey: resp})
}

func (g *gatewayService) listResourceHandler(w http.ResponseWriter, r *http.Request) {
	// bind
	req := &finding.ListResourceRequest{}
	bind(req, r)
	// validate
	if err := req.Validate(); err != nil {
		writeResponse(w, http.StatusBadRequest, map[string]interface{}{errorJSONKey: err.Error()})
		return
	}
	// call backend service
	resp, err := g.findingClient.ListResource(r.Context(), req)
	if err != nil {
		writeResponse(w, http.StatusInternalServerError, map[string]interface{}{errorJSONKey: err.Error()})
		return
	}
	writeResponse(w, http.StatusOK, map[string]interface{}{successJSONKey: resp})
}

func (g *gatewayService) getResourceHandler(w http.ResponseWriter, r *http.Request) {
	// bind
	req := &finding.GetResourceRequest{}
	bind(req, r)
	// validate
	if err := req.Validate(); err != nil {
		writeResponse(w, http.StatusBadRequest, map[string]interface{}{errorJSONKey: err.Error()})
		return
	}
	// call backend service
	resp, err := g.findingClient.GetResource(r.Context(), req)
	if err != nil {
		writeResponse(w, http.StatusInternalServerError, map[string]interface{}{errorJSONKey: err.Error()})
		return
	}
	writeResponse(w, http.StatusOK, map[string]interface{}{successJSONKey: resp})
}

func (g *gatewayService) putResourceHandler(w http.ResponseWriter, r *http.Request) {
	// bind
	req := &finding.PutResourceRequest{}
	bind(req, r)
	// validate
	if err := req.Validate(); err != nil {
		writeResponse(w, http.StatusBadRequest, map[string]interface{}{errorJSONKey: err.Error()})
		return
	}
	// call backend service
	resp, err := g.findingClient.PutResource(r.Context(), req)
	if err != nil {
		writeResponse(w, http.StatusInternalServerError, map[string]interface{}{errorJSONKey: err.Error()})
		return
	}
	writeResponse(w, http.StatusOK, map[string]interface{}{successJSONKey: resp})
}

func (g *gatewayService) deleteResourceHandler(w http.ResponseWriter, r *http.Request) {
	// bind
	req := &finding.DeleteResourceRequest{}
	bind(req, r)
	// validate
	if err := req.Validate(); err != nil {
		writeResponse(w, http.StatusBadRequest, map[string]interface{}{errorJSONKey: err.Error()})
		return
	}
	// call backend service
	resp, err := g.findingClient.DeleteResource(r.Context(), req)
	if err != nil {
		writeResponse(w, http.StatusInternalServerError, map[string]interface{}{errorJSONKey: err.Error()})
		return
	}
	writeResponse(w, http.StatusOK, map[string]interface{}{successJSONKey: resp})
}

func (g *gatewayService) listResourceTagHandler(w http.ResponseWriter, r *http.Request) {
	// bind
	req := &finding.ListResourceTagRequest{}
	bind(req, r)
	// validate
	if err := req.Validate(); err != nil {
		writeResponse(w, http.StatusBadRequest, map[string]interface{}{errorJSONKey: err.Error()})
		return
	}
	// call backend service
	resp, err := g.findingClient.ListResourceTag(r.Context(), req)
	if err != nil {
		writeResponse(w, http.StatusInternalServerError, map[string]interface{}{errorJSONKey: err.Error()})
		return
	}
	writeResponse(w, http.StatusOK, map[string]interface{}{successJSONKey: resp})
}

func (g *gatewayService) listResourceTagNameHandler(w http.ResponseWriter, r *http.Request) {
	// bind
	req := &finding.ListResourceTagNameRequest{}
	bind(req, r)
	// validate
	if err := req.Validate(); err != nil {
		writeResponse(w, http.StatusBadRequest, map[string]interface{}{errorJSONKey: err.Error()})
		return
	}
	// call backend service
	resp, err := g.findingClient.ListResourceTagName(r.Context(), req)
	if err != nil {
		writeResponse(w, http.StatusInternalServerError, map[string]interface{}{errorJSONKey: err.Error()})
		return
	}
	writeResponse(w, http.StatusOK, map[string]interface{}{successJSONKey: resp})
}

func (g *gatewayService) tagResourceHandler(w http.ResponseWriter, r *http.Request) {
	// bind
	req := &finding.TagResourceRequest{}
	bind(req, r)
	// validate
	if err := req.Validate(); err != nil {
		writeResponse(w, http.StatusBadRequest, map[string]interface{}{errorJSONKey: err.Error()})
		return
	}
	// call backend service
	resp, err := g.findingClient.TagResource(r.Context(), req)
	if err != nil {
		writeResponse(w, http.StatusInternalServerError, map[string]interface{}{errorJSONKey: err.Error()})
		return
	}
	writeResponse(w, http.StatusOK, map[string]interface{}{successJSONKey: resp})
}

func (g *gatewayService) untagResourceHandler(w http.ResponseWriter, r *http.Request) {
	// bind
	req := &finding.UntagResourceRequest{}
	bind(req, r)
	// validate
	if err := req.Validate(); err != nil {
		writeResponse(w, http.StatusBadRequest, map[string]interface{}{errorJSONKey: err.Error()})
		return
	}
	// call backend service
	resp, err := g.findingClient.UntagResource(r.Context(), req)
	if err != nil {
		writeResponse(w, http.StatusInternalServerError, map[string]interface{}{errorJSONKey: err.Error()})
		return
	}
	writeResponse(w, http.StatusOK, map[string]interface{}{successJSONKey: resp})
}

func (g *gatewayService) getPendFindingHandler(w http.ResponseWriter, r *http.Request) {
	req := &finding.GetPendFindingRequest{}
	bind(req, r)
	if err := req.Validate(); err != nil {
		writeResponse(w, http.StatusBadRequest, map[string]interface{}{errorJSONKey: err.Error()})
		return
	}
	resp, err := g.findingClient.GetPendFinding(r.Context(), req)
	if err != nil {
		writeResponse(w, http.StatusInternalServerError, map[string]interface{}{errorJSONKey: err.Error()})
		return
	}
	writeResponse(w, http.StatusOK, map[string]interface{}{successJSONKey: resp})
}

func (g *gatewayService) putPendFindingHandler(w http.ResponseWriter, r *http.Request) {
	// bind
	req := &finding.PutPendFindingRequest{}
	bind(req, r)
	// validate
	if err := req.Validate(); err != nil {
		writeResponse(w, http.StatusBadRequest, map[string]interface{}{errorJSONKey: err.Error()})
		return
	}
	// call backend service
	resp, err := g.findingClient.PutPendFinding(r.Context(), req)
	if err != nil {
		writeResponse(w, http.StatusInternalServerError, map[string]interface{}{errorJSONKey: err.Error()})
		return
	}
	writeResponse(w, http.StatusOK, map[string]interface{}{successJSONKey: resp})
}

func (g *gatewayService) deletePendFindingHandler(w http.ResponseWriter, r *http.Request) {
	// bind
	req := &finding.DeletePendFindingRequest{}
	bind(req, r)
	// validate
	if err := req.Validate(); err != nil {
		writeResponse(w, http.StatusBadRequest, map[string]interface{}{errorJSONKey: err.Error()})
		return
	}
	// call backend service
	resp, err := g.findingClient.DeletePendFinding(r.Context(), req)
	if err != nil {
		writeResponse(w, http.StatusInternalServerError, map[string]interface{}{errorJSONKey: err.Error()})
		return
	}
	writeResponse(w, http.StatusOK, map[string]interface{}{successJSONKey: resp})
}

func (g *gatewayService) listFindingSettingHandler(w http.ResponseWriter, r *http.Request) {
	req := &finding.ListFindingSettingRequest{}
	bind(req, r)
	if err := req.Validate(); err != nil {
		writeResponse(w, http.StatusBadRequest, map[string]interface{}{errorJSONKey: err.Error()})
		return
	}
	resp, err := g.findingClient.ListFindingSetting(r.Context(), req)
	if err != nil {
		writeResponse(w, http.StatusInternalServerError, map[string]interface{}{errorJSONKey: err.Error()})
		return
	}
	writeResponse(w, http.StatusOK, map[string]interface{}{successJSONKey: resp})
}

func (g *gatewayService) getFindingSettingHandler(w http.ResponseWriter, r *http.Request) {
	req := &finding.GetFindingSettingRequest{}
	bind(req, r)
	if err := req.Validate(); err != nil {
		writeResponse(w, http.StatusBadRequest, map[string]interface{}{errorJSONKey: err.Error()})
		return
	}
	resp, err := g.findingClient.GetFindingSetting(r.Context(), req)
	if err != nil {
		writeResponse(w, http.StatusInternalServerError, map[string]interface{}{errorJSONKey: err.Error()})
		return
	}
	writeResponse(w, http.StatusOK, map[string]interface{}{successJSONKey: resp})
}

func (g *gatewayService) putFindingSettingHandler(w http.ResponseWriter, r *http.Request) {
	req := &finding.PutFindingSettingRequest{}
	bind(req, r)
	if err := req.Validate(); err != nil {
		writeResponse(w, http.StatusBadRequest, map[string]interface{}{errorJSONKey: err.Error()})
		return
	}
	resp, err := g.findingClient.PutFindingSetting(r.Context(), req)
	if err != nil {
		writeResponse(w, http.StatusInternalServerError, map[string]interface{}{errorJSONKey: err.Error()})
		return
	}
	writeResponse(w, http.StatusOK, map[string]interface{}{successJSONKey: resp})
}

func (g *gatewayService) deleteFindingSettingHandler(w http.ResponseWriter, r *http.Request) {
	req := &finding.DeleteFindingSettingRequest{}
	bind(req, r)
	if err := req.Validate(); err != nil {
		writeResponse(w, http.StatusBadRequest, map[string]interface{}{errorJSONKey: err.Error()})
		return
	}
	resp, err := g.findingClient.DeleteFindingSetting(r.Context(), req)
	if err != nil {
		writeResponse(w, http.StatusInternalServerError, map[string]interface{}{errorJSONKey: err.Error()})
		return
	}
	writeResponse(w, http.StatusOK, map[string]interface{}{successJSONKey: resp})
}
