package tango

import (
	"os"
	"path/filepath"
	"net/http"
	"strings"
)

type Statics struct {
	Prefix string
	RootPath   string
	IndexFiles []string
}

func NewStatic(rootPath, prefix string, indexFiles []string) *Statics {
	return &Statics{
		prefix,
		rootPath,
		indexFiles,
	}
}

func (s *Statics) Handle(ctx *Context) {
	if ctx.Req().Method != "GET" && ctx.Req().Method != "HEAD" {
		ctx.Next()
		return
	}

	var rPath = ctx.Req().URL.Path
	// if defined prefix, then only check prefix
	if s.Prefix != "" {
		if !strings.HasPrefix(ctx.Req().URL.Path, "/"+s.Prefix) {
			ctx.Next()
			return
		} else {
			if len("/"+s.Prefix) == len(ctx.Req().URL.Path) {
				rPath = ""
			} else {
				rPath = ctx.Req().URL.Path[len("/"+s.Prefix):]
			}
		}
	}

	fPath, _ := filepath.Abs(filepath.Join(s.RootPath, rPath))
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
		if len(s.IndexFiles) > 0 {
			for _, index := range s.IndexFiles {
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