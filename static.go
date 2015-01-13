package tango

import (
	"os"
	"path/filepath"
	"strings"
)

type StaticOptions struct {
	RootPath string
	Prefix string
	IndexFiles []string
}

func prepareStaticOptions(options []StaticOptions) StaticOptions {
	var opt StaticOptions
	if len(options) > 0 {
		opt = options[0]
	}

	// Defaults
	if len(opt.RootPath) == 0 {
		opt.RootPath = "./public"
	}
	if len(opt.Prefix) == 0 {
		opt.Prefix = ""
	}
	if len(opt.IndexFiles) == 0 {
		opt.IndexFiles = []string{"index.html", "index.htm"}
	}

	return opt
}

func Static(opts ...StaticOptions) HandlerFunc {
	return func(ctx *Context) {
		if ctx.Req().Method != "GET" && ctx.Req().Method != "HEAD" {
			ctx.Next()
			return
		}

		opt := prepareStaticOptions(opts)

		var rPath = ctx.Req().URL.Path
		// if defined prefix, then only check prefix
		if opt.Prefix != "" {
			if !strings.HasPrefix(ctx.Req().URL.Path, "/"+opt.Prefix) {
				ctx.Next()
				return
			} else {
				if len("/"+opt.Prefix) == len(ctx.Req().URL.Path) {
					rPath = ""
				} else {
					rPath = ctx.Req().URL.Path[len("/"+opt.Prefix):]
				}
			}
		}

		fPath, _ := filepath.Abs(filepath.Join(opt.RootPath, rPath))
		finfo, err := os.Stat(fPath)
		if err != nil {
			if !os.IsNotExist(err) {
				ctx.Result = InternalServerError(err.Error())
				ctx.HandleError()
				return
			}
		} else if !finfo.IsDir() {
			err := ctx.ServeFile(fPath)
			if err != nil {
				ctx.Result = InternalServerError(err.Error())
				ctx.HandleError()
			}
			return
		} else {
			// try serving index.html or index.htm
			if len(opt.IndexFiles) > 0 {
				for _, index := range opt.IndexFiles {
					nPath := filepath.Join(fPath, index)
					finfo, err = os.Stat(nPath)
					if err != nil {
						if !os.IsNotExist(err) {
							ctx.Result = InternalServerError(err.Error())
							ctx.HandleError()
							return
						}
					} else if !finfo.IsDir() {
						err = ctx.ServeFile(nPath)
						if err != nil {
							ctx.Result = InternalServerError(err.Error())
							ctx.HandleError()
						}
						return
					}
				}
			}
		}

		ctx.Next()
	}
}