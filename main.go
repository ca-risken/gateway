package main

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/ca-risken/common/pkg/profiler"
	"github.com/ca-risken/common/pkg/tracer"
	"github.com/gassara-kys/envconfig"
)

const (
	namespace   = "gateway"
	serviceName = "gateway"
)

var (
	samplingRate float64 = 0.3000
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
	SessionCookieName  []string `split_words:"true" default:"AWSELBAuthSessionCookie-0,AWSELBAuthSessionCookie-1"`
	SessionTimeoutSec  int      `split_words:"true" default:"2592000"`
	VerifyIDToken      bool     `split_words:"true" default:"false"`
	UserIdpKey         string   `split_words:"true" default:"preferred_username"`
	Region             string   `default:"ap-northeast-1"`

	CoreAddr             string `required:"true" split_words:"true" default:"core.core.svc.cluster.local:8080"`
	DataSourceAPISvcAddr string `required:"true" split_words:"true" default:"datasource-api.datasource.svc.cluster.local:8081"`
	MaxRequestBodyBytes  int64  `split_words:"true" default:"1048576"`
	ReadHeaderTimeoutSec int    `split_words:"true" default:"5"`
	ReadTimeoutSec       int    `split_words:"true" default:"30"`
	WriteTimeoutSec      int    `split_words:"true" default:"30"`
	IdleTimeoutSec       int    `split_words:"true" default:"120"`
	MaxHeaderBytes       int    `split_words:"true" default:"1048576"`
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
		ServiceName:  getFullServiceName(),
		Environment:  appConfig.EnvName,
		Debug:        appConfig.TraceDebug,
		SamplingRate: &samplingRate,
	}
	tracer.Start(tc)
	defer tracer.Stop()

	svc, err := newGatewayService(ctx, &appConfig)
	if err != nil {
		appLogger.Fatal(ctx, err.Error())
	}
	maxRequestBodyBytes = appConfig.MaxRequestBodyBytes
	router := newRouter(svc)
	appLogger.Infof(ctx, "starting http server at :%s", svc.port)
	err = newHTTPServer(&appConfig, router).ListenAndServe()
	if err != nil {
		appLogger.Fatal(ctx, err.Error())
	}
}

func newHTTPServer(conf *AppConfig, handler http.Handler) *http.Server {
	return &http.Server{
		Addr:              ":" + conf.Port,
		Handler:           handler,
		ReadHeaderTimeout: time.Duration(conf.ReadHeaderTimeoutSec) * time.Second,
		ReadTimeout:       time.Duration(conf.ReadTimeoutSec) * time.Second,
		WriteTimeout:      time.Duration(conf.WriteTimeoutSec) * time.Second,
		IdleTimeout:       time.Duration(conf.IdleTimeoutSec) * time.Second,
		MaxHeaderBytes:    conf.MaxHeaderBytes,
	}
}
