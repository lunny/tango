Tango [![Build Status](https://drone.io/github.com/lunny/tango/status.png)](https://drone.io/github.com/lunny/tango/latest) [![](http://gocover.io/_badge/github.com/lunny/tango)](http://gocover.io/github.com/lunny/tango)
=======================

![Tango Logo](logo.png)

Package tango is a micro & pluggable web framework for Go.

##### Current version: 0.1.0

## Getting Started

To install Tango:

    go get github.com/lunny/tango

The very basic usage of Tango:

```go
package main

import "github.com/lunny/tango"

func main() {
    t := tango.Classic()
    t.Get("/", func() string {
        return "Hello tango!"
    })
    t.Run()
}
```

Then visit `http://localhost:8000` on your browser. Of course, tango support struct form also.

```go
package main

import "github.com/lunny/tango"

type Action struct {}

func (Action) Get() string {
    return "Hello tango!"
}

func main() {
    t := tango.Classic()
    t.Get("/", new(Action))
    t.Run()
}
```

More document, please see [godoc](http://godoc.org/github.com/lunny/tango) and [Wiki](https://github.com/lunny/tango/wiki)

## Features

- Powerful routing & Flexible routes combinations.
- Directly integrate with existing services.
- Easy to plugin/unplugin features with modular design.
- High Performance dependency injection embbed.

## Middlewares 

Middlewares allow you easily plugin/unplugin features for your Tango applications.

There are already many [middlewares](https://github.com/tango-contrib) to simplify your work:

- recovery - recover after panic
- logger - log the request
- compress - Gzip & Deflate compression
- static - Serves static files
- logger - Log the request & inject Logger to action struct
- param - get the router parameters
- return - Handle the returned value smartlly
- request - Inject request to action struct
- response - Inject response to action struct
- [session](https://github.com/tango-contrib/session) - Session manager
- [xsrf](https://github.com/tango-contrib/xsrf) - Generates and validates csrf tokens
- [bind](https://github.com/tango-contrib/bind) - Bind and validates forms
- [render](https://github.com/tango-contrib/render) - Go template engine
- [dispatch](https://github.com/tango-contrib/dispatch) - Multiple Application support on one server
- [tpongo2](https://github.com/tango-contrib/tpongo2) - Pango2 teamplte engine support

## Getting Help

- [API Reference](https://gowalker.org/github.com/lunny/tango)

## License

This project is under BSD License. See the [LICENSE](LICENSE) file for the full license text.