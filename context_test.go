// Copyright 2015 The Tango Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package tango

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
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

type CtxJSONAction struct {
	Ctx
}

func (p *CtxJSONAction) Get() error {
	return p.Ctx.ServeJson(map[string]string{
		"get": "ctx",
	})
}

func TestContext2(t *testing.T) {
	buff := bytes.NewBufferString("")
	recorder := httptest.NewRecorder()
	recorder.Body = buff

	o := Classic()
	o.Get("/", new(CtxJSONAction))

	req, err := http.NewRequest("GET", "http://localhost:8000/", nil)
	if err != nil {
		t.Error(err)
	}

	o.ServeHTTP(recorder, req)
	expect(t, recorder.Code, http.StatusOK)
	refute(t, len(buff.String()), 0)
	expect(t, recorder.Header().Get("Content-Type"), "application/json; charset=UTF-8")
	expect(t, strings.TrimSpace(buff.String()), `{"get":"ctx"}`)
}

type CtxXMLAction struct {
	Ctx
}

type XMLStruct struct {
	Content string
}

func (p *CtxXMLAction) Get() error {
	return p.Ctx.ServeXml(XMLStruct{"content"})
}

func TestContext3(t *testing.T) {
	buff := bytes.NewBufferString("")
	recorder := httptest.NewRecorder()
	recorder.Body = buff

	o := Classic()
	o.Get("/", new(CtxXMLAction))

	req, err := http.NewRequest("GET", "http://localhost:8000/", nil)
	if err != nil {
		t.Error(err)
	}

	o.ServeHTTP(recorder, req)
	expect(t, recorder.Code, http.StatusOK)
	expect(t, recorder.Header().Get("Content-Type"), "application/xml; charset=UTF-8")
	refute(t, len(buff.String()), 0)
	expect(t, strings.TrimSpace(buff.String()), `<XMLStruct><Content>content</Content></XMLStruct>`)
}

type CtxFileAction struct {
	Ctx
}

func (p *CtxFileAction) Get() error {
	return p.Ctx.ServeFile("./public/index.html")
}

func TestContext4(t *testing.T) {
	buff := bytes.NewBufferString("")
	recorder := httptest.NewRecorder()
	recorder.Body = buff

	o := Classic()
	o.Any("/", new(CtxFileAction))

	req, err := http.NewRequest("GET", "http://localhost:8000/", nil)
	if err != nil {
		t.Error(err)
	}

	o.ServeHTTP(recorder, req)
	expect(t, recorder.Code, http.StatusOK)
	refute(t, len(buff.String()), 0)
	expect(t, strings.TrimSpace(buff.String()), `this is index.html`)
}

func TestContext5(t *testing.T) {
	buff := bytes.NewBufferString("")
	recorder := httptest.NewRecorder()
	recorder.Body = buff

	o := Classic()
	o.Any("/2", func() string {
		return "2"
	})
	o.Any("/", func(ctx *Context) {
		ctx.Redirect("/2")
	})

	req, err := http.NewRequest("GET", "http://localhost:8000/", nil)
	if err != nil {
		t.Error(err)
	}

	o.ServeHTTP(recorder, req)
	expect(t, recorder.Code, http.StatusFound)
	refute(t, len(buff.String()), 0)
	expect(t, strings.TrimSpace(buff.String()), `<a href="/2">Found</a>.`)
}

type NotFoundAction struct {
	Ctx
}

func (n *NotFoundAction) Get() {
	n.NotFound("not found")
}

func TestContext6(t *testing.T) {
	buff := bytes.NewBufferString("")
	recorder := httptest.NewRecorder()
	recorder.Body = buff

	o := Classic()
	o.Any("/", new(NotFoundAction))

	req, err := http.NewRequest("GET", "http://localhost:8000/", nil)
	if err != nil {
		t.Error(err)
	}

	o.ServeHTTP(recorder, req)
	expect(t, recorder.Code, http.StatusNotFound)
	refute(t, len(buff.String()), 0)
	expect(t, strings.TrimSpace(buff.String()), "not found")
}

type NotModifidAction struct {
	Ctx
}

func (n *NotModifidAction) Get() {
	n.NotModified()
}

func TestContext7(t *testing.T) {
	buff := bytes.NewBufferString("")
	recorder := httptest.NewRecorder()
	recorder.Body = buff

	o := Classic()
	o.Any("/", new(NotModifidAction))

	req, err := http.NewRequest("GET", "http://localhost:8000/", nil)
	if err != nil {
		t.Error(err)
	}

	o.ServeHTTP(recorder, req)
	expect(t, recorder.Code, http.StatusNotModified)
	expect(t, len(buff.String()), 0)
}

type UnauthorizedAction struct {
	Ctx
}

func (n *UnauthorizedAction) Get() {
	n.Unauthorized()
}

func TestContext8(t *testing.T) {
	buff := bytes.NewBufferString("")
	recorder := httptest.NewRecorder()
	recorder.Body = buff

	o := Classic()
	o.Any("/", new(UnauthorizedAction))

	req, err := http.NewRequest("GET", "http://localhost:8000/", nil)
	if err != nil {
		t.Error(err)
	}

	o.ServeHTTP(recorder, req)
	expect(t, recorder.Code, http.StatusUnauthorized)
	expect(t, buff.String(), http.StatusText(http.StatusUnauthorized))
}

type DownloadAction struct {
	Ctx
}

func (n *DownloadAction) Get() {
	n.Download("./public/index.html")
}

func TestContext9(t *testing.T) {
	buff := bytes.NewBufferString("")
	recorder := httptest.NewRecorder()
	recorder.Body = buff

	o := Classic()
	o.Any("/", new(DownloadAction))

	req, err := http.NewRequest("GET", "http://localhost:8000/", nil)
	if err != nil {
		t.Error(err)
	}

	o.ServeHTTP(recorder, req)

	expect(t, recorder.Header().Get("Content-Disposition"), `attachment; filename="index.html"`)
	expect(t, recorder.Code, http.StatusOK)
	expect(t, buff.String(), "this is index.html")
}

// check unsupported function will panic
func TestContext10(t *testing.T) {
	buff := bytes.NewBufferString("")
	recorder := httptest.NewRecorder()
	recorder.Body = buff

	o := Classic()
	var ifPanic bool
	defer func() {
		if err := recover(); err != nil {
			ifPanic = true
		}

		expect(t, ifPanic, true)
	}()

	o.Any("/", func(i int) {
		fmt.Println(i)
	})
}

type DownloadAction2 struct {
	Ctx
}

func (n *DownloadAction2) Get() {
	n.ServeFile("./public/a.html")
}

func TestContext11(t *testing.T) {
	buff := bytes.NewBufferString("")
	recorder := httptest.NewRecorder()
	recorder.Body = buff

	o := Classic()
	o.Any("/", new(DownloadAction2))

	req, err := http.NewRequest("GET", "http://localhost:8000/", nil)
	if err != nil {
		t.Error(err)
	}

	o.ServeHTTP(recorder, req)

	expect(t, recorder.Code, http.StatusNotFound)
}

func TestContext12(t *testing.T) {
	buff := bytes.NewBufferString("")
	recorder := httptest.NewRecorder()
	recorder.Body = buff

	o := Classic()
	o.Any("/", func() string {
		return "text"
	})

	req, err := http.NewRequest("GET", "http://localhost:8000/", nil)
	if err != nil {
		t.Error(err)
	}

	o.ServeHTTP(recorder, req)

	expect(t, recorder.Code, http.StatusOK)
	expect(t, buff.String(), "text")
}

func TestContext13(t *testing.T) {
	buff := bytes.NewBufferString("")
	recorder := httptest.NewRecorder()
	recorder.Body = buff

	o := Classic()
	o.Get("/", func(responseWriter http.ResponseWriter, req *http.Request) {
		responseWriter.Write([]byte("text"))
	})

	req, err := http.NewRequest("GET", "http://localhost:8000/", nil)
	if err != nil {
		t.Error(err)
	}

	o.ServeHTTP(recorder, req)

	expect(t, recorder.Code, http.StatusOK)
	expect(t, buff.String(), "text")
}

func TestContext14(t *testing.T) {
	buff := bytes.NewBufferString("")
	recorder := httptest.NewRecorder()
	recorder.Body = buff

	o := Classic()
	o.Get("/", func(req *http.Request) string {
		return "text"
	})

	req, err := http.NewRequest("GET", "http://localhost:8000/", nil)
	if err != nil {
		t.Error(err)
	}

	o.ServeHTTP(recorder, req)

	expect(t, recorder.Code, http.StatusOK)
	expect(t, buff.String(), "text")
}

func TestContextDecodeJSON(t *testing.T) {
	buff := bytes.NewBufferString("")
	recorder := httptest.NewRecorder()
	recorder.Body = buff

	type Data struct {
		Name string
	}

	o := Classic()
	o.Post("/", func(ctx *Context) string {
		var data Data
		err := ctx.DecodeJSON(&data)
		if err != nil {
			return err.Error()
		}
		return data.Name
	})

	req, err := http.NewRequest("POST", "http://localhost:8000/", bytes.NewBufferString(`
		{
			"Name": "lunny"
		}
	`))
	if err != nil {
		t.Error(err)
	}

	o.ServeHTTP(recorder, req)

	expect(t, recorder.Code, http.StatusOK)
	expect(t, buff.String(), "lunny")
}

func TestContextDecodeXML(t *testing.T) {
	buff := bytes.NewBufferString("")
	recorder := httptest.NewRecorder()
	recorder.Body = buff

	type Data struct {
		Name string `xml:"name"`
	}

	o := Classic()
	o.Post("/", func(ctx *Context) string {
		var data Data
		err := ctx.DecodeXML(&data)
		if err != nil {
			return err.Error()
		}
		return data.Name
	})

	req, err := http.NewRequest("POST", "http://localhost:8000/", bytes.NewBufferString(`
		<Data><name>lunny</name></Data>
	`))
	if err != nil {
		t.Error(err)
	}

	o.ServeHTTP(recorder, req)

	expect(t, recorder.Code, http.StatusOK)
	expect(t, buff.String(), "lunny")
}
