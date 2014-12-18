package tango

import (
	"os"
	"testing"
)

func TestSwitch(t *testing.T) {
	logger := NewLogger(os.Stdout)

	t1 := NewWithLogger(logger)
	t1.Use(HandlerFunc(func(ctx *Context) {
		ctx.Write([]byte("tango 1"))
	}))

	t2 := NewWithLogger(logger)
	t2.Use(HandlerFunc(func(ctx *Context) {
		ctx.Write([]byte("tango 2"))
	}))

	t3 := NewWithLogger(logger)
	t3.Any("/(.*)", t1)
	t3.Any("/api/(.*)", t2)
}
