package main

import (
	"fmt"
	"net/http"

	"github.com/ca-risken/common/pkg/profiler"
	"github.com/ca-risken/common/pkg/tracer"
	"github.com/gassara-kys/envconfig"
)

const (
	namespace   = "gateway"
	serviceName = "gateway"
)

func getFullServiceName() string {
	return fmt.Sprintf("%s.%s", namespace, serviceName)
}

type AppConfig struct {
	Port            string   `default:"8000"`
	Debug           bool     `default:"false"`
	ProfileExporter string   `split_words:"true" default:"nop"`
	ProfileTypes    []string `split_words:"true"`

	EnvName    string `split_words:"true" default:"local"`
	TraceDebug bool   `split_words:"true" default:"false"`

	UserIdentityHeader string   `required:"true" split_words:"true" default:"x-amzn-oidc-identity"`
	OidcDataHeader     string   `required:"true" split_words:"true" default:"x-amzn-oidc-data"`
	IdpProviderName    []string `required:"true" split_words:"true" default:"YOUR_IDP1,YOUR_IDP2"`

	CoreAddr     string `required:"true" split_words:"true" default:"core.core.svc.cluster.local:8080"`
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

	pTypes, err := profiler.ConvertProfileTypeFrom(appConfig.ProfileTypes)
	if err != nil {
		appLogger.Fatal(err.Error())
	}
	pExporter, err := profiler.ConvertExporterTypeFrom(appConfig.ProfileExporter)
	if err != nil {
		appLogger.Fatal(err.Error())
	}
	pc := profiler.Config{
		ServiceName:  getFullServiceName(),
		EnvName:      appConfig.EnvName,
		ProfileTypes: pTypes,
		ExporterType: pExporter,
	}
	err = pc.Start()
	if err != nil {
		appLogger.Fatal(err.Error())
	}
	defer pc.Stop()

	tc := &tracer.Config{
		ServiceName: getFullServiceName(),
		Environment: appConfig.EnvName,
		Debug:       appConfig.TraceDebug,
	}
	tracer.Start(tc)
	defer tracer.Stop()

	svc, err := newGatewayService(&appConfig)
	if err != nil {
		appLogger.Fatal(err.Error())
	}
	router := newRouter(svc)
	appLogger.Infof("starting http server at :%s", svc.port)
	err = http.ListenAndServe(":"+svc.port, router)
	if err != nil {
		appLogger.Fatal(err.Error())
	}
}
