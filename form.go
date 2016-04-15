package tango

import (
	"errors"
	"html/template"
	"net/http"
	"net/url"
	"strconv"
)

// Forms a new enhancement of http.Request
type Forms http.Request

var _ Set = &Forms{}

func (f *Forms) Values() url.Values {
	return (*http.Request)(f).Form
}

func (f *Forms) String(key string) (string, error) {
	return (*http.Request)(f).FormValue(key), nil
}

func (f *Forms) Strings(key string) ([]string, error) {
	(*http.Request)(f).ParseMultipartForm(32 << 20)
	if v, ok := (*http.Request)(f).Form[key]; ok {
		return v, nil
	}
	return nil, errors.New("not exist")
}

func (f *Forms) Escape(key string) (string, error) {
	return template.HTMLEscapeString((*http.Request)(f).FormValue(key)), nil
}

func (f *Forms) Int(key string) (int, error) {
	return strconv.Atoi((*http.Request)(f).FormValue(key))
}

func (f *Forms) Int32(key string) (int32, error) {
	v, err := strconv.ParseInt((*http.Request)(f).FormValue(key), 10, 32)
	return int32(v), err
}

func (f *Forms) Int64(key string) (int64, error) {
	return strconv.ParseInt((*http.Request)(f).FormValue(key), 10, 64)
}

func (f *Forms) Uint(key string) (uint, error) {
	v, err := strconv.ParseUint((*http.Request)(f).FormValue(key), 10, 64)
	return uint(v), err
}

func (f *Forms) Uint32(key string) (uint32, error) {
	v, err := strconv.ParseUint((*http.Request)(f).FormValue(key), 10, 32)
	return uint32(v), err
}

func (f *Forms) Uint64(key string) (uint64, error) {
	return strconv.ParseUint((*http.Request)(f).FormValue(key), 10, 64)
}

func (f *Forms) Bool(key string) (bool, error) {
	return strconv.ParseBool((*http.Request)(f).FormValue(key))
}

func (f *Forms) Float32(key string) (float32, error) {
	v, err := strconv.ParseFloat((*http.Request)(f).FormValue(key), 64)
	return float32(v), err
}

func (f *Forms) Float64(key string) (float64, error) {
	return strconv.ParseFloat((*http.Request)(f).FormValue(key), 64)
}

func (f *Forms) MustString(key string, defaults ...string) string {
	if v := (*http.Request)(f).FormValue(key); len(v) > 0 {
		return v
	}
	if len(defaults) > 0 {
		return defaults[0]
	}
	return ""
}

func (f *Forms) MustStrings(key string, defaults ...[]string) []string {
	(*http.Request)(f).ParseMultipartForm(32 << 20)
	if v, ok := (*http.Request)(f).Form[key]; ok {
		return v
	}
	if len(defaults) > 0 {
		return defaults[0]
	}
	return []string{}
}

func (f *Forms) MustEscape(key string, defaults ...string) string {
	if v := (*http.Request)(f).FormValue(key); len(v) > 0 {
		return template.HTMLEscapeString(v)
	}
	if len(defaults) > 0 {
		return defaults[0]
	}
	return ""
}

func (f *Forms) MustInt(key string, defaults ...int) int {
	v, err := strconv.Atoi((*http.Request)(f).FormValue(key))
	if len(defaults) > 0 && err != nil {
		return defaults[0]
	}
	return v
}

func (f *Forms) MustInt32(key string, defaults ...int32) int32 {
	v, err := strconv.ParseInt((*http.Request)(f).FormValue(key), 10, 32)
	if len(defaults) > 0 && err != nil {
		return defaults[0]
	}
	return int32(v)
}

func (f *Forms) MustInt64(key string, defaults ...int64) int64 {
	v, err := strconv.ParseInt((*http.Request)(f).FormValue(key), 10, 64)
	if len(defaults) > 0 && err != nil {
		return defaults[0]
	}
	return v
}

func (f *Forms) MustUint(key string, defaults ...uint) uint {
	v, err := strconv.ParseUint((*http.Request)(f).FormValue(key), 10, 64)
	if len(defaults) > 0 && err != nil {
		return defaults[0]
	}
	return uint(v)
}

func (f *Forms) MustUint32(key string, defaults ...uint32) uint32 {
	v, err := strconv.ParseUint((*http.Request)(f).FormValue(key), 10, 32)
	if len(defaults) > 0 && err != nil {
		return defaults[0]
	}
	return uint32(v)
}

func (f *Forms) MustUint64(key string, defaults ...uint64) uint64 {
	v, err := strconv.ParseUint((*http.Request)(f).FormValue(key), 10, 64)
	if len(defaults) > 0 && err != nil {
		return defaults[0]
	}
	return v
}

func (f *Forms) MustFloat32(key string, defaults ...float32) float32 {
	v, err := strconv.ParseFloat((*http.Request)(f).FormValue(key), 32)
	if len(defaults) > 0 && err != nil {
		return defaults[0]
	}
	return float32(v)
}

func (f *Forms) MustFloat64(key string, defaults ...float64) float64 {
	v, err := strconv.ParseFloat((*http.Request)(f).FormValue(key), 64)
	if len(defaults) > 0 && err != nil {
		return defaults[0]
	}
	return v
}

func (f *Forms) MustBool(key string, defaults ...bool) bool {
	v, err := strconv.ParseBool((*http.Request)(f).FormValue(key))
	if len(defaults) > 0 && err != nil {
		return defaults[0]
	}
	return v
}

func (ctx *Context) Form(key string, defaults ...string) string {
	return (*Forms)(ctx.req).MustString(key, defaults...)
}

func (ctx *Context) FormStrings(key string, defaults ...[]string) []string {
	return (*Forms)(ctx.req).MustStrings(key, defaults...)
}

func (ctx *Context) FormEscape(key string, defaults ...string) string {
	return (*Forms)(ctx.req).MustEscape(key, defaults...)
}

func (ctx *Context) FormInt(key string, defaults ...int) int {
	return (*Forms)(ctx.req).MustInt(key, defaults...)
}

func (ctx *Context) FormInt32(key string, defaults ...int32) int32 {
	return (*Forms)(ctx.req).MustInt32(key, defaults...)
}

func (ctx *Context) FormInt64(key string, defaults ...int64) int64 {
	return (*Forms)(ctx.req).MustInt64(key, defaults...)
}

func (ctx *Context) FormUint(key string, defaults ...uint) uint {
	return (*Forms)(ctx.req).MustUint(key, defaults...)
}

func (ctx *Context) FormUint32(key string, defaults ...uint32) uint32 {
	return (*Forms)(ctx.req).MustUint32(key, defaults...)
}

func (ctx *Context) FormUint64(key string, defaults ...uint64) uint64 {
	return (*Forms)(ctx.req).MustUint64(key, defaults...)
}

func (ctx *Context) FormFloat32(key string, defaults ...float32) float32 {
	return (*Forms)(ctx.req).MustFloat32(key, defaults...)
}

func (ctx *Context) FormFloat64(key string, defaults ...float64) float64 {
	return (*Forms)(ctx.req).MustFloat64(key, defaults...)
}

func (ctx *Context) FormBool(key string, defaults ...bool) bool {
	return (*Forms)(ctx.req).MustBool(key, defaults...)
}
