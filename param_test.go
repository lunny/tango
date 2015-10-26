// Copyright 2015 The Tango Authors. All rights reserved.
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

type Param5Action struct {
	Params
}

func (p *Param5Action) Get() string {
	i, _ := p.Params.Int(":name")
	return fmt.Sprintf("%d", i)
}

func TestParams7(t *testing.T) {
	buff := bytes.NewBufferString("")
	recorder := httptest.NewRecorder()
	recorder.Body = buff

	o := Classic()
	o.Get("/(:name[0-9]+)", new(Param5Action))

	req, err := http.NewRequest("GET", "http://localhost:8000/1", nil)
	if err != nil {
		t.Error(err)
	}

	o.ServeHTTP(recorder, req)
	expect(t, recorder.Code, http.StatusOK)
	refute(t, len(buff.String()), 0)
	expect(t, buff.String(), "1")
}

type Param6Action struct {
	Params
}

func (p *Param6Action) Get() string {
	i, _ := p.Params.Int32(":name")
	return fmt.Sprintf("%d", i)
}

func TestParams8(t *testing.T) {
	buff := bytes.NewBufferString("")
	recorder := httptest.NewRecorder()
	recorder.Body = buff

	o := Classic()
	o.Get("/(:name[0-9]+)", new(Param6Action))

	req, err := http.NewRequest("GET", "http://localhost:8000/1", nil)
	if err != nil {
		t.Error(err)
	}

	o.ServeHTTP(recorder, req)
	expect(t, recorder.Code, http.StatusOK)
	refute(t, len(buff.String()), 0)
	expect(t, buff.String(), "1")
}

type Param7Action struct {
	Params
}

func (p *Param7Action) Get() string {
	i, _ := p.Params.Int64(":name")
	return fmt.Sprintf("%d", i)
}

func TestParams9(t *testing.T) {
	buff := bytes.NewBufferString("")
	recorder := httptest.NewRecorder()
	recorder.Body = buff

	o := Classic()
	o.Get("/(:name[0-9]+)", new(Param7Action))

	req, err := http.NewRequest("GET", "http://localhost:8000/1", nil)
	if err != nil {
		t.Error(err)
	}

	o.ServeHTTP(recorder, req)
	expect(t, recorder.Code, http.StatusOK)
	refute(t, len(buff.String()), 0)
	expect(t, buff.String(), "1")
}

type Param8Action struct {
	Params
}

func (p *Param8Action) Get() string {
	i, _ := p.Params.Float32(":name")
	return fmt.Sprintf("%.2f", i)
}

func TestParams10(t *testing.T) {
	buff := bytes.NewBufferString("")
	recorder := httptest.NewRecorder()
	recorder.Body = buff

	o := Classic()
	o.Get("/(:name[0-9]+)", new(Param8Action))

	req, err := http.NewRequest("GET", "http://localhost:8000/1", nil)
	if err != nil {
		t.Error(err)
	}

	o.ServeHTTP(recorder, req)
	expect(t, recorder.Code, http.StatusOK)
	refute(t, len(buff.String()), 0)
	expect(t, buff.String(), "1.00")
}

type Param9Action struct {
	Params
}

func (p *Param9Action) Get() string {
	i, _ := p.Params.Float64(":name")
	return fmt.Sprintf("%.2f", i)
}

func TestParams11(t *testing.T) {
	buff := bytes.NewBufferString("")
	recorder := httptest.NewRecorder()
	recorder.Body = buff

	o := Classic()
	o.Get("/(:name[0-9]+)", new(Param9Action))

	req, err := http.NewRequest("GET", "http://localhost:8000/1", nil)
	if err != nil {
		t.Error(err)
	}

	o.ServeHTTP(recorder, req)
	expect(t, recorder.Code, http.StatusOK)
	refute(t, len(buff.String()), 0)
	expect(t, buff.String(), "1.00")
}

type Param10Action struct {
	Params
}

func (p *Param10Action) Get() string {
	i, _ := p.Params.Uint(":name")
	return fmt.Sprintf("%d", i)
}

func TestParams12(t *testing.T) {
	buff := bytes.NewBufferString("")
	recorder := httptest.NewRecorder()
	recorder.Body = buff

	o := Classic()
	o.Get("/(:name[0-9]+)", new(Param10Action))

	req, err := http.NewRequest("GET", "http://localhost:8000/1", nil)
	if err != nil {
		t.Error(err)
	}

	o.ServeHTTP(recorder, req)
	expect(t, recorder.Code, http.StatusOK)
	refute(t, len(buff.String()), 0)
	expect(t, buff.String(), "1")
}

type Param11Action struct {
	Params
}

func (p *Param11Action) Get() string {
	i, _ := p.Params.Uint32(":name")
	return fmt.Sprintf("%d", i)
}

func TestParams13(t *testing.T) {
	buff := bytes.NewBufferString("")
	recorder := httptest.NewRecorder()
	recorder.Body = buff

	o := Classic()
	o.Get("/(:name[0-9]+)", new(Param11Action))

	req, err := http.NewRequest("GET", "http://localhost:8000/1", nil)
	if err != nil {
		t.Error(err)
	}

	o.ServeHTTP(recorder, req)
	expect(t, recorder.Code, http.StatusOK)
	refute(t, len(buff.String()), 0)
	expect(t, buff.String(), "1")
}

type Param12Action struct {
	Params
}

func (p *Param12Action) Get() string {
	i, _ := p.Params.Uint64(":name")
	return fmt.Sprintf("%d", i)
}

func TestParams14(t *testing.T) {
	buff := bytes.NewBufferString("")
	recorder := httptest.NewRecorder()
	recorder.Body = buff

	o := Classic()
	o.Get("/(:name[0-9]+)", new(Param12Action))

	req, err := http.NewRequest("GET", "http://localhost:8000/1", nil)
	if err != nil {
		t.Error(err)
	}

	o.ServeHTTP(recorder, req)
	expect(t, recorder.Code, http.StatusOK)
	refute(t, len(buff.String()), 0)
	expect(t, buff.String(), "1")
}

type Param13Action struct {
	Params
}

func (p *Param13Action) Get() string {
	i, _ := p.Params.Bool(":name")
	return fmt.Sprintf("%v", i)
}

func TestParams15(t *testing.T) {
	buff := bytes.NewBufferString("")
	recorder := httptest.NewRecorder()
	recorder.Body = buff

	o := Classic()
	o.Get("/(:name[0-9]+)", new(Param13Action))

	req, err := http.NewRequest("GET", "http://localhost:8000/1", nil)
	if err != nil {
		t.Error(err)
	}

	o.ServeHTTP(recorder, req)
	expect(t, recorder.Code, http.StatusOK)
	refute(t, len(buff.String()), 0)
	expect(t, buff.String(), "true")
}

type Param14Action struct {
	Params
}

func (p *Param14Action) Get() string {
	i := p.Params.MustInt(":name")
	return fmt.Sprintf("%d", i)
}

func TestParams16(t *testing.T) {
	buff := bytes.NewBufferString("")
	recorder := httptest.NewRecorder()
	recorder.Body = buff

	o := Classic()
	o.Get("/(:name[0-9]+)", new(Param14Action))

	req, err := http.NewRequest("GET", "http://localhost:8000/1", nil)
	if err != nil {
		t.Error(err)
	}

	o.ServeHTTP(recorder, req)
	expect(t, recorder.Code, http.StatusOK)
	refute(t, len(buff.String()), 0)
	expect(t, buff.String(), "1")
}

type Param15Action struct {
	Params
}

func (p *Param15Action) Get() string {
	i := p.Params.MustInt32(":name")
	return fmt.Sprintf("%d", i)
}

func TestParams17(t *testing.T) {
	buff := bytes.NewBufferString("")
	recorder := httptest.NewRecorder()
	recorder.Body = buff

	o := Classic()
	o.Get("/(:name[0-9]+)", new(Param15Action))

	req, err := http.NewRequest("GET", "http://localhost:8000/1", nil)
	if err != nil {
		t.Error(err)
	}

	o.ServeHTTP(recorder, req)
	expect(t, recorder.Code, http.StatusOK)
	refute(t, len(buff.String()), 0)
	expect(t, buff.String(), "1")
}

type Param16Action struct {
	Params
}

func (p *Param16Action) Get() string {
	i := p.Params.MustInt64(":name")
	return fmt.Sprintf("%d", i)
}

func TestParams18(t *testing.T) {
	buff := bytes.NewBufferString("")
	recorder := httptest.NewRecorder()
	recorder.Body = buff

	o := Classic()
	o.Get("/(:name[0-9]+)", new(Param16Action))

	req, err := http.NewRequest("GET", "http://localhost:8000/1", nil)
	if err != nil {
		t.Error(err)
	}

	o.ServeHTTP(recorder, req)
	expect(t, recorder.Code, http.StatusOK)
	refute(t, len(buff.String()), 0)
	expect(t, buff.String(), "1")
}

type Param17Action struct {
	Params
}

func (p *Param17Action) Get() string {
	i := p.Params.MustFloat32(":name")
	return fmt.Sprintf("%.2f", i)
}

func TestParams19(t *testing.T) {
	buff := bytes.NewBufferString("")
	recorder := httptest.NewRecorder()
	recorder.Body = buff

	o := Classic()
	o.Get("/(:name[0-9]+)", new(Param17Action))

	req, err := http.NewRequest("GET", "http://localhost:8000/1", nil)
	if err != nil {
		t.Error(err)
	}

	o.ServeHTTP(recorder, req)
	expect(t, recorder.Code, http.StatusOK)
	refute(t, len(buff.String()), 0)
	expect(t, buff.String(), "1.00")
}

type Param18Action struct {
	Params
}

func (p *Param18Action) Get() string {
	i := p.Params.MustFloat64(":name")
	return fmt.Sprintf("%.2f", i)
}

func TestParams20(t *testing.T) {
	buff := bytes.NewBufferString("")
	recorder := httptest.NewRecorder()
	recorder.Body = buff

	o := Classic()
	o.Get("/(:name[0-9]+)", new(Param18Action))

	req, err := http.NewRequest("GET", "http://localhost:8000/1", nil)
	if err != nil {
		t.Error(err)
	}

	o.ServeHTTP(recorder, req)
	expect(t, recorder.Code, http.StatusOK)
	refute(t, len(buff.String()), 0)
	expect(t, buff.String(), "1.00")
}

type Param19Action struct {
	Params
}

func (p *Param19Action) Get() string {
	i := p.Params.MustUint(":name")
	return fmt.Sprintf("%d", i)
}

func TestParams21(t *testing.T) {
	buff := bytes.NewBufferString("")
	recorder := httptest.NewRecorder()
	recorder.Body = buff

	o := Classic()
	o.Get("/(:name[0-9]+)", new(Param19Action))

	req, err := http.NewRequest("GET", "http://localhost:8000/1", nil)
	if err != nil {
		t.Error(err)
	}

	o.ServeHTTP(recorder, req)
	expect(t, recorder.Code, http.StatusOK)
	refute(t, len(buff.String()), 0)
	expect(t, buff.String(), "1")
}

type Param20Action struct {
	Params
}

func (p *Param20Action) Get() string {
	i := p.Params.MustUint32(":name")
	return fmt.Sprintf("%d", i)
}

func TestParams22(t *testing.T) {
	buff := bytes.NewBufferString("")
	recorder := httptest.NewRecorder()
	recorder.Body = buff

	o := Classic()
	o.Get("/(:name[0-9]+)", new(Param20Action))

	req, err := http.NewRequest("GET", "http://localhost:8000/1", nil)
	if err != nil {
		t.Error(err)
	}

	o.ServeHTTP(recorder, req)
	expect(t, recorder.Code, http.StatusOK)
	refute(t, len(buff.String()), 0)
	expect(t, buff.String(), "1")
}

type Param21Action struct {
	Params
}

func (p *Param21Action) Get() string {
	i := p.Params.MustUint64(":name")
	return fmt.Sprintf("%d", i)
}

func TestParams23(t *testing.T) {
	buff := bytes.NewBufferString("")
	recorder := httptest.NewRecorder()
	recorder.Body = buff

	o := Classic()
	o.Get("/(:name[0-9]+)", new(Param21Action))

	req, err := http.NewRequest("GET", "http://localhost:8000/1", nil)
	if err != nil {
		t.Error(err)
	}

	o.ServeHTTP(recorder, req)
	expect(t, recorder.Code, http.StatusOK)
	refute(t, len(buff.String()), 0)
	expect(t, buff.String(), "1")
}

type Param22Action struct {
	Params
}

func (p *Param22Action) Get() string {
	i := p.Params.MustBool(":name")
	return fmt.Sprintf("%v", i)
}

func TestParams24(t *testing.T) {
	buff := bytes.NewBufferString("")
	recorder := httptest.NewRecorder()
	recorder.Body = buff

	o := Classic()
	o.Get("/(:name[0-9]+)", new(Param22Action))

	req, err := http.NewRequest("GET", "http://localhost:8000/1", nil)
	if err != nil {
		t.Error(err)
	}

	o.ServeHTTP(recorder, req)
	expect(t, recorder.Code, http.StatusOK)
	refute(t, len(buff.String()), 0)
	expect(t, buff.String(), "true")
}

type Param23Action struct {
	Params
}

func (p *Param23Action) Get() string {
	i, _ := p.Params.String(":name")
	return fmt.Sprintf("%v", i)
}

func TestParams25(t *testing.T) {
	buff := bytes.NewBufferString("")
	recorder := httptest.NewRecorder()
	recorder.Body = buff

	o := Classic()
	o.Get("/(:name[0-9]+)", new(Param23Action))

	req, err := http.NewRequest("GET", "http://localhost:8000/1", nil)
	if err != nil {
		t.Error(err)
	}

	o.ServeHTTP(recorder, req)
	expect(t, recorder.Code, http.StatusOK)
	refute(t, len(buff.String()), 0)
	expect(t, buff.String(), "1")
}

type Param24Action struct {
	Params
}

func (p *Param24Action) Get() string {
	i := p.Params.MustString(":name")
	return fmt.Sprintf("%v", i)
}

func TestParams26(t *testing.T) {
	buff := bytes.NewBufferString("")
	recorder := httptest.NewRecorder()
	recorder.Body = buff

	o := Classic()
	o.Get("/(:name[0-9]+)", new(Param24Action))

	req, err := http.NewRequest("GET", "http://localhost:8000/1", nil)
	if err != nil {
		t.Error(err)
	}

	o.ServeHTTP(recorder, req)
	expect(t, recorder.Code, http.StatusOK)
	refute(t, len(buff.String()), 0)
	expect(t, buff.String(), "1")
}

type Param25Action struct {
	Ctx
}

func (p *Param25Action) Get() string {
	i := p.ParamInt(":name")
	return fmt.Sprintf("%d", i)
}

func TestParams27(t *testing.T) {
	buff := bytes.NewBufferString("")
	recorder := httptest.NewRecorder()
	recorder.Body = buff

	o := Classic()
	o.Get("/(:name[0-9]+)", new(Param25Action))

	req, err := http.NewRequest("GET", "http://localhost:8000/1", nil)
	if err != nil {
		t.Error(err)
	}

	o.ServeHTTP(recorder, req)
	expect(t, recorder.Code, http.StatusOK)
	refute(t, len(buff.String()), 0)
	expect(t, buff.String(), "1")
}

type Param26Action struct {
	Ctx
}

func (p *Param26Action) Get() string {
	i := p.ParamInt32(":name")
	return fmt.Sprintf("%d", i)
}

func TestParams28(t *testing.T) {
	buff := bytes.NewBufferString("")
	recorder := httptest.NewRecorder()
	recorder.Body = buff

	o := Classic()
	o.Get("/(:name[0-9]+)", new(Param26Action))

	req, err := http.NewRequest("GET", "http://localhost:8000/1", nil)
	if err != nil {
		t.Error(err)
	}

	o.ServeHTTP(recorder, req)
	expect(t, recorder.Code, http.StatusOK)
	refute(t, len(buff.String()), 0)
	expect(t, buff.String(), "1")
}

type Param27Action struct {
	Ctx
}

func (p *Param27Action) Get() string {
	i := p.ParamInt64(":name")
	return fmt.Sprintf("%d", i)
}

func TestParams29(t *testing.T) {
	buff := bytes.NewBufferString("")
	recorder := httptest.NewRecorder()
	recorder.Body = buff

	o := Classic()
	o.Get("/(:name[0-9]+)", new(Param27Action))

	req, err := http.NewRequest("GET", "http://localhost:8000/1", nil)
	if err != nil {
		t.Error(err)
	}

	o.ServeHTTP(recorder, req)
	expect(t, recorder.Code, http.StatusOK)
	refute(t, len(buff.String()), 0)
	expect(t, buff.String(), "1")
}

type Param28Action struct {
	Ctx
}

func (p *Param28Action) Get() string {
	i := p.ParamFloat32(":name")
	return fmt.Sprintf("%.2f", i)
}

func TestParams30(t *testing.T) {
	buff := bytes.NewBufferString("")
	recorder := httptest.NewRecorder()
	recorder.Body = buff

	o := Classic()
	o.Get("/(:name[0-9]+)", new(Param28Action))

	req, err := http.NewRequest("GET", "http://localhost:8000/1", nil)
	if err != nil {
		t.Error(err)
	}

	o.ServeHTTP(recorder, req)
	expect(t, recorder.Code, http.StatusOK)
	refute(t, len(buff.String()), 0)
	expect(t, buff.String(), "1.00")
}

type Param29Action struct {
	Ctx
}

func (p *Param29Action) Get() string {
	i := p.ParamFloat64(":name")
	return fmt.Sprintf("%.2f", i)
}

func TestParams31(t *testing.T) {
	buff := bytes.NewBufferString("")
	recorder := httptest.NewRecorder()
	recorder.Body = buff

	o := Classic()
	o.Get("/(:name[0-9]+)", new(Param29Action))

	req, err := http.NewRequest("GET", "http://localhost:8000/1", nil)
	if err != nil {
		t.Error(err)
	}

	o.ServeHTTP(recorder, req)
	expect(t, recorder.Code, http.StatusOK)
	refute(t, len(buff.String()), 0)
	expect(t, buff.String(), "1.00")
}

type Param30Action struct {
	Ctx
}

func (p *Param30Action) Get() string {
	i := p.ParamUint(":name")
	return fmt.Sprintf("%d", i)
}

func TestParams32(t *testing.T) {
	buff := bytes.NewBufferString("")
	recorder := httptest.NewRecorder()
	recorder.Body = buff

	o := Classic()
	o.Get("/(:name[0-9]+)", new(Param30Action))

	req, err := http.NewRequest("GET", "http://localhost:8000/1", nil)
	if err != nil {
		t.Error(err)
	}

	o.ServeHTTP(recorder, req)
	expect(t, recorder.Code, http.StatusOK)
	refute(t, len(buff.String()), 0)
	expect(t, buff.String(), "1")
}

type Param31Action struct {
	Ctx
}

func (p *Param31Action) Get() string {
	i := p.ParamUint32(":name")
	return fmt.Sprintf("%d", i)
}

func TestParams33(t *testing.T) {
	buff := bytes.NewBufferString("")
	recorder := httptest.NewRecorder()
	recorder.Body = buff

	o := Classic()
	o.Get("/(:name[0-9]+)", new(Param31Action))

	req, err := http.NewRequest("GET", "http://localhost:8000/1", nil)
	if err != nil {
		t.Error(err)
	}

	o.ServeHTTP(recorder, req)
	expect(t, recorder.Code, http.StatusOK)
	refute(t, len(buff.String()), 0)
	expect(t, buff.String(), "1")
}

type Param32Action struct {
	Ctx
}

func (p *Param32Action) Get() string {
	i := p.ParamUint64(":name")
	return fmt.Sprintf("%d", i)
}

func TestParams34(t *testing.T) {
	buff := bytes.NewBufferString("")
	recorder := httptest.NewRecorder()
	recorder.Body = buff

	o := Classic()
	o.Get("/(:name[0-9]+)", new(Param32Action))

	req, err := http.NewRequest("GET", "http://localhost:8000/1", nil)
	if err != nil {
		t.Error(err)
	}

	o.ServeHTTP(recorder, req)
	expect(t, recorder.Code, http.StatusOK)
	refute(t, len(buff.String()), 0)
	expect(t, buff.String(), "1")
}

type Param33Action struct {
	Ctx
}

func (p *Param33Action) Get() string {
	i := p.ParamBool(":name")
	return fmt.Sprintf("%v", i)
}

func TestParams35(t *testing.T) {
	buff := bytes.NewBufferString("")
	recorder := httptest.NewRecorder()
	recorder.Body = buff

	o := Classic()
	o.Get("/(:name[0-9]+)", new(Param33Action))

	req, err := http.NewRequest("GET", "http://localhost:8000/1", nil)
	if err != nil {
		t.Error(err)
	}

	o.ServeHTTP(recorder, req)
	expect(t, recorder.Code, http.StatusOK)
	refute(t, len(buff.String()), 0)
	expect(t, buff.String(), "true")
}

type Param34Action struct {
	Ctx
}

func (p *Param34Action) Get() string {
	i := p.Param(":name")
	return fmt.Sprintf("%v", i)
}

func TestParams36(t *testing.T) {
	buff := bytes.NewBufferString("")
	recorder := httptest.NewRecorder()
	recorder.Body = buff

	o := Classic()
	o.Get("/(:name[0-9]+)", new(Param34Action))

	req, err := http.NewRequest("GET", "http://localhost:8000/1", nil)
	if err != nil {
		t.Error(err)
	}

	o.ServeHTTP(recorder, req)
	expect(t, recorder.Code, http.StatusOK)
	refute(t, len(buff.String()), 0)
	expect(t, buff.String(), "1")
}
