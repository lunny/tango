package tango

import (
	"bytes"
	"reflect"
	"regexp"
	"strings"
	"net/url"
	"fmt"
	"net/http"
)

type RouteType int

const (
	FuncRoute RouteType = iota + 1 	// func ()
	FuncHttpRoute 					// func (http.ResponseWriter, *http.Request)
	FuncReqRoute 					// func (*http.Request)
	FuncResponseRoute 				// func (http.ResponseWriter)
	FuncCtxRoute 					// func (*tango.Context)
	StructRoute 					// func (st) Get()
	StructPtrRoute 					// func (*struct) Get()
)

type PathType int

const (
	StaticPath PathType = iota + 1
	NamedPath
	RegexpPath
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
	path       string          //path string
	regexp     *regexp.Regexp  //path regexp
	pathType   PathType
	method     reflect.Value
	routeType  RouteType
	pool       *pool
}

var specialBytes = []byte(`\.+*?()|[]{}^$`)

func pathType(s string) PathType {
	for i := 0; i < len(s); i++ {
		if s[i] == ':'{
			return NamedPath
		}
		if bytes.IndexByte(specialBytes, s[i]) >= 0 {
			return RegexpPath
		}
	}
	return StaticPath
}

func NewRoute(r string, t reflect.Type,
	method reflect.Value, tp RouteType) *Route {
	var cr *regexp.Regexp
	var err error
	var pathType = pathType(r)
	if pathType == RegexpPath {
		cr, err = regexp.Compile(r)
		if err != nil {
			panic("wrong route:"+err.Error())
			return nil
		}
	}

	var pool *pool
	if tp == StructRoute || tp == StructPtrRoute {
		pool = newPool(PoolSize, t)
	}
	return &Route{
		path: r,
		regexp: cr,
		pathType: pathType,
		method: method,
		routeType:  tp,
		pool: pool,
	}
}

func (r *Route) Method() reflect.Value {
	return r.method
}

func (r *Route) PathType() PathType {
	return r.pathType
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

func (r *Route) try(path string) (url.Values, bool) {
	p := make(url.Values)
	var i, j int
	for i < len(path) {
		switch {
		case j >= len(r.path):
			if r.path != "/" && len(r.path) > 0 && r.path[len(r.path)-1] == '/' {
				return p, true
			}
			return nil, false
		case r.path[j] == ':':
			var name, val string
			var nextc byte
			name, nextc, j = match(r.path, isAlnum, j+1)
			val, _, i = match(path, matchPart(nextc), i)
			p.Add(":"+name, val)
		case path[i] == r.path[j]:
			i++
			j++
		default:
			return nil, false
		}
	}
	if j != len(r.path) {
		return nil, false
	}
	return p, true
}

type Router interface {
	Route(methods []string, path string, handler interface{})
	Match(requestPath, method string) (*Route, url.Values)
}

type router struct {
	routes      map[string][]*Route
	routesEq    map[string]map[string]*Route
	routesName  map[string][]*Route
}

func NewRouter() Router {
	routesEq := make(map[string]map[string]*Route)
	for _, m := range SupportMethods {
		routesEq[m] = make(map[string]*Route)
	}

	routesName := make(map[string][]*Route)
	for _, m := range SupportMethods {
		routesName[m] = make([]*Route, 0)
	}

	routes := make(map[string][]*Route)
	for _, m := range SupportMethods {
		routes[m] = make([]*Route, 0)
	}

	return &router{
		routesEq: routesEq,
		routes: routes,
		routesName: routesName,
	}
}

func matchPart(b byte) func(byte) bool {
	return func(c byte) bool {
		return c != b && c != '/'
	}
}

func match(s string, f func(byte) bool, i int) (matched string, next byte, j int) {
	j = i
	for j < len(s) && f(s[j]) {
		j++
	}
	if j < len(s) {
		next = s[j]
	}
	return s[i:j], next, j
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

func tail(pat, path string) string {
	var i, j int
	for i < len(path) {
		switch {
		case j >= len(pat):
			if pat[len(pat)-1] == '/' {
				return path[i:]
			}
			return ""
		case pat[j] == ':':
			var nextc byte
			_, nextc, j = match(pat, isAlnum, j+1)
			_, _, i = match(path, matchPart(nextc), i)
		case path[i] == pat[j]:
			i++
			j++
		default:
			return ""
		}
	}
	return ""
}

func removeStick(uri string) string {
	uri = strings.TrimRight(uri, "/")
	if uri == "" {
		uri = "/"
	}
	return uri
}

func (router *router) Route(methods []string, url string, c interface{}) {
	vc := reflect.ValueOf(c)
	if vc.Kind() == reflect.Func {
		router.addFunc(methods, url, c)
	} else if vc.Kind() == reflect.Ptr && vc.Elem().Kind() == reflect.Struct {
		router.addStruct(methods, url, c)
	} else {
		panic("not support route type")
	}
}

func (router *router) addRoute(m string, route *Route) {
	switch route.pathType {
	case StaticPath:
		router.routesEq[m][route.path] = route
	case NamedPath:
		router.routesName[m] = append(router.routesName[m], route)
	case RegexpPath:
		router.routes[m] = append(router.routes[m], route)
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
	var r *Route

	if t.NumIn() == 0 {
		r = NewRoute(removeStick(url), t, vc, FuncRoute)
	} else if t.NumIn() == 1 {
		if t.In(0) == reflect.TypeOf(new(Context)) {
			r = NewRoute(removeStick(url), t, vc, FuncCtxRoute)
		} else if t.In(0) == reflect.TypeOf(new(http.Request)) {
			r = NewRoute(removeStick(url), t, vc, FuncReqRoute)
		} else if t.In(0).Kind() == reflect.Interface && t.In(0).Name() == "ResponseWriter" && 
			t.In(0).PkgPath() == "net/http" {
			r = NewRoute(removeStick(url), t, vc, FuncResponseRoute)
		} else {
			panic("no support function type")
		}
	} else if t.NumIn() == 2 && 
		(t.In(0).Kind() == reflect.Interface && t.In(0).Name() == "ResponseWriter" && 
			t.In(0).PkgPath() == "net/http") && 
		t.In(1) == reflect.TypeOf(new(http.Request)) {
		r = NewRoute(removeStick(url), t, vc, FuncHttpRoute)
	} else {
		panic("no support function type")
	}
	for _, m := range methods {
		router.addRoute(m, r)
	}
}

func (router *router) addStruct(methods []string, url string, c interface{}) {
	vc := reflect.ValueOf(c)
	t := vc.Type().Elem()

	// added a default method Get, Post
	for _, name := range methods {
		newName := strings.Title(strings.ToLower(name))
		if m, ok := t.MethodByName(newName); ok {
			router.addRoute(name, NewRoute(removeStick(url), t, m.Func, StructPtrRoute))
		} else if m, ok := vc.Type().MethodByName(newName); ok {
			router.addRoute(name, NewRoute(removeStick(url), t, m.Func, StructRoute))
		}
	}
}

// when a request ask, then match the correct route
func (router *router) Match(reqPath, allowMethod string) (*Route, url.Values) {
	// for non-regular path, search the map
	if routes, ok := router.routesEq[allowMethod]; ok {
		if route, ok := routes[reqPath]; ok {
			return route, make(url.Values)
		}
	}

	// name match
	routes := router.routesName[allowMethod]
	for _, r := range routes {
		if args, ok := r.try(reqPath); ok {
			return r, args
		}
	}

	// regex match
	routes = router.routes[allowMethod]
	for _, r := range routes {
		if !r.regexp.MatchString(reqPath) {
			continue
		}

		match := r.regexp.FindStringSubmatch(reqPath)
		if len(match[0]) != len(reqPath) {
			continue
		}

		var args = make(url.Values)
		// for regexp :0 -> first match param :1 -> the second
		for i, arg := range match[1:] {
			args.Add(fmt.Sprintf(":%d", i), arg)
		}

		return r, args
	}

	return nil, nil
}
