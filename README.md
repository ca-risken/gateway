# RISKEN Gateway

![Build Status](https://codebuild.ap-northeast-1.amazonaws.com/badges?uuid=eyJlbmNyeXB0ZWREYXRhIjoiVXhYVityeklNcTJBSEw5eVhIcUpvbVpFOGRoUDJiSTROMzJodmYxc2phRVI3dFA5Y2xnbFd6VUdUelczQms3YjFCNnEzYjczNDVCN250Wkd0T3dMdFVjPSIsIml2UGFyYW1ldGVyU3BlYyI6IjY5K3lhTUNERmF3SlBuYlciLCJtYXRlcmlhbFNldFNlcmlhbCI6MX0%3D&branch=master)

`RISKEN` is a monitoring tool for your cloud platforms, web-site, source-code... 
`RISKEN Gateway` is a API-Gateway layer for RISKEN backend APIs. Gateway is responsible for authentication & authorization control and request routing.

Please check [RISKEN Documentation](https://docs.security-hub.jp/).

## Installation

### Requirements

This module requires the following modules:

- [Go](https://go.dev/doc/install)
- [Docker](https://docs.docker.com/get-docker/)
- [Protocol Buffer](https://grpc.io/docs/protoc-installation/)

### Install packages

This module is developed in the `Go language`, please run the following command after installing the `Go`.

```bash
$ make install
```

### Building

Build the containers on your machine with the following command

```bash
$ make build
```

### Running Apps

Deploy the pre-built containers to the Kubernetes environment on your local machine.

- Follow the [documentation](https://docs.security-hub.jp/admin/infra_local/#risken) to download the Kubernetes manifest sample.
- Fix the Kubernetes object specs of the manifest file as follows and deploy it.

`k8s-sample/overlays/local/gateway.yaml`

| service | spec                                | before (public images)                         | after (pre-build images on your machine) |
| ------- | ----------------------------------- | ---------------------------------------------- | ---------------------------------------- |
| gateway | spec.template.spec.containers.image | `public.ecr.aws/risken/gateway/gateway:latest` | `gateway/gateway:latest`                 |

## Implements services

After change proto files in other repository, you must apply it to this repository. And you can do most of these easily.

### Case: only proxy HTTP to gRPC

Generate boilerplate code from proto files

```shell
make generate-service
```

Add routing code to router.go

```go
func newRouter(svc *gatewayService) *chi.Mux {
// ... 
r.Get("<routing path>", svc.<generated method>) // append
// ...
}
```

### Case: others, need most complex process

Add exclude configuration to hack/protoc-gen-service.yml

```yaml
excludes:
  # Ignore
  ## core
  - AlertService.AnalyzeAlertAll
  - FindingService.BatchListFinding
  - FindingService.PutFindingBatch
  - <ServiceName>.<MethodName> # append
```

Implements request handler and add it to router

## Remediation Proposal retrieval APIs

Both endpoints require authentication and project authorization. Query parameters are passed to the Core AI service.

| Method | Path | Query parameters |
| --- | --- | --- |
| GET | `/api/v1/ai/get-remediation-proposal` | Required: `project_id`, `remediation_proposal_id` |
| GET | `/api/v1/ai/list-remediation-proposal` | Required: `project_id`, `finding_id`. Optional: `status` (`PENDING`, `SUCCEEDED`, or `FAILED`) |

Get returns `{"data":{"remediation_proposal":{...}}}`. Each proposal includes its ID, Finding ID, project ID, status, status detail, remediation plan, and generation/creation/update timestamps when populated. The remediation plan is a string; timestamps are Unix seconds. If the proposal does not exist in the requested project, Get returns HTTP 200 with `{"data":{}}`.

List returns `{"data":{"remediation_proposals":[...]}}`, ordered by creation time descending. Omitting `status` returns all statuses. With no matches, the response is HTTP 200 with `{"data":{}}`; the empty repeated field is omitted by the existing JSON serializer.

Unauthenticated requests return HTTP 401, project authorization failures return HTTP 403, and invalid proposal/Finding IDs or statuses return HTTP 400 after authorization. Core errors use the existing gRPC-to-HTTP mapping.

## Community

Info on reporting bugs, getting help, finding roadmaps,
and more can be found in the [RISKEN Community](https://github.com/ca-risken/community).

## License

[MIT](LICENSE).
