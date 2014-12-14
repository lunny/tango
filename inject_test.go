package tango

import (
	"fmt"
	"testing"
)

type Tester struct {
}

type TestHandler struct {
	tester *Tester
}

func (t *TestHandler) SetTester(tester *Tester) {
	t.tester = tester
}

func (t *TestHandler) Handle(ctx *Context) {
	ctx.Next()
}

func TestInjector(t *testing.T) {
	c := NewInjector()
	c.Map(&Tester{})
	var itor TestHandler
	c.Inject(&itor)

	fmt.Println("itor.tester:", itor.tester)
	if itor.tester == nil {
		t.Error("itor.tester is nil")
	}
}
