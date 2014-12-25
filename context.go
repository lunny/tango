package tango

import (
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"net/url"
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
	params url.Values
	callArgs   []reflect.Value
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

func (ctx *Context) Params() url.Values {
	ctx.newAction()
	return ctx.params
}

func (ctx *Context) Action() interface{} {
	ctx.newAction()
	return ctx.action
}

func (ctx *Context) newAction() {
	if !ctx.matched {
		reqPath := removeStick(ctx.Req().URL.Path)
		ctx.route, ctx.params = ctx.router.Match(reqPath, ctx.Req().Method)
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
			case FuncHttpRoute:
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
			ret := ctx.route.method.Call(ctx.callArgs)
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

type Contexter interface {
	SetContext(*Context)
}

type Ctx struct {
	*Context
}

func (c *Ctx) SetContext(ctx *Context) {
	c.Context = ctx
}

func ContextHandler(ctx *Context) {
	if action := ctx.Action(); action != nil {
		if a, ok := action.(Contexter); ok {
			a.SetContext(ctx)
		}
	}
	ctx.Next()
}
