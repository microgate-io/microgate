package microgate

import (
	"context"
	"net"

	"github.com/emicklei/xconnect"
	mlog "github.com/microgate-io/microgate/v1/log"
	"github.com/vgough/grpc-proxy/proxy"
	"go.opencensus.io/plugin/ocgrpc"
	"google.golang.org/grpc"
)

// StartExternalProxyServer listens to gRPC requests send from a gRPC client.
func StartExternalProxyServer(config xconnect.Document) {
	ctx := context.Background()
	lis, err := net.Listen("tcp", ":9292")
	if err != nil {
		mlog.Fatalw(ctx, "failed to listen", "err", err)
	}

	// conditionally add the handler to serve metrics
	var statsHandler grpc.ServerOption = grpc.EmptyServerOption{}
	statsHandler = grpc.StatsHandler(new(ocgrpc.ServerHandler))

	// reusable gRPC connections
	pool := newConnectionPool()

	// Create gRPC server with interceptors
	director := newDirector(pool, config)
	grpcServer := grpc.NewServer(
		// special codec which allows the proxy to handle raw byte frames and pass them along without any serialization.
		grpc.CustomCodec(proxy.Codec()),
		statsHandler,
		grpc.UnknownServiceHandler(proxy.TransparentHandler(director)))

	mlog.Infow(ctx, "external serving gRPC", "addr", ":9292")
	grpcServer.Serve(lis)
}