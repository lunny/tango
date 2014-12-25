package tango

import (
	"fmt"
	"testing"
	"bytes"
	"net/http"
	"net/http/httptest"
)

type CompressExample struct {
}

// implemented this method for ask MUST use gzip, if no implemented, 
// then no compress
func (a *CompressExample) CompressType() string {
	return "gzip"
}

func (a *CompressExample) Get() string {
	return fmt.Sprintf("This is a gzip compress text")
}

func TestCompress(t *testing.T) {
	buff := bytes.NewBufferString("")
	recorder := httptest.NewRecorder()
	recorder.Body = buff

	o := Classic()
	o.Get("/", new(CompressExample))

	req, err := http.NewRequest("GET", "http://localhost:8000/", nil)
	if err != nil {
		t.Error(err)
	}
	req.Header.Add("Accept-Encoding", "gzip")

	o.ServeHTTP(recorder, req)
	expect(t, recorder.Code, http.StatusOK)
	refute(t, len(buff.String()), 0)
	expect(t, buff.String(), "This is a gzip compress text")
}
