module github.com/microgate-io/microgate

go 1.16

require (
	github.com/blendle/zapdriver v1.3.1
	github.com/emicklei/xconnect v0.10.1
	github.com/jhump/protoreflect v1.8.2
	github.com/microgate-io/microgate-lib-go v1.0.0
	github.com/pkg/errors v0.9.1 // indirect
	github.com/stretchr/objx v0.2.0 // indirect
	github.com/vgough/grpc-proxy v0.0.0-20210913231538-71832b651269
	go.opencensus.io v0.23.0
	go.uber.org/zap v1.19.1
	google.golang.org/grpc v1.41.0
	google.golang.org/protobuf v1.27.1
	gopkg.in/check.v1 v1.0.0-20190902080502-41f04d3bba15 // indirect
	gopkg.in/yaml.v2 v2.4.0
)

replace github.com/microgate-io/microgate-lib-go v1.0.0 => ../microgate-lib-go
