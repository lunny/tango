package tango

import (
	"net/http"
)

func ReturnHandler(ctx *Context) {
	ctx.Next()

	// if has been write, then return
	if ctx.Written() {
		return
	}

	if isNil(ctx.Result) {
		if ctx.Action() == nil {
			// if there is no action match
			ctx.Result = NotFound()
		} else {
			// there is an action but return nil, then we return blank page
			ctx.Result = ""
		}
	}

	if err, ok := ctx.Result.(AbortError); ok {
		ctx.WriteHeader(err.Code())
		ctx.Write([]byte(err.Error()))
	} else if err, ok := ctx.Result.(error); ok {
		ctx.WriteHeader(http.StatusInternalServerError)
		ctx.Write([]byte(err.Error()))
	} else if bs, ok := ctx.Result.([]byte); ok {
		ctx.WriteHeader(http.StatusOK)
		ctx.Write(bs)
	} else if s, ok := ctx.Result.(string); ok {
		ctx.WriteHeader(http.StatusOK)
		ctx.Write([]byte(s))
	}
}
