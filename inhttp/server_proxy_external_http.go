package inhttp

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/emicklei/xconnect"
	"github.com/jhump/protoreflect/dynamic"
	"github.com/jhump/protoreflect/grpcreflect"
	"github.com/microgate-io/microgate"
	mlog "github.com/microgate-io/microgate/v1/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection/grpc_reflection_v1alpha"
)

// StartExternalProxyHTTPServer listens to HTTP requests send from a HTTP client.
func StartExternalProxyHTTPServer(config xconnect.Document, reg microgate.ServicRegistry) {
	ctx := context.Background()

	handler := &HTTPHandler{
		messageFactory: dynamic.NewMessageFactoryWithDefaults(),
		registry:       reg,
	}
	handler.loadServiceRegistry()

	mlog.Infow(ctx, "external serving HTTP", "addr", ":8080")

	s := &http.Server{
		Addr:           fmt.Sprintf(":%d", 8080),
		Handler:        handler,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	if err := s.ListenAndServe(); err != nil {
		mlog.Infow(ctx, "stopped StartExternalProxyHTTPServer")
	}
}

type HTTPHandler struct {
	messageFactory *dynamic.MessageFactory
	registry       microgate.ServicRegistry
}

func (h *HTTPHandler) loadServiceRegistry() {
	conn, err := grpc.Dial("localhost:9090", grpc.WithInsecure())
	if err != nil {
		mlog.Infow(context.Background(), "unable to dail :9090", "err", err)
	}
	defer conn.Close()
	ctx := context.Background()
	c := grpcreflect.NewClient(ctx, grpc_reflection_v1alpha.NewServerReflectionClient(conn))
	names, err := c.ListServices()
	if err != nil {
		log.Println(err)
		// Dial does not connect so only when calling a service the connection is tried
		// if that fails then retry forever until we do.
		time.Sleep(5 * time.Second)
		h.loadServiceRegistry()
		return
	}
	for _, each := range names {
		desc, err := c.ResolveService(each)
		if err != nil {
			mlog.Errorw(ctx, "failed to resolve service", "name", each, "err", err)
			continue
		}
		h.registry.AddService(desc)
	}
}

func (h *HTTPHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

}

/**
func load(reg microgate.ServicRegistry) {

	df := dynamic.NewMessageFactoryWithDefaults()
	for _, each := range s {
		desc, _ := c.ResolveService(each)
		reg.AddService(desc)
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
**/
