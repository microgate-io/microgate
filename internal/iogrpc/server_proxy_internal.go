package iogrpc

import (
	"context"
	"net"

	"github.com/emicklei/xconnect"
	apiconfig "github.com/microgate-io/microgate-lib-go/v1/config"
	apilog "github.com/microgate-io/microgate-lib-go/v1/log"
	apiqueue "github.com/microgate-io/microgate-lib-go/v1/queue"
	"github.com/microgate-io/microgate/internal/common"
	mlog "github.com/microgate-io/microgate/v1/log"
	"github.com/vgough/grpc-proxy/proxy"
	"go.opencensus.io/plugin/ocgrpc"
	"google.golang.org/grpc"
)

type ProviderImplementations struct {
	Log      apilog.LogServiceServer
	Config   apiconfig.ConfigServiceServer
	Queueing apiqueue.QueueingServiceServer
}

// StartInternalProxyServer listens to gRPC requests send from the backend.

func StartInternalProxyServer(config xconnect.Document, provider ProviderImplementations, reg common.ServicRegistry) {
	ctx := context.Background()
	lis, err := net.Listen("tcp", ":9191")
	if err != nil {
		mlog.Fatalw(ctx, "failed to listen", "err", err)
	}

	// conditionally add the handler to serve metrics
	var statsHandler grpc.ServerOption = grpc.EmptyServerOption{}
	statsHandler = grpc.StatsHandler(new(ocgrpc.ServerHandler))

	// reusable gRPC connections
	pool := common.NewConnectionPool()

	// Create gRPC server with interceptors
	director := newDirector(pool, config, reg)
	grpcServer := grpc.NewServer(
		// special codec which allows the proxy to handle raw byte frames and pass them along without any serialization.
		grpc.CustomCodec(proxy.Codec()),
		statsHandler,
		grpc.UnknownServiceHandler(proxy.TransparentHandler(director)))

	if provider.Log != nil {
		apilog.RegisterLogServiceServer(grpcServer, provider.Log)
	}
	if provider.Config != nil {
		apiconfig.RegisterConfigServiceServer(grpcServer, provider.Config)
	}
	if provider.Queueing != nil {
		apiqueue.RegisterQueueingServiceServer(grpcServer, provider.Queueing)
	}

	mlog.Infow(ctx, "internal serving gRPC", "addr", ":9191")
	grpcServer.Serve(lis)
}
