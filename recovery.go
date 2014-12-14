package tango

import (
	"fmt"
	"log"
	"net/http"
	"runtime"
)

type Recovery struct {
	debug        bool
	logger       Logger
}

func (recovery *Recovery) SetLogger(logger Logger) {
	recovery.logger = logger
}

func NewRecovery(debug bool) *Recovery {
	return &Recovery{debug: debug}
}

func (recovery *Recovery) Handle(ctx *Context) {
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

			if recovery.logger != nil {
				recovery.logger.Error(content)
			} else {
				log.Println("[tango]", content)
			}

			if !ctx.Written() {
				ctx.WriteHeader(http.StatusInternalServerError)
				if !recovery.debug {
					content = statusText[http.StatusInternalServerError]
				}
				ctx.Write([]byte(content))
			}
		}
	}()

	ctx.Next()
}
