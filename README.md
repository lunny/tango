Tango [![Build Status](https://drone.io/github.com/lunny/tango/status.png)](https://drone.io/github.com/lunny/tango/latest) [![](http://gocover.io/_badge/github.com/lunny/tango)](http://gocover.io/github.com/lunny/tango)
=======================

![Tango Logo](tangologo.png)

Package tango is a micro & pluggable web framework in Go.

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
        return "Hello world!"
    })
    t.Run()
}
```

## Features

- Powerful routing & Flexible routes combinations.
- Directly integrate with existing services.
- Easy to plugin/unplugin features with modular design.
- High Performance dependency injection embbed.

## Middlewares 

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

## Use Cases


## Getting Help

- [API Reference](https://gowalker.org/github.com/lunny/tango)

## Credits


## License

This project is under BSD License. See the [LICENSE](LICENSE) file for the full license text.