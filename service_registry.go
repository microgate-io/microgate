package microgate

import (
	"fmt"
	"strings"

	"github.com/emicklei/xconnect"
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
	endpoints map[string]Endpoint
}

func NewServicRegistry(config xconnect.Document) ServicRegistry {
	sr := ServicRegistry{
		endpoints: map[string]Endpoint{},
	}
	for k, v := range config.XConnect.Connect {
		sr.endpoints[k] = NewEndpoint(v)
	}
	return sr
}

// Lookup returns the Endpoint for sending request to the service operation.
// Example: /UserService/CheckUser
func (r ServicRegistry) Lookup(fullMethodName string) (Endpoint, error) {
	s := strings.Split(fullMethodName, "/")
	if len(s) == 0 {
		return Endpoint{}, fmt.Errorf("failed to parse service:%s", fullMethodName)
	}
	e, ok := r.endpoints[s[1]]
	if !ok {
		return Endpoint{}, fmt.Errorf("unknown service:%s", s[1])
	}
	return e, nil
}
