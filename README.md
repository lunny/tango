Tango [![Build Status](https://drone.io/github.com/lunny/tango/status.png)](https://drone.io/github.com/lunny/tango/latest) [![](http://gocover.io/_badge/github.com/lunny/tango)](http://gocover.io/github.com/lunny/tango) [简体中文](README_CN.md)
=======================

[![Join the chat at https://gitter.im/lunny/tango](https://badges.gitter.im/Join%20Chat.svg)](https://gitter.im/lunny/tango?utm_source=badge&utm_medium=badge&utm_campaign=pr-badge&utm_content=badge)

![Tango Logo](logo.png)

Package tango is a micro & pluggable web framework for Go.

##### Current version: v0.4.6   [Version History](https://github.com/lunny/tango/releases)

## Getting Started

To install Tango:

    go get github.com/lunny/tango

A classic usage of Tango below:

```go
package main

import (
    "errors"
    "github.com/lunny/tango"
)

type Action struct {
    tango.Json
}

func (Action) Get() interface{} {
    if true {
        return map[string]string{
            "say": "Hello tango!",
        }
    }
    return errors.New("something error")
}

func main() {
    t := tango.Classic()
    t.Get("/", new(Action))
    t.Run()
}
```

Then visit `http://localhost:8000` on your browser. You will get
```
{"say":"Hello tango!"}
```

If you change `true` after `if` to `false`, then you will get
```
{"err":"something error"}
```

This code will automatically convert returned map or error to a json because we has an embedded struct `tango.Json`.

## Features

- Powerful routing & Flexible routes combinations.
- Directly integrate with existing services.
- Easy to plugin features with modular design.
- High performance dependency injection embedded.

## Middlewares 

Middlewares allow you easily plugin features for your Tango applications.

There are already many [middlewares](https://github.com/tango-contrib) to simplify your work:

- [recovery](https://github.com/lunny/tango/wiki/Recovery) - recover after panic
- [compress](https://github.com/lunny/tango/wiki/Compress) - Gzip & Deflate compression
- [static](https://github.com/lunny/tango/wiki/Static) - Serves static files
- [logger](https://github.com/lunny/tango/wiki/Logger) - Log the request & inject Logger to action struct
- [param](https://github.com/lunny/tango/wiki/Params) - get the router parameters
- [return](https://github.com/lunny/tango/wiki/Return) - Handle the returned value smartlly
- [context](https://github.com/lunny/tango/wiki/Context) - Inject context to action struct
- [session](https://github.com/tango-contrib/session) - [![Build Status](https://drone.io/github.com/tango-contrib/session/status.png)](https://drone.io/github.com/tango-contrib/session/latest) [![](http://gocover.io/_badge/github.com/tango-contrib/session)](http://gocover.io/github.com/tango-contrib/session) Session manager, [session-redis](http://github.com/tango-contrib/session-redis), [session-nodb](http://github.com/tango-contrib/session-nodb), [session-ledis](http://github.com/tango-contrib/session-ledis)
- [xsrf](https://github.com/tango-contrib/xsrf) - [![Build Status](https://drone.io/github.com/tango-contrib/xsrf/status.png)](https://drone.io/github.com/tango-contrib/xsrf/latest) [![](http://gocover.io/_badge/github.com/tango-contrib/xsrf)](http://gocover.io/github.com/tango-contrib/xsrf) Generates and validates csrf tokens
- [binding](https://github.com/tango-contrib/binding) - [![Build Status](https://drone.io/github.com/tango-contrib/binding/status.png)](https://drone.io/github.com/tango-contrib/binding/latest) [![](http://gocover.io/_badge/github.com/tango-contrib/binding)](http://gocover.io/github.com/tango-contrib/binding) Bind and validates forms
- [renders](https://github.com/tango-contrib/renders) - [![Build Status](https://drone.io/github.com/tango-contrib/renders/status.png)](https://drone.io/github.com/tango-contrib/renders/latest) [![](http://gocover.io/_badge/github.com/tango-contrib/renders)](http://gocover.io/github.com/tango-contrib/renders) Go template engine
- [dispatch](https://github.com/tango-contrib/dispatch) - [![Build Status](https://drone.io/github.com/tango-contrib/dispatch/status.png)](https://drone.io/github.com/tango-contrib/dispatch/latest) [![](http://gocover.io/_badge/github.com/tango-contrib/dispatch)](http://gocover.io/github.com/tango-contrib/dispatch) Multiple Application support on one server
- [tpongo2](https://github.com/tango-contrib/tpongo2) - [![Build Status](https://drone.io/github.com/tango-contrib/tpongo2/status.png)](https://drone.io/github.com/tango-contrib/tpongo2/latest) [![](http://gocover.io/_badge/github.com/tango-contrib/tpongo2)](http://gocover.io/github.com/tango-contrib/tpongo2) [Pongo2](https://github.com/flosch/pongo2) teamplte engine support
- [captcha](https://github.com/tango-contrib/captcha) - [![Build Status](https://drone.io/github.com/tango-contrib/captcha/status.png)](https://drone.io/github.com/tango-contrib/captcha/latest) [![](http://gocover.io/_badge/github.com/tango-contrib/captcha)](http://gocover.io/github.com/tango-contrib/captcha) Captcha
- [events](https://github.com/tango-contrib/events) - [![Build Status](https://drone.io/github.com/tango-contrib/events/status.png)](https://drone.io/github.com/tango-contrib/events/latest) [![](http://gocover.io/_badge/github.com/tango-contrib/events)](http://gocover.io/github.com/tango-contrib/events) Before and After
- [flash](https://github.com/tango-contrib/flash) - [![Build Status](https://drone.io/github.com/tango-contrib/flash/status.png)](https://drone.io/github.com/tango-contrib/flash/latest) [![](http://gocover.io/_badge/github.com/tango-contrib/flash)](http://gocover.io/github.com/tango-contrib/flash) Share data between requests
- [debug](https://github.com/tango-contrib/debug) - [![Build Status](https://drone.io/github.com/tango-contrib/debug/status.png)](https://drone.io/github.com/tango-contrib/debug/latest) [![](http://gocover.io/_badge/github.com/tango-contrib/debug)](http://gocover.io/github.com/tango-contrib/debug) show detail debug infomaton on log
- [basicauth](https://github.com/tango-contrib/basicauth) - [![Build Status](https://drone.io/github.com/tango-contrib/basicauth/status.png)](https://drone.io/github.com/tango-contrib/basicauth/latest) [![](http://gocover.io/_badge/github.com/tango-contrib/basicauth)](http://gocover.io/github.com/tango-contrib/basicauth) basicauth middleware

## Getting Help

- [Wiki](https://github.com/lunny/tango/wiki)
- [API Reference](https://gowalker.org/github.com/lunny/tango)
- [Discuss](https://groups.google.com/forum/#!forum/go-tango)

## Cases

- [Wego](https://github.com/go-tango/wego)
- [dbweb](https://github.com/go-xorm/dbweb)
- [Godaily](http://godaily.org) - [github](https://github.com/godaily/news)
- [FXH Blog](https://github.com/gofxh/blog)
- [Gos](https://github.com/go-tango/gos)

## License

This project is under BSD License. See the [LICENSE](LICENSE) file for the full license text.
