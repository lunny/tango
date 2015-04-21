package tango

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

type ParamAction struct {
	Params
}

func (p *ParamAction) Get() string {
	return p.Params.Get(":name")
}

func TestParams1(t *testing.T) {
	buff := bytes.NewBufferString("")
	recorder := httptest.NewRecorder()
	recorder.Body = buff

	o := Classic()
	o.Get("/:name", new(ParamAction))

	req, err := http.NewRequest("GET", "http://localhost:8000/foobar", nil)
	if err != nil {
		t.Error(err)
	}

	o.ServeHTTP(recorder, req)
	expect(t, recorder.Code, http.StatusOK)
	refute(t, len(buff.String()), 0)
	expect(t, buff.String(), "foobar")
}

type Param2Action struct {
	Params
}

func (p *Param2Action) Get() string {
	return p.Params.Get(":name")
}

func TestParams2(t *testing.T) {
	buff := bytes.NewBufferString("")
	recorder := httptest.NewRecorder()
	recorder.Body = buff

	o := Classic()
	o.Get("/(:name.*)", new(Param2Action))

	req, err := http.NewRequest("GET", "http://localhost:8000/foobar", nil)
	if err != nil {
		t.Error(err)
	}

	o.ServeHTTP(recorder, req)
	expect(t, recorder.Code, http.StatusOK)
	refute(t, len(buff.String()), 0)
	expect(t, buff.String(), "foobar")
}

type Param3Action struct {
	Ctx
}

func (p *Param3Action) Get() string {
	fmt.Println(p.params)
	p.Params().Set(":name", "name")
	fmt.Println(p.params)
	return p.Params().Get(":name")
}

func TestParams3(t *testing.T) {
	buff := bytes.NewBufferString("")
	recorder := httptest.NewRecorder()
	recorder.Body = buff

	o := Classic()
	o.Get("/(:name.*)", new(Param3Action))

	req, err := http.NewRequest("GET", "http://localhost:8000/foobar", nil)
	if err != nil {
		t.Error(err)
	}

	o.ServeHTTP(recorder, req)
	expect(t, recorder.Code, http.StatusOK)
	refute(t, len(buff.String()), 0)
	expect(t, buff.String(), "name")
}

type Param4Action struct {
	Params
}

func (p *Param4Action) Get() string {
	fmt.Println(p.Params)
	p.Params.Set(":name", "name")
	fmt.Println(p.Params)
	return p.Params.Get(":name")
}

func TestParams4(t *testing.T) {
	buff := bytes.NewBufferString("")
	recorder := httptest.NewRecorder()
	recorder.Body = buff

	o := Classic()
	o.Get("/(:name.*)", new(Param4Action))

	req, err := http.NewRequest("GET", "http://localhost:8000/foobar", nil)
	if err != nil {
		t.Error(err)
	}

	o.ServeHTTP(recorder, req)
	expect(t, recorder.Code, http.StatusOK)
	refute(t, len(buff.String()), 0)
	expect(t, buff.String(), "name")
}

func TestParams5(t *testing.T) {
	buff := bytes.NewBufferString("")
	recorder := httptest.NewRecorder()
	recorder.Body = buff

	o := Classic()
	o.Get("/", func(ctx *Context) {
		ctx.Params().Set(":name", "test")
		ctx.Write([]byte(ctx.Params().Get(":name")))
	})

	req, err := http.NewRequest("GET", "http://localhost:8000/", nil)
	if err != nil {
		t.Error(err)
	}

	o.ServeHTTP(recorder, req)
	expect(t, recorder.Code, http.StatusOK)
	expect(t, buff.String(), "test")
}

func TestParams6(t *testing.T) {
	buff := bytes.NewBufferString("")
	recorder := httptest.NewRecorder()
	recorder.Body = buff

	o := Classic()
	o.Get("/", func(ctx *Context) {
		ctx.Write([]byte(ctx.Params().Get(":name")))
	})

	req, err := http.NewRequest("GET", "http://localhost:8000/", nil)
	if err != nil {
		t.Error(err)
	}

	o.ServeHTTP(recorder, req)
	expect(t, recorder.Code, http.StatusOK)
	expect(t, buff.String(), "")
}
