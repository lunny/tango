// Copyright 2015 The Tango Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package tango

import (
	"errors"
	"html/template"
	"strconv"
)

type (
	param struct {
		Name  string
		Value string
	}
	// Params defines params of http request
	Params []param
)

var _ Set = &Params{}

// Get returns request form as string
func (p *Params) Get(key string) string {
	if len(key) == 0 {
		return ""
	}
	if key[0] != ':' && key[0] != '*' {
		key = ":" + key
	}

	for _, v := range *p {
		if v.Name == key {
			return v.Value
		}
	}
	return ""
}

// String returns request form as string
func (p *Params) String(key string) (string, error) {
	if len(key) == 0 {
		return "", errors.New("not exist")
	}
	if key[0] != ':' && key[0] != '*' {
		key = ":" + key
	}

	for _, v := range *p {
		if v.Name == key {
			return v.Value, nil
		}
	}
	return "", errors.New("not exist")
}

// Strings returns request form as slice of string
func (p *Params) Strings(key string) ([]string, error) {
	if len(key) == 0 {
		return nil, errors.New("not exist")
	}
	if key[0] != ':' && key[0] != '*' {
		key = ":" + key
	}

	var s = make([]string, 0)
	for _, v := range *p {
		if v.Name == key {
			s = append(s, v.Value)
		}
	}
	if len(s) > 0 {
		return s, nil
	}
	return nil, errors.New("not exist")
}

// Escape returns request form as escaped string
func (p *Params) Escape(key string) (string, error) {
	if len(key) == 0 {
		return "", errors.New("not exist")
	}
	if key[0] != ':' && key[0] != '*' {
		key = ":" + key
	}

	for _, v := range *p {
		if v.Name == key {
			return template.HTMLEscapeString(v.Value), nil
		}
	}
	return "", errors.New("not exist")
}

// Int returns request form as int
func (p *Params) Int(key string) (int, error) {
	return strconv.Atoi(p.Get(key))
}

// Int32 returns request form as int32
func (p *Params) Int32(key string) (int32, error) {
	v, err := strconv.ParseInt(p.Get(key), 10, 32)
	return int32(v), err
}

// Int64 returns request form as int64
func (p *Params) Int64(key string) (int64, error) {
	return strconv.ParseInt(p.Get(key), 10, 64)
}

// Uint returns request form as uint
func (p *Params) Uint(key string) (uint, error) {
	v, err := strconv.ParseUint(p.Get(key), 10, 64)
	return uint(v), err
}

// Uint32 returns request form as uint32
func (p *Params) Uint32(key string) (uint32, error) {
	v, err := strconv.ParseUint(p.Get(key), 10, 32)
	return uint32(v), err
}

// Uint64 returns request form as uint64
func (p *Params) Uint64(key string) (uint64, error) {
	return strconv.ParseUint(p.Get(key), 10, 64)
}

// Bool returns request form as bool
func (p *Params) Bool(key string) (bool, error) {
	return strconv.ParseBool(p.Get(key))
}

// Float32 returns request form as float32
func (p *Params) Float32(key string) (float32, error) {
	v, err := strconv.ParseFloat(p.Get(key), 32)
	return float32(v), err
}

// Float64 returns request form as float64
func (p *Params) Float64(key string) (float64, error) {
	return strconv.ParseFloat(p.Get(key), 64)
}

// MustString returns request form as slice of string with default
func (p *Params) MustString(key string, defaults ...string) string {
	if len(key) == 0 {
		return ""
	}
	if key[0] != ':' && key[0] != '*' {
		key = ":" + key
	}

	for _, v := range *p {
		if v.Name == key {
			return v.Value
		}
	}
	if len(defaults) > 0 {
		return defaults[0]
	}
	return ""
}

// MustStrings returns request form as slice of string with default
func (p *Params) MustStrings(key string, defaults ...[]string) []string {
	if len(key) == 0 {
		return []string{}
	}
	if key[0] != ':' && key[0] != '*' {
		key = ":" + key
	}

	var s = make([]string, 0)
	for _, v := range *p {
		if v.Name == key {
			s = append(s, v.Value)
		}
	}
	if len(s) > 0 {
		return s
	}
	if len(defaults) > 0 {
		return defaults[0]
	}
	return []string{}
}

// MustEscape returns request form as escaped string with default
func (p *Params) MustEscape(key string, defaults ...string) string {
	if len(key) == 0 {
		return ""
	}
	if key[0] != ':' && key[0] != '*' {
		key = ":" + key
	}

	for _, v := range *p {
		if v.Name == key {
			return template.HTMLEscapeString(v.Value)
		}
	}
	if len(defaults) > 0 {
		return defaults[0]
	}
	return ""
}

// MustInt returns request form as int with default
func (p *Params) MustInt(key string, defaults ...int) int {
	v, err := strconv.Atoi(p.Get(key))
	if len(defaults) > 0 && err != nil {
		return defaults[0]
	}
	return v
}

// MustInt32 returns request form as int32 with default
func (p *Params) MustInt32(key string, defaults ...int32) int32 {
	r, err := strconv.ParseInt(p.Get(key), 10, 32)
	if len(defaults) > 0 && err != nil {
		return defaults[0]
	}

	return int32(r)
}

// MustInt64 returns request form as int64 with default
func (p *Params) MustInt64(key string, defaults ...int64) int64 {
	r, err := strconv.ParseInt(p.Get(key), 10, 64)
	if len(defaults) > 0 && err != nil {
		return defaults[0]
	}
	return r
}

// MustUint returns request form as uint with default
func (p *Params) MustUint(key string, defaults ...uint) uint {
	v, err := strconv.ParseUint(p.Get(key), 10, 64)
	if len(defaults) > 0 && err != nil {
		return defaults[0]
	}
	return uint(v)
}

// MustUint32 returns request form as uint32 with default
func (p *Params) MustUint32(key string, defaults ...uint32) uint32 {
	r, err := strconv.ParseUint(p.Get(key), 10, 32)
	if len(defaults) > 0 && err != nil {
		return defaults[0]
	}

	return uint32(r)
}

// MustUint64 returns request form as uint64 with default
func (p *Params) MustUint64(key string, defaults ...uint64) uint64 {
	r, err := strconv.ParseUint(p.Get(key), 10, 64)
	if len(defaults) > 0 && err != nil {
		return defaults[0]
	}
	return r
}

// MustFloat32 returns request form as float32 with default
func (p *Params) MustFloat32(key string, defaults ...float32) float32 {
	r, err := strconv.ParseFloat(p.Get(key), 32)
	if len(defaults) > 0 && err != nil {
		return defaults[0]
	}
	return float32(r)
}

// MustFloat64 returns request form as float64 with default
func (p *Params) MustFloat64(key string, defaults ...float64) float64 {
	r, err := strconv.ParseFloat(p.Get(key), 64)
	if len(defaults) > 0 && err != nil {
		return defaults[0]
	}
	return r
}

// MustBool returns request form as bool with default
func (p *Params) MustBool(key string, defaults ...bool) bool {
	r, err := strconv.ParseBool(p.Get(key))
	if len(defaults) > 0 && err != nil {
		return defaults[0]
	}
	return r
}

// Param returns request form as string with default
func (ctx *Context) Param(key string, defaults ...string) string {
	return ctx.Params().MustString(key, defaults...)
}

// ParamStrings returns request form as slice of string with default
func (ctx *Context) ParamStrings(key string, defaults ...[]string) []string {
	return ctx.Params().MustStrings(key, defaults...)
}

// ParamEscape returns request form as escaped string with default
func (ctx *Context) ParamEscape(key string, defaults ...string) string {
	return ctx.Params().MustEscape(key, defaults...)
}

// ParamInt returns request form as int with default
func (ctx *Context) ParamInt(key string, defaults ...int) int {
	return ctx.Params().MustInt(key, defaults...)
}

// ParamInt32 returns request form as int32 with default
func (ctx *Context) ParamInt32(key string, defaults ...int32) int32 {
	return ctx.Params().MustInt32(key, defaults...)
}

// ParamInt64 returns request form as int64 with default
func (ctx *Context) ParamInt64(key string, defaults ...int64) int64 {
	return ctx.Params().MustInt64(key, defaults...)
}

// ParamUint returns request form as uint with default
func (ctx *Context) ParamUint(key string, defaults ...uint) uint {
	return ctx.Params().MustUint(key, defaults...)
}

// ParamUint32 returns request form as uint32 with default
func (ctx *Context) ParamUint32(key string, defaults ...uint32) uint32 {
	return ctx.Params().MustUint32(key, defaults...)
}

// ParamUint64 returns request form as uint64 with default
func (ctx *Context) ParamUint64(key string, defaults ...uint64) uint64 {
	return ctx.Params().MustUint64(key, defaults...)
}

// ParamFloat32 returns request form as float32 with default
func (ctx *Context) ParamFloat32(key string, defaults ...float32) float32 {
	return ctx.Params().MustFloat32(key, defaults...)
}

// ParamFloat64 returns request form as float64 with default
func (ctx *Context) ParamFloat64(key string, defaults ...float64) float64 {
	return ctx.Params().MustFloat64(key, defaults...)
}

// ParamBool returns request form as bool with default
func (ctx *Context) ParamBool(key string, defaults ...bool) bool {
	return ctx.Params().MustBool(key, defaults...)
}

// Set sets key/value to params
func (p *Params) Set(key, value string) {
	if len(key) == 0 {
		return
	}
	if key[0] != ':' && key[0] != '*' {
		key = ":" + key
	}

	for i, v := range *p {
		if v.Name == key {
			(*p)[i].Value = value
			return
		}
	}

	*p = append(*p, param{key, value})
}

// Paramer defines an interface to get params
type Paramer interface {
	SetParams([]param)
}

// SetParams implemented Paramer
func (p *Params) SetParams(params []param) {
	*p = params
}

// Param returns params handle to operate param
func Param() HandlerFunc {
	return func(ctx *Context) {
		if action := ctx.Action(); action != nil {
			if p, ok := action.(Paramer); ok {
				p.SetParams(*ctx.Params())
			}
		}
		ctx.Next()
	}
}
