package microgate

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/emicklei/xconnect"
	"github.com/microgate-io/microgate-lib-go/v1/httpjson"
	mlog "github.com/microgate-io/microgate/v1/log"
	"google.golang.org/protobuf/encoding/protojson"
)

// StartExternalProxyHTTPServer listens to HTTP requests send from a HTTP client.
func StartExternalProxyHTTPServer(config xconnect.Document) {
	ctx := context.Background()
	mlog.Debugw(ctx, "StartExternalProxyHTTPServer")

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
	pbr := new(httpjson.HandleRequest)
	pbr.Method = r.Method
	pbr.Path = r.URL.Path
	headers := map[string]string{}
	for k, v := range r.Header {
		if len(v) == 1 {
			headers[k] = v[0]
		}
	}
	pbr.Headers = headers
	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		mlog.Warnw(r.Context(), "cannot read payload", "err", err)
		return
	}
	r.Body.Close()
	pbr.Body = data
	// TEMP
	content, err := protojson.Marshal(pbr)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		mlog.Errorw(r.Context(), "cannot write payload", "err", err)
		return
	}
	w.Header().Set("content-type", "application/json")
	_, err = w.Write(content)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		mlog.Errorw(r.Context(), "cannot write payload", "err", err)
		return
	}
}
