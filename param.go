package tango

import (
	"net/url"
)

type Paramer interface {
	SetParams(url.Values)
}

type Params struct {
	url.Values
}

func (p *Params) SetParams(u url.Values) {
	p.Values = u
}

func Param() HandlerFunc {
	return func(ctx *Context) {
		if action := ctx.Action(); action != nil {
			if p, ok := action.(Paramer); ok {
				p.SetParams(ctx.Params())
			}
		}
		ctx.Next()
	}
}