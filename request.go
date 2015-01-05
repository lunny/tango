package tango

import (
	"net/http"
)

type Requester interface {
	SetRequest(*http.Request)
}

type Req struct {
	*http.Request
}

func (r *Req) SetRequest(req *http.Request) {
	r.Request = req
}

func Requests() HandlerFunc {
	return func(ctx *Context) {
		if action := ctx.Action(); action != nil {
			if s, ok := action.(Requester); ok {
				s.SetRequest(ctx.Req())
			}
		}

		ctx.Next()
	}
}
