package tango

import (
	"bytes"
	"fmt"
	"net/http"
	"reflect"
	"regexp"
	"strings"
)

type RouteType byte

const (
	FuncRoute         RouteType = iota + 1 // func ()
	FuncHttpRoute                          // func (http.ResponseWriter, *http.Request)
	FuncReqRoute                           // func (*http.Request)
	FuncResponseRoute                      // func (http.ResponseWriter)
	FuncCtxRoute                           // func (*tango.Context)
	StructRoute                            // func (st) <Get>()
	StructPtrRoute                         // func (*struct) <Get>()
)

var (
	SupportMethods = []string{
		"GET",
		"POST",
		"HEAD",
		"DELETE",
		"PUT",
		"OPTIONS",
		"TRACE",
		"PATCH",
	}

	PoolSize = 800
)

// Route
type Route struct {
	method    reflect.Value
	routeType RouteType
	pool      *pool
}

func NewRoute(t reflect.Type,
	method reflect.Value, tp RouteType) *Route {
	var pool *pool
	if tp == StructRoute || tp == StructPtrRoute {
		pool = newPool(PoolSize, t)
	}
	return &Route{
		routeType: tp,
		method:    method,
		pool:      pool,
	}
}

func (r *Route) Method() reflect.Value {
	return r.method
}

func (r *Route) RouteType() RouteType {
	return r.routeType
}

func (r *Route) IsStruct() bool {
	return r.routeType == StructRoute || r.routeType == StructPtrRoute
}

func (r *Route) newAction() reflect.Value {
	if !r.IsStruct() {
		return r.method
	}

	return r.pool.New()
}

type Router interface {
	Route(methods interface{}, path string, handler interface{})
	Match(requestPath, method string) (*Route, Params)
}

var specialBytes = []byte(`.\+*?|[]{}^$`)

func isSpecial(ch byte) bool {
	return bytes.IndexByte(specialBytes, ch) > -1
}

func isAlpha(ch byte) bool {
	return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z' || ch == '_'
}

func isDigit(ch byte) bool {
	return '0' <= ch && ch <= '9'
}

func isAlnum(ch byte) bool {
	return isAlpha(ch) || isDigit(ch)
}

type (
	router struct {
		trees map[string]*node
	}
	ntype byte
	node  struct {
		tp      ntype          // Type of node it contains
		handle  *Route         // executor
		regexp  *regexp.Regexp // regexp if tp is rnode
		content string         // static content or named
		edges   edges          // children
	}
	edges []*node
)

const (
	snode ntype = iota // static, should equal
	nnode              // named node, match a non-/ is ok
	anode              // catch-all node, match any
	rnode              // regex node, should match
)

func (n *node) equal(o *node) bool {
	if n.tp != o.tp || n.content != o.content {
		return false
	}
	return true
}

func NewRouter() (r *router) {
	r = &router{
		trees: make(map[string]*node),
	}
	for _, m := range SupportMethods {
		r.trees[m] = &node{
			edges: edges{},
		}
	}
	return
}

//   /:name1/:name2 /:name1-:name2 /(:name1)sss(:name2)
//   /(*name) /(:name[0-9]+) /(:name[a-z]+)
func parseNodes(path string) []*node {
	var i, j int
	l := len(path)
	var nodes = make([]*node, 0)
	var bracket int
	for ; i < l; i++ {
		if path[i] == ':' {
			nodes = append(nodes, &node{tp: snode, content: path[j : i-bracket]})
			j = i
			var regex string
			if bracket == 1 {
				var start = -1
				for ; i < l && ')' != path[i]; i++ {
					if start == -1 && isSpecial(path[i]) {
						start = i
					}
				}
				if path[i] != ')' {
					panic("lack of )")
				}
				if start > -1 {
					regex = path[start:i]
				}
			} else {
				i = i + 1
				for ; i < l && isAlnum(path[i]); i++ {
				}
			}

			if len(regex) > 0 {
				nodes = append(nodes, &node{tp: rnode,
					regexp:  regexp.MustCompile("(" + regex + ")"),
					content: path[j : i-len(regex)]})
			} else {
				nodes = append(nodes, &node{tp: nnode, content: path[j:i]})
			}
			i = i + bracket
			j = i
			if i == l {
				return nodes
			}
		} else if path[i] == '*' {
			nodes = append(nodes, &node{tp: snode, content: path[j:i]})
			j = i
			if bracket == 1 {
				for ; i < l && ')' == path[i]; i++ {
				}
			} else {
				i = i + 1
				for ; i < l && isAlnum(path[i]); i++ {
				}
			}
			nodes = append(nodes, &node{tp: anode, content: path[j:i]})
			j = i
			if i == l {
				return nodes
			}
		} else if path[i] == '(' {
			bracket = 1
		} else if path[i] == '/' {
			if bracket == 0 && i > j {
				nodes = append(nodes, &node{tp: snode, content: path[j:i]})
				j = i
			}
		} else {
			bracket = 0
		}
	}

	nodes = append(nodes, &node{
		tp:      snode,
		content: path[j:i],
	})

	return nodes
}

func printNode(i int, node *node) {
	for _, c := range node.edges {
		for j := 0; j < i; j++ {
			fmt.Print(" ")
		}
		fmt.Println(c.content, c.handle)
		printNode(i+1, c)
	}
}

func (r *router) printTrees() {
	for _, method := range SupportMethods {
		fmt.Print(method)
		printNode(1, r.trees[method])
		fmt.Println()
	}
}

func (r *router) addRoute(method, path string, h *Route) {
	nodes := parseNodes(path)
	nodes[len(nodes)-1].handle = h
	if !validNodes(nodes) {
		panic(fmt.Sprintln("express", path, "is not supported"))
	}
	r.addnodes(method, nodes)
}

func (r *router) matchNode(n *node, url string, params Params) (*Route, Params) {
	if n.tp == snode {
		if strings.HasPrefix(url, n.content) {
			if len(url) == len(n.content) {
				return n.handle, params
			}
			for _, c := range n.edges {
				e, p := r.matchNode(c, url[len(n.content):], params)
				if e != nil {
					return e, p
				}
			}
		}
	} else if n.tp == anode {
		if len(n.edges) == 0 {
			params = append(params, param{n.content, url})
			return n.handle, params
		}
		for _, c := range n.edges {
			if c.tp == snode {
				idx := strings.Index(url, c.content)
				if idx > -1 {
					params = append(params, param{n.content, url[:idx]})
					return r.matchNode(c, url[idx:], params)
				}
			} else {
				panic("should be snode")
			}
		}
	} else if n.tp == nnode {
		idx := strings.IndexByte(url, '/')
		if idx > -1 {
			params = append(params, param{n.content, url[:idx]})
			for _, c := range n.edges {
				h, p := r.matchNode(c, url[idx:], params)
				if h != nil {
					return h, p
				}
			}
			return nil, nil
		}

		if len(n.edges) == 0 {
			params = append(params, param{n.content, url})
			return n.handle, params
		}
		for _, c := range n.edges {
			if c.tp == snode {
				idx := strings.Index(url, c.content)
				if idx > -1 {
					params = append(params, param{n.content, url[:idx]})
					return r.matchNode(c, url[idx:], params)
				}
			} else {
				panic("should be snode")
			}
		}
	} else if n.tp == rnode {
		if len(n.edges) == 0 && n.regexp.MatchString(url) {
			params = append(params, param{n.content, url})
			return n.handle, params
		}
		for _, c := range n.edges {
			if c.tp == snode {
				idx := strings.Index(url, c.content)
				if idx > -1 && n.regexp.MatchString(url[:idx]) {
					params = append(params, param{n.content, url[:idx]})
					return r.matchNode(c, url[idx:], params)
				}
			} else {
				panic("should be snode")
			}
		}
	}
	return nil, nil
}

func (r *router) Match(url, method string) (*Route, Params) {
	cn := r.trees[method]
	var params = make(Params, 0)
	for _, n := range cn.edges {
		e, p := r.matchNode(n, url, params)
		if e != nil {
			return e, p
		}
	}
	return nil, nil
}

func (r *router) addnode(p *node, nodes []*node, i int) *node {
	if len(p.edges) == 0 {
		p.edges = make([]*node, 0)
	}
	for _, pc := range p.edges {
		if pc.equal(nodes[i]) {
			if i == len(nodes)-1 {
				pc.handle = nodes[i].handle
			}
			return pc
		}
	}
	p.edges = append(p.edges, nodes[i])
	return p.edges[len(p.edges)-1]
}

func validNodes(nodes []*node) bool {
	if len(nodes) == 0 {
		return false
	}
	var lastTp = nodes[0]
	for _, node := range nodes[1:] {
		if lastTp.tp != snode && node.tp != snode {
			return false
		}
	}
	return true
}

func (r *router) addnodes(method string, nodes []*node) {
	cn := r.trees[method]
	var p *node = cn

	for i, _ := range nodes {
		p = r.addnode(p, nodes, i)
	}
}

func removeStick(uri string) string {
	uri = strings.TrimRight(uri, "/")
	if uri == "" {
		uri = "/"
	}
	return uri
}

func (router *router) Route(ms interface{}, url string, c interface{}) {
	vc := reflect.ValueOf(c)
	if vc.Kind() == reflect.Func {
		switch ms.(type) {
		case string:
			router.addFunc([]string{ms.(string)}, url, c)
		case []string:
			router.addFunc(ms.([]string), url, c)
		default:
			panic("unknow methods format")
		}
	} else if vc.Kind() == reflect.Ptr && vc.Elem().Kind() == reflect.Struct {
		var methods = make(map[string]string)
		switch ms.(type) {
		case string:
			s := strings.Split(ms.(string), ":")
			if len(s) == 1 {
				methods[s[0]] = strings.Title(strings.ToLower(s[0]))
			} else if len(s) == 2 {
				methods[s[0]] = strings.TrimSpace(s[1])
			} else {
				panic("unknow methods format")
			}
		case []string:
			for _, m := range ms.([]string) {
				s := strings.Split(m, ":")
				if len(s) == 1 {
					methods[s[0]] = strings.Title(strings.ToLower(s[0]))
				} else if len(s) == 2 {
					methods[s[0]] = strings.TrimSpace(s[1])
				} else {
					panic("unknow format")
				}
			}
		case map[string]string:
			methods = ms.(map[string]string)
		default:
			panic("unsupported methods")
		}

		router.addStruct(methods, url, c)
	} else {
		panic("not support route type")
	}
}

/*
	Tango supports 5 form funcs

	func()
	func(*Context)
	func(http.ResponseWriter, *http.Request)
	func(http.ResponseWriter)
	func(*http.Request)

	it can has or has not return value
*/
func (router *router) addFunc(methods []string, url string, c interface{}) {
	vc := reflect.ValueOf(c)
	t := vc.Type()
	var rt RouteType

	if t.NumIn() == 0 {
		rt = FuncRoute
	} else if t.NumIn() == 1 {
		if t.In(0) == reflect.TypeOf(new(Context)) {
			rt = FuncCtxRoute
		} else if t.In(0) == reflect.TypeOf(new(http.Request)) {
			rt = FuncReqRoute
		} else if t.In(0).Kind() == reflect.Interface && t.In(0).Name() == "ResponseWriter" &&
			t.In(0).PkgPath() == "net/http" {
			rt = FuncResponseRoute
		} else {
			panic("no support function type")
		}
	} else if t.NumIn() == 2 &&
		(t.In(0).Kind() == reflect.Interface && t.In(0).Name() == "ResponseWriter" &&
			t.In(0).PkgPath() == "net/http") &&
		t.In(1) == reflect.TypeOf(new(http.Request)) {
		rt = FuncHttpRoute
	} else {
		panic("no support function type")
	}

	var r = NewRoute(t, vc, rt)
	url = removeStick(url)
	for _, m := range methods {
		router.addRoute(m, url, r)
	}
}

func (router *router) addStruct(methods map[string]string, url string, c interface{}) {
	vc := reflect.ValueOf(c)
	t := vc.Type().Elem()

	// added a default method Get, Post
	for name, method := range methods {
		if m, ok := t.MethodByName(method); ok {
			router.addRoute(name, removeStick(url), NewRoute(t, m.Func, StructPtrRoute))
		} else if m, ok := vc.Type().MethodByName(method); ok {
			router.addRoute(name, removeStick(url), NewRoute(t, m.Func, StructRoute))
		}
	}
}
