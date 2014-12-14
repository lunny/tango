package tango

import (
	"reflect"
	"regexp"
)

type RouteType int

const (
	FuncRoute = iota +1
	StructRoute
	StructPtrRoute
)

// Route
type Route struct {
	path           string          //path string
	regexp *regexp.Regexp  //path regexp
	methods    map[string]bool //GET POST HEAD DELETE etc.
	structType 	reflect.Type    //handler element
	method         reflect.Value
	routeType       RouteType
	pools *pools
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
	Any(string, interface{})
	Get(string, interface{})
	Post(string, interface{})
	Head(string, interface{})
	Match(requestPath, method string) (*Route, []reflect.Value)
}

type router struct {
	routes          []*Route
	routesEq        map[string]map[string]*Route
	defaultFunc string
	pools *pools
}

func NewRouter() *router {
	return &router{
		routes:          make([]*Route, 0),
		routesEq:        make(map[string]map[string]*Route),
		defaultFunc: "Do",
		pools: NewPools(2000),
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
		path:           r,
		regexp: cr,
		methods:    methods,
		structType: t,
		method:         method,
		routeType:       tp,
		pools: 	router.pools,
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
			method:         method,
			routeType:       tp,
			pools: 	router.pools,
		}
	}
}

var (
	defaultMethods = map[string]bool{
		"GET": true, 
		"POST": true,
		"HEAD": true,
	}
)

func (router *router) addRouter(methods map[string]bool, url string, c interface{}) {
	vc := reflect.ValueOf(c)
	if vc.Kind() == reflect.Func {
		router.addFuncRouter(methods, url, c)
	} else if vc.Kind() == reflect.Ptr && vc.Elem().Kind() == reflect.Struct {
		router.addStructRouter(methods, url, c)
	}
}

func (router *router) Get(url string, c interface{}) {
	methods := map[string]bool{"GET": true}
	router.addRouter(methods, url, c)
}

func (router *router) Post(url string, c interface{}) {
	methods := map[string]bool{"POST": true}
	router.addRouter(methods, url, c)
}

func (router *router) Head(url string, c interface{}) {
	methods := map[string]bool{"Head": true}
	router.addRouter(methods, url, c)
}

func (router *router) Any(url string, c interface{}) {
	router.addRouter(defaultMethods, url, c)
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