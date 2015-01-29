package tango

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"
)

type RouterNoMethodAction struct {
}

func TestRouter1(t *testing.T) {
	buff := bytes.NewBufferString("")
	recorder := httptest.NewRecorder()
	recorder.Body = buff

	o := Classic()
	o.Get("/:name", new(RouterNoMethodAction))

	req, err := http.NewRequest("GET", "http://localhost:8000/foobar", nil)
	if err != nil {
		t.Error(err)
	}

	o.ServeHTTP(recorder, req)
	expect(t, recorder.Code, http.StatusNotFound)
	refute(t, len(buff.String()), 0)
}

type RouterGetAction struct {
}

func (a *RouterGetAction) Get() string {
	return "get"
}

func (a *RouterGetAction) Post() string {
	return "post"
}

func TestRouter2(t *testing.T) {
	buff := bytes.NewBufferString("")
	recorder := httptest.NewRecorder()
	recorder.Body = buff

	o := Classic()
	o.Get("/:name", new(RouterGetAction))
	o.Post("/:name", new(RouterGetAction))

	req, err := http.NewRequest("GET", "http://localhost:8000/foobar", nil)
	if err != nil {
		t.Error(err)
	}

	o.ServeHTTP(recorder, req)
	expect(t, recorder.Code, http.StatusOK)
	refute(t, len(buff.String()), 0)
	expect(t, buff.String(), "get")

	buff.Reset()

	req, err = http.NewRequest("POST", "http://localhost:8000/foobar", nil)
	if err != nil {
		t.Error(err)
	}

	o.ServeHTTP(recorder, req)
	expect(t, recorder.Code, http.StatusOK)
	refute(t, len(buff.String()), 0)
	expect(t, buff.String(), "post")
}

type RouterSpecAction struct {
	a string
}

func (RouterSpecAction) Method1() string {
	return "1"
}

func (r *RouterSpecAction) Method2() string {
	return r.a
}

/*
func TestRouter3(t *testing.T) {
	buff := bytes.NewBufferString("")
	recorder := httptest.NewRecorder()
	recorder.Body = buff

	o := Classic()
	o.Get("/1", new(RouterSpecAction).Method1)
	o.Get("/2", new(RouterSpecAction).Method2)

	req, err := http.NewRequest("GET", "http://localhost:8000/1", nil)
	if err != nil {
		t.Error(err)
	}

	o.ServeHTTP(recorder, req)
	expect(t, recorder.Code, http.StatusOK)
	refute(t, len(buff.String()), 0)
}*/

func TestRouterFunc(t *testing.T) {
	buff := bytes.NewBufferString("")
	recorder := httptest.NewRecorder()
	recorder.Body = buff

	o := Classic()
	o.Get("/", func() string {
		return "func"
	})
	o.Post("/", func(ctx *Context) {
		ctx.Write([]byte("func(*Context)"))
	})
	o.Put("/", func(resp http.ResponseWriter, req *http.Request) {
		resp.Write([]byte("func(http.ResponseWriter, *http.Request)"))
	})
	o.Options("/", func(resp http.ResponseWriter) {
		resp.Write([]byte("func(http.ResponseWriter)"))
	})
	o.Delete("/", func(req *http.Request) string {
		return "func(*http.Request)"
	})

	// plain
	req, err := http.NewRequest("GET", "http://localhost:8000/", nil)
	if err != nil {
		t.Error(err)
	}

	o.ServeHTTP(recorder, req)
	expect(t, recorder.Code, http.StatusOK)
	refute(t, len(buff.String()), 0)
	expect(t, buff.String(), "func")

	// context
	buff.Reset()

	req, err = http.NewRequest("POST", "http://localhost:8000/", nil)
	if err != nil {
		t.Error(err)
	}

	o.ServeHTTP(recorder, req)
	expect(t, recorder.Code, http.StatusOK)
	refute(t, len(buff.String()), 0)
	expect(t, buff.String(), "func(*Context)")

	// http
	buff.Reset()

	req, err = http.NewRequest("PUT", "http://localhost:8000/", nil)
	if err != nil {
		t.Error(err)
	}

	o.ServeHTTP(recorder, req)
	expect(t, recorder.Code, http.StatusOK)
	refute(t, len(buff.String()), 0)
	expect(t, buff.String(), "func(http.ResponseWriter, *http.Request)")

	// response
	buff.Reset()

	req, err = http.NewRequest("OPTIONS", "http://localhost:8000/", nil)
	if err != nil {
		t.Error(err)
	}

	o.ServeHTTP(recorder, req)
	expect(t, recorder.Code, http.StatusOK)
	refute(t, len(buff.String()), 0)
	expect(t, buff.String(), "func(http.ResponseWriter)")

	// req
	buff.Reset()

	req, err = http.NewRequest("DELETE", "http://localhost:8000/", nil)
	if err != nil {
		t.Error(err)
	}

	o.ServeHTTP(recorder, req)
	expect(t, recorder.Code, http.StatusOK)
	refute(t, len(buff.String()), 0)
	expect(t, buff.String(), "func(*http.Request)")
}

type Router4Action struct {
	Params
}

func (r *Router4Action) Get() string {
	return r.Params.Get(":name1") + "-" + r.Params.Get(":name2")
}

func TestRouter4(t *testing.T) {
	buff := bytes.NewBufferString("")
	recorder := httptest.NewRecorder()
	recorder.Body = buff

	o := Classic()
	o.Get("/:name1-:name2", new(Router4Action))

	req, err := http.NewRequest("GET", "http://localhost:8000/foobar-foobar2", nil)
	if err != nil {
		t.Error(err)
	}

	o.ServeHTTP(recorder, req)
	expect(t, recorder.Code, http.StatusOK)
	refute(t, len(buff.String()), 0)
	expect(t, buff.String(), "foobar-foobar2")
}

type Router5Action struct {
}

func (r *Router5Action) Get() string {
	return "router5"
}

func TestRouter5(t *testing.T) {
	buff := bytes.NewBufferString("")
	recorder := httptest.NewRecorder()
	recorder.Body = buff

	o := Classic()
	o.Route("GET", "/", new(Router5Action))

	req, err := http.NewRequest("GET", "http://localhost:8000/", nil)
	if err != nil {
		t.Error(err)
	}

	o.ServeHTTP(recorder, req)
	expect(t, recorder.Code, http.StatusOK)
	refute(t, len(buff.String()), 0)
	expect(t, buff.String(), "router5")
}

type Router6Action struct {
}

func (r *Router6Action) MyMethod() string {
	return "router6"
}

func TestRouter6(t *testing.T) {
	buff := bytes.NewBufferString("")
	recorder := httptest.NewRecorder()
	recorder.Body = buff

	o := Classic()
	o.Route("GET:MyMethod", "/", new(Router6Action))

	req, err := http.NewRequest("GET", "http://localhost:8000/", nil)
	if err != nil {
		t.Error(err)
	}

	o.ServeHTTP(recorder, req)
	expect(t, recorder.Code, http.StatusOK)
	refute(t, len(buff.String()), 0)
	expect(t, buff.String(), "router6")
}

type Router7Action struct {
}

func (r *Router7Action) MyGet() string {
	return "router7-get"
}

func (r *Router7Action) Post() string {
	return "router7-post"
}

func TestRouter7(t *testing.T) {
	buff := bytes.NewBufferString("")
	recorder := httptest.NewRecorder()
	recorder.Body = buff

	o := Classic()
	o.Route([]string{"GET:MyGet", "POST"}, "/", new(Router7Action))

	req, err := http.NewRequest("GET", "http://localhost:8000/", nil)
	if err != nil {
		t.Error(err)
	}

	o.ServeHTTP(recorder, req)
	expect(t, recorder.Code, http.StatusOK)
	refute(t, len(buff.String()), 0)
	expect(t, buff.String(), "router7-get")

	buff.Reset()

	req, err = http.NewRequest("POST", "http://localhost:8000/", nil)
	if err != nil {
		t.Error(err)
	}

	o.ServeHTTP(recorder, req)
	expect(t, recorder.Code, http.StatusOK)
	refute(t, len(buff.String()), 0)
	expect(t, buff.String(), "router7-post")
}

type Router8Action struct {
}

func (r *Router8Action) MyGet() string {
	return "router8-get"
}

func (r *Router8Action) Post() string {
	return "router8-post"
}

func TestRouter8(t *testing.T) {
	buff := bytes.NewBufferString("")
	recorder := httptest.NewRecorder()
	recorder.Body = buff

	o := Classic()
	o.Route(map[string]string{"GET": "MyGet", "POST": "Post"}, "/", new(Router8Action))

	req, err := http.NewRequest("GET", "http://localhost:8000/", nil)
	if err != nil {
		t.Error(err)
	}

	o.ServeHTTP(recorder, req)
	expect(t, recorder.Code, http.StatusOK)
	refute(t, len(buff.String()), 0)
	expect(t, buff.String(), "router8-get")

	buff.Reset()

	req, err = http.NewRequest("POST", "http://localhost:8000/", nil)
	if err != nil {
		t.Error(err)
	}

	o.ServeHTTP(recorder, req)
	expect(t, recorder.Code, http.StatusOK)
	refute(t, len(buff.String()), 0)
	expect(t, buff.String(), "router8-post")
}
