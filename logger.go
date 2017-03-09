// Copyright 2015 The Tango Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package tango

import (
	"io"
	"time"

	"github.com/lunny/log"
)

// Logger defines the logger interface for tango use
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

// CompositeLogger defines a composite loggers
type CompositeLogger struct {
	loggers []Logger
}

// NewCompositeLogger creates a composite loggers
func NewCompositeLogger(logs ...Logger) Logger {
	return &CompositeLogger{loggers: logs}
}

// Debugf implementes Logger interface
func (l *CompositeLogger) Debugf(format string, v ...interface{}) {
	for _, log := range l.loggers {
		log.Debugf(format, v...)
	}
}

// Debug implementes Logger interface
func (l *CompositeLogger) Debug(v ...interface{}) {
	for _, log := range l.loggers {
		log.Debug(v...)
	}
}

// Infof implementes Logger interface
func (l *CompositeLogger) Infof(format string, v ...interface{}) {
	for _, log := range l.loggers {
		log.Infof(format, v...)
	}
}

// Info implementes Logger interface
func (l *CompositeLogger) Info(v ...interface{}) {
	for _, log := range l.loggers {
		log.Info(v...)
	}
}

// Warnf implementes Logger interface
func (l *CompositeLogger) Warnf(format string, v ...interface{}) {
	for _, log := range l.loggers {
		log.Warnf(format, v...)
	}
}

// Warn implementes Logger interface
func (l *CompositeLogger) Warn(v ...interface{}) {
	for _, log := range l.loggers {
		log.Warn(v...)
	}
}

// Errorf implementes Logger interface
func (l *CompositeLogger) Errorf(format string, v ...interface{}) {
	for _, log := range l.loggers {
		log.Errorf(format, v...)
	}
}

// Error implementes Logger interface
func (l *CompositeLogger) Error(v ...interface{}) {
	for _, log := range l.loggers {
		log.Error(v...)
	}
}

// NewLogger use the default logger with special writer
func NewLogger(out io.Writer) Logger {
	l := log.New(out, "[tango] ", log.Ldefault())
	l.SetOutputLevel(log.Ldebug)
	return l
}

// LogInterface defines logger interface to inject logger to struct
type LogInterface interface {
	SetLogger(Logger)
}

// Log implementes LogInterface
type Log struct {
	Logger
}

// SetLogger implementes LogInterface
func (l *Log) SetLogger(log Logger) {
	l.Logger = log
}

// Logging returns handler to log informations
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
