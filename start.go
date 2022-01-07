package microgate

import (
	"github.com/emicklei/xconnect"
	"github.com/microgate-io/microgate/internal/common"
	"github.com/microgate-io/microgate/internal/inhttp"
	"github.com/microgate-io/microgate/internal/iogrpc"
)

func Start(gateConfig xconnect.Document, provider ServiceProvider) {
	reg := common.NewServicRegistry(gateConfig)
	impl := iogrpc.ProviderImplementations{
		Log:      provider.Log,
		Config:   provider.Config,
		Queueing: provider.Queueing,
	}
	go iogrpc.StartInternalProxyServer(gateConfig, impl, reg)
	go inhttp.StartExternalProxyHTTPServer(gateConfig, reg)
	iogrpc.StartExternalProxyServer(gateConfig, reg)
}
