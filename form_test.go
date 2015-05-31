package tango

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

type FormAction struct {
	Ctx
}

func (a *FormAction) Get() string {
	v, _ := a.Forms().Int("test")
	return fmt.Sprintf("%d", v)
}

func TestForm1(t *testing.T) {
	buff := bytes.NewBufferString("")
	recorder := httptest.NewRecorder()
	recorder.Body = buff

	o := Classic()
	o.Get("/", new(FormAction))

	req, err := http.NewRequest("GET", "http://localhost:8000/?test=1", nil)
	if err != nil {
		t.Error(err)
	}

	o.ServeHTTP(recorder, req)
	expect(t, recorder.Code, http.StatusOK)
	refute(t, len(buff.String()), 0)
	expect(t, buff.String(), "1")
}
