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

func TestStatic(t *testing.T) {
	buff := bytes.NewBufferString("")
	recorder := httptest.NewRecorder()
	recorder.Body = buff

	tg := New()
	tg.Use(Static())

	req, err := http.NewRequest("GET", "http://localhost:8000/test.html", nil)
	if err != nil {
		t.Error(err)
	}

	tg.ServeHTTP(recorder, req)
	expect(t, recorder.Code, http.StatusOK)
	refute(t, len(buff.String()), 0)
	expect(t, buff.String(), "hello tango")

	buff.Reset()

	req, err = http.NewRequest("GET", "http://localhost:8000/", nil)
	if err != nil {
		t.Error(err)
	}

	tg.ServeHTTP(recorder, req)
	expect(t, recorder.Code, http.StatusOK)
	refute(t, len(buff.String()), 0)
	expect(t, buff.String(), "this is index.html")
}

func TestStatic2(t *testing.T) {
	buff := bytes.NewBufferString("")
	recorder := httptest.NewRecorder()
	recorder.Body = buff

	tg := New()
	tg.Use(Static())

	req, err := http.NewRequest("GET", "http://localhost:8000/test.png", nil)
	if err != nil {
		t.Error(err)
	}

	tg.ServeHTTP(recorder, req)
	expect(t, recorder.Code, http.StatusNotFound)
}

func TestStatic3(t *testing.T) {
	buff := bytes.NewBufferString("")
	recorder := httptest.NewRecorder()
	recorder.Body = buff

	tg := New()
	tg.Use(Static(StaticOptions{
		Prefix:   "/public",
		RootPath: "./public",
	}))

	req, err := http.NewRequest("GET", "http://localhost:8000/public/test.html", nil)
	if err != nil {
		t.Error(err)
	}

	tg.ServeHTTP(recorder, req)
	expect(t, recorder.Code, http.StatusOK)
	expect(t, buff.String(), "hello tango")
}

func TestStatic4(t *testing.T) {
	buff := bytes.NewBufferString("")
	recorder := httptest.NewRecorder()
	recorder.Body = buff

	tg := New()
	tg.Use(Static(StaticOptions{
		Prefix:   "/public",
		RootPath: "./public",
	}))

	req, err := http.NewRequest("GET", "http://localhost:8000/public/t.html", nil)
	if err != nil {
		t.Error(err)
	}

	tg.ServeHTTP(recorder, req)
	expect(t, recorder.Code, http.StatusNotFound)
	expect(t, buff.String(), NotFound().Error())
}

func TestStatic5(t *testing.T) {
	buff := bytes.NewBufferString("")
	recorder := httptest.NewRecorder()
	recorder.Body = buff

	tg := New()
	tg.Use(Static(StaticOptions{
		Prefix:     "/public",
		RootPath:   "./public",
		ListDir:    true,
		IndexFiles: []string{"a.html"},
	}))

	req, err := http.NewRequest("GET", "http://localhost:8000/public/", nil)
	if err != nil {
		t.Error(err)
	}

	tg.ServeHTTP(recorder, req)
	expect(t, recorder.Code, http.StatusOK)
	//expect(t, buff.String(), NotFound().Error())
}
