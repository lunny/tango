Tango [![Build Status](https://drone.io/github.com/lunny/tango/status.png)](https://drone.io/github.com/lunny/tango/latest) [![](http://gocover.io/_badge/github.com/lunny/tango)](http://gocover.io/github.com/lunny/tango) [English](README.md)
=======================

![Tango Logo](logo.png)

Tango 是一个微内核易扩展的Go语言Web框架.

##### 当前版本: v0.4.5   [版本更新记录](https://github.com/lunny/tango/releases)

## 简介

安装Tango：

    go get github.com/lunny/tango

一个经典的Tango例子如下：

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

然后在浏览器访问`http://localhost:8000`, 将会得到一个json返回
```
{"say":"Hello tango!"}
```

如果将上述例子中的 `true` 改为 `false`, 将会得到一个json返回
```
{"err":"something error"}
```

这段代码因为拥有一个内嵌的`tango.Json`，所以返回值会被自动的转成Json。具体返回可以参见以下文档。

## 特性

- 强大而灵活的路由设计
- 兼容已有的`http.Handler`
- 模块化设计，可以很容易写出自己的中间件
- 高性能的依赖注入方式

## 中间件 

中间件让你像AOP编程那样来操作你的Controller。

目前已有很多 [中间件 - github.com/tango-contrib](https://github.com/tango-contrib)，可以帮助你来简化工作:

- [recovery](https://github.com/lunny/tango/wiki/ZH_Recovery) - recover after panic
- [compress](https://github.com/lunny/tango/wiki/ZH_Compress) - Gzip & Deflate compression
- [static](https://github.com/lunny/tango/wiki/ZH_Static) - Serves static files
- [logger](https://github.com/lunny/tango/wiki/ZH_Logger) - Log the request & inject Logger to action struct
- [param](https://github.com/lunny/tango/wiki/ZH_Params) - get the router parameters
- [return](https://github.com/lunny/tango/wiki/ZH_Return) - Handle the returned value smartlly
- [context](https://github.com/lunny/tango/wiki/ZH_Context) - Inject context to action struct
- [session](https://github.com/tango-contrib/session) - [![Build Status](https://drone.io/github.com/tango-contrib/session/status.png)](https://drone.io/github.com/tango-contrib/session/latest) [![](http://gocover.io/_badge/github.com/tango-contrib/session)](http://gocover.io/github.com/tango-contrib/session) Session manager
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

## 获得帮助

- [Wiki](https://github.com/lunny/tango/wiki/ZH_Home)
- [API文档](https://gowalker.org/github.com/lunny/tango)
- [中文论坛](https://groups.google.com/forum/#!forum/go-tango)
- [英文论坛](https://groups.google.com/forum/#!forum/go-tango)

## 案例

- [Wego](https://github.com/go-tango/wego)
- [DBWeb](https://github.com/go-xorm/dbweb)
- [Godaily](http://godaily.org) - [github](https://github.com/godaily/news)
- [ABlog](https://github.com/fuxiaohei/ablog)
- [Gos](https://github.com/go-tango/gos)

## 讨论

QQ群：369240307

## License

This project is under BSD License. See the [LICENSE](LICENSE) file for the full license text.
