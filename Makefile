TARGETS = gateway
BUILD_TARGETS = $(TARGETS:=.build)
BUILD_CI_TARGETS = $(TARGETS:=.build-ci)
IMAGE_PUSH_TARGETS = $(TARGETS:=.push-image)
MANIFEST_CREATE_TARGETS = $(TARGETS:=.create-manifest)
MANIFEST_PUSH_TARGETS = $(TARGETS:=.push-manifest)
TEST_TARGETS = $(TARGETS:=.go-test)
BUILD_OPT=""
IMAGE_TAG=latest
MANIFEST_TAG=latest
IMAGE_PREFIX=gateway
IMAGE_REGISTRY=local

.PHONY: all
all: build

.PHONY: go-mod-tidy
go-mod-tidy:
	go mod tidy

.PHONY: go-mod-update
go-mod-update:
	go get -u \
			github.com/ca-risken/core/... \
			github.com/ca-risken/aws/... \
			github.com/CyberAgent/mimosa-diagnosis/... \
			github.com/ca-risken/code/... \
			github.com/ca-risken/google/...

.PHONY: doc
doc: go-test
	. env.sh && ls *.go | grep -v '_test.go' | xargs go run

PHONY: build $(BUILD_TARGETS)
build: $(BUILD_TARGETS)
%.build: %.go-test
	. env.sh && TARGET=$(*) IMAGE_TAG=$(IMAGE_TAG) IMAGE_PREFIX=$(IMAGE_PREFIX) BUILD_OPT="$(BUILD_OPT)" . hack/docker-build.sh

PHONY: build-ci $(BUILD_CI_TARGETS)
build-ci: $(BUILD_CI_TARGETS)
%.build-ci:
	TARGET=$(*) IMAGE_TAG=$(IMAGE_TAG) IMAGE_PREFIX=$(IMAGE_PREFIX) BUILD_OPT="$(BUILD_OPT)" . hack/docker-build.sh
	docker tag $(IMAGE_PREFIX)/$(*):$(IMAGE_TAG) $(IMAGE_REGISTRY)/$(IMAGE_PREFIX)/$(*):$(IMAGE_TAG)

PHONY: push-image $(IMAGE_PUSH_TARGETS)
push-image: $(IMAGE_PUSH_TARGETS)
%.push-image:
	docker push $(IMAGE_REGISTRY)/$(IMAGE_PREFIX)/$(*):$(IMAGE_TAG)

PHONY: create-manifest $(MANIFEST_CREATE_TARGETS)
create-manifest: $(MANIFEST_CREATE_TARGETS)
%.create-manifest:
	docker manifest create $(IMAGE_REGISTRY)/$(IMAGE_PREFIX)/$(*):$(MANIFEST_TAG) $(IMAGE_REGISTRY)/$(IMAGE_PREFIX)/$(*):$(IMAGE_TAG_BASE)_linux_amd64 $(IMAGE_REGISTRY)/$(IMAGE_PREFIX)/$(*):$(IMAGE_TAG_BASE)_linux_arm64
	docker manifest annotate --arch amd64 $(IMAGE_REGISTRY)/$(IMAGE_PREFIX)/$(*):$(MANIFEST_TAG) $(IMAGE_REGISTRY)/$(IMAGE_PREFIX)/$(*):$(IMAGE_TAG_BASE)_linux_amd64
	docker manifest annotate --arch arm64 $(IMAGE_REGISTRY)/$(IMAGE_PREFIX)/$(*):$(MANIFEST_TAG) $(IMAGE_REGISTRY)/$(IMAGE_PREFIX)/$(*):$(IMAGE_TAG_BASE)_linux_arm64

PHONY: push-manifest $(MANIFEST_PUSH_TARGETS)
push-manifest: $(MANIFEST_PUSH_TARGETS)
%.push-manifest:
	docker manifest push $(IMAGE_REGISTRY)/$(IMAGE_PREFIX)/$(*):$(MANIFEST_TAG)
	docker manifest inspect $(IMAGE_REGISTRY)/$(IMAGE_PREFIX)/$(*):$(MANIFEST_TAG)

PHONY: go-test $(TEST_TARGETS)
go-test: $(TEST_TARGETS)
%.go-test:
	go test ./...

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

.PHONY: get-pend-finding
get-pend-finding:
	curl -is -XGET \
		--header 'x-amzn-oidc-identity: alice' \
		--header 'X-XSRF-TOKEN: xxxxxxxxx' \
		--header 'Cookie: XSRF-TOKEN=xxxxxxxxx;' \
		'http://localhost:8000/api/v1/finding/get-pend-finding/?project_id=1001&finding_id=1001'

.PHONY: put-pend-finding
put-pend-finding:
	curl -is -XPOST \
		--header 'x-amzn-oidc-identity: alice' \
		--header 'X-XSRF-TOKEN: xxxxxxxxx' \
		--header 'Cookie: XSRF-TOKEN=xxxxxxxxx;' \
		--header 'Content-Type: application/json' \
		--data '{"project_id":1001, "pend_finding":{"finding_id":1001, "project_id":1001}}' \
		'http://localhost:8000/api/v1/finding/put-pend-finding/'

.PHONY: delete-pend-finding
delete-pend-finding:
	curl -is -XPOST \
		--header 'x-amzn-oidc-identity: alice' \
		--header 'X-XSRF-TOKEN: xxxxxxxxx' \
		--header 'Cookie: XSRF-TOKEN=xxxxxxxxx;' \
		--header 'Content-Type: application/json' \
		--data '{"project_id":1001, "finding_id":1001}' \
		'http://localhost:8000/api/v1/finding/delete-pend-finding/'

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
		'http://localhost:8000/api/v1/iam/list-user/?activated=true'

.PHONY: get-user
get-user:
	curl -is -XGET \
		--header 'x-amzn-oidc-identity: alice' \
		--header 'X-XSRF-TOKEN: xxxxxxxxx' \
		--header 'Cookie: XSRF-TOKEN=xxxxxxxxx;' \
		'http://localhost:8000/api/v1/iam/get-user/?user_id=1001'

.PHONY: is-admin
is-admin:
	curl -is -XGET \
		--header 'x-amzn-oidc-identity: alice' \
		--header 'X-XSRF-TOKEN: xxxxxxxxx' \
		--header 'Cookie: XSRF-TOKEN=xxxxxxxxx;' \
		'http://localhost:8000/api/v1/iam/is-admin/?user_id=1001'

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

.PHONY: list-access-token
list-access-token:
	curl -is -XGET \
		--header 'x-amzn-oidc-identity: alice' \
		--header 'X-XSRF-TOKEN: xxxxxxxxx' \
		--header 'Cookie: XSRF-TOKEN=xxxxxxxxx;' \
		'http://localhost:8000/api/v1/iam/list-access-token/?project_id=1001'

.PHONY: generate-access-token
generate-access-token:
	curl -is -XPOST \
		--header 'x-amzn-oidc-identity: alice' \
		--header 'X-XSRF-TOKEN: xxxxxxxxx' \
		--header 'Cookie: XSRF-TOKEN=xxxxxxxxx;' \
		--header 'Content-Type: application/json' \
		--data '{"project_id":1001, "access_token":{"plain_text_token":"plain_text_token", "project_id":1001, "name":"curl", "description":"test for curl"}}' \
		'http://localhost:8000/api/v1/iam/generate-access-token/'

.PHONY: update-access-token
update-access-token:
	curl -is -XPOST \
		--header 'x-amzn-oidc-identity: alice' \
		--header 'X-XSRF-TOKEN: xxxxxxxxx' \
		--header 'Cookie: XSRF-TOKEN=xxxxxxxxx;' \
		--header 'Content-Type: application/json' \
		--data '{"project_id":1001, "access_token":{"access_token_id":1010, "project_id":1001, "name":"curl", "description":"updated"}}' \
		'http://localhost:8000/api/v1/iam/update-access-token/'

.PHONY: delete-access-token
delete-access-token:
	curl -is -XPOST \
		--header 'x-amzn-oidc-identity: alice' \
		--header 'X-XSRF-TOKEN: xxxxxxxxx' \
		--header 'Cookie: XSRF-TOKEN=xxxxxxxxx;' \
		--header 'Content-Type: application/json' \
		--data '{"project_id":1001, "access_token_id":1009}' \
		'http://localhost:8000/api/v1/iam/delete-access-token/'

.PHONY: attach-access-token
attach-access-token:
	curl -is -XPOST \
		--header 'x-amzn-oidc-identity: alice' \
		--header 'X-XSRF-TOKEN: xxxxxxxxx' \
		--header 'Cookie: XSRF-TOKEN=xxxxxxxxx;' \
		--header 'Content-Type: application/json' \
		--data '{"project_id":1001, "role_id":1001, "access_token_id":1010}' \
		'http://localhost:8000/api/v1/iam/attach-access-token/'

.PHONY: detach-access-token
detach-access-token:
	curl -is -XPOST \
		--header 'x-amzn-oidc-identity: alice' \
		--header 'X-XSRF-TOKEN: xxxxxxxxx' \
		--header 'Cookie: XSRF-TOKEN=xxxxxxxxx;' \
		--header 'Content-Type: application/json' \
		--data '{"project_id":1001, "role_id":1001, "access_token_id":1010}' \
		'http://localhost:8000/api/v1/iam/detach-access-token/'

# access_token_id=1010
.PHONY: test-token-access
test-token-access:
	curl -is -XGET \
		--header 'Authorization: BEARER MTAwMUAxMDEwQHJKNmtpUkVsLWJmaUVlZFRkT0hmRVhtendjTmI2akhOUFJJT19ZRjVVRGZEU1dmZkdzcnRnSm5zVjBMT1RzWkY2T2FsVjBzdnNIZVNvZ2JWOFptbnFn' \
		'http://localhost:8000/api/v1/iam/list-access-token/?project_id=1001'

.PHONY: list-project
list-project:
	curl -is -XGET \
		--header 'x-amzn-oidc-identity: alice' \
		--header 'X-XSRF-TOKEN: xxxxxxxxx' \
		--header 'Cookie: XSRF-TOKEN=xxxxxxxxx;' \
		'http://localhost:8000/api/v1/project/list-project/'

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
		--data '{"project_id":1002, "name":"project-b2"}' \
		'http://localhost:8000/api/v1/project/update-project/'

.PHONY: delete-project
delete-project:
	curl -is -XPOST \
		--header 'x-amzn-oidc-identity: alice' \
		--header 'X-XSRF-TOKEN: xxxxxxxxx' \
		--header 'Cookie: XSRF-TOKEN=xxxxxxxxx;' \
		--header 'Content-Type: application/json' \
		--data '{"project_id":1001}' \
		'http://localhost:8000/api/v1/project/delete-project/'

.PHONY: tag-project
tag-project:
	curl -is -XPOST \
		--header 'x-amzn-oidc-identity: alice' \
		--header 'X-XSRF-TOKEN: xxxxxxxxx' \
		--header 'Cookie: XSRF-TOKEN=xxxxxxxxx;' \
		--header 'Content-Type: application/json' \
		--data '{"project_id":1001, "tag": "test:true"}' \
		'http://localhost:8000/api/v1/project/tag-project/'

.PHONY: untag-project
untag-project:
	curl -is -XPOST \
		--header 'x-amzn-oidc-identity: alice' \
		--header 'X-XSRF-TOKEN: xxxxxxxxx' \
		--header 'Cookie: XSRF-TOKEN=xxxxxxxxx;' \
		--header 'Content-Type: application/json' \
		--data '{"project_id":1001, "tag": "test:true"}' \
		'http://localhost:8000/api/v1/project/untag-project/'

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

.PHONY: get-report
get-report:
	curl -is -XGET \
		--header 'x-amzn-oidc-identity: alice' \
		--header 'X-XSRF-TOKEN: xxxxxxxxx' \
		--header 'Cookie: XSRF-TOKEN=xxxxxxxxx;' \
		'http://localhost:8000/api/v1/report/get-report/?project_id=1001'

get-report-all:
	curl -is -XGET \
		--header 'x-amzn-oidc-identity: alice' \
		--header 'X-XSRF-TOKEN: xxxxxxxxx' \
		--header 'Cookie: XSRF-TOKEN=xxxxxxxxx;' \
		'http://localhost:8000/api/v1/report/get-report-all/?project_id=1001'

.PHONY: list-aws
list-aws:
	curl -is -XGET \
		--header 'x-amzn-oidc-identity: bob' \
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

.PHONY: describe-arn
describe-arn:
	curl -is -XGET \
		--header 'x-amzn-oidc-identity: alice' \
		--header 'X-XSRF-TOKEN: xxxxxxxxx' \
		--header 'Cookie: XSRF-TOKEN=xxxxxxxxx;' \
		'http://localhost:8000/api/v1/aws/describe-arn/?arn=arn:aws:s3:::bucket_name'

.PHONY: list-cloudtrail
list-cloudtrail:
	curl -is -XGET \
		--header 'x-amzn-oidc-identity: alice' \
		--header 'X-XSRF-TOKEN: xxxxxxxxx' \
		--header 'Cookie: XSRF-TOKEN=xxxxxxxxx;' \
		'http://localhost:8000/api/v1/aws/list-cloudtrail/?project_id=1001&aws_id=1002&region=ap-northeast-1&resource_type=AWS::S3::Bucket&resource_name=risken-public-test.security-hub.jp'

.PHONY: list-config-history
list-config-history:
	curl -is -XGET \
		--header 'x-amzn-oidc-identity: alice' \
		--header 'X-XSRF-TOKEN: xxxxxxxxx' \
		--header 'Cookie: XSRF-TOKEN=xxxxxxxxx;' \
		'http://localhost:8000/api/v1/aws/list-config-history/?project_id=1001&aws_id=1002&region=ap-northeast-1&resource_type=AWS::S3::Bucket&resource_name=risken-public-test.security-hub.jp'

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
		'http://localhost:8000/api/v1/osint/get-osint/?project_id=1001&osint_id=1001'

.PHONY: put-osint
put-osint:
	curl -is -XPOST \
		--header 'x-amzn-oidc-identity: alice' \
		--header 'X-XSRF-TOKEN: xxxxxxxxx' \
		--header 'Cookie: XSRF-TOKEN=xxxxxxxxx;' \
		--header 'Content-Type: application/json' \
		--data '{"project_id":1001, "osint":{"project_id":1001, "resource_type":"Domain", "resource_name":"example.com"}}' \
		'http://localhost:8000/api/v1/osint/put-osint/'

.PHONY: delete-osint
delete-osint:
	curl -is -XPOST \
		--header 'x-amzn-oidc-identity: alice' \
		--header 'X-XSRF-TOKEN: xxxxxxxxx' \
		--header 'Cookie: XSRF-TOKEN=xxxxxxxxx;' \
		--header 'Content-Type: application/json' \
		--data '{"project_id":1001, "osint_id":1003}' \
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
		'http://localhost:8000/api/v1/osint/get-datasource/?project_id=1001&osint_data_source_id=1001'

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
		--data '{"project_id":1001, "osint_data_source_id":1002}' \
		'http://localhost:8000/api/v1/osint/delete-datasource/'

.PHONY: list-rel_osint_data_source
list-rel_osint_data_source:
	curl -is -XGET \
		--header 'x-amzn-oidc-identity: alice' \
		--header 'X-XSRF-TOKEN: xxxxxxxxx' \
		--header 'Cookie: XSRF-TOKEN=xxxxxxxxx;' \
		'http://localhost:8000/api/v1/osint/list-rel-datasource/?project_id=1001'

.PHONY: get-rel_osint_data_source
get-rel_osint_data_source:
	curl -is -XGET \
		--header 'x-amzn-oidc-identity: alice' \
		--header 'X-XSRF-TOKEN: xxxxxxxxx' \
		--header 'Cookie: XSRF-TOKEN=xxxxxxxxx;' \
		'http://localhost:8000/api/v1/osint/get-rel-datasource/?project_id=1001&rel_osint_data_source_id=1001'

.PHONY: put-rel_osint_data_source
put-rel_osint_data_source:
	curl -is -XPOST \
		--header 'x-amzn-oidc-identity: alice' \
		--header 'X-XSRF-TOKEN: xxxxxxxxx' \
		--header 'Cookie: XSRF-TOKEN=xxxxxxxxx;' \
		--header 'Content-Type: application/json' \
		--data '{"project_id":1001, "rel_osint_data_source":{"osint_id":1001, "osint_data_source_id":1001, "project_id":1001, "status": 1}}' \
		'http://localhost:8000/api/v1/osint/put-rel-datasource/'

.PHONY: delete-rel_osint_data_source
delete-rel_osint_data_source:
	curl -is -XPOST \
		--header 'x-amzn-oidc-identity: alice' \
		--header 'X-XSRF-TOKEN: xxxxxxxxx' \
		--header 'Cookie: XSRF-TOKEN=xxxxxxxxx;' \
		--header 'Content-Type: application/json' \
		--data '{"project_id":1001, "rel_osint_data_source_id":1003}' \
		'http://localhost:8000/api/v1/osint/delete-rel-datasource/'

.PHONY: list-osint_detect_word
list-osint_detect_word:
	curl -is -XGET \
		--header 'x-amzn-oidc-identity: alice' \
		--header 'X-XSRF-TOKEN: xxxxxxxxx' \
		--header 'Cookie: XSRF-TOKEN=xxxxxxxxx;' \
		'http://localhost:8000/api/v1/osint/list-word/?project_id=1001'

.PHONY: get-osint_detect_word
get-osint_detect_word:
	curl -is -XGET \
		--header 'x-amzn-oidc-identity: alice' \
		--header 'X-XSRF-TOKEN: xxxxxxxxx' \
		--header 'Cookie: XSRF-TOKEN=xxxxxxxxx;' \
		'http://localhost:8000/api/v1/osint/get-word/?project_id=1001&osint_detect_word_id=1001'

.PHONY: put-osint_detect_word
put-osint_detect_word:
	curl -is -XPOST \
		--header 'x-amzn-oidc-identity: alice' \
		--header 'X-XSRF-TOKEN: xxxxxxxxx' \
		--header 'Cookie: XSRF-TOKEN=xxxxxxxxx;' \
		--header 'Content-Type: application/json' \
		--data '{"project_id":1001, "osint_detect_word":{"rel_osint_data_source_id":1001, "word":"hoge", "project_id":1001}}' \
		'http://localhost:8000/api/v1/osint/put-word/'

.PHONY: delete-osint_detect_word
delete-osint_detect_word:
	curl -is -XPOST \
		--header 'x-amzn-oidc-identity: alice' \
		--header 'X-XSRF-TOKEN: xxxxxxxxx' \
		--header 'Cookie: XSRF-TOKEN=xxxxxxxxx;' \
		--header 'Content-Type: application/json' \
		--data '{"project_id":1001, "osint_detect_word_id":1003}' \
		'http://localhost:8000/api/v1/osint/delete-word/'


.PHONY: invoke-osint-scan
invoke-osint-scan:
	curl -is -XPOST \
		--header 'x-amzn-oidc-identity: alice' \
		--header 'X-XSRF-TOKEN: xxxxxxxxx' \
		--header 'Cookie: XSRF-TOKEN=xxxxxxxxx;' \
		--header 'Content-Type: application/json' \
		--data '{"project_id":1001, "rel_osint_data_source_id":1001}' \
		'http://localhost:8000/api/v1/osint/invoke-scan/'

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
		'http://localhost:8000/api/v1/diagnosis/get-datasource/?project_id=1001&diagnosis_data_source_id=1001'

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
		--data '{"project_id":1001, "diagnosis_data_source_id":1001}' \
		'http://localhost:8000/api/v1/diagnosis/delete-datasource/'

.PHONY: list-jira_setting
list-jira_setting:
	curl -is -XGET \
		--header 'x-amzn-oidc-identity: alice' \
		--header 'X-XSRF-TOKEN: xxxxxxxxx' \
		--header 'Cookie: XSRF-TOKEN=xxxxxxxxx;' \
		'http://localhost:8000/api/v1/diagnosis/list-jira-setting/?project_id=1001'

.PHONY: get-jira_setting
get-jira_setting:
	curl -is -XGET \
		--header 'x-amzn-oidc-identity: alice' \
		--header 'X-XSRF-TOKEN: xxxxxxxxx' \
		--header 'Cookie: XSRF-TOKEN=xxxxxxxxx;' \
		'http://localhost:8000/api/v1/diagnosis/get-jira-setting/?project_id=1001&jira_setting_id=1001'

.PHONY: put-jira_setting
put-jira_setting:
	curl -is -XPOST \
		--header 'x-amzn-oidc-identity: alice' \
		--header 'X-XSRF-TOKEN: xxxxxxxxx' \
		--header 'Cookie: XSRF-TOKEN=xxxxxxxxx;' \
		--header 'Content-Type: application/json' \
		--data '{"project_id":1001, "jira_setting":{"name":"name", "project_id":1001,"diagnosis_data_source_id":1001, "identity_field":"diagnosis_identity_field", "identity_value":"diagnosis_identity_value", "jira_id":"diagnosis_jira_id", "jira_key":"diagnosis_jira_key"}}' \
		'http://localhost:8000/api/v1/diagnosis/put-jira-setting/'

.PHONY: delete-jira_setting
delete-jira_setting:
	curl -is -XPOST \
		--header 'x-amzn-oidc-identity: alice' \
		--header 'X-XSRF-TOKEN: xxxxxxxxx' \
		--header 'Cookie: XSRF-TOKEN=xxxxxxxxx;' \
		--header 'Content-Type: application/json' \
		--data '{"project_id":1001, "jira_setting_id":1003}' \
		'http://localhost:8000/api/v1/diagnosis/delete-jira-setting/'

.PHONY: invoke-diagnosis-scan
invoke-diagnosis-scan:
	curl -is -XPOST \
		--header 'x-amzn-oidc-identity: alice' \
		--header 'X-XSRF-TOKEN: xxxxxxxxx' \
		--header 'Cookie: XSRF-TOKEN=xxxxxxxxx;' \
		--header 'Content-Type: application/json' \
		--data '{"project_id":1001, "jira_setting_id":1001}' \
		'http://localhost:8000/api/v1/diagnosis/invoke-scan/'

.PHONY: list-code-datasource
list-code-datasource:
	curl -is -XGET \
		--header 'x-amzn-oidc-identity: alice' \
		--header 'X-XSRF-TOKEN: xxxxxxxxx' \
		--header 'Cookie: XSRF-TOKEN=xxxxxxxxx;' \
		'http://localhost:8000/api/v1/code/list-datasource/?project_id=1001'

.PHONY: list-gitleaks
list-gitleaks:
	curl -is -XGET \
		--header 'x-amzn-oidc-identity: alice' \
		--header 'X-XSRF-TOKEN: xxxxxxxxx' \
		--header 'Cookie: XSRF-TOKEN=xxxxxxxxx;' \
		'http://localhost:8000/api/v1/code/list-gitleaks/?project_id=1001'

.PHONY: put-gitleaks
put-gitleaks:
	curl -is -XPOST \
		--header 'x-amzn-oidc-identity: alice' \
		--header 'X-XSRF-TOKEN: xxxxxxxxx' \
		--header 'Cookie: XSRF-TOKEN=xxxxxxxxx;' \
		--header 'Content-Type: application/json' \
		--data '{"project_id":1001, "gitleaks": {"gitleaks_id":1001, "code_data_source_id":1001, "name":"test-gitleaks", "project_id":1001, "type":2, "target_resource":"gitleakstest"}}' \
		'http://localhost:8000/api/v1/code/put-gitleaks/'

.PHONY: delete-gitleaks
delete-gitleaks:
	curl -is -XPOST \
		--header 'x-amzn-oidc-identity: alice' \
		--header 'X-XSRF-TOKEN: xxxxxxxxx' \
		--header 'Cookie: XSRF-TOKEN=xxxxxxxxx;' \
		--header 'Content-Type: application/json' \
		--data '{"project_id":1001, "gitleaks_id":1001}' \
		'http://localhost:8000/api/v1/code/delete-gitleaks/'

.PHONY: invoke-scan-gitleaks
invoke-scan-gitleaks:
	curl -is -XPOST \
		--header 'x-amzn-oidc-identity: alice' \
		--header 'X-XSRF-TOKEN: xxxxxxxxx' \
		--header 'Cookie: XSRF-TOKEN=xxxxxxxxx;' \
		--header 'Content-Type: application/json' \
		--data '{"project_id":1001, "gitleaks_id":1001}' \
		'http://localhost:8000/api/v1/code/invoke-scan-gitleaks/'

.PHONY: list-google-datasource
list-google-datasource:
	curl -is -XGET \
		--header 'x-amzn-oidc-identity: alice' \
		--header 'X-XSRF-TOKEN: xxxxxxxxx' \
		--header 'Cookie: XSRF-TOKEN=xxxxxxxxx;' \
		'http://localhost:8000/api/v1/google/list-google-datasource/?project_id=1001'

.PHONY: list-gcp
list-gcp:
	curl -is -XGET \
		--header 'x-amzn-oidc-identity: alice' \
		--header 'X-XSRF-TOKEN: xxxxxxxxx' \
		--header 'Cookie: XSRF-TOKEN=xxxxxxxxx;' \
		'http://localhost:8000/api/v1/google/list-gcp/?project_id=1001'

.PHONY: get-gcp
get-gcp:
	curl -is -XGET \
		--header 'x-amzn-oidc-identity: alice' \
		--header 'X-XSRF-TOKEN: xxxxxxxxx' \
		--header 'Cookie: XSRF-TOKEN=xxxxxxxxx;' \
		'http://localhost:8000/api/v1/google/get-gcp/?project_id=1001&gcp_id=1001'

.PHONY: put-gcp
put-gcp:
	curl -is -XPOST \
		--header 'x-amzn-oidc-identity: alice' \
		--header 'X-XSRF-TOKEN: xxxxxxxxx' \
		--header 'Cookie: XSRF-TOKEN=xxxxxxxxx;' \
		--header 'Content-Type: application/json' \
		--data '{"project_id":1001, "gcp": {"gcp_id":1002, "name":"test", "project_id":1001, "gcp_project_id":"test"}}' \
		'http://localhost:8000/api/v1/google/put-gcp/'

.PHONY: delete-gcp
delete-gcp:
	curl -is -XPOST \
		--header 'x-amzn-oidc-identity: alice' \
		--header 'X-XSRF-TOKEN: xxxxxxxxx' \
		--header 'Cookie: XSRF-TOKEN=xxxxxxxxx;' \
		--header 'Content-Type: application/json' \
		--data '{"project_id":1001, "gcp_id":1002}' \
		'http://localhost:8000/api/v1/google/delete-gcp/'

.PHONY: invoke-scan-gcp
invoke-scan-gcp:
	curl -is -XPOST \
		--header 'x-amzn-oidc-identity: alice' \
		--header 'X-XSRF-TOKEN: xxxxxxxxx' \
		--header 'Cookie: XSRF-TOKEN=xxxxxxxxx;' \
		--header 'Content-Type: application/json' \
		--data '{"project_id":1001, "gcp_id":1001, "google_data_source_id":1001}' \
		'http://localhost:8000/api/v1/google/invoke-scan-gcp/'

.PHONY: list-finding-setting
list-finding-setting:
	curl -is -XGET \
		--header 'x-amzn-oidc-identity: alice' \
		--header 'X-XSRF-TOKEN: xxxxxxxxx' \
		--header 'Cookie: XSRF-TOKEN=xxxxxxxxx;' \
		'http://localhost:8000/api/v1/finding/list-finding-setting/?project_id=1001&status=1'

.PHONY: get-finding-setting
get-finding-setting:
	curl -is -XGET \
		--header 'x-amzn-oidc-identity: alice' \
		--header 'X-XSRF-TOKEN: xxxxxxxxx' \
		--header 'Cookie: XSRF-TOKEN=xxxxxxxxx;' \
		'http://localhost:8000/api/v1/finding/get-finding-setting/?project_id=1001&finding_setting_id=1001'

.PHONY: put-finding-setting
put-finding-setting:
	curl -is -XPOST \
		--header 'x-amzn-oidc-identity: alice' \
		--header 'X-XSRF-TOKEN: xxxxxxxxx' \
		--header 'Cookie: XSRF-TOKEN=xxxxxxxxx;' \
		--header 'Content-Type: application/json' \
		--data '{"project_id":1001, "finding_setting":{"project_id":1001, "resource_name":"RN", "status":1, "setting":"{}"}}' \
		'http://localhost:8000/api/v1/finding/put-finding-setting/'

.PHONY: delete-finding-setting
delete-finding-setting:
	curl -is -XPOST \
		--header 'x-amzn-oidc-identity: alice' \
		--header 'X-XSRF-TOKEN: xxxxxxxxx' \
		--header 'Cookie: XSRF-TOKEN=xxxxxxxxx;' \
		--header 'Content-Type: application/json' \
		--data '{"project_id":1001, "finding_setting_id":1001}' \
		'http://localhost:8000/api/v1/finding/delete-finding-setting/'
