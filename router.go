package main

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

const (
	contenTypeJSON = "application/json"
)

func newRouter(svc *gatewayService) *chi.Mux {
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Recoverer)
	r.Use(httpLogger)
	r.Use(middleware.StripSlashes)
	r.Use(svc.authn)

	r.Route("/api/v1", func(r chi.Router) {
		r.Get("/signin", signinHandler)

		r.Route("/finding", func(r chi.Router) {
			r.Use(svc.authzWithProject)
			r.Get("/list-finding", svc.listFindingHandler)
			r.Get("/get-finding", svc.getFindingHandler)
			r.Get("/list-finding-tag", svc.listFindingTagHandler)
			r.Get("/list-finding-tag-name", svc.listFindingTagNameHandler)
			r.Get("/list-resource", svc.listResourceHandler)
			r.Get("/get-resource", svc.getResourceHandler)
			r.Get("/list-resource-tag", svc.listResourceTagHandler)
			r.Get("/list-resource-tag-name", svc.listResourceTagNameHandler)
			r.Group(func(r chi.Router) {
				r.Use(middleware.AllowContentType(contenTypeJSON))
				r.Post("/put-finding", svc.putFindingHandler)
				r.Post("/delete-finding", svc.deleteFindingHandler)
				r.Post("/tag-finding", svc.tagFindingHandler)
				r.Post("/untag-finding", svc.untagFindingHandler)
				r.Post("/put-resource", svc.putResourceHandler)
				r.Post("/delete-resource", svc.deleteResourceHandler)
				r.Post("/tag-resource", svc.tagResourceHandler)
				r.Post("/untag-resource", svc.untagResourceHandler)
			})
		})

		r.Route("/iam", func(r chi.Router) {
			r.Get("/list-user", svc.listUserHandler)
			r.Get("/get-user", svc.getUserHandler)
			r.Group(func(r chi.Router) {
				r.Use(middleware.AllowContentType(contenTypeJSON))
				r.Post("/put-user", svc.putUserHandler)
			})
			r.Group(func(r chi.Router) {
				r.Use(svc.authzWithProject)
				r.Get("/list-role", svc.listRoleHandler)
				r.Get("/get-role", svc.getRoleHandler)
				r.Get("/list-policy", svc.listPolicyHandler)
				r.Get("/get-policy", svc.getPolicyHandler)
				r.Group(func(r chi.Router) {
					r.Use(middleware.AllowContentType(contenTypeJSON))
					r.Post("/put-role", svc.putRoleHandler)
					r.Post("/delete-role", svc.deleteRoleHandler)
					r.Post("/attach-role", svc.attachRoleHandler)
					r.Post("/detach-role", svc.detachRoleHandler)
					r.Post("/put-policy", svc.putPolicyHandler)
					r.Post("/delete-policy", svc.deletePolicyHandler)
					r.Post("/attach-policy", svc.attachPolicyHandler)
					r.Post("/detach-policy", svc.detachPolicyHandler)
				})
			})
		})

		r.Route("/project", func(r chi.Router) {
			r.Get("/list-project", svc.listProjectHandler)
			r.Group(func(r chi.Router) {
				r.Use(middleware.AllowContentType(contenTypeJSON))
				r.Post("/create-project", svc.createProjectHandler)
			})
			r.Group(func(r chi.Router) {
				r.Use(svc.authzWithProject)
				r.Use(middleware.AllowContentType(contenTypeJSON))
				r.Post("/update-project", svc.updateProjectHandler)
				r.Post("/delete-project", svc.deleteProjectHandler)
			})
		})

		r.Route("/alert", func(r chi.Router) {
			r.Use(svc.authzWithProject)
			r.Get("/list-alert", svc.listAlertHandler)
			r.Get("/list-history", svc.listAlertHistoryHandler)
			r.Get("/list-rel_alert_finding", svc.listRelAlertFindingHandler)
			r.Get("/list-condition", svc.listAlertConditionHandler)
			r.Get("/list-rule", svc.listAlertRuleHandler)
			r.Get("/list-condition_rule", svc.listAlertCondRuleHandler)
			r.Get("/list-notification", svc.listNotificationHandler)
			r.Get("/list-condition_notification", svc.listAlertCondNotificationHandler)
			r.Get("/get-alert", svc.getAlertHandler)
			r.Get("/get-history", svc.getAlertHistoryHandler)
			r.Get("/get-rel_alert_finding", svc.getRelAlertFindingHandler)
			r.Get("/get-condition", svc.getAlertConditionHandler)
			r.Get("/get-rule", svc.getAlertRuleHandler)
			r.Get("/get-condition_rule", svc.getAlertCondRuleHandler)
			r.Get("/get-notification", svc.getNotificationHandler)
			r.Get("/get-condition_notification", svc.getAlertCondNotificationHandler)
			r.Group(func(r chi.Router) {
				r.Use(middleware.AllowContentType(contenTypeJSON))
				r.Post("/put-alert", svc.putAlertHandler)
				r.Post("/put-history", svc.putAlertHistoryHandler)
				r.Post("/put-rel_alert_finding", svc.putRelAlertFindingHandler)
				r.Post("/put-condition", svc.putAlertConditionHandler)
				r.Post("/put-rule", svc.putAlertRuleHandler)
				r.Post("/put-condition_rule", svc.putAlertCondRuleHandler)
				r.Post("/put-notification", svc.putNotificationHandler)
				r.Post("/put-condition_notification", svc.putAlertCondNotificationHandler)
				r.Post("/delete-alert", svc.deleteAlertHandler)
				r.Post("/delete-history", svc.deleteAlertHistoryHandler)
				r.Post("/delete-rel_alert_finding", svc.deleteRelAlertFindingHandler)
				r.Post("/delete-condition", svc.deleteAlertConditionHandler)
				r.Post("/delete-rule", svc.deleteAlertRuleHandler)
				r.Post("/delete-condition_rule", svc.deleteAlertCondRuleHandler)
				r.Post("/delete-notification", svc.deleteNotificationHandler)
				r.Post("/delete-condition_notification", svc.deleteAlertCondNotificationHandler)
				r.Post("/analyze-alert", svc.analyzeAlertHandler)
			})
		})

		r.Route("/aws", func(r chi.Router) {
			r.Use(svc.authzWithProject)
			r.Get("/list-aws", svc.listAWSHandler)
			r.Get("/list-datasource", svc.listDataSourceHandler)
			r.Group(func(r chi.Router) {
				r.Use(middleware.AllowContentType(contenTypeJSON))
				r.Post("/put-aws", svc.putAWSHandler)
				r.Post("/delete-aws", svc.deleteAWSHandler)
				r.Post("/attach-datasource", svc.attachDataSourceHandler)
				r.Post("/detach-datasource", svc.detachDataSourceHandler)
				r.Post("/invoke-scan", svc.invokeScanHandler)
			})
		})

		r.Route("/osint", func(r chi.Router) {
			r.Use(svc.authzWithProject)
			r.Get("/list-osint", svc.listOSINTHandler)
			r.Get("/list-datasource", svc.listOSINTDataSourceHandler)
			r.Get("/list-rel-datasource", svc.listOSINTResultHandler)
			r.Get("/get-osint", svc.getOSINTHandler)
			r.Get("/get-datasource", svc.getOSINTDataSourceHandler)
			r.Get("/get-rel-datasource", svc.getOSINTResultHandler)
			r.Group(func(r chi.Router) {
				r.Use(middleware.AllowContentType(contenTypeJSON))
				r.Post("/put-osint", svc.putOSINTHandler)
				r.Post("/delete-osint", svc.deleteOSINTHandler)
				r.Post("/put-datasource", svc.putOSINTDataSourceHandler)
				r.Post("/delete-datasource", svc.deleteOSINTDataSourceHandler)
				r.Post("/put-rel-datasource", svc.putOSINTResultHandler)
				r.Post("/delete-rel-datasource", svc.deleteOSINTResultHandler)
				r.Post("/start-osint", svc.startOSINTHandler)
			})
		})

		r.Route("/diagnosis", func(r chi.Router) {
			r.Use(svc.authzWithProject)
			r.Get("/list-datasource", svc.listDiagnosisDataSourceHandler)
			r.Get("/list-jira-setting", svc.listJiraSettingHandler)
			r.Get("/get-datasource", svc.getDiagnosisDataSourceHandler)
			r.Get("/get-jira-setting", svc.getJiraSettingHandler)
			r.Group(func(r chi.Router) {
				r.Use(middleware.AllowContentType(contenTypeJSON))
				r.Post("/put-datasource", svc.putDiagnosisDataSourceHandler)
				r.Post("/delete-datasource", svc.deleteDiagnosisDataSourceHandler)
				r.Post("/put-jira-setting", svc.putJiraSettingHandler)
				r.Post("/delete-jira-setting", svc.deleteJiraSettingHandler)
				r.Post("/start-diagnosis", svc.startDiagnosisHandler)
			})
		})
	})

	r.Get("/healthz", func(w http.ResponseWriter, r *http.Request) { w.WriteHeader(http.StatusOK) })
	return r
}
