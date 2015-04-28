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

type PrefixAction struct {
}

func (PrefixAction) Get() string {
	return "Prefix"
}

type NoPrefixAction struct {
}

func (NoPrefixAction) Get() string {
	return "NoPrefix"
}

func TestPrefix(t *testing.T) {
	buff := bytes.NewBufferString("")
	recorder := httptest.NewRecorder()
	recorder.Body = buff

	var isPrefix bool

	o := Classic()
	o.Use(Prefix("/prefix", HandlerFunc(func(ctx *Context) {
		isPrefix = true
		ctx.Next()
	})))
	o.Get("/prefix/t", new(PrefixAction))
	o.Get("/t", new(NoPrefixAction))

	req, err := http.NewRequest("GET", "http://localhost:8000/prefix/t", nil)
	if err != nil {
		t.Error(err)
	}

	o.ServeHTTP(recorder, req)
	expect(t, recorder.Code, http.StatusOK)
	refute(t, len(buff.String()), 0)
	expect(t, buff.String(), "Prefix")
	expect(t, isPrefix, true)

	isPrefix = false
	buff.Reset()

	req, err = http.NewRequest("GET", "http://localhost:8000/t", nil)
	if err != nil {
		t.Error(err)
	}

	o.ServeHTTP(recorder, req)
	expect(t, recorder.Code, http.StatusOK)
	refute(t, len(buff.String()), 0)
	expect(t, buff.String(), "NoPrefix")
	expect(t, isPrefix, false)
}
