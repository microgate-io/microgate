package main

import (
	"github.com/microgate-io/microgate"
	apilog "github.com/microgate-io/microgate-lib-go/v1/log"
	"github.com/microgate-io/microgate/v1/config"
	mconfig "github.com/microgate-io/microgate/v1/config"
	mlog "github.com/microgate-io/microgate/v1/log"
	mqueue "github.com/microgate-io/microgate/v1/queue"
)

func main() {
	mlog.Init()

	gateConfig := config.Load("config.yaml")
	apilog.GlobalDebug, _ = gateConfig.FindBool("global_debug")

	// these are the gRPC services provided to the backend
	provider := microgate.ServiceProvider{
		Log:      mlog.NewLogService(),
		Config:   mconfig.NewConfigServiceImpl(),
		Queueing: mqueue.NewQueueingServiceImpl(gateConfig),
	}

	go microgate.StartInternalProxyServer(gateConfig, provider)
	go microgate.StartExternalProxyHTTPServer(gateConfig)
	microgate.StartExternalProxyServer(gateConfig)
}
