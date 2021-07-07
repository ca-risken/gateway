package main

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/CyberAgent/mimosa-aws/proto/activity"
	"github.com/CyberAgent/mimosa-aws/proto/aws"
	"github.com/CyberAgent/mimosa-code/proto/code"
	"github.com/CyberAgent/mimosa-core/proto/alert"
	"github.com/CyberAgent/mimosa-core/proto/finding"
	"github.com/CyberAgent/mimosa-core/proto/iam"
	"github.com/CyberAgent/mimosa-core/proto/project"
	"github.com/CyberAgent/mimosa-core/proto/report"
	"github.com/CyberAgent/mimosa-diagnosis/proto/diagnosis"
	"github.com/CyberAgent/mimosa-google/proto/google"
	"github.com/CyberAgent/mimosa-osint/proto/osint"
	"github.com/aws/aws-xray-sdk-go/xray"
	"github.com/kelseyhightower/envconfig"
	"github.com/sirupsen/logrus"
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
	EnvName            string   `default:"default" split_words:"true"`
	Port               string   `default:"8000"`
	Debug              bool     `default:"false"`
	UserIdentityHeader string   `required:"true" split_words:"true"`
	OidcDataHeader     string   `required:"true" split_words:"true"`
	IdpProviderName    []string `required:"true" split_words:"true"`
	FindingSvcAddr     string   `required:"true" split_words:"true"`
	IAMSvcAddr         string   `required:"true" split_words:"true"`
	ProjectSvcAddr     string   `required:"true" split_words:"true"`
	AlertSvcAddr       string   `required:"true" split_words:"true"`
	ReportSvcAddr      string   `required:"true" split_words:"true"`
	AWSSvcAddr         string   `required:"true" split_words:"true"`
	AWSActivitySvcAddr string   `required:"true" split_words:"true"`
	OSINTSvcAddr       string   `required:"true" split_words:"true"`
	DiagnosisSvcAddr   string   `required:"true" split_words:"true"`
	CodeSvcAddr        string   `required:"true" split_words:"true"`
	GoogleSvcAddr      string   `required:"true" split_words:"true"`
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
	// TODO refactor for x-ray
	conn, err := grpc.DialContext(ctx, addr,
		grpc.WithUnaryInterceptor(xray.UnaryClientInterceptor()), grpc.WithInsecure(), grpc.WithTimeout(time.Second*3))
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
