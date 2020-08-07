module githug.com/CyberAgent/mimosa-gateway

go 1.13

require (
	github.com/CyberAgent/mimosa-aws/proto/aws v0.0.0-20200804164209-f0eae226d636
	github.com/CyberAgent/mimosa-core/proto/finding v0.0.0-20200807055414-eefbfe4897b2
	github.com/CyberAgent/mimosa-core/proto/iam v0.0.0-20200807055414-eefbfe4897b2
	github.com/CyberAgent/mimosa-core/proto/project v0.0.0-20200807055414-eefbfe4897b2
	github.com/CyberAgent/mimosa-gateway v0.0.0-20200807154329-5f5d52b7079a // indirect
	github.com/go-chi/chi v4.1.2+incompatible
	github.com/golang/protobuf v1.4.2
	github.com/gorilla/schema v1.1.0
	github.com/kelseyhightower/envconfig v1.4.0
	github.com/sirupsen/logrus v1.6.0
	github.com/stretchr/testify v1.6.1
	github.com/vikyd/zero v0.0.0-20190921142904-0f738d0bc858
	google.golang.org/grpc v1.31.0
)
