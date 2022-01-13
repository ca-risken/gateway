package main

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/aws/aws-xray-sdk-go/xray"
	"github.com/ca-risken/aws/proto/activity"
	"github.com/ca-risken/aws/proto/aws"
	"github.com/ca-risken/code/proto/code"
	"github.com/ca-risken/common/pkg/logging"
	"github.com/ca-risken/core/proto/alert"
	"github.com/ca-risken/core/proto/finding"
	"github.com/ca-risken/core/proto/iam"
	"github.com/ca-risken/core/proto/project"
	"github.com/ca-risken/core/proto/report"
	"github.com/ca-risken/diagnosis/proto/diagnosis"
	"github.com/ca-risken/google/proto/google"
	"github.com/ca-risken/osint/proto/osint"
	"github.com/gassara-kys/envconfig"
	"google.golang.org/grpc"
)

const (
	successJSONKey = "data"
	errorJSONKey   = "error"
)

type gatewayService struct {
	envName           string
	port              string
	uidHeader         string
	oidcDataHeader    string
	idpProviderName   []string
	findingClient     finding.FindingServiceClient
	iamClient         iam.IAMServiceClient
	projectClient     project.ProjectServiceClient
	alertClient       alert.AlertServiceClient
	reportClient      report.ReportServiceClient
	awsClient         aws.AWSServiceClient
	awsActivityClient activity.ActivityServiceClient
	osintClient       osint.OsintServiceClient
	diagnosisClient   diagnosis.DiagnosisServiceClient
	codeClient        code.CodeServiceClient
	googleClient      google.GoogleServiceClient
}

type gatewayConf struct {
	EnvName string `default:"local" split_words:"true"`
	Port    string `default:"8000"`
	Debug   bool   `default:"false"`

	UserIdentityHeader string   `required:"true" split_words:"true" default:"x-amzn-oidc-identity"`
	OidcDataHeader     string   `required:"true" split_words:"true" default:"x-amzn-oidc-data"`
	IdpProviderName    []string `required:"true" split_words:"true" default:"YOUR_IDP1,YOUR_IDP2"`

	FindingSvcAddr     string `required:"true" split_words:"true" default:"finding.core.svc.cluster.local:8001"`
	IAMSvcAddr         string `required:"true" split_words:"true" default:"iam.core.svc.cluster.local:8002"`
	ProjectSvcAddr     string `required:"true" split_words:"true" default:"project.core.svc.cluster.local:8003"`
	AlertSvcAddr       string `required:"true" split_words:"true" default:"alert.core.svc.cluster.local:8004"`
	ReportSvcAddr      string `required:"true" split_words:"true" default:"report.core.svc.cluster.local:8005"`
	AWSSvcAddr         string `required:"true" split_words:"true" default:"aws.aws.svc.cluster.local:9001"`
	AWSActivitySvcAddr string `required:"true" split_words:"true" default:"activity.aws.svc.cluster.local:9007"`
	OSINTSvcAddr       string `required:"true" split_words:"true" default:"osint.osint.svc.cluster.local:18081"`
	DiagnosisSvcAddr   string `required:"true" split_words:"true" default:"diagnosis.diagnosis.svc.cluster.local:19001"`
	CodeSvcAddr        string `required:"true" split_words:"true" default:"code.code.svc.cluster.local:10001"`
	GoogleSvcAddr      string `required:"true" split_words:"true" default:"google.google.svc.cluster.local:11001"`
}

func newGatewayService() (*gatewayService, error) {
	var conf gatewayConf
	err := envconfig.Process("", &conf)
	if err != nil {
		return nil, err
	}

	if conf.Debug {
		appLogger.Level(logging.DebugLevel)
	}

	ctx := context.Background()
	return &gatewayService{
		envName:           conf.EnvName,
		port:              conf.Port,
		uidHeader:         conf.UserIdentityHeader,
		oidcDataHeader:    conf.OidcDataHeader,
		idpProviderName:   conf.IdpProviderName,
		findingClient:     finding.NewFindingServiceClient(getGRPCConn(ctx, conf.FindingSvcAddr)),
		iamClient:         iam.NewIAMServiceClient(getGRPCConn(ctx, conf.IAMSvcAddr)),
		projectClient:     project.NewProjectServiceClient(getGRPCConn(ctx, conf.ProjectSvcAddr)),
		alertClient:       alert.NewAlertServiceClient(getGRPCConn(ctx, conf.AlertSvcAddr)),
		reportClient:      report.NewReportServiceClient(getGRPCConn(ctx, conf.ReportSvcAddr)),
		awsClient:         aws.NewAWSServiceClient(getGRPCConn(ctx, conf.AWSSvcAddr)),
		awsActivityClient: activity.NewActivityServiceClient(getGRPCConn(ctx, conf.AWSActivitySvcAddr)),
		osintClient:       osint.NewOsintServiceClient(getGRPCConn(ctx, conf.OSINTSvcAddr)),
		diagnosisClient:   diagnosis.NewDiagnosisServiceClient(getGRPCConn(ctx, conf.DiagnosisSvcAddr)),
		codeClient:        code.NewCodeServiceClient(getGRPCConn(ctx, conf.CodeSvcAddr)),
		googleClient:      google.NewGoogleServiceClient(getGRPCConn(ctx, conf.GoogleSvcAddr)),
	}, nil
}

func getGRPCConn(ctx context.Context, addr string) *grpc.ClientConn {
	ctx, cancel := context.WithTimeout(ctx, time.Second*3)
	defer cancel()
	conn, err := grpc.DialContext(ctx, addr,
		grpc.WithUnaryInterceptor(xray.UnaryClientInterceptor()), grpc.WithInsecure())
	if err != nil {
		appLogger.Fatalf("Failed to connect backend gRPC server, addr=%s, err=%+v", addr, err)
	}
	return conn
}

func notFoundHandler(w http.ResponseWriter, r *http.Request) {
	writeResponse(w, http.StatusNotFound, nil)
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
	// appLogger.Debugf("buf %s", buf)
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(status)
	_, _ = w.Write(buf)
}
