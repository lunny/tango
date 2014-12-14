Tango [![Build Status](https://drone.io/github.com/lunny/tango/status.png)](https://drone.io/github.com/lunny/tango/latest) [![](http://gocover.io/_badge/github.com/lunny/tango)](http://gocover.io/github.com/lunny/tango)
=======================

![Tango Logo](tangologo.png)

Package tango is a high productive and modular design web framework in Go.

##### Current version: 0.1.0

## Getting Started

To install Tango:

    go get github.com/unny/tango

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

- Powerful routing with suburl.
- Flexible routes combinations.
- Unlimited nested group routers.
- Directly integrate with existing services.
- Easy to plugin/unplugin features with modular design.
- Handy & Performance dependency injection embbed.
- Good performance.

## Middlewares 

Middlewares allow you easily plugin/unplugin features for your Tango applications.

There are already many [middlewares](https://github.com/tango-contrib) to simplify your work:

- compress - Gzip compression to all requests
- render - Go template engine
- static - Serves static files
- [session](https://github.com/tango-contrib/session) - Session manager
- [csrf](https://github.com/tango-contrib/csrf) - Generates and validates csrf tokens

## Use Cases


## Getting Help

- [API Reference](https://gowalker.org/github.com/lunny/tango)

## Credits


## License

This project is under BSD License. See the [LICENSE](LICENSE) file for the full license text.