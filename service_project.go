package main

import (
	"errors"
	"net/http"

	"github.com/ca-risken/core/proto/project"
)

func (g *gatewayService) listProjectHandler(w http.ResponseWriter, r *http.Request) {
	req := &project.ListProjectRequest{}
	bind(req, r)
	if err := req.Validate(); err != nil {
		writeResponse(w, http.StatusBadRequest, map[string]interface{}{errorJSONKey: err.Error()})
		return
	}
	resp, err := g.projectClient.ListProject(r.Context(), req)
	if err != nil {
		writeResponse(w, http.StatusInternalServerError, map[string]interface{}{errorJSONKey: err.Error()})
		return
	}
	writeResponse(w, http.StatusOK, map[string]interface{}{successJSONKey: resp})
}

func (g *gatewayService) createProjectHandler(w http.ResponseWriter, r *http.Request) {
	user, err := getRequestUser(r)
	if err != nil {
		writeResponse(w, http.StatusUnauthorized, map[string]interface{}{errorJSONKey: errors.New("InvalidUser")})
	}
	req := &project.CreateProjectRequest{}
	req.UserId = user.userID // force update by own userID
	bind(req, r)
	if err := req.Validate(); err != nil {
		writeResponse(w, http.StatusBadRequest, map[string]interface{}{errorJSONKey: err.Error()})
		return
	}
	resp, err := g.projectClient.CreateProject(r.Context(), req)
	if err != nil {
		writeResponse(w, http.StatusInternalServerError, map[string]interface{}{errorJSONKey: err.Error()})
		return
	}
	writeResponse(w, http.StatusOK, map[string]interface{}{successJSONKey: resp})
}

func (g *gatewayService) updateProjectHandler(w http.ResponseWriter, r *http.Request) {
	req := &project.UpdateProjectRequest{}
	bind(req, r)
	if err := req.Validate(); err != nil {
		writeResponse(w, http.StatusBadRequest, map[string]interface{}{errorJSONKey: err.Error()})
		return
	}
	resp, err := g.projectClient.UpdateProject(r.Context(), req)
	if err != nil {
		writeResponse(w, http.StatusInternalServerError, map[string]interface{}{errorJSONKey: err.Error()})
		return
	}
	writeResponse(w, http.StatusOK, map[string]interface{}{successJSONKey: resp})
}

func (g *gatewayService) deleteProjectHandler(w http.ResponseWriter, r *http.Request) {
	req := &project.DeleteProjectRequest{}
	bind(req, r)
	if err := req.Validate(); err != nil {
		writeResponse(w, http.StatusBadRequest, map[string]interface{}{errorJSONKey: err.Error()})
		return
	}
	resp, err := g.projectClient.DeleteProject(r.Context(), req)
	if err != nil {
		writeResponse(w, http.StatusInternalServerError, map[string]interface{}{errorJSONKey: err.Error()})
		return
	}
	writeResponse(w, http.StatusOK, map[string]interface{}{successJSONKey: resp})
}

func (g *gatewayService) tagProjectHandler(w http.ResponseWriter, r *http.Request) {
	req := &project.TagProjectRequest{}
	bind(req, r)
	if err := req.Validate(); err != nil {
		writeResponse(w, http.StatusBadRequest, map[string]interface{}{errorJSONKey: err.Error()})
		return
	}
	resp, err := g.projectClient.TagProject(r.Context(), req)
	if err != nil {
		writeResponse(w, http.StatusInternalServerError, map[string]interface{}{errorJSONKey: err.Error()})
		return
	}
	writeResponse(w, http.StatusOK, map[string]interface{}{successJSONKey: resp})
}

func (g *gatewayService) untagProjectHandler(w http.ResponseWriter, r *http.Request) {
	req := &project.UntagProjectRequest{}
	bind(req, r)
	if err := req.Validate(); err != nil {
		writeResponse(w, http.StatusBadRequest, map[string]interface{}{errorJSONKey: err.Error()})
		return
	}
	resp, err := g.projectClient.UntagProject(r.Context(), req)
	if err != nil {
		writeResponse(w, http.StatusInternalServerError, map[string]interface{}{errorJSONKey: err.Error()})
		return
	}
	writeResponse(w, http.StatusOK, map[string]interface{}{successJSONKey: resp})
}
