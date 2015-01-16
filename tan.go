package tango

import (
	"net/http"
	"os"
)

const (
	Dev = iota
	Prod
)

var (
	Env = Dev

	modes = []string{
		"Dev",
		"Product",
	}
)

func Version() string {
	return "0.2.6.0116"
}

type Tango struct {
	Router
	Mode     int
	handlers []Handler
	logger   Logger
	ErrHandler Handler
}

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
	for _, handler := range handlers {
		t.handlers = append(t.handlers, handler)
	}
}

func (t *Tango) Run(addrs ...string) {
	var addr string
	if len(addrs) == 0 {
		addr = ":8000"
	} else {
		addr = addrs[0]
	}

	t.logger.Info("listening on", addr, modes[t.Mode])

	err := http.ListenAndServe(addr, t)
	if err != nil {
		t.logger.Error(err)
	}
}

func (t *Tango) RunTLS(certFile, keyFile string, addrs ...string) {
	var addr string
	if len(addrs) == 0 {
		addr = ":8000"
	} else {
		addr = addrs[0]
	}

	t.logger.Info("listening on https", addr, modes[t.Mode])

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
	ctx := NewContext(
		t.Router,
		t.handlers,
		req,
		NewResponseWriter(w),
		t.logger,
		t.ErrHandler,
	)

	ctx.Invoke()

	// if there is no logging or error handle, so the last written check.
	if !ctx.Written() {
		if ctx.Result == nil {
			ctx.Result = NotFound()
		}
		ctx.HandleError()
		p := req.URL.Path 
		if len(req.URL.RawQuery) > 0 {
			p = p + "?"+req.URL.RawQuery
		}

		t.logger.Error(req.Method, ctx.Status(), p)
	}
}

func NewWithLog(logger Logger, handlers ...Handler) *Tango {
	tango := &Tango{
		Router:   NewRouter(),
		Mode:     Env,
		logger:   logger,
		handlers: make([]Handler, 0),
		ErrHandler: Errors(),
	}

	tango.Use(handlers...)

	return tango
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
		Logging(),
		Recovery(true),
		Compresses([]string{".js", ".css", ".html", ".htm"}),
		Static(StaticOptions{Prefix:"public"}),
		Return(),
		Responses(),
		Requests(),
		Param(),
		Contexts(),
	)
}