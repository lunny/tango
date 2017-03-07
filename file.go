// Copyright 2015 The Tango Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package tango

import "path/filepath"

// File returns a handle to serve a file
func File(path string) func(ctx *Context) {
	return func(ctx *Context) {
		ctx.ServeFile(path)
	}
}

// Dir returns a handle to serve a directory
func Dir(dir string) func(ctx *Context) {
	return func(ctx *Context) {
		params := ctx.Params()
		if len(*params) <= 0 {
			ctx.Result = NotFound()
			ctx.HandleError()
			return
		}
		ctx.ServeFile(filepath.Join(dir, (*params)[0].Value))
	}
}
