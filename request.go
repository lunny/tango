package tango

import (
	"net/http"
)

type Requester interface {
	SetRequest(*http.Request)
}

func RequestHandler(ctx *Context) {
	if action := ctx.Action(); action != nil {
		if s, ok := action.(Requester); ok {
			s.SetRequest(ctx.Req())
		}
	}

	ctx.Next()
}
