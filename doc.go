// Copyright 2014 lunny. All rights reserved.
// Use of this source code is governed by a BSD
// license that can be found in the LICENSE file.
/*
Tango is a micro & pluggable web framework for Go language.


package main

import "github.com/lunny/tango"

type Action struct {
}

func (Action) Get() string {
    return "Hello tango!"
}

func main() {
    t := tango.Classic()
    t.Get("/", new(Action))
    t.Run()
}

Middlewares allow you easily plugin/unplugin features for your Tango applications.

There are already many [middlewares](https://github.com/tango-contrib) to simplify your work:

- compress - Gzip & Deflate compression
- static - Serves static files
- logger - Log the request & inject Logger to action struct
- return - Handle the returned value smartlly
- request - Inject request to action struct
- response - Inject response to action struct

- [session](https://github.com/tango-contrib/session) - Session manager
- [xsrf](https://github.com/tango-contrib/xsrf) - Generates and validates csrf tokens
- [bind](https://github.com/tango-contrib/bind) - Bind and validates forms
- [render](https://github.com/tango-contrib/render) - Go template engine
- [dispatch](https://github.com/tango-contrib/dispatch) - Multiple Application support on one server
- [pongo2](https://github.com/tango-contrib/pongo2) - Pongo2 teamplte engine support

*/

package tango