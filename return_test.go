// Copyright 2015 The Tango Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package tango

import (
	"bytes"
	"encoding/xml"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

type MyReturn struct {
}

func (m MyReturn) Get() string {
	return "string return"
}

func (m MyReturn) Post() []byte {
	return []byte("bytes return")
}

func (m MyReturn) Put() error {
	return errors.New("error return")
}

func TestReturn(t *testing.T) {
	buff := bytes.NewBufferString("")
	recorder := httptest.NewRecorder()
	recorder.Body = buff

	o := Classic()
	o.Any("/", new(MyReturn))

	req, err := http.NewRequest("GET", "http://localhost:8000/", nil)
	if err != nil {
		t.Error(err)
	}

	o.ServeHTTP(recorder, req)
	expect(t, recorder.Code, http.StatusOK)
	refute(t, len(buff.String()), 0)
	expect(t, buff.String(), "string return")

	buff.Reset()
	req, err = http.NewRequest("POST", "http://localhost:8000/", nil)
	if err != nil {
		t.Error(err)
	}

	o.ServeHTTP(recorder, req)
	expect(t, recorder.Code, http.StatusOK)
	refute(t, len(buff.String()), 0)
	expect(t, buff.String(), "bytes return")
}

func TestReturnPut(t *testing.T) {
	buff := bytes.NewBufferString("")
	recorder := httptest.NewRecorder()
	recorder.Body = buff

	o := Classic()
	o.Any("/", new(MyReturn))

	req, err := http.NewRequest("PUT", "http://localhost:8000/", nil)
	if err != nil {
		t.Error(err)
	}

	o.ServeHTTP(recorder, req)
	expect(t, recorder.Code, http.StatusInternalServerError)
	refute(t, len(buff.String()), 0)
	expect(t, buff.String(), "error return")
}

type JsonReturn struct {
	Json
}

func (JsonReturn) Get() interface{} {
	return map[string]interface{}{
		"test1": 1,
		"test2": "2",
		"test3": true,
	}
}

func TestReturnJson1(t *testing.T) {
	buff := bytes.NewBufferString("")
	recorder := httptest.NewRecorder()
	recorder.Body = buff

	o := Classic()
	o.Get("/", new(JsonReturn))

	req, err := http.NewRequest("GET", "http://localhost:8000/", nil)
	if err != nil {
		t.Error(err)
	}

	o.ServeHTTP(recorder, req)
	expect(t, recorder.Code, http.StatusOK)
	refute(t, len(buff.String()), 0)
	expect(t, strings.TrimSpace(buff.String()), `{"test1":1,"test2":"2","test3":true}`)
}

type JsonErrReturn struct {
	Json
}

func (JsonErrReturn) Get() error {
	return errors.New("error")
}

func TestReturnJsonError(t *testing.T) {
	buff := bytes.NewBufferString("")
	recorder := httptest.NewRecorder()
	recorder.Body = buff

	o := Classic()
	o.Get("/", new(JsonErrReturn))

	req, err := http.NewRequest("GET", "http://localhost:8000/", nil)
	if err != nil {
		t.Error(err)
	}

	o.ServeHTTP(recorder, req)
	expect(t, recorder.Code, http.StatusOK)
	refute(t, len(buff.String()), 0)
	expect(t, strings.TrimSpace(buff.String()), `{"err":"error"}`)
}

type JsonErrReturn2 struct {
	Json
}

func (JsonErrReturn2) Get() error {
	return Abort(http.StatusInternalServerError, "error")
}

func TestReturnJsonError2(t *testing.T) {
	buff := bytes.NewBufferString("")
	recorder := httptest.NewRecorder()
	recorder.Body = buff

	o := Classic()
	o.Get("/", new(JsonErrReturn2))

	req, err := http.NewRequest("GET", "http://localhost:8000/", nil)
	if err != nil {
		t.Error(err)
	}

	o.ServeHTTP(recorder, req)
	expect(t, recorder.Code, http.StatusInternalServerError)
	refute(t, len(buff.String()), 0)
	expect(t, strings.TrimSpace(buff.String()), `{"err":"error"}`)
}

type JsonReturn1 struct {
	Json
}

func (JsonReturn1) Get() string {
	return "return"
}

func TestReturnJson2(t *testing.T) {
	buff := bytes.NewBufferString("")
	recorder := httptest.NewRecorder()
	recorder.Body = buff

	o := Classic()
	o.Get("/", new(JsonReturn1))

	req, err := http.NewRequest("GET", "http://localhost:8000/", nil)
	if err != nil {
		t.Error(err)
	}

	o.ServeHTTP(recorder, req)
	expect(t, recorder.Code, http.StatusOK)
	refute(t, len(buff.String()), 0)
	expect(t, strings.TrimSpace(buff.String()), `{"content":"return"}`)
}

type JsonReturn2 struct {
	Json
}

func (JsonReturn2) Get() []byte {
	return []byte("return")
}

func TestReturnJson3(t *testing.T) {
	buff := bytes.NewBufferString("")
	recorder := httptest.NewRecorder()
	recorder.Body = buff

	o := Classic()
	o.Get("/", new(JsonReturn2))

	req, err := http.NewRequest("GET", "http://localhost:8000/", nil)
	if err != nil {
		t.Error(err)
	}

	o.ServeHTTP(recorder, req)
	expect(t, recorder.Code, http.StatusOK)
	refute(t, len(buff.String()), 0)
	expect(t, strings.TrimSpace(buff.String()), `{"content":"return"}`)
}

type JsonReturn3 struct {
    Json
}

func (JsonReturn3) Get() (int, interface{}) {
    if true {
        return 201, map[string]string{
            "say": "Hello tango!",
        }
    }
    return 500, errors.New("something error")
}

func TestReturnJson4(t *testing.T) {
	buff := bytes.NewBufferString("")
	recorder := httptest.NewRecorder()
	recorder.Body = buff

	o := Classic()
	o.Get("/", new(JsonReturn3))

	req, err := http.NewRequest("GET", "http://localhost:8000/", nil)
	if err != nil {
		t.Error(err)
	}

	o.ServeHTTP(recorder, req)
	expect(t, recorder.Code, 201)
	refute(t, len(buff.String()), 0)
	expect(t, strings.TrimSpace(buff.String()), `{"say":"Hello tango!"}`)
}

type XmlReturn struct {
	Xml
}

type Address struct {
	City, State string
}
type Person struct {
	XMLName   xml.Name `xml:"person"`
	Id        int      `xml:"id,attr"`
	FirstName string   `xml:"name>first"`
	LastName  string   `xml:"name>last"`
	Age       int      `xml:"age"`
	Height    float32  `xml:"height,omitempty"`
	Married   bool
	Address
	Comment string `xml:",comment"`
}

func (XmlReturn) Get() interface{} {
	v := &Person{Id: 13, FirstName: "John", LastName: "Doe", Age: 42}
	v.Comment = " Need more details. "
	v.Address = Address{"Hanga Roa", "Easter Island"}
	return v
}

func TestReturnXml(t *testing.T) {
	buff := bytes.NewBufferString("")
	recorder := httptest.NewRecorder()
	recorder.Body = buff

	o := Classic()
	o.Get("/", new(XmlReturn))

	req, err := http.NewRequest("GET", "http://localhost:8000/", nil)
	if err != nil {
		t.Error(err)
	}

	o.ServeHTTP(recorder, req)
	expect(t, recorder.Code, http.StatusOK)
	refute(t, len(buff.String()), 0)
	expect(t, buff.String(), `<person id="13"><name><first>John</first><last>Doe</last></name><age>42</age><Married>false</Married><City>Hanga Roa</City><State>Easter Island</State><!-- Need more details. --></person>`)
}

type XmlErrReturn struct {
	Xml
}

func (XmlErrReturn) Get() error {
	return errors.New("error")
}

func TestReturnXmlError(t *testing.T) {
	buff := bytes.NewBufferString("")
	recorder := httptest.NewRecorder()
	recorder.Body = buff

	o := Classic()
	o.Get("/", new(XmlErrReturn))

	req, err := http.NewRequest("GET", "http://localhost:8000/", nil)
	if err != nil {
		t.Error(err)
	}

	o.ServeHTTP(recorder, req)
	expect(t, recorder.Code, http.StatusOK)
	refute(t, len(buff.String()), 0)
	expect(t, strings.TrimSpace(buff.String()), `<err><content>error</content></err>`)
}
