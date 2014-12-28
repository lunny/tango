package tango

import (
	"net/http"
	"os"
	"log"
	"time"
)

const (
	Dev = iota
	Test
	Prod
)

var (
	Env = Dev

	modes = []string{
		"Dev",
		"Test",
		"Product",
	}
)

func Version() string {
	return "0.1.0.1228"
}

type Tango struct {
	Router
	Mode     int
	handlers []Handler
	injectors []Injector
	logger   Logger
}

func (t *Tango) Get(url string, c interface{}) {
	t.AddRouter(url, []string{"GET"}, c)
	t.AddRouter(url, []string{"HEAD"}, c)
}

func (t *Tango) Post(url string, c interface{}) {
	t.AddRouter(url, []string{"POST"}, c)
}

func (t *Tango) Head(url string, c interface{}) {
	t.AddRouter(url, []string{"HEAD"}, c)
}

func (t *Tango) Options(url string, c interface{}) {
	t.AddRouter(url, []string{"OPTIONS"}, c)
}

func (t *Tango) Trace(url string, c interface{}) {
	t.AddRouter(url, []string{"TRACE"}, c)
}

func (t *Tango) Patch(url string, c interface{}) {
	t.AddRouter(url, []string{"PATCH"}, c)
}

func (t *Tango) Delete(url string, c interface{}) {
	t.AddRouter(url, []string{"DELETE"}, c)
}

func (t *Tango) Put(url string, c interface{}) {
	t.AddRouter(url, []string{"PUT"}, c)
}

func (t *Tango) Any(url string, c interface{}) {
	t.AddRouter(url, SupportMethods, c)
}

/*
func (t *Tango) Group(url string, cs map[string]interface{}) {
	for k, v :=
}*/

func (t *Tango) Use(handlers ...Handler) {
	for _, handler := range handlers {
		t.handlers = append(t.handlers, handler)
		if i, ok := interface{}(handler).(Injector); ok {
			t.injectors = append(t.injectors, i)
		}
	}
	t.injectAll()
}

func (t *Tango) Run(addrs ...string) {
	var addr string
	if len(addrs) == 0 {
		addr = ":8000"
	} else {
		addr = addrs[0]
	}

	if t.logger != nil {
		t.logger.Info("listening on", addr, modes[t.Mode])
	}else {
		log.Println("listening on", addr, modes[t.Mode])
	}

	err := http.ListenAndServe(addr, t)
	if err != nil {
		if t.logger != nil {
			t.logger.Error(err)
		} else {
			log.Println(err)
		}
	}
}

type HandlerFunc func(ctx *Context)

func (h HandlerFunc) Handle(ctx *Context) {
	h(ctx)
}

func WrapBefore(handler http.Handler) HandlerFunc {
	return HandlerFunc(func(ctx *Context) {
		handler.ServeHTTP(ctx.ResponseWriter, ctx.Req())

		ctx.Next()
	})
}

func WrapAfter(handler http.Handler) HandlerFunc {
	return HandlerFunc(func(ctx *Context) {
		ctx.Next()

		handler.ServeHTTP(ctx.ResponseWriter, ctx.Req())
	})
}

func (t *Tango) UseHandler(handler http.Handler) {
	t.Use(WrapBefore(handler))
}

func (t *Tango) injectAll() {
	for _, injector := range t.injectors {
		for _, handler := range t.handlers {
			injector.Inject(handler)
		}
	}
}

func (t *Tango) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	start := time.Now()

	ctx := NewContext(
		t.Router,
		t.handlers,
		req,
		NewResponseWriter(w),
	)

	ctx.Invoke()

	if !ctx.Written() {
		if t.logger != nil {
			ctx.WriteHeader(http.StatusNotFound)
			escape := time.Now().Sub(start)
			t.logger.Error(ctx.Req().Method, http.StatusNotFound, escape, req.URL.Path)
		}
	}
}

func (t *Tango) Handle(ctx *Context) {
	ctx.Invoke()
}

func New(handlers ...Handler) *Tango {
	return NewWithLog(NewLogger(os.Stdout), handlers...)
}

func NewWithLog(logger Logger, handlers ...Handler) *Tango {
	tango := &Tango{
		Router:   NewRouter(),
		Mode:     Env,
		logger:   logger,
		handlers: make([]Handler, 0),
		injectors: make([]Injector, 0),
	}

	tango.Use(handlers...)

	return tango
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
		NewLogging(logger),
		NewRecovery(true),
		NewCompress([]string{".js", ".css", ".html", ".htm"}),
		NewStatic("./public", "public", []string{"index.html", "index.htm"}),
		HandlerFunc(ReturnHandler),
		HandlerFunc(ResponseHandler),
		HandlerFunc(RequestHandler),
		HandlerFunc(ParamHandler),
		HandlerFunc(ContextHandler),
	)
}

func Static(l ...Logger) *Tango {
	var logger Logger
	if len(l) == 0 {
		logger = NewLogger(os.Stdout)
	} else {
		logger = l[0]
	}

	return NewWithLog(
		logger,
		NewLogging(logger),
		NewCompress([]string{".js", ".css", ".html", ".htm"}),
		NewStatic("./public", "", []string{"index.html", "index.htm"}),
	)
}
