module github.com/microgate-io/microgate

go 1.16

require (
	github.com/blendle/zapdriver v1.3.1
	github.com/microgate-io/microgate-lib-go v1.0.0
	github.com/emicklei/proto-contrib v0.9.2 // indirect
	github.com/emicklei/xconnect v0.9.6
	github.com/google/uuid v1.3.0
	github.com/jackc/pgtype v1.8.1
	github.com/jackc/pgx v3.6.2+incompatible
	github.com/jackc/pgx/v4 v4.13.0
	github.com/vgough/grpc-proxy v0.0.0-20210913231538-71832b651269
	go.opencensus.io v0.23.0
	go.uber.org/zap v1.19.1
	google.golang.org/grpc v1.41.0
	google.golang.org/protobuf v1.27.1
	gopkg.in/yaml.v2 v2.4.0
)

replace github.com/microgate-io/microgate-lib-go v1.0.0 => ../microgate-lib-go
