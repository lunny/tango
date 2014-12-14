package tango

import (
	"reflect"
	"strings"
)

type Injector struct {
	objs map[reflect.Type]interface{}
}

func NewInjector() *Injector {
	return &Injector{
		objs: make(map[reflect.Type]interface{}),
	}
}

func (c *Injector) Map(obj interface{}) {
	c.objs[reflect.TypeOf(obj)] = obj
}

func (c *Injector) Inject(objs ...interface{}) {
	for _, obj := range objs {
		c.inject(obj)
	}
}

func (c *Injector) inject(obj interface{}) {
	vv := reflect.ValueOf(obj)
	t := vv.Type()
	for k, v := range c.objs {
		for i := 0; i < t.NumMethod(); i++ {
			m := t.Method(i)
			if m.Type.NumIn() == 2 &&
				strings.HasPrefix(m.Name, "Set") &&
				m.Type.In(1) == k {
				m.Func.Call([]reflect.Value{vv, reflect.ValueOf(v)})
			}
		}
	}
}
