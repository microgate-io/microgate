package microgate

// Author: E.Micklei

import (
	"context"
	"time"

	mlog "github.com/microgate-io/microgate/v1/log"
	"github.com/vgough/grpc-proxy/connector"
	"github.com/vgough/grpc-proxy/proxy"
	"go.opencensus.io/plugin/ocgrpc"
	"google.golang.org/grpc"
)

func NewConnectionPool() *connector.CachingConnector {
	pool := connector.NewCachingConnector(connector.WithDialer(dialInsecure))
	go func() {
		for {
			time.Sleep(10 * time.Second)
			expired := pool.Expire()
			if len(expired) > 0 {
				mlog.Debugw(context.Background(), "cleaning expired connections", "count", len(expired))
			}
		}
	}()
	return pool
}

// dialInsecure adds the grpc Insecure dial option
func dialInsecure(ctx context.Context, target string, opts ...grpc.DialOption) (*grpc.ClientConn, error) {
	options := append(opts,
		// backend does not accept TLS
		grpc.WithInsecure(),
		// special codec which allows the proxy to handle raw byte frames and pass them along without any serialization.
		grpc.WithCodec(proxy.Codec()),
		// Use the open-census handler to pass tracing context.
		grpc.WithStatsHandler(new(ocgrpc.ClientHandler)))
	if true { //*config.Verbose {
		mlog.Debugw(ctx, "dialing for new backend connection", "address", target)
	}
	return grpc.DialContext(ctx, target, options...)
}
