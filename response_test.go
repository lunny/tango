package tango

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"
)

type MyResponse struct {
	Resp
}

func (m MyResponse) Get() {
	m.ResponseWriter.Write([]byte("my response"))
}

func TestResponse(t *testing.T) {
	buff := bytes.NewBufferString("")
	recorder := httptest.NewRecorder()
	recorder.Body = buff

	o := Classic()
	o.Get("/", new(MyResponse))
	
	req, err := http.NewRequest("GET", "http://localhost:8000/", nil)
	if err != nil {
		t.Error(err)
	}

	o.ServeHTTP(recorder, req)
	expect(t, recorder.Code, http.StatusOK)
	refute(t, len(buff.String()), 0)
	expect(t, buff.String(), "my response")
}
