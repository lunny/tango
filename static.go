// Copyright 2015 The Tango Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package tango

import (
	"os"
	"path"
	"path/filepath"
	"strings"
)

// StaticOptions defines Static middleware's options
type StaticOptions struct {
	RootPath   string
	Prefix     string
	IndexFiles []string
	ListDir    bool
	FilterExts []string
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

	if len(opt.Prefix) > 0 {
		if opt.Prefix[0] != '/' {
			opt.Prefix = "/" + opt.Prefix
		}
	}

	if len(opt.IndexFiles) == 0 {
		opt.IndexFiles = []string{"index.html", "index.htm"}
	}

	return opt
}

// Static return a middleware for serving static files
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
			if !strings.HasPrefix(ctx.Req().URL.Path, opt.Prefix) {
				ctx.Next()
				return
			}

			if len(opt.Prefix) == len(ctx.Req().URL.Path) {
				rPath = ""
			} else {
				rPath = ctx.Req().URL.Path[len(opt.Prefix):]
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
			if len(opt.FilterExts) > 0 {
				var matched bool
				for _, ext := range opt.FilterExts {
					if filepath.Ext(fPath) == ext {
						matched = true
						break
					}
				}
				if !matched {
					ctx.Next()
					return
				}
			}

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

			// list dir files
			if opt.ListDir {
				ctx.Header().Set("Content-Type", "text/html; charset=UTF-8")
				ctx.WriteString(`<ul style="list-style-type:none;line-height:32px;">`)
				rootPath, _ := filepath.Abs(opt.RootPath)
				rPath, _ := filepath.Rel(rootPath, fPath)
				if fPath != rootPath {
					ctx.WriteString(`<li>&nbsp; &nbsp; <a href="` + path.Join("/", opt.Prefix, filepath.Dir(rPath)) + `">..</a></li>`)
				}
				err = filepath.Walk(fPath, func(p string, fi os.FileInfo, err error) error {
					rPath, _ := filepath.Rel(fPath, p)
					if rPath == "." || len(strings.Split(rPath, string(filepath.Separator))) > 1 {
						return nil
					}
					rPath, _ = filepath.Rel(rootPath, p)
					ps, _ := os.Stat(p)
					if ps.IsDir() {
						ctx.WriteString(`<li>┖ <a href="` + path.Join("/", opt.Prefix, rPath) + `">` + filepath.Base(p) + `</a></li>`)
					} else {
						if len(opt.FilterExts) > 0 {
							var matched bool
							for _, ext := range opt.FilterExts {
								if filepath.Ext(p) == ext {
									matched = true
									break
								}
							}
							if !matched {
								return nil
							}
						}

						ctx.WriteString(`<li>&nbsp; &nbsp; <a href="` + path.Join("/", opt.Prefix, rPath) + `">` + filepath.Base(p) + `</a></li>`)
					}
					return nil
				})
				ctx.WriteString("</ul>")
				return
			}
		}

		ctx.Next()
	}
}
