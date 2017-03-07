package tango

import (
	"bytes"
	"fmt"
	"net/http"
	"reflect"
	"regexp"
	"sort"
	"strings"
)

// RouteType defines route types
type RouteType byte

// enumerates route types
const (
	FuncRoute         RouteType = iota + 1 // 1 func ()
	FuncHTTPRoute                          // 2 func (http.ResponseWriter, *http.Request)
	FuncReqRoute                           // 3 func (*http.Request)
	FuncResponseRoute                      // 4 func (http.ResponseWriter)
	FuncCtxRoute                           // 5 func (*tango.Context)
	StructRoute                            // 6 func (st) <Get>()
	StructPtrRoute                         // 7 func (*struct) <Get>()
)

// enumerates all supported HTTP methods
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

	PoolSize = 10
)

// Route defines HTTP route
type Route struct {
	raw       interface{}
	method    reflect.Value
	handlers  []Handler
	routeType RouteType
	pool      *pool
}

// NewRoute returns a route
func NewRoute(v interface{}, t reflect.Type,
	method reflect.Value, tp RouteType, handlers []Handler) *Route {
	var pool *pool
	if tp == StructRoute || tp == StructPtrRoute {
		pool = newPool(PoolSize, t)
	}
	return &Route{
		raw:       v,
		routeType: tp,
		method:    method,
		pool:      pool,
		handlers:  handlers,
	}
}

// Raw returns raw data to define route.
func (r *Route) Raw() interface{} {
	return r.raw
}

// Method returns finalize execute method.
func (r *Route) Method() reflect.Value {
	return r.method
}

// RouteType returns route type.
func (r *Route) RouteType() RouteType {
	return r.routeType
}

// IsStruct returns if the execute is a struct
func (r *Route) IsStruct() bool {
	return r.routeType == StructRoute || r.routeType == StructPtrRoute
}

func (r *Route) newAction() reflect.Value {
	if !r.IsStruct() {
		return r.method
	}

	return r.pool.New()
}

// Router describes the interface of route
type Router interface {
	Route(methods interface{}, path string, handler interface{}, middlewares ...Handler)
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
		path    string         // executor path
	}
	edges []*node
)

func (e edges) Len() int { return len(e) }

func (e edges) Swap(i, j int) { e[i], e[j] = e[j], e[i] }

// static route will be put the first, so it will be match first.
// two static route, content longer is first.
func (e edges) Less(i, j int) bool {
	if e[i].tp == snode {
		if e[j].tp == snode {
			return len(e[i].content) > len(e[j].content)
		}
		return true
	}
	if e[j].tp == snode {
		return false
	}
	return i < j
}

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

// newRouter return a new router
func newRouter() (r *router) {
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
			bracket = 0
			if i == l {
				return nodes
			}
		} else if path[i] == '*' {
			nodes = append(nodes, &node{tp: snode, content: path[j : i-bracket]})
			j = i
			if bracket == 1 {
				for ; i < l && ')' != path[i]; i++ {
				}
			} else {
				i = i + 1
				for ; i < l && isAlnum(path[i]); i++ {
				}
			}
			nodes = append(nodes, &node{tp: anode, content: path[j:i]})
			i = i + bracket
			bracket = 0
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
			fmt.Print("  ")
		}
		if i > 1 {
			fmt.Print("┗", "  ")
		}

		fmt.Print(c.content)
		if c.handle != nil {
			fmt.Print("  ", c.handle.method.Type())
			fmt.Printf("  %p", c.handle.method.Interface())
		}
		fmt.Println()
		printNode(i+1, c)
	}
}

func (r *router) printTrees() {
	for _, method := range SupportMethods {
		if len(r.trees[method].edges) > 0 {
			fmt.Println(method)
			printNode(1, r.trees[method])
			fmt.Println()
		}
	}
}

func (r *router) addRoute(method, path string, h *Route) {
	nodes := parseNodes(path)
	nodes[len(nodes)-1].handle = h
	nodes[len(nodes)-1].path = path
	if !validNodes(nodes) {
		panic(fmt.Sprintln("express", path, "is not supported"))
	}
	r.addnodes(method, nodes)
	//r.printTrees()
}

func (r *router) matchNode(n *node, url string, params Params) (*node, Params) {
	if n.tp == snode {
		if strings.HasPrefix(url, n.content) {
			if len(url) == len(n.content) {
				return n, params
			}
			for _, c := range n.edges {
				e, newParams := r.matchNode(c, url[len(n.content):], params)
				if e != nil {
					return e, newParams
				}
			}
		}
	} else if n.tp == anode {
		for _, c := range n.edges {
			idx := strings.LastIndex(url, c.content)
			if idx > -1 {
				params = append(params, param{n.content, url[:idx]})
				return r.matchNode(c, url[idx:], params)
			}
		}
		return n, append(params, param{n.content, url})
	} else if n.tp == nnode {
		for _, c := range n.edges {
			idx := strings.Index(url, c.content)
			if idx > -1 {
				params = append(params, param{n.content, url[:idx]})
				return r.matchNode(c, url[idx:], params)
			}
		}
		idx := strings.IndexByte(url, '/')
		if idx < 0 {
			params = append(params, param{n.content, url})
			return n, params
		}
	} else if n.tp == rnode {
		idx := strings.IndexByte(url, '/')
		if idx > -1 {
			if n.regexp.MatchString(url[:idx]) {
				for _, c := range n.edges {
					h, newParams := r.matchNode(c, url[idx:], params)
					if h != nil {
						return h, append([]param{param{n.content, url[:idx]}}, newParams...)
					}
				}
			}
			return nil, params
		}

		for _, c := range n.edges {
			idx := strings.Index(url, c.content)
			if idx > -1 && n.regexp.MatchString(url[:idx]) {
				params = append(params, param{n.content, url[:idx]})
				return r.matchNode(c, url[idx:], params)
			}
		}

		if n.regexp.MatchString(url) {
			params = append(params, param{n.content, url})
			return n, params
		}
	}
	return nil, params
}

// Match for request url, match router
func (r *router) Match(url, method string) (*Route, Params) {
	cn, ok := r.trees[method]
	if !ok {
		return nil, nil
	}
	var params = make(Params, 0, strings.Count(url, "/"))
	for _, n := range cn.edges {
		e, newParams := r.matchNode(n, url, params)
		if e != nil {
			return e.handle, newParams
		}
	}
	return nil, nil
}

// addnode adds node nodes[i] to parent node p
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
	sort.Sort(p.edges)
	return nodes[i]
}

// validNodes validates parsed nodes, all non-static route should have static route children.
func validNodes(nodes []*node) bool {
	if len(nodes) == 0 {
		return false
	}
	var lastTp = nodes[0]
	for _, node := range nodes[1:] {
		if lastTp.tp != snode && node.tp != snode {
			return false
		}
		lastTp = node
	}
	return true
}

// addnodes adds nodes to trees
func (r *router) addnodes(method string, nodes []*node) {
	cn := r.trees[method]
	var p = cn
	for i := 0; i < len(nodes); i++ {
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

// Route adds route
func (r *router) Route(ms interface{}, url string, c interface{}, handlers ...Handler) {
	vc := reflect.ValueOf(c)
	if vc.Kind() == reflect.Func {
		switch ms.(type) {
		case string:
			s := strings.Split(ms.(string), ":")
			r.addFunc([]string{s[0]}, url, c, handlers)
		case []string:
			var newSlice []string
			for _, m := range ms.([]string) {
				s := strings.Split(m, ":")
				newSlice = append(newSlice, s[0])
			}
			r.addFunc(newSlice, url, c, handlers)
		default:
			panic("unknow methods format")
		}
	} else if vc.Kind() == reflect.Ptr && vc.Elem().Kind() == reflect.Struct {
		if handler, ok := vc.Interface().(http.Handler); ok {
			r.Route(ms, url, handler.ServeHTTP, handlers...)
			return
		}
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

		r.addStruct(methods, url, c, handlers)
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
func (r *router) addFunc(methods []string, url string, c interface{}, handlers []Handler) {
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
			panic(fmt.Sprintln("no support function type", methods, url, c))
		}
	} else if t.NumIn() == 2 &&
		(t.In(0).Kind() == reflect.Interface && t.In(0).Name() == "ResponseWriter" &&
			t.In(0).PkgPath() == "net/http") &&
		t.In(1) == reflect.TypeOf(new(http.Request)) {
		rt = FuncHTTPRoute
	} else {
		panic(fmt.Sprintln("no support function type", methods, url, c))
	}

	url = removeStick(url)
	for _, m := range methods {
		r.addRoute(m, url, NewRoute(c, t, vc, rt, handlers))
	}
}

func (r *router) addStruct(methods map[string]string, url string, c interface{}, handlers []Handler) {
	vc := reflect.ValueOf(c)
	t := vc.Type().Elem()

	// added a default method Get, Post
	for name, method := range methods {
		if m, ok := t.MethodByName(method); ok {
			r.addRoute(name, removeStick(url), NewRoute(c, t, m.Func, StructPtrRoute, handlers))
		} else if m, ok := vc.Type().MethodByName(method); ok {
			r.addRoute(name, removeStick(url), NewRoute(c, t, m.Func, StructRoute, handlers))
		} else if m, ok := t.MethodByName("Any"); ok {
			r.addRoute(name, removeStick(url), NewRoute(c, t, m.Func, StructPtrRoute, handlers))
		} else if m, ok := vc.Type().MethodByName("Any"); ok {
			r.addRoute(name, removeStick(url), NewRoute(c, t, m.Func, StructRoute, handlers))
		}
	}
}
