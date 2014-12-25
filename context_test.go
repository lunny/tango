package tango


import (
	"testing"
	"bytes"
	"net/http/httptest"
	"net/http"
)

type CtxAction struct {
	Ctx
}

func (p *CtxAction) Get() {
	p.Ctx.Write([]byte("context"))
}

func TestContext1(t *testing.T) {
	buff := bytes.NewBufferString("")
	recorder := httptest.NewRecorder()
	recorder.Body = buff

	o := Classic()
	o.Get("/", new(CtxAction))

	req, err := http.NewRequest("GET", "http://localhost:8000/", nil)
	if err != nil {
		t.Error(err)
	}

	o.ServeHTTP(recorder, req)
	expect(t, recorder.Code, http.StatusOK)
	refute(t, len(buff.String()), 0)
	expect(t, buff.String(), "context")
}