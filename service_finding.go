package main

import (
	"net/http"

	"github.com/ca-risken/core/proto/finding"
)

func (g *gatewayService) listFindingHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	// bind
	req := &finding.ListFindingRequest{}
	bind(req, r)
	// validate
	if err := req.Validate(); err != nil {
		writeResponse(ctx, w, http.StatusBadRequest, map[string]interface{}{errorJSONKey: err.Error()})
		return
	}
	// call backend service
	resp, err := g.findingClient.ListFinding(ctx, req)
	if err != nil {
		writeResponse(ctx, w, http.StatusInternalServerError, map[string]interface{}{errorJSONKey: err.Error()})
		return
	}
	writeResponse(ctx, w, http.StatusOK, map[string]interface{}{successJSONKey: resp})
}

func (g *gatewayService) getFindingHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	// bind
	req := &finding.GetFindingRequest{}
	bind(req, r)
	// validate
	if err := req.Validate(); err != nil {
		writeResponse(ctx, w, http.StatusBadRequest, map[string]interface{}{errorJSONKey: err.Error()})
		return
	}
	// call backend service
	resp, err := g.findingClient.GetFinding(ctx, req)
	if err != nil {
		writeResponse(ctx, w, http.StatusInternalServerError, map[string]interface{}{errorJSONKey: err.Error()})
		return
	}
	writeResponse(ctx, w, http.StatusOK, map[string]interface{}{successJSONKey: resp})
}

func (g *gatewayService) putFindingHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	// bind
	req := &finding.PutFindingRequest{}
	bind(req, r)
	// validate
	if err := req.Validate(); err != nil {
		writeResponse(ctx, w, http.StatusBadRequest, map[string]interface{}{errorJSONKey: err.Error()})
		return
	}
	// call backend service
	resp, err := g.findingClient.PutFinding(ctx, req)
	if err != nil {
		writeResponse(ctx, w, http.StatusInternalServerError, map[string]interface{}{errorJSONKey: err.Error()})
		return
	}
	writeResponse(ctx, w, http.StatusOK, map[string]interface{}{successJSONKey: resp})
}

func (g *gatewayService) deleteFindingHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	// bind
	req := &finding.DeleteFindingRequest{}
	bind(req, r)
	// validate
	if err := req.Validate(); err != nil {
		writeResponse(ctx, w, http.StatusBadRequest, map[string]interface{}{errorJSONKey: err.Error()})
		return
	}
	// call backend service
	resp, err := g.findingClient.DeleteFinding(ctx, req)
	if err != nil {
		writeResponse(ctx, w, http.StatusInternalServerError, map[string]interface{}{errorJSONKey: err.Error()})
		return
	}
	writeResponse(ctx, w, http.StatusOK, map[string]interface{}{successJSONKey: resp})
}

func (g *gatewayService) listFindingTagHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	// bind
	req := &finding.ListFindingTagRequest{}
	bind(req, r)
	// validate
	if err := req.Validate(); err != nil {
		writeResponse(ctx, w, http.StatusBadRequest, map[string]interface{}{errorJSONKey: err.Error()})
		return
	}
	// call backend service
	resp, err := g.findingClient.ListFindingTag(ctx, req)
	if err != nil {
		writeResponse(ctx, w, http.StatusInternalServerError, map[string]interface{}{errorJSONKey: err.Error()})
		return
	}
	writeResponse(ctx, w, http.StatusOK, map[string]interface{}{successJSONKey: resp})
}

func (g *gatewayService) listFindingTagNameHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	// bind
	req := &finding.ListFindingTagNameRequest{}
	bind(req, r)
	// validate
	if err := req.Validate(); err != nil {
		writeResponse(ctx, w, http.StatusBadRequest, map[string]interface{}{errorJSONKey: err.Error()})
		return
	}
	// call backend service
	resp, err := g.findingClient.ListFindingTagName(ctx, req)
	if err != nil {
		writeResponse(ctx, w, http.StatusInternalServerError, map[string]interface{}{errorJSONKey: err.Error()})
		return
	}
	writeResponse(ctx, w, http.StatusOK, map[string]interface{}{successJSONKey: resp})
}

func (g *gatewayService) tagFindingHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	// bind
	req := &finding.TagFindingRequest{}
	bind(req, r)
	// validate
	if err := req.Validate(); err != nil {
		writeResponse(ctx, w, http.StatusBadRequest, map[string]interface{}{errorJSONKey: err.Error()})
		return
	}
	// call backend service
	resp, err := g.findingClient.TagFinding(ctx, req)
	if err != nil {
		writeResponse(ctx, w, http.StatusInternalServerError, map[string]interface{}{errorJSONKey: err.Error()})
		return
	}
	writeResponse(ctx, w, http.StatusOK, map[string]interface{}{successJSONKey: resp})
}

func (g *gatewayService) untagFindingHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	// bind
	req := &finding.UntagFindingRequest{}
	bind(req, r)
	// validate
	if err := req.Validate(); err != nil {
		writeResponse(ctx, w, http.StatusBadRequest, map[string]interface{}{errorJSONKey: err.Error()})
		return
	}
	// call backend service
	resp, err := g.findingClient.UntagFinding(ctx, req)
	if err != nil {
		writeResponse(ctx, w, http.StatusInternalServerError, map[string]interface{}{errorJSONKey: err.Error()})
		return
	}
	writeResponse(ctx, w, http.StatusOK, map[string]interface{}{successJSONKey: resp})
}

func (g *gatewayService) listResourceHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	// bind
	req := &finding.ListResourceRequest{}
	bind(req, r)
	// validate
	if err := req.Validate(); err != nil {
		writeResponse(ctx, w, http.StatusBadRequest, map[string]interface{}{errorJSONKey: err.Error()})
		return
	}
	// call backend service
	resp, err := g.findingClient.ListResource(ctx, req)
	if err != nil {
		writeResponse(ctx, w, http.StatusInternalServerError, map[string]interface{}{errorJSONKey: err.Error()})
		return
	}
	writeResponse(ctx, w, http.StatusOK, map[string]interface{}{successJSONKey: resp})
}

func (g *gatewayService) getResourceHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	// bind
	req := &finding.GetResourceRequest{}
	bind(req, r)
	// validate
	if err := req.Validate(); err != nil {
		writeResponse(ctx, w, http.StatusBadRequest, map[string]interface{}{errorJSONKey: err.Error()})
		return
	}
	// call backend service
	resp, err := g.findingClient.GetResource(ctx, req)
	if err != nil {
		writeResponse(ctx, w, http.StatusInternalServerError, map[string]interface{}{errorJSONKey: err.Error()})
		return
	}
	writeResponse(ctx, w, http.StatusOK, map[string]interface{}{successJSONKey: resp})
}

func (g *gatewayService) putResourceHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	// bind
	req := &finding.PutResourceRequest{}
	bind(req, r)
	// validate
	if err := req.Validate(); err != nil {
		writeResponse(ctx, w, http.StatusBadRequest, map[string]interface{}{errorJSONKey: err.Error()})
		return
	}
	// call backend service
	resp, err := g.findingClient.PutResource(ctx, req)
	if err != nil {
		writeResponse(ctx, w, http.StatusInternalServerError, map[string]interface{}{errorJSONKey: err.Error()})
		return
	}
	writeResponse(ctx, w, http.StatusOK, map[string]interface{}{successJSONKey: resp})
}

func (g *gatewayService) deleteResourceHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	// bind
	req := &finding.DeleteResourceRequest{}
	bind(req, r)
	// validate
	if err := req.Validate(); err != nil {
		writeResponse(ctx, w, http.StatusBadRequest, map[string]interface{}{errorJSONKey: err.Error()})
		return
	}
	// call backend service
	resp, err := g.findingClient.DeleteResource(ctx, req)
	if err != nil {
		writeResponse(ctx, w, http.StatusInternalServerError, map[string]interface{}{errorJSONKey: err.Error()})
		return
	}
	writeResponse(ctx, w, http.StatusOK, map[string]interface{}{successJSONKey: resp})
}

func (g *gatewayService) listResourceTagHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	// bind
	req := &finding.ListResourceTagRequest{}
	bind(req, r)
	// validate
	if err := req.Validate(); err != nil {
		writeResponse(ctx, w, http.StatusBadRequest, map[string]interface{}{errorJSONKey: err.Error()})
		return
	}
	// call backend service
	resp, err := g.findingClient.ListResourceTag(ctx, req)
	if err != nil {
		writeResponse(ctx, w, http.StatusInternalServerError, map[string]interface{}{errorJSONKey: err.Error()})
		return
	}
	writeResponse(ctx, w, http.StatusOK, map[string]interface{}{successJSONKey: resp})
}

func (g *gatewayService) listResourceTagNameHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	req := &finding.ListResourceTagNameRequest{}
	bind(req, r)
	if err := req.Validate(); err != nil {
		writeResponse(ctx, w, http.StatusBadRequest, map[string]interface{}{errorJSONKey: err.Error()})
		return
	}
	resp, err := g.findingClient.ListResourceTagName(ctx, req)
	if err != nil {
		writeResponse(ctx, w, http.StatusInternalServerError, map[string]interface{}{errorJSONKey: err.Error()})
		return
	}
	writeResponse(ctx, w, http.StatusOK, map[string]interface{}{successJSONKey: resp})
}

func (g *gatewayService) tagResourceHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	// bind
	req := &finding.TagResourceRequest{}
	bind(req, r)
	// validate
	if err := req.Validate(); err != nil {
		writeResponse(ctx, w, http.StatusBadRequest, map[string]interface{}{errorJSONKey: err.Error()})
		return
	}
	// call backend service
	resp, err := g.findingClient.TagResource(ctx, req)
	if err != nil {
		writeResponse(ctx, w, http.StatusInternalServerError, map[string]interface{}{errorJSONKey: err.Error()})
		return
	}
	writeResponse(ctx, w, http.StatusOK, map[string]interface{}{successJSONKey: resp})
}

func (g *gatewayService) untagResourceHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	// bind
	req := &finding.UntagResourceRequest{}
	bind(req, r)
	// validate
	if err := req.Validate(); err != nil {
		writeResponse(ctx, w, http.StatusBadRequest, map[string]interface{}{errorJSONKey: err.Error()})
		return
	}
	// call backend service
	resp, err := g.findingClient.UntagResource(ctx, req)
	if err != nil {
		writeResponse(ctx, w, http.StatusInternalServerError, map[string]interface{}{errorJSONKey: err.Error()})
		return
	}
	writeResponse(ctx, w, http.StatusOK, map[string]interface{}{successJSONKey: resp})
}

func (g *gatewayService) getPendFindingHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	req := &finding.GetPendFindingRequest{}
	bind(req, r)
	if err := req.Validate(); err != nil {
		writeResponse(ctx, w, http.StatusBadRequest, map[string]interface{}{errorJSONKey: err.Error()})
		return
	}
	resp, err := g.findingClient.GetPendFinding(ctx, req)
	if err != nil {
		writeResponse(ctx, w, http.StatusInternalServerError, map[string]interface{}{errorJSONKey: err.Error()})
		return
	}
	writeResponse(ctx, w, http.StatusOK, map[string]interface{}{successJSONKey: resp})
}

func (g *gatewayService) putPendFindingHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	// bind
	req := &finding.PutPendFindingRequest{}
	bind(req, r)
	// validate
	if err := req.Validate(); err != nil {
		writeResponse(ctx, w, http.StatusBadRequest, map[string]interface{}{errorJSONKey: err.Error()})
		return
	}
	// call backend service
	resp, err := g.findingClient.PutPendFinding(ctx, req)
	if err != nil {
		writeResponse(ctx, w, http.StatusInternalServerError, map[string]interface{}{errorJSONKey: err.Error()})
		return
	}
	writeResponse(ctx, w, http.StatusOK, map[string]interface{}{successJSONKey: resp})
}

func (g *gatewayService) deletePendFindingHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	// bind
	req := &finding.DeletePendFindingRequest{}
	bind(req, r)
	// validate
	if err := req.Validate(); err != nil {
		writeResponse(ctx, w, http.StatusBadRequest, map[string]interface{}{errorJSONKey: err.Error()})
		return
	}
	// call backend service
	resp, err := g.findingClient.DeletePendFinding(ctx, req)
	if err != nil {
		writeResponse(ctx, w, http.StatusInternalServerError, map[string]interface{}{errorJSONKey: err.Error()})
		return
	}
	writeResponse(ctx, w, http.StatusOK, map[string]interface{}{successJSONKey: resp})
}

func (g *gatewayService) listFindingSettingHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	req := &finding.ListFindingSettingRequest{}
	bind(req, r)
	if err := req.Validate(); err != nil {
		writeResponse(ctx, w, http.StatusBadRequest, map[string]interface{}{errorJSONKey: err.Error()})
		return
	}
	resp, err := g.findingClient.ListFindingSetting(ctx, req)
	if err != nil {
		writeResponse(ctx, w, http.StatusInternalServerError, map[string]interface{}{errorJSONKey: err.Error()})
		return
	}
	writeResponse(ctx, w, http.StatusOK, map[string]interface{}{successJSONKey: resp})
}

func (g *gatewayService) getFindingSettingHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	req := &finding.GetFindingSettingRequest{}
	bind(req, r)
	if err := req.Validate(); err != nil {
		writeResponse(ctx, w, http.StatusBadRequest, map[string]interface{}{errorJSONKey: err.Error()})
		return
	}
	resp, err := g.findingClient.GetFindingSetting(ctx, req)
	if err != nil {
		writeResponse(ctx, w, http.StatusInternalServerError, map[string]interface{}{errorJSONKey: err.Error()})
		return
	}
	writeResponse(ctx, w, http.StatusOK, map[string]interface{}{successJSONKey: resp})
}

func (g *gatewayService) putFindingSettingHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	req := &finding.PutFindingSettingRequest{}
	bind(req, r)
	if err := req.Validate(); err != nil {
		writeResponse(ctx, w, http.StatusBadRequest, map[string]interface{}{errorJSONKey: err.Error()})
		return
	}
	resp, err := g.findingClient.PutFindingSetting(ctx, req)
	if err != nil {
		writeResponse(ctx, w, http.StatusInternalServerError, map[string]interface{}{errorJSONKey: err.Error()})
		return
	}
	writeResponse(ctx, w, http.StatusOK, map[string]interface{}{successJSONKey: resp})
}

func (g *gatewayService) deleteFindingSettingHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	req := &finding.DeleteFindingSettingRequest{}
	bind(req, r)
	if err := req.Validate(); err != nil {
		writeResponse(ctx, w, http.StatusBadRequest, map[string]interface{}{errorJSONKey: err.Error()})
		return
	}
	resp, err := g.findingClient.DeleteFindingSetting(ctx, req)
	if err != nil {
		writeResponse(ctx, w, http.StatusInternalServerError, map[string]interface{}{errorJSONKey: err.Error()})
		return
	}
	writeResponse(ctx, w, http.StatusOK, map[string]interface{}{successJSONKey: resp})
}

func (g *gatewayService) getRecommendHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	req := &finding.GetRecommendRequest{}
	bind(req, r)
	if err := req.Validate(); err != nil {
		writeResponse(ctx, w, http.StatusBadRequest, map[string]interface{}{errorJSONKey: err.Error()})
		return
	}
	resp, err := g.findingClient.GetRecommend(ctx, req)
	if err != nil {
		writeResponse(ctx, w, http.StatusInternalServerError, map[string]interface{}{errorJSONKey: err.Error()})
		return
	}
	writeResponse(ctx, w, http.StatusOK, map[string]interface{}{successJSONKey: resp})
}

func (g *gatewayService) putRecommendHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	req := &finding.PutRecommendRequest{}
	bind(req, r)
	if err := req.Validate(); err != nil {
		writeResponse(ctx, w, http.StatusBadRequest, map[string]interface{}{errorJSONKey: err.Error()})
		return
	}
	resp, err := g.findingClient.PutRecommend(ctx, req)
	if err != nil {
		writeResponse(ctx, w, http.StatusInternalServerError, map[string]interface{}{errorJSONKey: err.Error()})
		return
	}
	writeResponse(ctx, w, http.StatusOK, map[string]interface{}{successJSONKey: resp})
}
