package tango

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"
	"errors"
	"encoding/xml"
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
	o.Get("/", new(MyResponse))
	
	req, err := http.NewRequest("GET", "http://localhost:8000/", nil)
	if err != nil {
		t.Error(err)
	}

	o.ServeHTTP(recorder, req)
	expect(t, recorder.Code, http.StatusOK)
	refute(t, len(buff.String()), 0)
	expect(t, buff.String(), "my response")
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

func TestReturnJson(t *testing.T) {
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
	expect(t, buff.String(), `{"test1":1,"test2":"2","test3":true}`)
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
