package microgate

import (
	"testing"

	"github.com/microgate-io/microgate-lib-go/v1/httpjson"
	"google.golang.org/protobuf/encoding/protojson"
)

func TestMarshalHandleRequest(t *testing.T) {
	r := new(httpjson.HandleRequest)
	r.Method = "POST"
	r.Body = []byte("hello")
	data, err := protojson.Marshal(r)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(string(data))
}
