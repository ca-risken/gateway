package main

import (
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
