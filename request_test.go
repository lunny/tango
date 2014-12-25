package tango


import (
	"testing"
	"bytes"
	"net/http/httptest"
	"net/http"
)

type ReqAction struct {
	Req
}

func (p *ReqAction) Get() string {
	return p.Request.FormValue("id")
}

func TestRequest(t *testing.T) {
	buff := bytes.NewBufferString("")
	recorder := httptest.NewRecorder()
	recorder.Body = buff

	o := Classic()
	o.Get("/", new(ReqAction))

	req, err := http.NewRequest("GET", "http://localhost:8000/?id=1", nil)
	if err != nil {
		t.Error(err)
	}

	o.ServeHTTP(recorder, req)
	expect(t, recorder.Code, http.StatusOK)
	refute(t, len(buff.String()), 0)
	expect(t, buff.String(), "1")
}