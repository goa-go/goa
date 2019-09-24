# Goa

[![Build Status](https://travis-ci.org/goa-go/goa.svg?branch=master)](https://travis-ci.org/goa-go/goa)
[![Codecov](https://codecov.io/gh/goa-go/goa/branch/master/graph/badge.svg)](https://codecov.io/github/goa-go/goa?branch=master)
[![Go Doc](https://godoc.org/github.com/goa-go/goa?status.svg)](http://godoc.org/github.com/goa-go/goa)
[![Go Report](https://goreportcard.com/badge/github.com/goa-go/goa)](https://goreportcard.com/report/github.com/goa-go/goa)
[![Mentioned in Awesome Go](https://awesome.re/mentioned-badge.svg)](https://github.com/avelino/awesome-go)
[![PR's Welcome](https://img.shields.io/badge/PRs-welcome-brightgreen.svg?style=flat)](https://github.com/goa-go/goa/pull/new) 

Goa is under construction, if you are familiar with [koa](https://github.com/koajs/koa) or go and interested in this project, please join us.

## What is goa?

goa = go + [koa](https://github.com/koajs/koa)

Just like koa, goa is also not bundled with any middleware. But you can expand functionality to meet your needs at will by middlware. It is flexible, light, high-performance and extensible.

## Installation

```bash
$ go get -u github.com/goa-go/goa
```

##  Hello Goa

```go
func main() {
  app := goa.New()

  app.Use(func(c *goa.Context) {
    c.String("Hello Goa!")
  })
  log.Fatal(app.Listen(":3000"))
}
```

## Middleware

Goa is a web framework based on middleware.
Here is an example of using goa-router and logger.
```go
package main

import (
  "fmt"
  "log"
  "time"

  "github.com/goa-go/goa"
  "github.com/goa-go/router"
)

func logger(c *goa.Context) {
  start := time.Now()

  fmt.Printf(
    "[%s] <-- %s %s\n",
    start.Format("2006-01-02 15:04:05"),
    c.Method,
    c.URL,
  )
  c.Next()
  fmt.Printf(
    "[%s] --> %s %s %d %s\n",
    time.Now().Format("2006-01-02 15:04:05"),
    c.Method,
    c.URL,
    time.Since(start).Nanoseconds()/1e6,
    "ms",
  )
}

func main() {
  app := goa.New()
  r := router.New()

  r.GET("/", func(c *goa.Context) {
    c.String("Hello Goa!")
  })

  app.Use(logger)
  app.Use(r.Routes())
  log.Fatal(app.Listen(":3000"))
}
```

If you are unwilling to use goa-router, you can make a custom router middleware as you like.

## Maintainers

[@NicholasCao](https://github.com/NicholasCao).

## License

[MIT](https://github.com/goa-go/goa/blob/master/LICENSE)
