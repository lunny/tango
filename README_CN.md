Tango [![Build Status](https://drone.io/github.com/lunny/tango/status.png)](https://drone.io/github.com/lunny/tango/latest) [![](http://gocover.io/_badge/github.com/lunny/tango)](http://gocover.io/github.com/lunny/tango) [English](README.md)
=======================

![Tango Logo](logo.png)

Tango 是一个微内核易扩展的Go语言Web框架.

##### 当前版本: 0.2.7

## 简介

安装Tango:

    go get github.com/lunny/tango

最简单的例子:

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

然后在浏览器访问`http://localhost:8000`即可。当然了，tango其实对struct形式的支持更好。比如：

```go
package main

import "github.com/lunny/tango"

type Action struct {
    tango.Json
}

func (Action) Get() map[string]string {
    return map[string]string{
        "say": "Hello tango!",
    }
}

func main() {
    t := tango.Classic()
    t.Get("/", new(Action))
    t.Run()
}
```

这段代码因为拥有一个内嵌的`tango.Json`，所以返回值会被自动的转成Json。具体返回可以参见以下文档。

源码文档 [godoc](http://godoc.org/github.com/lunny/tango) 和 [Wiki](https://github.com/lunny/tango/wiki)

## 特性

- 强大而灵活的路由设计
- 兼容已有的`http.Handler`
- 模块化设计，可以很容易写出自定义插件
- 高性能的依赖注入方式

## 中间件 

中间件让你像切面编程那样来操作你的Controller。

目前已有很多 [中间件](https://github.com/tango-contrib)，可以帮助你来简化工作:

- recovery - recover after panic
- logger - log the request
- compress - Gzip & Deflate compression
- static - Serves static files
- logger - Log the request & inject Logger to action struct
- param - get the router parameters
- return - Handle the returned value smartlly
- request - Inject request to action struct
- response - Inject response to action struct
- [session](https://github.com/tango-contrib/session) - [![Build Status](https://drone.io/github.com/tango-contrib/session/status.png)](https://drone.io/github.com/tango-contrib/session/latest) [![](http://gocover.io/_badge/github.com/tango-contrib/session)](http://gocover.io/github.com/tango-contrib/session) Session manager
- [xsrf](https://github.com/tango-contrib/xsrf) - [![Build Status](https://drone.io/github.com/tango-contrib/xsrf/status.png)](https://drone.io/github.com/tango-contrib/xsrf/latest) [![](http://gocover.io/_badge/github.com/tango-contrib/xsrf)](http://gocover.io/github.com/tango-contrib/xsrf) Generates and validates csrf tokens
- [binding](https://github.com/tango-contrib/binding) - [![Build Status](https://drone.io/github.com/tango-contrib/binding/status.png)](https://drone.io/github.com/tango-contrib/binding/latest) [![](http://gocover.io/_badge/github.com/tango-contrib/binding)](http://gocover.io/github.com/tango-contrib/binding) Bind and validates forms
- [renders](https://github.com/tango-contrib/renders) - [![Build Status](https://drone.io/github.com/tango-contrib/renders/status.png)](https://drone.io/github.com/tango-contrib/renders/latest) [![](http://gocover.io/_badge/github.com/tango-contrib/renders)](http://gocover.io/github.com/tango-contrib/renders) Go template engine
- [dispatch](https://github.com/tango-contrib/dispatch) - [![Build Status](https://drone.io/github.com/tango-contrib/dispatch/status.png)](https://drone.io/github.com/tango-contrib/dispatch/latest) [![](http://gocover.io/_badge/github.com/tango-contrib/dispatch)](http://gocover.io/github.com/tango-contrib/dispatch) Multiple Application support on one server
- [tpongo2](https://github.com/tango-contrib/tpongo2) - [![Build Status](https://drone.io/github.com/tango-contrib/tpongo2/status.png)](https://drone.io/github.com/tango-contrib/tpongo2/latest) [![](http://gocover.io/_badge/github.com/tango-contrib/tpongo2)](http://gocover.io/github.com/tango-contrib/tpongo2) [Pongo2](https://github.com/flosch/pongo2) teamplte engine support
- [captcha](https://github.com/tango-contrib/captcha) - [![Build Status](https://drone.io/github.com/tango-contrib/captcha/status.png)](https://drone.io/github.com/tango-contrib/captcha/latest) [![](http://gocover.io/_badge/github.com/tango-contrib/captcha)](http://gocover.io/github.com/tango-contrib/captcha) Captcha
- [events](https://github.com/tango-contrib/events) - [![Build Status](https://drone.io/github.com/tango-contrib/events/status.png)](https://drone.io/github.com/tango-contrib/events/latest) [![](http://gocover.io/_badge/github.com/tango-contrib/events)](http://gocover.io/github.com/tango-contrib/events) Before and After
- [flash](https://github.com/tango-contrib/flash) - [![Build Status](https://drone.io/github.com/tango-contrib/flash/status.png)](https://drone.io/github.com/tango-contrib/flash/latest) [![](http://gocover.io/_badge/github.com/tango-contrib/flash)](http://gocover.io/github.com/tango-contrib/flash) Share data between requests

## 获得帮助

- [API文档](https://gowalker.org/github.com/lunny/tango)

## 案例

- [Wego](https://github.com/go-tango/wego)
- [ABlog](https://github.com/fuxiaohei/ablog)

## 讨论

QQ群：369240307

## License

This project is under BSD License. See the [LICENSE](LICENSE) file for the full license text.