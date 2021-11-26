package microgate

import (
	apiconfig "github.com/microgate-io/microgate-lib-go/v1/config"
	apidb "github.com/microgate-io/microgate-lib-go/v1/db"
	apilog "github.com/microgate-io/microgate-lib-go/v1/log"
	apiqueue "github.com/microgate-io/microgate-lib-go/v1/queue"
)

type ServiceProvider struct {
	Log      apilog.LogServiceServer
	Config   apiconfig.ConfigServiceServer
	Database apidb.DatabaseServiceServer
	Queueing apiqueue.QueueingServiceServer
}
