package microgate

import (
	"context"

	"github.com/danielvladco/go-proto-gql/pkg/generator"
	"github.com/danielvladco/go-proto-gql/pkg/server"
	"github.com/emicklei/xconnect"
	mlog "github.com/microgate-io/microgate/v1/log"
)

const (
	optMergeSchemas    = true
	optGenServiceDescr = true
)

var (
	optEmptyGoRef generator.GoRef = nil
)

// StartExternalProxyServer listens to gRPC requests send from a gRPC client.
func StartExternalProxyGraphQLServer(config xconnect.Document) {

	caller, descs, _, err := server.NewReflectCaller([]string{"localhost:9090"})

	mlog.Debugw(context.Background(), "NewReflectCaller", "caller", caller, "desc", descs, "err", err)
	for _, each := range descs {
		mlog.Debugw(context.Background(), "descriptor", "name", each.GetName())
	}

	gqlDesc, err := generator.NewSchemas(descs, optMergeSchemas, optGenServiceDescr, optEmptyGoRef)
	mlog.Debugw(context.Background(), "NewSchemas", "gqlDesc", gqlDesc, "err", err)

	/**
	// start listener
	mux := http.NewServeMux()
	mux.HandleFunc("/query", g.GraphQLHandler)
	if true {
		mux.HandleFunc("/playground", g.PlaygroundHandler)
	}
	log.Fatal(http.ListenAndServe(":8080", mux))
	**/
}
