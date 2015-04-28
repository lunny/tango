// Copyright 2015 The Tango Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

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
	return fmt.Sprintf("%v", a.content)
}

func Abort(code int, content ...string) AbortError {
	if len(content) >= 1 {
		return &abortError{code, content[0]}
	}
	return &abortError{code, http.StatusText(code)}
}

func NotFound(content ...string) AbortError {
	return Abort(http.StatusNotFound, content...)
}

func NotSupported(content ...string) AbortError {
	return Abort(http.StatusMethodNotAllowed, content...)
}

func InternalServerError(content ...string) AbortError {
	return Abort(http.StatusInternalServerError, content...)
}

func Forbidden(content ...string) AbortError {
	return Abort(http.StatusForbidden, content...)
}

func Unauthorized(content ...string) AbortError {
	return Abort(http.StatusUnauthorized, content...)
}

// default errorhandler, you can use your self handler
func Errors() HandlerFunc {
	return func(ctx *Context) {
		switch res := ctx.Result.(type) {
		case AbortError:
			ctx.WriteHeader(res.Code())
			ctx.Write([]byte(res.Error()))
		case error:
			ctx.WriteHeader(http.StatusInternalServerError)
			ctx.Write([]byte(res.Error()))
		default:
			ctx.WriteHeader(http.StatusInternalServerError)
			ctx.Write([]byte(http.StatusText(http.StatusInternalServerError)))
		}
	}
}
