package tango

import (
	"reflect"
	"regexp"
	"strings"
	"net/url"
	"fmt"
)

type RouteType int

const (
	FuncRoute RouteType = iota + 1 	// func () string
	FuncHttpRoute 					// func (response, request) string
	FuncCtxRoute 					// func (*tango.Context) string
	StructRoute 					// func (st) Get() string
	StructPtrRoute 					// func (*struct) Get() string
)

type UrlType int

const (
	StaticUrl UrlType = iota + 1
	NamedUrl
	RegexpUrl
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
	urlType    UrlType
	method     reflect.Value
	routeType  RouteType
	pool       *pool
}

func NewRoute(r string, t reflect.Type,
	method reflect.Value, tp RouteType) *Route {
	var urlType UrlType = StaticUrl
	var cr *regexp.Regexp
	var err error
	if regexp.QuoteMeta(r) != r {
		urlType = RegexpUrl
		cr, err = regexp.Compile(r)
		if err != nil {
			panic("wrong route:"+err.Error())
			return nil
		}
	} else if strings.Contains(r, ":") {
		urlType = NamedUrl
	}

	return &Route{
		path: r,
		regexp: cr,
		urlType: urlType,
		method: method,
		routeType:  tp,
		pool: newPool(PoolSize, t),
	}
}

func (route *Route) Method() reflect.Value {
	return route.method
}

func (route *Route) UrlType() UrlType {
	return route.urlType
}

func (route *Route) RouteType() RouteType {
	return route.routeType
}

func (route *Route) IsStruct() bool {
	return route.routeType == StructRoute || route.routeType == StructPtrRoute
}

func (route *Route) newAction() reflect.Value {
	if !route.IsStruct() {
		return route.method
	}

	return route.pool.New()
}

func (route *Route) try(path string) (url.Values, bool) {
	p := make(url.Values)
	var i, j int
	for i < len(path) {
		switch {
		case j >= len(route.path):
			if route.path != "/" && len(route.path) > 0 && route.path[len(route.path)-1] == '/' {
				return p, true
			}
			return nil, false
		case route.path[j] == ':':
			var name, val string
			var nextc byte
			name, nextc, j = match(route.path, isAlnum, j+1)
			val, _, i = match(path, matchPart(nextc), i)
			p.Add(":"+name, val)
		case path[i] == route.path[j]:
			i++
			j++
		default:
			return nil, false
		}
	}
	if j != len(route.path) {
		return nil, false
	}
	return p, true
}

type Router interface {
	AddRouter(path string, methods []string, handler interface{})
	Match(requestPath, method string) (*Route, url.Values)
}

type router struct {
	routes      map[string][]*Route
	routesEq    map[string]map[string]*Route
	routesName  map[string][]*Route
}

func NewRouter() *router {
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

func Tail(pat, path string) string {
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

func (router *router) AddRouter(url string, methods []string, c interface{}) {
	vc := reflect.ValueOf(c)
	if vc.Kind() == reflect.Func {
		router.addFunc(methods, url, c)
	} else if vc.Kind() == reflect.Ptr && vc.Elem().Kind() == reflect.Struct {
		router.addStruct(methods, url, c)
	}
}

func (router *router) AddRoute(m string, route *Route) {
	switch route.urlType {
	case StaticUrl:
		router.routesEq[m][route.path] = route
	case NamedUrl:
		router.routesName[m] = append(router.routesName[m], route)
	case RegexpUrl:
		router.routes[m] = append(router.routes[m], route)
	}
}

func (router *router) addFunc(methods []string, url string, c interface{}) {
	vc := reflect.ValueOf(c)
	t := vc.Type()
	r := NewRoute(removeStick(url), t, vc, FuncRoute)
	for _, m := range methods {
		router.AddRoute(m, r)
	}
}

func (router *router) addStruct(methods []string, url string, c interface{}) {
	vc := reflect.ValueOf(c)
	t := vc.Type().Elem()

	// added a default method Get, Post
	for _, name := range methods {
		newName := strings.Title(strings.ToLower(name))
		if m, ok := t.MethodByName(newName); ok {
			router.AddRoute(name, NewRoute(removeStick(url), t, m.Func, StructPtrRoute))
		} else if m, ok := vc.Type().MethodByName(newName); ok {
			router.AddRoute(name, NewRoute(removeStick(url), t, m.Func, StructRoute))
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
