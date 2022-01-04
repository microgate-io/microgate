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
	e, err := r.LookupEndpoint("/TestService/Doit")
	if err != nil {
		t.Fail()
	}
	if e.HostPort != "local:1234" {
		t.Fail()
	}
}

func TestHTTPPath(t *testing.T) {
	p := toHTTPPath("todo.v1.TodoService", "CreateTodo")
	if got, want := p, "todo/v1/todo-service/create-todo"; got != want {
		t.Errorf("got [%v:%T] want [%v:%T]", got, got, want, want)
	}
}
