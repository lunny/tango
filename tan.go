// Copyright 2015 The Tango Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package tango

import (
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
)

// Version returns tango's version
func Version() string {
	return "0.5.4.0517"
}

// Tango describes tango object
type Tango struct {
	Router
	handlers   []Handler
	logger     Logger
	ErrHandler Handler
	ctxPool    sync.Pool
	respPool   sync.Pool
}

var (
	// ClassicHandlers the default handlers
	ClassicHandlers = []Handler{
		Logging(),
		Recovery(false),
		Compresses([]string{}),
		Static(StaticOptions{Prefix: "public"}),
		Return(),
		Param(),
		Contexts(),
	}
)

// Logger returns logger interface
func (t *Tango) Logger() Logger {
	return t.logger
}

// Get sets a route with GET method
func (t *Tango) Get(url string, c interface{}, middlewares ...Handler) {
	t.Route([]string{"GET", "HEAD:Get"}, url, c, middlewares...)
}

// Post sets a route with POST method
func (t *Tango) Post(url string, c interface{}, middlewares ...Handler) {
	t.Route([]string{"POST"}, url, c, middlewares...)
}

// Head sets a route with HEAD method
func (t *Tango) Head(url string, c interface{}, middlewares ...Handler) {
	t.Route([]string{"HEAD"}, url, c, middlewares...)
}

// Options sets a route with OPTIONS method
func (t *Tango) Options(url string, c interface{}, middlewares ...Handler) {
	t.Route([]string{"OPTIONS"}, url, c, middlewares...)
}

// Trace sets a route with TRACE method
func (t *Tango) Trace(url string, c interface{}, middlewares ...Handler) {
	t.Route([]string{"TRACE"}, url, c, middlewares...)
}

// Patch sets a route with PATCH method
func (t *Tango) Patch(url string, c interface{}, middlewares ...Handler) {
	t.Route([]string{"PATCH"}, url, c, middlewares...)
}

// Delete sets a route with DELETE method
func (t *Tango) Delete(url string, c interface{}, middlewares ...Handler) {
	t.Route([]string{"DELETE"}, url, c, middlewares...)
}

// Put sets a route with PUT method
func (t *Tango) Put(url string, c interface{}, middlewares ...Handler) {
	t.Route([]string{"PUT"}, url, c, middlewares...)
}

// Any sets a route every support method is OK.
func (t *Tango) Any(url string, c interface{}, middlewares ...Handler) {
	t.Route(SupportMethods, url, c, middlewares...)
	t.Route([]string{"HEAD:Get"}, url, c, middlewares...)
}

// Use addes some global handlers
func (t *Tango) Use(handlers ...Handler) {
	t.handlers = append(t.handlers, handlers...)
}

// GetAddress parses address
func getAddress(args ...interface{}) string {
	var host string
	var port int

	if len(args) == 1 {
		switch arg := args[0].(type) {
		case string:
			addrs := strings.Split(args[0].(string), ":")
			if len(addrs) == 1 {
				host = addrs[0]
			} else if len(addrs) >= 2 {
				host = addrs[0]
				_port, _ := strconv.ParseInt(addrs[1], 10, 0)
				port = int(_port)
			}
		case int:
			port = arg
		}
	} else if len(args) >= 2 {
		if arg, ok := args[0].(string); ok {
			host = arg
		}
		if arg, ok := args[1].(int); ok {
			port = arg
		}
	}

	if envHost := os.Getenv("HOST"); len(envHost) != 0 {
		host = envHost
	} else if len(host) == 0 {
		host = "0.0.0.0"
	}

	if envPort, _ := strconv.ParseInt(os.Getenv("PORT"), 10, 32); envPort != 0 {
		port = int(envPort)
	} else if port == 0 {
		port = 8000
	}

	addr := host + ":" + strconv.FormatInt(int64(port), 10)

	return addr
}

// Run the http server. Listening on os.GetEnv("PORT") or 8000 by default.
func (t *Tango) Run(args ...interface{}) {
	addr := getAddress(args...)
	t.logger.Info("Listening on http://" + addr)

	err := http.ListenAndServe(addr, t)
	if err != nil {
		t.logger.Error(err)
	}
}

// RunTLS runs the https server with special cert and key files
func (t *Tango) RunTLS(certFile, keyFile string, args ...interface{}) {
	addr := getAddress(args...)

	t.logger.Info("Listening on https://" + addr)

	err := http.ListenAndServeTLS(addr, certFile, keyFile, t)
	if err != nil {
		t.logger.Error(err)
	}
}

// HandlerFunc describes the handle function
type HandlerFunc func(ctx *Context)

// Handle executes the handler
func (h HandlerFunc) Handle(ctx *Context) {
	h(ctx)
}

// WrapBefore wraps a http standard handler to tango's before action executes
func WrapBefore(handler http.Handler) HandlerFunc {
	return func(ctx *Context) {
		handler.ServeHTTP(ctx.ResponseWriter, ctx.Req())

		ctx.Next()
	}
}

// WrapAfter wraps a http standard handler to tango's after action executes
func WrapAfter(handler http.Handler) HandlerFunc {
	return func(ctx *Context) {
		ctx.Next()

		handler.ServeHTTP(ctx.ResponseWriter, ctx.Req())
	}
}

// UseHandler adds a standard http handler to tango's
func (t *Tango) UseHandler(handler http.Handler) {
	t.Use(WrapBefore(handler))
}

// ServeHTTP implementes net/http interface so that it could run with net/http
func (t *Tango) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	resp := t.respPool.Get().(*responseWriter)
	resp.reset(w)

	ctx := t.ctxPool.Get().(*Context)
	ctx.tan = t
	ctx.reset(req, resp)

	ctx.invoke()

	// if there is no logging or error handle, so the last written check.
	if !ctx.Written() {
		p := req.URL.Path
		if len(req.URL.RawQuery) > 0 {
			p = p + "?" + req.URL.RawQuery
		}

		if ctx.Route() != nil {
			if ctx.Result == nil {
				ctx.WriteString("")
				t.logger.Info(req.Method, ctx.Status(), p)
				t.ctxPool.Put(ctx)
				t.respPool.Put(resp)
				return
			}
			panic("result should be handler before")
		}

		if ctx.Result == nil {
			ctx.Result = NotFound()
		}

		ctx.HandleError()

		t.logger.Error(req.Method, ctx.Status(), p)
	}

	t.ctxPool.Put(ctx)
	t.respPool.Put(resp)
}

// NewWithLog creates tango with the special logger and handlers
func NewWithLog(logger Logger, handlers ...Handler) *Tango {
	tan := &Tango{
		Router:     newRouter(),
		logger:     logger,
		handlers:   make([]Handler, 0),
		ErrHandler: Errors(),
	}

	tan.ctxPool.New = func() interface{} {
		return &Context{
			tan:    tan,
			Logger: tan.logger,
		}
	}

	tan.respPool.New = func() interface{} {
		return &responseWriter{}
	}

	tan.Use(handlers...)

	return tan
}

// New creates tango with the default logger and handlers
func New(handlers ...Handler) *Tango {
	return NewWithLog(NewLogger(os.Stdout), handlers...)
}

// Classic returns the tango with default handlers and logger
func Classic(l ...Logger) *Tango {
	var logger Logger
	if len(l) == 0 {
		logger = NewLogger(os.Stdout)
	} else {
		logger = l[0]
	}

	return NewWithLog(
		logger,
		ClassicHandlers...,
	)
}
