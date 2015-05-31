package tango

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

type Form1Action struct {
	Ctx
}

func (a *Form1Action) Get() string {
	v, _ := a.Forms().Int("test")
	return fmt.Sprintf("%d", v)
}

func TestForm1(t *testing.T) {
	buff := bytes.NewBufferString("")
	recorder := httptest.NewRecorder()
	recorder.Body = buff

	o := Classic()
	o.Get("/", new(Form1Action))

	req, err := http.NewRequest("GET", "http://localhost:8000/?test=1", nil)
	if err != nil {
		t.Error(err)
	}

	o.ServeHTTP(recorder, req)
	expect(t, recorder.Code, http.StatusOK)
	refute(t, len(buff.String()), 0)
	expect(t, buff.String(), "1")
}

type Form2Action struct {
	Ctx
}

func (a *Form2Action) Get() string {
	v, _ := a.Forms().Int32("test")
	return fmt.Sprintf("%d", v)
}

func TestForm2(t *testing.T) {
	buff := bytes.NewBufferString("")
	recorder := httptest.NewRecorder()
	recorder.Body = buff

	o := Classic()
	o.Get("/", new(Form2Action))

	req, err := http.NewRequest("GET", "http://localhost:8000/?test=1", nil)
	if err != nil {
		t.Error(err)
	}

	o.ServeHTTP(recorder, req)
	expect(t, recorder.Code, http.StatusOK)
	refute(t, len(buff.String()), 0)
	expect(t, buff.String(), "1")
}

type Form3Action struct {
	Ctx
}

func (a *Form3Action) Get() string {
	v, _ := a.Forms().Int64("test")
	return fmt.Sprintf("%d", v)
}

func TestForm3(t *testing.T) {
	buff := bytes.NewBufferString("")
	recorder := httptest.NewRecorder()
	recorder.Body = buff

	o := Classic()
	o.Get("/", new(Form3Action))

	req, err := http.NewRequest("GET", "http://localhost:8000/?test=1", nil)
	if err != nil {
		t.Error(err)
	}

	o.ServeHTTP(recorder, req)
	expect(t, recorder.Code, http.StatusOK)
	refute(t, len(buff.String()), 0)
	expect(t, buff.String(), "1")
}

type Form4Action struct {
	Ctx
}

func (a *Form4Action) Get() string {
	v, _ := a.Forms().Uint("test")
	return fmt.Sprintf("%d", v)
}

func TestForm4(t *testing.T) {
	buff := bytes.NewBufferString("")
	recorder := httptest.NewRecorder()
	recorder.Body = buff

	o := Classic()
	o.Get("/", new(Form4Action))

	req, err := http.NewRequest("GET", "http://localhost:8000/?test=1", nil)
	if err != nil {
		t.Error(err)
	}

	o.ServeHTTP(recorder, req)
	expect(t, recorder.Code, http.StatusOK)
	refute(t, len(buff.String()), 0)
	expect(t, buff.String(), "1")
}

type Form5Action struct {
	Ctx
}

func (a *Form5Action) Get() string {
	v, _ := a.Forms().Uint32("test")
	return fmt.Sprintf("%d", v)
}

func TestForm5(t *testing.T) {
	buff := bytes.NewBufferString("")
	recorder := httptest.NewRecorder()
	recorder.Body = buff

	o := Classic()
	o.Get("/", new(Form5Action))

	req, err := http.NewRequest("GET", "http://localhost:8000/?test=1", nil)
	if err != nil {
		t.Error(err)
	}

	o.ServeHTTP(recorder, req)
	expect(t, recorder.Code, http.StatusOK)
	refute(t, len(buff.String()), 0)
	expect(t, buff.String(), "1")
}

type Form6Action struct {
	Ctx
}

func (a *Form6Action) Get() string {
	v, _ := a.Forms().Uint64("test")
	return fmt.Sprintf("%d", v)
}

func TestForm6(t *testing.T) {
	buff := bytes.NewBufferString("")
	recorder := httptest.NewRecorder()
	recorder.Body = buff

	o := Classic()
	o.Get("/", new(Form6Action))

	req, err := http.NewRequest("GET", "http://localhost:8000/?test=1", nil)
	if err != nil {
		t.Error(err)
	}

	o.ServeHTTP(recorder, req)
	expect(t, recorder.Code, http.StatusOK)
	refute(t, len(buff.String()), 0)
	expect(t, buff.String(), "1")
}

type Form7Action struct {
	Ctx
}

func (a *Form7Action) Get() string {
	v, _ := a.Forms().Float32("test")
	return fmt.Sprintf("%.2f", v)
}

func TestForm7(t *testing.T) {
	buff := bytes.NewBufferString("")
	recorder := httptest.NewRecorder()
	recorder.Body = buff

	o := Classic()
	o.Get("/", new(Form7Action))

	req, err := http.NewRequest("GET", "http://localhost:8000/?test=1", nil)
	if err != nil {
		t.Error(err)
	}

	o.ServeHTTP(recorder, req)
	expect(t, recorder.Code, http.StatusOK)
	refute(t, len(buff.String()), 0)
	expect(t, buff.String(), "1.00")
}

type Form8Action struct {
	Ctx
}

func (a *Form8Action) Get() string {
	v, _ := a.Forms().Float64("test")
	return fmt.Sprintf("%.2f", v)
}

func TestForm8(t *testing.T) {
	buff := bytes.NewBufferString("")
	recorder := httptest.NewRecorder()
	recorder.Body = buff

	o := Classic()
	o.Get("/", new(Form8Action))

	req, err := http.NewRequest("GET", "http://localhost:8000/?test=1", nil)
	if err != nil {
		t.Error(err)
	}

	o.ServeHTTP(recorder, req)
	expect(t, recorder.Code, http.StatusOK)
	refute(t, len(buff.String()), 0)
	expect(t, buff.String(), "1.00")
}

type Form9Action struct {
	Ctx
}

func (a *Form9Action) Get() string {
	v, _ := a.Forms().Bool("test")
	return fmt.Sprintf("%v", v)
}

func TestForm9(t *testing.T) {
	buff := bytes.NewBufferString("")
	recorder := httptest.NewRecorder()
	recorder.Body = buff

	o := Classic()
	o.Get("/", new(Form9Action))

	req, err := http.NewRequest("GET", "http://localhost:8000/?test=1", nil)
	if err != nil {
		t.Error(err)
	}

	o.ServeHTTP(recorder, req)
	expect(t, recorder.Code, http.StatusOK)
	refute(t, len(buff.String()), 0)
	expect(t, buff.String(), "true")
}

type Form10Action struct {
	Ctx
}

func (a *Form10Action) Get() string {
	v, _ := a.Forms().String("test")
	return fmt.Sprintf("%v", v)
}

func TestForm10(t *testing.T) {
	buff := bytes.NewBufferString("")
	recorder := httptest.NewRecorder()
	recorder.Body = buff

	o := Classic()
	o.Get("/", new(Form10Action))

	req, err := http.NewRequest("GET", "http://localhost:8000/?test=1", nil)
	if err != nil {
		t.Error(err)
	}

	o.ServeHTTP(recorder, req)
	expect(t, recorder.Code, http.StatusOK)
	refute(t, len(buff.String()), 0)
	expect(t, buff.String(), "1")
}

type Form11Action struct {
	Ctx
}

func (a *Form11Action) Get() string {
	v := a.Forms().MustInt("test")
	return fmt.Sprintf("%d", v)
}

func TestForm11(t *testing.T) {
	buff := bytes.NewBufferString("")
	recorder := httptest.NewRecorder()
	recorder.Body = buff

	o := Classic()
	o.Get("/", new(Form11Action))

	req, err := http.NewRequest("GET", "http://localhost:8000/?test=1", nil)
	if err != nil {
		t.Error(err)
	}

	o.ServeHTTP(recorder, req)
	expect(t, recorder.Code, http.StatusOK)
	refute(t, len(buff.String()), 0)
	expect(t, buff.String(), "1")
}

type Form12Action struct {
	Ctx
}

func (a *Form12Action) Get() string {
	v := a.Forms().MustInt32("test")
	return fmt.Sprintf("%d", v)
}

func TestForm12(t *testing.T) {
	buff := bytes.NewBufferString("")
	recorder := httptest.NewRecorder()
	recorder.Body = buff

	o := Classic()
	o.Get("/", new(Form12Action))

	req, err := http.NewRequest("GET", "http://localhost:8000/?test=1", nil)
	if err != nil {
		t.Error(err)
	}

	o.ServeHTTP(recorder, req)
	expect(t, recorder.Code, http.StatusOK)
	refute(t, len(buff.String()), 0)
	expect(t, buff.String(), "1")
}

type Form13Action struct {
	Ctx
}

func (a *Form13Action) Get() string {
	v := a.Forms().MustInt64("test")
	return fmt.Sprintf("%d", v)
}

func TestForm13(t *testing.T) {
	buff := bytes.NewBufferString("")
	recorder := httptest.NewRecorder()
	recorder.Body = buff

	o := Classic()
	o.Get("/", new(Form13Action))

	req, err := http.NewRequest("GET", "http://localhost:8000/?test=1", nil)
	if err != nil {
		t.Error(err)
	}

	o.ServeHTTP(recorder, req)
	expect(t, recorder.Code, http.StatusOK)
	refute(t, len(buff.String()), 0)
	expect(t, buff.String(), "1")
}

type Form14Action struct {
	Ctx
}

func (a *Form14Action) Get() string {
	v := a.Forms().MustUint("test")
	return fmt.Sprintf("%d", v)
}

func TestForm14(t *testing.T) {
	buff := bytes.NewBufferString("")
	recorder := httptest.NewRecorder()
	recorder.Body = buff

	o := Classic()
	o.Get("/", new(Form14Action))

	req, err := http.NewRequest("GET", "http://localhost:8000/?test=1", nil)
	if err != nil {
		t.Error(err)
	}

	o.ServeHTTP(recorder, req)
	expect(t, recorder.Code, http.StatusOK)
	refute(t, len(buff.String()), 0)
	expect(t, buff.String(), "1")
}

type Form15Action struct {
	Ctx
}

func (a *Form15Action) Get() string {
	v := a.Forms().MustUint32("test")
	return fmt.Sprintf("%d", v)
}

func TestForm15(t *testing.T) {
	buff := bytes.NewBufferString("")
	recorder := httptest.NewRecorder()
	recorder.Body = buff

	o := Classic()
	o.Get("/", new(Form15Action))

	req, err := http.NewRequest("GET", "http://localhost:8000/?test=1", nil)
	if err != nil {
		t.Error(err)
	}

	o.ServeHTTP(recorder, req)
	expect(t, recorder.Code, http.StatusOK)
	refute(t, len(buff.String()), 0)
	expect(t, buff.String(), "1")
}

type Form16Action struct {
	Ctx
}

func (a *Form16Action) Get() string {
	v := a.Forms().MustUint64("test")
	return fmt.Sprintf("%d", v)
}

func TestForm16(t *testing.T) {
	buff := bytes.NewBufferString("")
	recorder := httptest.NewRecorder()
	recorder.Body = buff

	o := Classic()
	o.Get("/", new(Form16Action))

	req, err := http.NewRequest("GET", "http://localhost:8000/?test=1", nil)
	if err != nil {
		t.Error(err)
	}

	o.ServeHTTP(recorder, req)
	expect(t, recorder.Code, http.StatusOK)
	refute(t, len(buff.String()), 0)
	expect(t, buff.String(), "1")
}

type Form17Action struct {
	Ctx
}

func (a *Form17Action) Get() string {
	v := a.Forms().MustFloat32("test")
	return fmt.Sprintf("%.2f", v)
}

func TestForm17(t *testing.T) {
	buff := bytes.NewBufferString("")
	recorder := httptest.NewRecorder()
	recorder.Body = buff

	o := Classic()
	o.Get("/", new(Form17Action))

	req, err := http.NewRequest("GET", "http://localhost:8000/?test=1", nil)
	if err != nil {
		t.Error(err)
	}

	o.ServeHTTP(recorder, req)
	expect(t, recorder.Code, http.StatusOK)
	refute(t, len(buff.String()), 0)
	expect(t, buff.String(), "1.00")
}

type Form18Action struct {
	Ctx
}

func (a *Form18Action) Get() string {
	v := a.Forms().MustFloat64("test")
	return fmt.Sprintf("%.2f", v)
}

func TestForm18(t *testing.T) {
	buff := bytes.NewBufferString("")
	recorder := httptest.NewRecorder()
	recorder.Body = buff

	o := Classic()
	o.Get("/", new(Form18Action))

	req, err := http.NewRequest("GET", "http://localhost:8000/?test=1", nil)
	if err != nil {
		t.Error(err)
	}

	o.ServeHTTP(recorder, req)
	expect(t, recorder.Code, http.StatusOK)
	refute(t, len(buff.String()), 0)
	expect(t, buff.String(), "1.00")
}

type Form19Action struct {
	Ctx
}

func (a *Form19Action) Get() string {
	v := a.Forms().MustBool("test")
	return fmt.Sprintf("%v", v)
}

func TestForm19(t *testing.T) {
	buff := bytes.NewBufferString("")
	recorder := httptest.NewRecorder()
	recorder.Body = buff

	o := Classic()
	o.Get("/", new(Form19Action))

	req, err := http.NewRequest("GET", "http://localhost:8000/?test=1", nil)
	if err != nil {
		t.Error(err)
	}

	o.ServeHTTP(recorder, req)
	expect(t, recorder.Code, http.StatusOK)
	refute(t, len(buff.String()), 0)
	expect(t, buff.String(), "true")
}

type Form20Action struct {
	Ctx
}

func (a *Form20Action) Get() string {
	v := a.Forms().MustString("test")
	return fmt.Sprintf("%s", v)
}

func TestForm20(t *testing.T) {
	buff := bytes.NewBufferString("")
	recorder := httptest.NewRecorder()
	recorder.Body = buff

	o := Classic()
	o.Get("/", new(Form20Action))

	req, err := http.NewRequest("GET", "http://localhost:8000/?test=1", nil)
	if err != nil {
		t.Error(err)
	}

	o.ServeHTTP(recorder, req)
	expect(t, recorder.Code, http.StatusOK)
	refute(t, len(buff.String()), 0)
	expect(t, buff.String(), "1")
}