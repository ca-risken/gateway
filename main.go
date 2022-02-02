package main

import (
	"context"
	"net/http"

	"github.com/aws/aws-xray-sdk-go/xray"
	"github.com/ca-risken/common/pkg/trace"
	mimosaxray "github.com/ca-risken/common/pkg/xray"
	"github.com/gassara-kys/envconfig"
)

const (
	namespace   = "gateway"
	serviceName = "gateway"
)

type AppConfig struct {
	Port  string `default:"8000"`
	Debug bool   `default:"false"`

	EnvName       string `split_words:"true" default:"local"`
	TraceExporter string `split_words:"true" default:"nop"`

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

func main() {
	var appConfig AppConfig
	err := envconfig.Process("", &appConfig)
	if err != nil {
		appLogger.Fatal(err.Error())
	}

	err = mimosaxray.InitXRay(xray.Config{})
	if err != nil {
		appLogger.Fatal(err.Error())
	}

	traceConfig := &trace.Config{
		Namespace:    namespace,
		ServiceName:  serviceName,
		Environment:  appConfig.EnvName,
		ExporterType: trace.GetExporterType(appConfig.TraceExporter),
	}
	ctx := context.Background()
	tp, err := trace.Init(ctx, traceConfig)
	if err != nil {
		appLogger.Fatal(err.Error())
	}
	defer func() {
		if err := tp.Shutdown(ctx); err != nil {
			appLogger.Fatal(err.Error())
		}
	}()

	svc, err := newGatewayService(&appConfig)
	if err != nil {
		appLogger.Fatal(err.Error())
	}
	router := newRouter(traceConfig.GetFullServiceName(), svc)
	if err := updateDoc(router); err != nil {
		appLogger.Errorf("Failed to generate document, err=%+v", err)
	}

	appLogger.Infof("starting http server at :%s", svc.port)
	err = http.ListenAndServe(":"+svc.port, router)
	if err != nil {
		appLogger.Fatal(err.Error())
	}
}
