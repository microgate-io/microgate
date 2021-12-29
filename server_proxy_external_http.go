package microgate

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/emicklei/xconnect"
	"github.com/jhump/protoreflect/grpcreflect"
	mlog "github.com/microgate-io/microgate/v1/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection/grpc_reflection_v1alpha"
)

// StartExternalProxyHTTPServer listens to HTTP requests send from a HTTP client.
func StartExternalProxyHTTPServer(config xconnect.Document) {
	ctx := context.Background()
	load()

	mlog.Infow(ctx, "StartExternalProxyHTTPServer")

	s := &http.Server{
		Addr:           fmt.Sprintf(":%d", 8080),
		Handler:        HTTPHandler{},
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	if err := s.ListenAndServe(); err != nil {
		mlog.Infow(ctx, "stopped StartExternalProxyHTTPServer")
	}
}

type HTTPHandler struct{}

func (h HTTPHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

}

func load() {
	conn, err := grpc.Dial("localhost:9090", grpc.WithInsecure())
	if err != nil {
		log.Println(err)
		return
	}
	defer conn.Close()
	c := grpcreflect.NewClient(context.Background(), grpc_reflection_v1alpha.NewServerReflectionClient(conn))
	s, err := c.ListServices()
	if err != nil {
		log.Println(err)
		return
	}
	for _, each := range s {
		log.Println(each)
		desc, _ := c.ResolveService(each)
		for _, other := range desc.GetMethods() {
			log.Println(other.GetName())
		}
	}
}
