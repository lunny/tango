// Copyright 2015 The Tango Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Tango is a micro & pluggable web framework for Go language.

// 	package main

// 	import "github.com/lunny/tango"

// 	type Action struct {
// 	}

// 	func (Action) Get() string {
// 	    return "Hello tango!"
// 	}

// 	func main() {
// 	    t := tango.Classic()
// 	    t.Get("/", new(Action))
// 	    t.Run()
// 	}

// Middlewares allow you easily plugin/unplugin features for your Tango applications.

// There are already many [middlewares](https://github.com/tango-contrib) to simplify your work:

// - recovery - recover after panic
// - compress - Gzip & Deflate compression
// - static - Serves static files
// - logger - Log the request & inject Logger to action struct
// - param - get the router parameters
// - return - Handle the returned value smartlly
// - ctx - Inject context to action struct

// - [session](https://github.com/tango-contrib/session) - Session manager, with stores support:
//   * Memory - memory as a session store
//   * [Redis](https://github.com/tango-contrib/session-redis) - redis server as a session store
//   * [nodb](https://github.com/tango-contrib/session-nodb) - nodb as a session store
//   * [ledis](https://github.com/tango-contrib/session-ledis) - ledis server as a session store)
// - [xsrf](https://github.com/tango-contrib/xsrf) - Generates and validates csrf tokens
// - [binding](https://github.com/tango-contrib/binding) - Bind and validates forms
// - [renders](https://github.com/tango-contrib/renders) - Go template engine
// - [dispatch](https://github.com/tango-contrib/dispatch) - Multiple Application support on one server
// - [tpongo2](https://github.com/tango-contrib/tpongo2) - Pongo2 teamplte engine support
// - [captcha](https://github.com/tango-contrib/captcha) - Captcha
// - [events](https://github.com/tango-contrib/events) - Before and After
// - [flash](https://github.com/tango-contrib/flash) - Share data between requests
// - [debug](https://github.com/tango-contrib/debug) - Show detail debug infomaton on log

package tango
