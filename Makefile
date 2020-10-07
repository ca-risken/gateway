####################################################
## curl example
####################################################
.PHONY: all
all: run

.PHONY: help
help:
	@echo "Usage: make <sub-command>"
	@echo "\n---------------- sub-command list ----------------"
	@cat Makefile | grep -e "^.PHONY:" | grep -v "all" | cut -f2 -d' '

# @see https://github.com/CyberAgent/mimosa-common/tree/master/local
.PHONY: network
network:
	@if [ -z "`docker network ls | grep local-shared`" ]; then docker network create local-shared; fi

.PHONY: go-test
go-test: 
	go test ./...

.PHONY: go-mod-update
go-mod-update:
	go get -u \
			github.com/CyberAgent/mimosa-core/proto/finding \
			github.com/CyberAgent/mimosa-core/proto/iam \
			github.com/CyberAgent/mimosa-core/proto/project \
			github.com/CyberAgent/mimosa-aws/proto/aws

.PHONY: run
run: go-test network
	. env.sh && docker-compose up -d --build

.PHONY: log
log:
	. env.sh && docker-compose logs -f

.PHONY: stop
stop:
	. env.sh && docker-compose down


.PHONY: health-check
health-check:
	curl -is -XGET \
		'http://localhost:8000/healthz/'

.PHONY: signin
signin:
	curl -is -XGET \
		--header 'x-amzn-oidc-identity: alice' \
		'http://localhost:8000/api/v1/signin/'

.PHONY: list-finding
list-finding:
	curl -is -XGET \
		--header 'x-amzn-oidc-identity: alice' \
		--header 'X-XSRF-TOKEN: xxxxxxxxx' \
		--header 'Cookie: XSRF-TOKEN=xxxxxxxxx;' \
		'http://localhost:8000/api/v1/finding/list-finding/?project_id=1001&data_source=aws:guardduty,aws:access-analizer'

.PHONY: get-finding
get-finding:
	curl -is -XGET \
		--header 'x-amzn-oidc-identity: alice' \
		--header 'X-XSRF-TOKEN: xxxxxxxxx' \
		--header 'Cookie: XSRF-TOKEN=xxxxxxxxx;' \
		'http://localhost:8000/api/v1/finding/get-finding/?project_id=1001&finding_id=1001'

.PHONY: put-finding
put-finding:
	curl -is -XPOST \
		--header 'x-amzn-oidc-identity: alice' \
		--header 'X-XSRF-TOKEN: xxxxxxxxx' \
		--header 'Cookie: XSRF-TOKEN=xxxxxxxxx;' \
		--header 'Content-Type: application/json' \
		--data '{"project_id":1001, "finding":{"description":"desc", "data_source":"ds", "data_source_id":"ds-004", "resource_name":"rn", "project_id":1001, "original_score":55.51, "original_max_score":100.0, "data":"{\"key\":\"value\"}"}}' \
		'http://localhost:8000/api/v1/finding/put-finding/'

.PHONY: delete-finding
delete-finding:
	curl -is -XPOST \
		--header 'x-amzn-oidc-identity: alice' \
		--header 'X-XSRF-TOKEN: xxxxxxxxx' \
		--header 'Cookie: XSRF-TOKEN=xxxxxxxxx;' \
		--header 'Content-Type: application/json' \
		--data '{"project_id":1001, "finding_id":1005}' \
		'http://localhost:8000/api/v1/finding/delete-finding/'

.PHONY: list-finding-tag
list-finding-tag:
	curl -is -XGET \
		--header 'x-amzn-oidc-identity: alice' \
		--header 'X-XSRF-TOKEN: xxxxxxxxx' \
		--header 'Cookie: XSRF-TOKEN=xxxxxxxxx;' \
		'http://localhost:8000/api/v1/finding/list-finding-tag/?project_id=1001&finding_id=1001'

.PHONY: tag-finding
tag-finding:
	curl -is -XPOST \
		--header 'x-amzn-oidc-identity: alice' \
		--header 'X-XSRF-TOKEN: xxxxxxxxx' \
		--header 'Cookie: XSRF-TOKEN=xxxxxxxxx;' \
		--header 'Content-Type: application/json' \
		--data '{"project_id":1001, "tag":{"finding_id":1001, "project_id":1001, "tag":"tag"}}' \
		'http://localhost:8000/api/v1/finding/tag-finding/'

.PHONY: untag-finding
untag-finding:
	curl -is -XPOST \
		--header 'x-amzn-oidc-identity: alice' \
		--header 'X-XSRF-TOKEN: xxxxxxxxx' \
		--header 'Cookie: XSRF-TOKEN=xxxxxxxxx;' \
		--header 'Content-Type: application/json' \
		--data '{"project_id":1001, "finding_tag_id":1002}' \
		'http://localhost:8000/api/v1/finding/untag-finding/'

.PHONY: list-resource
list-resource:
	curl -is -XGET \
		--header 'x-amzn-oidc-identity: alice' \
		--header 'X-XSRF-TOKEN: xxxxxxxxx' \
		--header 'Cookie: XSRF-TOKEN=xxxxxxxxx;' \
		'http://localhost:8000/api/v1/finding/list-resource/?project_id=1001'

.PHONY: get-resource
get-resource:
	curl -is -XGET \
		--header 'x-amzn-oidc-identity: alice' \
		--header 'X-XSRF-TOKEN: xxxxxxxxx' \
		--header 'Cookie: XSRF-TOKEN=xxxxxxxxx;' \
		'http://localhost:8000/api/v1/finding/get-resource/?project_id=1001&resource_id=1001'

.PHONY: put-resource
put-resource:
	curl -is -XPOST \
		--header 'x-amzn-oidc-identity: alice' \
		--header 'X-XSRF-TOKEN: xxxxxxxxx' \
		--header 'Cookie: XSRF-TOKEN=xxxxxxxxx;' \
		--header 'Content-Type: application/json' \
		--data '{"project_id":1001, "resource":{"resource_name":"rn", "project_id":1001}}' \
		'http://localhost:8000/api/v1/finding/put-resource/'

.PHONY: delete-resource
delete-resource:
	curl -is -XPOST \
		--header 'x-amzn-oidc-identity: alice' \
		--header 'X-XSRF-TOKEN: xxxxxxxxx' \
		--header 'Cookie: XSRF-TOKEN=xxxxxxxxx;' \
		--header 'Content-Type: application/json' \
		--data '{"project_id":1001, "resource_id":1003}' \
		'http://localhost:8000/api/v1/finding/delete-resource/'

.PHONY: list-resource-tag
list-resource-tag:
	curl -is -XGET \
		--header 'x-amzn-oidc-identity: alice' \
		--header 'X-XSRF-TOKEN: xxxxxxxxx' \
		--header 'Cookie: XSRF-TOKEN=xxxxxxxxx;' \
		'http://localhost:8000/api/v1/finding/list-resource-tag/?project_id=1001&resource_id=1001'

.PHONY: tag-resource
tag-resource:
	curl -is -XPOST \
		--header 'x-amzn-oidc-identity: alice' \
		--header 'X-XSRF-TOKEN: xxxxxxxxx' \
		--header 'Cookie: XSRF-TOKEN=xxxxxxxxx;' \
		--header 'Content-Type: application/json' \
		--data '{"project_id":1001, "tag":{"resource_id":1001, "project_id":1001, "tag":"tag"}}' \
		'http://localhost:8000/api/v1/finding/tag-resource/'

.PHONY: untag-resource
untag-resource:
	curl -is -XPOST \
		--header 'x-amzn-oidc-identity: alice' \
		--header 'X-XSRF-TOKEN: xxxxxxxxx' \
		--header 'Cookie: XSRF-TOKEN=xxxxxxxxx;' \
		--header 'Content-Type: application/json' \
		--data '{"project_id":1001, "resource_tag_id":1004}' \
		'http://localhost:8000/api/v1/finding/untag-resource/'

.PHONY: list-user
list-user:
	curl -is -XGET \
		--header 'x-amzn-oidc-identity: alice' \
		--header 'X-XSRF-TOKEN: xxxxxxxxx' \
		--header 'Cookie: XSRF-TOKEN=xxxxxxxxx;' \
		'http://localhost:8000/api/v1/iam/list-user/'

.PHONY: get-user
get-user:
	curl -is -XGET \
		--header 'x-amzn-oidc-identity: alice' \
		--header 'X-XSRF-TOKEN: xxxxxxxxx' \
		--header 'Cookie: XSRF-TOKEN=xxxxxxxxx;' \
		'http://localhost:8000/api/v1/iam/get-user/?user_id=1001'

.PHONY: put-user
put-user:
	curl -is -XPOST \
		--header 'x-amzn-oidc-identity: alice' \
		--header 'X-XSRF-TOKEN: xxxxxxxxx' \
		--header 'Cookie: XSRF-TOKEN=xxxxxxxxx;' \
		--header 'Content-Type: application/json' \
		--data '{"user":{"sub":"sub", "name":"nm", "activated":true}}' \
		'http://localhost:8000/api/v1/iam/put-user/'

.PHONY: list-role
list-role:
	curl -is -XGET \
		--header 'x-amzn-oidc-identity: alice' \
		--header 'X-XSRF-TOKEN: xxxxxxxxx' \
		--header 'Cookie: XSRF-TOKEN=xxxxxxxxx;' \
		'http://localhost:8000/api/v1/iam/list-role/?project_id=1001'

.PHONY: get-role
get-role:
	curl -is -XGET \
		--header 'x-amzn-oidc-identity: alice' \
		--header 'X-XSRF-TOKEN: xxxxxxxxx' \
		--header 'Cookie: XSRF-TOKEN=xxxxxxxxx;' \
		'http://localhost:8000/api/v1/iam/get-role/?project_id=1001&role_id=1001'

.PHONY: put-role
put-role:
	curl -is -XPOST \
		--header 'x-amzn-oidc-identity: alice' \
		--header 'X-XSRF-TOKEN: xxxxxxxxx' \
		--header 'Cookie: XSRF-TOKEN=xxxxxxxxx;' \
		--header 'Content-Type: application/json' \
		--data '{"project_id":1001, "role":{"project_id":1001, "name":"nm"}}' \
		'http://localhost:8000/api/v1/iam/put-role/'

.PHONY: delete-role
delete-role:
	curl -is -XPOST \
		--header 'x-amzn-oidc-identity: alice' \
		--header 'X-XSRF-TOKEN: xxxxxxxxx' \
		--header 'Cookie: XSRF-TOKEN=xxxxxxxxx;' \
		--header 'Content-Type: application/json' \
		--data '{"project_id":1001, "role_id":1008}' \
		'http://localhost:8000/api/v1/iam/delete-role/'

.PHONY: attach-role
attach-role:
	curl -is -XPOST \
		--header 'x-amzn-oidc-identity: alice' \
		--header 'X-XSRF-TOKEN: xxxxxxxxx' \
		--header 'Cookie: XSRF-TOKEN=xxxxxxxxx;' \
		--header 'Content-Type: application/json' \
		--data '{"project_id":1001, "user_id":1001, "role_id":1005}' \
		'http://localhost:8000/api/v1/iam/attach-role/'

.PHONY: detach-role
detach-role:
	curl -is -XPOST \
		--header 'x-amzn-oidc-identity: alice' \
		--header 'X-XSRF-TOKEN: xxxxxxxxx' \
		--header 'Cookie: XSRF-TOKEN=xxxxxxxxx;' \
		--header 'Content-Type: application/json' \
		--data '{"project_id":1001, "user_id":1001, "role_id":1005}' \
		'http://localhost:8000/api/v1/iam/detach-role/'

.PHONY: list-policy
list-policy:
	curl -is -XGET \
		--header 'x-amzn-oidc-identity: alice' \
		--header 'X-XSRF-TOKEN: xxxxxxxxx' \
		--header 'Cookie: XSRF-TOKEN=xxxxxxxxx;' \
		'http://localhost:8000/api/v1/iam/list-policy/?project_id=1001'

.PHONY: get-policy
get-policy:
	curl -is -XGET \
		--header 'x-amzn-oidc-identity: alice' \
		--header 'X-XSRF-TOKEN: xxxxxxxxx' \
		--header 'Cookie: XSRF-TOKEN=xxxxxxxxx;' \
		'http://localhost:8000/api/v1/iam/get-policy/?project_id=1001&policy_id=1001'

.PHONY: put-policy
put-policy:
	curl -is -XPOST \
		--header 'x-amzn-oidc-identity: alice' \
		--header 'X-XSRF-TOKEN: xxxxxxxxx' \
		--header 'Cookie: XSRF-TOKEN=xxxxxxxxx;' \
		--header 'Content-Type: application/json' \
		--data '{"project_id":1001, "policy":{"name":"nm", "project_id":1001, "action_ptn":".*", "resource_ptn":".*"}}' \
		'http://localhost:8000/api/v1/iam/put-policy/'

.PHONY: delete-policy
delete-policy:
	curl -is -XPOST \
		--header 'x-amzn-oidc-identity: alice' \
		--header 'X-XSRF-TOKEN: xxxxxxxxx' \
		--header 'Cookie: XSRF-TOKEN=xxxxxxxxx;' \
		--header 'Content-Type: application/json' \
		--data '{"project_id":1001, "policy_id":1008}' \
		'http://localhost:8000/api/v1/iam/delete-policy/'

.PHONY: attach-policy
attach-policy:
	curl -is -XPOST \
		--header 'x-amzn-oidc-identity: alice' \
		--header 'X-XSRF-TOKEN: xxxxxxxxx' \
		--header 'Cookie: XSRF-TOKEN=xxxxxxxxx;' \
		--header 'Content-Type: application/json' \
		--data '{"project_id":1001, "role_id":1001, "policy_id":1005}' \
		'http://localhost:8000/api/v1/iam/attach-policy/'

.PHONY: detach-policy
detach-policy:
	curl -is -XPOST \
		--header 'x-amzn-oidc-identity: alice' \
		--header 'X-XSRF-TOKEN: xxxxxxxxx' \
		--header 'Cookie: XSRF-TOKEN=xxxxxxxxx;' \
		--header 'Content-Type: application/json' \
		--data '{"project_id":1001, "role_id":1001, "policy_id":1005}' \
		'http://localhost:8000/api/v1/iam/detach-policy/'

.PHONY: list-project
list-project:
	curl -is -XGET \
		--header 'x-amzn-oidc-identity: alice' \
		--header 'X-XSRF-TOKEN: xxxxxxxxx' \
		--header 'Cookie: XSRF-TOKEN=xxxxxxxxx;' \
		'http://localhost:8000/api/v1/project/list-project/?user_id=1001'

.PHONY: create-project
create-project:
	curl -is -XPOST \
		--header 'x-amzn-oidc-identity: alice' \
		--header 'X-XSRF-TOKEN: xxxxxxxxx' \
		--header 'Cookie: XSRF-TOKEN=xxxxxxxxx;' \
		--header 'Content-Type: application/json' \
		--data '{"user_id":1001, "name":"test-pj"}' \
		'http://localhost:8000/api/v1/project/create-project/'

.PHONY: update-project
update-project:
	curl -is -XPOST \
		--header 'x-amzn-oidc-identity: alice' \
		--header 'X-XSRF-TOKEN: xxxxxxxxx' \
		--header 'Cookie: XSRF-TOKEN=xxxxxxxxx;' \
		--header 'Content-Type: application/json' \
		--data '{"project_id":1006, "name":"test-pj-x"}' \
		'http://localhost:8000/api/v1/project/update-project/'

.PHONY: delete-project
delete-project:
	curl -is -XPOST \
		--header 'x-amzn-oidc-identity: alice' \
		--header 'X-XSRF-TOKEN: xxxxxxxxx' \
		--header 'Cookie: XSRF-TOKEN=xxxxxxxxx;' \
		--header 'Content-Type: application/json' \
		--data '{"project_id":1007}' \
		'http://localhost:8000/api/v1/project/delete-project/'

.PHONY: list-alert
list-alert:
	curl -is -XGET \
		--header 'x-amzn-oidc-identity: alice' \
		--header 'X-XSRF-TOKEN: xxxxxxxxx' \
		--header 'Cookie: XSRF-TOKEN=xxxxxxxxx;' \
		'http://localhost:8000/api/v1/alert/list-alert/?project_id=1001'

.PHONY: get-alert
get-alert:
	curl -is -XGET \
		--header 'x-amzn-oidc-identity: alice' \
		--header 'X-XSRF-TOKEN: xxxxxxxxx' \
		--header 'Cookie: XSRF-TOKEN=xxxxxxxxx;' \
		'http://localhost:8000/api/v1/alert/get-alert/?project_id=1001&alert_id=1001'

.PHONY: upsert-alert
upsert-alert:
	curl -is -XPOST \
		--header 'x-amzn-oidc-identity: alice' \
		--header 'X-XSRF-TOKEN: xxxxxxxxx' \
		--header 'Cookie: XSRF-TOKEN=xxxxxxxxx;' \
		--header 'Content-Type: application/json' \
		--data '{"project_id":1001,"alert":{"alert_condition_id":1001,"description":"test_desc","severity":"high","project_id":1001,"activated":true}}' \
		'http://localhost:8000/api/v1/alert/put-alert/'

.PHONY: delete-alert
delete-alert:
	curl -is -XPOST \
		--header 'x-amzn-oidc-identity: alice' \
		--header 'X-XSRF-TOKEN: xxxxxxxxx' \
		--header 'Cookie: XSRF-TOKEN=xxxxxxxxx;' \
		--header 'Content-Type: application/json' \
		--data '{"alert_id":1001,"project_id":1001}' \
		'http://localhost:8000/api/v1/alert/delete-alert/'

.PHONY: list-alert_history
list-alert_history:
	curl -is -XGET \
		--header 'x-amzn-oidc-identity: alice' \
		--header 'X-XSRF-TOKEN: xxxxxxxxx' \
		--header 'Cookie: XSRF-TOKEN=xxxxxxxxx;' \
		'http://localhost:8000/api/v1/alert/list-history/?project_id=1001'

.PHONY: get-alert_history
get-alert_history:
	curl -is -XGET \
		--header 'x-amzn-oidc-identity: alice' \
		--header 'X-XSRF-TOKEN: xxxxxxxxx' \
		--header 'Cookie: XSRF-TOKEN=xxxxxxxxx;' \
		'http://localhost:8000/api/v1/alert/get-history/?project_id=1001&alert_history_id=1001'

.PHONY: upsert-alert_history
upsert-alert_history:
	curl -is -XPOST \
		--header 'x-amzn-oidc-identity: alice' \
		--header 'X-XSRF-TOKEN: xxxxxxxxx' \
		--header 'Cookie: XSRF-TOKEN=xxxxxxxxx;' \
		--header 'Content-Type: application/json' \
		--data '{"project_id":1001,"alert_history":{"alert_id":1021,"description":"test_desc","severity":"high","project_id":1001,"history_type":"created"}}' \
		'http://localhost:8000/api/v1/alert/put-history/'

.PHONY: delete-alert_history
delete-alert_history:
	curl -is -XPOST \
		--header 'x-amzn-oidc-identity: alice' \
		--header 'X-XSRF-TOKEN: xxxxxxxxx' \
		--header 'Cookie: XSRF-TOKEN=xxxxxxxxx;' \
		--header 'Content-Type: application/json' \
		--data '{"alert_history_id":1001,"project_id":1001}' \
		'http://localhost:8000/api/v1/alert/delete-history/'

.PHONY: list-rel_alert_finding
list-rel_alert_finding:
	curl -is -XGET \
		--header 'x-amzn-oidc-identity: alice' \
		--header 'X-XSRF-TOKEN: xxxxxxxxx' \
		--header 'Cookie: XSRF-TOKEN=xxxxxxxxx;' \
		'http://localhost:8000/api/v1/alert/list-rel_alert_finding/?project_id=1001'

.PHONY: get-rel_alert_finding
get-rel_alert_finding:
	curl -is -XGET \
		--header 'x-amzn-oidc-identity: alice' \
		--header 'X-XSRF-TOKEN: xxxxxxxxx' \
		--header 'Cookie: XSRF-TOKEN=xxxxxxxxx;' \
		'http://localhost:8000/api/v1/alert/get-rel_alert_finding/?project_id=1001&alert_id=1001&finding_id=1001'

.PHONY: upsert-rel_alert_finding
upsert-rel_alert_finding:
	curl -is -XPOST \
		--header 'x-amzn-oidc-identity: alice' \
		--header 'X-XSRF-TOKEN: xxxxxxxxx' \
		--header 'Cookie: XSRF-TOKEN=xxxxxxxxx;' \
		--header 'Content-Type: application/json' \
		--data '{"project_id":1001,"rel_alert_finding":{"alert_id":1021,"finding_id":1001,"project_id":1001}}' \
		'http://localhost:8000/api/v1/alert/put-rel_alert_finding/'

.PHONY: delete-rel_alert_finding
delete-rel_alert_finding:
	curl -is -XPOST \
		--header 'x-amzn-oidc-identity: alice' \
		--header 'X-XSRF-TOKEN: xxxxxxxxx' \
		--header 'Cookie: XSRF-TOKEN=xxxxxxxxx;' \
		--header 'Content-Type: application/json' \
		--data '{"alert_id":1001,"finding_id":1001,"project_id":1001}' \
		'http://localhost:8000/api/v1/alert/delete-rel_alert_finding/'

.PHONY: list-alert_condition
list-alert_condition:
	curl -is -XGET \
		--header 'x-amzn-oidc-identity: alice' \
		--header 'X-XSRF-TOKEN: xxxxxxxxx' \
		--header 'Cookie: XSRF-TOKEN=xxxxxxxxx;' \
		'http://localhost:8000/api/v1/alert/list-condition/?project_id=1001'

.PHONY: get-alert_condition
get-alert_condition:
	curl -is -XGET \
		--header 'x-amzn-oidc-identity: alice' \
		--header 'X-XSRF-TOKEN: xxxxxxxxx' \
		--header 'Cookie: XSRF-TOKEN=xxxxxxxxx;' \
		'http://localhost:8000/api/v1/alert/get-condition/?project_id=1001&alert_condition_id=1001'

.PHONY: upsert-alert_condition
upsert-alert_condition:
	curl -is -XPOST \
		--header 'x-amzn-oidc-identity: alice' \
		--header 'X-XSRF-TOKEN: xxxxxxxxx' \
		--header 'Cookie: XSRF-TOKEN=xxxxxxxxx;' \
		--header 'Content-Type: application/json' \
		--data '{"project_id":1001,"alert_condition":{"description":"test_desc","severity":"high","project_id":1001,"and_or":"and","enabled":true}}' \
		'http://localhost:8000/api/v1/alert/put-condition/'

.PHONY: delete-alert_condition
delete-alert_condition:
	curl -is -XPOST \
		--header 'x-amzn-oidc-identity: alice' \
		--header 'X-XSRF-TOKEN: xxxxxxxxx' \
		--header 'Cookie: XSRF-TOKEN=xxxxxxxxx;' \
		--header 'Content-Type: application/json' \
		--data '{"alert_condition_id":1001,"project_id":1001}' \
		'http://localhost:8000/api/v1/alert/delete-condition/'

.PHONY: list-alert_rule
list-alert_rule:
	curl -is -XGET \
		--header 'x-amzn-oidc-identity: alice' \
		--header 'X-XSRF-TOKEN: xxxxxxxxx' \
		--header 'Cookie: XSRF-TOKEN=xxxxxxxxx;' \
		'http://localhost:8000/api/v1/alert/list-rule/?project_id=1001'

.PHONY: get-alert_rule
get-alert_rule:
	curl -is -XGET \
		--header 'x-amzn-oidc-identity: alice' \
		--header 'X-XSRF-TOKEN: xxxxxxxxx' \
		--header 'Cookie: XSRF-TOKEN=xxxxxxxxx;' \
		'http://localhost:8000/api/v1/alert/get-rule/?project_id=1001&alert_rule_id=1001'

.PHONY: upsert-alert_rule
upsert-alert_rule:
	curl -is -XPOST \
		--header 'x-amzn-oidc-identity: alice' \
		--header 'X-XSRF-TOKEN: xxxxxxxxx' \
		--header 'Cookie: XSRF-TOKEN=xxxxxxxxx;' \
		--header 'Content-Type: application/json' \
		--data '{"project_id":1001,"alert_rule":{"name":"test_rule","score":1.0,"resource_name":"test_resource","tag":"test_tag","finding_cnt": 1,"project_id":1001}}' \
		'http://localhost:8000/api/v1/alert/put-rule/'

.PHONY: delete-alert_rule
delete-alert_rule:
	curl -is -XPOST \
		--header 'x-amzn-oidc-identity: alice' \
		--header 'X-XSRF-TOKEN: xxxxxxxxx' \
		--header 'Cookie: XSRF-TOKEN=xxxxxxxxx;' \
		--header 'Content-Type: application/json' \
		--data '{"alert_rule_id":1001,"project_id":1001}' \
		'http://localhost:8000/api/v1/alert/delete-rule/'

.PHONY: list-alert_cond_rule
list-alert_cond_rule:
	curl -is -XGET \
		--header 'x-amzn-oidc-identity: alice' \
		--header 'X-XSRF-TOKEN: xxxxxxxxx' \
		--header 'Cookie: XSRF-TOKEN=xxxxxxxxx;' \
		'http://localhost:8000/api/v1/alert/list-condition_rule/?project_id=1001'

.PHONY: get-alert_cond_rule
get-alert_cond_rule:
	curl -is -XGET \
		--header 'x-amzn-oidc-identity: alice' \
		--header 'X-XSRF-TOKEN: xxxxxxxxx' \
		--header 'Cookie: XSRF-TOKEN=xxxxxxxxx;' \
		'http://localhost:8000/api/v1/alert/get-condition_rule/?project_id=1001&alert_condition_id=1001&alert_rule_id=1001'

.PHONY: upsert-alert_cond_rule
upsert-alert_cond_rule:
	curl -is -XPOST \
		--header 'x-amzn-oidc-identity: alice' \
		--header 'X-XSRF-TOKEN: xxxxxxxxx' \
		--header 'Cookie: XSRF-TOKEN=xxxxxxxxx;' \
		--header 'Content-Type: application/json' \
		--data '{"project_id":1001,"alert_cond_rule":{"alert_condition_id":1001,"alert_rule_id":1001,"project_id":1001}}' \
		'http://localhost:8000/api/v1/alert/put-condition_rule/'

.PHONY: delete-alert_cond_rule
delete-alert_cond_rule:
	curl -is -XPOST \
		--header 'x-amzn-oidc-identity: alice' \
		--header 'X-XSRF-TOKEN: xxxxxxxxx' \
		--header 'Cookie: XSRF-TOKEN=xxxxxxxxx;' \
		--header 'Content-Type: application/json' \
		--data '{"alert_condition_id":1001,"alert_rule_id":1001,"project_id":1001}' \
		'http://localhost:8000/api/v1/alert/delete-condition_rule/'

.PHONY: list-notification
list-notification:
	curl -is -XGET \
		--header 'x-amzn-oidc-identity: alice' \
		--header 'X-XSRF-TOKEN: xxxxxxxxx' \
		--header 'Cookie: XSRF-TOKEN=xxxxxxxxx;' \
		'http://localhost:8000/api/v1/alert/list-notification/?project_id=1001'

.PHONY: get-notification
get-notification:
	curl -is -XGET \
		--header 'x-amzn-oidc-identity: alice' \
		--header 'X-XSRF-TOKEN: xxxxxxxxx' \
		--header 'Cookie: XSRF-TOKEN=xxxxxxxxx;' \
		'http://localhost:8000/api/v1/alert/get-notification/?project_id=1001&notification_id=1001'

.PHONY: upsert-notification
upsert-notification:
	curl -is -XPOST \
		--header 'x-amzn-oidc-identity: alice' \
		--header 'X-XSRF-TOKEN: xxxxxxxxx' \
		--header 'Cookie: XSRF-TOKEN=xxxxxxxxx;' \
		--header 'Content-Type: application/json' \
		--data '{"project_id":1001,"notification":{"name":"test_notification","type":"slack","notify_setting":"{\"test_key\":\"test_value\"}","project_id":1001}}' \
		'http://localhost:8000/api/v1/alert/put-notification/'

.PHONY: delete-notification
delete-notification:
	curl -is -XPOST \
		--header 'x-amzn-oidc-identity: alice' \
		--header 'X-XSRF-TOKEN: xxxxxxxxx' \
		--header 'Cookie: XSRF-TOKEN=xxxxxxxxx;' \
		--header 'Content-Type: application/json' \
		--data '{"notification_id":1001,"project_id":1001}' \
		'http://localhost:8000/api/v1/alert/delete-notification/'

.PHONY: list-alert_cond_notification
list-alert_cond_notification:
	curl -is -XGET \
		--header 'x-amzn-oidc-identity: alice' \
		--header 'X-XSRF-TOKEN: xxxxxxxxx' \
		--header 'Cookie: XSRF-TOKEN=xxxxxxxxx;' \
		'http://localhost:8000/api/v1/alert/list-condition_notification/?project_id=1001'

.PHONY: get-alert_cond_notification
get-alert_cond_notification:
	curl -is -XGET \
		--header 'x-amzn-oidc-identity: alice' \
		--header 'X-XSRF-TOKEN: xxxxxxxxx' \
		--header 'Cookie: XSRF-TOKEN=xxxxxxxxx;' \
		'http://localhost:8000/api/v1/alert/get-condition_notification/?project_id=1001&alert_condition_id=1001&notification_id=1001'

.PHONY: upsert-alert_cond_notification
upsert-alert_cond_notification:
	curl -is -XPOST \
		--header 'x-amzn-oidc-identity: alice' \
		--header 'X-XSRF-TOKEN: xxxxxxxxx' \
		--header 'Cookie: XSRF-TOKEN=xxxxxxxxx;' \
		--header 'Content-Type: application/json' \
		--data '{"project_id":1001,"alert_cond_notification":{"alert_condition_id":1001,"notification_id":1001,"project_id":1001}}' \
		'http://localhost:8000/api/v1/alert/put-condition_notification/'

.PHONY: delete-alert_cond_notification
delete-alert_cond_notification:
	curl -is -XPOST \
		--header 'x-amzn-oidc-identity: alice' \
		--header 'X-XSRF-TOKEN: xxxxxxxxx' \
		--header 'Cookie: XSRF-TOKEN=xxxxxxxxx;' \
		--header 'Content-Type: application/json' \
		--data '{"alert_condition_id":1001,"notification_id":1001,"project_id":1001}' \
		'http://localhost:8000/api/v1/alert/delete-condition_notification/'

.PHONY: list-aws
list-aws:
	curl -is -XGET \
		--header 'x-amzn-oidc-identity: alice' \
		--header 'X-XSRF-TOKEN: xxxxxxxxx' \
		--header 'Cookie: XSRF-TOKEN=xxxxxxxxx;' \
		'http://localhost:8000/api/v1/aws/list-aws/?project_id=1001'

.PHONY: put-aws
put-aws:
	curl -is -XPOST \
		--header 'x-amzn-oidc-identity: alice' \
		--header 'X-XSRF-TOKEN: xxxxxxxxx' \
		--header 'Cookie: XSRF-TOKEN=xxxxxxxxx;' \
		--header 'Content-Type: application/json' \
		--data '{"project_id":1001, "aws":{"project_id":1001, "name":"hoge-aws", "aws_account_id":"123456789012"}}' \
		'http://localhost:8000/api/v1/aws/put-aws/'

.PHONY: delete-aws
delete-aws:
	curl -is -XPOST \
		--header 'x-amzn-oidc-identity: alice' \
		--header 'X-XSRF-TOKEN: xxxxxxxxx' \
		--header 'Cookie: XSRF-TOKEN=xxxxxxxxx;' \
		--header 'Content-Type: application/json' \
		--data '{"project_id":1001, "aws_id":1002}' \
		'http://localhost:8000/api/v1/aws/delete-aws/'

.PHONY: list-datasource
list-datasource:
	curl -is -XGET \
		--header 'x-amzn-oidc-identity: alice' \
		--header 'X-XSRF-TOKEN: xxxxxxxxx' \
		--header 'Cookie: XSRF-TOKEN=xxxxxxxxx;' \
		'http://localhost:8000/api/v1/aws/list-datasource/?project_id=1001&aws_id=1003'

.PHONY: attach-datasource
attach-datasource:
	curl -is -XPOST \
		--header 'x-amzn-oidc-identity: alice' \
		--header 'X-XSRF-TOKEN: xxxxxxxxx' \
		--header 'Cookie: XSRF-TOKEN=xxxxxxxxx;' \
		--header 'Content-Type: application/json' \
		--data '{"project_id":1001, "attach_data_source":{"aws_id":1003, "aws_data_source_id":1001, "project_id":1001, "assume_role_arn":"hoge-role", "external_id":""}}' \
		'http://localhost:8000/api/v1/aws/attach-datasource/'

.PHONY: detach-datasource
detach-datasource:
	curl -is -XPOST \
		--header 'x-amzn-oidc-identity: alice' \
		--header 'X-XSRF-TOKEN: xxxxxxxxx' \
		--header 'Cookie: XSRF-TOKEN=xxxxxxxxx;' \
		--header 'Content-Type: application/json' \
		--data '{"project_id":1001, "aws_id":1003, "aws_data_source_id":1001}' \
		'http://localhost:8000/api/v1/aws/detach-datasource/'

.PHONY: invoke-scan
invoke-scan:
	curl -is -XPOST \
		--header 'x-amzn-oidc-identity: alice' \
		--header 'X-XSRF-TOKEN: xxxxxxxxx' \
		--header 'Cookie: XSRF-TOKEN=xxxxxxxxx;' \
		--header 'Content-Type: application/json' \
		--data '{"project_id":1001, "aws_id":1003, "aws_data_source_id":1001}' \
		'http://localhost:8000/api/v1/aws/invoke-scan/'

.PHONY: list-osint
list-osint:
	curl -is -XGET \
		--header 'x-amzn-oidc-identity: alice' \
		--header 'X-XSRF-TOKEN: xxxxxxxxx' \
		--header 'Cookie: XSRF-TOKEN=xxxxxxxxx;' \
		'http://localhost:8000/api/v1/osint/list-osint/?project_id=1001'

.PHONY: get-osint
get-osint:
	curl -is -XGET \
		--header 'x-amzn-oidc-identity: alice' \
		--header 'X-XSRF-TOKEN: xxxxxxxxx' \
		--header 'Cookie: XSRF-TOKEN=xxxxxxxxx;' \
		'http://localhost:8000/api/v1/osint/get-osint/?project_id=1001&osint_id=1'

.PHONY: put-osint
put-osint:
	curl -is -XPOST \
		--header 'x-amzn-oidc-identity: alice' \
		--header 'X-XSRF-TOKEN: xxxxxxxxx' \
		--header 'Cookie: XSRF-TOKEN=xxxxxxxxx;' \
		--header 'Content-Type: application/json' \
		--data '{"project_id":1001, "osint":{"project_id":1001, "name":"hoge-osint"}}' \
		'http://localhost:8000/api/v1/osint/put-osint/'

.PHONY: delete-osint
delete-osint:
	curl -is -XPOST \
		--header 'x-amzn-oidc-identity: alice' \
		--header 'X-XSRF-TOKEN: xxxxxxxxx' \
		--header 'Cookie: XSRF-TOKEN=xxxxxxxxx;' \
		--header 'Content-Type: application/json' \
		--data '{"project_id":1001, "osint_id":3}' \
		'http://localhost:8000/api/v1/osint/delete-osint/'

.PHONY: list-osint_data_source
list-osint_data_source:
	curl -is -XGET \
		--header 'x-amzn-oidc-identity: alice' \
		--header 'X-XSRF-TOKEN: xxxxxxxxx' \
		--header 'Cookie: XSRF-TOKEN=xxxxxxxxx;' \
		'http://localhost:8000/api/v1/osint/list-datasource/?project_id=1001'

.PHONY: get-osint_data_source
get-osint_data_source:
	curl -is -XGET \
		--header 'x-amzn-oidc-identity: alice' \
		--header 'X-XSRF-TOKEN: xxxxxxxxx' \
		--header 'Cookie: XSRF-TOKEN=xxxxxxxxx;' \
		'http://localhost:8000/api/v1/osint/get-datasource/?project_id=1001&osint_data_source_id=1'

.PHONY: put-osint_data_source
put-osint_data_source:
	curl -is -XPOST \
		--header 'x-amzn-oidc-identity: alice' \
		--header 'X-XSRF-TOKEN: xxxxxxxxx' \
		--header 'Cookie: XSRF-TOKEN=xxxxxxxxx;' \
		--header 'Content-Type: application/json' \
		--data '{"project_id":1001, "osint_data_source":{"project_id":1001, "name":"hoge-osint", "description":"osint-datasource", "max_score": 10.0}}' \
		'http://localhost:8000/api/v1/osint/put-datasource/'

.PHONY: delete-osint_data_source
delete-osint_data_source:
	curl -is -XPOST \
		--header 'x-amzn-oidc-identity: alice' \
		--header 'X-XSRF-TOKEN: xxxxxxxxx' \
		--header 'Cookie: XSRF-TOKEN=xxxxxxxxx;' \
		--header 'Content-Type: application/json' \
		--data '{"project_id":1001, "osint_data_source_id":3}' \
		'http://localhost:8000/api/v1/osint/delete-datasource/'

.PHONY: list-osint_result
list-osint_result:
	curl -is -XGET \
		--header 'x-amzn-oidc-identity: alice' \
		--header 'X-XSRF-TOKEN: xxxxxxxxx' \
		--header 'Cookie: XSRF-TOKEN=xxxxxxxxx;' \
		'http://localhost:8000/api/v1/osint/list-rel-datasource/?project_id=1001'

.PHONY: get-osint_result
get-osint_result:
	curl -is -XGET \
		--header 'x-amzn-oidc-identity: alice' \
		--header 'X-XSRF-TOKEN: xxxxxxxxx' \
		--header 'Cookie: XSRF-TOKEN=xxxxxxxxx;' \
		'http://localhost:8000/api/v1/osint/get-rel-datasource/?project_id=1001&osint_result_id=1'

.PHONY: put-osint_result
put-osint_result:
	curl -is -XPOST \
		--header 'x-amzn-oidc-identity: alice' \
		--header 'X-XSRF-TOKEN: xxxxxxxxx' \
		--header 'Cookie: XSRF-TOKEN=xxxxxxxxx;' \
		--header 'Content-Type: application/json' \
		--data '{"project_id":1001, "osint_result":{"osint_id":1, "osint_data_source_id":1, "resource_type":"osint_resource_type", "resource_name":"osint_resource_name"}}' \
		'http://localhost:8000/api/v1/osint/put-rel-datasource/'

.PHONY: delete-osint_result
delete-osint_result:
	curl -is -XPOST \
		--header 'x-amzn-oidc-identity: alice' \
		--header 'X-XSRF-TOKEN: xxxxxxxxx' \
		--header 'Cookie: XSRF-TOKEN=xxxxxxxxx;' \
		--header 'Content-Type: application/json' \
		--data '{"project_id":1001, "osint_result_id":3}' \
		'http://localhost:8000/api/v1/osint/delete-rel-datasource/'

.PHONY: start-osint
start-osint:
	curl -is -XPOST \
		--header 'x-amzn-oidc-identity: alice' \
		--header 'X-XSRF-TOKEN: xxxxxxxxx' \
		--header 'Cookie: XSRF-TOKEN=xxxxxxxxx;' \
		--header 'Content-Type: application/json' \
		--data '{"project_id":1001, "osint_result_id":1}' \
		'http://localhost:8000/api/v1/osint/start-osint/'

.PHONY: list-diagnosis
list-diagnosis:
	curl -is -XGET \
		--header 'x-amzn-oidc-identity: alice' \
		--header 'X-XSRF-TOKEN: xxxxxxxxx' \
		--header 'Cookie: XSRF-TOKEN=xxxxxxxxx;' \
		'http://localhost:8000/api/v1/diagnosis/list-diagnosis/?project_id=1001'

.PHONY: get-diagnosis
get-diagnosis:
	curl -is -XGET \
		--header 'x-amzn-oidc-identity: alice' \
		--header 'X-XSRF-TOKEN: xxxxxxxxx' \
		--header 'Cookie: XSRF-TOKEN=xxxxxxxxx;' \
		'http://localhost:8000/api/v1/diagnosis/get-diagnosis/?project_id=1001&diagnosis_id=1'

.PHONY: put-diagnosis
put-diagnosis:
	curl -is -XPOST \
		--header 'x-amzn-oidc-identity: alice' \
		--header 'X-XSRF-TOKEN: xxxxxxxxx' \
		--header 'Cookie: XSRF-TOKEN=xxxxxxxxx;' \
		--header 'Content-Type: application/json' \
		--data '{"project_id":1001, "diagnosis":{"name":"hoge-diagnosis"}}' \
		'http://localhost:8000/api/v1/diagnosis/put-diagnosis/'

.PHONY: delete-diagnosis
delete-diagnosis:
	curl -is -XPOST \
		--header 'x-amzn-oidc-identity: alice' \
		--header 'X-XSRF-TOKEN: xxxxxxxxx' \
		--header 'Cookie: XSRF-TOKEN=xxxxxxxxx;' \
		--header 'Content-Type: application/json' \
		--data '{"project_id":1001, "diagnosis_id":3}' \
		'http://localhost:8000/api/v1/diagnosis/delete-diagnosis/'

.PHONY: list-diagnosis_data_source
list-diagnosis_data_source:
	curl -is -XGET \
		--header 'x-amzn-oidc-identity: alice' \
		--header 'X-XSRF-TOKEN: xxxxxxxxx' \
		--header 'Cookie: XSRF-TOKEN=xxxxxxxxx;' \
		'http://localhost:8000/api/v1/diagnosis/list-datasource/?project_id=1001'

.PHONY: get-diagnosis_data_source
get-diagnosis_data_source:
	curl -is -XGET \
		--header 'x-amzn-oidc-identity: alice' \
		--header 'X-XSRF-TOKEN: xxxxxxxxx' \
		--header 'Cookie: XSRF-TOKEN=xxxxxxxxx;' \
		'http://localhost:8000/api/v1/diagnosis/get-datasource/?project_id=1001&diagnosis_data_source_id=2'

.PHONY: put-diagnosis_data_source
put-diagnosis_data_source:
	curl -is -XPOST \
		--header 'x-amzn-oidc-identity: alice' \
		--header 'X-XSRF-TOKEN: xxxxxxxxx' \
		--header 'Cookie: XSRF-TOKEN=xxxxxxxxx;' \
		--header 'Content-Type: application/json' \
		--data '{"project_id":1001, "diagnosis_data_source":{"name":"hoge-diagnosis", "description":"diagnosis-datasource", "max_score": 10.0}}' \
		'http://localhost:8000/api/v1/diagnosis/put-datasource/'

.PHONY: delete-diagnosis_data_source
delete-diagnosis_data_source:
	curl -is -XPOST \
		--header 'x-amzn-oidc-identity: alice' \
		--header 'X-XSRF-TOKEN: xxxxxxxxx' \
		--header 'Cookie: XSRF-TOKEN=xxxxxxxxx;' \
		--header 'Content-Type: application/json' \
		--data '{"project_id":1001, "diagnosis_data_source_id":3}' \
		'http://localhost:8000/api/v1/diagnosis/delete-datasource/'

.PHONY: list-rel_diagnosis_datasource
list-rel_diagnosis_datasource:
	curl -is -XGET \
		--header 'x-amzn-oidc-identity: alice' \
		--header 'X-XSRF-TOKEN: xxxxxxxxx' \
		--header 'Cookie: XSRF-TOKEN=xxxxxxxxx;' \
		'http://localhost:8000/api/v1/diagnosis/list-rel-datasource/?project_id=1001'

.PHONY: get-rel_diagnosis_datasource
get-rel_diagnosis_datasource:
	curl -is -XGET \
		--header 'x-amzn-oidc-identity: alice' \
		--header 'X-XSRF-TOKEN: xxxxxxxxx' \
		--header 'Cookie: XSRF-TOKEN=xxxxxxxxx;' \
		'http://localhost:8000/api/v1/diagnosis/get-rel-datasource/?project_id=1001&rel_diagnosis_data_source_id=1'

.PHONY: put-rel_diagnosis_datasource
put-rel_diagnosis_datasource:
	curl -is -XPOST \
		--header 'x-amzn-oidc-identity: alice' \
		--header 'X-XSRF-TOKEN: xxxxxxxxx' \
		--header 'Cookie: XSRF-TOKEN=xxxxxxxxx;' \
		--header 'Content-Type: application/json' \
		--data '{"project_id":1001, "rel_diagnosis_data_source":{"diagnosis_id":1, "diagnosis_data_source_id":1, "record_id":"diagnosis_record_id", "jira_id":"diagnosis_jira_id", "jira_key":"diagnosis_jira_key"}}' \
		'http://localhost:8000/api/v1/diagnosis/put-rel-datasource/'

.PHONY: delete-rel_diagnosis_datasource
delete-rel_diagnosis_datasource:
	curl -is -XPOST \
		--header 'x-amzn-oidc-identity: alice' \
		--header 'X-XSRF-TOKEN: xxxxxxxxx' \
		--header 'Cookie: XSRF-TOKEN=xxxxxxxxx;' \
		--header 'Content-Type: application/json' \
		--data '{"project_id":1001, "rel_diagnosis_data_source_id":3}' \
		'http://localhost:8000/api/v1/diagnosis/delete-rel-datasource/'

.PHONY: start-diagnosis
start-diagnosis:
	curl -is -XPOST \
		--header 'x-amzn-oidc-identity: alice' \
		--header 'X-XSRF-TOKEN: xxxxxxxxx' \
		--header 'Cookie: XSRF-TOKEN=xxxxxxxxx;' \
		--header 'Content-Type: application/json' \
		--data '{"project_id":1001, "rel_diagnosis_data_source_id":1}' \
		'http://localhost:8000/api/v1/diagnosis/start-diagnosis/'
