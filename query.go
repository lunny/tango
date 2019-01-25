// Copyright 2019 The Tango Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package tango

import (
	"html/template"
	"net/http"
	"net/url"
	"strconv"
)

// Queries a new enhancement of http.Request
type Queries http.Request

var _ Set = &Queries{}

// Values returns http.Request values
func (f *Queries) Values() url.Values {
	return (*http.Request)(f).URL.Query()
}

// String returns request form as string
func (f *Queries) String(key string) (string, error) {
	if v, ok := f.Values()[key]; ok {
		return v[0], nil
	}
	return "", ErrorKeyIsNotExist{key}
}

// Strings returns request form as strings
func (f *Queries) Strings(key string) ([]string, error) {
	if v, ok := f.Values()[key]; ok {
		return v, nil
	}
	return nil, ErrorKeyIsNotExist{key}
}

// Escape returns request form as escaped string
func (f *Queries) Escape(key string) (string, error) {
	s, err := f.String(key)
	if err != nil {
		return "", err
	}
	return template.HTMLEscapeString(s), nil
}

// Int returns request form as int
func (f *Queries) Int(key string) (int, error) {
	s, err := f.String(key)
	if err != nil {
		return 0, err
	}
	return strconv.Atoi(s)
}

// Int32 returns request form as int32
func (f *Queries) Int32(key string) (int32, error) {
	s, err := f.String(key)
	if err != nil {
		return 0, err
	}
	v, err := strconv.ParseInt(s, 10, 32)
	return int32(v), err
}

// Int64 returns request form as int64
func (f *Queries) Int64(key string) (int64, error) {
	s, err := f.String(key)
	if err != nil {
		return 0, err
	}
	v, err := strconv.ParseInt(s, 10, 64)
	return v, err
}

// Uint returns request form as uint
func (f *Queries) Uint(key string) (uint, error) {
	s, err := f.String(key)
	if err != nil {
		return 0, err
	}
	v, err := strconv.ParseUint(s, 10, 64)
	return uint(v), err
}

// Uint32 returns request form as uint32
func (f *Queries) Uint32(key string) (uint32, error) {
	s, err := f.String(key)
	if err != nil {
		return 0, err
	}
	v, err := strconv.ParseUint(s, 10, 32)
	return uint32(v), err
}

// Uint64 returns request form as uint64
func (f *Queries) Uint64(key string) (uint64, error) {
	s, err := f.String(key)
	if err != nil {
		return 0, err
	}
	return strconv.ParseUint(s, 10, 64)
}

// Bool returns request form as bool
func (f *Queries) Bool(key string) (bool, error) {
	s, err := f.String(key)
	if err != nil {
		return false, err
	}
	return strconv.ParseBool(s)
}

// Float32 returns request form as float32
func (f *Queries) Float32(key string) (float32, error) {
	s, err := f.String(key)
	if err != nil {
		return 0, err
	}
	v, err := strconv.ParseFloat(s, 64)
	return float32(v), err
}

// Float64 returns request form as float64
func (f *Queries) Float64(key string) (float64, error) {
	s, err := f.String(key)
	if err != nil {
		return 0, err
	}
	return strconv.ParseFloat(s, 64)
}

// MustString returns request form as string with default
func (f *Queries) MustString(key string, defaults ...string) string {
	if v, ok := f.Values()[key]; ok {
		return v[0]
	}
	if len(defaults) > 0 {
		return defaults[0]
	}
	return ""
}

// MustStrings returns request form as strings with default
func (f *Queries) MustStrings(key string, defaults ...[]string) []string {
	(*http.Request)(f).ParseMultipartForm(32 << 20)
	if v, ok := f.Values()[key]; ok {
		return v
	}
	if len(defaults) > 0 {
		return defaults[0]
	}
	return []string{}
}

// MustEscape returns request form as escaped string with default
func (f *Queries) MustEscape(key string, defaults ...string) string {
	if v, ok := f.Values()[key]; ok {
		return template.HTMLEscapeString(v[0])
	}
	if len(defaults) > 0 {
		return defaults[0]
	}
	return ""
}

// MustInt returns request form as int with default
func (f *Queries) MustInt(key string, defaults ...int) int {
	if v, ok := f.Values()[key]; ok {
		i, err := strconv.Atoi(v[0])
		if err == nil {
			return i
		}
	}
	if len(defaults) > 0 {
		return defaults[0]
	}
	return 0
}

// MustInt32 returns request form as int32 with default
func (f *Queries) MustInt32(key string, defaults ...int32) int32 {
	if v, ok := f.Values()[key]; ok {
		i, err := strconv.ParseInt(v[0], 10, 32)
		if err == nil {
			return int32(i)
		}
	}
	if len(defaults) > 0 {
		return defaults[0]
	}
	return 0
}

// MustInt64 returns request form as int64 with default
func (f *Queries) MustInt64(key string, defaults ...int64) int64 {
	if v, ok := f.Values()[key]; ok {
		i, err := strconv.ParseInt(v[0], 10, 64)
		if err == nil {
			return i
		}
	}
	if len(defaults) > 0 {
		return defaults[0]
	}
	return 0
}

// MustUint returns request form as uint with default
func (f *Queries) MustUint(key string, defaults ...uint) uint {
	if v, ok := f.Values()[key]; ok {
		i, err := strconv.ParseUint(v[0], 10, 64)
		if err == nil {
			return uint(i)
		}
	}
	if len(defaults) > 0 {
		return defaults[0]
	}
	return 0
}

// MustUint32 returns request form as uint32 with default
func (f *Queries) MustUint32(key string, defaults ...uint32) uint32 {
	if v, ok := f.Values()[key]; ok {
		i, err := strconv.ParseUint(v[0], 10, 32)
		if err == nil {
			return uint32(i)
		}
	}
	if len(defaults) > 0 {
		return defaults[0]
	}
	return 0
}

// MustUint64 returns request form as uint64 with default
func (f *Queries) MustUint64(key string, defaults ...uint64) uint64 {
	if v, ok := f.Values()[key]; ok {
		i, err := strconv.ParseUint(v[0], 10, 64)
		if err == nil {
			return i
		}
	}
	if len(defaults) > 0 {
		return defaults[0]
	}
	return 0
}

// MustFloat32 returns request form as float32 with default
func (f *Queries) MustFloat32(key string, defaults ...float32) float32 {
	if v, ok := f.Values()[key]; ok {
		i, err := strconv.ParseFloat(v[0], 32)
		if err == nil {
			return float32(i)
		}
	}
	if len(defaults) > 0 {
		return defaults[0]
	}
	return 0
}

// MustFloat64 returns request form as float64 with default
func (f *Queries) MustFloat64(key string, defaults ...float64) float64 {
	if v, ok := f.Values()[key]; ok {
		i, err := strconv.ParseFloat(v[0], 64)
		if err == nil {
			return i
		}
	}
	if len(defaults) > 0 {
		return defaults[0]
	}
	return 0
}

// MustBool returns request form as bool with default
func (f *Queries) MustBool(key string, defaults ...bool) bool {
	if v, ok := f.Values()[key]; ok {
		i, err := strconv.ParseBool(v[0])
		if err == nil {
			return i
		}
	}
	if len(defaults) > 0 {
		return defaults[0]
	}
	return false
}

// Query returns request form as string with default
func (ctx *Context) Query(key string, defaults ...string) string {
	return (*Queries)(ctx.req).MustString(key, defaults...)
}

// QueryStrings returns request form as strings with default
func (ctx *Context) QueryStrings(key string, defaults ...[]string) []string {
	return (*Queries)(ctx.req).MustStrings(key, defaults...)
}

// QueryEscape returns request form as escaped string with default
func (ctx *Context) QueryEscape(key string, defaults ...string) string {
	return (*Queries)(ctx.req).MustEscape(key, defaults...)
}

// QueryInt returns request form as int with default
func (ctx *Context) QueryInt(key string, defaults ...int) int {
	return (*Queries)(ctx.req).MustInt(key, defaults...)
}

// QueryInt32 returns request form as int32 with default
func (ctx *Context) QueryInt32(key string, defaults ...int32) int32 {
	return (*Queries)(ctx.req).MustInt32(key, defaults...)
}

// QueryInt64 returns request form as int64 with default
func (ctx *Context) QueryInt64(key string, defaults ...int64) int64 {
	return (*Queries)(ctx.req).MustInt64(key, defaults...)
}

// QueryUint returns request form as uint with default
func (ctx *Context) QueryUint(key string, defaults ...uint) uint {
	return (*Queries)(ctx.req).MustUint(key, defaults...)
}

// QueryUint32 returns request form as uint32 with default
func (ctx *Context) QueryUint32(key string, defaults ...uint32) uint32 {
	return (*Queries)(ctx.req).MustUint32(key, defaults...)
}

// QueryUint64 returns request form as uint64 with default
func (ctx *Context) QueryUint64(key string, defaults ...uint64) uint64 {
	return (*Queries)(ctx.req).MustUint64(key, defaults...)
}

// QueryFloat32 returns request form as float32 with default
func (ctx *Context) QueryFloat32(key string, defaults ...float32) float32 {
	return (*Queries)(ctx.req).MustFloat32(key, defaults...)
}

// FormFloat64 returns request form as float64 with default
func (ctx *Context) QueryFloat64(key string, defaults ...float64) float64 {
	return (*Queries)(ctx.req).MustFloat64(key, defaults...)
}

// FormBool returns request form as bool with default
func (ctx *Context) QueryBool(key string, defaults ...bool) bool {
	return (*Queries)(ctx.req).MustBool(key, defaults...)
}
