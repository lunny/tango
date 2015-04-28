// Copyright 2015 The Tango Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package tango

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"regexp"
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

type Regex1Action struct {
	Params
}

func (r *Regex1Action) Get() string {
	return r.Params.Get(":name")
}

func TestRouter9(t *testing.T) {
	buff := bytes.NewBufferString("")
	recorder := httptest.NewRecorder()
	recorder.Body = buff

	o := Classic()
	o.Get("/(:name[a-zA-Z]+)", new(Regex1Action))

	req, err := http.NewRequest("GET", "http://localhost:8000/foobar", nil)
	if err != nil {
		t.Error(err)
	}

	o.ServeHTTP(recorder, req)
	expect(t, recorder.Code, http.StatusOK)
	expect(t, buff.String(), "foobar")
	refute(t, len(buff.String()), 0)
}

var (
	parsedResult = map[string][]*node{
		"/": []*node{
			{content: "/", tp: snode},
		},
		"/static/css/bootstrap.css": []*node{
			{content: "/static", tp: snode},
			{content: "/css", tp: snode},
			{content: "/bootstrap.css", tp: snode},
		},
		"/:name": []*node{
			{content: "/", tp: snode},
			{content: ":name", tp: nnode},
		},
		"/sss:name": []*node{
			{content: "/sss", tp: snode},
			{content: ":name", tp: nnode},
		},
		"/(:name)": []*node{
			{content: "/", tp: snode},
			{content: ":name", tp: nnode},
		},
		"/(:name)/sss": []*node{
			{content: "/", tp: snode},
			{content: ":name", tp: nnode},
			{content: "/sss", tp: snode},
		},
		"/:name-:value": []*node{
			{content: "/", tp: snode},
			{content: ":name", tp: nnode},
			{content: "-", tp: snode},
			{content: ":value", tp: nnode},
		},
		"/(:name)ssss(:value)": []*node{
			{content: "/", tp: snode},
			{content: ":name", tp: nnode},
			{content: "ssss", tp: snode},
			{content: ":value", tp: nnode},
		},
		"/(:name[0-9]+)": []*node{
			{content: "/", tp: snode},
			{content: ":name", tp: rnode, regexp: regexp.MustCompile("([0-9]+)")},
		},
		"/*name": []*node{
			{content: "/", tp: snode},
			{content: "*name", tp: anode},
		},
		"/*name/ssss": []*node{
			{content: "/", tp: snode},
			{content: "*name", tp: anode},
			{content: "/ssss", tp: snode},
		},
		"/(*name)ssss": []*node{
			{content: "/", tp: snode},
			{content: "*name", tp: anode},
			{content: "ssss", tp: snode},
		},
		"/:name-(:name2[a-z]+)": []*node{
			{content: "/", tp: snode},
			{content: ":name", tp: nnode},
			{content: "-", tp: snode},
			{content: ":name2", tp: rnode, regexp: regexp.MustCompile("([a-z]+)")},
		},
	}
)

func TestParseNode(t *testing.T) {
	for p, r := range parsedResult {
		res := parseNodes(p)
		if len(r) != len(res) {
			t.Fatalf("%v 's result %v is not equal %v", p, r, res)
		}
		for i := 0; i < len(r); i++ {
			if r[i].content != res[i].content ||
				r[i].tp != res[i].tp {
				t.Fatalf("%v 's %d result %v is not equal %v, %v", p, i, r[i], res[i])
			}

			if r[i].tp != rnode {
				if r[i].regexp != nil {
					t.Fatalf("%v 's %d result %v is not equal %v, %v", p, i, r[i], res[i])
				}
			} else {
				if r[i].regexp == nil {
					t.Fatalf("%v 's %d result %v is not equal %v, %v", p, i, r[i], res[i])
				}
			}
		}
	}
}

type result struct {
	url    string
	match  bool
	params Params
}

var (
	matchResult = map[string][]result{
		"/": []result{
			{"/", true, Params{}},
			{"/s", false, Params{}},
			{"/123", false, Params{}},
		},
		"/ss/tt": []result{
			{"/ss/tt", true, Params{}},
			{"/s", false, Params{}},
			{"/ss", false, Params{}},
		},
		"/:name": []result{
			{"/s", true, Params{{":name", "s"}}},
			{"/", false, Params{}},
			{"/123/s", false, Params{}},
		},
		"/:name1/:name2/:name3": []result{
			{"/1/2/3", true, Params{{":name1", "1"}, {":name2", "2"}, {":name3", "3"}}},
			{"/1/2", false, Params{}},
			{"/1/2/3/", false, Params{}},
		},
		"/*name": []result{
			{"/s", true, Params{{"*name", "s"}}},
			{"/123/s", true, Params{{"*name", "123/s"}}},
			{"/", false, Params{}},
		},
		"/(*name)ssss": []result{
			{"/sssss", true, Params{{"*name", "s"}}},
			{"/123/ssss", true, Params{{"*name", "123/"}}},
			{"/", false, Params{}},
			{"/ss", false, Params{}},
		},
		"/111(*name)ssss": []result{
			{"/111sssss", true, Params{{"*name", "s"}}},
			{"/111/123/ssss", true, Params{{"*name", "/123/"}}},
			{"/", false, Params{}},
			{"/ss", false, Params{}},
		},
		"/(:name[0-9]+)": []result{
			{"/123", true, Params{{":name", "123"}}},
			{"/sss", false, Params{}},
		},
		"/ss(:name[0-9]+)": []result{
			{"/ss123", true, Params{{":name", "123"}}},
			{"/sss", false, Params{}},
		},
		"/ss(:name[0-9]+)tt": []result{
			{"/ss123tt", true, Params{{":name", "123"}}},
			{"/sss", false, Params{}},
		},
		"/:name1-(:name2[0-9]+)": []result{
			{"/ss-123", true, Params{{":name1", "ss"}, {":name2", "123"}}},
			{"/sss", false, Params{}},
		},
		"/(:name1)00(:name2[0-9]+)": []result{
			{"/ss00123", true, Params{{":name1", "ss"}, {":name2", "123"}}},
			{"/sss", false, Params{}},
		},
		"/(:name1)!(:name2[0-9]+)!(:name3.*)": []result{
			{"/ss!123!456", true, Params{{":name1", "ss"}, {":name2", "123"}, {":name3", "456"}}},
			{"/sss", false, Params{}},
		},
	}
)

type Action struct {
}

func (Action) Get() string {
	return "get"
}

func TestRouterSingle(t *testing.T) {
	for k, m := range matchResult {
		r := New()
		r.Route("GET", k, new(Action))

		for _, res := range m {
			handler, params := r.Match(res.url, "GET")
			if res.match {
				if handler == nil {
					t.Fatal(k, res, "handler", handler, "should not be nil")
				}
				for i, v := range params {
					if res.params[i].Name != v.Name {
						t.Fatal(k, res, "params name", v, "not equal", res.params[i])
					}
					if res.params[i].Value != v.Value {
						t.Fatal(k, res, "params value", v, "not equal", res.params[i])
					}
				}
			} else {
				if handler != nil {
					t.Fatal(k, res, "handler", handler, "should be nil")
				}
			}
		}
	}
}

type testCase struct {
	routers []string
	results []result
}

var (
	matchResult2 = []testCase{
		{
			[]string{"/"},
			[]result{
				{"/", true, Params{}},
				{"/s", false, Params{}},
				{"/123", false, Params{}},
			},
		},

		{
			[]string{"/admin", "/:name"},
			[]result{
				{"/", false, Params{}},
				{"/admin", true, Params{}},
				{"/s", true, Params{param{":name", "s"}}},
				{"/123", true, Params{param{":name", "123"}}},
			},
		},

		{
			[]string{"/:name", "/admin"},
			[]result{
				{"/", false, Params{}},
				{"/admin", true, Params{}},
				{"/s", true, Params{param{":name", "s"}}},
				{"/123", true, Params{param{":name", "123"}}},
			},
		},

		{
			[]string{"/admin", "/*name"},
			[]result{
				{"/", false, Params{}},
				{"/admin", true, Params{}},
				{"/s", true, Params{param{"*name", "s"}}},
				{"/123", true, Params{param{"*name", "123"}}},
			},
		},

		{
			[]string{"/*name", "/admin"},
			[]result{
				{"/", false, Params{}},
				{"/admin", true, Params{}},
				{"/s", true, Params{param{"*name", "s"}}},
				{"/123", true, Params{param{"*name", "123"}}},
			},
		},

		{
			[]string{"/*name", "/:name"},
			[]result{
				{"/", false, Params{}},
				{"/s", true, Params{param{"*name", "s"}}},
				{"/123", true, Params{param{"*name", "123"}}},
			},
		},

		{
			[]string{"/:name", "/*name"},
			[]result{
				{"/", false, Params{}},
				{"/s", true, Params{param{":name", "s"}}},
				{"/123", true, Params{param{":name", "123"}}},
				{"/123/1", true, Params{param{"*name", "123/1"}}},
			},
		},

		{
			[]string{"/*name", "/*name/123"},
			[]result{
				{"/", false, Params{}},
				{"/123", true, Params{param{"*name", "123"}}},
				{"/s", true, Params{param{"*name", "s"}}},
				{"/abc/123", true, Params{param{"*name", "abc"}}},
				{"/name1/name2/123", true, Params{param{"*name", "name1/name2"}}},
			},
		},

		{
			[]string{"/admin/ui", "/*name", "/:name"},
			[]result{
				{"/", false, Params{}},
				{"/admin/ui", true, Params{}},
				{"/s", true, Params{param{"*name", "s"}}},
				{"/123", true, Params{param{"*name", "123"}}},
			},
		},

		{
			[]string{"/(:id[0-9]+)", "/(:id[0-9]+)/edit", "/(:id[0-9]+)/del"},
			[]result{
				{"/1", true, Params{param{":id", "1"}}},
				{"/admin/ui", false, Params{}},
				{"/2/edit", true, Params{param{":id", "2"}}},
				{"/3/del", true, Params{param{":id", "3"}}},
			},
		},

		{
			[]string{"/admin/ui", "/:name1/:name2"},
			[]result{
				{"/", false, Params{}},
				{"/admin/ui", true, Params{}},
				{"/s", false, Params{}},
				{"/admin/ui2", true, Params{param{":name1", "admin"}, param{":name2", "ui2"}}},
				{"/123/s", true, Params{param{":name1", "123"}, param{":name2", "s"}}},
			},
		},
	}
)

func TestRouterMultiple(t *testing.T) {
	for _, kase := range matchResult2 {
		r := New()
		for _, k := range kase.routers {
			r.Route("GET", k, new(Action))
		}

		for _, res := range kase.results {
			handler, params := r.Match(res.url, "GET")
			if res.match {
				if handler == nil {
					t.Fatal(kase.routers, res, "handler", handler, "should not be nil")
				}

				if len(res.params) != len(params) {
					t.Fatal(kase.routers, res, "params", params, "not equal", res.params)
				}
				for i, v := range params {
					if res.params[i].Name != v.Name {
						t.Fatal(kase.routers, res, "params name", v, "not equal", res.params[i])
					}
					if res.params[i].Value != v.Value {
						t.Fatal(kase.routers, res, "params value", v, "not equal", res.params[i])
					}
				}
			} else {
				if handler != nil {
					t.Fatal(kase.routers, res, "handler", handler, "should be nil")
				}
			}
		}
	}
}

func TestRouter10(t *testing.T) {
	buff := bytes.NewBufferString("")
	recorder := httptest.NewRecorder()
	recorder.Body = buff

	r := New()
	r.Get("/", func(ctx *Context) {
		ctx.Write([]byte("test"))
	})

	req, err := http.NewRequest("GET", "http://localhost:8000/", nil)
	if err != nil {
		t.Error(err)
	}

	r.ServeHTTP(recorder, req)
	expect(t, recorder.Code, http.StatusOK)
	expect(t, buff.String(), "test")
	refute(t, len(buff.String()), 0)
}
