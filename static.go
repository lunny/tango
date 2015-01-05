package tango

import (
	"os"
	"path/filepath"
	"net/http"
	"strings"
)

func Static(rootPath, prefix string, indexFiles []string) HandlerFunc {
	return func(ctx *Context) {
		if ctx.Req().Method != "GET" && ctx.Req().Method != "HEAD" {
			ctx.Next()
			return
		}

		var rPath = ctx.Req().URL.Path
		// if defined prefix, then only check prefix
		if prefix != "" {
			if !strings.HasPrefix(ctx.Req().URL.Path, "/"+prefix) {
				ctx.Next()
				return
			} else {
				if len("/"+prefix) == len(ctx.Req().URL.Path) {
					rPath = ""
				} else {
					rPath = ctx.Req().URL.Path[len("/"+prefix):]
				}
			}
		}

		fPath, _ := filepath.Abs(filepath.Join(rootPath, rPath))
		finfo, err := os.Stat(fPath)
		if err != nil {
			if !os.IsNotExist(err) {
				ctx.WriteHeader(http.StatusInternalServerError)
				ctx.Write([]byte(err.Error()))
				return
			}
		} else if !finfo.IsDir() {
			err := ctx.ServeFile(fPath)
			if err != nil {
				ctx.WriteHeader(http.StatusInternalServerError)
				ctx.Write([]byte(err.Error()))
			}
			return
		} else {
			// try serving index.html or index.htm
			if len(indexFiles) > 0 {
				for _, index := range indexFiles {
					nPath := filepath.Join(fPath, index)
					finfo, err = os.Stat(nPath)
					if err != nil {
						if !os.IsNotExist(err) {
							ctx.WriteHeader(http.StatusInternalServerError)
							ctx.Write([]byte(err.Error()))
							return
						}
					} else if !finfo.IsDir() {
						err = ctx.ServeFile(nPath)
						if err != nil {
							ctx.WriteHeader(http.StatusInternalServerError)
							ctx.Write([]byte(err.Error()))
						}
						return
					}
				}
			}
		}

		ctx.Next()
	}
}