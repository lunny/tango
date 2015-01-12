package tango

import (
	"fmt"
	"net/http"
)

type AbortError interface {
	error
	Code() int
}

type abortError struct {
	code    int
	content string
}

func (a *abortError) Code() int {
	return a.code
}

func (a *abortError) Error() string {
	return fmt.Sprintf("%v %v", a.code, a.content)
}

func Abort(code int, content ...string) error {
	if len(content) >= 1 {
		return &abortError{code, content[0]}
	}
	return &abortError{code, http.StatusText(code)}
}

func NotFound(content ...string) error {
	return Abort(http.StatusNotFound, content...)
}

func NotSupported(content ...string) error {
	return Abort(http.StatusMethodNotAllowed, content...)
}

func InternalServerError(content ...string) error {
	return Abort(http.StatusInternalServerError, content...)
}

func Forbidden(content ...string) error {
	return Abort(http.StatusForbidden, content...)
}

func Unauthorized(content ...string) error {
	return Abort(http.StatusUnauthorized, content...)
}

func Errors() HandlerFunc {
	return func(ctx *Context) {
		switch res := ctx.Result.(type) {
		case AbortError:
			ctx.WriteHeader(res.Code())
			ctx.Write([]byte(res.Error()))
		case error:
			ctx.WriteHeader(http.StatusInternalServerError)
			ctx.Write([]byte(res.Error()))
		case []byte:
			ctx.WriteHeader(http.StatusInternalServerError)
			ctx.Write(res)
		case string:
			ctx.WriteHeader(http.StatusInternalServerError)
			ctx.Write([]byte(res))
		default:
			ctx.WriteHeader(http.StatusInternalServerError)
			ctx.Write([]byte(http.StatusText(http.StatusInternalServerError)))
		}
	}
}