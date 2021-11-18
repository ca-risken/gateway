module githug.com/ca-risken/gateway

go 1.16

require (
	github.com/andybalholm/brotli v1.0.3 // indirect
	github.com/aws/aws-xray-sdk-go v1.6.0
	github.com/ca-risken/aws/proto/activity v0.0.0-20210909100807-3b4b7805f6fa
	github.com/ca-risken/aws/proto/aws v0.0.0-20210909100807-3b4b7805f6fa
	github.com/ca-risken/code/proto/code v0.0.0-20210910090912-21759a7829ac
	github.com/ca-risken/common/pkg/xray v0.0.0-20211118071101-9855266b50a1
	github.com/ca-risken/core/proto/alert v0.0.0-20210917123127-86fbc1daa83f
	github.com/ca-risken/core/proto/finding v0.0.0-20210917123127-86fbc1daa83f
	github.com/ca-risken/core/proto/iam v0.0.0-20210917123127-86fbc1daa83f
	github.com/ca-risken/core/proto/project v0.0.0-20210917123127-86fbc1daa83f
	github.com/ca-risken/core/proto/report v0.0.0-20210917123127-86fbc1daa83f
	github.com/ca-risken/diagnosis/proto/diagnosis v0.0.0-20211015052841-3e9c68166866
	github.com/ca-risken/google/proto/google v0.0.0-20210909140724-ed07f4a6bd06
	github.com/ca-risken/osint/proto/osint v0.0.0-20210909141026-cc5c5b5d6de0
	github.com/gassara-kys/envconfig v1.4.4
	github.com/go-chi/chi/v5 v5.0.4
	github.com/go-chi/docgen v1.2.0
	github.com/golang/protobuf v1.5.2
	github.com/gorilla/schema v1.2.0
	github.com/klauspost/compress v1.13.5 // indirect
	github.com/sirupsen/logrus v1.8.1
	github.com/stretchr/testify v1.7.0
	github.com/valyala/fasthttp v1.30.0 // indirect
	github.com/vikyd/zero v0.0.0-20190921142904-0f738d0bc858
	golang.org/x/net v0.0.0-20211014222326-fd004c51d1d6 // indirect
	google.golang.org/grpc v1.41.0
)
