// Copyright 2015 The Tango Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package tango

import (
	"net/http"
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
	// FileSystem is the interface for supporting any implmentation of file system.
	FileSystem http.FileSystem
}

// IsFilterExt decribes if rPath's ext match filter ext
func (s *StaticOptions) IsFilterExt(rPath string) bool {
	rext := path.Ext(rPath)
	for _, ext := range s.FilterExts {
		if rext == ext {
			return true
		}
	}
	return false
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

	if opt.FileSystem == nil {
		ps, _ := filepath.Abs(opt.RootPath)
		opt.FileSystem = http.Dir(ps)
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

		f, err := opt.FileSystem.Open(strings.TrimLeft(rPath, "/"))
		if err != nil {
			if os.IsNotExist(err) {
				if opt.Prefix != "" {
					ctx.Result = NotFound()
				} else {
					ctx.Next()
					return
				}
			} else {
				ctx.Result = InternalServerError(err.Error())
			}
			ctx.HandleError()
			return
		}
		defer f.Close()

		finfo, err := f.Stat()
		if err != nil {
			ctx.Result = InternalServerError(err.Error())
			ctx.HandleError()
			return
		}

		if !finfo.IsDir() {
			if len(opt.FilterExts) > 0 && !opt.IsFilterExt(rPath) {
				ctx.Next()
				return
			}

			http.ServeContent(ctx, ctx.Req(), finfo.Name(), finfo.ModTime(), f)
			return
		}

		// try serving index.html or index.htm
		if len(opt.IndexFiles) > 0 {
			for _, index := range opt.IndexFiles {
				fi, err := opt.FileSystem.Open(strings.TrimLeft(path.Join(rPath, index), "/"))
				if err != nil {
					if !os.IsNotExist(err) {
						ctx.Result = InternalServerError(err.Error())
						ctx.HandleError()
						return
					}
				} else {
					finfo, err = fi.Stat()
					if err != nil {
						fi.Close()
						ctx.Result = InternalServerError(err.Error())
						ctx.HandleError()
						return
					}
					if !finfo.IsDir() {
						http.ServeContent(ctx, ctx.Req(), finfo.Name(), finfo.ModTime(), fi)
						fi.Close()
						return
					}
					fi.Close()
				}
			}
		}

		// list dir files
		if opt.ListDir {
			ctx.Header().Set("Content-Type", "text/html; charset=UTF-8")
			ctx.WriteString(`<ul style="list-style-type:none;line-height:32px;">`)
			if rPath != "/" {
				ctx.WriteString(`<li>&nbsp; &nbsp; <a href="` + path.Join("/", opt.Prefix, filepath.Dir(rPath)) + `">..</a></li>`)
			}

			fs, err := f.Readdir(0)
			if err != nil {
				ctx.Result = InternalServerError(err.Error())
				ctx.HandleError()
				return
			}

			for _, fi := range fs {
				if fi.IsDir() {
					ctx.WriteString(`<li>┖ <a href="` + path.Join("/", opt.Prefix, rPath, fi.Name()) + `">` + path.Base(fi.Name()) + `</a></li>`)
				} else {
					if len(opt.FilterExts) > 0 && !opt.IsFilterExt(fi.Name()) {
						continue
					}

					ctx.WriteString(`<li>&nbsp; &nbsp; <a href="` + path.Join("/", opt.Prefix, rPath, fi.Name()) + `">` + filepath.Base(fi.Name()) + `</a></li>`)
				}
			}
			ctx.WriteString("</ul>")
			return
		}

		ctx.Next()
	}
}
