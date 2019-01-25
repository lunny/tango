// Copyright 2016 The Tango Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

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

// Values returns http.Request values
func (f *Forms) Values() url.Values {
	return (*http.Request)(f).Form
}

// String returns request form as string
func (f *Forms) String(key string) (string, error) {
	return (*http.Request)(f).FormValue(key), nil
}

// Strings returns request form as strings
func (f *Forms) Strings(key string) ([]string, error) {
	(*http.Request)(f).ParseMultipartForm(32 << 20)
	if v, ok := (*http.Request)(f).Form[key]; ok {
		return v, nil
	}
	return nil, errors.New("not exist")
}

// Escape returns request form as escaped string
func (f *Forms) Escape(key string) (string, error) {
	return template.HTMLEscapeString((*http.Request)(f).FormValue(key)), nil
}

// Int returns request form as int
func (f *Forms) Int(key string) (int, error) {
	return strconv.Atoi((*http.Request)(f).FormValue(key))
}

// Int32 returns request form as int32
func (f *Forms) Int32(key string) (int32, error) {
	v, err := strconv.ParseInt((*http.Request)(f).FormValue(key), 10, 32)
	return int32(v), err
}

// Int64 returns request form as int64
func (f *Forms) Int64(key string) (int64, error) {
	return strconv.ParseInt((*http.Request)(f).FormValue(key), 10, 64)
}

// Uint returns request form as uint
func (f *Forms) Uint(key string) (uint, error) {
	v, err := strconv.ParseUint((*http.Request)(f).FormValue(key), 10, 64)
	return uint(v), err
}

// Uint32 returns request form as uint32
func (f *Forms) Uint32(key string) (uint32, error) {
	v, err := strconv.ParseUint((*http.Request)(f).FormValue(key), 10, 32)
	return uint32(v), err
}

// Uint64 returns request form as uint64
func (f *Forms) Uint64(key string) (uint64, error) {
	return strconv.ParseUint((*http.Request)(f).FormValue(key), 10, 64)
}

// Bool returns request form as bool
func (f *Forms) Bool(key string) (bool, error) {
	return strconv.ParseBool((*http.Request)(f).FormValue(key))
}

// Float32 returns request form as float32
func (f *Forms) Float32(key string) (float32, error) {
	v, err := strconv.ParseFloat((*http.Request)(f).FormValue(key), 64)
	return float32(v), err
}

// Float64 returns request form as float64
func (f *Forms) Float64(key string) (float64, error) {
	return strconv.ParseFloat((*http.Request)(f).FormValue(key), 64)
}

// MustString returns request form as string with default
func (f *Forms) MustString(key string, defaults ...string) string {
	if v := (*http.Request)(f).FormValue(key); len(v) > 0 {
		return v
	}
	if len(defaults) > 0 {
		return defaults[0]
	}
	return ""
}

// MustStrings returns request form as strings with default
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

// MustEscape returns request form as escaped string with default
func (f *Forms) MustEscape(key string, defaults ...string) string {
	if v := (*http.Request)(f).FormValue(key); len(v) > 0 {
		return template.HTMLEscapeString(v)
	}
	if len(defaults) > 0 {
		return defaults[0]
	}
	return ""
}

// MustInt returns request form as int with default
func (f *Forms) MustInt(key string, defaults ...int) int {
	v, err := strconv.Atoi((*http.Request)(f).FormValue(key))
	if len(defaults) > 0 && err != nil {
		return defaults[0]
	}
	return v
}

// MustInt32 returns request form as int32 with default
func (f *Forms) MustInt32(key string, defaults ...int32) int32 {
	v, err := strconv.ParseInt((*http.Request)(f).FormValue(key), 10, 32)
	if len(defaults) > 0 && err != nil {
		return defaults[0]
	}
	return int32(v)
}

// MustInt64 returns request form as int64 with default
func (f *Forms) MustInt64(key string, defaults ...int64) int64 {
	v, err := strconv.ParseInt((*http.Request)(f).FormValue(key), 10, 64)
	if len(defaults) > 0 && err != nil {
		return defaults[0]
	}
	return v
}

// MustUint returns request form as uint with default
func (f *Forms) MustUint(key string, defaults ...uint) uint {
	v, err := strconv.ParseUint((*http.Request)(f).FormValue(key), 10, 64)
	if len(defaults) > 0 && err != nil {
		return defaults[0]
	}
	return uint(v)
}

// MustUint32 returns request form as uint32 with default
func (f *Forms) MustUint32(key string, defaults ...uint32) uint32 {
	v, err := strconv.ParseUint((*http.Request)(f).FormValue(key), 10, 32)
	if len(defaults) > 0 && err != nil {
		return defaults[0]
	}
	return uint32(v)
}

// MustUint64 returns request form as uint64 with default
func (f *Forms) MustUint64(key string, defaults ...uint64) uint64 {
	v, err := strconv.ParseUint((*http.Request)(f).FormValue(key), 10, 64)
	if len(defaults) > 0 && err != nil {
		return defaults[0]
	}
	return v
}

// MustFloat32 returns request form as float32 with default
func (f *Forms) MustFloat32(key string, defaults ...float32) float32 {
	v, err := strconv.ParseFloat((*http.Request)(f).FormValue(key), 32)
	if len(defaults) > 0 && err != nil {
		return defaults[0]
	}
	return float32(v)
}

// MustFloat64 returns request form as float64 with default
func (f *Forms) MustFloat64(key string, defaults ...float64) float64 {
	v, err := strconv.ParseFloat((*http.Request)(f).FormValue(key), 64)
	if len(defaults) > 0 && err != nil {
		return defaults[0]
	}
	return v
}

// MustBool returns request form as bool with default
func (f *Forms) MustBool(key string, defaults ...bool) bool {
	v, err := strconv.ParseBool((*http.Request)(f).FormValue(key))
	if len(defaults) > 0 && err != nil {
		return defaults[0]
	}
	return v
}

// Form returns request form as string with default
func (ctx *Context) Form(key string, defaults ...string) string {
	return (*Forms)(ctx.req).MustString(key, defaults...)
}

// FormStrings returns request form as strings with default
func (ctx *Context) FormStrings(key string, defaults ...[]string) []string {
	return (*Forms)(ctx.req).MustStrings(key, defaults...)
}

// FormEscape returns request form as escaped string with default
func (ctx *Context) FormEscape(key string, defaults ...string) string {
	return (*Forms)(ctx.req).MustEscape(key, defaults...)
}

// FormInt returns request form as int with default
func (ctx *Context) FormInt(key string, defaults ...int) int {
	return (*Forms)(ctx.req).MustInt(key, defaults...)
}

// FormInt32 returns request form as int32 with default
func (ctx *Context) FormInt32(key string, defaults ...int32) int32 {
	return (*Forms)(ctx.req).MustInt32(key, defaults...)
}

// FormInt64 returns request form as int64 with default
func (ctx *Context) FormInt64(key string, defaults ...int64) int64 {
	return (*Forms)(ctx.req).MustInt64(key, defaults...)
}

// FormUint returns request form as uint with default
func (ctx *Context) FormUint(key string, defaults ...uint) uint {
	return (*Forms)(ctx.req).MustUint(key, defaults...)
}

// FormUint32 returns request form as uint32 with default
func (ctx *Context) FormUint32(key string, defaults ...uint32) uint32 {
	return (*Forms)(ctx.req).MustUint32(key, defaults...)
}

// FormUint64 returns request form as uint64 with default
func (ctx *Context) FormUint64(key string, defaults ...uint64) uint64 {
	return (*Forms)(ctx.req).MustUint64(key, defaults...)
}

// FormFloat32 returns request form as float32 with default
func (ctx *Context) FormFloat32(key string, defaults ...float32) float32 {
	return (*Forms)(ctx.req).MustFloat32(key, defaults...)
}

// FormFloat64 returns request form as float64 with default
func (ctx *Context) FormFloat64(key string, defaults ...float64) float64 {
	return (*Forms)(ctx.req).MustFloat64(key, defaults...)
}

// FormBool returns request form as bool with default
func (ctx *Context) FormBool(key string, defaults ...bool) bool {
	return (*Forms)(ctx.req).MustBool(key, defaults...)
}
