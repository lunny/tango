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
	g.Get("/1", func() string {
		return "/1"
	})
	g.Post("/2", func() string {
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
	req, err = http.NewRequest("POST", "http://localhost:8000/api/2", nil)
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
	        cg.Get("/1", func() string {
		    return "/1"
	        })
	        cg.Post("/2", func() string {
		    return "/2"
	        })
	    })
	})

	req, err := http.NewRequest("GET", "http://localhost:8000/api/v1/1", nil)
	if err != nil {
		t.Error(err)
	}

	o.ServeHTTP(recorder, req)
	expect(t, recorder.Code, http.StatusOK)
	refute(t, len(buff.String()), 0)
	expect(t, buff.String(), "/1")

	buff.Reset()
	req, err = http.NewRequest("POST", "http://localhost:8000/api/v1/2", nil)
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
		g.Get("/api/1", func() string {
			return "/1"
		})
		g.Post("/api/2", func() string {
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