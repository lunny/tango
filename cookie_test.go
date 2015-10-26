// Copyright 2015 The Tango Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package tango

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestCookie1(t *testing.T) {
	buff := bytes.NewBufferString("")
	recorder := httptest.NewRecorder()
	recorder.Body = buff

	o := Classic()
	o.Get("/", func(ctx *Context) string {
		ck := ctx.Cookies().Get("name")
		if ck != nil {
			return ck.Value
		}
		return ""
	})

	req, err := http.NewRequest("GET", "http://localhost:8000/", nil)
	if err != nil {
		t.Error(err)
	}
	req.AddCookie(NewCookie("name", "test"))

	o.ServeHTTP(recorder, req)
	expect(t, recorder.Code, http.StatusOK)
	refute(t, len(buff.String()), 0)
	expect(t, buff.String(), "test")
}

func TestCookie2(t *testing.T) {
	buff := bytes.NewBufferString("")
	recorder := httptest.NewRecorder()
	recorder.Body = buff

	o := Classic()
	o.Get("/", func(ctx *Context) string {
		ctx.Cookies().Set(NewCookie("name", "test"))
		return "test"
	})

	req, err := http.NewRequest("GET", "http://localhost:8000/", nil)
	if err != nil {
		t.Error(err)
	}

	o.ServeHTTP(recorder, req)
	expect(t, recorder.Code, http.StatusOK)
	refute(t, len(buff.String()), 0)
	expect(t, buff.String(), "test")
	expect(t, strings.Split(recorder.Header().Get("Set-Cookie"), ";")[0], "name=test")
}

func TestCookie3(t *testing.T) {
	buff := bytes.NewBufferString("")
	recorder := httptest.NewRecorder()
	recorder.Body = buff

	o := Classic()
	o.Get("/", func(ctx *Context) string {
		ctx.Cookies().Expire("expire", time.Now())
		return "test"
	})

	req, err := http.NewRequest("GET", "http://localhost:8000/", nil)
	if err != nil {
		t.Error(err)
	}

	req.AddCookie(NewCookie("expire", "test"))

	o.ServeHTTP(recorder, req)
	expect(t, recorder.Code, http.StatusOK)
	refute(t, len(buff.String()), 0)
	expect(t, buff.String(), "test")
	fmt.Println(recorder.Header().Get("Set-Cookie"))
	expect(t, strings.Split(recorder.Header().Get("Set-Cookie"), ";")[0], "expire=test")
}

func TestCookie4(t *testing.T) {
	buff := bytes.NewBufferString("")
	recorder := httptest.NewRecorder()
	recorder.Body = buff

	o := Classic()
	o.Get("/", func(ctx *Context) string {
		ctx.Cookies().Del("ttttt")
		return "test"
	})

	req, err := http.NewRequest("GET", "http://localhost:8000/", nil)
	if err != nil {
		t.Error(err)
	}

	req.AddCookie(NewCookie("ttttt", "test"))

	o.ServeHTTP(recorder, req)
	expect(t, recorder.Code, http.StatusOK)
	refute(t, len(buff.String()), 0)
	expect(t, buff.String(), "test")
	expect(t, recorder.Header().Get("Set-Cookie"), "ttttt=test; Max-Age=0")
}

func TestCookie5(t *testing.T) {
	buff := bytes.NewBufferString("")
	recorder := httptest.NewRecorder()
	recorder.Body = buff

	o := Classic()
	o.Get("/", func(ctx *Context) string {
		ck := ctx.SecureCookies("sssss").Get("name")
		if ck != nil {
			return ck.Value
		}
		return ""
	})

	req, err := http.NewRequest("GET", "http://localhost:8000/", nil)
	if err != nil {
		t.Error(err)
	}
	req.AddCookie(NewSecureCookie("sssss", "name", "test"))

	o.ServeHTTP(recorder, req)
	expect(t, recorder.Code, http.StatusOK)
	refute(t, len(buff.String()), 0)
	expect(t, buff.String(), "test")
}

func TestCookie6(t *testing.T) {
	buff := bytes.NewBufferString("")
	recorder := httptest.NewRecorder()
	recorder.Body = buff

	o := Classic()
	o.Get("/", func(ctx *Context) string {
		ctx.SecureCookies("ttttt").Set(NewSecureCookie("ttttt", "name", "test"))
		return "test"
	})

	req, err := http.NewRequest("GET", "http://localhost:8000/", nil)
	if err != nil {
		t.Error(err)
	}

	o.ServeHTTP(recorder, req)
	expect(t, recorder.Code, http.StatusOK)
	refute(t, len(buff.String()), 0)
	expect(t, buff.String(), "test")

	r := strings.Split(recorder.Header().Get("Set-Cookie"), ";")[0]
	s := strings.SplitN(r, "=", 2)
	name, value := s[0], s[1]
	expect(t, name, "name")
	expect(t, parseSecureCookie("ttttt", value), "test")
}

func TestCookie7(t *testing.T) {
	buff := bytes.NewBufferString("")
	recorder := httptest.NewRecorder()
	recorder.Body = buff

	o := Classic()
	o.Get("/", func(ctx *Context) string {
		ctx.SecureCookies("ttttt").Expire("expire", time.Now())
		return "test"
	})

	req, err := http.NewRequest("GET", "http://localhost:8000/", nil)
	if err != nil {
		t.Error(err)
	}

	req.AddCookie(NewSecureCookie("ttttt", "expire", "test"))

	o.ServeHTTP(recorder, req)
	expect(t, recorder.Code, http.StatusOK)
	refute(t, len(buff.String()), 0)
	expect(t, buff.String(), "test")
	expect(t, strings.Split(recorder.Header().Get("Set-Cookie"), "|")[0], "expire=test")
}

func TestCookie8(t *testing.T) {
	buff := bytes.NewBufferString("")
	recorder := httptest.NewRecorder()
	recorder.Body = buff

	o := Classic()
	o.Get("/", func(ctx *Context) string {
		ctx.SecureCookies("ttttt").Del("ttttt")
		return "test"
	})

	req, err := http.NewRequest("GET", "http://localhost:8000/", nil)
	if err != nil {
		t.Error(err)
	}

	req.AddCookie(NewSecureCookie("ttttt", "ttttt", "test"))

	o.ServeHTTP(recorder, req)
	expect(t, recorder.Code, http.StatusOK)
	refute(t, len(buff.String()), 0)
	expect(t, buff.String(), "test")
	expect(t, strings.Split(recorder.Header().Get("Set-Cookie"), "|")[0], "ttttt=test")
}

type Cookie11Action struct {
	Ctx
}

func (a *Cookie11Action) Get() string {
	v, _ := a.Cookies().Int("test")
	return fmt.Sprintf("%d", v)
}

func TestCookie11(t *testing.T) {
	buff := bytes.NewBufferString("")
	recorder := httptest.NewRecorder()
	recorder.Body = buff

	o := Classic()
	o.Get("/", new(Cookie11Action))

	req, err := http.NewRequest("GET", "http://localhost:8000", nil)
	if err != nil {
		t.Error(err)
	}
	req.AddCookie(NewCookie("test", "1"))

	o.ServeHTTP(recorder, req)
	expect(t, recorder.Code, http.StatusOK)
	refute(t, len(buff.String()), 0)
	expect(t, buff.String(), "1")
}

type Cookie12Action struct {
	Ctx
}

func (a *Cookie12Action) Get() string {
	v, _ := a.Cookies().Int32("test")
	return fmt.Sprintf("%d", v)
}

func TestCookie12(t *testing.T) {
	buff := bytes.NewBufferString("")
	recorder := httptest.NewRecorder()
	recorder.Body = buff

	o := Classic()
	o.Get("/", new(Cookie12Action))

	req, err := http.NewRequest("GET", "http://localhost:8000/", nil)
	if err != nil {
		t.Error(err)
	}
	req.AddCookie(NewCookie("test", "1"))

	o.ServeHTTP(recorder, req)
	expect(t, recorder.Code, http.StatusOK)
	refute(t, len(buff.String()), 0)
	expect(t, buff.String(), "1")
}

type Cookie13Action struct {
	Ctx
}

func (a *Cookie13Action) Get() string {
	v, _ := a.Cookies().Int64("test")
	return fmt.Sprintf("%d", v)
}

func TestCookie13(t *testing.T) {
	buff := bytes.NewBufferString("")
	recorder := httptest.NewRecorder()
	recorder.Body = buff

	o := Classic()
	o.Get("/", new(Cookie13Action))

	req, err := http.NewRequest("GET", "http://localhost:8000/", nil)
	if err != nil {
		t.Error(err)
	}
	req.AddCookie(NewCookie("test", "1"))

	o.ServeHTTP(recorder, req)
	expect(t, recorder.Code, http.StatusOK)
	refute(t, len(buff.String()), 0)
	expect(t, buff.String(), "1")
}

type Cookie14Action struct {
	Ctx
}

func (a *Cookie14Action) Get() string {
	v, _ := a.Cookies().Uint("test")
	return fmt.Sprintf("%d", v)
}

func TestCookie14(t *testing.T) {
	buff := bytes.NewBufferString("")
	recorder := httptest.NewRecorder()
	recorder.Body = buff

	o := Classic()
	o.Get("/", new(Cookie14Action))

	req, err := http.NewRequest("GET", "http://localhost:8000/", nil)
	if err != nil {
		t.Error(err)
	}
	req.AddCookie(NewCookie("test", "1"))

	o.ServeHTTP(recorder, req)
	expect(t, recorder.Code, http.StatusOK)
	refute(t, len(buff.String()), 0)
	expect(t, buff.String(), "1")
}

type Cookie15Action struct {
	Ctx
}

func (a *Cookie15Action) Get() string {
	v, _ := a.Cookies().Uint32("test")
	return fmt.Sprintf("%d", v)
}

func TestCookie15(t *testing.T) {
	buff := bytes.NewBufferString("")
	recorder := httptest.NewRecorder()
	recorder.Body = buff

	o := Classic()
	o.Get("/", new(Cookie15Action))

	req, err := http.NewRequest("GET", "http://localhost:8000/", nil)
	if err != nil {
		t.Error(err)
	}
	req.AddCookie(NewCookie("test", "1"))

	o.ServeHTTP(recorder, req)
	expect(t, recorder.Code, http.StatusOK)
	refute(t, len(buff.String()), 0)
	expect(t, buff.String(), "1")
}

type Cookie16Action struct {
	Ctx
}

func (a *Cookie16Action) Get() string {
	v, _ := a.Cookies().Uint64("test")
	return fmt.Sprintf("%d", v)
}

func TestCookie16(t *testing.T) {
	buff := bytes.NewBufferString("")
	recorder := httptest.NewRecorder()
	recorder.Body = buff

	o := Classic()
	o.Get("/", new(Cookie16Action))

	req, err := http.NewRequest("GET", "http://localhost:8000/", nil)
	if err != nil {
		t.Error(err)
	}
	req.AddCookie(NewCookie("test", "1"))

	o.ServeHTTP(recorder, req)
	expect(t, recorder.Code, http.StatusOK)
	refute(t, len(buff.String()), 0)
	expect(t, buff.String(), "1")
}

type Cookie17Action struct {
	Ctx
}

func (a *Cookie17Action) Get() string {
	v, _ := a.Cookies().Float32("test")
	return fmt.Sprintf("%.2f", v)
}

func TestCookie17(t *testing.T) {
	buff := bytes.NewBufferString("")
	recorder := httptest.NewRecorder()
	recorder.Body = buff

	o := Classic()
	o.Get("/", new(Cookie17Action))

	req, err := http.NewRequest("GET", "http://localhost:8000/", nil)
	if err != nil {
		t.Error(err)
	}
	req.AddCookie(NewCookie("test", "1"))

	o.ServeHTTP(recorder, req)
	expect(t, recorder.Code, http.StatusOK)
	refute(t, len(buff.String()), 0)
	expect(t, buff.String(), "1.00")
}

type Cookie18Action struct {
	Ctx
}

func (a *Cookie18Action) Get() string {
	v, _ := a.Cookies().Float64("test")
	return fmt.Sprintf("%.2f", v)
}

func TestCookie18(t *testing.T) {
	buff := bytes.NewBufferString("")
	recorder := httptest.NewRecorder()
	recorder.Body = buff

	o := Classic()
	o.Get("/", new(Cookie18Action))

	req, err := http.NewRequest("GET", "http://localhost:8000/", nil)
	if err != nil {
		t.Error(err)
	}
	req.AddCookie(NewCookie("test", "1"))

	o.ServeHTTP(recorder, req)
	expect(t, recorder.Code, http.StatusOK)
	refute(t, len(buff.String()), 0)
	expect(t, buff.String(), "1.00")
}

type Cookie19Action struct {
	Ctx
}

func (a *Cookie19Action) Get() string {
	v, _ := a.Cookies().Bool("test")
	return fmt.Sprintf("%v", v)
}

func TestCookie19(t *testing.T) {
	buff := bytes.NewBufferString("")
	recorder := httptest.NewRecorder()
	recorder.Body = buff

	o := Classic()
	o.Get("/", new(Cookie19Action))

	req, err := http.NewRequest("GET", "http://localhost:8000/", nil)
	if err != nil {
		t.Error(err)
	}
	req.AddCookie(NewCookie("test", "1"))

	o.ServeHTTP(recorder, req)
	expect(t, recorder.Code, http.StatusOK)
	refute(t, len(buff.String()), 0)
	expect(t, buff.String(), "true")
}

type Cookie20Action struct {
	Ctx
}

func (a *Cookie20Action) Get() string {
	v, _ := a.Cookies().String("test")
	return fmt.Sprintf("%v", v)
}

func TestCookie20(t *testing.T) {
	buff := bytes.NewBufferString("")
	recorder := httptest.NewRecorder()
	recorder.Body = buff

	o := Classic()
	o.Get("/", new(Cookie20Action))

	req, err := http.NewRequest("GET", "http://localhost:8000/", nil)
	if err != nil {
		t.Error(err)
	}
	req.AddCookie(NewCookie("test", "1"))

	o.ServeHTTP(recorder, req)
	expect(t, recorder.Code, http.StatusOK)
	refute(t, len(buff.String()), 0)
	expect(t, buff.String(), "1")
}

type Cookie21Action struct {
	Ctx
}

func (a *Cookie21Action) Get() string {
	v := a.Cookies().MustInt("test")
	return fmt.Sprintf("%d", v)
}

func TestCookie21(t *testing.T) {
	buff := bytes.NewBufferString("")
	recorder := httptest.NewRecorder()
	recorder.Body = buff

	o := Classic()
	o.Get("/", new(Cookie21Action))

	req, err := http.NewRequest("GET", "http://localhost:8000/", nil)
	if err != nil {
		t.Error(err)
	}
	req.AddCookie(NewCookie("test", "1"))

	o.ServeHTTP(recorder, req)
	expect(t, recorder.Code, http.StatusOK)
	refute(t, len(buff.String()), 0)
	expect(t, buff.String(), "1")
}

type Cookie22Action struct {
	Ctx
}

func (a *Cookie22Action) Get() string {
	v := a.Cookies().MustInt32("test")
	return fmt.Sprintf("%d", v)
}

func TestCookie22(t *testing.T) {
	buff := bytes.NewBufferString("")
	recorder := httptest.NewRecorder()
	recorder.Body = buff

	o := Classic()
	o.Get("/", new(Cookie22Action))

	req, err := http.NewRequest("GET", "http://localhost:8000/", nil)
	if err != nil {
		t.Error(err)
	}
	req.AddCookie(NewCookie("test", "1"))

	o.ServeHTTP(recorder, req)
	expect(t, recorder.Code, http.StatusOK)
	refute(t, len(buff.String()), 0)
	expect(t, buff.String(), "1")
}

type Cookie23Action struct {
	Ctx
}

func (a *Cookie23Action) Get() string {
	v := a.Cookies().MustInt64("test")
	return fmt.Sprintf("%d", v)
}

func TestCookie23(t *testing.T) {
	buff := bytes.NewBufferString("")
	recorder := httptest.NewRecorder()
	recorder.Body = buff

	o := Classic()
	o.Get("/", new(Cookie23Action))

	req, err := http.NewRequest("GET", "http://localhost:8000/", nil)
	if err != nil {
		t.Error(err)
	}
	req.AddCookie(NewCookie("test", "1"))

	o.ServeHTTP(recorder, req)
	expect(t, recorder.Code, http.StatusOK)
	refute(t, len(buff.String()), 0)
	expect(t, buff.String(), "1")
}

type Cookie24Action struct {
	Ctx
}

func (a *Cookie24Action) Get() string {
	v := a.Cookies().MustUint("test")
	return fmt.Sprintf("%d", v)
}

func TestCookie24(t *testing.T) {
	buff := bytes.NewBufferString("")
	recorder := httptest.NewRecorder()
	recorder.Body = buff

	o := Classic()
	o.Get("/", new(Cookie24Action))

	req, err := http.NewRequest("GET", "http://localhost:8000/", nil)
	if err != nil {
		t.Error(err)
	}
	req.AddCookie(NewCookie("test", "1"))

	o.ServeHTTP(recorder, req)
	expect(t, recorder.Code, http.StatusOK)
	refute(t, len(buff.String()), 0)
	expect(t, buff.String(), "1")
}

type Cookie25Action struct {
	Ctx
}

func (a *Cookie25Action) Get() string {
	v := a.Cookies().MustUint32("test")
	return fmt.Sprintf("%d", v)
}

func TestCookie25(t *testing.T) {
	buff := bytes.NewBufferString("")
	recorder := httptest.NewRecorder()
	recorder.Body = buff

	o := Classic()
	o.Get("/", new(Cookie25Action))

	req, err := http.NewRequest("GET", "http://localhost:8000/", nil)
	if err != nil {
		t.Error(err)
	}
	req.AddCookie(NewCookie("test", "1"))

	o.ServeHTTP(recorder, req)
	expect(t, recorder.Code, http.StatusOK)
	refute(t, len(buff.String()), 0)
	expect(t, buff.String(), "1")
}

type Cookie26Action struct {
	Ctx
}

func (a *Cookie26Action) Get() string {
	v := a.Cookies().MustUint64("test")
	return fmt.Sprintf("%d", v)
}

func TestCookie26(t *testing.T) {
	buff := bytes.NewBufferString("")
	recorder := httptest.NewRecorder()
	recorder.Body = buff

	o := Classic()
	o.Get("/", new(Cookie26Action))

	req, err := http.NewRequest("GET", "http://localhost:8000/", nil)
	if err != nil {
		t.Error(err)
	}
	req.AddCookie(NewCookie("test", "1"))

	o.ServeHTTP(recorder, req)
	expect(t, recorder.Code, http.StatusOK)
	refute(t, len(buff.String()), 0)
	expect(t, buff.String(), "1")
}

type Cookie27Action struct {
	Ctx
}

func (a *Cookie27Action) Get() string {
	v := a.Cookies().MustFloat32("test")
	return fmt.Sprintf("%.2f", v)
}

func TestCookie27(t *testing.T) {
	buff := bytes.NewBufferString("")
	recorder := httptest.NewRecorder()
	recorder.Body = buff

	o := Classic()
	o.Get("/", new(Cookie27Action))

	req, err := http.NewRequest("GET", "http://localhost:8000/", nil)
	if err != nil {
		t.Error(err)
	}
	req.AddCookie(NewCookie("test", "1"))

	o.ServeHTTP(recorder, req)
	expect(t, recorder.Code, http.StatusOK)
	refute(t, len(buff.String()), 0)
	expect(t, buff.String(), "1.00")
}

type Cookie28Action struct {
	Ctx
}

func (a *Cookie28Action) Get() string {
	v := a.Cookies().MustFloat64("test")
	return fmt.Sprintf("%.2f", v)
}

func TestCookie28(t *testing.T) {
	buff := bytes.NewBufferString("")
	recorder := httptest.NewRecorder()
	recorder.Body = buff

	o := Classic()
	o.Get("/", new(Cookie28Action))

	req, err := http.NewRequest("GET", "http://localhost:8000/", nil)
	if err != nil {
		t.Error(err)
	}
	req.AddCookie(NewCookie("test", "1"))

	o.ServeHTTP(recorder, req)
	expect(t, recorder.Code, http.StatusOK)
	refute(t, len(buff.String()), 0)
	expect(t, buff.String(), "1.00")
}

type Cookie29Action struct {
	Ctx
}

func (a *Cookie29Action) Get() string {
	v := a.Cookies().MustBool("test")
	return fmt.Sprintf("%v", v)
}

func TestCookie29(t *testing.T) {
	buff := bytes.NewBufferString("")
	recorder := httptest.NewRecorder()
	recorder.Body = buff

	o := Classic()
	o.Get("/", new(Cookie29Action))

	req, err := http.NewRequest("GET", "http://localhost:8000/", nil)
	if err != nil {
		t.Error(err)
	}
	req.AddCookie(NewCookie("test", "1"))

	o.ServeHTTP(recorder, req)
	expect(t, recorder.Code, http.StatusOK)
	refute(t, len(buff.String()), 0)
	expect(t, buff.String(), "true")
}

type Cookie30Action struct {
	Ctx
}

func (a *Cookie30Action) Get() string {
	v := a.Cookies().MustString("test")
	return fmt.Sprintf("%s", v)
}

func TestCookie30(t *testing.T) {
	buff := bytes.NewBufferString("")
	recorder := httptest.NewRecorder()
	recorder.Body = buff

	o := Classic()
	o.Get("/", new(Cookie30Action))

	req, err := http.NewRequest("GET", "http://localhost:8000/", nil)
	if err != nil {
		t.Error(err)
	}
	req.AddCookie(NewCookie("test", "1"))

	o.ServeHTTP(recorder, req)
	expect(t, recorder.Code, http.StatusOK)
	refute(t, len(buff.String()), 0)
	expect(t, buff.String(), "1")
}

const (
	secrectCookie = "ssss"
)

type Cookie31Action struct {
	Ctx
}

func (a *Cookie31Action) Get() string {
	v, _ := a.SecureCookies(secrectCookie).Int("test")
	return fmt.Sprintf("%d", v)
}

func TestCookie31(t *testing.T) {
	buff := bytes.NewBufferString("")
	recorder := httptest.NewRecorder()
	recorder.Body = buff

	o := Classic()
	o.Get("/", new(Cookie31Action))

	req, err := http.NewRequest("GET", "http://localhost:8000", nil)
	if err != nil {
		t.Error(err)
	}
	req.AddCookie(NewSecureCookie(secrectCookie, "test", "1"))

	o.ServeHTTP(recorder, req)
	expect(t, recorder.Code, http.StatusOK)
	refute(t, len(buff.String()), 0)
	expect(t, buff.String(), "1")
}

type Cookie32Action struct {
	Ctx
}

func (a *Cookie32Action) Get() string {
	v, _ := a.SecureCookies(secrectCookie).Int32("test")
	return fmt.Sprintf("%d", v)
}

func TestCookie32(t *testing.T) {
	buff := bytes.NewBufferString("")
	recorder := httptest.NewRecorder()
	recorder.Body = buff

	o := Classic()
	o.Get("/", new(Cookie32Action))

	req, err := http.NewRequest("GET", "http://localhost:8000/", nil)
	if err != nil {
		t.Error(err)
	}
	req.AddCookie(NewSecureCookie(secrectCookie, "test", "1"))

	o.ServeHTTP(recorder, req)
	expect(t, recorder.Code, http.StatusOK)
	refute(t, len(buff.String()), 0)
	expect(t, buff.String(), "1")
}

type Cookie33Action struct {
	Ctx
}

func (a *Cookie33Action) Get() string {
	v, _ := a.SecureCookies(secrectCookie).Int64("test")
	return fmt.Sprintf("%d", v)
}

func TestCookie33(t *testing.T) {
	buff := bytes.NewBufferString("")
	recorder := httptest.NewRecorder()
	recorder.Body = buff

	o := Classic()
	o.Get("/", new(Cookie33Action))

	req, err := http.NewRequest("GET", "http://localhost:8000/", nil)
	if err != nil {
		t.Error(err)
	}
	req.AddCookie(NewSecureCookie(secrectCookie, "test", "1"))

	o.ServeHTTP(recorder, req)
	expect(t, recorder.Code, http.StatusOK)
	refute(t, len(buff.String()), 0)
	expect(t, buff.String(), "1")
}

type Cookie34Action struct {
	Ctx
}

func (a *Cookie34Action) Get() string {
	v, _ := a.SecureCookies(secrectCookie).Uint("test")
	return fmt.Sprintf("%d", v)
}

func TestCookie34(t *testing.T) {
	buff := bytes.NewBufferString("")
	recorder := httptest.NewRecorder()
	recorder.Body = buff

	o := Classic()
	o.Get("/", new(Cookie34Action))

	req, err := http.NewRequest("GET", "http://localhost:8000/", nil)
	if err != nil {
		t.Error(err)
	}
	req.AddCookie(NewSecureCookie(secrectCookie, "test", "1"))

	o.ServeHTTP(recorder, req)
	expect(t, recorder.Code, http.StatusOK)
	refute(t, len(buff.String()), 0)
	expect(t, buff.String(), "1")
}

type Cookie35Action struct {
	Ctx
}

func (a *Cookie35Action) Get() string {
	v, _ := a.SecureCookies(secrectCookie).Uint32("test")
	return fmt.Sprintf("%d", v)
}

func TestCookie35(t *testing.T) {
	buff := bytes.NewBufferString("")
	recorder := httptest.NewRecorder()
	recorder.Body = buff

	o := Classic()
	o.Get("/", new(Cookie35Action))

	req, err := http.NewRequest("GET", "http://localhost:8000/", nil)
	if err != nil {
		t.Error(err)
	}
	req.AddCookie(NewSecureCookie(secrectCookie, "test", "1"))

	o.ServeHTTP(recorder, req)
	expect(t, recorder.Code, http.StatusOK)
	refute(t, len(buff.String()), 0)
	expect(t, buff.String(), "1")
}

type Cookie36Action struct {
	Ctx
}

func (a *Cookie36Action) Get() string {
	v, _ := a.SecureCookies(secrectCookie).Uint64("test")
	return fmt.Sprintf("%d", v)
}

func TestCookie36(t *testing.T) {
	buff := bytes.NewBufferString("")
	recorder := httptest.NewRecorder()
	recorder.Body = buff

	o := Classic()
	o.Get("/", new(Cookie36Action))

	req, err := http.NewRequest("GET", "http://localhost:8000/", nil)
	if err != nil {
		t.Error(err)
	}
	req.AddCookie(NewSecureCookie(secrectCookie, "test", "1"))

	o.ServeHTTP(recorder, req)
	expect(t, recorder.Code, http.StatusOK)
	refute(t, len(buff.String()), 0)
	expect(t, buff.String(), "1")
}

type Cookie37Action struct {
	Ctx
}

func (a *Cookie37Action) Get() string {
	v, _ := a.SecureCookies(secrectCookie).Float32("test")
	return fmt.Sprintf("%.2f", v)
}

func TestCookie37(t *testing.T) {
	buff := bytes.NewBufferString("")
	recorder := httptest.NewRecorder()
	recorder.Body = buff

	o := Classic()
	o.Get("/", new(Cookie37Action))

	req, err := http.NewRequest("GET", "http://localhost:8000/", nil)
	if err != nil {
		t.Error(err)
	}
	req.AddCookie(NewSecureCookie(secrectCookie, "test", "1"))

	o.ServeHTTP(recorder, req)
	expect(t, recorder.Code, http.StatusOK)
	refute(t, len(buff.String()), 0)
	expect(t, buff.String(), "1.00")
}

type Cookie38Action struct {
	Ctx
}

func (a *Cookie38Action) Get() string {
	v, _ := a.SecureCookies(secrectCookie).Float64("test")
	return fmt.Sprintf("%.2f", v)
}

func TestCookie38(t *testing.T) {
	buff := bytes.NewBufferString("")
	recorder := httptest.NewRecorder()
	recorder.Body = buff

	o := Classic()
	o.Get("/", new(Cookie38Action))

	req, err := http.NewRequest("GET", "http://localhost:8000/", nil)
	if err != nil {
		t.Error(err)
	}
	req.AddCookie(NewSecureCookie(secrectCookie, "test", "1"))

	o.ServeHTTP(recorder, req)
	expect(t, recorder.Code, http.StatusOK)
	refute(t, len(buff.String()), 0)
	expect(t, buff.String(), "1.00")
}

type Cookie39Action struct {
	Ctx
}

func (a *Cookie39Action) Get() string {
	v, _ := a.SecureCookies(secrectCookie).Bool("test")
	return fmt.Sprintf("%v", v)
}

func TestCookie39(t *testing.T) {
	buff := bytes.NewBufferString("")
	recorder := httptest.NewRecorder()
	recorder.Body = buff

	o := Classic()
	o.Get("/", new(Cookie39Action))

	req, err := http.NewRequest("GET", "http://localhost:8000/", nil)
	if err != nil {
		t.Error(err)
	}
	req.AddCookie(NewSecureCookie(secrectCookie, "test", "1"))

	o.ServeHTTP(recorder, req)
	expect(t, recorder.Code, http.StatusOK)
	refute(t, len(buff.String()), 0)
	expect(t, buff.String(), "true")
}

type Cookie40Action struct {
	Ctx
}

func (a *Cookie40Action) Get() string {
	v, _ := a.SecureCookies(secrectCookie).String("test")
	return fmt.Sprintf("%v", v)
}

func TestCookie40(t *testing.T) {
	buff := bytes.NewBufferString("")
	recorder := httptest.NewRecorder()
	recorder.Body = buff

	o := Classic()
	o.Get("/", new(Cookie40Action))

	req, err := http.NewRequest("GET", "http://localhost:8000/", nil)
	if err != nil {
		t.Error(err)
	}
	req.AddCookie(NewSecureCookie(secrectCookie, "test", "1"))

	o.ServeHTTP(recorder, req)
	expect(t, recorder.Code, http.StatusOK)
	refute(t, len(buff.String()), 0)
	expect(t, buff.String(), "1")
}

type Cookie41Action struct {
	Ctx
}

func (a *Cookie41Action) Get() string {
	v := a.SecureCookies(secrectCookie).MustInt("test")
	return fmt.Sprintf("%d", v)
}

func TestCookie41(t *testing.T) {
	buff := bytes.NewBufferString("")
	recorder := httptest.NewRecorder()
	recorder.Body = buff

	o := Classic()
	o.Get("/", new(Cookie41Action))

	req, err := http.NewRequest("GET", "http://localhost:8000/", nil)
	if err != nil {
		t.Error(err)
	}
	req.AddCookie(NewSecureCookie(secrectCookie, "test", "1"))

	o.ServeHTTP(recorder, req)
	expect(t, recorder.Code, http.StatusOK)
	refute(t, len(buff.String()), 0)
	expect(t, buff.String(), "1")
}

type Cookie42Action struct {
	Ctx
}

func (a *Cookie42Action) Get() string {
	v := a.SecureCookies(secrectCookie).MustInt32("test")
	return fmt.Sprintf("%d", v)
}

func TestCookie42(t *testing.T) {
	buff := bytes.NewBufferString("")
	recorder := httptest.NewRecorder()
	recorder.Body = buff

	o := Classic()
	o.Get("/", new(Cookie42Action))

	req, err := http.NewRequest("GET", "http://localhost:8000/", nil)
	if err != nil {
		t.Error(err)
	}
	req.AddCookie(NewSecureCookie(secrectCookie, "test", "1"))

	o.ServeHTTP(recorder, req)
	expect(t, recorder.Code, http.StatusOK)
	refute(t, len(buff.String()), 0)
	expect(t, buff.String(), "1")
}

type Cookie43Action struct {
	Ctx
}

func (a *Cookie43Action) Get() string {
	v := a.SecureCookies(secrectCookie).MustInt64("test")
	return fmt.Sprintf("%d", v)
}

func TestCookie43(t *testing.T) {
	buff := bytes.NewBufferString("")
	recorder := httptest.NewRecorder()
	recorder.Body = buff

	o := Classic()
	o.Get("/", new(Cookie43Action))

	req, err := http.NewRequest("GET", "http://localhost:8000/", nil)
	if err != nil {
		t.Error(err)
	}
	req.AddCookie(NewSecureCookie(secrectCookie, "test", "1"))

	o.ServeHTTP(recorder, req)
	expect(t, recorder.Code, http.StatusOK)
	refute(t, len(buff.String()), 0)
	expect(t, buff.String(), "1")
}

type Cookie44Action struct {
	Ctx
}

func (a *Cookie44Action) Get() string {
	v := a.SecureCookies(secrectCookie).MustUint("test")
	return fmt.Sprintf("%d", v)
}

func TestCookie44(t *testing.T) {
	buff := bytes.NewBufferString("")
	recorder := httptest.NewRecorder()
	recorder.Body = buff

	o := Classic()
	o.Get("/", new(Cookie44Action))

	req, err := http.NewRequest("GET", "http://localhost:8000/", nil)
	if err != nil {
		t.Error(err)
	}
	req.AddCookie(NewSecureCookie(secrectCookie, "test", "1"))

	o.ServeHTTP(recorder, req)
	expect(t, recorder.Code, http.StatusOK)
	refute(t, len(buff.String()), 0)
	expect(t, buff.String(), "1")
}

type Cookie45Action struct {
	Ctx
}

func (a *Cookie45Action) Get() string {
	v := a.SecureCookies(secrectCookie).MustUint32("test")
	return fmt.Sprintf("%d", v)
}

func TestCookie45(t *testing.T) {
	buff := bytes.NewBufferString("")
	recorder := httptest.NewRecorder()
	recorder.Body = buff

	o := Classic()
	o.Get("/", new(Cookie45Action))

	req, err := http.NewRequest("GET", "http://localhost:8000/", nil)
	if err != nil {
		t.Error(err)
	}
	req.AddCookie(NewSecureCookie(secrectCookie, "test", "1"))

	o.ServeHTTP(recorder, req)
	expect(t, recorder.Code, http.StatusOK)
	refute(t, len(buff.String()), 0)
	expect(t, buff.String(), "1")
}

type Cookie46Action struct {
	Ctx
}

func (a *Cookie46Action) Get() string {
	v := a.SecureCookies(secrectCookie).MustUint64("test")
	return fmt.Sprintf("%d", v)
}

func TestCookie46(t *testing.T) {
	buff := bytes.NewBufferString("")
	recorder := httptest.NewRecorder()
	recorder.Body = buff

	o := Classic()
	o.Get("/", new(Cookie46Action))

	req, err := http.NewRequest("GET", "http://localhost:8000/", nil)
	if err != nil {
		t.Error(err)
	}
	req.AddCookie(NewSecureCookie(secrectCookie, "test", "1"))

	o.ServeHTTP(recorder, req)
	expect(t, recorder.Code, http.StatusOK)
	refute(t, len(buff.String()), 0)
	expect(t, buff.String(), "1")
}

type Cookie47Action struct {
	Ctx
}

func (a *Cookie47Action) Get() string {
	v := a.SecureCookies(secrectCookie).MustFloat32("test")
	return fmt.Sprintf("%.2f", v)
}

func TestCookie47(t *testing.T) {
	buff := bytes.NewBufferString("")
	recorder := httptest.NewRecorder()
	recorder.Body = buff

	o := Classic()
	o.Get("/", new(Cookie47Action))

	req, err := http.NewRequest("GET", "http://localhost:8000/", nil)
	if err != nil {
		t.Error(err)
	}
	req.AddCookie(NewSecureCookie(secrectCookie, "test", "1"))

	o.ServeHTTP(recorder, req)
	expect(t, recorder.Code, http.StatusOK)
	refute(t, len(buff.String()), 0)
	expect(t, buff.String(), "1.00")
}

type Cookie48Action struct {
	Ctx
}

func (a *Cookie48Action) Get() string {
	v := a.SecureCookies(secrectCookie).MustFloat64("test")
	return fmt.Sprintf("%.2f", v)
}

func TestCookie48(t *testing.T) {
	buff := bytes.NewBufferString("")
	recorder := httptest.NewRecorder()
	recorder.Body = buff

	o := Classic()
	o.Get("/", new(Cookie48Action))

	req, err := http.NewRequest("GET", "http://localhost:8000/", nil)
	if err != nil {
		t.Error(err)
	}
	req.AddCookie(NewSecureCookie(secrectCookie, "test", "1"))

	o.ServeHTTP(recorder, req)
	expect(t, recorder.Code, http.StatusOK)
	refute(t, len(buff.String()), 0)
	expect(t, buff.String(), "1.00")
}

type Cookie49Action struct {
	Ctx
}

func (a *Cookie49Action) Get() string {
	v := a.SecureCookies(secrectCookie).MustBool("test")
	return fmt.Sprintf("%v", v)
}

func TestCookie49(t *testing.T) {
	buff := bytes.NewBufferString("")
	recorder := httptest.NewRecorder()
	recorder.Body = buff

	o := Classic()
	o.Get("/", new(Cookie49Action))

	req, err := http.NewRequest("GET", "http://localhost:8000/", nil)
	if err != nil {
		t.Error(err)
	}
	req.AddCookie(NewSecureCookie(secrectCookie, "test", "1"))

	o.ServeHTTP(recorder, req)
	expect(t, recorder.Code, http.StatusOK)
	refute(t, len(buff.String()), 0)
	expect(t, buff.String(), "true")
}

type Cookie50Action struct {
	Ctx
}

func (a *Cookie50Action) Get() string {
	v := a.SecureCookies(secrectCookie).MustString("test")
	return fmt.Sprintf("%s", v)
}

func TestCookie50(t *testing.T) {
	buff := bytes.NewBufferString("")
	recorder := httptest.NewRecorder()
	recorder.Body = buff

	o := Classic()
	o.Get("/", new(Cookie50Action))

	req, err := http.NewRequest("GET", "http://localhost:8000/", nil)
	if err != nil {
		t.Error(err)
	}
	req.AddCookie(NewSecureCookie(secrectCookie, "test", "1"))

	o.ServeHTTP(recorder, req)
	expect(t, recorder.Code, http.StatusOK)
	refute(t, len(buff.String()), 0)
	expect(t, buff.String(), "1")
}

type Cookie51Action struct {
	Ctx
}

func (a *Cookie51Action) Get() string {
	v := a.CookieInt("test")
	return fmt.Sprintf("%d", v)
}

func TestCookie51(t *testing.T) {
	buff := bytes.NewBufferString("")
	recorder := httptest.NewRecorder()
	recorder.Body = buff

	o := Classic()
	o.Get("/", new(Cookie51Action))

	req, err := http.NewRequest("GET", "http://localhost:8000/", nil)
	if err != nil {
		t.Error(err)
	}
	req.AddCookie(NewCookie("test", "1"))

	o.ServeHTTP(recorder, req)
	expect(t, recorder.Code, http.StatusOK)
	refute(t, len(buff.String()), 0)
	expect(t, buff.String(), "1")
}

type Cookie52Action struct {
	Ctx
}

func (a *Cookie52Action) Get() string {
	v := a.CookieInt32("test")
	return fmt.Sprintf("%d", v)
}

func TestCookie52(t *testing.T) {
	buff := bytes.NewBufferString("")
	recorder := httptest.NewRecorder()
	recorder.Body = buff

	o := Classic()
	o.Get("/", new(Cookie52Action))

	req, err := http.NewRequest("GET", "http://localhost:8000/", nil)
	if err != nil {
		t.Error(err)
	}
	req.AddCookie(NewCookie("test", "1"))

	o.ServeHTTP(recorder, req)
	expect(t, recorder.Code, http.StatusOK)
	refute(t, len(buff.String()), 0)
	expect(t, buff.String(), "1")
}

type Cookie53Action struct {
	Ctx
}

func (a *Cookie53Action) Get() string {
	v := a.CookieInt64("test")
	return fmt.Sprintf("%d", v)
}

func TestCookie53(t *testing.T) {
	buff := bytes.NewBufferString("")
	recorder := httptest.NewRecorder()
	recorder.Body = buff

	o := Classic()
	o.Get("/", new(Cookie53Action))

	req, err := http.NewRequest("GET", "http://localhost:8000/", nil)
	if err != nil {
		t.Error(err)
	}
	req.AddCookie(NewCookie("test", "1"))

	o.ServeHTTP(recorder, req)
	expect(t, recorder.Code, http.StatusOK)
	refute(t, len(buff.String()), 0)
	expect(t, buff.String(), "1")
}

type Cookie54Action struct {
	Ctx
}

func (a *Cookie54Action) Get() string {
	v := a.CookieUint("test")
	return fmt.Sprintf("%d", v)
}

func TestCookie54(t *testing.T) {
	buff := bytes.NewBufferString("")
	recorder := httptest.NewRecorder()
	recorder.Body = buff

	o := Classic()
	o.Get("/", new(Cookie54Action))

	req, err := http.NewRequest("GET", "http://localhost:8000/", nil)
	if err != nil {
		t.Error(err)
	}
	req.AddCookie(NewCookie("test", "1"))

	o.ServeHTTP(recorder, req)
	expect(t, recorder.Code, http.StatusOK)
	refute(t, len(buff.String()), 0)
	expect(t, buff.String(), "1")
}

type Cookie55Action struct {
	Ctx
}

func (a *Cookie55Action) Get() string {
	v := a.CookieUint32("test")
	return fmt.Sprintf("%d", v)
}

func TestCookie55(t *testing.T) {
	buff := bytes.NewBufferString("")
	recorder := httptest.NewRecorder()
	recorder.Body = buff

	o := Classic()
	o.Get("/", new(Cookie55Action))

	req, err := http.NewRequest("GET", "http://localhost:8000/", nil)
	if err != nil {
		t.Error(err)
	}
	req.AddCookie(NewCookie("test", "1"))

	o.ServeHTTP(recorder, req)
	expect(t, recorder.Code, http.StatusOK)
	refute(t, len(buff.String()), 0)
	expect(t, buff.String(), "1")
}

type Cookie56Action struct {
	Ctx
}

func (a *Cookie56Action) Get() string {
	v := a.CookieUint64("test")
	return fmt.Sprintf("%d", v)
}

func TestCookie56(t *testing.T) {
	buff := bytes.NewBufferString("")
	recorder := httptest.NewRecorder()
	recorder.Body = buff

	o := Classic()
	o.Get("/", new(Cookie56Action))

	req, err := http.NewRequest("GET", "http://localhost:8000/", nil)
	if err != nil {
		t.Error(err)
	}
	req.AddCookie(NewCookie("test", "1"))

	o.ServeHTTP(recorder, req)
	expect(t, recorder.Code, http.StatusOK)
	refute(t, len(buff.String()), 0)
	expect(t, buff.String(), "1")
}

type Cookie57Action struct {
	Ctx
}

func (a *Cookie57Action) Get() string {
	v := a.CookieFloat32("test")
	return fmt.Sprintf("%.2f", v)
}

func TestCookie57(t *testing.T) {
	buff := bytes.NewBufferString("")
	recorder := httptest.NewRecorder()
	recorder.Body = buff

	o := Classic()
	o.Get("/", new(Cookie57Action))

	req, err := http.NewRequest("GET", "http://localhost:8000/", nil)
	if err != nil {
		t.Error(err)
	}
	req.AddCookie(NewCookie("test", "1"))

	o.ServeHTTP(recorder, req)
	expect(t, recorder.Code, http.StatusOK)
	refute(t, len(buff.String()), 0)
	expect(t, buff.String(), "1.00")
}

type Cookie58Action struct {
	Ctx
}

func (a *Cookie58Action) Get() string {
	v := a.CookieFloat64("test")
	return fmt.Sprintf("%.2f", v)
}

func TestCookie58(t *testing.T) {
	buff := bytes.NewBufferString("")
	recorder := httptest.NewRecorder()
	recorder.Body = buff

	o := Classic()
	o.Get("/", new(Cookie58Action))

	req, err := http.NewRequest("GET", "http://localhost:8000/", nil)
	if err != nil {
		t.Error(err)
	}
	req.AddCookie(NewCookie("test", "1"))

	o.ServeHTTP(recorder, req)
	expect(t, recorder.Code, http.StatusOK)
	refute(t, len(buff.String()), 0)
	expect(t, buff.String(), "1.00")
}

type Cookie59Action struct {
	Ctx
}

func (a *Cookie59Action) Get() string {
	v := a.CookieBool("test")
	return fmt.Sprintf("%v", v)
}

func TestCookie59(t *testing.T) {
	buff := bytes.NewBufferString("")
	recorder := httptest.NewRecorder()
	recorder.Body = buff

	o := Classic()
	o.Get("/", new(Cookie59Action))

	req, err := http.NewRequest("GET", "http://localhost:8000/", nil)
	if err != nil {
		t.Error(err)
	}
	req.AddCookie(NewCookie("test", "1"))

	o.ServeHTTP(recorder, req)
	expect(t, recorder.Code, http.StatusOK)
	refute(t, len(buff.String()), 0)
	expect(t, buff.String(), "true")
}

type Cookie60Action struct {
	Ctx
}

func (a *Cookie60Action) Get() string {
	v := a.Cookie("test")
	return fmt.Sprintf("%s", v)
}

func TestCookie60(t *testing.T) {
	buff := bytes.NewBufferString("")
	recorder := httptest.NewRecorder()
	recorder.Body = buff

	o := Classic()
	o.Get("/", new(Cookie60Action))

	req, err := http.NewRequest("GET", "http://localhost:8000/", nil)
	if err != nil {
		t.Error(err)
	}
	req.AddCookie(NewCookie("test", "1"))

	o.ServeHTTP(recorder, req)
	expect(t, recorder.Code, http.StatusOK)
	refute(t, len(buff.String()), 0)
	expect(t, buff.String(), "1")
}
