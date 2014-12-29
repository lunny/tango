package tango

import (
	"reflect"
	"path"
)

type groupRouter struct {
	methods []string
	url string
	c interface{}
}

type Group struct {
	routers []groupRouter
}

func NewGroup() *Group {
	return &Group{
		routers : make([]groupRouter, 0),
	}
}

func (g *Group) Get(url string, c interface{}) {
	g.Route([]string{"GET", "HEAD"}, url, c)
}

func (g *Group) Post(url string, c interface{}) {
	g.Route([]string{"POST"}, url, c)
}

func (g *Group) Head(url string, c interface{}) {
	g.Route([]string{"HEAD"}, url, c)
}

func (g *Group) Options(url string, c interface{}) {
	g.Route([]string{"OPTIONS"}, url, c)
}

func (g *Group) Trace(url string, c interface{}) {
	g.Route([]string{"TRACE"}, url, c)
}

func (g *Group) Patch(url string, c interface{}) {
	g.Route([]string{"PATCH"}, url, c)
}

func (g *Group) Delete(url string, c interface{}) {
	g.Route([]string{"DELETE"}, url, c)
}

func (g *Group) Put(url string, c interface{}) {
	g.Route([]string{"PUT"}, url, c)
}

func (g *Group) Any(url string, c interface{}) {
	g.Route(SupportMethods, url, c)
}

func (g *Group) Route(methods []string, url string, c interface{}) {
	g.routers = append(g.routers, groupRouter{methods, url, c})
}

func (t *Tango) addGroup(p string, g *Group) {
	for _, r := range g.routers {
		t.Route(r.methods, path.Join(p, r.url), r.c)
	}
}

var (
	gt = reflect.TypeOf(new(Group))
)

func (t *Tango) Group(p string, o interface{}) {
	vc := reflect.ValueOf(o)
	tp := vc.Type()
	var g *Group
	if tp == gt {
		g = o.(*Group)
	} else if tp.Kind() == reflect.Func && 
		tp.NumIn() == 1 && tp.In(0) == gt {
		g = NewGroup()
		vc.Call([]reflect.Value{reflect.ValueOf(g)})
	} else {
		panic("not allowed group parameter")
	}
	t.addGroup(p, g)
}
