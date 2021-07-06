package main

import (
	"net/http"

	"github.com/aws/aws-xray-sdk-go/xray"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

const (
	contenTypeJSON = "application/json"
)

func newRouter(svc *gatewayService) *chi.Mux {
	r := chi.NewRouter()
	// TODO refactor
	r.Use(
		func(next http.Handler) http.Handler {
			return xray.Handler(xray.NewFixedSegmentNamer("gateway"), next)
		})
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(httpLogger)
	r.Use(middleware.Recoverer)
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
			r.Get("/get-pend-finding", svc.getPendFindingHandler)
			r.Get("/list-finding-setting", svc.listFindingSettingHandler)
			r.Get("/get-finding-setting", svc.getFindingSettingHandler)
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
				r.Post("/put-pend-finding", svc.putPendFindingHandler)
				r.Post("/delete-pend-finding", svc.deletePendFindingHandler)
				r.Post("/put-finding-setting", svc.putFindingSettingHandler)
				r.Post("/delete-finding-setting", svc.deleteFindingSettingHandler)
			})
		})

		r.Route("/iam", func(r chi.Router) {
			r.Get("/list-user", svc.listUserHandler)
			r.Get("/get-user", svc.getUserHandler)
			r.Get("/is-admin", svc.isAdminHandler)
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
				r.Post("/tag-project", svc.tagProjectHandler)
				r.Post("/untag-project", svc.untagProjectHandler)
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
				r.Post("/test-notification", svc.testNotificationHandler)
			})
		})

		r.Route("/report", func(r chi.Router) {
			r.Group(func(r chi.Router) {
				r.Use(svc.authzWithProject)
				r.Get("/get-report", svc.getReportFindingHandler)
			})
			r.Group(func(r chi.Router) {
				// Admin API
				r.Use(svc.authzOnlyAdmin)
				r.Get("/get-report-all", svc.getReportFindingAllHandler)
			})
		})

		r.Route("/aws", func(r chi.Router) {
			r.Group(func(r chi.Router) {
				r.Use(svc.authzWithProject)
				r.Get("/list-aws", svc.listAWSHandler)
				r.Get("/list-datasource", svc.listDataSourceHandler)
				r.Get("/describe-arn", svc.describeARNHandler)
				r.Get("/list-cloudtrail", svc.listCloudTrailHandler)
				r.Get("/list-config-history", svc.listConfigHistoryHandler)
				r.Group(func(r chi.Router) {
					r.Use(middleware.AllowContentType(contenTypeJSON))
					r.Post("/put-aws", svc.putAWSHandler)
					r.Post("/delete-aws", svc.deleteAWSHandler)
					r.Post("/invoke-scan", svc.invokeScanHandler)
					r.Post("/attach-datasource", svc.attachDataSourceHandler)
					r.Post("/detach-datasource", svc.detachDataSourceHandler)
				})
			})
		})

		r.Route("/osint", func(r chi.Router) {
			r.Group(func(r chi.Router) {
				r.Use(svc.authzWithProject)
				r.Get("/list-osint", svc.listOsintHandler)
				r.Get("/list-datasource", svc.listOsintDataSourceHandler)
				r.Get("/list-rel-datasource", svc.listRelOsintDataSourceHandler)
				r.Get("/list-word", svc.listOsintDetectWordHandler)
				r.Get("/get-osint", svc.getOsintHandler)
				r.Get("/get-datasource", svc.getOsintDataSourceHandler)
				r.Get("/get-rel-datasource", svc.getRelOsintDataSourceHandler)
				r.Get("/get-word", svc.getOsintDetectWordHandler)
				r.Group(func(r chi.Router) {
					r.Use(middleware.AllowContentType(contenTypeJSON))
					r.Post("/invoke-scan", svc.invokeOsintScanHandler)
					r.Post("/put-osint", svc.putOsintHandler)
					r.Post("/delete-osint", svc.deleteOsintHandler)
					r.Post("/put-rel-datasource", svc.putRelOsintDataSourceHandler)
					r.Post("/delete-rel-datasource", svc.deleteRelOsintDataSourceHandler)
					r.Post("/put-word", svc.putOsintDetectWordHandler)
					r.Post("/delete-word", svc.deleteOsintDetectWordHandler)
				})
			})
			r.Group(func(r chi.Router) {
				// Admin API
				r.Use(svc.authzOnlyAdmin)
				r.Use(middleware.AllowContentType(contenTypeJSON))
				r.Post("/put-datasource", svc.putOsintDataSourceHandler)
				r.Post("/delete-datasource", svc.deleteOsintDataSourceHandler)
			})
		})

		r.Route("/diagnosis", func(r chi.Router) {
			r.Group(func(r chi.Router) {
				r.Use(svc.authzWithProject)
				r.Get("/list-datasource", svc.listDiagnosisDataSourceHandler)
				r.Get("/list-jira-setting", svc.listJiraSettingHandler)
				r.Get("/list-wpscan-setting", svc.listWpscanSettingHandler)
				r.Get("/get-datasource", svc.getDiagnosisDataSourceHandler)
				r.Get("/get-jira-setting", svc.getJiraSettingHandler)
				r.Get("/get-wpscan-setting", svc.getWpscanSettingHandler)
				r.Group(func(r chi.Router) {
					r.Use(middleware.AllowContentType(contenTypeJSON))
					r.Post("/invoke-scan", svc.invokeDiagnosisScanHandler)
					r.Post("/put-wpscan-setting", svc.putWpscanSettingHandler)
					r.Post("/delete-wpscan-setting", svc.deleteWpscanSettingHandler)
				})
			})
			r.Group(func(r chi.Router) {
				// Admin API
				r.Use(svc.authzOnlyAdmin)
				r.Use(middleware.AllowContentType(contenTypeJSON))
				r.Post("/put-datasource", svc.putDiagnosisDataSourceHandler)
				r.Post("/delete-datasource", svc.deleteDiagnosisDataSourceHandler)
				r.Post("/put-jira-setting", svc.putJiraSettingHandler)
				r.Post("/delete-jira-setting", svc.deleteJiraSettingHandler)
			})
		})

		r.Route("/code", func(r chi.Router) {
			r.Group(func(r chi.Router) {
				// project any
				r.Get("/list-datasource", svc.listCodeDataSourceHandler)
			})
			r.Group(func(r chi.Router) {
				r.Use(svc.authzWithProject)
				r.Get("/list-gitleaks", svc.listGitleaksHandler)
				r.Group(func(r chi.Router) {
					r.Use(middleware.AllowContentType(contenTypeJSON))
					r.Post("/put-gitleaks", svc.putGitleaksHandler)
					r.Post("/delete-gitleaks", svc.deleteGitleaksHandler)
					r.Post("/invoke-scan-gitleaks", svc.invokeScanGitleaksHandler)
				})
			})
		})

		r.Route("/google", func(r chi.Router) {
			r.Group(func(r chi.Router) {
				// project any
				r.Get("/list-google-datasource", svc.listGoogleDataSourceHandler)
			})
			r.Group(func(r chi.Router) {
				r.Use(svc.authzWithProject)
				r.Get("/list-gcp", svc.listGCPHandler)
				r.Get("/list-gcp-datasource", svc.listGCPDataSourceHandler)
				r.Get("/get-gcp", svc.getGCPHandler)
				r.Get("/get-gcp-datasource", svc.getGCPDataSourceHandler)
				r.Group(func(r chi.Router) {
					r.Use(middleware.AllowContentType(contenTypeJSON))
					r.Post("/put-gcp", svc.putGCPHandler)
					r.Post("/delete-gcp", svc.deleteGCPHandler)
					r.Post("/attach-gcp-datasource", svc.attachGCPDataSourceHandler)
					r.Post("/detach-gcp-datasource", svc.detachGCPDataSourceHandler)
					r.Post("/invoke-scan-gcp", svc.invokeScanGCPHandler)
				})
			})
		})

	})
	r.Get("/healthz", func(w http.ResponseWriter, r *http.Request) { w.WriteHeader(http.StatusOK) })
	return r
}
