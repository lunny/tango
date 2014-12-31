package tango

import (
	"testing"
	"bytes"
	"strings"
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

type CtxJsonAction struct {
	Ctx
}

func (p *CtxJsonAction) Get() error {
	return p.Ctx.ServeJson(map[string]string{
		"get": "ctx",
		})
}

func TestContext2(t *testing.T) {
	buff := bytes.NewBufferString("")
	recorder := httptest.NewRecorder()
	recorder.Body = buff

	o := Classic()
	o.Get("/", new(CtxJsonAction))

	req, err := http.NewRequest("GET", "http://localhost:8000/", nil)
	if err != nil {
		t.Error(err)
	}

	o.ServeHTTP(recorder, req)
	expect(t, recorder.Code, http.StatusOK)
	refute(t, len(buff.String()), 0)
	expect(t, strings.TrimSpace(buff.String()), `{"get":"ctx"}`)
}

type CtxXmlAction struct {
	Ctx
}

type XmlStruct struct {
	Content string
}

func (p *CtxXmlAction) Get() error {
	return p.Ctx.ServeXml(XmlStruct{"content"})
}

func TestContext3(t *testing.T) {
	buff := bytes.NewBufferString("")
	recorder := httptest.NewRecorder()
	recorder.Body = buff

	o := Classic()
	o.Get("/", new(CtxXmlAction))

	req, err := http.NewRequest("GET", "http://localhost:8000/", nil)
	if err != nil {
		t.Error(err)
	}

	o.ServeHTTP(recorder, req)
	expect(t, recorder.Code, http.StatusOK)
	refute(t, len(buff.String()), 0)
	expect(t, strings.TrimSpace(buff.String()), `<XmlStruct><Content>content</Content></XmlStruct>`)
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
	o.Get("/", new(CtxFileAction))

	req, err := http.NewRequest("GET", "http://localhost:8000/", nil)
	if err != nil {
		t.Error(err)
	}

	o.ServeHTTP(recorder, req)
	expect(t, recorder.Code, http.StatusOK)
	refute(t, len(buff.String()), 0)
	expect(t, strings.TrimSpace(buff.String()), `this is index.html`)
}