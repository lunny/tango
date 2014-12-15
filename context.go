package tango

import (
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"reflect"
)

type Handler interface {
	Handle(*Context)
}

type Context struct {
	router   Router
	handlers []Handler

	idx int
	req *http.Request
	ResponseWriter
	route   *Route
	args    []reflect.Value
	matched bool

	action interface{}
	Result interface{}
}

func NewContext(
	router Router,
	handlers []Handler,
	req *http.Request,
	resp ResponseWriter) *Context {
	return &Context{
		router:         router,
		handlers:       handlers,
		idx:            0,
		req:            req,
		ResponseWriter: resp,
	}
}

func (ctx *Context) Req() *http.Request {
	return ctx.req
}

func (ctx *Context) Route() *Route {
	ctx.newAction()
	return ctx.route
}

func (ctx *Context) Action() interface{} {
	ctx.newAction()
	return ctx.action
}

func (ctx *Context) newAction() {
	if !ctx.matched {
		reqPath := removeStick(ctx.Req().URL.Path)
		var allowedMethod = ctx.Req().Method
		if allowedMethod == "HEAD" {
			allowedMethod = "GET"
		}

		route, args := ctx.router.Match(reqPath, allowedMethod)
		if route != nil {
			ctx.route = route
			vc := route.newAction()
			ctx.action = vc.Interface()
			switch route.routeType {
			case StructPtrRoute:
				ctx.args = append([]reflect.Value{vc.Elem()}, args...)
			case StructRoute:
				ctx.args = append([]reflect.Value{vc}, args...)
			case FuncRoute:
				ctx.args = args
			}
		}
		ctx.matched = true
	}
}

func (ctx *Context) Next() {
	ctx.idx += 1
	ctx.Invoke()
}

func (ctx *Context) Invoke() {
	if ctx.idx < len(ctx.handlers) {
		ctx.handlers[ctx.idx].Handle(ctx)
	} else {
		ctx.newAction()
		if ctx.action != nil {
			ret := ctx.route.method.Call(ctx.args)
			if len(ret) > 0 {
				ctx.Result = ret[0].Interface()
			}
		}
	}
}

func (ctx *Context) ServeFile(path string) error {
	http.ServeFile(ctx, ctx.Req(), path)
	if ctx.Status() != http.StatusOK {
		return errors.New("serve file failed")
	}
	return nil
}

func (ctx *Context) ServeReader(rd io.Reader) error {
	_, err := io.Copy(ctx.ResponseWriter, rd)
	return err
}

func (ctx *Context) ServeXml(obj interface{}) error {
	dt, err := xml.Marshal(obj)
	if err != nil {
		return err
	}
	ctx.Header().Set("Content-Type", "application/xml")
	_, err = ctx.Write(dt)
	return err
}

func (ctx *Context) ServeJson(obj interface{}) error {
	dt, err := json.Marshal(obj)
	if err != nil {
		return err
	}
	ctx.Header().Set("Content-Type", "application/json")
	_, err = ctx.Write(dt)
	return err
}

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

func (ctx *Context) Redirect(url string) error {
	return redirect(ctx, url)
}

func redirect(w http.ResponseWriter, url string, status ...int) error {
	s := 302
	if len(status) > 0 {
		s = status[0]
	}
	w.Header().Set("Location", url)
	w.WriteHeader(s)
	_, err := w.Write([]byte("Redirecting to: " + url))
	return err
}

// Notmodified writes a 304 HTTP response
func (ctx *Context) NotModified() {
	ctx.WriteHeader(http.StatusNotModified)
}

// NotFound writes a 404 HTTP response
func (ctx *Context) NotFound(message ...string) error {
	if len(message) == 0 {
		return ctx.Abort(http.StatusNotFound, "Not Found")
	}
	return ctx.Abort(http.StatusNotFound, message[0])
}

// Abort is a helper method that sends an HTTP header and an optional
// body. It is useful for returning 4xx or 5xx errors.
// Once it has been called, any return value from the handler will
// not be written to the response.
func (ctx *Context) Abort(status int, body string) error {
	ctx.WriteHeader(status)
	ctx.Write([]byte(body))
	return nil
}

// SetHeader sets a response header. the current value
// of that header will be overwritten .
func (ctx *Context) SetHeader(key string, value string) {
	ctx.Header().Set(key, value)
}
