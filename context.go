// Copyright 2015 The Tango Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package tango

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"reflect"
	"strings"
)

// Handler defines middleware interface
type Handler interface {
	Handle(*Context)
}

// Context defines request and response context
type Context struct {
	tan *Tango
	Logger

	idx int
	req *http.Request
	ResponseWriter
	route    *Route
	params   Params
	callArgs []reflect.Value
	matched  bool
	stage    byte

	action interface{}
	Result interface{}
}

func (ctx *Context) reset(req *http.Request, resp ResponseWriter) {
	ctx.req = req
	ctx.ResponseWriter = resp
	ctx.idx = 0
	ctx.stage = 0
	ctx.route = nil
	ctx.params = nil
	ctx.callArgs = nil
	ctx.matched = false
	ctx.action = nil
	ctx.Result = nil
}

// HandleError handles errors
func (ctx *Context) HandleError() {
	ctx.tan.ErrHandler.Handle(ctx)
}

// Req returns current HTTP Request information
func (ctx *Context) Req() *http.Request {
	return ctx.req
}

// IsAjax returns if the request is an ajax request
func (ctx *Context) IsAjax() bool {
	return ctx.Req().Header.Get("X-Requested-With") == "XMLHttpRequest"
}

// SecureCookies generates a secret cookie
func (ctx *Context) SecureCookies(secret string) Cookies {
	return &secureCookies{
		(*cookies)(ctx),
		secret,
	}
}

// Cookies returns the cookies
func (ctx *Context) Cookies() Cookies {
	return (*cookies)(ctx)
}

// Forms returns the query names and values
func (ctx *Context) Forms() *Forms {
	ctx.req.ParseForm()
	return (*Forms)(ctx.req)
}

// Route returns route
func (ctx *Context) Route() *Route {
	ctx.newAction()
	return ctx.route
}

// Params returns the URL params
func (ctx *Context) Params() *Params {
	ctx.newAction()
	return &ctx.params
}

// IP returns remote IP
func (ctx *Context) IP() string {
	proxy := []string{}
	if ips := ctx.Req().Header.Get("X-Forwarded-For"); ips != "" {
		proxy = strings.Split(ips, ",")
	}
	if len(proxy) > 0 && proxy[0] != "" {
		return proxy[0]
	}
	ip := strings.Split(ctx.Req().RemoteAddr, ":")
	if len(ip) > 0 {
		if ip[0] != "[" {
			return ip[0]
		}
	}
	return "127.0.0.1"
}

// Action returns action
func (ctx *Context) Action() interface{} {
	ctx.newAction()
	return ctx.action
}

// ActionValue returns action value
func (ctx *Context) ActionValue() reflect.Value {
	ctx.newAction()
	return ctx.callArgs[0]
}

// ActionTag returns field tag on action struct
// TODO: cache the name
func (ctx *Context) ActionTag(fieldName string) string {
	ctx.newAction()
	if ctx.route.routeType == StructPtrRoute || ctx.route.routeType == StructRoute {
		tp := ctx.callArgs[0].Type()
		if tp.Kind() == reflect.Ptr {
			tp = tp.Elem()
		}
		field, ok := tp.FieldByName(fieldName)
		if !ok {
			return ""
		}
		return string(field.Tag)
	}
	return ""
}

// WriteString writes a string to response write
func (ctx *Context) WriteString(content string) (int, error) {
	return io.WriteString(ctx.ResponseWriter, content)
}

func (ctx *Context) newAction() {
	if !ctx.matched {
		reqPath := removeStick(ctx.Req().URL.Path)
		ctx.route, ctx.params = ctx.tan.Match(reqPath, ctx.Req().Method)
		if ctx.route != nil {
			vc := ctx.route.newAction()
			ctx.action = vc.Interface()
			switch ctx.route.routeType {
			case StructPtrRoute:
				ctx.callArgs = []reflect.Value{vc.Elem()}
			case StructRoute:
				ctx.callArgs = []reflect.Value{vc}
			case FuncRoute:
				ctx.callArgs = []reflect.Value{}
			case FuncHTTPRoute:
				ctx.callArgs = []reflect.Value{reflect.ValueOf(ctx.ResponseWriter),
					reflect.ValueOf(ctx.Req())}
			case FuncReqRoute:
				ctx.callArgs = []reflect.Value{reflect.ValueOf(ctx.Req())}
			case FuncResponseRoute:
				ctx.callArgs = []reflect.Value{reflect.ValueOf(ctx.ResponseWriter)}
			case FuncCtxRoute:
				ctx.callArgs = []reflect.Value{reflect.ValueOf(ctx)}
			default:
				panic("routeType error")
			}
		}
		ctx.matched = true
	}
}

// Next call next middleware or action
// WARNING: don't invoke this method on action
func (ctx *Context) Next() {
	ctx.idx++
	ctx.invoke()
}

func (ctx *Context) execute() {
	ctx.newAction()
	// route is matched
	if ctx.action != nil {
		if len(ctx.route.handlers) > 0 && ctx.stage == 0 {
			ctx.idx = 0
			ctx.stage = 1
			ctx.invoke()
			return
		}

		var ret []reflect.Value
		switch fn := ctx.route.raw.(type) {
		case func(*Context):
			fn(ctx)
		case func(*http.Request, http.ResponseWriter):
			fn(ctx.req, ctx.ResponseWriter)
		case func():
			fn()
		case func(*http.Request):
			fn(ctx.req)
		case func(http.ResponseWriter):
			fn(ctx.ResponseWriter)
		default:
			ret = ctx.route.method.Call(ctx.callArgs)
		}

		if len(ret) == 1 {
			ctx.Result = ret[0].Interface()
		} else if len(ret) == 2 {
			if code, ok := ret[0].Interface().(int); ok {
				ctx.Result = &StatusResult{code, ret[1].Interface()}
			}
		}
		// not route matched
	} else {
		if !ctx.Written() {
			ctx.NotFound()
		}
	}
}

func (ctx *Context) invoke() {
	if ctx.stage == 0 {
		if ctx.idx < len(ctx.tan.handlers) {
			ctx.tan.handlers[ctx.idx].Handle(ctx)
		} else {
			ctx.execute()
		}
	} else if ctx.stage == 1 {
		if ctx.idx < len(ctx.route.handlers) {
			ctx.route.handlers[ctx.idx].Handle(ctx)
		} else {
			ctx.execute()
		}
	}
}

func toHTTPError(err error) (msg string, httpStatus int) {
	if os.IsNotExist(err) {
		return "404 page not found", http.StatusNotFound
	}
	if os.IsPermission(err) {
		return "403 Forbidden", http.StatusForbidden
	}
	// Default:
	return "500 Internal Server Error", http.StatusInternalServerError
}

// ServeFile serves a file
func (ctx *Context) ServeFile(path string) error {
	f, err := os.Open(path)
	if err != nil {
		msg, code := toHTTPError(err)
		http.Error(ctx, msg, code)
		return nil
	}
	defer f.Close()

	d, err := f.Stat()
	if err != nil {
		msg, code := toHTTPError(err)
		http.Error(ctx, msg, code)
		return nil
	}

	if d.IsDir() {
		http.Error(ctx, http.StatusText(http.StatusForbidden), http.StatusForbidden)
		return nil
	}

	http.ServeContent(ctx, ctx.Req(), d.Name(), d.ModTime(), f)
	return nil
}

// ServeXml serves marshaled XML content from obj
// Deprecated: use ServeXML instead
func (ctx *Context) ServeXml(obj interface{}) error {
	return ctx.ServeXML(obj)
}

// ServeXML serves marshaled XML content from obj
func (ctx *Context) ServeXML(obj interface{}) error {
	encoder := xml.NewEncoder(ctx)
	ctx.Header().Set("Content-Type", "application/xml; charset=UTF-8")
	err := encoder.Encode(obj)
	if err != nil {
		ctx.Header().Del("Content-Type")
	}
	return err
}

// ServeJson serves marshaled JSON content from obj
// Deprecated: use ServeJSON instead
func (ctx *Context) ServeJson(obj interface{}) error {
	return ctx.ServeJSON(obj)
}

// ServeJSON serves marshaled JSON content from obj
func (ctx *Context) ServeJSON(obj interface{}) error {
	encoder := json.NewEncoder(ctx)
	ctx.Header().Set("Content-Type", "application/json; charset=UTF-8")
	err := encoder.Encode(obj)
	if err != nil {
		ctx.Header().Del("Content-Type")
	}
	return err
}

// Body returns body's content
func (ctx *Context) Body() ([]byte, error) {
	body, err := ioutil.ReadAll(ctx.req.Body)
	if err != nil {
		return nil, err
	}

	ctx.req.Body.Close()
	ctx.req.Body = ioutil.NopCloser(bytes.NewBuffer(body))

	return body, nil
}

// DecodeJson decodes body as JSON format to obj
// Deprecated: use DecodeJSON instead
func (ctx *Context) DecodeJson(obj interface{}) error {
	return ctx.DecodeJSON(obj)
}

// DecodeJSON decodes body as JSON format to obj
func (ctx *Context) DecodeJSON(obj interface{}) error {
	body, err := ctx.Body()
	if err != nil {
		return err
	}

	return json.Unmarshal(body, obj)
}

// DecodeXml decodes body as XML format to obj
// Deprecated: use DecodeXML instead
func (ctx *Context) DecodeXml(obj interface{}) error {
	return ctx.DecodeXML(obj)
}

// DecodeXML decodes body as XML format to obj
func (ctx *Context) DecodeXML(obj interface{}) error {
	body, err := ctx.Body()
	if err != nil {
		return err
	}

	return xml.Unmarshal(body, obj)
}

// Download provides a locale file to http client
func (ctx *Context) Download(fpath string) error {
	f, err := os.Open(fpath)
	if err != nil {
		return err
	}
	defer f.Close()

	fName := filepath.Base(fpath)
	ctx.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"%v\"", fName))
	_, err = io.Copy(ctx, f)
	return err
}

// SaveToFile saves the HTTP post file form to local file path
func (ctx *Context) SaveToFile(formName, savePath string) error {
	file, _, err := ctx.Req().FormFile(formName)
	if err != nil {
		return err
	}
	defer file.Close()

	f, err := os.OpenFile(savePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = io.Copy(f, file)
	return err
}

// Redirect redirects the request to another URL
func (ctx *Context) Redirect(url string, status ...int) {
	s := http.StatusFound
	if len(status) > 0 {
		s = status[0]
	}
	http.Redirect(ctx.ResponseWriter, ctx.Req(), url, s)
}

// NotModified writes a 304 HTTP response
func (ctx *Context) NotModified() {
	ctx.WriteHeader(http.StatusNotModified)
}

// Unauthorized writes a 401 HTTP response
func (ctx *Context) Unauthorized() {
	ctx.Abort(http.StatusUnauthorized, http.StatusText(http.StatusUnauthorized))
}

// NotFound writes a 404 HTTP response
func (ctx *Context) NotFound(message ...string) {
	if len(message) == 0 {
		ctx.Abort(http.StatusNotFound, http.StatusText(http.StatusNotFound))
		return
	}
	ctx.Abort(http.StatusNotFound, message[0])
}

// Abort is a helper method that sends an HTTP header and an optional
// body. It is useful for returning 4xx or 5xx errors.
// Once it has been called, any return value from the handler will
// not be written to the response.
func (ctx *Context) Abort(status int, body ...string) {
	ctx.Result = Abort(status, body...)
	ctx.HandleError()
}

// Contexter describes an interface to set *Context
type Contexter interface {
	SetContext(*Context)
}

// Ctx implements Contexter
type Ctx struct {
	*Context
}

// SetContext set *Context to action struct
func (c *Ctx) SetContext(ctx *Context) {
	c.Context = ctx
}

// Contexts returns a middleware to inject Context to action struct
func Contexts() HandlerFunc {
	return func(ctx *Context) {
		if action := ctx.Action(); action != nil {
			if a, ok := action.(Contexter); ok {
				a.SetContext(ctx)
			}
		}
		ctx.Next()
	}
}
