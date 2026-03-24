# AGENTS.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Commands

```bash
make install              # go mod download
make go-test              # Run all tests (go test ./...)
go test -run TestXxx .    # Run a single test
make lint                 # golangci-lint (requires golangci-lint installed)
make build                # Docker image build (runs tests first)
make generate-service     # Generate HTTP handlers from proto files (requires protoc + sibling repos)
make go-mod-update        # Update ca-risken/core and ca-risken/datasource-api dependencies
```

## Architecture

This is the **RISKEN API Gateway** — an HTTP-to-gRPC reverse proxy that sits in front of all RISKEN backend services. It handles authentication, authorization, CSRF protection, and request routing via a chi router middleware stack.

### Key Flow

```
HTTP Request → chi middleware (auth, CSRF, logging, Datadog tracing)
             → handler (bind request → gRPC call → write JSON response)
```

### Code Generation

Most handler code is **auto-generated** from proto files in sibling repositories (`ca-risken/core`, `ca-risken/datasource-api`). The custom `protoc-gen-service` plugin in `hack/` reads proto service definitions and generates `service_*_generated.go` files. **Do not manually edit `*_generated.go` files.**

To add a new API that simply proxies HTTP to gRPC:
1. Run `make generate-service`
2. Add routing in `router.go`

To add a custom handler (non-trivial logic):
1. Add the method to `hack/protoc-gen-service.yml` excludes list
2. Implement the handler manually
3. Add routing in `router.go`

### Core Files

- `main.go` — Entry point, config loading, gRPC client setup
- `router.go` — All HTTP routing and middleware wiring
- `service.go` — `gatewayService` struct holding all gRPC clients
- `authorizer.go` — Authentication (OIDC/JWT) and authorization (project/organization level)
- `claims.go` — JWT claims parsing from ALB OIDC headers
- `bind.go` — Request parameter binding (query params for GET, JSON body for POST/PUT/DELETE)
- `error.go` — gRPC-to-HTTP status code mapping
- `access_token.go` — Project/organization access token encode/decode

### Authentication Model

Two auth paths:
1. **User auth**: OIDC via ALB headers (`x-amzn-oidc-data`, `x-amzn-oidc-identity`) → JWT claims → `requestUser` in context
2. **Token auth**: Bearer token in Authorization header → project or organization access token

### Authorization Levels

- `authzWithProject` — Project-scoped authorization via `IAMService.IsAuthorized`
- `authzWithOrganization` — Organization-scoped authorization via `OrganizationIAMService.IsAuthorizedOrganization`

### Dependencies

- `ca-risken/core` — Proto definitions for core services (finding, alert, IAM, project, report, organization, AI)
- `ca-risken/datasource-api` — Proto definitions for datasource services (AWS, Azure, Google, Code, OSINT, Diagnosis)
- `go-chi/chi/v5` — HTTP router
- `golang-jwt/jwt/v4` — JWT processing
- Datadog (`dd-trace-go`) — APM tracing and profiling

## Routing Path Naming Convention

- When combining resource names in API paths, use hyphens (`organization-alert`) instead of slashes (`organization/alert`)
  - Reason: Based on ActionName specification. `getActionNameFromURI` generates ActionName as `{service}/{path1}` from `/api/v1/{service}/{path1}`, so the service name itself must be hyphen-joined
  - Example: `/api/v1/organization-alert/list-notification` → ActionName: `organization-alert/list-notification`

## Conventions

- Go 1.23.3
- All code lives in a single `main` package (no sub-packages)
- Generated files follow the pattern `service_<name>_generated.go`
- Custom (non-generated) service handlers are in `service_<name>.go`
- Tests use `httptest.NewRecorder` + `http.NewRequest` with mocked gRPC clients

## CLAUDE.md Placement Rule

Every directory that contains an AGENTS.md file must also have a CLAUDE.md file that includes `@AGENTS.md` to ensure the agent guidelines are loaded automatically.
