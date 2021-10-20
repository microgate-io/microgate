package microgate

// Author: E.Micklei

import (
	"context"
	"testing"

	"github.com/emicklei/xconnect"
	"github.com/vgough/grpc-proxy/connector"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

func Test_backendDirector_Connect(t *testing.T) {
	// create
	cfg := xconnect.Document{ExtraFields: make(map[string]interface{})}
	cfg.ExtraFields["masked_headers"] = "mask-1,mask-2"
	cfg.ExtraFields["verbose"] = true
	cfg.ExtraFields["accesslog_enabled"] = true
	cc := connector.NewCachingConnector(connector.WithDialer(testDialer))
	dir := newDirector(cc, cfg)
	bg := context.Background()
	md := metadata.MD{}
	md.Set("mask-1", "MASK-1")
	md.Set("x-api-key", "god")
	ctx, con, err := dir.Connect(metadata.NewIncomingContext(bg, md), "test")
	if err != nil {
		t.Fatal(err)
	}
	if ctx == nil {
		t.Fatal()
	}
	if con == nil {
		t.Error()
	}
}

func testDialer(ctx context.Context, target string, opts ...grpc.DialOption) (*grpc.ClientConn, error) {
	return new(grpc.ClientConn), nil
}
