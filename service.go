package main

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/ca-risken/common/pkg/logging"
	"github.com/ca-risken/core/proto/ai"
	"github.com/ca-risken/core/proto/alert"
	"github.com/ca-risken/core/proto/finding"
	"github.com/ca-risken/core/proto/iam"
	"github.com/ca-risken/core/proto/organization"
	"github.com/ca-risken/core/proto/organization_iam"
	"github.com/ca-risken/core/proto/project"
	"github.com/ca-risken/core/proto/report"
	"github.com/ca-risken/datasource-api/proto/aws"
	"github.com/ca-risken/datasource-api/proto/azure"
	"github.com/ca-risken/datasource-api/proto/code"
	"github.com/ca-risken/datasource-api/proto/datasource"
	"github.com/ca-risken/datasource-api/proto/diagnosis"
	"github.com/ca-risken/datasource-api/proto/google"
	"github.com/ca-risken/datasource-api/proto/osint"
	grpcmiddleware "github.com/grpc-ecosystem/go-grpc-middleware"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	grpctrace "gopkg.in/DataDog/dd-trace-go.v1/contrib/google.golang.org/grpc"
)

const (
	successJSONKey = "data"
	errorJSONKey   = "error"
)

type gatewayService struct {
	envName                string
	port                   string
	uidHeader              string
	oidcDataHeader         string
	sessionCookieName      []string
	sessionTimeoutSec      int
	findingClient          finding.FindingServiceClient
	iamClient              iam.IAMServiceClient
	projectClient          project.ProjectServiceClient
	alertClient            alert.AlertServiceClient
	reportClient           report.ReportServiceClient
	organizationClient     organization.OrganizationServiceClient
	organization_iamClient organization_iam.OrganizationIAMServiceClient
	awsClient              aws.AWSServiceClient
	osintClient            osint.OsintServiceClient
	diagnosisClient        diagnosis.DiagnosisServiceClient
	codeClient             code.CodeServiceClient
	googleClient           google.GoogleServiceClient
	azureClient            azure.AzureServiceClient
	aiClient               ai.AIServiceClient
	claimsClient           claimsInterface
	datasourceClient       datasource.DataSourceServiceClient
}

func newGatewayService(ctx context.Context, conf *AppConfig) (*gatewayService, error) {
	if conf.Debug {
		appLogger.Level(logging.DebugLevel)
	}

	coreConn, err := getGRPCConn(ctx, conf.CoreAddr)
	if err != nil {
		appLogger.Errorf(ctx, "Failed to get grpc connection to core service, err=%+v", err)
		return nil, err
	}
	datasourceConn, err := getGRPCConn(ctx, conf.DataSourceAPISvcAddr)
	if err != nil {
		appLogger.Errorf(ctx, "Failed to get grpc connection to datasource api service, err=%+v", err)
		return nil, err
	}
	return &gatewayService{
		envName:                conf.EnvName,
		port:                   conf.Port,
		uidHeader:              conf.UserIdentityHeader,
		oidcDataHeader:         conf.OidcDataHeader,
		sessionCookieName:      conf.SessionCookieName,
		sessionTimeoutSec:      conf.SessionTimeoutSec,
		findingClient:          finding.NewFindingServiceClient(coreConn),
		iamClient:              iam.NewIAMServiceClient(coreConn),
		projectClient:          project.NewProjectServiceClient(coreConn),
		alertClient:            alert.NewAlertServiceClient(coreConn),
		reportClient:           report.NewReportServiceClient(coreConn),
		organizationClient:     organization.NewOrganizationServiceClient(coreConn),
		organization_iamClient: organization_iam.NewOrganizationIAMServiceClient(coreConn),
		awsClient:              aws.NewAWSServiceClient(datasourceConn),
		osintClient:            osint.NewOsintServiceClient(datasourceConn),
		diagnosisClient:        diagnosis.NewDiagnosisServiceClient(datasourceConn),
		codeClient:             code.NewCodeServiceClient(datasourceConn),
		googleClient:           google.NewGoogleServiceClient(datasourceConn),
		azureClient:            azure.NewAzureServiceClient(datasourceConn),
		aiClient:               ai.NewAIServiceClient(coreConn),
		claimsClient:           newClaimsClient(conf.Region, conf.UserIdpKey, conf.IdpProviderName, conf.VerifyIDToken),
		datasourceClient:       datasource.NewDataSourceServiceClient(datasourceConn),
	}, nil
}

func getGRPCConn(ctx context.Context, addr string) (*grpc.ClientConn, error) {
	ctx, cancel := context.WithTimeout(ctx, time.Second*3)
	defer cancel()
	conn, err := grpc.DialContext(ctx, addr,
		grpc.WithUnaryInterceptor(
			grpcmiddleware.ChainUnaryClient(
				grpctrace.UnaryClientInterceptor())),
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}
	return conn, nil
}

func notFoundHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	writeResponse(ctx, w, http.StatusNotFound, nil)
}

func commonHeader(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Cache-control", "no-store")
		w.Header().Add("Pragma", "no-cache")
		w.Header().Add("X-Frame-Options", "SAMEORIGIN")
		next.ServeHTTP(w, r)
	}
	return http.HandlerFunc(fn)
}

func writeResponse(ctx context.Context, w http.ResponseWriter, status int, body map[string]interface{}) {
	if body == nil {
		w.WriteHeader(status)
		return
	}
	buf, err := json.Marshal(body)
	if err != nil {
		appLogger.Errorf(ctx, "Response body JSON marshal error: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(status)
	_, err = w.Write(buf)
	if err != nil {
		appLogger.Errorf(ctx, "failed to write response, err=%+v", err)
	}
}
