// Copyright 2015 The Tango Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package tango

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func TestError1(t *testing.T) {
	buff := bytes.NewBufferString("")
	recorder := httptest.NewRecorder()
	recorder.Body = buff

	l := NewLogger(os.Stdout)
	o := Classic(l)
	o.Get("/", func(ctx *Context) error {
		return NotFound()
	})

	req, err := http.NewRequest("GET", "http://localhost:8000/", nil)
	if err != nil {
		t.Error(err)
	}

	o.ServeHTTP(recorder, req)
	expect(t, recorder.Code, http.StatusNotFound)
	refute(t, len(buff.String()), 0)
}

func TestError2(t *testing.T) {
	buff := bytes.NewBufferString("")
	recorder := httptest.NewRecorder()
	recorder.Body = buff

	o := Classic()
	o.Get("/", func(ctx *Context) error {
		return NotSupported()
	})

	req, err := http.NewRequest("GET", "http://localhost:8000/", nil)
	if err != nil {
		t.Error(err)
	}

	o.ServeHTTP(recorder, req)
	expect(t, recorder.Code, http.StatusMethodNotAllowed)
	refute(t, len(buff.String()), 0)
}

func TestError3(t *testing.T) {
	buff := bytes.NewBufferString("")
	recorder := httptest.NewRecorder()
	recorder.Body = buff

	o := Classic()
	o.Get("/", func(ctx *Context) error {
		return InternalServerError()
	})

	req, err := http.NewRequest("GET", "http://localhost:8000/", nil)
	if err != nil {
		t.Error(err)
	}

	o.ServeHTTP(recorder, req)
	expect(t, recorder.Code, http.StatusInternalServerError)
	refute(t, len(buff.String()), 0)
}

func TestError4(t *testing.T) {
	buff := bytes.NewBufferString("")
	recorder := httptest.NewRecorder()
	recorder.Body = buff

	o := Classic()
	o.Patch("/", func(ctx *Context) error {
		return Forbidden()
	})

	req, err := http.NewRequest("PATCH", "http://localhost:8000/", nil)
	if err != nil {
		t.Error(err)
	}

	o.ServeHTTP(recorder, req)
	expect(t, recorder.Code, http.StatusForbidden)
	refute(t, len(buff.String()), 0)
}

func TestError5(t *testing.T) {
	buff := bytes.NewBufferString("")
	recorder := httptest.NewRecorder()
	recorder.Body = buff

	o := Classic()
	o.Trace("/", func(ctx *Context) error {
		return Unauthorized()
	})

	req, err := http.NewRequest("TRACE", "http://localhost:8000/", nil)
	if err != nil {
		t.Error(err)
	}

	o.ServeHTTP(recorder, req)
	expect(t, recorder.Code, http.StatusUnauthorized)
	refute(t, len(buff.String()), 0)
}

var err500 = Abort(500, "error")

func TestError6(t *testing.T) {
	buff := bytes.NewBufferString("")
	recorder := httptest.NewRecorder()
	recorder.Body = buff

	o := Classic()
	o.Head("/", func(ctx *Context) error {
		return err500
	})

	req, err := http.NewRequest("HEAD", "http://localhost:8000/", nil)
	if err != nil {
		t.Error(err)
	}

	o.ServeHTTP(recorder, req)
	expect(t, recorder.Code, err500.Code())
	refute(t, len(buff.String()), 0)
	expect(t, buff.String(), err500.Error())
}

func TestError7(t *testing.T) {
	buff := bytes.NewBufferString("")
	recorder := httptest.NewRecorder()
	recorder.Body = buff

	o := Classic()
	o.Head("/", func(ctx *Context) {
		return
	})

	req, err := http.NewRequest("HEAD", "http://localhost:8000/11?==", nil)
	if err != nil {
		t.Error(err)
	}

	o.ServeHTTP(recorder, req)
	expect(t, recorder.Code, http.StatusNotFound)
	refute(t, len(buff.String()), 0)
}

var (
	prefix = "<html><head>tango</head><body><div>"
	suffix = fmt.Sprintf("</div><div>version: %s</div></body></html>", Version())
)

func TestError8(t *testing.T) {
	buff := bytes.NewBufferString("")
	recorder := httptest.NewRecorder()
	recorder.Body = buff

	o := Classic()
	o.ErrHandler = HandlerFunc(func(ctx *Context) {
		switch res := ctx.Result.(type) {
		case AbortError:
			ctx.WriteHeader(res.Code())
			ctx.Write([]byte(prefix))
			ctx.Write([]byte(res.Error()))
		case error:
			ctx.WriteHeader(http.StatusInternalServerError)
			ctx.Write([]byte(prefix))
			ctx.Write([]byte(res.Error()))
		default:
			ctx.WriteHeader(http.StatusInternalServerError)
			ctx.Write([]byte(prefix))
			ctx.Write([]byte(http.StatusText(http.StatusInternalServerError)))
		}
		ctx.Write([]byte(suffix))
	})

	o.Get("/", func() error {
		return NotFound()
	})

	req, err := http.NewRequest("HEAD", "http://localhost:8000/", nil)
	if err != nil {
		t.Error(err)
	}

	o.ServeHTTP(recorder, req)
	expect(t, recorder.Code, http.StatusNotFound)
	refute(t, len(buff.String()), 0)
	expect(t, buff.String(), prefix+"Not Found"+suffix)
}
