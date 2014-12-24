package tango

import (
	"testing"
	"bytes"
	"net/http/httptest"
	"net/http"
	"fmt"
)

type ParamAction struct {
	Params
}

func (p *ParamAction) Get() string {
	return p.Params.Get(":name")
}

func TestParams(t *testing.T) {
	buff := bytes.NewBufferString("")
	recorder := httptest.NewRecorder()
	recorder.Body = buff

	o := Classic()
	o.Get("/:name", new(ParamAction))

	req, err := http.NewRequest("GET", "http://localhost:3000/foobar", nil)
	if err != nil {
		t.Error(err)
	}

	o.ServeHTTP(recorder, req)
	expect(t, recorder.Code, http.StatusOK)
	refute(t, len(buff.String()), 0)
	fmt.Println(recorder.Body)
}
