package main

import (
	"net/http"

	"github.com/aws/aws-xray-sdk-go/xray"
	mimosaxray "github.com/ca-risken/common/pkg/xray"
)

func main() {
	mimosaxray.InitXRay(xray.Config{})

	svc, err := newGatewayService()
	if err != nil {
		appLogger.Fatal(err.Error())
	}
	router := newRouter(svc)
	if err := updateDoc(router); err != nil {
		appLogger.Errorf("Failed to generate document, err=%+v", err)
	}

	appLogger.Infof("starting http server at :%s", svc.port)
	err = http.ListenAndServe(":"+svc.port, router)
	if err != nil {
		appLogger.Fatal(err.Error())
	}
}
