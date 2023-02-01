package main

import (
	"context"
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
	VerifyIDToken      bool     `split_words:"true" default:"false"`
	UserIdpKey         string   `split_words:"true" default:"preferred_username"`
	Region             string   `default:"ap-northeast-1"`

	CoreAddr             string `required:"true" split_words:"true" default:"core.core.svc.cluster.local:8080"`
	DataSourceAPISvcAddr string `required:"true" split_words:"true" default:"datasource-api.core.svc.cluster.local:8081"`
}

func main() {
	ctx := context.Background()
	var appConfig AppConfig
	err := envconfig.Process("", &appConfig)
	if err != nil {
		appLogger.Fatal(ctx, err.Error())
	}

	pTypes, err := profiler.ConvertProfileTypeFrom(appConfig.ProfileTypes)
	if err != nil {
		appLogger.Fatal(ctx, err.Error())
	}
	pExporter, err := profiler.ConvertExporterTypeFrom(appConfig.ProfileExporter)
	if err != nil {
		appLogger.Fatal(ctx, err.Error())
	}
	pc := profiler.Config{
		ServiceName:  getFullServiceName(),
		EnvName:      appConfig.EnvName,
		ProfileTypes: pTypes,
		ExporterType: pExporter,
	}
	err = pc.Start()
	if err != nil {
		appLogger.Fatal(ctx, err.Error())
	}
	defer pc.Stop()

	tc := &tracer.Config{
		ServiceName: getFullServiceName(),
		Environment: appConfig.EnvName,
		Debug:       appConfig.TraceDebug,
	}
	tracer.Start(tc)
	defer tracer.Stop()

	svc, err := newGatewayService(ctx, &appConfig)
	if err != nil {
		appLogger.Fatal(ctx, err.Error())
	}
	router := newRouter(svc)
	appLogger.Infof(ctx, "starting http server at :%s", svc.port)
	err = http.ListenAndServe(":"+svc.port, router)
	if err != nil {
		appLogger.Fatal(ctx, err.Error())
	}
}
