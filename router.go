package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	chitrace "gopkg.in/DataDog/dd-trace-go.v1/contrib/go-chi/chi.v5"
)

const (
	contenTypeJSON = "application/json"
	healthzPath    = "/healthz"
)

func newRouter(svc *gatewayService) *chi.Mux {
	r := chi.NewRouter()
	r.Use(commonHeader)
	r.Use(chitrace.Middleware(chitrace.WithIgnoreRequest(isTraceSkip)))
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(httpLogger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.StripSlashes)
	r.Use(svc.authn)
	r.Use(svc.authnToken)
	r.Use(svc.verifyCSRF)
	r.NotFound(notFoundHandler)

	r.Route("/api/v1", func(r chi.Router) {
		r.Route("/signin", func(r chi.Router) {
			r.Use(svc.UpdateUserFromIdp)
			r.Get("/", signinHandler)
		})

		r.Route("/finding", func(r chi.Router) {
			r.Use(svc.authzWithProject)
			r.Get("/list-finding", svc.listFindingFindingHandler)
			r.Get("/get-finding", svc.getFindingFindingHandler)
			r.Get("/list-finding-tag", svc.listFindingTagFindingHandler)
			r.Get("/list-finding-tag-name", svc.listFindingTagNameFindingHandler)
			r.Get("/list-resource", svc.listResourceFindingHandler)
			r.Get("/list-resource-tag", svc.listResourceTagFindingHandler)
			r.Get("/list-resource-tag-name", svc.listResourceTagNameFindingHandler)
			r.Get("/get-resource", svc.getResourceFindingHandler)
			r.Get("/get-pend-finding", svc.getPendFindingFindingHandler)
			r.Get("/list-finding-setting", svc.listFindingSettingFindingHandler)
			r.Get("/get-recommend", svc.getRecommendFindingHandler)
			r.Get("/get-ai-summary", svc.getAISummaryFindingHandler)
			r.Get("/get-ai-summary-stream", svc.getAISummaryStreamHandler)
			r.Group(func(r chi.Router) {
				r.Use(middleware.AllowContentType(contenTypeJSON))
				r.Post("/put-finding", svc.putFindingFindingHandler)
				r.Post("/delete-finding", svc.deleteFindingFindingHandler)
				r.Post("/tag-finding", svc.tagFindingFindingHandler)
				r.Post("/untag-finding", svc.untagFindingFindingHandler)
				r.Post("/put-resource", svc.putResourceFindingHandler)
				r.Post("/delete-resource", svc.deleteResourceFindingHandler)
				r.Post("/put-pend-finding", svc.putPendFindingFindingHandler)
				r.Post("/delete-pend-finding", svc.deletePendFindingFindingHandler)
				r.Post("/put-finding-setting", svc.putFindingSettingFindingHandler)
				r.Post("/delete-finding-setting", svc.deleteFindingSettingFindingHandler)
				r.Post("/put-recommend", svc.putRecommendFindingHandler)
			})
		})

		r.Route("/iam", func(r chi.Router) {
			r.Get("/list-user", svc.listUserIamHandler)
			r.Get("/get-user", svc.getUserIamHandler)
			r.Get("/is-admin", svc.isAdminIamHandler)
			r.Group(func(r chi.Router) {
				r.Use(middleware.AllowContentType(contenTypeJSON))
				r.Post("/put-user", svc.putUserHandler)
			})
			r.Group(func(r chi.Router) {
				r.Use(svc.authzWithProject)
				r.Get("/list-role", svc.listRoleIamHandler)
				r.Get("/get-role", svc.getRoleIamHandler)
				r.Get("/list-policy", svc.listPolicyIamHandler)
				r.Get("/get-policy", svc.getPolicyIamHandler)
				r.Get("/list-access-token", svc.listAccessTokenIamHandler)
				r.Get("/list-user-reserved", svc.listUserReservedIamHandler)
				r.Group(func(r chi.Router) {
					r.Use(middleware.AllowContentType(contenTypeJSON))
					r.Post("/put-role", svc.putRoleIamHandler)
					r.Post("/delete-role", svc.deleteRoleIamHandler)
					r.Post("/attach-role", svc.attachRoleIamHandler)
					r.Post("/detach-role", svc.detachRoleIamHandler)
					r.Post("/put-policy", svc.putPolicyIamHandler)
					r.Post("/delete-policy", svc.deletePolicyIamHandler)
					r.Post("/attach-policy", svc.attachPolicyIamHandler)
					r.Post("/detach-policy", svc.detachPolicyIamHandler)
					r.Post("/generate-access-token", svc.generateAccessTokenHandler)
					r.Post("/update-access-token", svc.updateAccessTokenHandler)
					r.Post("/delete-access-token", svc.deleteAccessTokenIamHandler)
					r.Post("/attach-access-token", svc.attachAccessTokenRoleIamHandler)
					r.Post("/detach-access-token", svc.detachAccessTokenRoleIamHandler)
					r.Post("/put-user-reserved", svc.putUserReservedIamHandler)
					r.Post("/delete-user-reserved", svc.deleteUserReservedIamHandler)
				})
			})
		})

		r.Route("/project", func(r chi.Router) {
			r.Get("/list-project", svc.listProjectProjectHandler)
			r.Group(func(r chi.Router) {
				r.Use(middleware.AllowContentType(contenTypeJSON))
				r.Post("/create-project", svc.createProjectHandler)
			})
			r.Group(func(r chi.Router) {
				r.Use(svc.authzWithProject)
				r.Use(middleware.AllowContentType(contenTypeJSON))
				r.Post("/update-project", svc.updateProjectProjectHandler)
				r.Post("/delete-project", svc.deleteProjectProjectHandler)
				r.Post("/tag-project", svc.tagProjectProjectHandler)
				r.Post("/untag-project", svc.untagProjectProjectHandler)
			})
		})

		r.Route("/alert", func(r chi.Router) {
			r.Use(svc.authzWithProject)
			r.Get("/list-alert", svc.listAlertAlertHandler)
			r.Get("/list-history", svc.listAlertHistoryAlertHandler)
			r.Get("/list-rel_alert_finding", svc.listRelAlertFindingAlertHandler)
			r.Get("/list-condition", svc.listAlertConditionAlertHandler)
			r.Get("/list-rule", svc.listAlertRuleAlertHandler)
			r.Get("/list-condition_rule", svc.listAlertCondRuleAlertHandler)
			r.Get("/list-notification", svc.listNotificationAlertHandler)
			r.Get("/list-condition_notification", svc.listAlertCondNotificationAlertHandler)
			r.Group(func(r chi.Router) {
				r.Use(middleware.AllowContentType(contenTypeJSON))
				r.Post("/put-alert", svc.putAlertAlertHandler)
				r.Post("/put-alert-first-viewed-at", svc.putAlertFirstViewedAtAlertHandler)
				r.Post("/put-condition", svc.putAlertConditionAlertHandler)
				r.Post("/put-rule", svc.putAlertRuleAlertHandler)
				r.Post("/put-condition_rule", svc.putAlertCondRuleAlertHandler)
				r.Post("/put-notification", svc.putNotificationAlertHandler)
				r.Post("/put-condition_notification", svc.putAlertCondNotificationAlertHandler)
				r.Post("/delete-condition", svc.deleteAlertConditionAlertHandler)
				r.Post("/delete-rule", svc.deleteAlertRuleAlertHandler)
				r.Post("/delete-condition_rule", svc.deleteAlertCondRuleAlertHandler)
				r.Post("/delete-notification", svc.deleteNotificationAlertHandler)
				r.Post("/delete-condition_notification", svc.deleteAlertCondNotificationAlertHandler)
				r.Post("/analyze-alert", svc.analyzeAlertAlertHandler)
				r.Post("/test-notification", svc.testNotificationAlertHandler)
			})
		})

		r.Route("/report", func(r chi.Router) {
			r.Group(func(r chi.Router) {
				r.Use(svc.authzWithProject)
				r.Get("/get-report", svc.getReportFindingReportHandler)
			})
			r.Group(func(r chi.Router) {
				// Admin API
				r.Use(svc.authzOnlyAdmin)
				r.Get("/get-report-all", svc.getReportFindingAllReportHandler)
			})
		})

		r.Route("/aws", func(r chi.Router) {
			r.Group(func(r chi.Router) {
				r.Use(svc.authzWithProject)
				r.Get("/list-aws", svc.listAWSAwsHandler)
				r.Get("/list-datasource", svc.listDataSourceAwsHandler)
				r.Group(func(r chi.Router) {
					r.Use(middleware.AllowContentType(contenTypeJSON))
					r.Post("/put-aws", svc.putAWSAwsHandler)
					r.Post("/delete-aws", svc.deleteAWSAwsHandler)
					r.Post("/invoke-scan", svc.invokeScanAwsHandler)
					r.Post("/attach-datasource", svc.attachDataSourceHandler)
					r.Post("/detach-datasource", svc.detachDataSourceAwsHandler)
				})
			})
		})

		r.Route("/osint", func(r chi.Router) {
			r.Group(func(r chi.Router) {
				r.Use(svc.authzWithProject)
				r.Get("/list-osint", svc.listOsintOsintHandler)
				r.Get("/list-datasource", svc.listOsintDataSourceOsintHandler)
				r.Get("/list-rel-datasource", svc.listRelOsintDataSourceOsintHandler)
				r.Get("/list-word", svc.listOsintDetectWordOsintHandler)
				r.Get("/get-osint", svc.getOsintOsintHandler)
				r.Group(func(r chi.Router) {
					r.Use(middleware.AllowContentType(contenTypeJSON))
					r.Post("/invoke-scan", svc.invokeScanOsintHandler)
					r.Post("/put-osint", svc.putOsintOsintHandler)
					r.Post("/delete-osint", svc.deleteOsintOsintHandler)
					r.Post("/put-rel-datasource", svc.putRelOsintDataSourceOsintHandler)
					r.Post("/delete-rel-datasource", svc.deleteRelOsintDataSourceOsintHandler)
					r.Post("/put-word", svc.putOsintDetectWordOsintHandler)
					r.Post("/delete-word", svc.deleteOsintDetectWordOsintHandler)
				})
			})
		})

		r.Route("/diagnosis", func(r chi.Router) {
			r.Group(func(r chi.Router) {
				r.Use(svc.authzWithProject)
				r.Get("/list-wpscan-setting", svc.listWpscanSettingDiagnosisHandler)
				r.Get("/list-portscan-setting", svc.listPortscanSettingDiagnosisHandler)
				r.Get("/list-portscan-target", svc.listPortscanTargetDiagnosisHandler)
				r.Get("/list-application-scan", svc.listApplicationScanDiagnosisHandler)
				r.Get("/get-application-scan-basic-setting", svc.getApplicationScanBasicSettingDiagnosisHandler)
				r.Get("/get-datasource", svc.getDiagnosisDataSourceDiagnosisHandler)
				r.Group(func(r chi.Router) {
					r.Use(middleware.AllowContentType(contenTypeJSON))
					r.Post("/invoke-scan", svc.invokeScanDiagnosisHandler)
					r.Post("/put-wpscan-setting", svc.putWpscanSettingDiagnosisHandler)
					r.Post("/delete-wpscan-setting", svc.deleteWpscanSettingDiagnosisHandler)
					r.Post("/put-application-scan", svc.putApplicationScanDiagnosisHandler)
					r.Post("/delete-application-scan", svc.deleteApplicationScanDiagnosisHandler)
					r.Post("/put-application-scan-basic-setting", svc.putApplicationScanBasicSettingDiagnosisHandler)
					r.Post("/delete-application-scan-basic-setting", svc.deleteApplicationScanBasicSettingDiagnosisHandler)
					r.Post("/put-portscan-setting", svc.putPortscanSettingDiagnosisHandler)
					r.Post("/put-portscan-target", svc.putPortscanTargetDiagnosisHandler)
					r.Post("/delete-portscan-setting", svc.deletePortscanSettingDiagnosisHandler)
					r.Post("/delete-portscan-target", svc.deletePortscanTargetDiagnosisHandler)
				})
			})
		})

		r.Route("/code", func(r chi.Router) {
			r.Group(func(r chi.Router) {
				// project any
				r.Get("/list-datasource", svc.listDataSourceCodeHandler)
			})
			r.Group(func(r chi.Router) {
				r.Use(svc.authzWithProject)
				r.Get("/list-github-setting", svc.listGitHubSettingCodeHandler)
				r.Get("/list-gitleaks-cache", svc.listGitleaksCacheCodeHandler)
				r.Group(func(r chi.Router) {
					r.Use(middleware.AllowContentType(contenTypeJSON))
					r.Post("/put-github-setting", svc.putGitHubSettingCodeHandler)
					r.Post("/delete-github-setting", svc.deleteGitHubSettingCodeHandler)
					r.Post("/put-gitleaks-setting", svc.putGitleaksSettingCodeHandler)
					r.Post("/delete-gitleaks-setting", svc.deleteGitleaksSettingCodeHandler)
					r.Post("/put-dependency-setting", svc.putDependencySettingCodeHandler)
					r.Post("/delete-dependency-setting", svc.deleteDependencySettingCodeHandler)
					r.Post("/invoke-scan-gitleaks", svc.invokeScanGitleaksCodeHandler)
					r.Post("/invoke-scan-dependency", svc.invokeScanDependencyCodeHandler)
				})
			})
		})

		r.Route("/google", func(r chi.Router) {
			r.Group(func(r chi.Router) {
				// project any
				r.Get("/list-google-datasource", svc.listGoogleDataSourceGoogleHandler)
			})
			r.Group(func(r chi.Router) {
				r.Use(svc.authzWithProject)
				r.Get("/list-gcp", svc.listGCPGoogleHandler)
				r.Get("/get-gcp-datasource", svc.getGCPDataSourceGoogleHandler)
				r.Group(func(r chi.Router) {
					r.Use(middleware.AllowContentType(contenTypeJSON))
					r.Post("/put-gcp", svc.putGCPGoogleHandler)
					r.Post("/delete-gcp", svc.deleteGCPGoogleHandler)
					r.Post("/attach-gcp-datasource", svc.attachGCPDataSourceGoogleHandler)
					r.Post("/detach-gcp-datasource", svc.detachGCPDataSourceGoogleHandler)
					r.Post("/invoke-scan-gcp", svc.invokeScanGCPGoogleHandler)
				})
			})
		})

		r.Route("/datasource", func(r chi.Router) {
			r.Group(func(r chi.Router) {
				r.Use(svc.authzWithProject)
				r.Get("/get-attack-flow-analysis", svc.analyzeAttackFlowDatasourceHandler)
			})
		})

		r.Route("/admin", func(r chi.Router) {
			// only admin
			r.Use(svc.authzOnlyAdmin)
			r.Group(func(r chi.Router) {
				r.Get("/list-admin-role", svc.listAdminRoleHandler)
				r.Get("/get-admin-role", svc.getAdminRoleHandler)
			})
			r.Group(func(r chi.Router) {
				r.Use(middleware.AllowContentType(contenTypeJSON))
				r.Post("/attach-admin-role", svc.attachAdminRoleHandler)
				r.Post("/detach-admin-role", svc.detachAdminRoleHandler)
			})
		})
	})
	r.Get(healthzPath, func(w http.ResponseWriter, r *http.Request) { w.WriteHeader(http.StatusOK) })
	return r
}

func isTraceSkip(r *http.Request) bool {
	if r == nil || r.URL == nil {
		return true
	}
	return r.URL.Path == healthzPath
}
