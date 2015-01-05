package tango

import (
	"testing"
	"bytes"
	"net/http/httptest"
	"net/http"

	"github.com/lunny/log"
)

func TestLogger(t *testing.T) {
	buff := bytes.NewBufferString("")
	recorder := httptest.NewRecorder()

	n := NewWithLog(log.New(buff, "[tango] ", 0))
	n.Use(Logging())
	n.UseHandler(http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		rw.WriteHeader(http.StatusNotFound)
	}))

	req, err := http.NewRequest("GET", "http://localhost:3000/foobar", nil)
	if err != nil {
		t.Error(err)
	}

	n.ServeHTTP(recorder, req)
	expect(t, recorder.Code, http.StatusNotFound)
	refute(t, len(buff.String()), 0)
}