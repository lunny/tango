package tango

import (
	"net/http"
)

type HttpRequestInterface interface {
	SetRequest(*http.Request)
}

func RequestHandler(ctx *Context) {
	if action := ctx.Action(); action != nil {
		if s, ok := action.(HttpRequestInterface); ok {
			s.SetRequest(ctx.Req())
		}
	}

	ctx.Next()
}
