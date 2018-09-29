Tango [简体中文](README_CN.md)
=======================

[![GitCI](https://gitci.cn/api/badges/lunny/tango/status.svg)](https://gitci.cn/gh/lunny/tango) [![codecov](https://codecov.io/gh/lunny/tango/branch/master/graph/badge.svg)](https://codecov.io/gh/lunny/tango)
[![](https://goreportcard.com/badge/github.com/lunny/tango)](https://goreportcard.com/report/github.com/lunny/tango)
[![Join the chat at https://img.shields.io/discord/323705316027924491.svg](https://img.shields.io/discord/323705316027924491.svg)](https://discord.gg/7Ckxjwu)

![Tango Logo](logo.png)

Package tango is a micro & pluggable web framework for Go.

##### Current version: v0.5.0   [Version History](https://github.com/lunny/tango/releases)

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
    tango.JSON
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

This code will automatically convert returned map or error to a json because we has an embedded struct `tango.JSON`.

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
- [session](https://github.com/tango-contrib/session) - [![CircleCI](https://circleci.com/gh/tango-contrib/session/tree/master.svg?style=svg)](https://circleci.com/gh/tango-contrib/session/tree/master) [![codecov](https://codecov.io/gh/tango-contrib/session/branch/master/graph/badge.svg)](https://codecov.io/gh/tango-contrib/session) Session manager, [session-redis](http://github.com/tango-contrib/session-redis), [session-nodb](http://github.com/tango-contrib/session-nodb), [session-ledis](http://github.com/tango-contrib/session-ledis), [session-ssdb](http://github.com/tango-contrib/session-ssdb)
- [xsrf](https://github.com/tango-contrib/xsrf) - [![CircleCI](https://circleci.com/gh/tango-contrib/xsrf/tree/master.svg?style=svg)](https://circleci.com/gh/tango-contrib/xsrf/tree/master) [![codecov](https://codecov.io/gh/tango-contrib/xsrf/branch/master/graph/badge.svg)](https://codecov.io/gh/tango-contrib/xsrf) Generates and validates csrf tokens
- [binding](https://github.com/tango-contrib/binding) - [![CircleCI](https://circleci.com/gh/tango-contrib/binding/tree/master.svg?style=svg)](https://circleci.com/gh/tango-contrib/binding/tree/master) [![codecov](https://codecov.io/gh/tango-contrib/binding/branch/master/graph/badge.svg)](https://codecov.io/gh/tango-contrib/binding) Bind and validates forms
- [renders](https://github.com/tango-contrib/renders) - [![CircleCI](https://circleci.com/gh/tango-contrib/renders/tree/master.svg?style=svg)](https://circleci.com/gh/tango-contrib/renders/tree/master) [![codecov](https://codecov.io/gh/tango-contrib/renders/branch/master/graph/badge.svg)](https://codecov.io/gh/tango-contrib/renders) Go template engine
- [dispatch](https://github.com/tango-contrib/dispatch) - [![CircleCI](https://circleci.com/gh/tango-contrib/dispatch/tree/master.svg?style=svg)](https://circleci.com/gh/tango-contrib/dispatch/tree/master) [![codecov](https://codecov.io/gh/tango-contrib/dispatch/branch/master/graph/badge.svg)](https://codecov.io/gh/tango-contrib/dispatch) Multiple Application support on one server
- [tpongo2](https://github.com/tango-contrib/tpongo2) - [![CircleCI](https://circleci.com/gh/tango-contrib/tpongo2/tree/master.svg?style=svg)](https://circleci.com/gh/tango-contrib/tpongo2/tree/master) [![codecov](https://codecov.io/gh/tango-contrib/tpongo2/branch/master/graph/badge.svg)](https://codecov.io/gh/tango-contrib/tpongo2) [Pongo2](https://github.com/flosch/pongo2) teamplte engine support
- [captcha](https://github.com/tango-contrib/captcha) - [![CircleCI](https://circleci.com/gh/tango-contrib/captcha/tree/master.svg?style=svg)](https://circleci.com/gh/tango-contrib/captcha/tree/master) [![codecov](https://codecov.io/gh/tango-contrib/captcha/branch/master/graph/badge.svg)](https://codecov.io/gh/tango-contrib/captcha) Captcha
- [events](https://github.com/tango-contrib/events) - [![CircleCI](https://circleci.com/gh/tango-contrib/events/tree/master.svg?style=svg)](https://circleci.com/gh/tango-contrib/events/tree/master) [![codecov](https://codecov.io/gh/tango-contrib/events/branch/master/graph/badge.svg)](https://codecov.io/gh/tango-contrib/events) Before and After
- [flash](https://github.com/tango-contrib/flash) - [![CircleCI](https://circleci.com/gh/tango-contrib/flash/tree/master.svg?style=svg)](https://circleci.com/gh/tango-contrib/flash/tree/master) [![codecov](https://codecov.io/gh/tango-contrib/flash/branch/master/graph/badge.svg)](https://codecov.io/gh/tango-contrib/flash) Share data between requests
- [debug](https://github.com/tango-contrib/debug) - [![CircleCI](https://circleci.com/gh/tango-contrib/debug/tree/master.svg?style=svg)](https://circleci.com/gh/tango-contrib/debug/tree/master) [![codecov](https://codecov.io/gh/tango-contrib/debug/branch/master/graph/badge.svg)](https://codecov.io/gh/tango-contrib/debug) show detail debug infomaton on log
- [basicauth](https://github.com/tango-contrib/basicauth) - [![CircleCI](https://circleci.com/gh/tango-contrib/basicauth/tree/master.svg?style=svg)](https://circleci.com/gh/tango-contrib/basicauth/tree/master) [![codecov](https://codecov.io/gh/tango-contrib/basicauth/branch/master/graph/badge.svg)](https://codecov.io/gh/tango-contrib/basicauth) basicauth middleware
- [authz](https://github.com/tango-contrib/authz) - [![Build Status](https://travis-ci.org/tango-contrib/authz.svg?branch=master)](https://travis-ci.org/tango-contrib/authz) [![Coverage Status](https://coveralls.io/repos/github/tango-contrib/authz/badge.svg?branch=master)](https://coveralls.io/github/tango-contrib/authz?branch=master) manage permissions via ACL, RBAC, ABAC
- [cache](https://github.com/tango-contrib/cache) - [![CircleCI](https://circleci.com/gh/tango-contrib/cache/tree/master.svg?style=svg)](https://circleci.com/gh/tango-contrib/cache/tree/master) [![codecov](https://codecov.io/gh/tango-contrib/cache/branch/master/graph/badge.svg)](https://codecov.io/gh/tango-contrib/cache) cache middleware - cache-memory, cache-file, [cache-ledis](https://github.com/tango-contrib/cache-ledis), [cache-nodb](https://github.com/tango-contrib/cache-nodb), [cache-mysql](https://github.com/tango-contrib/cache-mysql), [cache-postgres](https://github.com/tango-contrib/cache-postgres), [cache-memcache](https://github.com/tango-contrib/cache-memcache), [cache-redis](https://github.com/tango-contrib/cache-redis)
- [rbac](https://github.com/tango-contrib/rbac) - [![CircleCI](https://circleci.com/gh/tango-contrib/rbac/tree/master.svg?style=svg)](https://circleci.com/gh/tango-contrib/rbac/tree/master) [![codecov](https://codecov.io/gh/tango-contrib/rbac/branch/master/graph/badge.svg)](https://codecov.io/gh/tango-contrib/rbac) rbac control

## Documentation

- [Manual](http://gobook.io/read/github.com/go-tango/manual-en-US/), And you are welcome to contribue for the book by git PR to [github.com/go-tango/manual-en-US](https://github.com/go-tango/manual-en-US)
- [操作手册](http://gobook.io/read/github.com/go-tango/manual-zh-CN/)，您也可以访问 [github.com/go-tango/manual-zh-CN](https://github.com/go-tango/manual-zh-CN)为本手册进行贡献
- [API Reference](https://gowalker.org/github.com/lunny/tango)

## Discuss

- [Google Group - English](https://groups.google.com/forum/#!forum/go-tango)
- QQ Group - 简体中文 #369240307

## Cases

- [GopherTC](https://github.com/jimmykuu/gopher/tree/2.0) - China Discuss Forum
- [Wego](https://github.com/go-tango/wego) - Discuss Forum
- [dbweb](https://github.com/go-xorm/dbweb) - DB management web UI
- [Godaily](http://godaily.org) - [github](https://github.com/godaily/news)
- [Pugo](https://github.com/go-xiaohei/pugo) - A pugo blog
- [Gos](https://github.com/go-tango/gos) - Static web server
- [GoFtpd](https://github.com/goftp/ftpd) - Pure Go cross-platform ftp server

## License

This project is under BSD License. See the [LICENSE](LICENSE) file for the full license text.
