package microgate

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/emicklei/xconnect"
	"github.com/jhump/protoreflect/dynamic"
	"github.com/jhump/protoreflect/dynamic/grpcdynamic"
	"github.com/jhump/protoreflect/grpcreflect"
	mlog "github.com/microgate-io/microgate/v1/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
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

// TEMP
func load() {
	conn, err := grpc.Dial("localhost:9090", grpc.WithInsecure())
	if err != nil {
		mlog.Infow(context.Background(), "unable to dail :9090", "err", err)
	}
	defer conn.Close()

	c := grpcreflect.NewClient(context.Background(), grpc_reflection_v1alpha.NewServerReflectionClient(conn))
	s, err := c.ListServices()
	if err != nil {
		log.Println(err)
		time.Sleep(5 * time.Second)
		load()
		return
	}
	df := dynamic.NewMessageFactoryWithDefaults()
	for _, each := range s {
		desc, _ := c.ResolveService(each)
		for _, other := range desc.GetMethods() {
			log.Println("full:", desc.GetFullyQualifiedName(), "service:", each, "method:", other.GetName(), "input:", other.GetInputType().GetName(), "output:", other.GetOutputType().GetName())

			if other.GetInputType().GetName() == "CreateTodoRequest" {
				dm := df.NewDynamicMessage(other.GetInputType())

				data := []byte(`{"title":"test-title"}`)
				if err := dm.UnmarshalJSON(data); err != nil {
					log.Println(err)
				} else {
					log.Println(dm.String())
				}

				// call backend
				c := grpcdynamic.NewStub(conn)
				withKey := metadata.AppendToOutgoingContext(context.Background(), "x-api-key", "goldenkey", "x-cloud-debug", "DEBUG")
				resp, err := c.InvokeRpc(withKey, other, dm)
				if err != nil {
					log.Println(err)
				} else {
					log.Println(resp)
				}
			}
		}
	}
}
