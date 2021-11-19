# githug.com/ca-risken/gateway

MIMOSA API document by go-chi.

## Routes

<details>
<summary>`/api/v1/admin/attach-admin-role`</summary>

- [IgnoreHealthCheckTracingMiddleware.func1]()
- [AnnotateEnvTracingMiddleware.func1]()
- [RequestID]()
- [RealIP]()
- [RequestLogger.func1]()
- [Recoverer]()
- [StripSlashes]()
- [main.(*gatewayService).authn-fm]()
- [main.(*gatewayService).authnToken-fm]()
- **/api/v1**
	- **/admin**
		- [main.(*gatewayService).authzOnlyAdmin-fm]()
		- **/attach-admin-role**
			- _POST_
				- [AllowContentType.func1]()
				- [main.(*gatewayService).attachAdminRoleHandler-fm]()

</details>
<details>
<summary>`/api/v1/admin/detach-admin-role`</summary>

- [IgnoreHealthCheckTracingMiddleware.func1]()
- [AnnotateEnvTracingMiddleware.func1]()
- [RequestID]()
- [RealIP]()
- [RequestLogger.func1]()
- [Recoverer]()
- [StripSlashes]()
- [main.(*gatewayService).authn-fm]()
- [main.(*gatewayService).authnToken-fm]()
- **/api/v1**
	- **/admin**
		- [main.(*gatewayService).authzOnlyAdmin-fm]()
		- **/detach-admin-role**
			- _POST_
				- [AllowContentType.func1]()
				- [main.(*gatewayService).detachAdminRoleHandler-fm]()

</details>
<details>
<summary>`/api/v1/admin/get-admin-role`</summary>

- [IgnoreHealthCheckTracingMiddleware.func1]()
- [AnnotateEnvTracingMiddleware.func1]()
- [RequestID]()
- [RealIP]()
- [RequestLogger.func1]()
- [Recoverer]()
- [StripSlashes]()
- [main.(*gatewayService).authn-fm]()
- [main.(*gatewayService).authnToken-fm]()
- **/api/v1**
	- **/admin**
		- [main.(*gatewayService).authzOnlyAdmin-fm]()
		- **/get-admin-role**
			- _GET_
				- [main.(*gatewayService).getAdminRoleHandler-fm]()

</details>
<details>
<summary>`/api/v1/admin/list-admin-role`</summary>

- [IgnoreHealthCheckTracingMiddleware.func1]()
- [AnnotateEnvTracingMiddleware.func1]()
- [RequestID]()
- [RealIP]()
- [RequestLogger.func1]()
- [Recoverer]()
- [StripSlashes]()
- [main.(*gatewayService).authn-fm]()
- [main.(*gatewayService).authnToken-fm]()
- **/api/v1**
	- **/admin**
		- [main.(*gatewayService).authzOnlyAdmin-fm]()
		- **/list-admin-role**
			- _GET_
				- [main.(*gatewayService).listAdminRoleHandler-fm]()

</details>
<details>
<summary>`/api/v1/alert/analyze-alert`</summary>

- [IgnoreHealthCheckTracingMiddleware.func1]()
- [AnnotateEnvTracingMiddleware.func1]()
- [RequestID]()
- [RealIP]()
- [RequestLogger.func1]()
- [Recoverer]()
- [StripSlashes]()
- [main.(*gatewayService).authn-fm]()
- [main.(*gatewayService).authnToken-fm]()
- **/api/v1**
	- **/alert**
		- [main.(*gatewayService).authzWithProject-fm]()
		- **/analyze-alert**
			- _POST_
				- [AllowContentType.func1]()
				- [main.(*gatewayService).analyzeAlertHandler-fm]()

</details>
<details>
<summary>`/api/v1/alert/delete-condition`</summary>

- [IgnoreHealthCheckTracingMiddleware.func1]()
- [AnnotateEnvTracingMiddleware.func1]()
- [RequestID]()
- [RealIP]()
- [RequestLogger.func1]()
- [Recoverer]()
- [StripSlashes]()
- [main.(*gatewayService).authn-fm]()
- [main.(*gatewayService).authnToken-fm]()
- **/api/v1**
	- **/alert**
		- [main.(*gatewayService).authzWithProject-fm]()
		- **/delete-condition**
			- _POST_
				- [AllowContentType.func1]()
				- [main.(*gatewayService).deleteAlertConditionHandler-fm]()

</details>
<details>
<summary>`/api/v1/alert/delete-condition_notification`</summary>

- [IgnoreHealthCheckTracingMiddleware.func1]()
- [AnnotateEnvTracingMiddleware.func1]()
- [RequestID]()
- [RealIP]()
- [RequestLogger.func1]()
- [Recoverer]()
- [StripSlashes]()
- [main.(*gatewayService).authn-fm]()
- [main.(*gatewayService).authnToken-fm]()
- **/api/v1**
	- **/alert**
		- [main.(*gatewayService).authzWithProject-fm]()
		- **/delete-condition_notification**
			- _POST_
				- [AllowContentType.func1]()
				- [main.(*gatewayService).deleteAlertCondNotificationHandler-fm]()

</details>
<details>
<summary>`/api/v1/alert/delete-condition_rule`</summary>

- [IgnoreHealthCheckTracingMiddleware.func1]()
- [AnnotateEnvTracingMiddleware.func1]()
- [RequestID]()
- [RealIP]()
- [RequestLogger.func1]()
- [Recoverer]()
- [StripSlashes]()
- [main.(*gatewayService).authn-fm]()
- [main.(*gatewayService).authnToken-fm]()
- **/api/v1**
	- **/alert**
		- [main.(*gatewayService).authzWithProject-fm]()
		- **/delete-condition_rule**
			- _POST_
				- [AllowContentType.func1]()
				- [main.(*gatewayService).deleteAlertCondRuleHandler-fm]()

</details>
<details>
<summary>`/api/v1/alert/delete-notification`</summary>

- [IgnoreHealthCheckTracingMiddleware.func1]()
- [AnnotateEnvTracingMiddleware.func1]()
- [RequestID]()
- [RealIP]()
- [RequestLogger.func1]()
- [Recoverer]()
- [StripSlashes]()
- [main.(*gatewayService).authn-fm]()
- [main.(*gatewayService).authnToken-fm]()
- **/api/v1**
	- **/alert**
		- [main.(*gatewayService).authzWithProject-fm]()
		- **/delete-notification**
			- _POST_
				- [AllowContentType.func1]()
				- [main.(*gatewayService).deleteNotificationHandler-fm]()

</details>
<details>
<summary>`/api/v1/alert/delete-rule`</summary>

- [IgnoreHealthCheckTracingMiddleware.func1]()
- [AnnotateEnvTracingMiddleware.func1]()
- [RequestID]()
- [RealIP]()
- [RequestLogger.func1]()
- [Recoverer]()
- [StripSlashes]()
- [main.(*gatewayService).authn-fm]()
- [main.(*gatewayService).authnToken-fm]()
- **/api/v1**
	- **/alert**
		- [main.(*gatewayService).authzWithProject-fm]()
		- **/delete-rule**
			- _POST_
				- [AllowContentType.func1]()
				- [main.(*gatewayService).deleteAlertRuleHandler-fm]()

</details>
<details>
<summary>`/api/v1/alert/list-alert`</summary>

- [IgnoreHealthCheckTracingMiddleware.func1]()
- [AnnotateEnvTracingMiddleware.func1]()
- [RequestID]()
- [RealIP]()
- [RequestLogger.func1]()
- [Recoverer]()
- [StripSlashes]()
- [main.(*gatewayService).authn-fm]()
- [main.(*gatewayService).authnToken-fm]()
- **/api/v1**
	- **/alert**
		- [main.(*gatewayService).authzWithProject-fm]()
		- **/list-alert**
			- _GET_
				- [main.(*gatewayService).listAlertHandler-fm]()

</details>
<details>
<summary>`/api/v1/alert/list-condition`</summary>

- [IgnoreHealthCheckTracingMiddleware.func1]()
- [AnnotateEnvTracingMiddleware.func1]()
- [RequestID]()
- [RealIP]()
- [RequestLogger.func1]()
- [Recoverer]()
- [StripSlashes]()
- [main.(*gatewayService).authn-fm]()
- [main.(*gatewayService).authnToken-fm]()
- **/api/v1**
	- **/alert**
		- [main.(*gatewayService).authzWithProject-fm]()
		- **/list-condition**
			- _GET_
				- [main.(*gatewayService).listAlertConditionHandler-fm]()

</details>
<details>
<summary>`/api/v1/alert/list-condition_notification`</summary>

- [IgnoreHealthCheckTracingMiddleware.func1]()
- [AnnotateEnvTracingMiddleware.func1]()
- [RequestID]()
- [RealIP]()
- [RequestLogger.func1]()
- [Recoverer]()
- [StripSlashes]()
- [main.(*gatewayService).authn-fm]()
- [main.(*gatewayService).authnToken-fm]()
- **/api/v1**
	- **/alert**
		- [main.(*gatewayService).authzWithProject-fm]()
		- **/list-condition_notification**
			- _GET_
				- [main.(*gatewayService).listAlertCondNotificationHandler-fm]()

</details>
<details>
<summary>`/api/v1/alert/list-condition_rule`</summary>

- [IgnoreHealthCheckTracingMiddleware.func1]()
- [AnnotateEnvTracingMiddleware.func1]()
- [RequestID]()
- [RealIP]()
- [RequestLogger.func1]()
- [Recoverer]()
- [StripSlashes]()
- [main.(*gatewayService).authn-fm]()
- [main.(*gatewayService).authnToken-fm]()
- **/api/v1**
	- **/alert**
		- [main.(*gatewayService).authzWithProject-fm]()
		- **/list-condition_rule**
			- _GET_
				- [main.(*gatewayService).listAlertCondRuleHandler-fm]()

</details>
<details>
<summary>`/api/v1/alert/list-history`</summary>

- [IgnoreHealthCheckTracingMiddleware.func1]()
- [AnnotateEnvTracingMiddleware.func1]()
- [RequestID]()
- [RealIP]()
- [RequestLogger.func1]()
- [Recoverer]()
- [StripSlashes]()
- [main.(*gatewayService).authn-fm]()
- [main.(*gatewayService).authnToken-fm]()
- **/api/v1**
	- **/alert**
		- [main.(*gatewayService).authzWithProject-fm]()
		- **/list-history**
			- _GET_
				- [main.(*gatewayService).listAlertHistoryHandler-fm]()

</details>
<details>
<summary>`/api/v1/alert/list-notification`</summary>

- [IgnoreHealthCheckTracingMiddleware.func1]()
- [AnnotateEnvTracingMiddleware.func1]()
- [RequestID]()
- [RealIP]()
- [RequestLogger.func1]()
- [Recoverer]()
- [StripSlashes]()
- [main.(*gatewayService).authn-fm]()
- [main.(*gatewayService).authnToken-fm]()
- **/api/v1**
	- **/alert**
		- [main.(*gatewayService).authzWithProject-fm]()
		- **/list-notification**
			- _GET_
				- [main.(*gatewayService).listNotificationHandler-fm]()

</details>
<details>
<summary>`/api/v1/alert/list-rel_alert_finding`</summary>

- [IgnoreHealthCheckTracingMiddleware.func1]()
- [AnnotateEnvTracingMiddleware.func1]()
- [RequestID]()
- [RealIP]()
- [RequestLogger.func1]()
- [Recoverer]()
- [StripSlashes]()
- [main.(*gatewayService).authn-fm]()
- [main.(*gatewayService).authnToken-fm]()
- **/api/v1**
	- **/alert**
		- [main.(*gatewayService).authzWithProject-fm]()
		- **/list-rel_alert_finding**
			- _GET_
				- [main.(*gatewayService).listRelAlertFindingHandler-fm]()

</details>
<details>
<summary>`/api/v1/alert/list-rule`</summary>

- [IgnoreHealthCheckTracingMiddleware.func1]()
- [AnnotateEnvTracingMiddleware.func1]()
- [RequestID]()
- [RealIP]()
- [RequestLogger.func1]()
- [Recoverer]()
- [StripSlashes]()
- [main.(*gatewayService).authn-fm]()
- [main.(*gatewayService).authnToken-fm]()
- **/api/v1**
	- **/alert**
		- [main.(*gatewayService).authzWithProject-fm]()
		- **/list-rule**
			- _GET_
				- [main.(*gatewayService).listAlertRuleHandler-fm]()

</details>
<details>
<summary>`/api/v1/alert/put-alert`</summary>

- [IgnoreHealthCheckTracingMiddleware.func1]()
- [AnnotateEnvTracingMiddleware.func1]()
- [RequestID]()
- [RealIP]()
- [RequestLogger.func1]()
- [Recoverer]()
- [StripSlashes]()
- [main.(*gatewayService).authn-fm]()
- [main.(*gatewayService).authnToken-fm]()
- **/api/v1**
	- **/alert**
		- [main.(*gatewayService).authzWithProject-fm]()
		- **/put-alert**
			- _POST_
				- [AllowContentType.func1]()
				- [main.(*gatewayService).putAlertHandler-fm]()

</details>
<details>
<summary>`/api/v1/alert/put-condition`</summary>

- [IgnoreHealthCheckTracingMiddleware.func1]()
- [AnnotateEnvTracingMiddleware.func1]()
- [RequestID]()
- [RealIP]()
- [RequestLogger.func1]()
- [Recoverer]()
- [StripSlashes]()
- [main.(*gatewayService).authn-fm]()
- [main.(*gatewayService).authnToken-fm]()
- **/api/v1**
	- **/alert**
		- [main.(*gatewayService).authzWithProject-fm]()
		- **/put-condition**
			- _POST_
				- [AllowContentType.func1]()
				- [main.(*gatewayService).putAlertConditionHandler-fm]()

</details>
<details>
<summary>`/api/v1/alert/put-condition_notification`</summary>

- [IgnoreHealthCheckTracingMiddleware.func1]()
- [AnnotateEnvTracingMiddleware.func1]()
- [RequestID]()
- [RealIP]()
- [RequestLogger.func1]()
- [Recoverer]()
- [StripSlashes]()
- [main.(*gatewayService).authn-fm]()
- [main.(*gatewayService).authnToken-fm]()
- **/api/v1**
	- **/alert**
		- [main.(*gatewayService).authzWithProject-fm]()
		- **/put-condition_notification**
			- _POST_
				- [AllowContentType.func1]()
				- [main.(*gatewayService).putAlertCondNotificationHandler-fm]()

</details>
<details>
<summary>`/api/v1/alert/put-condition_rule`</summary>

- [IgnoreHealthCheckTracingMiddleware.func1]()
- [AnnotateEnvTracingMiddleware.func1]()
- [RequestID]()
- [RealIP]()
- [RequestLogger.func1]()
- [Recoverer]()
- [StripSlashes]()
- [main.(*gatewayService).authn-fm]()
- [main.(*gatewayService).authnToken-fm]()
- **/api/v1**
	- **/alert**
		- [main.(*gatewayService).authzWithProject-fm]()
		- **/put-condition_rule**
			- _POST_
				- [AllowContentType.func1]()
				- [main.(*gatewayService).putAlertCondRuleHandler-fm]()

</details>
<details>
<summary>`/api/v1/alert/put-notification`</summary>

- [IgnoreHealthCheckTracingMiddleware.func1]()
- [AnnotateEnvTracingMiddleware.func1]()
- [RequestID]()
- [RealIP]()
- [RequestLogger.func1]()
- [Recoverer]()
- [StripSlashes]()
- [main.(*gatewayService).authn-fm]()
- [main.(*gatewayService).authnToken-fm]()
- **/api/v1**
	- **/alert**
		- [main.(*gatewayService).authzWithProject-fm]()
		- **/put-notification**
			- _POST_
				- [AllowContentType.func1]()
				- [main.(*gatewayService).putNotificationHandler-fm]()

</details>
<details>
<summary>`/api/v1/alert/put-rule`</summary>

- [IgnoreHealthCheckTracingMiddleware.func1]()
- [AnnotateEnvTracingMiddleware.func1]()
- [RequestID]()
- [RealIP]()
- [RequestLogger.func1]()
- [Recoverer]()
- [StripSlashes]()
- [main.(*gatewayService).authn-fm]()
- [main.(*gatewayService).authnToken-fm]()
- **/api/v1**
	- **/alert**
		- [main.(*gatewayService).authzWithProject-fm]()
		- **/put-rule**
			- _POST_
				- [AllowContentType.func1]()
				- [main.(*gatewayService).putAlertRuleHandler-fm]()

</details>
<details>
<summary>`/api/v1/alert/test-notification`</summary>

- [IgnoreHealthCheckTracingMiddleware.func1]()
- [AnnotateEnvTracingMiddleware.func1]()
- [RequestID]()
- [RealIP]()
- [RequestLogger.func1]()
- [Recoverer]()
- [StripSlashes]()
- [main.(*gatewayService).authn-fm]()
- [main.(*gatewayService).authnToken-fm]()
- **/api/v1**
	- **/alert**
		- [main.(*gatewayService).authzWithProject-fm]()
		- **/test-notification**
			- _POST_
				- [AllowContentType.func1]()
				- [main.(*gatewayService).testNotificationHandler-fm]()

</details>
<details>
<summary>`/api/v1/aws/attach-datasource`</summary>

- [IgnoreHealthCheckTracingMiddleware.func1]()
- [AnnotateEnvTracingMiddleware.func1]()
- [RequestID]()
- [RealIP]()
- [RequestLogger.func1]()
- [Recoverer]()
- [StripSlashes]()
- [main.(*gatewayService).authn-fm]()
- [main.(*gatewayService).authnToken-fm]()
- **/api/v1**
	- **/aws**
		- **/attach-datasource**
			- _POST_
				- [main.(*gatewayService).authzWithProject-fm]()
				- [AllowContentType.func1]()
				- [main.(*gatewayService).attachDataSourceHandler-fm]()

</details>
<details>
<summary>`/api/v1/aws/delete-aws`</summary>

- [IgnoreHealthCheckTracingMiddleware.func1]()
- [AnnotateEnvTracingMiddleware.func1]()
- [RequestID]()
- [RealIP]()
- [RequestLogger.func1]()
- [Recoverer]()
- [StripSlashes]()
- [main.(*gatewayService).authn-fm]()
- [main.(*gatewayService).authnToken-fm]()
- **/api/v1**
	- **/aws**
		- **/delete-aws**
			- _POST_
				- [main.(*gatewayService).authzWithProject-fm]()
				- [AllowContentType.func1]()
				- [main.(*gatewayService).deleteAWSHandler-fm]()

</details>
<details>
<summary>`/api/v1/aws/describe-arn`</summary>

- [IgnoreHealthCheckTracingMiddleware.func1]()
- [AnnotateEnvTracingMiddleware.func1]()
- [RequestID]()
- [RealIP]()
- [RequestLogger.func1]()
- [Recoverer]()
- [StripSlashes]()
- [main.(*gatewayService).authn-fm]()
- [main.(*gatewayService).authnToken-fm]()
- **/api/v1**
	- **/aws**
		- **/describe-arn**
			- _GET_
				- [main.(*gatewayService).authzWithProject-fm]()
				- [main.(*gatewayService).describeARNHandler-fm]()

</details>
<details>
<summary>`/api/v1/aws/detach-datasource`</summary>

- [IgnoreHealthCheckTracingMiddleware.func1]()
- [AnnotateEnvTracingMiddleware.func1]()
- [RequestID]()
- [RealIP]()
- [RequestLogger.func1]()
- [Recoverer]()
- [StripSlashes]()
- [main.(*gatewayService).authn-fm]()
- [main.(*gatewayService).authnToken-fm]()
- **/api/v1**
	- **/aws**
		- **/detach-datasource**
			- _POST_
				- [main.(*gatewayService).authzWithProject-fm]()
				- [AllowContentType.func1]()
				- [main.(*gatewayService).detachDataSourceHandler-fm]()

</details>
<details>
<summary>`/api/v1/aws/invoke-scan`</summary>

- [IgnoreHealthCheckTracingMiddleware.func1]()
- [AnnotateEnvTracingMiddleware.func1]()
- [RequestID]()
- [RealIP]()
- [RequestLogger.func1]()
- [Recoverer]()
- [StripSlashes]()
- [main.(*gatewayService).authn-fm]()
- [main.(*gatewayService).authnToken-fm]()
- **/api/v1**
	- **/aws**
		- **/invoke-scan**
			- _POST_
				- [main.(*gatewayService).authzWithProject-fm]()
				- [AllowContentType.func1]()
				- [main.(*gatewayService).invokeScanHandler-fm]()

</details>
<details>
<summary>`/api/v1/aws/list-aws`</summary>

- [IgnoreHealthCheckTracingMiddleware.func1]()
- [AnnotateEnvTracingMiddleware.func1]()
- [RequestID]()
- [RealIP]()
- [RequestLogger.func1]()
- [Recoverer]()
- [StripSlashes]()
- [main.(*gatewayService).authn-fm]()
- [main.(*gatewayService).authnToken-fm]()
- **/api/v1**
	- **/aws**
		- **/list-aws**
			- _GET_
				- [main.(*gatewayService).authzWithProject-fm]()
				- [main.(*gatewayService).listAWSHandler-fm]()

</details>
<details>
<summary>`/api/v1/aws/list-cloudtrail`</summary>

- [IgnoreHealthCheckTracingMiddleware.func1]()
- [AnnotateEnvTracingMiddleware.func1]()
- [RequestID]()
- [RealIP]()
- [RequestLogger.func1]()
- [Recoverer]()
- [StripSlashes]()
- [main.(*gatewayService).authn-fm]()
- [main.(*gatewayService).authnToken-fm]()
- **/api/v1**
	- **/aws**
		- **/list-cloudtrail**
			- _GET_
				- [main.(*gatewayService).authzWithProject-fm]()
				- [main.(*gatewayService).listCloudTrailHandler-fm]()

</details>
<details>
<summary>`/api/v1/aws/list-config-history`</summary>

- [IgnoreHealthCheckTracingMiddleware.func1]()
- [AnnotateEnvTracingMiddleware.func1]()
- [RequestID]()
- [RealIP]()
- [RequestLogger.func1]()
- [Recoverer]()
- [StripSlashes]()
- [main.(*gatewayService).authn-fm]()
- [main.(*gatewayService).authnToken-fm]()
- **/api/v1**
	- **/aws**
		- **/list-config-history**
			- _GET_
				- [main.(*gatewayService).authzWithProject-fm]()
				- [main.(*gatewayService).listConfigHistoryHandler-fm]()

</details>
<details>
<summary>`/api/v1/aws/list-datasource`</summary>

- [IgnoreHealthCheckTracingMiddleware.func1]()
- [AnnotateEnvTracingMiddleware.func1]()
- [RequestID]()
- [RealIP]()
- [RequestLogger.func1]()
- [Recoverer]()
- [StripSlashes]()
- [main.(*gatewayService).authn-fm]()
- [main.(*gatewayService).authnToken-fm]()
- **/api/v1**
	- **/aws**
		- **/list-datasource**
			- _GET_
				- [main.(*gatewayService).authzWithProject-fm]()
				- [main.(*gatewayService).listDataSourceHandler-fm]()

</details>
<details>
<summary>`/api/v1/aws/put-aws`</summary>

- [IgnoreHealthCheckTracingMiddleware.func1]()
- [AnnotateEnvTracingMiddleware.func1]()
- [RequestID]()
- [RealIP]()
- [RequestLogger.func1]()
- [Recoverer]()
- [StripSlashes]()
- [main.(*gatewayService).authn-fm]()
- [main.(*gatewayService).authnToken-fm]()
- **/api/v1**
	- **/aws**
		- **/put-aws**
			- _POST_
				- [main.(*gatewayService).authzWithProject-fm]()
				- [AllowContentType.func1]()
				- [main.(*gatewayService).putAWSHandler-fm]()

</details>
<details>
<summary>`/api/v1/code/delete-gitleaks`</summary>

- [IgnoreHealthCheckTracingMiddleware.func1]()
- [AnnotateEnvTracingMiddleware.func1]()
- [RequestID]()
- [RealIP]()
- [RequestLogger.func1]()
- [Recoverer]()
- [StripSlashes]()
- [main.(*gatewayService).authn-fm]()
- [main.(*gatewayService).authnToken-fm]()
- **/api/v1**
	- **/code**
		- **/delete-gitleaks**
			- _POST_
				- [main.(*gatewayService).authzWithProject-fm]()
				- [AllowContentType.func1]()
				- [main.(*gatewayService).deleteGitleaksHandler-fm]()

</details>
<details>
<summary>`/api/v1/code/invoke-scan-gitleaks`</summary>

- [IgnoreHealthCheckTracingMiddleware.func1]()
- [AnnotateEnvTracingMiddleware.func1]()
- [RequestID]()
- [RealIP]()
- [RequestLogger.func1]()
- [Recoverer]()
- [StripSlashes]()
- [main.(*gatewayService).authn-fm]()
- [main.(*gatewayService).authnToken-fm]()
- **/api/v1**
	- **/code**
		- **/invoke-scan-gitleaks**
			- _POST_
				- [main.(*gatewayService).authzWithProject-fm]()
				- [AllowContentType.func1]()
				- [main.(*gatewayService).invokeScanGitleaksHandler-fm]()

</details>
<details>
<summary>`/api/v1/code/list-datasource`</summary>

- [IgnoreHealthCheckTracingMiddleware.func1]()
- [AnnotateEnvTracingMiddleware.func1]()
- [RequestID]()
- [RealIP]()
- [RequestLogger.func1]()
- [Recoverer]()
- [StripSlashes]()
- [main.(*gatewayService).authn-fm]()
- [main.(*gatewayService).authnToken-fm]()
- **/api/v1**
	- **/code**
		- **/list-datasource**
			- _GET_
				- [main.(*gatewayService).listCodeDataSourceHandler-fm]()

</details>
<details>
<summary>`/api/v1/code/list-gitleaks`</summary>

- [IgnoreHealthCheckTracingMiddleware.func1]()
- [AnnotateEnvTracingMiddleware.func1]()
- [RequestID]()
- [RealIP]()
- [RequestLogger.func1]()
- [Recoverer]()
- [StripSlashes]()
- [main.(*gatewayService).authn-fm]()
- [main.(*gatewayService).authnToken-fm]()
- **/api/v1**
	- **/code**
		- **/list-gitleaks**
			- _GET_
				- [main.(*gatewayService).authzWithProject-fm]()
				- [main.(*gatewayService).listGitleaksHandler-fm]()

</details>
<details>
<summary>`/api/v1/code/put-gitleaks`</summary>

- [IgnoreHealthCheckTracingMiddleware.func1]()
- [AnnotateEnvTracingMiddleware.func1]()
- [RequestID]()
- [RealIP]()
- [RequestLogger.func1]()
- [Recoverer]()
- [StripSlashes]()
- [main.(*gatewayService).authn-fm]()
- [main.(*gatewayService).authnToken-fm]()
- **/api/v1**
	- **/code**
		- **/put-gitleaks**
			- _POST_
				- [main.(*gatewayService).authzWithProject-fm]()
				- [AllowContentType.func1]()
				- [main.(*gatewayService).putGitleaksHandler-fm]()

</details>
<details>
<summary>`/api/v1/diagnosis/delete-application-scan`</summary>

- [IgnoreHealthCheckTracingMiddleware.func1]()
- [AnnotateEnvTracingMiddleware.func1]()
- [RequestID]()
- [RealIP]()
- [RequestLogger.func1]()
- [Recoverer]()
- [StripSlashes]()
- [main.(*gatewayService).authn-fm]()
- [main.(*gatewayService).authnToken-fm]()
- **/api/v1**
	- **/diagnosis**
		- **/delete-application-scan**
			- _POST_
				- [main.(*gatewayService).authzWithProject-fm]()
				- [AllowContentType.func1]()
				- [main.(*gatewayService).deleteApplicationScanHandler-fm]()

</details>
<details>
<summary>`/api/v1/diagnosis/delete-application-scan-basic-setting`</summary>

- [IgnoreHealthCheckTracingMiddleware.func1]()
- [AnnotateEnvTracingMiddleware.func1]()
- [RequestID]()
- [RealIP]()
- [RequestLogger.func1]()
- [Recoverer]()
- [StripSlashes]()
- [main.(*gatewayService).authn-fm]()
- [main.(*gatewayService).authnToken-fm]()
- **/api/v1**
	- **/diagnosis**
		- **/delete-application-scan-basic-setting**
			- _POST_
				- [main.(*gatewayService).authzWithProject-fm]()
				- [AllowContentType.func1]()
				- [main.(*gatewayService).deleteApplicationScanBasicSettingHandler-fm]()

</details>
<details>
<summary>`/api/v1/diagnosis/delete-jira-setting`</summary>

- [IgnoreHealthCheckTracingMiddleware.func1]()
- [AnnotateEnvTracingMiddleware.func1]()
- [RequestID]()
- [RealIP]()
- [RequestLogger.func1]()
- [Recoverer]()
- [StripSlashes]()
- [main.(*gatewayService).authn-fm]()
- [main.(*gatewayService).authnToken-fm]()
- **/api/v1**
	- **/diagnosis**
		- **/delete-jira-setting**
			- _POST_
				- [main.(*gatewayService).authzOnlyAdmin-fm]()
				- [AllowContentType.func1]()
				- [main.(*gatewayService).deleteJiraSettingHandler-fm]()

</details>
<details>
<summary>`/api/v1/diagnosis/delete-portscan-setting`</summary>

- [IgnoreHealthCheckTracingMiddleware.func1]()
- [AnnotateEnvTracingMiddleware.func1]()
- [RequestID]()
- [RealIP]()
- [RequestLogger.func1]()
- [Recoverer]()
- [StripSlashes]()
- [main.(*gatewayService).authn-fm]()
- [main.(*gatewayService).authnToken-fm]()
- **/api/v1**
	- **/diagnosis**
		- **/delete-portscan-setting**
			- _POST_
				- [main.(*gatewayService).authzWithProject-fm]()
				- [AllowContentType.func1]()
				- [main.(*gatewayService).deletePortscanSettingHandler-fm]()

</details>
<details>
<summary>`/api/v1/diagnosis/delete-portscan-target`</summary>

- [IgnoreHealthCheckTracingMiddleware.func1]()
- [AnnotateEnvTracingMiddleware.func1]()
- [RequestID]()
- [RealIP]()
- [RequestLogger.func1]()
- [Recoverer]()
- [StripSlashes]()
- [main.(*gatewayService).authn-fm]()
- [main.(*gatewayService).authnToken-fm]()
- **/api/v1**
	- **/diagnosis**
		- **/delete-portscan-target**
			- _POST_
				- [main.(*gatewayService).authzWithProject-fm]()
				- [AllowContentType.func1]()
				- [main.(*gatewayService).deletePortscanTargetHandler-fm]()

</details>
<details>
<summary>`/api/v1/diagnosis/delete-wpscan-setting`</summary>

- [IgnoreHealthCheckTracingMiddleware.func1]()
- [AnnotateEnvTracingMiddleware.func1]()
- [RequestID]()
- [RealIP]()
- [RequestLogger.func1]()
- [Recoverer]()
- [StripSlashes]()
- [main.(*gatewayService).authn-fm]()
- [main.(*gatewayService).authnToken-fm]()
- **/api/v1**
	- **/diagnosis**
		- **/delete-wpscan-setting**
			- _POST_
				- [main.(*gatewayService).authzWithProject-fm]()
				- [AllowContentType.func1]()
				- [main.(*gatewayService).deleteWpscanSettingHandler-fm]()

</details>
<details>
<summary>`/api/v1/diagnosis/get-application-scan-basic-setting`</summary>

- [IgnoreHealthCheckTracingMiddleware.func1]()
- [AnnotateEnvTracingMiddleware.func1]()
- [RequestID]()
- [RealIP]()
- [RequestLogger.func1]()
- [Recoverer]()
- [StripSlashes]()
- [main.(*gatewayService).authn-fm]()
- [main.(*gatewayService).authnToken-fm]()
- **/api/v1**
	- **/diagnosis**
		- **/get-application-scan-basic-setting**
			- _GET_
				- [main.(*gatewayService).authzWithProject-fm]()
				- [main.(*gatewayService).getApplicationScanBasicSettingHandler-fm]()

</details>
<details>
<summary>`/api/v1/diagnosis/get-datasource`</summary>

- [IgnoreHealthCheckTracingMiddleware.func1]()
- [AnnotateEnvTracingMiddleware.func1]()
- [RequestID]()
- [RealIP]()
- [RequestLogger.func1]()
- [Recoverer]()
- [StripSlashes]()
- [main.(*gatewayService).authn-fm]()
- [main.(*gatewayService).authnToken-fm]()
- **/api/v1**
	- **/diagnosis**
		- **/get-datasource**
			- _GET_
				- [main.(*gatewayService).authzWithProject-fm]()
				- [main.(*gatewayService).getDiagnosisDataSourceHandler-fm]()

</details>
<details>
<summary>`/api/v1/diagnosis/invoke-scan`</summary>

- [IgnoreHealthCheckTracingMiddleware.func1]()
- [AnnotateEnvTracingMiddleware.func1]()
- [RequestID]()
- [RealIP]()
- [RequestLogger.func1]()
- [Recoverer]()
- [StripSlashes]()
- [main.(*gatewayService).authn-fm]()
- [main.(*gatewayService).authnToken-fm]()
- **/api/v1**
	- **/diagnosis**
		- **/invoke-scan**
			- _POST_
				- [main.(*gatewayService).authzWithProject-fm]()
				- [AllowContentType.func1]()
				- [main.(*gatewayService).invokeDiagnosisScanHandler-fm]()

</details>
<details>
<summary>`/api/v1/diagnosis/list-application-scan`</summary>

- [IgnoreHealthCheckTracingMiddleware.func1]()
- [AnnotateEnvTracingMiddleware.func1]()
- [RequestID]()
- [RealIP]()
- [RequestLogger.func1]()
- [Recoverer]()
- [StripSlashes]()
- [main.(*gatewayService).authn-fm]()
- [main.(*gatewayService).authnToken-fm]()
- **/api/v1**
	- **/diagnosis**
		- **/list-application-scan**
			- _GET_
				- [main.(*gatewayService).authzWithProject-fm]()
				- [main.(*gatewayService).listApplicationScanHandler-fm]()

</details>
<details>
<summary>`/api/v1/diagnosis/list-jira-setting`</summary>

- [IgnoreHealthCheckTracingMiddleware.func1]()
- [AnnotateEnvTracingMiddleware.func1]()
- [RequestID]()
- [RealIP]()
- [RequestLogger.func1]()
- [Recoverer]()
- [StripSlashes]()
- [main.(*gatewayService).authn-fm]()
- [main.(*gatewayService).authnToken-fm]()
- **/api/v1**
	- **/diagnosis**
		- **/list-jira-setting**
			- _GET_
				- [main.(*gatewayService).authzWithProject-fm]()
				- [main.(*gatewayService).listJiraSettingHandler-fm]()

</details>
<details>
<summary>`/api/v1/diagnosis/list-portscan-setting`</summary>

- [IgnoreHealthCheckTracingMiddleware.func1]()
- [AnnotateEnvTracingMiddleware.func1]()
- [RequestID]()
- [RealIP]()
- [RequestLogger.func1]()
- [Recoverer]()
- [StripSlashes]()
- [main.(*gatewayService).authn-fm]()
- [main.(*gatewayService).authnToken-fm]()
- **/api/v1**
	- **/diagnosis**
		- **/list-portscan-setting**
			- _GET_
				- [main.(*gatewayService).authzWithProject-fm]()
				- [main.(*gatewayService).listPortscanSettingHandler-fm]()

</details>
<details>
<summary>`/api/v1/diagnosis/list-portscan-target`</summary>

- [IgnoreHealthCheckTracingMiddleware.func1]()
- [AnnotateEnvTracingMiddleware.func1]()
- [RequestID]()
- [RealIP]()
- [RequestLogger.func1]()
- [Recoverer]()
- [StripSlashes]()
- [main.(*gatewayService).authn-fm]()
- [main.(*gatewayService).authnToken-fm]()
- **/api/v1**
	- **/diagnosis**
		- **/list-portscan-target**
			- _GET_
				- [main.(*gatewayService).authzWithProject-fm]()
				- [main.(*gatewayService).listPortscanTargetHandler-fm]()

</details>
<details>
<summary>`/api/v1/diagnosis/list-wpscan-setting`</summary>

- [IgnoreHealthCheckTracingMiddleware.func1]()
- [AnnotateEnvTracingMiddleware.func1]()
- [RequestID]()
- [RealIP]()
- [RequestLogger.func1]()
- [Recoverer]()
- [StripSlashes]()
- [main.(*gatewayService).authn-fm]()
- [main.(*gatewayService).authnToken-fm]()
- **/api/v1**
	- **/diagnosis**
		- **/list-wpscan-setting**
			- _GET_
				- [main.(*gatewayService).authzWithProject-fm]()
				- [main.(*gatewayService).listWpscanSettingHandler-fm]()

</details>
<details>
<summary>`/api/v1/diagnosis/put-application-scan`</summary>

- [IgnoreHealthCheckTracingMiddleware.func1]()
- [AnnotateEnvTracingMiddleware.func1]()
- [RequestID]()
- [RealIP]()
- [RequestLogger.func1]()
- [Recoverer]()
- [StripSlashes]()
- [main.(*gatewayService).authn-fm]()
- [main.(*gatewayService).authnToken-fm]()
- **/api/v1**
	- **/diagnosis**
		- **/put-application-scan**
			- _POST_
				- [main.(*gatewayService).authzWithProject-fm]()
				- [AllowContentType.func1]()
				- [main.(*gatewayService).putApplicationScanHandler-fm]()

</details>
<details>
<summary>`/api/v1/diagnosis/put-application-scan-basic-setting`</summary>

- [IgnoreHealthCheckTracingMiddleware.func1]()
- [AnnotateEnvTracingMiddleware.func1]()
- [RequestID]()
- [RealIP]()
- [RequestLogger.func1]()
- [Recoverer]()
- [StripSlashes]()
- [main.(*gatewayService).authn-fm]()
- [main.(*gatewayService).authnToken-fm]()
- **/api/v1**
	- **/diagnosis**
		- **/put-application-scan-basic-setting**
			- _POST_
				- [main.(*gatewayService).authzWithProject-fm]()
				- [AllowContentType.func1]()
				- [main.(*gatewayService).putApplicationScanBasicSettingHandler-fm]()

</details>
<details>
<summary>`/api/v1/diagnosis/put-jira-setting`</summary>

- [IgnoreHealthCheckTracingMiddleware.func1]()
- [AnnotateEnvTracingMiddleware.func1]()
- [RequestID]()
- [RealIP]()
- [RequestLogger.func1]()
- [Recoverer]()
- [StripSlashes]()
- [main.(*gatewayService).authn-fm]()
- [main.(*gatewayService).authnToken-fm]()
- **/api/v1**
	- **/diagnosis**
		- **/put-jira-setting**
			- _POST_
				- [main.(*gatewayService).authzOnlyAdmin-fm]()
				- [AllowContentType.func1]()
				- [main.(*gatewayService).putJiraSettingHandler-fm]()

</details>
<details>
<summary>`/api/v1/diagnosis/put-portscan-setting`</summary>

- [IgnoreHealthCheckTracingMiddleware.func1]()
- [AnnotateEnvTracingMiddleware.func1]()
- [RequestID]()
- [RealIP]()
- [RequestLogger.func1]()
- [Recoverer]()
- [StripSlashes]()
- [main.(*gatewayService).authn-fm]()
- [main.(*gatewayService).authnToken-fm]()
- **/api/v1**
	- **/diagnosis**
		- **/put-portscan-setting**
			- _POST_
				- [main.(*gatewayService).authzWithProject-fm]()
				- [AllowContentType.func1]()
				- [main.(*gatewayService).putPortscanSettingHandler-fm]()

</details>
<details>
<summary>`/api/v1/diagnosis/put-portscan-target`</summary>

- [IgnoreHealthCheckTracingMiddleware.func1]()
- [AnnotateEnvTracingMiddleware.func1]()
- [RequestID]()
- [RealIP]()
- [RequestLogger.func1]()
- [Recoverer]()
- [StripSlashes]()
- [main.(*gatewayService).authn-fm]()
- [main.(*gatewayService).authnToken-fm]()
- **/api/v1**
	- **/diagnosis**
		- **/put-portscan-target**
			- _POST_
				- [main.(*gatewayService).authzWithProject-fm]()
				- [AllowContentType.func1]()
				- [main.(*gatewayService).putPortscanTargetHandler-fm]()

</details>
<details>
<summary>`/api/v1/diagnosis/put-wpscan-setting`</summary>

- [IgnoreHealthCheckTracingMiddleware.func1]()
- [AnnotateEnvTracingMiddleware.func1]()
- [RequestID]()
- [RealIP]()
- [RequestLogger.func1]()
- [Recoverer]()
- [StripSlashes]()
- [main.(*gatewayService).authn-fm]()
- [main.(*gatewayService).authnToken-fm]()
- **/api/v1**
	- **/diagnosis**
		- **/put-wpscan-setting**
			- _POST_
				- [main.(*gatewayService).authzWithProject-fm]()
				- [AllowContentType.func1]()
				- [main.(*gatewayService).putWpscanSettingHandler-fm]()

</details>
<details>
<summary>`/api/v1/finding/delete-finding`</summary>

- [IgnoreHealthCheckTracingMiddleware.func1]()
- [AnnotateEnvTracingMiddleware.func1]()
- [RequestID]()
- [RealIP]()
- [RequestLogger.func1]()
- [Recoverer]()
- [StripSlashes]()
- [main.(*gatewayService).authn-fm]()
- [main.(*gatewayService).authnToken-fm]()
- **/api/v1**
	- **/finding**
		- [main.(*gatewayService).authzWithProject-fm]()
		- **/delete-finding**
			- _POST_
				- [AllowContentType.func1]()
				- [main.(*gatewayService).deleteFindingHandler-fm]()

</details>
<details>
<summary>`/api/v1/finding/delete-finding-setting`</summary>

- [IgnoreHealthCheckTracingMiddleware.func1]()
- [AnnotateEnvTracingMiddleware.func1]()
- [RequestID]()
- [RealIP]()
- [RequestLogger.func1]()
- [Recoverer]()
- [StripSlashes]()
- [main.(*gatewayService).authn-fm]()
- [main.(*gatewayService).authnToken-fm]()
- **/api/v1**
	- **/finding**
		- [main.(*gatewayService).authzWithProject-fm]()
		- **/delete-finding-setting**
			- _POST_
				- [AllowContentType.func1]()
				- [main.(*gatewayService).deleteFindingSettingHandler-fm]()

</details>
<details>
<summary>`/api/v1/finding/delete-pend-finding`</summary>

- [IgnoreHealthCheckTracingMiddleware.func1]()
- [AnnotateEnvTracingMiddleware.func1]()
- [RequestID]()
- [RealIP]()
- [RequestLogger.func1]()
- [Recoverer]()
- [StripSlashes]()
- [main.(*gatewayService).authn-fm]()
- [main.(*gatewayService).authnToken-fm]()
- **/api/v1**
	- **/finding**
		- [main.(*gatewayService).authzWithProject-fm]()
		- **/delete-pend-finding**
			- _POST_
				- [AllowContentType.func1]()
				- [main.(*gatewayService).deletePendFindingHandler-fm]()

</details>
<details>
<summary>`/api/v1/finding/delete-resource`</summary>

- [IgnoreHealthCheckTracingMiddleware.func1]()
- [AnnotateEnvTracingMiddleware.func1]()
- [RequestID]()
- [RealIP]()
- [RequestLogger.func1]()
- [Recoverer]()
- [StripSlashes]()
- [main.(*gatewayService).authn-fm]()
- [main.(*gatewayService).authnToken-fm]()
- **/api/v1**
	- **/finding**
		- [main.(*gatewayService).authzWithProject-fm]()
		- **/delete-resource**
			- _POST_
				- [AllowContentType.func1]()
				- [main.(*gatewayService).deleteResourceHandler-fm]()

</details>
<details>
<summary>`/api/v1/finding/get-finding`</summary>

- [IgnoreHealthCheckTracingMiddleware.func1]()
- [AnnotateEnvTracingMiddleware.func1]()
- [RequestID]()
- [RealIP]()
- [RequestLogger.func1]()
- [Recoverer]()
- [StripSlashes]()
- [main.(*gatewayService).authn-fm]()
- [main.(*gatewayService).authnToken-fm]()
- **/api/v1**
	- **/finding**
		- [main.(*gatewayService).authzWithProject-fm]()
		- **/get-finding**
			- _GET_
				- [main.(*gatewayService).getFindingHandler-fm]()

</details>
<details>
<summary>`/api/v1/finding/get-pend-finding`</summary>

- [IgnoreHealthCheckTracingMiddleware.func1]()
- [AnnotateEnvTracingMiddleware.func1]()
- [RequestID]()
- [RealIP]()
- [RequestLogger.func1]()
- [Recoverer]()
- [StripSlashes]()
- [main.(*gatewayService).authn-fm]()
- [main.(*gatewayService).authnToken-fm]()
- **/api/v1**
	- **/finding**
		- [main.(*gatewayService).authzWithProject-fm]()
		- **/get-pend-finding**
			- _GET_
				- [main.(*gatewayService).getPendFindingHandler-fm]()

</details>
<details>
<summary>`/api/v1/finding/get-resource`</summary>

- [IgnoreHealthCheckTracingMiddleware.func1]()
- [AnnotateEnvTracingMiddleware.func1]()
- [RequestID]()
- [RealIP]()
- [RequestLogger.func1]()
- [Recoverer]()
- [StripSlashes]()
- [main.(*gatewayService).authn-fm]()
- [main.(*gatewayService).authnToken-fm]()
- **/api/v1**
	- **/finding**
		- [main.(*gatewayService).authzWithProject-fm]()
		- **/get-resource**
			- _GET_
				- [main.(*gatewayService).getResourceHandler-fm]()

</details>
<details>
<summary>`/api/v1/finding/list-finding`</summary>

- [IgnoreHealthCheckTracingMiddleware.func1]()
- [AnnotateEnvTracingMiddleware.func1]()
- [RequestID]()
- [RealIP]()
- [RequestLogger.func1]()
- [Recoverer]()
- [StripSlashes]()
- [main.(*gatewayService).authn-fm]()
- [main.(*gatewayService).authnToken-fm]()
- **/api/v1**
	- **/finding**
		- [main.(*gatewayService).authzWithProject-fm]()
		- **/list-finding**
			- _GET_
				- [main.(*gatewayService).listFindingHandler-fm]()

</details>
<details>
<summary>`/api/v1/finding/list-finding-setting`</summary>

- [IgnoreHealthCheckTracingMiddleware.func1]()
- [AnnotateEnvTracingMiddleware.func1]()
- [RequestID]()
- [RealIP]()
- [RequestLogger.func1]()
- [Recoverer]()
- [StripSlashes]()
- [main.(*gatewayService).authn-fm]()
- [main.(*gatewayService).authnToken-fm]()
- **/api/v1**
	- **/finding**
		- [main.(*gatewayService).authzWithProject-fm]()
		- **/list-finding-setting**
			- _GET_
				- [main.(*gatewayService).listFindingSettingHandler-fm]()

</details>
<details>
<summary>`/api/v1/finding/list-finding-tag`</summary>

- [IgnoreHealthCheckTracingMiddleware.func1]()
- [AnnotateEnvTracingMiddleware.func1]()
- [RequestID]()
- [RealIP]()
- [RequestLogger.func1]()
- [Recoverer]()
- [StripSlashes]()
- [main.(*gatewayService).authn-fm]()
- [main.(*gatewayService).authnToken-fm]()
- **/api/v1**
	- **/finding**
		- [main.(*gatewayService).authzWithProject-fm]()
		- **/list-finding-tag**
			- _GET_
				- [main.(*gatewayService).listFindingTagHandler-fm]()

</details>
<details>
<summary>`/api/v1/finding/list-finding-tag-name`</summary>

- [IgnoreHealthCheckTracingMiddleware.func1]()
- [AnnotateEnvTracingMiddleware.func1]()
- [RequestID]()
- [RealIP]()
- [RequestLogger.func1]()
- [Recoverer]()
- [StripSlashes]()
- [main.(*gatewayService).authn-fm]()
- [main.(*gatewayService).authnToken-fm]()
- **/api/v1**
	- **/finding**
		- [main.(*gatewayService).authzWithProject-fm]()
		- **/list-finding-tag-name**
			- _GET_
				- [main.(*gatewayService).listFindingTagNameHandler-fm]()

</details>
<details>
<summary>`/api/v1/finding/list-resource`</summary>

- [IgnoreHealthCheckTracingMiddleware.func1]()
- [AnnotateEnvTracingMiddleware.func1]()
- [RequestID]()
- [RealIP]()
- [RequestLogger.func1]()
- [Recoverer]()
- [StripSlashes]()
- [main.(*gatewayService).authn-fm]()
- [main.(*gatewayService).authnToken-fm]()
- **/api/v1**
	- **/finding**
		- [main.(*gatewayService).authzWithProject-fm]()
		- **/list-resource**
			- _GET_
				- [main.(*gatewayService).listResourceHandler-fm]()

</details>
<details>
<summary>`/api/v1/finding/put-finding-setting`</summary>

- [IgnoreHealthCheckTracingMiddleware.func1]()
- [AnnotateEnvTracingMiddleware.func1]()
- [RequestID]()
- [RealIP]()
- [RequestLogger.func1]()
- [Recoverer]()
- [StripSlashes]()
- [main.(*gatewayService).authn-fm]()
- [main.(*gatewayService).authnToken-fm]()
- **/api/v1**
	- **/finding**
		- [main.(*gatewayService).authzWithProject-fm]()
		- **/put-finding-setting**
			- _POST_
				- [AllowContentType.func1]()
				- [main.(*gatewayService).putFindingSettingHandler-fm]()

</details>
<details>
<summary>`/api/v1/finding/put-pend-finding`</summary>

- [IgnoreHealthCheckTracingMiddleware.func1]()
- [AnnotateEnvTracingMiddleware.func1]()
- [RequestID]()
- [RealIP]()
- [RequestLogger.func1]()
- [Recoverer]()
- [StripSlashes]()
- [main.(*gatewayService).authn-fm]()
- [main.(*gatewayService).authnToken-fm]()
- **/api/v1**
	- **/finding**
		- [main.(*gatewayService).authzWithProject-fm]()
		- **/put-pend-finding**
			- _POST_
				- [AllowContentType.func1]()
				- [main.(*gatewayService).putPendFindingHandler-fm]()

</details>
<details>
<summary>`/api/v1/finding/tag-finding`</summary>

- [IgnoreHealthCheckTracingMiddleware.func1]()
- [AnnotateEnvTracingMiddleware.func1]()
- [RequestID]()
- [RealIP]()
- [RequestLogger.func1]()
- [Recoverer]()
- [StripSlashes]()
- [main.(*gatewayService).authn-fm]()
- [main.(*gatewayService).authnToken-fm]()
- **/api/v1**
	- **/finding**
		- [main.(*gatewayService).authzWithProject-fm]()
		- **/tag-finding**
			- _POST_
				- [AllowContentType.func1]()
				- [main.(*gatewayService).tagFindingHandler-fm]()

</details>
<details>
<summary>`/api/v1/finding/untag-finding`</summary>

- [IgnoreHealthCheckTracingMiddleware.func1]()
- [AnnotateEnvTracingMiddleware.func1]()
- [RequestID]()
- [RealIP]()
- [RequestLogger.func1]()
- [Recoverer]()
- [StripSlashes]()
- [main.(*gatewayService).authn-fm]()
- [main.(*gatewayService).authnToken-fm]()
- **/api/v1**
	- **/finding**
		- [main.(*gatewayService).authzWithProject-fm]()
		- **/untag-finding**
			- _POST_
				- [AllowContentType.func1]()
				- [main.(*gatewayService).untagFindingHandler-fm]()

</details>
<details>
<summary>`/api/v1/google/attach-gcp-datasource`</summary>

- [IgnoreHealthCheckTracingMiddleware.func1]()
- [AnnotateEnvTracingMiddleware.func1]()
- [RequestID]()
- [RealIP]()
- [RequestLogger.func1]()
- [Recoverer]()
- [StripSlashes]()
- [main.(*gatewayService).authn-fm]()
- [main.(*gatewayService).authnToken-fm]()
- **/api/v1**
	- **/google**
		- **/attach-gcp-datasource**
			- _POST_
				- [main.(*gatewayService).authzWithProject-fm]()
				- [AllowContentType.func1]()
				- [main.(*gatewayService).attachGCPDataSourceHandler-fm]()

</details>
<details>
<summary>`/api/v1/google/delete-gcp`</summary>

- [IgnoreHealthCheckTracingMiddleware.func1]()
- [AnnotateEnvTracingMiddleware.func1]()
- [RequestID]()
- [RealIP]()
- [RequestLogger.func1]()
- [Recoverer]()
- [StripSlashes]()
- [main.(*gatewayService).authn-fm]()
- [main.(*gatewayService).authnToken-fm]()
- **/api/v1**
	- **/google**
		- **/delete-gcp**
			- _POST_
				- [main.(*gatewayService).authzWithProject-fm]()
				- [AllowContentType.func1]()
				- [main.(*gatewayService).deleteGCPHandler-fm]()

</details>
<details>
<summary>`/api/v1/google/detach-gcp-datasource`</summary>

- [IgnoreHealthCheckTracingMiddleware.func1]()
- [AnnotateEnvTracingMiddleware.func1]()
- [RequestID]()
- [RealIP]()
- [RequestLogger.func1]()
- [Recoverer]()
- [StripSlashes]()
- [main.(*gatewayService).authn-fm]()
- [main.(*gatewayService).authnToken-fm]()
- **/api/v1**
	- **/google**
		- **/detach-gcp-datasource**
			- _POST_
				- [main.(*gatewayService).authzWithProject-fm]()
				- [AllowContentType.func1]()
				- [main.(*gatewayService).detachGCPDataSourceHandler-fm]()

</details>
<details>
<summary>`/api/v1/google/get-gcp-datasource`</summary>

- [IgnoreHealthCheckTracingMiddleware.func1]()
- [AnnotateEnvTracingMiddleware.func1]()
- [RequestID]()
- [RealIP]()
- [RequestLogger.func1]()
- [Recoverer]()
- [StripSlashes]()
- [main.(*gatewayService).authn-fm]()
- [main.(*gatewayService).authnToken-fm]()
- **/api/v1**
	- **/google**
		- **/get-gcp-datasource**
			- _GET_
				- [main.(*gatewayService).authzWithProject-fm]()
				- [main.(*gatewayService).getGCPDataSourceHandler-fm]()

</details>
<details>
<summary>`/api/v1/google/invoke-scan-gcp`</summary>

- [IgnoreHealthCheckTracingMiddleware.func1]()
- [AnnotateEnvTracingMiddleware.func1]()
- [RequestID]()
- [RealIP]()
- [RequestLogger.func1]()
- [Recoverer]()
- [StripSlashes]()
- [main.(*gatewayService).authn-fm]()
- [main.(*gatewayService).authnToken-fm]()
- **/api/v1**
	- **/google**
		- **/invoke-scan-gcp**
			- _POST_
				- [main.(*gatewayService).authzWithProject-fm]()
				- [AllowContentType.func1]()
				- [main.(*gatewayService).invokeScanGCPHandler-fm]()

</details>
<details>
<summary>`/api/v1/google/list-gcp`</summary>

- [IgnoreHealthCheckTracingMiddleware.func1]()
- [AnnotateEnvTracingMiddleware.func1]()
- [RequestID]()
- [RealIP]()
- [RequestLogger.func1]()
- [Recoverer]()
- [StripSlashes]()
- [main.(*gatewayService).authn-fm]()
- [main.(*gatewayService).authnToken-fm]()
- **/api/v1**
	- **/google**
		- **/list-gcp**
			- _GET_
				- [main.(*gatewayService).authzWithProject-fm]()
				- [main.(*gatewayService).listGCPHandler-fm]()

</details>
<details>
<summary>`/api/v1/google/list-google-datasource`</summary>

- [IgnoreHealthCheckTracingMiddleware.func1]()
- [AnnotateEnvTracingMiddleware.func1]()
- [RequestID]()
- [RealIP]()
- [RequestLogger.func1]()
- [Recoverer]()
- [StripSlashes]()
- [main.(*gatewayService).authn-fm]()
- [main.(*gatewayService).authnToken-fm]()
- **/api/v1**
	- **/google**
		- **/list-google-datasource**
			- _GET_
				- [main.(*gatewayService).listGoogleDataSourceHandler-fm]()

</details>
<details>
<summary>`/api/v1/google/put-gcp`</summary>

- [IgnoreHealthCheckTracingMiddleware.func1]()
- [AnnotateEnvTracingMiddleware.func1]()
- [RequestID]()
- [RealIP]()
- [RequestLogger.func1]()
- [Recoverer]()
- [StripSlashes]()
- [main.(*gatewayService).authn-fm]()
- [main.(*gatewayService).authnToken-fm]()
- **/api/v1**
	- **/google**
		- **/put-gcp**
			- _POST_
				- [main.(*gatewayService).authzWithProject-fm]()
				- [AllowContentType.func1]()
				- [main.(*gatewayService).putGCPHandler-fm]()

</details>
<details>
<summary>`/api/v1/iam/attach-access-token`</summary>

- [IgnoreHealthCheckTracingMiddleware.func1]()
- [AnnotateEnvTracingMiddleware.func1]()
- [RequestID]()
- [RealIP]()
- [RequestLogger.func1]()
- [Recoverer]()
- [StripSlashes]()
- [main.(*gatewayService).authn-fm]()
- [main.(*gatewayService).authnToken-fm]()
- **/api/v1**
	- **/iam**
		- **/attach-access-token**
			- _POST_
				- [main.(*gatewayService).authzWithProject-fm]()
				- [AllowContentType.func1]()
				- [main.(*gatewayService).attachAccessTokenRoleHandler-fm]()

</details>
<details>
<summary>`/api/v1/iam/attach-policy`</summary>

- [IgnoreHealthCheckTracingMiddleware.func1]()
- [AnnotateEnvTracingMiddleware.func1]()
- [RequestID]()
- [RealIP]()
- [RequestLogger.func1]()
- [Recoverer]()
- [StripSlashes]()
- [main.(*gatewayService).authn-fm]()
- [main.(*gatewayService).authnToken-fm]()
- **/api/v1**
	- **/iam**
		- **/attach-policy**
			- _POST_
				- [main.(*gatewayService).authzWithProject-fm]()
				- [AllowContentType.func1]()
				- [main.(*gatewayService).attachPolicyHandler-fm]()

</details>
<details>
<summary>`/api/v1/iam/attach-role`</summary>

- [IgnoreHealthCheckTracingMiddleware.func1]()
- [AnnotateEnvTracingMiddleware.func1]()
- [RequestID]()
- [RealIP]()
- [RequestLogger.func1]()
- [Recoverer]()
- [StripSlashes]()
- [main.(*gatewayService).authn-fm]()
- [main.(*gatewayService).authnToken-fm]()
- **/api/v1**
	- **/iam**
		- **/attach-role**
			- _POST_
				- [main.(*gatewayService).authzWithProject-fm]()
				- [AllowContentType.func1]()
				- [main.(*gatewayService).attachRoleHandler-fm]()

</details>
<details>
<summary>`/api/v1/iam/delete-access-token`</summary>

- [IgnoreHealthCheckTracingMiddleware.func1]()
- [AnnotateEnvTracingMiddleware.func1]()
- [RequestID]()
- [RealIP]()
- [RequestLogger.func1]()
- [Recoverer]()
- [StripSlashes]()
- [main.(*gatewayService).authn-fm]()
- [main.(*gatewayService).authnToken-fm]()
- **/api/v1**
	- **/iam**
		- **/delete-access-token**
			- _POST_
				- [main.(*gatewayService).authzWithProject-fm]()
				- [AllowContentType.func1]()
				- [main.(*gatewayService).deleteAccessTokenHandler-fm]()

</details>
<details>
<summary>`/api/v1/iam/delete-policy`</summary>

- [IgnoreHealthCheckTracingMiddleware.func1]()
- [AnnotateEnvTracingMiddleware.func1]()
- [RequestID]()
- [RealIP]()
- [RequestLogger.func1]()
- [Recoverer]()
- [StripSlashes]()
- [main.(*gatewayService).authn-fm]()
- [main.(*gatewayService).authnToken-fm]()
- **/api/v1**
	- **/iam**
		- **/delete-policy**
			- _POST_
				- [main.(*gatewayService).authzWithProject-fm]()
				- [AllowContentType.func1]()
				- [main.(*gatewayService).deletePolicyHandler-fm]()

</details>
<details>
<summary>`/api/v1/iam/delete-role`</summary>

- [IgnoreHealthCheckTracingMiddleware.func1]()
- [AnnotateEnvTracingMiddleware.func1]()
- [RequestID]()
- [RealIP]()
- [RequestLogger.func1]()
- [Recoverer]()
- [StripSlashes]()
- [main.(*gatewayService).authn-fm]()
- [main.(*gatewayService).authnToken-fm]()
- **/api/v1**
	- **/iam**
		- **/delete-role**
			- _POST_
				- [main.(*gatewayService).authzWithProject-fm]()
				- [AllowContentType.func1]()
				- [main.(*gatewayService).deleteRoleHandler-fm]()

</details>
<details>
<summary>`/api/v1/iam/detach-access-token`</summary>

- [IgnoreHealthCheckTracingMiddleware.func1]()
- [AnnotateEnvTracingMiddleware.func1]()
- [RequestID]()
- [RealIP]()
- [RequestLogger.func1]()
- [Recoverer]()
- [StripSlashes]()
- [main.(*gatewayService).authn-fm]()
- [main.(*gatewayService).authnToken-fm]()
- **/api/v1**
	- **/iam**
		- **/detach-access-token**
			- _POST_
				- [main.(*gatewayService).authzWithProject-fm]()
				- [AllowContentType.func1]()
				- [main.(*gatewayService).detachAccessTokenRoleHandler-fm]()

</details>
<details>
<summary>`/api/v1/iam/detach-policy`</summary>

- [IgnoreHealthCheckTracingMiddleware.func1]()
- [AnnotateEnvTracingMiddleware.func1]()
- [RequestID]()
- [RealIP]()
- [RequestLogger.func1]()
- [Recoverer]()
- [StripSlashes]()
- [main.(*gatewayService).authn-fm]()
- [main.(*gatewayService).authnToken-fm]()
- **/api/v1**
	- **/iam**
		- **/detach-policy**
			- _POST_
				- [main.(*gatewayService).authzWithProject-fm]()
				- [AllowContentType.func1]()
				- [main.(*gatewayService).detachPolicyHandler-fm]()

</details>
<details>
<summary>`/api/v1/iam/detach-role`</summary>

- [IgnoreHealthCheckTracingMiddleware.func1]()
- [AnnotateEnvTracingMiddleware.func1]()
- [RequestID]()
- [RealIP]()
- [RequestLogger.func1]()
- [Recoverer]()
- [StripSlashes]()
- [main.(*gatewayService).authn-fm]()
- [main.(*gatewayService).authnToken-fm]()
- **/api/v1**
	- **/iam**
		- **/detach-role**
			- _POST_
				- [main.(*gatewayService).authzWithProject-fm]()
				- [AllowContentType.func1]()
				- [main.(*gatewayService).detachRoleHandler-fm]()

</details>
<details>
<summary>`/api/v1/iam/generate-access-token`</summary>

- [IgnoreHealthCheckTracingMiddleware.func1]()
- [AnnotateEnvTracingMiddleware.func1]()
- [RequestID]()
- [RealIP]()
- [RequestLogger.func1]()
- [Recoverer]()
- [StripSlashes]()
- [main.(*gatewayService).authn-fm]()
- [main.(*gatewayService).authnToken-fm]()
- **/api/v1**
	- **/iam**
		- **/generate-access-token**
			- _POST_
				- [main.(*gatewayService).authzWithProject-fm]()
				- [AllowContentType.func1]()
				- [main.(*gatewayService).generateAccessTokenHandler-fm]()

</details>
<details>
<summary>`/api/v1/iam/get-policy`</summary>

- [IgnoreHealthCheckTracingMiddleware.func1]()
- [AnnotateEnvTracingMiddleware.func1]()
- [RequestID]()
- [RealIP]()
- [RequestLogger.func1]()
- [Recoverer]()
- [StripSlashes]()
- [main.(*gatewayService).authn-fm]()
- [main.(*gatewayService).authnToken-fm]()
- **/api/v1**
	- **/iam**
		- **/get-policy**
			- _GET_
				- [main.(*gatewayService).authzWithProject-fm]()
				- [main.(*gatewayService).getPolicyHandler-fm]()

</details>
<details>
<summary>`/api/v1/iam/get-role`</summary>

- [IgnoreHealthCheckTracingMiddleware.func1]()
- [AnnotateEnvTracingMiddleware.func1]()
- [RequestID]()
- [RealIP]()
- [RequestLogger.func1]()
- [Recoverer]()
- [StripSlashes]()
- [main.(*gatewayService).authn-fm]()
- [main.(*gatewayService).authnToken-fm]()
- **/api/v1**
	- **/iam**
		- **/get-role**
			- _GET_
				- [main.(*gatewayService).authzWithProject-fm]()
				- [main.(*gatewayService).getRoleHandler-fm]()

</details>
<details>
<summary>`/api/v1/iam/get-user`</summary>

- [IgnoreHealthCheckTracingMiddleware.func1]()
- [AnnotateEnvTracingMiddleware.func1]()
- [RequestID]()
- [RealIP]()
- [RequestLogger.func1]()
- [Recoverer]()
- [StripSlashes]()
- [main.(*gatewayService).authn-fm]()
- [main.(*gatewayService).authnToken-fm]()
- **/api/v1**
	- **/iam**
		- **/get-user**
			- _GET_
				- [main.(*gatewayService).getUserHandler-fm]()

</details>
<details>
<summary>`/api/v1/iam/is-admin`</summary>

- [IgnoreHealthCheckTracingMiddleware.func1]()
- [AnnotateEnvTracingMiddleware.func1]()
- [RequestID]()
- [RealIP]()
- [RequestLogger.func1]()
- [Recoverer]()
- [StripSlashes]()
- [main.(*gatewayService).authn-fm]()
- [main.(*gatewayService).authnToken-fm]()
- **/api/v1**
	- **/iam**
		- **/is-admin**
			- _GET_
				- [main.(*gatewayService).isAdminHandler-fm]()

</details>
<details>
<summary>`/api/v1/iam/list-access-token`</summary>

- [IgnoreHealthCheckTracingMiddleware.func1]()
- [AnnotateEnvTracingMiddleware.func1]()
- [RequestID]()
- [RealIP]()
- [RequestLogger.func1]()
- [Recoverer]()
- [StripSlashes]()
- [main.(*gatewayService).authn-fm]()
- [main.(*gatewayService).authnToken-fm]()
- **/api/v1**
	- **/iam**
		- **/list-access-token**
			- _GET_
				- [main.(*gatewayService).authzWithProject-fm]()
				- [main.(*gatewayService).listAccessTokenHandler-fm]()

</details>
<details>
<summary>`/api/v1/iam/list-policy`</summary>

- [IgnoreHealthCheckTracingMiddleware.func1]()
- [AnnotateEnvTracingMiddleware.func1]()
- [RequestID]()
- [RealIP]()
- [RequestLogger.func1]()
- [Recoverer]()
- [StripSlashes]()
- [main.(*gatewayService).authn-fm]()
- [main.(*gatewayService).authnToken-fm]()
- **/api/v1**
	- **/iam**
		- **/list-policy**
			- _GET_
				- [main.(*gatewayService).authzWithProject-fm]()
				- [main.(*gatewayService).listPolicyHandler-fm]()

</details>
<details>
<summary>`/api/v1/iam/list-role`</summary>

- [IgnoreHealthCheckTracingMiddleware.func1]()
- [AnnotateEnvTracingMiddleware.func1]()
- [RequestID]()
- [RealIP]()
- [RequestLogger.func1]()
- [Recoverer]()
- [StripSlashes]()
- [main.(*gatewayService).authn-fm]()
- [main.(*gatewayService).authnToken-fm]()
- **/api/v1**
	- **/iam**
		- **/list-role**
			- _GET_
				- [main.(*gatewayService).authzWithProject-fm]()
				- [main.(*gatewayService).listRoleHandler-fm]()

</details>
<details>
<summary>`/api/v1/iam/list-user`</summary>

- [IgnoreHealthCheckTracingMiddleware.func1]()
- [AnnotateEnvTracingMiddleware.func1]()
- [RequestID]()
- [RealIP]()
- [RequestLogger.func1]()
- [Recoverer]()
- [StripSlashes]()
- [main.(*gatewayService).authn-fm]()
- [main.(*gatewayService).authnToken-fm]()
- **/api/v1**
	- **/iam**
		- **/list-user**
			- _GET_
				- [main.(*gatewayService).listUserHandler-fm]()

</details>
<details>
<summary>`/api/v1/iam/put-policy`</summary>

- [IgnoreHealthCheckTracingMiddleware.func1]()
- [AnnotateEnvTracingMiddleware.func1]()
- [RequestID]()
- [RealIP]()
- [RequestLogger.func1]()
- [Recoverer]()
- [StripSlashes]()
- [main.(*gatewayService).authn-fm]()
- [main.(*gatewayService).authnToken-fm]()
- **/api/v1**
	- **/iam**
		- **/put-policy**
			- _POST_
				- [main.(*gatewayService).authzWithProject-fm]()
				- [AllowContentType.func1]()
				- [main.(*gatewayService).putPolicyHandler-fm]()

</details>
<details>
<summary>`/api/v1/iam/put-role`</summary>

- [IgnoreHealthCheckTracingMiddleware.func1]()
- [AnnotateEnvTracingMiddleware.func1]()
- [RequestID]()
- [RealIP]()
- [RequestLogger.func1]()
- [Recoverer]()
- [StripSlashes]()
- [main.(*gatewayService).authn-fm]()
- [main.(*gatewayService).authnToken-fm]()
- **/api/v1**
	- **/iam**
		- **/put-role**
			- _POST_
				- [main.(*gatewayService).authzWithProject-fm]()
				- [AllowContentType.func1]()
				- [main.(*gatewayService).putRoleHandler-fm]()

</details>
<details>
<summary>`/api/v1/iam/put-user`</summary>

- [IgnoreHealthCheckTracingMiddleware.func1]()
- [AnnotateEnvTracingMiddleware.func1]()
- [RequestID]()
- [RealIP]()
- [RequestLogger.func1]()
- [Recoverer]()
- [StripSlashes]()
- [main.(*gatewayService).authn-fm]()
- [main.(*gatewayService).authnToken-fm]()
- **/api/v1**
	- **/iam**
		- **/put-user**
			- _POST_
				- [AllowContentType.func1]()
				- [main.(*gatewayService).putUserHandler-fm]()

</details>
<details>
<summary>`/api/v1/iam/update-access-token`</summary>

- [IgnoreHealthCheckTracingMiddleware.func1]()
- [AnnotateEnvTracingMiddleware.func1]()
- [RequestID]()
- [RealIP]()
- [RequestLogger.func1]()
- [Recoverer]()
- [StripSlashes]()
- [main.(*gatewayService).authn-fm]()
- [main.(*gatewayService).authnToken-fm]()
- **/api/v1**
	- **/iam**
		- **/update-access-token**
			- _POST_
				- [main.(*gatewayService).authzWithProject-fm]()
				- [AllowContentType.func1]()
				- [main.(*gatewayService).updateAccessTokenHandler-fm]()

</details>
<details>
<summary>`/api/v1/osint/delete-osint`</summary>

- [IgnoreHealthCheckTracingMiddleware.func1]()
- [AnnotateEnvTracingMiddleware.func1]()
- [RequestID]()
- [RealIP]()
- [RequestLogger.func1]()
- [Recoverer]()
- [StripSlashes]()
- [main.(*gatewayService).authn-fm]()
- [main.(*gatewayService).authnToken-fm]()
- **/api/v1**
	- **/osint**
		- **/delete-osint**
			- _POST_
				- [main.(*gatewayService).authzWithProject-fm]()
				- [AllowContentType.func1]()
				- [main.(*gatewayService).deleteOsintHandler-fm]()

</details>
<details>
<summary>`/api/v1/osint/delete-rel-datasource`</summary>

- [IgnoreHealthCheckTracingMiddleware.func1]()
- [AnnotateEnvTracingMiddleware.func1]()
- [RequestID]()
- [RealIP]()
- [RequestLogger.func1]()
- [Recoverer]()
- [StripSlashes]()
- [main.(*gatewayService).authn-fm]()
- [main.(*gatewayService).authnToken-fm]()
- **/api/v1**
	- **/osint**
		- **/delete-rel-datasource**
			- _POST_
				- [main.(*gatewayService).authzWithProject-fm]()
				- [AllowContentType.func1]()
				- [main.(*gatewayService).deleteRelOsintDataSourceHandler-fm]()

</details>
<details>
<summary>`/api/v1/osint/delete-word`</summary>

- [IgnoreHealthCheckTracingMiddleware.func1]()
- [AnnotateEnvTracingMiddleware.func1]()
- [RequestID]()
- [RealIP]()
- [RequestLogger.func1]()
- [Recoverer]()
- [StripSlashes]()
- [main.(*gatewayService).authn-fm]()
- [main.(*gatewayService).authnToken-fm]()
- **/api/v1**
	- **/osint**
		- **/delete-word**
			- _POST_
				- [main.(*gatewayService).authzWithProject-fm]()
				- [AllowContentType.func1]()
				- [main.(*gatewayService).deleteOsintDetectWordHandler-fm]()

</details>
<details>
<summary>`/api/v1/osint/get-osint`</summary>

- [IgnoreHealthCheckTracingMiddleware.func1]()
- [AnnotateEnvTracingMiddleware.func1]()
- [RequestID]()
- [RealIP]()
- [RequestLogger.func1]()
- [Recoverer]()
- [StripSlashes]()
- [main.(*gatewayService).authn-fm]()
- [main.(*gatewayService).authnToken-fm]()
- **/api/v1**
	- **/osint**
		- **/get-osint**
			- _GET_
				- [main.(*gatewayService).authzWithProject-fm]()
				- [main.(*gatewayService).getOsintHandler-fm]()

</details>
<details>
<summary>`/api/v1/osint/invoke-scan`</summary>

- [IgnoreHealthCheckTracingMiddleware.func1]()
- [AnnotateEnvTracingMiddleware.func1]()
- [RequestID]()
- [RealIP]()
- [RequestLogger.func1]()
- [Recoverer]()
- [StripSlashes]()
- [main.(*gatewayService).authn-fm]()
- [main.(*gatewayService).authnToken-fm]()
- **/api/v1**
	- **/osint**
		- **/invoke-scan**
			- _POST_
				- [main.(*gatewayService).authzWithProject-fm]()
				- [AllowContentType.func1]()
				- [main.(*gatewayService).invokeOsintScanHandler-fm]()

</details>
<details>
<summary>`/api/v1/osint/list-datasource`</summary>

- [IgnoreHealthCheckTracingMiddleware.func1]()
- [AnnotateEnvTracingMiddleware.func1]()
- [RequestID]()
- [RealIP]()
- [RequestLogger.func1]()
- [Recoverer]()
- [StripSlashes]()
- [main.(*gatewayService).authn-fm]()
- [main.(*gatewayService).authnToken-fm]()
- **/api/v1**
	- **/osint**
		- **/list-datasource**
			- _GET_
				- [main.(*gatewayService).authzWithProject-fm]()
				- [main.(*gatewayService).listOsintDataSourceHandler-fm]()

</details>
<details>
<summary>`/api/v1/osint/list-osint`</summary>

- [IgnoreHealthCheckTracingMiddleware.func1]()
- [AnnotateEnvTracingMiddleware.func1]()
- [RequestID]()
- [RealIP]()
- [RequestLogger.func1]()
- [Recoverer]()
- [StripSlashes]()
- [main.(*gatewayService).authn-fm]()
- [main.(*gatewayService).authnToken-fm]()
- **/api/v1**
	- **/osint**
		- **/list-osint**
			- _GET_
				- [main.(*gatewayService).authzWithProject-fm]()
				- [main.(*gatewayService).listOsintHandler-fm]()

</details>
<details>
<summary>`/api/v1/osint/list-rel-datasource`</summary>

- [IgnoreHealthCheckTracingMiddleware.func1]()
- [AnnotateEnvTracingMiddleware.func1]()
- [RequestID]()
- [RealIP]()
- [RequestLogger.func1]()
- [Recoverer]()
- [StripSlashes]()
- [main.(*gatewayService).authn-fm]()
- [main.(*gatewayService).authnToken-fm]()
- **/api/v1**
	- **/osint**
		- **/list-rel-datasource**
			- _GET_
				- [main.(*gatewayService).authzWithProject-fm]()
				- [main.(*gatewayService).listRelOsintDataSourceHandler-fm]()

</details>
<details>
<summary>`/api/v1/osint/list-word`</summary>

- [IgnoreHealthCheckTracingMiddleware.func1]()
- [AnnotateEnvTracingMiddleware.func1]()
- [RequestID]()
- [RealIP]()
- [RequestLogger.func1]()
- [Recoverer]()
- [StripSlashes]()
- [main.(*gatewayService).authn-fm]()
- [main.(*gatewayService).authnToken-fm]()
- **/api/v1**
	- **/osint**
		- **/list-word**
			- _GET_
				- [main.(*gatewayService).authzWithProject-fm]()
				- [main.(*gatewayService).listOsintDetectWordHandler-fm]()

</details>
<details>
<summary>`/api/v1/osint/put-osint`</summary>

- [IgnoreHealthCheckTracingMiddleware.func1]()
- [AnnotateEnvTracingMiddleware.func1]()
- [RequestID]()
- [RealIP]()
- [RequestLogger.func1]()
- [Recoverer]()
- [StripSlashes]()
- [main.(*gatewayService).authn-fm]()
- [main.(*gatewayService).authnToken-fm]()
- **/api/v1**
	- **/osint**
		- **/put-osint**
			- _POST_
				- [main.(*gatewayService).authzWithProject-fm]()
				- [AllowContentType.func1]()
				- [main.(*gatewayService).putOsintHandler-fm]()

</details>
<details>
<summary>`/api/v1/osint/put-rel-datasource`</summary>

- [IgnoreHealthCheckTracingMiddleware.func1]()
- [AnnotateEnvTracingMiddleware.func1]()
- [RequestID]()
- [RealIP]()
- [RequestLogger.func1]()
- [Recoverer]()
- [StripSlashes]()
- [main.(*gatewayService).authn-fm]()
- [main.(*gatewayService).authnToken-fm]()
- **/api/v1**
	- **/osint**
		- **/put-rel-datasource**
			- _POST_
				- [main.(*gatewayService).authzWithProject-fm]()
				- [AllowContentType.func1]()
				- [main.(*gatewayService).putRelOsintDataSourceHandler-fm]()

</details>
<details>
<summary>`/api/v1/osint/put-word`</summary>

- [IgnoreHealthCheckTracingMiddleware.func1]()
- [AnnotateEnvTracingMiddleware.func1]()
- [RequestID]()
- [RealIP]()
- [RequestLogger.func1]()
- [Recoverer]()
- [StripSlashes]()
- [main.(*gatewayService).authn-fm]()
- [main.(*gatewayService).authnToken-fm]()
- **/api/v1**
	- **/osint**
		- **/put-word**
			- _POST_
				- [main.(*gatewayService).authzWithProject-fm]()
				- [AllowContentType.func1]()
				- [main.(*gatewayService).putOsintDetectWordHandler-fm]()

</details>
<details>
<summary>`/api/v1/project/create-project`</summary>

- [IgnoreHealthCheckTracingMiddleware.func1]()
- [AnnotateEnvTracingMiddleware.func1]()
- [RequestID]()
- [RealIP]()
- [RequestLogger.func1]()
- [Recoverer]()
- [StripSlashes]()
- [main.(*gatewayService).authn-fm]()
- [main.(*gatewayService).authnToken-fm]()
- **/api/v1**
	- **/project**
		- **/create-project**
			- _POST_
				- [AllowContentType.func1]()
				- [main.(*gatewayService).createProjectHandler-fm]()

</details>
<details>
<summary>`/api/v1/project/delete-project`</summary>

- [IgnoreHealthCheckTracingMiddleware.func1]()
- [AnnotateEnvTracingMiddleware.func1]()
- [RequestID]()
- [RealIP]()
- [RequestLogger.func1]()
- [Recoverer]()
- [StripSlashes]()
- [main.(*gatewayService).authn-fm]()
- [main.(*gatewayService).authnToken-fm]()
- **/api/v1**
	- **/project**
		- **/delete-project**
			- _POST_
				- [main.(*gatewayService).authzWithProject-fm]()
				- [AllowContentType.func1]()
				- [main.(*gatewayService).deleteProjectHandler-fm]()

</details>
<details>
<summary>`/api/v1/project/list-project`</summary>

- [IgnoreHealthCheckTracingMiddleware.func1]()
- [AnnotateEnvTracingMiddleware.func1]()
- [RequestID]()
- [RealIP]()
- [RequestLogger.func1]()
- [Recoverer]()
- [StripSlashes]()
- [main.(*gatewayService).authn-fm]()
- [main.(*gatewayService).authnToken-fm]()
- **/api/v1**
	- **/project**
		- **/list-project**
			- _GET_
				- [main.(*gatewayService).listProjectHandler-fm]()

</details>
<details>
<summary>`/api/v1/project/tag-project`</summary>

- [IgnoreHealthCheckTracingMiddleware.func1]()
- [AnnotateEnvTracingMiddleware.func1]()
- [RequestID]()
- [RealIP]()
- [RequestLogger.func1]()
- [Recoverer]()
- [StripSlashes]()
- [main.(*gatewayService).authn-fm]()
- [main.(*gatewayService).authnToken-fm]()
- **/api/v1**
	- **/project**
		- **/tag-project**
			- _POST_
				- [main.(*gatewayService).authzWithProject-fm]()
				- [AllowContentType.func1]()
				- [main.(*gatewayService).tagProjectHandler-fm]()

</details>
<details>
<summary>`/api/v1/project/untag-project`</summary>

- [IgnoreHealthCheckTracingMiddleware.func1]()
- [AnnotateEnvTracingMiddleware.func1]()
- [RequestID]()
- [RealIP]()
- [RequestLogger.func1]()
- [Recoverer]()
- [StripSlashes]()
- [main.(*gatewayService).authn-fm]()
- [main.(*gatewayService).authnToken-fm]()
- **/api/v1**
	- **/project**
		- **/untag-project**
			- _POST_
				- [main.(*gatewayService).authzWithProject-fm]()
				- [AllowContentType.func1]()
				- [main.(*gatewayService).untagProjectHandler-fm]()

</details>
<details>
<summary>`/api/v1/project/update-project`</summary>

- [IgnoreHealthCheckTracingMiddleware.func1]()
- [AnnotateEnvTracingMiddleware.func1]()
- [RequestID]()
- [RealIP]()
- [RequestLogger.func1]()
- [Recoverer]()
- [StripSlashes]()
- [main.(*gatewayService).authn-fm]()
- [main.(*gatewayService).authnToken-fm]()
- **/api/v1**
	- **/project**
		- **/update-project**
			- _POST_
				- [main.(*gatewayService).authzWithProject-fm]()
				- [AllowContentType.func1]()
				- [main.(*gatewayService).updateProjectHandler-fm]()

</details>
<details>
<summary>`/api/v1/report/get-report`</summary>

- [IgnoreHealthCheckTracingMiddleware.func1]()
- [AnnotateEnvTracingMiddleware.func1]()
- [RequestID]()
- [RealIP]()
- [RequestLogger.func1]()
- [Recoverer]()
- [StripSlashes]()
- [main.(*gatewayService).authn-fm]()
- [main.(*gatewayService).authnToken-fm]()
- **/api/v1**
	- **/report**
		- **/get-report**
			- _GET_
				- [main.(*gatewayService).authzWithProject-fm]()
				- [main.(*gatewayService).getReportFindingHandler-fm]()

</details>
<details>
<summary>`/api/v1/report/get-report-all`</summary>

- [IgnoreHealthCheckTracingMiddleware.func1]()
- [AnnotateEnvTracingMiddleware.func1]()
- [RequestID]()
- [RealIP]()
- [RequestLogger.func1]()
- [Recoverer]()
- [StripSlashes]()
- [main.(*gatewayService).authn-fm]()
- [main.(*gatewayService).authnToken-fm]()
- **/api/v1**
	- **/report**
		- **/get-report-all**
			- _GET_
				- [main.(*gatewayService).authzOnlyAdmin-fm]()
				- [main.(*gatewayService).getReportFindingAllHandler-fm]()

</details>
<details>
<summary>`/api/v1/signin`</summary>

- [IgnoreHealthCheckTracingMiddleware.func1]()
- [AnnotateEnvTracingMiddleware.func1]()
- [RequestID]()
- [RealIP]()
- [RequestLogger.func1]()
- [Recoverer]()
- [StripSlashes]()
- [main.(*gatewayService).authn-fm]()
- [main.(*gatewayService).authnToken-fm]()
- **/api/v1**
	- **/signin**
		- _GET_
			- [main.signinHandler]()

</details>
<details>
<summary>`/healthz`</summary>

- [IgnoreHealthCheckTracingMiddleware.func1]()
- [AnnotateEnvTracingMiddleware.func1]()
- [RequestID]()
- [RealIP]()
- [RequestLogger.func1]()
- [Recoverer]()
- [StripSlashes]()
- [main.(*gatewayService).authn-fm]()
- [main.(*gatewayService).authnToken-fm]()
- **/healthz**
	- _GET_
		- [main.newRouter.func2]()

</details>

Total # of routes: 128
