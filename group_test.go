// Copyright 2015 The Tango Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package tango

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGroup1(t *testing.T) {
	buff := bytes.NewBufferString("")
	recorder := httptest.NewRecorder()
	recorder.Body = buff

	o := Classic()
	o.Group("/api", func(g *Group) {
		g.Get("/1", func() string {
			return "/1"
		})
		g.Post("/2", func() string {
			return "/2"
		})
	})

	req, err := http.NewRequest("GET", "http://localhost:8000/api/1", nil)
	if err != nil {
		t.Error(err)
	}

	o.ServeHTTP(recorder, req)
	expect(t, recorder.Code, http.StatusOK)
	refute(t, len(buff.String()), 0)
	expect(t, buff.String(), "/1")

	buff.Reset()
	req, err = http.NewRequest("POST", "http://localhost:8000/api/2", nil)
	if err != nil {
		t.Error(err)
	}

	o.ServeHTTP(recorder, req)
	expect(t, recorder.Code, http.StatusOK)
	refute(t, len(buff.String()), 0)
	expect(t, buff.String(), "/2")
}

func TestGroup2(t *testing.T) {
	buff := bytes.NewBufferString("")
	recorder := httptest.NewRecorder()
	recorder.Body = buff

	g := NewGroup()
	g.Any("/1", func() string {
		return "/1"
	})
	g.Options("/2", func() string {
		return "/2"
	})

	o := Classic()
	o.Group("/api", g)

	req, err := http.NewRequest("GET", "http://localhost:8000/api/1", nil)
	if err != nil {
		t.Error(err)
	}

	o.ServeHTTP(recorder, req)
	expect(t, recorder.Code, http.StatusOK)
	refute(t, len(buff.String()), 0)
	expect(t, buff.String(), "/1")

	buff.Reset()
	req, err = http.NewRequest("OPTIONS", "http://localhost:8000/api/2", nil)
	if err != nil {
		t.Error(err)
	}

	o.ServeHTTP(recorder, req)
	expect(t, recorder.Code, http.StatusOK)
	refute(t, len(buff.String()), 0)
	expect(t, buff.String(), "/2")
}

func TestGroup3(t *testing.T) {
	buff := bytes.NewBufferString("")
	recorder := httptest.NewRecorder()
	recorder.Body = buff

	o := Classic()
	o.Group("/api", func(g *Group) {
		g.Group("/v1", func(cg *Group) {
			cg.Trace("/1", func() string {
				return "/1"
			})
			cg.Patch("/2", func() string {
				return "/2"
			})
		})
	})

	req, err := http.NewRequest("TRACE", "http://localhost:8000/api/v1/1", nil)
	if err != nil {
		t.Error(err)
	}

	o.ServeHTTP(recorder, req)
	expect(t, recorder.Code, http.StatusOK)
	refute(t, len(buff.String()), 0)
	expect(t, buff.String(), "/1")

	buff.Reset()
	req, err = http.NewRequest("PATCH", "http://localhost:8000/api/v1/2", nil)
	if err != nil {
		t.Error(err)
	}

	o.ServeHTTP(recorder, req)
	expect(t, recorder.Code, http.StatusOK)
	refute(t, len(buff.String()), 0)
	expect(t, buff.String(), "/2")
}

func TestGroup4(t *testing.T) {
	buff := bytes.NewBufferString("")
	recorder := httptest.NewRecorder()
	recorder.Body = buff

	o := Classic()
	o.Group("", func(g *Group) {
		g.Delete("/api/1", func() string {
			return "/1"
		})
		g.Head("/api/2", func() string {
			return "/2"
		})
	})

	req, err := http.NewRequest("DELETE", "http://localhost:8000/api/1", nil)
	if err != nil {
		t.Error(err)
	}

	o.ServeHTTP(recorder, req)
	expect(t, recorder.Code, http.StatusOK)
	refute(t, len(buff.String()), 0)
	expect(t, buff.String(), "/1")

	buff.Reset()
	req, err = http.NewRequest("HEAD", "http://localhost:8000/api/2", nil)
	if err != nil {
		t.Error(err)
	}

	o.ServeHTTP(recorder, req)
	expect(t, recorder.Code, http.StatusOK)
	refute(t, len(buff.String()), 0)
	expect(t, buff.String(), "/2")
}

func TestGroup5(t *testing.T) {
	buff := bytes.NewBufferString("")
	recorder := httptest.NewRecorder()
	recorder.Body = buff

	o := Classic()
	var handlerGroup bool
	o.Group("/api", func(g *Group) {
		g.Use(HandlerFunc(func(ctx *Context) {
			handlerGroup = true
			ctx.Next()
		}))
		g.Put("/1", func() string {
			return "/1"
		})
	})
	o.Post("/2", func() string {
		return "/2"
	})

	req, err := http.NewRequest("PUT", "http://localhost:8000/api/1", nil)
	if err != nil {
		t.Error(err)
	}

	o.ServeHTTP(recorder, req)
	expect(t, recorder.Code, http.StatusOK)
	refute(t, len(buff.String()), 0)
	expect(t, buff.String(), "/1")
	expect(t, handlerGroup, true)

	handlerGroup = false
	buff.Reset()
	req, err = http.NewRequest("POST", "http://localhost:8000/2", nil)
	if err != nil {
		t.Error(err)
	}

	o.ServeHTTP(recorder, req)
	expect(t, recorder.Code, http.StatusOK)
	refute(t, len(buff.String()), 0)
	expect(t, buff.String(), "/2")
	expect(t, handlerGroup, false)
}

func TestGroup6(t *testing.T) {
	buff := bytes.NewBufferString("")
	recorder := httptest.NewRecorder()
	recorder.Body = buff

	var handlerGroup bool
	g := NewGroup()
	g.Use(HandlerFunc(func(ctx *Context) {
		handlerGroup = true
		ctx.Next()
	}))
	g.Get("/1", func() string {
		return "/1"
	})

	o := Classic()
	o.Group("/api", g)
	o.Post("/2", func() string {
		return "/2"
	})

	req, err := http.NewRequest("GET", "http://localhost:8000/api/1", nil)
	if err != nil {
		t.Error(err)
	}

	o.ServeHTTP(recorder, req)
	expect(t, recorder.Code, http.StatusOK)
	refute(t, len(buff.String()), 0)
	expect(t, buff.String(), "/1")
	expect(t, handlerGroup, true)

	handlerGroup = false
	buff.Reset()
	req, err = http.NewRequest("POST", "http://localhost:8000/2", nil)
	if err != nil {
		t.Error(err)
	}

	o.ServeHTTP(recorder, req)
	expect(t, recorder.Code, http.StatusOK)
	refute(t, len(buff.String()), 0)
	expect(t, buff.String(), "/2")
	expect(t, handlerGroup, false)
}

func TestGroup7(t *testing.T) {
	buff := bytes.NewBufferString("")
	recorder := httptest.NewRecorder()
	recorder.Body = buff

	var isPanic bool
	defer func() {
		if err := recover(); err != nil {
			isPanic = true
		}
		expect(t, isPanic, true)
	}()

	o := Classic()
	o.Group("/api", func() {
	})

	req, err := http.NewRequest("GET", "http://localhost:8000/api/1", nil)
	if err != nil {
		t.Error(err)
	}

	o.ServeHTTP(recorder, req)
}
