package tango

import (
	"testing"
)

type MyResponse struct {
	Resp
}

func (m MyResponse) Do() {
	m.ResponseWriter.Write([]byte("my response"))
}

func TestResponse(t *testing.T) {
	o := Classic()
	o.Get("/", new(MyResponse))
	o.Run(":8000")
}
