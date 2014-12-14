package tango

import (
	"net/http"
	"os"
	"time"
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
	Mode int
	handlers []Handler
	logger Logger
}

func NewWithLogger(logger Logger, handlers ...Handler) *Tango {
	tango := &Tango{
		Injector: NewInjector(),
		Router: NewRouter(),
		Mode: Dev,
		logger: logger,
	}

	tango.Use(handlers...)

	return tango
}

func (t *Tango) Use(handlers ...Handler) {
	for _, handler := range handlers {
		t.handlers = append(t.handlers, handler)
		t.Map(handler)
	}
}

func (t *Tango) Run(addr string) {
	if t.logger != nil {
		t.logger.Info("listening on", addr, modes[t.Mode])
	}
	http.ListenAndServe(addr, t)
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

func New(handlers ...Handler) *Tango {
	return NewWithLogger(NewLogger(os.Stdout), handlers...)
}

func Classic() *Tango {
	logger := NewLogger(os.Stdout)
	return NewWithLogger(
		logger,
		NewRecovery(true),
		NewLogging(logger),
		NewCompress([]string{".js", ".css", ".html", ".htm"}),
		HandlerFunc(ReturnHandler),
		NewStatic("public", "public", []string{"index.html", "index.htm"}),
		HandlerFunc(ResponseHandler),
		HandlerFunc(RequestHandler),
	)
}

func Full() *Tango {
	logger := NewLogger(os.Stdout)
	return NewWithLogger(
		logger,
		NewRecovery(true),
		NewLogging(logger),
		NewCompress([]string{".js", ".css", ".html", ".htm"}),
		HandlerFunc(ReturnHandler),
		NewStatic("public", "public", []string{"index.html", "index.htm"}),
		HandlerFunc(ResponseHandler),
		HandlerFunc(RequestHandler),
		NewSessions(time.Minute * 20),
		NewRender("templates", true, true),
	)
}

func Static() *Tango {
	logger := NewLogger(os.Stdout)
	return NewWithLogger(
		logger,
		NewLogging(logger),
		NewCompress([]string{".js", ".css", ".html", ".htm"}),
		NewStatic("public", "", []string{"index.html", "index.htm"}),
	)
}