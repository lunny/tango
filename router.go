package tango

import (
	"reflect"
	"regexp"
	"strings"
)

type RouteType int

const (
	FuncRoute = iota + 1
	StructRoute
	StructPtrRoute
)

func removeStick(uri string) string {
	uri = strings.TrimRight(uri, "/")
	if uri == "" {
		uri = "/"
	}
	return uri
}

// Route
type Route struct {
	path       string          //path string
	regexp     *regexp.Regexp  //path regexp
	methods    map[string]bool //GET POST HEAD DELETE etc.
	structType reflect.Type    //handler element
	method     reflect.Value
	routeType  RouteType
	pools      *pools
}

func (route *Route) Method() reflect.Value {
	return route.method
}

func (route *Route) StructType() reflect.Type {
	return route.structType
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

	return route.pools.New(route.structType)
}

type Router interface {
	AddRouter(path string, methods []string, handler interface{})
	Match(requestPath, method string) (*Route, []reflect.Value)
}

type router struct {
	routes      []*Route
	routesEq    map[string]map[string]*Route
	defaultFunc string
	pools       *pools
}

func NewRouter() *router {
	return &router{
		routes:      make([]*Route, 0),
		routesEq:    make(map[string]map[string]*Route),
		defaultFunc: "Do",
		pools:       NewPools(800),
	}
}

func (router *router) addRoute(r string, methods map[string]bool,
	t reflect.Type, handler string,
	method reflect.Value, tp RouteType) error {
	cr, err := regexp.Compile(r)
	if err != nil {
		return err
	}
	router.routes = append(router.routes, &Route{
		path:       r,
		regexp:     cr,
		methods:    methods,
		structType: t,
		method:     method,
		routeType:  tp,
		pools:      router.pools,
	})
	return nil
}

func (router *router) addEqRoute(r string, methods map[string]bool,
	t reflect.Type, handler string,
	method reflect.Value, tp RouteType) {
	if _, ok := router.routesEq[r]; !ok {
		router.routesEq[r] = make(map[string]*Route)
	}
	for v, _ := range methods {
		router.routesEq[r][v] = &Route{
			structType: t,
			method:     method,
			routeType:  tp,
			pools:      router.pools,
		}
	}
}

var (
	defaultMethods = []string{
		"GET",
		"POST",
		"HEAD",
		"DELETE",
		"PUT",
		"OPTIONS",
		"TRACE",
		"PATCH",
	}
)

func slice2map(slice []string) map[string]bool {
	var res = make(map[string]bool, len(slice))
	for _, s := range slice {
		res[s] = true
	}
	return res
}

func (router *router) AddRouter(url string, methods []string, c interface{}) {
	vc := reflect.ValueOf(c)
	methodsMap := slice2map(methods)
	if vc.Kind() == reflect.Func {
		router.addFuncRouter(methodsMap, url, c)
	} else if vc.Kind() == reflect.Ptr && vc.Elem().Kind() == reflect.Struct {
		router.addStructRouter(methodsMap, url, c)
	}
}

func (router *router) addFuncRouter(methods map[string]bool, url string, c interface{}) {
	vc := reflect.ValueOf(c)
	t := vc.Type()
	router.addEqRoute(removeStick(url), methods, t, "", vc, FuncRoute)
}

func (router *router) addStructRouter(methods map[string]bool, url string, c interface{}) {
	vc := reflect.ValueOf(c)
	t := vc.Type().Elem()

	// added a default method Do as /
	if m, ok := t.MethodByName(router.defaultFunc); ok {
		router.addEqRoute(
			removeStick(url), methods,
			t, router.defaultFunc,
			m.Func, StructPtrRoute,
		)
	} else if m, ok := vc.Type().MethodByName(router.defaultFunc); ok {
		router.addEqRoute(
			removeStick(url), methods,
			t, router.defaultFunc,
			m.Func, StructRoute,
		)
	}
}

// when a request ask, then match the correct route
func (router *router) Match(reqPath, allowMethod string) (*Route, []reflect.Value) {
	var route *Route
	var args = make([]reflect.Value, 0)

	// for non-regular path, search the map
	if routes, ok := router.routesEq[reqPath]; ok {
		if route, ok = routes[allowMethod]; ok {
			return route, args
		}
	}

	for _, r := range router.routes {
		//if the methods don't match, skip this handler (except HEAD can be used in place of GET)
		if _, ok := r.methods[allowMethod]; !ok {
			continue
		}

		if !r.regexp.MatchString(reqPath) {
			continue
		}

		match := r.regexp.FindStringSubmatch(reqPath)
		if len(match[0]) != len(reqPath) {
			continue
		}

		for _, arg := range match[1:] {
			args = append(args, reflect.ValueOf(arg))
		}

		return route, args
	}

	return nil, nil
}
