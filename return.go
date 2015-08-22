// Copyright 2015 The Tango Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package tango

import (
	"encoding/json"
	"encoding/xml"
	"net/http"
	"reflect"
)

type StatusResult struct {
	Code int
	Result interface{}
}

const (
	AutoResponse = iota
	JsonResponse
	XmlResponse
)

type ResponseTyper interface {
	ResponseType() int
}

type Json struct{}

func (Json) ResponseType() int {
	return JsonResponse
}

type Xml struct{}

func (Xml) ResponseType() int {
	return XmlResponse
}

func isNil(a interface{}) bool {
	if a == nil {
		return true
	}
	aa := reflect.ValueOf(a)
	return !aa.IsValid() || (aa.Type().Kind() == reflect.Ptr && aa.IsNil())
}

type XmlError struct {
	XMLName xml.Name `xml:"err"`
	Content string   `xml:"content"`
}

type XmlString struct {
	XMLName xml.Name `xml:"string"`
	Content string   `xml:"content"`
}

func Return() HandlerFunc {
	return func(ctx *Context) {
		var rt int
		action := ctx.Action()
		if action != nil {
			if i, ok := action.(ResponseTyper); ok {
				rt = i.ResponseType()
			}
		}

		ctx.Next()

		// if no route match or has been write, then return
		if action == nil || ctx.Written() {
			return
		}

		// if there is no return value or return nil
		if isNil(ctx.Result) {
			// then we return blank page
			ctx.Result = ""
		}

		var result = ctx.Result
		if res, ok := ctx.Result.(*StatusResult); ok {
			ctx.WriteHeader(res.Code)
			result = res.Result
		}

		if rt == JsonResponse {
			encoder := json.NewEncoder(ctx)
			ctx.Header().Set("Content-Type", "application/json; charset=UTF-8")

			switch res := result.(type) {
			case AbortError:
				ctx.WriteHeader(res.Code())
				encoder.Encode(map[string]string{
					"err": res.Error(),
				})
			case error:
				encoder.Encode(map[string]string{
					"err": res.Error(),
				})
			case string:
				encoder.Encode(map[string]string{
					"content": res,
				})
			case []byte:
				encoder.Encode(map[string]string{
					"content": string(res),
				})
			default:
				err := encoder.Encode(result)
				if err != nil {
					ctx.Result = err
					encoder.Encode(map[string]string{
						"err": err.Error(),
					})
				}
			}

			return
		} else if rt == XmlResponse {
			encoder := xml.NewEncoder(ctx)
			ctx.Header().Set("Content-Type", "application/xml; charset=UTF-8")
			switch res := result.(type) {
			case AbortError:
				ctx.WriteHeader(res.Code())
				encoder.Encode(XmlError{
					Content: res.Error(),
				})
			case error:
				encoder.Encode(XmlError{
					Content: res.Error(),
				})
			case string:
				encoder.Encode(XmlString{
					Content: res,
				})
			case []byte:
				encoder.Encode(XmlString{
					Content: string(res),
				})
			default:
				err := encoder.Encode(result)
				if err != nil {
					ctx.Result = err
					encoder.Encode(XmlError{
						Content: err.Error(),
					})
				}
			}
			return
		}

		switch res := result.(type) {
		case AbortError, error:
			ctx.HandleError()
		case []byte:
			ctx.WriteHeader(http.StatusOK)
			ctx.Write(res)
		case string:
			ctx.WriteHeader(http.StatusOK)
			ctx.Write([]byte(res))
		}
	}
}
