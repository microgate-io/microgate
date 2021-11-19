package microgate

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/danielvladco/go-proto-gql/pkg/generator"
	"github.com/danielvladco/go-proto-gql/pkg/server"
	"github.com/emicklei/xconnect"
	"github.com/jhump/protoreflect/desc"
	mlog "github.com/microgate-io/microgate/v1/log"
	"github.com/nautilus/gateway"
	"github.com/nautilus/graphql"
)

const (
	optMergeSchemas    = true
	optGenServiceDescr = true
	myProtoPackage     = "microgate"
)

var (
	optEmptyGoRef generator.GoRef = nil
)

// StartExternalProxyServer listens to gRPC requests send from a gRPC client.
func StartExternalProxyGraphQLServer(config xconnect.Document) {
	ctx := context.Background()

	caller, descs := waitForDescriptors("localhost:9090")

	gqlDesc, err := generator.NewSchemas(rejectPackage(descs, myProtoPackage), optMergeSchemas, optGenServiceDescr, optEmptyGoRef)
	mlog.Debugw(ctx, "NewSchemas", "gqlDesc", gqlDesc, "err", err)

	registry := generator.NewRegistry(gqlDesc)

	queryFactory := gateway.QueryerFactory(func(ctx *gateway.PlanningContext, url string) graphql.Queryer {
		return server.NewQueryer(registry, caller)
	})

	// ?????
	sources := []*graphql.RemoteSchema{{URL: "url1"}}
	sources[0].Schema = gqlDesc.AsGraphql()[0]

	g, err := gateway.New(sources, gateway.WithQueryerFactory(&queryFactory))
	if err != nil {
		mlog.Fatalw(ctx, "new gateway failed", "err", err)
	}

	// start listener
	mux := http.NewServeMux()
	mux.HandleFunc("/query", g.GraphQLHandler)
	if true {
		mux.HandleFunc("/playground", g.PlaygroundHandler)
	}
	mlog.Infow(ctx, "external serving GraphQL", "addr", ":8080", "path", "/query", "play", "/playground")
	log.Fatal(http.ListenAndServe(":8080", mux))
}

func waitForDescriptors(hostPort string) (caller server.Caller, list []*desc.FileDescriptor) {
	ctx := context.Background()
	for {
		caller, descs, _, err := server.NewReflectCaller([]string{"localhost:9090"})
		if err != nil {
			mlog.Warnw(ctx, "NewReflectCaller", "wait", "5 sec", "err", err)
			// better pass timeout to newreflectcaller
			// TODO config
			time.Sleep(5 * time.Second)
		} else {
			mlog.Debugw(ctx, "NewReflectCaller", "caller", caller, "desc", descs)
			for _, each := range descs {
				mlog.Debugw(ctx, "descriptor", "name", each.GetName(), "pkg", each.GetPackage())
			}
			return caller, descs
		}
	}
}

func rejectPackage(descs []*desc.FileDescriptor, pkg string) (list []*desc.FileDescriptor) {
	for _, each := range descs {
		if each.GetPackage() != pkg {
			list = append(list, each)
		}
	}
	return
}
