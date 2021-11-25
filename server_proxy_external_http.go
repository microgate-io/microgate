package microgate

import (
	"context"

	"github.com/emicklei/xconnect"
	mlog "github.com/microgate-io/microgate/v1/log"
)

// StartExternalProxyHTTPServer listens to HTTP requests send from a HTTP client.
func StartExternalProxyHTTPServer(config xconnect.Document) {
	ctx := context.Background()
	mlog.Debugw(ctx, "StartExternalProxyHTTPServer")

	// Need to fetch backend filescriptors
	// to setup up available HTTP routes such as
	// POST /v1/TodoService/CreateTodo
	// { ... }
}
