package tango

import (
	"strings"
)

func Prefix(prefix string, handler Handler) HandlerFunc {
	return func(ctx *Context) {
		if strings.HasPrefix(ctx.Req().URL.Path, prefix) {
			handler.Handle(ctx)
			return
		}

		ctx.Next()
	}
}
