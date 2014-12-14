package tango

import (
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

/* Test Helpers */
func expect(t *testing.T, a interface{}, b interface{}) {
	if a != b {
		t.Errorf("Expected %v (type %v) - Got %v (type %v)", b, reflect.TypeOf(b), a, reflect.TypeOf(a))
	}
}

func refute(t *testing.T, a interface{}, b interface{}) {
	if a == b {
		t.Errorf("Did not expect %v (type %v) - Got %v (type %v)", b, reflect.TypeOf(b), a, reflect.TypeOf(a))
	}
}

func TestTangoRun(t *testing.T) {
	// just test that Run doesn't bomb
	go New().Run(":8000")
}

func TestTangoServeHTTP(t *testing.T) {
	result := ""
	response := httptest.NewRecorder()

	n := New()
	n.Use(HandlerFunc(func(ctx *Context) {
		result += "foo"
		ctx.Next()
		result += "ban"
	}))
	n.Use(HandlerFunc(func(ctx *Context) {
		result += "bar"
		ctx.Next()
		result += "baz"
	}))
	n.Use(HandlerFunc(func(ctx *Context) {
		result += "bat"
		ctx.WriteHeader(http.StatusBadRequest)
	}))

	n.ServeHTTP(response, (*http.Request)(nil))

	expect(t, result, "foobarbatbazban")
	expect(t, response.Code, http.StatusBadRequest)
}