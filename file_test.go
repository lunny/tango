package tango

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestDir1(t *testing.T) {
	buff := bytes.NewBufferString("")
	recorder := httptest.NewRecorder()
	recorder.Body = buff

	tg := New()
	//tg.Get("/test.html", File("./public/test.html"))
	tg.Get("/:name", Dir("./public"))

	req, err := http.NewRequest("GET", "http://localhost:8000/test.html", nil)
	if err != nil {
		t.Error(err)
	}

	tg.ServeHTTP(recorder, req)
	expect(t, recorder.Code, http.StatusOK)
	refute(t, len(buff.String()), 0)
	expect(t, buff.String(), "hello tango")
}

func TestFile1(t *testing.T) {
	buff := bytes.NewBufferString("")
	recorder := httptest.NewRecorder()
	recorder.Body = buff

	tg := New()
	tg.Get("/test.html", File("./public/test.html"))

	req, err := http.NewRequest("GET", "http://localhost:8000/test.html", nil)
	if err != nil {
		t.Error(err)
	}

	tg.ServeHTTP(recorder, req)
	expect(t, recorder.Code, http.StatusOK)
	refute(t, len(buff.String()), 0)
	expect(t, buff.String(), "hello tango")
}
