package tango

import (
	"fmt"
	"net/http"
	"runtime"
)

func Recovery(debug bool) HandlerFunc {
	return func(ctx *Context) {
		defer func() {
			if e := recover(); e != nil {
				content := fmt.Sprintf("Handler crashed with error: %v", e)
				for i := 1; ;i += 1 {
					_, file, line, ok := runtime.Caller(i)
					if !ok {
						break
					} else {
						content += "\n"
					}
					content += fmt.Sprintf("%v %v", file, line)
				}

				ctx.Logger.Error(content)

				if !ctx.Written() {
					ctx.WriteHeader(http.StatusInternalServerError)
					if !debug {
						content = http.StatusText(http.StatusInternalServerError)
					}
					ctx.Write([]byte(content))
				}
			}
		}()

		ctx.Next()
	}
}
