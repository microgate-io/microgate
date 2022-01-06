package microgate

import (
	"context"
	"fmt"
	"io"
	"strings"

	"github.com/emicklei/xconnect"
	"github.com/iancoleman/strcase"
	"github.com/jhump/protoreflect/desc"
	mlog "github.com/microgate-io/microgate/v1/log"
)

type Endpoint struct {
	HostPort string
	Secure   bool
	ApiKey   string
}

func NewEndpoint(e xconnect.ConnectEntry) Endpoint {
	p := 443 // default for secure
	if e.Port != nil {
		p = *e.Port
	}
	ep := Endpoint{
		HostPort: fmt.Sprintf("%s:%d", e.Host, p),
	}
	if e.Secure != nil {
		ep.Secure = *e.Secure
	}
	if v, ok := e.ExtraFields["api-key"]; ok {
		ep.ApiKey = v.(string)
	}
	return ep
}

type ServicRegistry struct {
	endpoints         map[string]Endpoint
	methodDescriptors map[string]*desc.MethodDescriptor
}

func NewServicRegistry(config xconnect.Document) ServicRegistry {
	sr := ServicRegistry{
		endpoints:         map[string]Endpoint{},
		methodDescriptors: map[string]*desc.MethodDescriptor{},
	}
	for k, v := range config.XConnect.Connect {
		sr.endpoints[k] = NewEndpoint(v)
	}
	return sr
}

// LookupEndpoint returns the Endpoint for sending request to the service operation.
// Example: /UserService/CheckUser
func (r ServicRegistry) LookupEndpoint(fullMethodName string) (Endpoint, error) {
	s := strings.Split(fullMethodName, "/")
	if len(s) < 2 {
		return Endpoint{}, fmt.Errorf("failed to parse service:%s", fullMethodName)
	}
	e, ok := r.endpoints[s[1]]
	if !ok {
		return Endpoint{}, fmt.Errorf("unknown service:%s", s[1])
	}
	return e, nil
}

func (r ServicRegistry) LookupMethod(path string) (*desc.MethodDescriptor, bool) {
	md, ok := r.methodDescriptors[path]
	return md, ok
}

func (r ServicRegistry) AddService(d *desc.ServiceDescriptor) {
	ctx := context.Background()
	if !canExposeAsHTTP(d.GetFullyQualifiedName()) {
		mlog.Debugw(ctx, "reject service as HTTP", "name", d.GetFullyQualifiedName())
		return
	}
	for _, each := range d.GetMethods() {
		path := toHTTPPath(d.GetFullyQualifiedName(), each.GetName())
		mlog.Debugw(ctx, "rpc", "path", path, "method:", each.GetName(), "input:", each.GetInputType().GetName(), "output:", each.GetOutputType().GetName())
		r.methodDescriptors[path] = each
	}
}

// "todo.v1.TodoService", "CreateTodo" -> "/todo/v1/todo-service/create-todo"
func toHTTPPath(service, method string) string {
	sb := new(strings.Builder)
	io.WriteString(sb, "/")
	sp := strings.Split(service, ".")
	for _, each := range sp {
		if strings.HasPrefix(each, "v") { // do not transform version
			io.WriteString(sb, each)
		} else {
			io.WriteString(sb, strcase.ToKebab(each))
		}
		io.WriteString(sb, "/")
	}
	io.WriteString(sb, strcase.ToKebab(method))
	return sb.String()
}

func canExposeAsHTTP(service string) bool {
	if strings.HasPrefix(service, "grpc") {
		return false
	}
	if strings.HasPrefix(service, "microgate") {
		return false
	}
	return true
}
