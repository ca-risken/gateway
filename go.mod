module githug.com/ca-risken/gateway

go 1.17

require (
	github.com/aws/aws-xray-sdk-go v1.6.0
	github.com/ca-risken/aws/proto/activity v0.0.0-20210909100807-3b4b7805f6fa
	github.com/ca-risken/aws/proto/aws v0.0.0-20210909100807-3b4b7805f6fa
	github.com/ca-risken/code/proto/code v0.0.0-20210910090912-21759a7829ac
	github.com/ca-risken/common/pkg/logging v0.0.0-20220113015330-0e8462d52b5b
	github.com/ca-risken/common/pkg/profiler v0.0.0-20220304031727-c94e2c463b27
	github.com/ca-risken/common/pkg/trace v0.0.0-20220208092527-a794114e4840
	github.com/ca-risken/common/pkg/xray v0.0.0-20211118071101-9855266b50a1
	github.com/ca-risken/core/proto/alert v0.0.0-20211124103540-8199b00d6465
	github.com/ca-risken/core/proto/finding v0.0.0-20211124103540-8199b00d6465
	github.com/ca-risken/core/proto/iam v0.0.0-20211124103540-8199b00d6465
	github.com/ca-risken/core/proto/project v0.0.0-20220317065728-a1d0959e3170
	github.com/ca-risken/core/proto/report v0.0.0-20211124103540-8199b00d6465
	github.com/ca-risken/diagnosis/proto/diagnosis v0.0.0-20220221070251-d51f9add7f73
	github.com/ca-risken/google/proto/google v0.0.0-20210909140724-ed07f4a6bd06
	github.com/ca-risken/osint/proto/osint v0.0.0-20210909141026-cc5c5b5d6de0
	github.com/gassara-kys/envconfig v1.4.4
	github.com/go-chi/chi/v5 v5.0.7
	github.com/go-chi/docgen v1.2.0
	github.com/golang/protobuf v1.5.2
	github.com/gorilla/schema v1.2.0
	github.com/grpc-ecosystem/go-grpc-middleware v1.2.2
	github.com/sirupsen/logrus v1.8.1
	github.com/stretchr/testify v1.7.0
	github.com/vikyd/zero v0.0.0-20190921142904-0f738d0bc858
	go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc v0.28.0
	google.golang.org/grpc v1.45.0
)

require (
	github.com/DataDog/datadog-go/v5 v5.0.2 // indirect
	github.com/DataDog/gostackparse v0.5.0 // indirect
	github.com/Microsoft/go-winio v0.5.1 // indirect
	github.com/andybalholm/brotli v1.0.3 // indirect
	github.com/asaskevich/govalidator v0.0.0-20210307081110-f21760c49a8d // indirect
	github.com/aws/aws-sdk-go v1.40.48 // indirect
	github.com/ca-risken/common/pkg/sqs v0.0.0-20220113015330-0e8462d52b5b // indirect
	github.com/cenkalti/backoff/v4 v4.1.2 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/envoyproxy/protoc-gen-validate v0.6.7 // indirect
	github.com/felixge/httpsnoop v1.0.2 // indirect
	github.com/gassara-kys/go-sqs-poller/worker/v4 v4.0.0-20210215110542-0be358599a2f // indirect
	github.com/go-logr/logr v1.2.2 // indirect
	github.com/go-logr/stdr v1.2.2 // indirect
	github.com/go-ozzo/ozzo-validation/v4 v4.3.0 // indirect
	github.com/google/pprof v0.0.0-20210423192551-a2663126120b // indirect
	github.com/google/uuid v1.3.0 // indirect
	github.com/grpc-ecosystem/grpc-gateway v1.16.0 // indirect
	github.com/jmespath/go-jmespath v0.4.0 // indirect
	github.com/klauspost/compress v1.13.5 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/stretchr/objx v0.1.1 // indirect
	github.com/valyala/bytebufferpool v1.0.0 // indirect
	github.com/valyala/fasthttp v1.30.0 // indirect
	go.opentelemetry.io/contrib v1.3.0 // indirect
	go.opentelemetry.io/otel v1.3.0 // indirect
	go.opentelemetry.io/otel/exporters/otlp/internal/retry v1.3.0 // indirect
	go.opentelemetry.io/otel/exporters/otlp/otlptrace v1.3.0 // indirect
	go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp v1.3.0 // indirect
	go.opentelemetry.io/otel/exporters/stdout/stdouttrace v1.3.0 // indirect
	go.opentelemetry.io/otel/sdk v1.3.0 // indirect
	go.opentelemetry.io/otel/trace v1.3.0 // indirect
	go.opentelemetry.io/proto/otlp v0.11.0 // indirect
	golang.org/x/net v0.0.0-20220403103023-749bd193bc2b // indirect
	golang.org/x/sys v0.0.0-20220405210540-1e041c57c461 // indirect
	golang.org/x/text v0.3.7 // indirect
	google.golang.org/genproto v0.0.0-20220405205423-9d709892a2bf // indirect
	google.golang.org/protobuf v1.28.0 // indirect
	gopkg.in/DataDog/dd-trace-go.v1 v1.36.0 // indirect
	gopkg.in/yaml.v3 v3.0.0-20200615113413-eeeca48fe776 // indirect
)
