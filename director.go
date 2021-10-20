package microgate

// Author: E.Micklei

import (
	"context"
	"fmt"
	"strings"

	"github.com/emicklei/xconnect"
	mlog "github.com/microgate-io/microgate/v1/log"
	"github.com/vgough/grpc-proxy/connector"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

var (
	apikeykey = "x-api-key"
	debugkey  = "x-cloud-debug"
)

// backendDirector implements a director which uses the caching connector for reusing connections.
// for access logs it knows how to mask certain header values.
type backendDirector struct {
	connector           *connector.CachingConnector // use grpc.Pool?
	apichecker          APIChecker
	accessMaskHeaderMap map[string]bool
	verbose             bool
	accessLogEnabled    bool
	registry            ServicRegistry
}

func newDirector(c *connector.CachingConnector, config xconnect.Document) *backendDirector {
	hmp := map[string]bool{}
	commaSeparatedStringOfMaskedHeaders := config.FindString("masked_headers")
	for _, each := range strings.Split(commaSeparatedStringOfMaskedHeaders, ",") {
		hmp[strings.TrimSpace(each)] = true
	}
	var checker APIChecker = AllowAll{}
	return &backendDirector{
		connector:           c,
		accessMaskHeaderMap: hmp,
		apichecker:          checker,
		verbose:             config.FindBool("verbose"),
		accessLogEnabled:    config.FindBool("accesslog_enabled"),
		registry:            NewServicRegistry(config),
	}
}

func (d *backendDirector) Connect(ctx context.Context, fullMethodName string) (context.Context, *grpc.ClientConn, error) {
	mlog.Debugw(ctx, "Connect", "fullMethodName", fullMethodName)

	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		mlog.Warnw(ctx, "requested connection to call", "method", fullMethodName, "error", "no incoming metadata in context")
	}

	// find endpoint and fallback to backend
	endpoint, err := d.registry.Lookup(fullMethodName)
	if err != nil {
		//mlog.Errorw("failed to resolve service", "fullMethodName", fullMethodName, "err", err)
		//return ctx, nil, fmt.Errorf("failed to resolve service:%v", err)

		// fallback to backend so we do not have to register all backend services too
		endpoint = Endpoint{HostPort: "localhost:9090", Secure: false}
		if keys := md[apikeykey]; len(keys) > 0 {
			endpoint.ApiKey = keys[0]
		}
	}

	if len(endpoint.ApiKey) == 0 {
		return ctx, nil, fmt.Errorf("access denied because missing value for header:%s", apikeykey)
	}

	allowed, reason, err := d.apichecker.Check(fullMethodName, endpoint.ApiKey)
	if !allowed {
		// either error or some other reason
		if err != nil {
			return ctx, nil, fmt.Errorf("access denied because of error:%v", err)
		}
		return ctx, nil, fmt.Errorf("access denied because:%s", reason)
	}

	// pass all metadata to outgoing calls
	contextOut := metadata.NewOutgoingContext(ctx, md)
	conn, err := d.connector.Dial(contextOut, endpoint.HostPort)

	if err != nil || d.accessLogEnabled || d.verbose || isCloudDebug(md[debugkey]) {
		kv := []interface{}{"address", endpoint, "method", fullMethodName}
		if err != nil {
			// if dial failed then log that will all details
			kv = append(kv, "error", err)
		}
		if len(endpoint.ApiKey) > 0 {
			// mask the api key value in the log
			kv = append(kv, apikeykey, masked(endpoint.ApiKey))
		}
		for k, v := range md {
			// do not pass apikey
			if k == apikeykey {
				continue
			}
			// replace v for logging of masked header values
			if _, ok := d.accessMaskHeaderMap[k]; ok && len(v) > 0 {
				kv = append(kv, k, masked(v[0]))
			} else {
				// unmasked, take first element if that's the only one
				if len(v) == 1 {
					kv = append(kv, k, v[0])
				} else {
					kv = append(kv, k, v)
				}
			}
		}
		// user proper log level
		if err != nil {
			mlog.Warnw(ctx, "request details", kv...)
		} else {
			mlog.Infow(ctx, "request details", kv...)
		}
	}
	// store the actual endpoint for release
	return context.WithValue(ctx, addressKey, endpoint), conn, err
}

var addressKey struct{}

func (d *backendDirector) Release(ctx context.Context, conn *grpc.ClientConn) {
	// fetch the actual endpoint used to create the connection
	e := ctx.Value(addressKey).(Endpoint)
	d.connector.Release(e.HostPort, conn)
}

func masked(s string) string {
	return fmt.Sprintf("*** masked %d chars ***", len(s))
}

func isCloudDebug(vals []string) bool {
	if len(vals) == 0 {
		return false
	}
	val := vals[0]
	if len(val) == 0 {
		return false
	}
	if val == "null" || val == "nil" || val == "undefined" {
		return false
	}
	return true
}
