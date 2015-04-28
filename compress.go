// Copyright 2015 The Tango Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package tango

import (
	"bufio"
	"compress/flate"
	"compress/gzip"
	"fmt"
	"io"
	"net"
	"net/http"
	"path"
	"strings"
)

const (
	HeaderAcceptEncoding  = "Accept-Encoding"
	HeaderContentEncoding = "Content-Encoding"
	HeaderContentLength   = "Content-Length"
	HeaderContentType     = "Content-Type"
	HeaderVary            = "Vary"
)

type Compresser interface {
	CompressType() string
}

type GZip struct{}

func (GZip) CompressType() string {
	return "gzip"
}

type Deflate struct{}

func (Deflate) CompressType() string {
	return "deflate"
}

type Compress struct{}

func (Compress) CompressType() string {
	return "auto"
}

func Compresses(exts []string) HandlerFunc {
	extsmap := make(map[string]bool)
	for _, ext := range exts {
		extsmap[strings.ToLower(ext)] = true
	}

	return func(ctx *Context) {
		ae := ctx.Req().Header.Get("Accept-Encoding")
		if ae == "" {
			ctx.Next()
			return
		}

		if len(extsmap) > 0 {
			ext := strings.ToLower(path.Ext(ctx.Req().URL.Path))
			if _, ok := extsmap[ext]; ok {
				compress(ctx, "auto")
				return
			}
		}

		if action := ctx.Action(); action != nil {
			if c, ok := action.(Compresser); ok {
				compress(ctx, c.CompressType())
				return
			}
		}

		// if blank, then no compress
		ctx.Next()
	}
}

func compress(ctx *Context, compressType string) {
	ae := ctx.Req().Header.Get("Accept-Encoding")
	acceptCompress := strings.SplitN(ae, ",", -1)
	var writer io.Writer
	var val string

	for _, val = range acceptCompress {
		val = strings.TrimSpace(val)
		if compressType == "auto" || val == compressType {
			if val == "gzip" {
				ctx.Header().Set("Content-Encoding", "gzip")
				writer = gzip.NewWriter(ctx.ResponseWriter)
				break
			} else if val == "deflate" {
				ctx.Header().Set("Content-Encoding", "deflate")
				writer, _ = flate.NewWriter(ctx.ResponseWriter, flate.BestSpeed)
				break
			}
		}
	}

	// not supported compress method, then ignore
	if writer == nil {
		ctx.Next()
		return
	}

	// for cache server
	ctx.Header().Add(HeaderVary, "Accept-Encoding")

	gzw := &compressWriter{writer, ctx.ResponseWriter}
	ctx.ResponseWriter = gzw

	ctx.Next()

	// delete content length after we know we have been written to
	gzw.Header().Del(HeaderContentLength)
	ctx.ResponseWriter = gzw.ResponseWriter

	switch writer.(type) {
	case *gzip.Writer:
		writer.(*gzip.Writer).Close()
	case *flate.Writer:
		writer.(*flate.Writer).Close()
	}
}

type compressWriter struct {
	w io.Writer
	ResponseWriter
}

func (grw *compressWriter) Write(p []byte) (int, error) {
	if len(grw.Header().Get(HeaderContentType)) == 0 {
		grw.Header().Set(HeaderContentType, http.DetectContentType(p))
	}
	return grw.w.Write(p)
}

func (grw *compressWriter) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	hijacker, ok := grw.ResponseWriter.(http.Hijacker)
	if !ok {
		return nil, nil, fmt.Errorf("the ResponseWriter doesn't support the Hijacker interface")
	}
	return hijacker.Hijack()
}
