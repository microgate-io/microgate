package microgate

import (
	"context"
	"fmt"
	"net"
	"os"

	"github.com/emicklei/xconnect"
	"github.com/jackc/pgx/v4"
	apiconfig "github.com/microgate-io/microgate-lib-go/v1/config"
	apidb "github.com/microgate-io/microgate-lib-go/v1/db"
	apilog "github.com/microgate-io/microgate-lib-go/v1/log"
	apiqueue "github.com/microgate-io/microgate-lib-go/v1/queue"
	mconfig "github.com/microgate-io/microgate/v1/config"
	mdb "github.com/microgate-io/microgate/v1/db"
	mlog "github.com/microgate-io/microgate/v1/log"
	mqueue "github.com/microgate-io/microgate/v1/queue"
	"github.com/vgough/grpc-proxy/proxy"
	"go.opencensus.io/plugin/ocgrpc"
	"google.golang.org/grpc"
)

// StartInternalProxyServer listens to gRPC requests send from the backend.
func StartInternalProxyServer(config xconnect.Document) {
	ctx := context.Background()
	lis, err := net.Listen("tcp", ":9191")
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

	// log
	l := mlog.NewLogService()
	apilog.RegisterLogServiceServer(grpcServer, l)

	// config
	c := new(mconfig.ConfigServiceImpl)
	apiconfig.RegisterConfigServiceServer(grpcServer, c)

	// db
	// TODO: temp, check config
	conn, err := pgx.Connect(context.Background(), config.FindString("postgres_connect"))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	d := mdb.NewDatabaseServiceImpl(conn)
	// TODO: end temp
	apidb.RegisterDatabaseServiceServer(grpcServer, d)

	// queue
	q := new(mqueue.QueueingServiceImpl)
	apiqueue.RegisterQueueingServiceServer(grpcServer, q)

	mlog.Infow(ctx, "external serving gRPC", "addr", ":9191")
	grpcServer.Serve(lis)
}
