package main

import (
	"net/http"

	"github.com/ca-risken/core/proto/org_alert"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

type updateOrgAlertCondNotificationCachePayload struct {
	OrganizationID   uint32  `json:"organization_id"`
	ProjectID        uint32  `json:"project_id"`
	AlertConditionID uint32  `json:"alert_condition_id"`
	NotificationID   uint32  `json:"notification_id"`
	CacheSecond      *uint32 `json:"cache_second"`
}

func (p *updateOrgAlertCondNotificationCachePayload) toRequest() *org_alert.UpdateOrgAlertCondNotificationCacheRequest {
	req := &org_alert.UpdateOrgAlertCondNotificationCacheRequest{
		OrganizationId:   p.OrganizationID,
		ProjectId:        p.ProjectID,
		AlertConditionId: p.AlertConditionID,
		NotificationId:   p.NotificationID,
	}
	if p.CacheSecond != nil {
		req.CacheSecond = wrapperspb.UInt32(*p.CacheSecond)
	}
	return req
}

type updateOrgAlertProjectNotificationEnabledPayload struct {
	OrganizationID uint32 `json:"organization_id"`
	ProjectID      uint32 `json:"project_id"`
	NotificationID uint32 `json:"notification_id"`
	Enabled        *bool  `json:"enabled"`
}

func (p *updateOrgAlertProjectNotificationEnabledPayload) toRequest() *org_alert.UpdateOrgAlertProjectNotificationEnabledRequest {
	req := &org_alert.UpdateOrgAlertProjectNotificationEnabledRequest{
		OrganizationId: p.OrganizationID,
		ProjectId:      p.ProjectID,
		NotificationId: p.NotificationID,
	}
	if p.Enabled != nil {
		req.Enabled = wrapperspb.Bool(*p.Enabled)
	}
	return req
}

type updateOrgAlertProjectNotificationCachePayload struct {
	OrganizationID uint32  `json:"organization_id"`
	ProjectID      uint32  `json:"project_id"`
	NotificationID uint32  `json:"notification_id"`
	CacheSecond    *uint32 `json:"cache_second"`
}

func (p *updateOrgAlertProjectNotificationCachePayload) toRequest() *org_alert.UpdateOrgAlertProjectNotificationCacheRequest {
	req := &org_alert.UpdateOrgAlertProjectNotificationCacheRequest{
		OrganizationId: p.OrganizationID,
		ProjectId:      p.ProjectID,
		NotificationId: p.NotificationID,
	}
	if p.CacheSecond != nil {
		req.CacheSecond = wrapperspb.UInt32(*p.CacheSecond)
	}
	return req
}

func (g *gatewayService) listOrgAlertCondNotificationOrg_alertHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	req := &org_alert.ListOrgAlertCondNotificationRequest{}
	if err := bind(req, r); err != nil {
		appLogger.Warnf(ctx, "Failed to bind request, req=%s, err=%+v", "ListOrgAlertCondNotificationRequest", err)
		writeResponse(ctx, w, http.StatusBadRequest, map[string]any{errorJSONKey: err.Error()})
		return
	}
	if err := req.Validate(); err != nil {
		writeResponse(ctx, w, http.StatusBadRequest, map[string]any{errorJSONKey: err.Error()})
		return
	}
	resp, err := g.org_alertClient.ListOrgAlertCondNotification(ctx, req)
	if err != nil {
		if handleErr := handleGRPCError(ctx, w, err); handleErr != nil {
			appLogger.Errorf(ctx, "HandleGRPCError: %+v", handleErr)
			writeResponse(ctx, w, http.StatusInternalServerError, map[string]any{errorJSONKey: "InternalServerError"})
		}
		return
	}
	writeResponse(ctx, w, http.StatusOK, map[string]any{successJSONKey: resp})
}

func (g *gatewayService) getOrgAlertCondNotificationOrg_alertHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	req := &org_alert.GetOrgAlertCondNotificationRequest{}
	if err := bind(req, r); err != nil {
		appLogger.Warnf(ctx, "Failed to bind request, req=%s, err=%+v", "GetOrgAlertCondNotificationRequest", err)
		writeResponse(ctx, w, http.StatusBadRequest, map[string]any{errorJSONKey: err.Error()})
		return
	}
	if err := req.Validate(); err != nil {
		writeResponse(ctx, w, http.StatusBadRequest, map[string]any{errorJSONKey: err.Error()})
		return
	}
	resp, err := g.org_alertClient.GetOrgAlertCondNotification(ctx, req)
	if err != nil {
		if handleErr := handleGRPCError(ctx, w, err); handleErr != nil {
			appLogger.Errorf(ctx, "HandleGRPCError: %+v", handleErr)
			writeResponse(ctx, w, http.StatusInternalServerError, map[string]any{errorJSONKey: "InternalServerError"})
		}
		return
	}
	writeResponse(ctx, w, http.StatusOK, map[string]any{successJSONKey: resp})
}

func (g *gatewayService) updateOrgAlertCondNotificationCacheOrg_alertHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	payload := &updateOrgAlertCondNotificationCachePayload{}
	if err := bind(payload, r); err != nil {
		appLogger.Warnf(ctx, "Failed to bind request, req=%s, err=%+v", "UpdateOrgAlertCondNotificationCacheRequest", err)
		writeResponse(ctx, w, http.StatusBadRequest, map[string]any{errorJSONKey: err.Error()})
		return
	}
	req := payload.toRequest()
	if err := req.Validate(); err != nil {
		writeResponse(ctx, w, http.StatusBadRequest, map[string]any{errorJSONKey: err.Error()})
		return
	}
	resp, err := g.org_alertClient.UpdateOrgAlertCondNotificationCache(ctx, req)
	if err != nil {
		if handleErr := handleGRPCError(ctx, w, err); handleErr != nil {
			appLogger.Errorf(ctx, "HandleGRPCError: %+v", handleErr)
			writeResponse(ctx, w, http.StatusInternalServerError, map[string]any{errorJSONKey: "InternalServerError"})
		}
		return
	}
	writeResponse(ctx, w, http.StatusOK, map[string]any{successJSONKey: resp})
}

func (g *gatewayService) updateOrgAlertProjectNotificationEnabledOrg_alertHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	payload := &updateOrgAlertProjectNotificationEnabledPayload{}
	if err := bind(payload, r); err != nil {
		appLogger.Warnf(ctx, "Failed to bind request, req=%s, err=%+v", "UpdateOrgAlertProjectNotificationEnabledRequest", err)
		writeResponse(ctx, w, http.StatusBadRequest, map[string]any{errorJSONKey: err.Error()})
		return
	}
	req := payload.toRequest()
	if err := req.Validate(); err != nil {
		writeResponse(ctx, w, http.StatusBadRequest, map[string]any{errorJSONKey: err.Error()})
		return
	}
	resp, err := g.org_alertClient.UpdateOrgAlertProjectNotificationEnabled(ctx, req)
	if err != nil {
		if handleErr := handleGRPCError(ctx, w, err); handleErr != nil {
			appLogger.Errorf(ctx, "HandleGRPCError: %+v", handleErr)
			writeResponse(ctx, w, http.StatusInternalServerError, map[string]any{errorJSONKey: "InternalServerError"})
		}
		return
	}
	writeResponse(ctx, w, http.StatusOK, map[string]any{successJSONKey: resp})
}

func (g *gatewayService) updateOrgAlertProjectNotificationCacheOrg_alertHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	payload := &updateOrgAlertProjectNotificationCachePayload{}
	if err := bind(payload, r); err != nil {
		appLogger.Warnf(ctx, "Failed to bind request, req=%s, err=%+v", "UpdateOrgAlertProjectNotificationCacheRequest", err)
		writeResponse(ctx, w, http.StatusBadRequest, map[string]any{errorJSONKey: err.Error()})
		return
	}
	req := payload.toRequest()
	if err := req.Validate(); err != nil {
		writeResponse(ctx, w, http.StatusBadRequest, map[string]any{errorJSONKey: err.Error()})
		return
	}
	resp, err := g.org_alertClient.UpdateOrgAlertProjectNotificationCache(ctx, req)
	if err != nil {
		if handleErr := handleGRPCError(ctx, w, err); handleErr != nil {
			appLogger.Errorf(ctx, "HandleGRPCError: %+v", handleErr)
			writeResponse(ctx, w, http.StatusInternalServerError, map[string]any{errorJSONKey: "InternalServerError"})
		}
		return
	}
	writeResponse(ctx, w, http.StatusOK, map[string]any{successJSONKey: resp})
}
