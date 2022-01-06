package inhttp

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/emicklei/xconnect"
	"github.com/jhump/protoreflect/dynamic"
	"github.com/jhump/protoreflect/dynamic/grpcdynamic"
	"github.com/jhump/protoreflect/grpcreflect"
	"github.com/microgate-io/microgate"
	mlog "github.com/microgate-io/microgate/v1/log"
	grpcpool "github.com/processout/grpc-go-pool"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection/grpc_reflection_v1alpha"
)

// StartExternalProxyHTTPServer listens to HTTP requests send from a HTTP client.
func StartExternalProxyHTTPServer(config xconnect.Document, reg microgate.ServicRegistry) {
	ctx := context.Background()
	fact := func() (*grpc.ClientConn, error) {
		return grpc.Dial("localhost:9090", grpc.WithInsecure())
	}
	grpcPool, err := grpcpool.New(fact, 1, 4, 5*time.Second, 10*time.Second)
	if err != nil {
		mlog.Errorw(ctx, "unable to create grpc pool", "err", err)
		return
	}
	handler := &HTTPHandler{
		messageFactory: dynamic.NewMessageFactoryWithDefaults(),
		registry:       reg,
		pool:           grpcPool,
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
	pool           *grpcpool.Pool
}

func (h *HTTPHandler) loadServiceRegistry() {
	conn, err := h.pool.Get(context.Background())
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
	if ct := r.Header.Get("content-type"); ct != "application/json" {
		w.WriteHeader(http.StatusNotAcceptable)
		mlog.Debugw(r.Context(), "not JSON content", "content-type", ct)
		return
	}
	// consume payload
	payload, err := ioutil.ReadAll(r.Body)
	r.Body.Close()
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		mlog.Errorw(r.Context(), "unable to read payload", "err", err)
		return
	}
	// find service for path
	md, ok := h.registry.LookupMethod(r.URL.Path)
	if !ok {
		w.WriteHeader(http.StatusNotFound)
		mlog.Debugw(r.Context(), "method not found", "path", r.URL.Path)
		return
	}
	// get connection
	conn, err := h.pool.Get(r.Context())
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		mlog.Errorw(r.Context(), "unable to dial backend", "err", err)
		return
	}
	defer conn.Close()
	// build proto request
	request := h.messageFactory.NewDynamicMessage(md.GetInputType())
	if err := request.UnmarshalJSON(payload); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		mlog.Errorw(r.Context(), "unable to unmarshal payload",
			"err", err, "payload", string(payload),
			"inputtype", md.GetInputType().GetName(),
			"rpc", md.GetFullyQualifiedName())
		return
	}
	// call function
	client := grpcdynamic.NewStub(conn)
	resp, err := client.InvokeRpc(r.Context(), md, request)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		mlog.Errorw(r.Context(), "unable to dial backend", "err", err)
		return
	}
	// create response payload
	dresp := resp.(*dynamic.Message)
	data, err := dresp.MarshalJSON()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		mlog.Errorw(r.Context(), "unable to write response", "err", err)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Header().Set("content-type", "application/json")
	w.Write(data)
}
