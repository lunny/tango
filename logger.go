// Copyright 2015 The Tango Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package tango

import (
	"io"
	"time"

	"github.com/lunny/log"
)

type Logger interface {
	Debugf(format string, v ...interface{})
	Debug(v ...interface{})
	Infof(format string, v ...interface{})
	Info(v ...interface{})
	Warnf(format string, v ...interface{})
	Warn(v ...interface{})
	Errorf(format string, v ...interface{})
	Error(v ...interface{})
}

func NewLogger(out io.Writer) Logger {
	l := log.New(out, "[tango] ", log.Ldefault())
	l.SetOutputLevel(log.Ldebug)
	return l
}

type LogInterface interface {
	SetLogger(Logger)
}

type Log struct {
	Logger
}

func (l *Log) SetLogger(log Logger) {
	l.Logger = log
}

func Logging() HandlerFunc {
	return func(ctx *Context) {
		start := time.Now()
		p := ctx.Req().URL.Path
		if len(ctx.Req().URL.RawQuery) > 0 {
			p = p + "?" + ctx.Req().URL.RawQuery
		}

		ctx.Debug("Started", ctx.Req().Method, p, "for", ctx.IP())

		if action := ctx.Action(); action != nil {
			if l, ok := action.(LogInterface); ok {
				l.SetLogger(ctx.Logger)
			}
		}

		ctx.Next()

		if !ctx.Written() {
			if ctx.Result == nil {
				ctx.Result = NotFound()
			}
			ctx.HandleError()
		}

		statusCode := ctx.Status()

		if statusCode >= 200 && statusCode < 400 {
			ctx.Info(ctx.Req().Method, statusCode, time.Since(start), p)
		} else {
			ctx.Error(ctx.Req().Method, statusCode, time.Since(start), p, ctx.Result)
		}
	}
}
