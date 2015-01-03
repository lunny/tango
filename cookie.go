package tango

import (
	"net/http"
	"time"
	"fmt"
	"strings"
	"crypto/sha1"
	"crypto/hmac"
	"strconv"
	"bytes"
	"encoding/base64"
	"io/ioutil"
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

type Cookies interface {
	Get(string) *http.Cookie
	Set(*http.Cookie)
	Expire(string, time.Time)
	Del(string)
}

type cookies struct {
	req *http.Request
	resp http.ResponseWriter
}

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

func (c *cookies) Get(key string) *http.Cookie {
	ck, err := c.req.Cookie(key)
	if err != nil {
		return nil
	}
	return ck
}

func (c *cookies) Set(ck *http.Cookie) {
 	http.SetCookie(c.resp, ck)
}

func (c *cookies) Expire(key string, expire time.Time) {
	ck := c.Get(key)
	if ck != nil {
		ck.Expires = expire
		c.Set(ck)
	}
}

func (c *cookies) Del(key string) {
	c.Expire(key, time.Date(1900, 1, 1, 0, 0, 0, 0, time.Local))
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

func parseSecureCookie(secret string, value string) string {
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

func NewSecureCookie(secret, name string, val string, age ...int64) *http.Cookie {
	var buf bytes.Buffer
	encoder := base64.NewEncoder(base64.StdEncoding, &buf)
	encoder.Write([]byte(val))
	encoder.Close()

	vs := buf.String()
	vb := buf.Bytes()
	timestamp := strconv.FormatInt(time.Now().Unix(), 10)
	sig := getCookieSig(secret, vb, timestamp)
	cookie := strings.Join([]string{vs, timestamp, sig}, "|")
	return NewCookie(name, cookie, age...)
}

func (c *secureCookies) Expire(key string, expire time.Time) {
	ck := c.Get(key)
	if ck != nil {
		ck.Expires = expire
		c.Set(ck)
	}
}

func (c *secureCookies) Del(key string) {
	c.Expire(key, time.Date(1900, 1, 1, 0, 0, 0, 0, time.Local))
}