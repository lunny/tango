// Copyright 2019 The Tango Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package tango

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

type Query1Action struct {
	Ctx
}

func (a *Query1Action) Get() string {
	v, _ := a.Queries().Int("test")
	return fmt.Sprintf("%d", v)
}

func TestQuery1(t *testing.T) {
	buff := bytes.NewBufferString("")
	recorder := httptest.NewRecorder()
	recorder.Body = buff

	o := Classic()
	o.Get("/", new(Query1Action))

	req, err := http.NewRequest("GET", "http://localhost:8000/?test=1", nil)
	if err != nil {
		t.Error(err)
	}

	o.ServeHTTP(recorder, req)
	expect(t, recorder.Code, http.StatusOK)
	refute(t, len(buff.String()), 0)
	expect(t, buff.String(), "1")
}

type Query2Action struct {
	Ctx
}

func (a *Query2Action) Get() string {
	v, _ := a.Queries().Int32("test")
	return fmt.Sprintf("%d", v)
}

func TestQuery2(t *testing.T) {
	buff := bytes.NewBufferString("")
	recorder := httptest.NewRecorder()
	recorder.Body = buff

	o := Classic()
	o.Get("/", new(Query2Action))

	req, err := http.NewRequest("GET", "http://localhost:8000/?test=1", nil)
	if err != nil {
		t.Error(err)
	}

	o.ServeHTTP(recorder, req)
	expect(t, recorder.Code, http.StatusOK)
	refute(t, len(buff.String()), 0)
	expect(t, buff.String(), "1")
}

type Query3Action struct {
	Ctx
}

func (a *Query3Action) Get() string {
	v, _ := a.Queries().Int64("test")
	return fmt.Sprintf("%d", v)
}

func TestQuery3(t *testing.T) {
	buff := bytes.NewBufferString("")
	recorder := httptest.NewRecorder()
	recorder.Body = buff

	o := Classic()
	o.Get("/", new(Query3Action))

	req, err := http.NewRequest("GET", "http://localhost:8000/?test=1", nil)
	if err != nil {
		t.Error(err)
	}

	o.ServeHTTP(recorder, req)
	expect(t, recorder.Code, http.StatusOK)
	refute(t, len(buff.String()), 0)
	expect(t, buff.String(), "1")
}

type Query4Action struct {
	Ctx
}

func (a *Query4Action) Get() string {
	v, _ := a.Queries().Uint("test")
	return fmt.Sprintf("%d", v)
}

func TestQuery4(t *testing.T) {
	buff := bytes.NewBufferString("")
	recorder := httptest.NewRecorder()
	recorder.Body = buff

	o := Classic()
	o.Get("/", new(Query4Action))

	req, err := http.NewRequest("GET", "http://localhost:8000/?test=1", nil)
	if err != nil {
		t.Error(err)
	}

	o.ServeHTTP(recorder, req)
	expect(t, recorder.Code, http.StatusOK)
	refute(t, len(buff.String()), 0)
	expect(t, buff.String(), "1")
}

type Query5Action struct {
	Ctx
}

func (a *Query5Action) Get() string {
	v, _ := a.Queries().Uint32("test")
	return fmt.Sprintf("%d", v)
}

func TestQuery5(t *testing.T) {
	buff := bytes.NewBufferString("")
	recorder := httptest.NewRecorder()
	recorder.Body = buff

	o := Classic()
	o.Get("/", new(Query5Action))

	req, err := http.NewRequest("GET", "http://localhost:8000/?test=1", nil)
	if err != nil {
		t.Error(err)
	}

	o.ServeHTTP(recorder, req)
	expect(t, recorder.Code, http.StatusOK)
	refute(t, len(buff.String()), 0)
	expect(t, buff.String(), "1")
}

type Query6Action struct {
	Ctx
}

func (a *Query6Action) Get() string {
	v, _ := a.Queries().Uint64("test")
	return fmt.Sprintf("%d", v)
}

func TestQuery6(t *testing.T) {
	buff := bytes.NewBufferString("")
	recorder := httptest.NewRecorder()
	recorder.Body = buff

	o := Classic()
	o.Get("/", new(Query6Action))

	req, err := http.NewRequest("GET", "http://localhost:8000/?test=1", nil)
	if err != nil {
		t.Error(err)
	}

	o.ServeHTTP(recorder, req)
	expect(t, recorder.Code, http.StatusOK)
	refute(t, len(buff.String()), 0)
	expect(t, buff.String(), "1")
}

type Query7Action struct {
	Ctx
}

func (a *Query7Action) Get() string {
	v, _ := a.Queries().Float32("test")
	return fmt.Sprintf("%.2f", v)
}

func TestQuery7(t *testing.T) {
	buff := bytes.NewBufferString("")
	recorder := httptest.NewRecorder()
	recorder.Body = buff

	o := Classic()
	o.Get("/", new(Query7Action))

	req, err := http.NewRequest("GET", "http://localhost:8000/?test=1", nil)
	if err != nil {
		t.Error(err)
	}

	o.ServeHTTP(recorder, req)
	expect(t, recorder.Code, http.StatusOK)
	refute(t, len(buff.String()), 0)
	expect(t, buff.String(), "1.00")
}

type Query8Action struct {
	Ctx
}

func (a *Query8Action) Get() string {
	v, _ := a.Queries().Float64("test")
	return fmt.Sprintf("%.2f", v)
}

func TestQuery8(t *testing.T) {
	buff := bytes.NewBufferString("")
	recorder := httptest.NewRecorder()
	recorder.Body = buff

	o := Classic()
	o.Get("/", new(Query8Action))

	req, err := http.NewRequest("GET", "http://localhost:8000/?test=1", nil)
	if err != nil {
		t.Error(err)
	}

	o.ServeHTTP(recorder, req)
	expect(t, recorder.Code, http.StatusOK)
	refute(t, len(buff.String()), 0)
	expect(t, buff.String(), "1.00")
}

type Query9Action struct {
	Ctx
}

func (a *Query9Action) Get() string {
	v, _ := a.Queries().Bool("test")
	return fmt.Sprintf("%v", v)
}

func TestQuery9(t *testing.T) {
	buff := bytes.NewBufferString("")
	recorder := httptest.NewRecorder()
	recorder.Body = buff

	o := Classic()
	o.Get("/", new(Query9Action))

	req, err := http.NewRequest("GET", "http://localhost:8000/?test=1", nil)
	if err != nil {
		t.Error(err)
	}

	o.ServeHTTP(recorder, req)
	expect(t, recorder.Code, http.StatusOK)
	refute(t, len(buff.String()), 0)
	expect(t, buff.String(), "true")
}

type Query10Action struct {
	Ctx
}

func (a *Query10Action) Get() string {
	v, _ := a.Queries().String("test")
	return fmt.Sprintf("%v", v)
}

func TestQuery10(t *testing.T) {
	buff := bytes.NewBufferString("")
	recorder := httptest.NewRecorder()
	recorder.Body = buff

	o := Classic()
	o.Get("/", new(Query10Action))

	req, err := http.NewRequest("GET", "http://localhost:8000/?test=1", nil)
	if err != nil {
		t.Error(err)
	}

	o.ServeHTTP(recorder, req)
	expect(t, recorder.Code, http.StatusOK)
	refute(t, len(buff.String()), 0)
	expect(t, buff.String(), "1")
}

type Query11Action struct {
	Ctx
}

func (a *Query11Action) Get() string {
	v := a.Queries().MustInt("test")
	return fmt.Sprintf("%d", v)
}

func TestQuery11(t *testing.T) {
	buff := bytes.NewBufferString("")
	recorder := httptest.NewRecorder()
	recorder.Body = buff

	o := Classic()
	o.Get("/", new(Query11Action))

	req, err := http.NewRequest("GET", "http://localhost:8000/?test=1", nil)
	if err != nil {
		t.Error(err)
	}

	o.ServeHTTP(recorder, req)
	expect(t, recorder.Code, http.StatusOK)
	refute(t, len(buff.String()), 0)
	expect(t, buff.String(), "1")
}

type Query12Action struct {
	Ctx
}

func (a *Query12Action) Get() string {
	v := a.Queries().MustInt32("test")
	return fmt.Sprintf("%d", v)
}

func TestQuery12(t *testing.T) {
	buff := bytes.NewBufferString("")
	recorder := httptest.NewRecorder()
	recorder.Body = buff

	o := Classic()
	o.Get("/", new(Query12Action))

	req, err := http.NewRequest("GET", "http://localhost:8000/?test=1", nil)
	if err != nil {
		t.Error(err)
	}

	o.ServeHTTP(recorder, req)
	expect(t, recorder.Code, http.StatusOK)
	refute(t, len(buff.String()), 0)
	expect(t, buff.String(), "1")
}

type Query13Action struct {
	Ctx
}

func (a *Query13Action) Get() string {
	v := a.Queries().MustInt64("test")
	return fmt.Sprintf("%d", v)
}

func TestQuery13(t *testing.T) {
	buff := bytes.NewBufferString("")
	recorder := httptest.NewRecorder()
	recorder.Body = buff

	o := Classic()
	o.Get("/", new(Query13Action))

	req, err := http.NewRequest("GET", "http://localhost:8000/?test=1", nil)
	if err != nil {
		t.Error(err)
	}

	o.ServeHTTP(recorder, req)
	expect(t, recorder.Code, http.StatusOK)
	refute(t, len(buff.String()), 0)
	expect(t, buff.String(), "1")
}

type Query14Action struct {
	Ctx
}

func (a *Query14Action) Get() string {
	v := a.Queries().MustUint("test")
	return fmt.Sprintf("%d", v)
}

func TestQuery14(t *testing.T) {
	buff := bytes.NewBufferString("")
	recorder := httptest.NewRecorder()
	recorder.Body = buff

	o := Classic()
	o.Get("/", new(Query14Action))

	req, err := http.NewRequest("GET", "http://localhost:8000/?test=1", nil)
	if err != nil {
		t.Error(err)
	}

	o.ServeHTTP(recorder, req)
	expect(t, recorder.Code, http.StatusOK)
	refute(t, len(buff.String()), 0)
	expect(t, buff.String(), "1")
}

type Query15Action struct {
	Ctx
}

func (a *Query15Action) Get() string {
	v := a.Queries().MustUint32("test")
	return fmt.Sprintf("%d", v)
}

func TestQuery15(t *testing.T) {
	buff := bytes.NewBufferString("")
	recorder := httptest.NewRecorder()
	recorder.Body = buff

	o := Classic()
	o.Get("/", new(Query15Action))

	req, err := http.NewRequest("GET", "http://localhost:8000/?test=1", nil)
	if err != nil {
		t.Error(err)
	}

	o.ServeHTTP(recorder, req)
	expect(t, recorder.Code, http.StatusOK)
	refute(t, len(buff.String()), 0)
	expect(t, buff.String(), "1")
}

type Query16Action struct {
	Ctx
}

func (a *Query16Action) Get() string {
	v := a.Queries().MustUint64("test")
	return fmt.Sprintf("%d", v)
}

func TestQuery16(t *testing.T) {
	buff := bytes.NewBufferString("")
	recorder := httptest.NewRecorder()
	recorder.Body = buff

	o := Classic()
	o.Get("/", new(Query16Action))

	req, err := http.NewRequest("GET", "http://localhost:8000/?test=1", nil)
	if err != nil {
		t.Error(err)
	}

	o.ServeHTTP(recorder, req)
	expect(t, recorder.Code, http.StatusOK)
	refute(t, len(buff.String()), 0)
	expect(t, buff.String(), "1")
}

type Query17Action struct {
	Ctx
}

func (a *Query17Action) Get() string {
	v := a.Queries().MustFloat32("test")
	return fmt.Sprintf("%.2f", v)
}

func TestQuery17(t *testing.T) {
	buff := bytes.NewBufferString("")
	recorder := httptest.NewRecorder()
	recorder.Body = buff

	o := Classic()
	o.Get("/", new(Query17Action))

	req, err := http.NewRequest("GET", "http://localhost:8000/?test=1", nil)
	if err != nil {
		t.Error(err)
	}

	o.ServeHTTP(recorder, req)
	expect(t, recorder.Code, http.StatusOK)
	refute(t, len(buff.String()), 0)
	expect(t, buff.String(), "1.00")
}

type Query18Action struct {
	Ctx
}

func (a *Query18Action) Get() string {
	v := a.Queries().MustFloat64("test")
	return fmt.Sprintf("%.2f", v)
}

func TestQuery18(t *testing.T) {
	buff := bytes.NewBufferString("")
	recorder := httptest.NewRecorder()
	recorder.Body = buff

	o := Classic()
	o.Get("/", new(Query18Action))

	req, err := http.NewRequest("GET", "http://localhost:8000/?test=1", nil)
	if err != nil {
		t.Error(err)
	}

	o.ServeHTTP(recorder, req)
	expect(t, recorder.Code, http.StatusOK)
	refute(t, len(buff.String()), 0)
	expect(t, buff.String(), "1.00")
}

type Query19Action struct {
	Ctx
}

func (a *Query19Action) Get() string {
	v := a.Queries().MustBool("test")
	return fmt.Sprintf("%v", v)
}

func TestQuery19(t *testing.T) {
	buff := bytes.NewBufferString("")
	recorder := httptest.NewRecorder()
	recorder.Body = buff

	o := Classic()
	o.Get("/", new(Query19Action))

	req, err := http.NewRequest("GET", "http://localhost:8000/?test=1", nil)
	if err != nil {
		t.Error(err)
	}

	o.ServeHTTP(recorder, req)
	expect(t, recorder.Code, http.StatusOK)
	refute(t, len(buff.String()), 0)
	expect(t, buff.String(), "true")
}

type Query20Action struct {
	Ctx
}

func (a *Query20Action) Get() string {
	v := a.Queries().MustString("test")
	return fmt.Sprintf("%s", v)
}

func TestQuery20(t *testing.T) {
	buff := bytes.NewBufferString("")
	recorder := httptest.NewRecorder()
	recorder.Body = buff

	o := Classic()
	o.Get("/", new(Query20Action))

	req, err := http.NewRequest("GET", "http://localhost:8000/?test=1", nil)
	if err != nil {
		t.Error(err)
	}

	o.ServeHTTP(recorder, req)
	expect(t, recorder.Code, http.StatusOK)
	refute(t, len(buff.String()), 0)
	expect(t, buff.String(), "1")
}

type Query21Action struct {
	Ctx
}

func (a *Query21Action) Get() string {
	v := a.QueryInt("test")
	return fmt.Sprintf("%d", v)
}

func TestQuery21(t *testing.T) {
	buff := bytes.NewBufferString("")
	recorder := httptest.NewRecorder()
	recorder.Body = buff

	o := Classic()
	o.Get("/", new(Query21Action))

	req, err := http.NewRequest("GET", "http://localhost:8000/?test=1", nil)
	if err != nil {
		t.Error(err)
	}

	o.ServeHTTP(recorder, req)
	expect(t, recorder.Code, http.StatusOK)
	refute(t, len(buff.String()), 0)
	expect(t, buff.String(), "1")
}

type Query22Action struct {
	Ctx
}

func (a *Query22Action) Get() string {
	v := a.QueryInt32("test")
	return fmt.Sprintf("%d", v)
}

func TestQuery22(t *testing.T) {
	buff := bytes.NewBufferString("")
	recorder := httptest.NewRecorder()
	recorder.Body = buff

	o := Classic()
	o.Get("/", new(Query22Action))

	req, err := http.NewRequest("GET", "http://localhost:8000/?test=1", nil)
	if err != nil {
		t.Error(err)
	}

	o.ServeHTTP(recorder, req)
	expect(t, recorder.Code, http.StatusOK)
	refute(t, len(buff.String()), 0)
	expect(t, buff.String(), "1")
}

type Query23Action struct {
	Ctx
}

func (a *Query23Action) Get() string {
	v := a.QueryInt64("test")
	return fmt.Sprintf("%d", v)
}

func TestQuery23(t *testing.T) {
	buff := bytes.NewBufferString("")
	recorder := httptest.NewRecorder()
	recorder.Body = buff

	o := Classic()
	o.Get("/", new(Query23Action))

	req, err := http.NewRequest("GET", "http://localhost:8000/?test=1", nil)
	if err != nil {
		t.Error(err)
	}

	o.ServeHTTP(recorder, req)
	expect(t, recorder.Code, http.StatusOK)
	refute(t, len(buff.String()), 0)
	expect(t, buff.String(), "1")
}

type Query24Action struct {
	Ctx
}

func (a *Query24Action) Get() string {
	v := a.QueryUint("test")
	return fmt.Sprintf("%d", v)
}

func TestQuery24(t *testing.T) {
	buff := bytes.NewBufferString("")
	recorder := httptest.NewRecorder()
	recorder.Body = buff

	o := Classic()
	o.Get("/", new(Query24Action))

	req, err := http.NewRequest("GET", "http://localhost:8000/?test=1", nil)
	if err != nil {
		t.Error(err)
	}

	o.ServeHTTP(recorder, req)
	expect(t, recorder.Code, http.StatusOK)
	refute(t, len(buff.String()), 0)
	expect(t, buff.String(), "1")
}

type Query25Action struct {
	Ctx
}

func (a *Query25Action) Get() string {
	v := a.QueryUint32("test")
	return fmt.Sprintf("%d", v)
}

func TestQuery25(t *testing.T) {
	buff := bytes.NewBufferString("")
	recorder := httptest.NewRecorder()
	recorder.Body = buff

	o := Classic()
	o.Get("/", new(Query25Action))

	req, err := http.NewRequest("GET", "http://localhost:8000/?test=1", nil)
	if err != nil {
		t.Error(err)
	}

	o.ServeHTTP(recorder, req)
	expect(t, recorder.Code, http.StatusOK)
	refute(t, len(buff.String()), 0)
	expect(t, buff.String(), "1")
}

type Query26Action struct {
	Ctx
}

func (a *Query26Action) Get() string {
	v := a.QueryUint64("test")
	return fmt.Sprintf("%d", v)
}

func TestQuery26(t *testing.T) {
	buff := bytes.NewBufferString("")
	recorder := httptest.NewRecorder()
	recorder.Body = buff

	o := Classic()
	o.Get("/", new(Query26Action))

	req, err := http.NewRequest("GET", "http://localhost:8000/?test=1", nil)
	if err != nil {
		t.Error(err)
	}

	o.ServeHTTP(recorder, req)
	expect(t, recorder.Code, http.StatusOK)
	refute(t, len(buff.String()), 0)
	expect(t, buff.String(), "1")
}

type Query27Action struct {
	Ctx
}

func (a *Query27Action) Get() string {
	v := a.QueryFloat32("test")
	return fmt.Sprintf("%.2f", v)
}

func TestQuery27(t *testing.T) {
	buff := bytes.NewBufferString("")
	recorder := httptest.NewRecorder()
	recorder.Body = buff

	o := Classic()
	o.Get("/", new(Query27Action))

	req, err := http.NewRequest("GET", "http://localhost:8000/?test=1", nil)
	if err != nil {
		t.Error(err)
	}

	o.ServeHTTP(recorder, req)
	expect(t, recorder.Code, http.StatusOK)
	refute(t, len(buff.String()), 0)
	expect(t, buff.String(), "1.00")
}

type Query28Action struct {
	Ctx
}

func (a *Query28Action) Get() string {
	v := a.QueryFloat64("test")
	return fmt.Sprintf("%.2f", v)
}

func TestQuery28(t *testing.T) {
	buff := bytes.NewBufferString("")
	recorder := httptest.NewRecorder()
	recorder.Body = buff

	o := Classic()
	o.Get("/", new(Query28Action))

	req, err := http.NewRequest("GET", "http://localhost:8000/?test=1", nil)
	if err != nil {
		t.Error(err)
	}

	o.ServeHTTP(recorder, req)
	expect(t, recorder.Code, http.StatusOK)
	refute(t, len(buff.String()), 0)
	expect(t, buff.String(), "1.00")
}

type Query29Action struct {
	Ctx
}

func (a *Query29Action) Get() string {
	v := a.QueryBool("test")
	return fmt.Sprintf("%v", v)
}

func TestQuery29(t *testing.T) {
	buff := bytes.NewBufferString("")
	recorder := httptest.NewRecorder()
	recorder.Body = buff

	o := Classic()
	o.Get("/", new(Query29Action))

	req, err := http.NewRequest("GET", "http://localhost:8000/?test=1", nil)
	if err != nil {
		t.Error(err)
	}

	o.ServeHTTP(recorder, req)
	expect(t, recorder.Code, http.StatusOK)
	refute(t, len(buff.String()), 0)
	expect(t, buff.String(), "true")
}

type Query30Action struct {
	Ctx
}

func (a *Query30Action) Get() string {
	v := a.Query("test")
	return fmt.Sprintf("%s", v)
}

func TestQuery30(t *testing.T) {
	buff := bytes.NewBufferString("")
	recorder := httptest.NewRecorder()
	recorder.Body = buff

	o := Classic()
	o.Get("/", new(Query30Action))

	req, err := http.NewRequest("GET", "http://localhost:8000/?test=1", nil)
	if err != nil {
		t.Error(err)
	}

	o.ServeHTTP(recorder, req)
	expect(t, recorder.Code, http.StatusOK)
	refute(t, len(buff.String()), 0)
	expect(t, buff.String(), "1")
}
