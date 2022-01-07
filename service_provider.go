package microgate

import (
	apiconfig "github.com/microgate-io/microgate-lib-go/v1/config"
	apilog "github.com/microgate-io/microgate-lib-go/v1/log"
	apiqueue "github.com/microgate-io/microgate-lib-go/v1/queue"
)

type ServiceProvider struct {
	Log      apilog.LogServiceServer
	Config   apiconfig.ConfigServiceServer
	Queueing apiqueue.QueueingServiceServer
}
