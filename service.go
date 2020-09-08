package main

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/CyberAgent/mimosa-aws/proto/aws"
	"github.com/CyberAgent/mimosa-core/proto/finding"
	"github.com/CyberAgent/mimosa-core/proto/iam"
	"github.com/CyberAgent/mimosa-core/proto/project"
	"github.com/CyberAgent/mimosa-osint/pkg/pb/osint"
	"github.com/kelseyhightower/envconfig"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

const (
	successJSONKey = "data"
	errorJSONKey   = "error"
)

type gatewayService struct {
	port          string
	uidHeader     string
	findingClient finding.FindingServiceClient
	iamClient     iam.IAMServiceClient
	projectClient project.ProjectServiceClient
	awsClient     aws.AWSServiceClient
	osintClient   osint.OSINTServiceClient
}

type gatewayConf struct {
	Port               string `default:"8000"`
	Debug              bool   `default:"false"`
	UserIdentityHeader string `required:"true" split_words:"true"`
	FindingSvcAddr     string `required:"true" split_words:"true"`
	IAMSvcAddr         string `required:"true" split_words:"true"`
	ProjectSvcAddr     string `required:"true" split_words:"true"`
	AWSSvcAddr         string `required:"true" split_words:"true"`
	OSINTSvcAddr       string `required:"true" split_words:"true"`
}

func newGatewayService() (*gatewayService, error) {
	var conf gatewayConf
	err := envconfig.Process("", &conf)
	if err != nil {
		return nil, err
	}

	if conf.Debug {
		appLogger.SetLevel(logrus.DebugLevel)
	}

	ctx := context.Background()
	return &gatewayService{
		port:          conf.Port,
		uidHeader:     conf.UserIdentityHeader,
		findingClient: finding.NewFindingServiceClient(getGRPCConn(ctx, conf.FindingSvcAddr)),
		iamClient:     iam.NewIAMServiceClient(getGRPCConn(ctx, conf.IAMSvcAddr)),
		projectClient: project.NewProjectServiceClient(getGRPCConn(ctx, conf.ProjectSvcAddr)),
		awsClient:     aws.NewAWSServiceClient(getGRPCConn(ctx, conf.AWSSvcAddr)),
		osintClient:   osint.NewOSINTServiceClient(getGRPCConn(ctx, conf.OSINTSvcAddr)),
	}, nil
}

func getGRPCConn(ctx context.Context, addr string) *grpc.ClientConn {
	conn, err := grpc.DialContext(ctx, addr, grpc.WithInsecure(), grpc.WithTimeout(time.Second*3))
	if err != nil {
		appLogger.Fatalf("Failed to connect backend gRPC server, addr=%s, err=%+v", addr, err)
	}
	return conn
}

func writeResponse(w http.ResponseWriter, status int, body map[string]interface{}) {
	if body == nil {
		w.WriteHeader(status)
		return
	}
	buf, err := json.Marshal(body)
	if err != nil {
		appLogger.Errorf("Response body JSON marshal error: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(buf)
}
