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

// StatusResult describes http response
type StatusResult struct {
	Code   int
	Result interface{}
}

// enumerate all the return response types
const (
	autoResponse = iota
	jsonResponse
	xmlResponse
)

// ResponseTyper describes reponse type
type ResponseTyper interface {
	ResponseType() int
}

// Json describes return JSON type
// Deprecated: use JSON instead
type Json struct{}

// ResponseType implementes ResponseTyper
func (Json) ResponseType() int {
	return jsonResponse
}

// JSON describes return JSON type
type JSON struct{}

// ResponseType implementes ResponseTyper
func (JSON) ResponseType() int {
	return jsonResponse
}

// Xml descirbes return XML type
// Deprecated: use XML instead
type Xml struct{}

// ResponseType implementes ResponseTyper
func (Xml) ResponseType() int {
	return xmlResponse
}

// XML descirbes return XML type
type XML struct{}

// ResponseType implementes ResponseTyper
func (XML) ResponseType() int {
	return xmlResponse
}

func isNil(a interface{}) bool {
	if a == nil {
		return true
	}
	aa := reflect.ValueOf(a)
	return !aa.IsValid() || (aa.Type().Kind() == reflect.Ptr && aa.IsNil())
}

// XMLError describes return xml error
type XMLError struct {
	XMLName xml.Name `xml:"err"`
	Content string   `xml:"content"`
}

// XMLString describes return xml string
type XMLString struct {
	XMLName xml.Name `xml:"string"`
	Content string   `xml:"content"`
}

// Return returns a tango middleware to handler return values
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
		var statusCode = 0
		if res, ok := ctx.Result.(*StatusResult); ok {
			statusCode = res.Code
			result = res.Result
		}

		if rt == jsonResponse {
			encoder := json.NewEncoder(ctx)
			if len(ctx.Header().Get("Content-Type")) <= 0 {
				ctx.Header().Set("Content-Type", "application/json; charset=UTF-8")
			}

			switch res := result.(type) {
			case AbortError:
				if statusCode == 0 {
					statusCode = res.Code()
				}
				ctx.WriteHeader(statusCode)
				encoder.Encode(map[string]string{
					"err": res.Error(),
				})
			case error:
				if statusCode == 0 {
					statusCode = http.StatusOK
				}
				ctx.WriteHeader(statusCode)
				encoder.Encode(map[string]string{
					"err": res.Error(),
				})
			case string:
				if statusCode == 0 {
					statusCode = http.StatusOK
				}
				ctx.WriteHeader(statusCode)
				encoder.Encode(map[string]string{
					"content": res,
				})
			case []byte:
				if statusCode == 0 {
					statusCode = http.StatusOK
				}
				ctx.WriteHeader(statusCode)
				encoder.Encode(map[string]string{
					"content": string(res),
				})
			default:
				if statusCode == 0 {
					statusCode = http.StatusOK
				}
				ctx.WriteHeader(statusCode)
				err := encoder.Encode(result)
				if err != nil {
					ctx.Result = err
					encoder.Encode(map[string]string{
						"err": err.Error(),
					})
				}
			}

			return
		} else if rt == xmlResponse {
			encoder := xml.NewEncoder(ctx)
			if len(ctx.Header().Get("Content-Type")) <= 0 {
				ctx.Header().Set("Content-Type", "application/xml; charset=UTF-8")
			}
			switch res := result.(type) {
			case AbortError:
				if statusCode == 0 {
					statusCode = res.Code()
				}
				ctx.WriteHeader(statusCode)
				encoder.Encode(XMLError{
					Content: res.Error(),
				})
			case error:
				if statusCode == 0 {
					statusCode = http.StatusOK
				}
				ctx.WriteHeader(statusCode)
				encoder.Encode(XMLError{
					Content: res.Error(),
				})
			case string:
				if statusCode == 0 {
					statusCode = http.StatusOK
				}
				ctx.WriteHeader(statusCode)
				encoder.Encode(XMLString{
					Content: res,
				})
			case []byte:
				if statusCode == 0 {
					statusCode = http.StatusOK
				}
				ctx.WriteHeader(statusCode)
				encoder.Encode(XMLString{
					Content: string(res),
				})
			default:
				if statusCode == 0 {
					statusCode = http.StatusOK
				}
				ctx.WriteHeader(statusCode)
				err := encoder.Encode(result)
				if err != nil {
					ctx.Result = err
					encoder.Encode(XMLError{
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
			if statusCode == 0 {
				statusCode = http.StatusOK
			}
			ctx.WriteHeader(statusCode)
			ctx.Write(res)
		case string:
			if statusCode == 0 {
				statusCode = http.StatusOK
			}
			ctx.WriteHeader(statusCode)
			ctx.WriteString(res)
		}
	}
}
