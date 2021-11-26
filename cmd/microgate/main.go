package main

import (
	"github.com/microgate-io/microgate"
	apilog "github.com/microgate-io/microgate-lib-go/v1/log"
	"github.com/microgate-io/microgate/v1/config"
	mconfig "github.com/microgate-io/microgate/v1/config"
	mdb "github.com/microgate-io/microgate/v1/db"
	mlog "github.com/microgate-io/microgate/v1/log"
	mqueue "github.com/microgate-io/microgate/v1/queue"
)

func main() {
	mlog.Init()

	gateConfig := config.Load("config.yaml")
	apilog.GlobalDebug = gateConfig.FindBool("global_debug")

	// these are the gRPC services provided to the backend
	provider := microgate.ServiceProvider{
		Log:      mlog.NewLogService(),
		Config:   mconfig.NewConfigServiceImpl(),
		Database: mdb.NewDatabaseServiceImpl(gateConfig),
		Queueing: mqueue.NewQueueingServiceImpl(gateConfig),
	}

	go microgate.StartInternalProxyServer(gateConfig, provider)

	microgate.StartExternalProxyServer(gateConfig)
}
