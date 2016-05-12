// Copyright 2015 The Tango Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package tango

type groupRouter struct {
	methods  interface{}
	url      string
	c        interface{}
	handlers []Handler
}

type Group struct {
	routers  []groupRouter
	handlers []Handler
}

func NewGroup() *Group {
	return &Group{
		routers:  make([]groupRouter, 0),
		handlers: make([]Handler, 0),
	}
}

func (g *Group) Use(handlers ...Handler) {
	g.handlers = append(g.handlers, handlers...)
}

func (g *Group) Get(url string, c interface{}, middlewares ...Handler) {
	g.Route([]string{"GET", "HEAD:Get"}, url, c, middlewares...)
}

func (g *Group) Post(url string, c interface{}, middlewares ...Handler) {
	g.Route([]string{"POST"}, url, c, middlewares...)
}

func (g *Group) Head(url string, c interface{}, middlewares ...Handler) {
	g.Route([]string{"HEAD"}, url, c, middlewares...)
}

func (g *Group) Options(url string, c interface{}, middlewares ...Handler) {
	g.Route([]string{"OPTIONS"}, url, c, middlewares...)
}

func (g *Group) Trace(url string, c interface{}, middlewares ...Handler) {
	g.Route([]string{"TRACE"}, url, c, middlewares...)
}

func (g *Group) Patch(url string, c interface{}, middlewares ...Handler) {
	g.Route([]string{"PATCH"}, url, c, middlewares...)
}

func (g *Group) Delete(url string, c interface{}, middlewares ...Handler) {
	g.Route([]string{"DELETE"}, url, c, middlewares...)
}

func (g *Group) Put(url string, c interface{}, middlewares ...Handler) {
	g.Route([]string{"PUT"}, url, c, middlewares...)
}

func (g *Group) Any(url string, c interface{}, middlewares ...Handler) {
	g.Route(SupportMethods, url, c, middlewares...)
	g.Route([]string{"HEAD:Get"}, url, c, middlewares...)
}

func (g *Group) Route(methods interface{}, url string, c interface{}, middlewares ...Handler) {
	g.routers = append(g.routers, groupRouter{methods, url, c, middlewares})
}

func (g *Group) Group(p string, o interface{}) {
	gr := getGroup(o)
	for _, gchild := range gr.routers {
		g.Route(gchild.methods, joinRoute(p, gchild.url), gchild.c, gchild.handlers...)
	}
}

func getGroup(o interface{}) *Group {
	var g *Group
	var gf func(*Group)
	var ok bool
	if g, ok = o.(*Group); ok {
	} else if gf, ok = o.(func(*Group)); ok {
		g = NewGroup()
		gf(g)
	} else {
		panic("not allowed group parameter")
	}
	return g
}

func joinRoute(p, url string) string {
	if len(p) == 0 || p == "/" {
		return url
	}
	return p + url
}

func (t *Tango) addGroup(p string, g *Group) {
	for _, r := range g.routers {
		t.Route(r.methods, joinRoute(p, r.url), r.c, append(g.handlers, r.handlers...)...)
	}
}

func (t *Tango) Group(p string, o interface{}) {
	t.addGroup(p, getGroup(o))
}
