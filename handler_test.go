package tango

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGroupGroup(t *testing.T) {
	var a int
	tg := Classic()
	tg.Group("/api", func(api *Group) {
		api.Group("/v1", func(v1 *Group) {
			v1.Use(HandlerFunc(func(ctx *Context) {
				a = 1
				ctx.Next()
			}))

			v1.Get("/", func(ctx *Context) {
				fmt.Println("context")
			})

		})
	})

	buff := bytes.NewBufferString("")
	recorder := httptest.NewRecorder()
	recorder.Body = buff

	req, err := http.NewRequest("GET", "http://localhost:8000/api/v1", nil)
	if err != nil {
		t.Error(err)
	}
	tg.ServeHTTP(recorder, req)
	expect(t, recorder.Code, http.StatusOK)
	expect(t, buff.String(), "")
	expect(t, 1, a)

}
