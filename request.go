package tango

import (
	"net/http"
)

type Requester interface {
	SetRequest(*http.Request)
}

type Req struct {
	req *http.Request
}

func (r *Req) SetRequest(req *http.Request) {
	r.req = req
}

func (r *Req) Req() *http.Request {
	return r.req
}

func RequestHandler(ctx *Context) {
	if action := ctx.Action(); action != nil {
		if s, ok := action.(Requester); ok {
			s.SetRequest(ctx.Req())
		}
	}

	ctx.Next()
}
