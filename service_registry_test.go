package microgate

import (
	"testing"

	"github.com/emicklei/xconnect"
)

func TestServicRegistry_Lookup(t *testing.T) {
	p := 1234
	c := xconnect.Document{
		XConnect: xconnect.XConnect{
			Connect: map[string]xconnect.ConnectEntry{
				"TestService": {
					Host: "local",
					Port: &p,
				},
			},
		},
	}
	r := NewServicRegistry(c)
	e, err := r.Lookup("/TestService/Doit")
	if err != nil {
		t.Fail()
	}
	if e.HostPort != "local:1234" {
		t.Fail()
	}
}
