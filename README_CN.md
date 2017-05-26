Tango [![CircleCI](https://circleci.com/gh/lunny/tango/tree/master.svg?style=svg)](https://circleci.com/gh/lunny/tango/tree/master)  [![codecov](https://codecov.io/gh/lunny/tango/branch/master/graph/badge.svg)](https://codecov.io/gh/lunny/tango)
 [![](https://goreportcard.com/badge/github.com/lunny/tango)](https://goreportcard.com/report/github.com/lunny/tango) [![Join the chat at https://gitter.im/lunny/tango](https://badges.gitter.im/Join%20Chat.svg)](https://gitter.im/lunny/tango?utm_source=badge&utm_medium=badge&utm_campaign=pr-badge&utm_content=badge) [English](README.md)
=======================

![Tango Logo](logo.png)

Tango 是一个微内核的Go语言Web框架，采用模块化和注入式的设计理念。开发者可根据自身业务逻辑来选择性的装卸框架的功能，甚至利用丰富的中间件来搭建一个全栈式Web开发框架。

## 最近更新
- [2016-5-12] 开放Route级别中间件支持
- [2016-3-16] Group完善中间件支持，Route支持中间件
- [2016-2-1] 新增 session-ssdb，支持将ssdb作为session的后端存储
- [2015-10-23] 更新[renders](https://github.com/tango-contrib/renders)插件，解决模板修改后需要刷新两次才能生效的问题

## 特性
- 强大而灵活的路由设计
- 兼容已有的 `http.Handler`
- 基于中间件的模块化设计，灵活定制框架功能
- 高性能的依赖注入方式

## 安装Tango：

    go get github.com/lunny/tango

## 快速入门

一个经典的Tango例子如下：

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

然后在浏览器访问`http://localhost:8000`, 将会得到一个json返回

```
{"say":"Hello tango!"}
```

如果将上述例子中的 `true` 改为 `false`, 将会得到一个json返回

```
{"err":"something error"}
```

这段代码因为拥有一个内嵌的`tango.JSON`，所以返回值会被自动的转成Json

## 文档

- [Manual](http://gobook.io/read/github.com/go-tango/manual-en-US/), And you are welcome to contribue for the book by git PR to [github.com/go-tango/manual-en-US](https://github.com/go-tango/manual-en-US)
- [操作手册](http://gobook.io/read/github.com/go-tango/manual-zh-CN/)，您也可以访问 [github.com/go-tango/manual-zh-CN](https://github.com/go-tango/manual-zh-CN)为本手册进行贡献
- [API Reference](https://gowalker.org/github.com/lunny/tango)

## 交流讨论

- QQ群：369240307
- [论坛](https://groups.google.com/forum/#!forum/go-tango)

## 使用案例
- [Wego](https://github.com/go-tango/wego)  tango结合[xorm](http://www.xorm.io/)开发的论坛
- [Pugo](https://github.com/go-xiaohei/pugo) 博客
- [DBWeb](https://github.com/go-xorm/dbweb) 基于Web的数据库管理工具
- [Godaily](http://godaily.org) - [github](https://github.com/godaily/news) RSS聚合工具
- [Gos](https://github.com/go-tango/gos)  简易的Web静态文件服务端
- [GoFtpd](https://github.com/goftp/ftpd) - 纯Go的跨平台FTP服务器

## 中间件列表

[中间件](https://github.com/tango-contrib)可以重用代码并且简化工作：

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

## License
This project is under BSD License. See the [LICENSE](LICENSE) file for the full license text.
