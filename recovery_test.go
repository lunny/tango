// Copyright 2015 The Tango Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package tango

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func TestRecovery(t *testing.T) {
	buff := bytes.NewBufferString("")
	recorder := httptest.NewRecorder()
	recorder.Body = buff

	n := NewWithLog(NewLogger(os.Stdout))
	n.Use(Recovery(true))
	n.UseHandler(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		panic("here is a panic!")
	}))

	req, err := http.NewRequest("GET", "http://localhost:8000/", nil)
	if err != nil {
		t.Error(err)
	}

	n.ServeHTTP(recorder, req)
	expect(t, recorder.Code, http.StatusInternalServerError)
	refute(t, recorder.Body.Len(), 0)
	refute(t, len(buff.String()), 0)
}
