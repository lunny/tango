package tango

import (
	"testing"
	"bytes"
	"strings"
	"net/http"
	"net/http/httptest"
	"fmt"
)

func TestCookie1(t *testing.T) {
	buff := bytes.NewBufferString("")
	recorder := httptest.NewRecorder()
	recorder.Body = buff

	o := Classic()
	o.Get("/", func(ctx *Context) string {
		ck := ctx.Cookies().Get("name")
		if ck != nil {
			return ck.Value
		}
		return ""
	})

	req, err := http.NewRequest("GET", "http://localhost:8000/", nil)
	if err != nil {
		t.Error(err)
	}
	req.AddCookie(NewCookie("name", "test"))

	o.ServeHTTP(recorder, req)
	expect(t, recorder.Code, http.StatusOK)
	refute(t, len(buff.String()), 0)
	expect(t, buff.String(), "test")
}

func TestCookie2(t *testing.T) {
	buff := bytes.NewBufferString("")
	recorder := httptest.NewRecorder()
	recorder.Body = buff

	o := Classic()
	o.Get("/", func(ctx *Context) string {
		ctx.Cookies().Set(NewCookie("name", "test"))
		return "test"
	})

	req, err := http.NewRequest("GET", "http://localhost:8000/", nil)
	if err != nil {
		t.Error(err)
	}

	o.ServeHTTP(recorder, req)
	expect(t, recorder.Code, http.StatusOK)
	refute(t, len(buff.String()), 0)
	expect(t, buff.String(), "test")
	expect(t, strings.Split(recorder.Header().Get("Set-Cookie"), ";")[0], "name=test")
}

func TestCookie3(t *testing.T) {
	buff := bytes.NewBufferString("")
	recorder := httptest.NewRecorder()
	recorder.Body = buff

	o := Classic()
	o.Get("/", func(ctx *Context) string {
		ck := ctx.SecureCookies("sssss").Get("name")
		if ck != nil {
			return ck.Value
		}
		return ""
	})

	req, err := http.NewRequest("GET", "http://localhost:8000/", nil)
	if err != nil {
		t.Error(err)
	}
	req.AddCookie(NewSecureCookie("sssss", "name", "test"))

	o.ServeHTTP(recorder, req)
	expect(t, recorder.Code, http.StatusOK)
	refute(t, len(buff.String()), 0)
	expect(t, buff.String(), "test")
}

func TestCookie4(t *testing.T) {
	buff := bytes.NewBufferString("")
	recorder := httptest.NewRecorder()
	recorder.Body = buff

	o := Classic()
	o.Get("/", func(ctx *Context) string {
		ctx.SecureCookies("ttttt").Set(NewSecureCookie("ttttt", "name", "test"))
		return "test"
	})

	req, err := http.NewRequest("GET", "http://localhost:8000/", nil)
	if err != nil {
		t.Error(err)
	}

	o.ServeHTTP(recorder, req)
	expect(t, recorder.Code, http.StatusOK)
	refute(t, len(buff.String()), 0)
	expect(t, buff.String(), "test")

	r := strings.Split(recorder.Header().Get("Set-Cookie"), ";")[0]
	s := strings.SplitN(r, "=", 2)
	name, value := s[0], s[1]
	expect(t, name, "name")
	expect(t, parseSecureCookie("ttttt", value), "test")
}
