// Copyright 2015 The Tango Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package tango

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func isValidCookieValue(p []byte) bool {
	for _, b := range p {
		if b <= ' ' ||
			b >= 127 ||
			b == '"' ||
			b == ',' ||
			b == ';' ||
			b == '\\' {
			return false
		}
	}
	return true
}

func isValidCookieName(s string) bool {
	for _, r := range s {
		if r <= ' ' ||
			r >= 127 ||
			strings.ContainsRune(" \t\"(),/:;<=>?@[]\\{}", r) {
			return false
		}
	}
	return true
}

// Set describes a set interface
type Set interface {
	String(key string) (string, error)
	Int(key string) (int, error)
	Int32(key string) (int32, error)
	Int64(key string) (int64, error)
	Uint(key string) (uint, error)
	Uint32(key string) (uint32, error)
	Uint64(key string) (uint64, error)
	Float32(key string) (float32, error)
	Float64(key string) (float64, error)
	Bool(key string) (bool, error)

	MustString(key string, defaults ...string) string
	MustEscape(key string, defaults ...string) string
	MustInt(key string, defaults ...int) int
	MustInt32(key string, defaults ...int32) int32
	MustInt64(key string, defaults ...int64) int64
	MustUint(key string, defaults ...uint) uint
	MustUint32(key string, defaults ...uint32) uint32
	MustUint64(key string, defaults ...uint64) uint64
	MustFloat32(key string, defaults ...float32) float32
	MustFloat64(key string, defaults ...float64) float64
	MustBool(key string, defaults ...bool) bool
}

// Cookies describes cookie interface
type Cookies interface {
	Set
	Get(string) *http.Cookie
	Set(*http.Cookie)
	Expire(string, time.Time)
	Del(string)
}

type cookies Context

var _ Cookies = &cookies{}

// NewCookie return a http.Cookie via give name and value
func NewCookie(name string, value string, age ...int64) *http.Cookie {
	if !isValidCookieName(name) || !isValidCookieValue([]byte(value)) {
		return nil
	}

	var utctime time.Time
	if len(age) == 0 {
		// 2^31 - 1 seconds (roughly 2038)
		utctime = time.Unix(2147483647, 0)
	} else {
		utctime = time.Unix(time.Now().Unix()+age[0], 0)
	}
	return &http.Cookie{Name: name, Value: value, Expires: utctime}
}

// Get return http.Cookie via key
func (c *cookies) Get(key string) *http.Cookie {
	ck, err := c.req.Cookie(key)
	if err != nil {
		return nil
	}
	return ck
}

// Set set a http.Cookie
func (c *cookies) Set(ck *http.Cookie) {
	http.SetCookie(c.ResponseWriter, ck)
}

// Expire let a cookie named key when expire address
func (c *cookies) Expire(key string, expire time.Time) {
	ck := c.Get(key)
	if ck != nil {
		ck.Expires = expire
		ck.MaxAge = int(expire.Sub(time.Now()).Seconds())
		c.Set(ck)
	}
}

// Del del cookie by key
func (c *cookies) Del(key string) {
	c.Expire(key, time.Date(1900, 1, 1, 0, 0, 0, 0, time.Local))
}

// String get cookie as string
func (c *cookies) String(key string) (string, error) {
	ck, err := c.req.Cookie(key)
	if err != nil {
		return "", err
	}
	return ck.Value, nil
}

// Int get cookie as int
func (c *cookies) Int(key string) (int, error) {
	ck, err := c.req.Cookie(key)
	if err != nil {
		return 0, err
	}
	return strconv.Atoi(ck.Value)
}

// Int32 get cookie as int32
func (c *cookies) Int32(key string) (int32, error) {
	ck, err := c.req.Cookie(key)
	if err != nil {
		return 0, err
	}
	v, err := strconv.ParseInt(ck.Value, 10, 32)
	return int32(v), err
}

// Int64 get cookie as int64
func (c *cookies) Int64(key string) (int64, error) {
	ck, err := c.req.Cookie(key)
	if err != nil {
		return 0, err
	}
	return strconv.ParseInt(ck.Value, 10, 64)
}

// Uint get cookie as uint
func (c *cookies) Uint(key string) (uint, error) {
	ck, err := c.req.Cookie(key)
	if err != nil {
		return 0, err
	}
	v, err := strconv.ParseUint(ck.Value, 10, 64)
	return uint(v), err
}

// Uint32 get cookie as uint32
func (c *cookies) Uint32(key string) (uint32, error) {
	ck, err := c.req.Cookie(key)
	if err != nil {
		return 0, err
	}
	v, err := strconv.ParseUint(ck.Value, 10, 32)
	return uint32(v), err
}

// Uint64 get cookie as uint64
func (c *cookies) Uint64(key string) (uint64, error) {
	ck, err := c.req.Cookie(key)
	if err != nil {
		return 0, err
	}
	return strconv.ParseUint(ck.Value, 10, 64)
}

// Float32 get cookie as float32
func (c *cookies) Float32(key string) (float32, error) {
	ck, err := c.req.Cookie(key)
	if err != nil {
		return 0, err
	}
	v, err := strconv.ParseFloat(ck.Value, 32)
	return float32(v), err
}

// Float64 get cookie as float64
func (c *cookies) Float64(key string) (float64, error) {
	ck, err := c.req.Cookie(key)
	if err != nil {
		return 0, err
	}
	return strconv.ParseFloat(ck.Value, 32)
}

// Bool get cookie as bool
func (c *cookies) Bool(key string) (bool, error) {
	ck, err := c.req.Cookie(key)
	if err != nil {
		return false, err
	}
	return strconv.ParseBool(ck.Value)
}

// MustString get cookie as string with default
func (c *cookies) MustString(key string, defaults ...string) string {
	ck, err := c.req.Cookie(key)
	if err != nil {
		if len(defaults) > 0 {
			return defaults[0]
		}
		return ""
	}
	return ck.Value
}

// MustEscape get cookie as escaped string with default
func (c *cookies) MustEscape(key string, defaults ...string) string {
	ck, err := c.req.Cookie(key)
	if err != nil {
		if len(defaults) > 0 {
			return defaults[0]
		}
		return ""
	}
	return template.HTMLEscapeString(ck.Value)
}

// MustInt get cookie as int with default
func (c *cookies) MustInt(key string, defaults ...int) int {
	ck, err := c.req.Cookie(key)
	if err != nil {
		if len(defaults) > 0 {
			return defaults[0]
		}
		return 0
	}
	v, err := strconv.Atoi(ck.Value)
	if len(defaults) > 0 && err != nil {
		return defaults[0]
	}
	return v
}

// MustInt32 get cookie as int32 with default
func (c *cookies) MustInt32(key string, defaults ...int32) int32 {
	ck, err := c.req.Cookie(key)
	if err != nil {
		if len(defaults) > 0 {
			return defaults[0]
		}
		return 0
	}
	v, err := strconv.ParseInt(ck.Value, 10, 32)
	if len(defaults) > 0 && err != nil {
		return defaults[0]
	}
	return int32(v)
}

// MustInt64 get cookie as int64 with default
func (c *cookies) MustInt64(key string, defaults ...int64) int64 {
	ck, err := c.req.Cookie(key)
	if err != nil {
		if len(defaults) > 0 {
			return defaults[0]
		}
		return 0
	}
	v, err := strconv.ParseInt(ck.Value, 10, 64)
	if len(defaults) > 0 && err != nil {
		return defaults[0]
	}
	return v
}

// MustUint get cookie as uint with default
func (c *cookies) MustUint(key string, defaults ...uint) uint {
	ck, err := c.req.Cookie(key)
	if err != nil {
		if len(defaults) > 0 {
			return defaults[0]
		}
		return 0
	}
	v, err := strconv.ParseUint(ck.Value, 10, 64)
	if len(defaults) > 0 && err != nil {
		return defaults[0]
	}
	return uint(v)
}

// MustUint32 get cookie as uint32 with default
func (c *cookies) MustUint32(key string, defaults ...uint32) uint32 {
	ck, err := c.req.Cookie(key)
	if err != nil {
		if len(defaults) > 0 {
			return defaults[0]
		}
		return 0
	}
	v, err := strconv.ParseUint(ck.Value, 10, 32)
	if len(defaults) > 0 && err != nil {
		return defaults[0]
	}
	return uint32(v)
}

// MustUint64 get cookie as uint64 with default
func (c *cookies) MustUint64(key string, defaults ...uint64) uint64 {
	ck, err := c.req.Cookie(key)
	if err != nil {
		if len(defaults) > 0 {
			return defaults[0]
		}
		return 0
	}
	v, err := strconv.ParseUint(ck.Value, 10, 64)
	if len(defaults) > 0 && err != nil {
		return defaults[0]
	}
	return v
}

// MustFloat32 get cookie as float32 with default
func (c *cookies) MustFloat32(key string, defaults ...float32) float32 {
	ck, err := c.req.Cookie(key)
	if err != nil {
		if len(defaults) > 0 {
			return defaults[0]
		}
		return 0
	}
	v, err := strconv.ParseFloat(ck.Value, 32)
	if len(defaults) > 0 && err != nil {
		return defaults[0]
	}
	return float32(v)
}

// MustFloat64 get cookie as float64 with default
func (c *cookies) MustFloat64(key string, defaults ...float64) float64 {
	ck, err := c.req.Cookie(key)
	if err != nil {
		if len(defaults) > 0 {
			return defaults[0]
		}
		return 0
	}
	v, err := strconv.ParseFloat(ck.Value, 32)
	if len(defaults) > 0 && err != nil {
		return defaults[0]
	}
	return v
}

// MustBool get cookie as bool with default
func (c *cookies) MustBool(key string, defaults ...bool) bool {
	ck, err := c.req.Cookie(key)
	if err != nil {
		if len(defaults) > 0 {
			return defaults[0]
		}
		return false
	}
	v, err := strconv.ParseBool(ck.Value)
	if len(defaults) > 0 && err != nil {
		return defaults[0]
	}
	return v
}

// Cookie returns cookie as string with default
func (ctx *Context) Cookie(key string, defaults ...string) string {
	return ctx.Cookies().MustString(key, defaults...)
}

// CookieEscape returns cookie as escaped string with default
func (ctx *Context) CookieEscape(key string, defaults ...string) string {
	return ctx.Cookies().MustEscape(key, defaults...)
}

// CookieInt returns cookie as int with default
func (ctx *Context) CookieInt(key string, defaults ...int) int {
	return ctx.Cookies().MustInt(key, defaults...)
}

// CookieInt32 returns cookie as int32 with default
func (ctx *Context) CookieInt32(key string, defaults ...int32) int32 {
	return ctx.Cookies().MustInt32(key, defaults...)
}

// CookieInt64 returns cookie as int64 with default
func (ctx *Context) CookieInt64(key string, defaults ...int64) int64 {
	return ctx.Cookies().MustInt64(key, defaults...)
}

// CookieUint returns cookie as uint with default
func (ctx *Context) CookieUint(key string, defaults ...uint) uint {
	return ctx.Cookies().MustUint(key, defaults...)
}

// CookieUint32 returns cookie as uint32 with default
func (ctx *Context) CookieUint32(key string, defaults ...uint32) uint32 {
	return ctx.Cookies().MustUint32(key, defaults...)
}

// CookieUint64 returns cookie as uint64 with default
func (ctx *Context) CookieUint64(key string, defaults ...uint64) uint64 {
	return ctx.Cookies().MustUint64(key, defaults...)
}

// CookieFloat32 returns cookie as float32 with default
func (ctx *Context) CookieFloat32(key string, defaults ...float32) float32 {
	return ctx.Cookies().MustFloat32(key, defaults...)
}

// CookieFloat64 returns cookie as float64 with default
func (ctx *Context) CookieFloat64(key string, defaults ...float64) float64 {
	return ctx.Cookies().MustFloat64(key, defaults...)
}

// CookieBool returns cookie as bool with default
func (ctx *Context) CookieBool(key string, defaults ...bool) bool {
	return ctx.Cookies().MustBool(key, defaults...)
}

func getCookieSig(key string, val []byte, timestamp string) string {
	hm := hmac.New(sha1.New, []byte(key))

	hm.Write(val)
	hm.Write([]byte(timestamp))

	hex := fmt.Sprintf("%02x", hm.Sum(nil))
	return hex
}

// secure cookies
type secureCookies struct {
	*cookies
	secret string
}

var _ Cookies = &secureCookies{}

func parseSecureCookie(secret string, value string) string {
	if strings.Count(value, "|") < 2 {
		return ""
	}
	parts := strings.SplitN(value, "|", 3)
	val, timestamp, sig := parts[0], parts[1], parts[2]
	if getCookieSig(secret, []byte(val), timestamp) != sig {
		return ""
	}

	ts, _ := strconv.ParseInt(timestamp, 0, 64)
	if time.Now().Unix()-31*86400 > ts {
		return ""
	}

	buf := bytes.NewBufferString(val)
	encoder := base64.NewDecoder(base64.StdEncoding, buf)

	res, _ := ioutil.ReadAll(encoder)
	return string(res)
}

func (c *secureCookies) Get(key string) *http.Cookie {
	ck := c.cookies.Get(key)
	if ck == nil {
		return nil
	}

	v := parseSecureCookie(c.secret, ck.Value)
	if v == "" {
		return nil
	}
	ck.Value = v
	return ck
}

func (c *secureCookies) String(key string) (string, error) {
	ck, err := c.req.Cookie(key)
	if err != nil {
		return "", err
	}
	v := parseSecureCookie(c.secret, ck.Value)
	return v, nil
}

func (c *secureCookies) Int(key string) (int, error) {
	ck, err := c.req.Cookie(key)
	if err != nil {
		return 0, err
	}
	s := parseSecureCookie(c.secret, ck.Value)
	return strconv.Atoi(s)
}

// Int32 gets int32 data
func (c *secureCookies) Int32(key string) (int32, error) {
	ck, err := c.req.Cookie(key)
	if err != nil {
		return 0, err
	}
	s := parseSecureCookie(c.secret, ck.Value)
	v, err := strconv.ParseInt(s, 10, 32)
	return int32(v), err
}

// Int64 gets int64 data
func (c *secureCookies) Int64(key string) (int64, error) {
	ck, err := c.req.Cookie(key)
	if err != nil {
		return 0, err
	}
	s := parseSecureCookie(c.secret, ck.Value)
	return strconv.ParseInt(s, 10, 64)
}

// Uint gets uint data
func (c *secureCookies) Uint(key string) (uint, error) {
	ck, err := c.req.Cookie(key)
	if err != nil {
		return 0, err
	}
	s := parseSecureCookie(c.secret, ck.Value)
	v, err := strconv.ParseUint(s, 10, 64)
	return uint(v), err
}

// Uint32 gets uint32 data
func (c *secureCookies) Uint32(key string) (uint32, error) {
	ck, err := c.req.Cookie(key)
	if err != nil {
		return 0, err
	}
	s := parseSecureCookie(c.secret, ck.Value)
	v, err := strconv.ParseUint(s, 10, 32)
	return uint32(v), err
}

// Uint64 gets unit64 data
func (c *secureCookies) Uint64(key string) (uint64, error) {
	ck, err := c.req.Cookie(key)
	if err != nil {
		return 0, err
	}
	s := parseSecureCookie(c.secret, ck.Value)
	return strconv.ParseUint(s, 10, 64)
}

// Float32 gets float32 data
func (c *secureCookies) Float32(key string) (float32, error) {
	ck, err := c.req.Cookie(key)
	if err != nil {
		return 0, err
	}
	s := parseSecureCookie(c.secret, ck.Value)
	v, err := strconv.ParseFloat(s, 32)
	return float32(v), err
}

// Float64 gets float64 data
func (c *secureCookies) Float64(key string) (float64, error) {
	ck, err := c.req.Cookie(key)
	if err != nil {
		return 0, err
	}
	s := parseSecureCookie(c.secret, ck.Value)
	return strconv.ParseFloat(s, 32)
}

// Bool gets bool data
func (c *secureCookies) Bool(key string) (bool, error) {
	ck, err := c.req.Cookie(key)
	if err != nil {
		return false, err
	}
	s := parseSecureCookie(c.secret, ck.Value)
	return strconv.ParseBool(s)
}

// MustString gets data
func (c *secureCookies) MustString(key string, defaults ...string) string {
	ck, err := c.req.Cookie(key)
	if err != nil {
		if len(defaults) > 0 {
			return defaults[0]
		}
		return ""
	}
	s := parseSecureCookie(c.secret, ck.Value)
	return s
}

// MustEscape gets data escaped
func (c *secureCookies) MustEscape(key string, defaults ...string) string {
	ck, err := c.req.Cookie(key)
	if err != nil {
		if len(defaults) > 0 {
			return defaults[0]
		}
		return ""
	}
	s := parseSecureCookie(c.secret, ck.Value)
	return template.HTMLEscapeString(s)
}

// MustInt gets data int
func (c *secureCookies) MustInt(key string, defaults ...int) int {
	ck, err := c.req.Cookie(key)
	if err != nil {
		if len(defaults) > 0 {
			return defaults[0]
		}
		return 0
	}
	s := parseSecureCookie(c.secret, ck.Value)
	v, err := strconv.Atoi(s)
	if len(defaults) > 0 && err != nil {
		return defaults[0]
	}
	return v
}

// MustInt32 gets data int32
func (c *secureCookies) MustInt32(key string, defaults ...int32) int32 {
	ck, err := c.req.Cookie(key)
	if err != nil {
		if len(defaults) > 0 {
			return defaults[0]
		}
		return 0
	}
	s := parseSecureCookie(c.secret, ck.Value)
	v, err := strconv.ParseInt(s, 10, 32)
	if len(defaults) > 0 && err != nil {
		return defaults[0]
	}
	return int32(v)
}

// MustInt64 gets data int64 type
func (c *secureCookies) MustInt64(key string, defaults ...int64) int64 {
	ck, err := c.req.Cookie(key)
	if err != nil {
		if len(defaults) > 0 {
			return defaults[0]
		}
		return 0
	}
	s := parseSecureCookie(c.secret, ck.Value)
	v, err := strconv.ParseInt(s, 10, 64)
	if len(defaults) > 0 && err != nil {
		return defaults[0]
	}
	return v
}

// MustUint gets data unit type
func (c *secureCookies) MustUint(key string, defaults ...uint) uint {
	ck, err := c.req.Cookie(key)
	if err != nil {
		if len(defaults) > 0 {
			return defaults[0]
		}
		return 0
	}
	s := parseSecureCookie(c.secret, ck.Value)
	v, err := strconv.ParseUint(s, 10, 64)
	if len(defaults) > 0 && err != nil {
		return defaults[0]
	}
	return uint(v)
}

// MustUint32 gets data uint32 type
func (c *secureCookies) MustUint32(key string, defaults ...uint32) uint32 {
	ck, err := c.req.Cookie(key)
	if err != nil {
		if len(defaults) > 0 {
			return defaults[0]
		}
		return 0
	}
	s := parseSecureCookie(c.secret, ck.Value)
	v, err := strconv.ParseUint(s, 10, 32)
	if len(defaults) > 0 && err != nil {
		return defaults[0]
	}
	return uint32(v)
}

// MustUint64 gets data unit64 type
func (c *secureCookies) MustUint64(key string, defaults ...uint64) uint64 {
	ck, err := c.req.Cookie(key)
	if err != nil {
		if len(defaults) > 0 {
			return defaults[0]
		}
		return 0
	}
	s := parseSecureCookie(c.secret, ck.Value)
	v, err := strconv.ParseUint(s, 10, 64)
	if len(defaults) > 0 && err != nil {
		return defaults[0]
	}
	return v
}

// MustFloat32 gets data float32 type
func (c *secureCookies) MustFloat32(key string, defaults ...float32) float32 {
	ck, err := c.req.Cookie(key)
	if err != nil {
		if len(defaults) > 0 {
			return defaults[0]
		}
		return 0
	}
	s := parseSecureCookie(c.secret, ck.Value)
	v, err := strconv.ParseFloat(s, 32)
	if len(defaults) > 0 && err != nil {
		return defaults[0]
	}
	return float32(v)
}

// MustFloat64 gets data float64 type
func (c *secureCookies) MustFloat64(key string, defaults ...float64) float64 {
	ck, err := c.req.Cookie(key)
	if err != nil {
		if len(defaults) > 0 {
			return defaults[0]
		}
		return 0
	}
	s := parseSecureCookie(c.secret, ck.Value)
	v, err := strconv.ParseFloat(s, 32)
	if len(defaults) > 0 && err != nil {
		return defaults[0]
	}
	return v
}

// MustBool gets data bool type
func (c *secureCookies) MustBool(key string, defaults ...bool) bool {
	ck, err := c.req.Cookie(key)
	if err != nil {
		if len(defaults) > 0 {
			return defaults[0]
		}
		return false
	}
	s := parseSecureCookie(c.secret, ck.Value)
	v, err := strconv.ParseBool(s)
	if len(defaults) > 0 && err != nil {
		return defaults[0]
	}
	return v
}

func secCookieValue(secret string, vb []byte) string {
	timestamp := strconv.FormatInt(time.Now().Unix(), 10)
	sig := getCookieSig(secret, vb, timestamp)
	return strings.Join([]string{string(vb), timestamp, sig}, "|")
}

// NewSecureCookie generates a new secure cookie
func NewSecureCookie(secret, name string, val string, age ...int64) *http.Cookie {
	var buf bytes.Buffer
	encoder := base64.NewEncoder(base64.StdEncoding, &buf)
	encoder.Write([]byte(val))
	encoder.Close()

	cookie := secCookieValue(secret, buf.Bytes())
	return NewCookie(name, cookie, age...)
}

// Expire sets key expire time
func (c *secureCookies) Expire(key string, expire time.Time) {
	ck := c.Get(key)
	if ck != nil {
		ck.Expires = expire
		ck.Value = secCookieValue(c.secret, []byte(ck.Value))
		c.Set(ck)
	}
}

// Del deletes key
func (c *secureCookies) Del(key string) {
	c.Expire(key, time.Date(1900, 1, 1, 0, 0, 0, 0, time.Local))
}
