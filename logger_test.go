// Copyright 2015 The Tango Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package tango

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"gitea.com/lunny/log"
)

func TestLogger1(t *testing.T) {
	buff := bytes.NewBufferString("")
	recorder := httptest.NewRecorder()

	n := NewWithLog(log.New(buff, "[tango] ", 0))
	n.Use(Logging())
	n.UseHandler(http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		rw.WriteHeader(http.StatusNotFound)
	}))

	req, err := http.NewRequest("GET", "http://localhost:3000/foobar", nil)
	if err != nil {
		t.Error(err)
	}

	n.ServeHTTP(recorder, req)
	expect(t, recorder.Code, http.StatusNotFound)
	refute(t, len(buff.String()), 0)
}

type LoggerAction struct {
	Log
}

func (l *LoggerAction) Get() string {
	l.Warn("this is a warn")
	l.Warnf("This is a %s", "warnf")
	l.Error("this is an error")
	l.Errorf("This is a %s", "errorf")
	l.Infof("This is a %s", "infof")
	l.Debugf("This is a %s", "debuf")
	return "log"
}

func TestLogger2(t *testing.T) {
	buff := bytes.NewBufferString("")
	recorder := httptest.NewRecorder()
	recorder.Body = buff

	n := Classic()
	n.Get("/", new(LoggerAction))

	req, err := http.NewRequest("GET", "http://localhost:3000/", nil)
	if err != nil {
		t.Error(err)
	}

	n.ServeHTTP(recorder, req)
	expect(t, recorder.Code, http.StatusOK)
	refute(t, len(buff.String()), 0)
	expect(t, buff.String(), "log")
}

func TestLogger3(t *testing.T) {
	buff := bytes.NewBufferString("")
	recorder := httptest.NewRecorder()
	recorder.Body = buff

	logger := NewCompositeLogger(log.Std, log.New(log.NewFileWriter(log.FileOptions{
		Dir:    "./",
		ByType: log.ByDay,
	}), "file", log.Ldefault()))

	n := Classic(logger)
	n.Get("/", new(LoggerAction))

	req, err := http.NewRequest("GET", "http://localhost:3000/", nil)
	if err != nil {
		t.Error(err)
	}

	n.ServeHTTP(recorder, req)
	expect(t, recorder.Code, http.StatusOK)
	refute(t, len(buff.String()), 0)
	expect(t, buff.String(), "log")
}

type Logger4Action struct {
}

func (l *Logger4Action) Get() {
}

func TestLogger4(t *testing.T) {
	buff := bytes.NewBufferString("")
	recorder := httptest.NewRecorder()
	recorder.Body = buff

	n := Classic()
	n.Get("/", new(Logger4Action))

	req, err := http.NewRequest("GET", "http://localhost:3000/", nil)
	if err != nil {
		t.Error(err)
	}

	n.ServeHTTP(recorder, req)
	expect(t, recorder.Code, http.StatusOK)
}
