package tango

import (
	"net/http"
	"os"
	"log"
)

const (
	Dev = iota
	Test
	Product
)

var (
	modes = []string{
		"Dev",
		"Test",
		"Product",
	}
)

func Version() string {
	return "0.1.0.1213"
}

type Tango struct {
	*Injector
	Router
	Mode     int
	handlers []Handler
	logger   Logger
}

func (tango *Tango) Get(url string, c interface{}) {
	tango.AddRouter(url, []string{"GET"}, c)
}

func (tango *Tango) Post(url string, c interface{}) {
	tango.AddRouter(url, []string{"POST"}, c)
}

func (tango *Tango) Head(url string, c interface{}) {
	tango.AddRouter(url, []string{"HEAD"}, c)
}

func (tango *Tango) Options(url string, c interface{}) {
	tango.AddRouter(url, []string{"OPTIONS"}, c)
}

func (tango *Tango) Trace(url string, c interface{}) {
	tango.AddRouter(url, []string{"TRACE"}, c)
}

func (tango *Tango) Patch(url string, c interface{}) {
	tango.AddRouter(url, []string{"PATCH"}, c)
}

func (tango *Tango) Delete(url string, c interface{}) {
	tango.AddRouter(url, []string{"DELETE"}, c)
}

func (tango *Tango) Put(url string, c interface{}) {
	tango.AddRouter(url, []string{"PUT"}, c)
}

func (tango *Tango) Any(url string, c interface{}) {
	tango.AddRouter(url, SupportMethods, c)
}

func (t *Tango) Use(handlers ...Handler) {
	for _, handler := range handlers {
		t.handlers = append(t.handlers, handler)
		t.Map(handler)
	}
	t.injectAll()
}

func (t *Tango) Run(addr string) {
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
	for _, handler := range t.handlers {
		t.Inject(handler)
	}
}

func (t *Tango) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	ctx := NewContext(
		t.Router,
		t.handlers,
		req,
		NewResponseWriter(w),
	)

	ctx.Invoke()
}

func (t *Tango) Handle(ctx *Context) {
	ctx.Invoke()
}

func New(handlers ...Handler) *Tango {
	return NewWithLog(NewLogger(os.Stdout), handlers...)
}

func NewWithLog(logger Logger, handlers ...Handler) *Tango {
	tango := &Tango{
		Injector: NewInjector(),
		Router:   NewRouter(),
		Mode:     Dev,
		logger:   logger,
	}

	tango.Map(logger)
	tango.Use(handlers...)

	return tango
}

func Classic() *Tango {
	logger := NewLogger(os.Stdout)
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

func Static() *Tango {
	logger := NewLogger(os.Stdout)
	return NewWithLog(
		logger,
		NewLogging(logger),
		NewCompress([]string{".js", ".css", ".html", ".htm"}),
		NewStatic("./public", "", []string{"index.html", "index.htm"}),
	)
}
