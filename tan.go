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

func Version() string {
	return "0.4.6.0807"
}

type Tango struct {
	Router
	handlers   []Handler
	logger     Logger
	ErrHandler Handler
	ctxPool    sync.Pool
	respPool   sync.Pool
}

var (
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

func (t *Tango) Logger() Logger {
	return t.logger
}

func (t *Tango) Get(url string, c interface{}) {
	t.Route([]string{"GET", "HEAD"}, url, c)
}

func (t *Tango) Post(url string, c interface{}) {
	t.Route([]string{"POST"}, url, c)
}

func (t *Tango) Head(url string, c interface{}) {
	t.Route([]string{"HEAD"}, url, c)
}

func (t *Tango) Options(url string, c interface{}) {
	t.Route([]string{"OPTIONS"}, url, c)
}

func (t *Tango) Trace(url string, c interface{}) {
	t.Route([]string{"TRACE"}, url, c)
}

func (t *Tango) Patch(url string, c interface{}) {
	t.Route([]string{"PATCH"}, url, c)
}

func (t *Tango) Delete(url string, c interface{}) {
	t.Route([]string{"DELETE"}, url, c)
}

func (t *Tango) Put(url string, c interface{}) {
	t.Route([]string{"PUT"}, url, c)
}

func (t *Tango) Any(url string, c interface{}) {
	t.Route(SupportMethods, url, c)
}

func (t *Tango) Use(handlers ...Handler) {
	t.handlers = append(t.handlers, handlers...)
}

func GetDefaultListenInfo() (string, int) {
	host := os.Getenv("HOST")
	if len(host) == 0 {
		host = "0.0.0.0"
	}
	_port, _ := strconv.ParseInt(os.Getenv("PORT"), 10, 32)
	port := int(_port)
	if port == 0 {
		port = 8000
	}
	return host, port
}

func GetAddress(args ...interface{}) string {
	host, port := GetDefaultListenInfo()

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

	addr := host + ":" + strconv.FormatInt(int64(port), 10)

	return addr
}

// Run the http server. Listening on os.GetEnv("PORT") or 8000 by default.
func (t *Tango) Run(args ...interface{}) {
	addr := GetAddress(args...)
	t.logger.Info("Listening on http", addr)

	err := http.ListenAndServe(addr, t)
	if err != nil {
		t.logger.Error(err)
	}
}

func (t *Tango) RunTLS(certFile, keyFile string, args ...interface{}) {
	addr := GetAddress(args...)

	t.logger.Info("Listening on https", addr)

	err := http.ListenAndServeTLS(addr, certFile, keyFile, t)
	if err != nil {
		t.logger.Error(err)
	}
}

type HandlerFunc func(ctx *Context)

func (h HandlerFunc) Handle(ctx *Context) {
	h(ctx)
}

func WrapBefore(handler http.Handler) HandlerFunc {
	return func(ctx *Context) {
		handler.ServeHTTP(ctx.ResponseWriter, ctx.Req())

		ctx.Next()
	}
}

func WrapAfter(handler http.Handler) HandlerFunc {
	return func(ctx *Context) {
		ctx.Next()

		handler.ServeHTTP(ctx.ResponseWriter, ctx.Req())
	}
}

func (t *Tango) UseHandler(handler http.Handler) {
	t.Use(WrapBefore(handler))
}

func (t *Tango) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	resp := t.respPool.Get().(*responseWriter)
	resp.reset(w)

	ctx := t.ctxPool.Get().(*Context)
	ctx.tan = t
	ctx.reset(req, resp)

	ctx.Invoke()

	// if there is no logging or error handle, so the last written check.
	if !ctx.Written() {
		p := req.URL.Path
		if len(req.URL.RawQuery) > 0 {
			p = p + "?" + req.URL.RawQuery
		}

		if ctx.Route() != nil {
			if ctx.Result == nil {
				ctx.Write([]byte(""))
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

func NewWithLog(logger Logger, handlers ...Handler) *Tango {
	tan := &Tango{
		Router:     NewRouter(),
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

func New(handlers ...Handler) *Tango {
	return NewWithLog(NewLogger(os.Stdout), handlers...)
}

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
