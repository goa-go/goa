# Goa

[![Build Status](https://travis-ci.org/goa-go/goa.svg?branch=master)](https://travis-ci.org/goa-go/goa)
[![Codecov](https://codecov.io/gh/goa-go/goa/branch/master/graph/badge.svg)](https://codecov.io/github/goa-go/goa?branch=master)
[![Go Doc](https://godoc.org/github.com/goa-go/goa?status.svg)](http://godoc.org/github.com/goa-go/goa)
[![Go Report](https://goreportcard.com/badge/github.com/goa-go/goa)](https://goreportcard.com/report/github.com/goa-go/goa)
[![PR's Welcome](https://img.shields.io/badge/PRs-welcome-brightgreen.svg?style=flat)](https://github.com/goa-go/goa/pull/new)

Goa is under construction, if you are familiar with [koa](https://github.com/koajs/koa) or go and interested in this project, please join us.

## What is goa?
goa = go + [koa](https://github.com/koajs/koa)

Just like koa, goa is also not bundled with any middleware. But you can expand functionality to meet your needs at will by middlware. It is flexible, light, high-performance and extensible.

## Installation

```bash
$ go get github.com/goa-go/goa
```

##  Hello Goa

```go
func main() {
  app := goa.New()

  app.Use(func(c *goa.Context, next func()) {
    c.String("Hello Goa!")
    next()
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
  "time"

  "github.com/goa-go/goa"
  "github.com/goa-go/goa/router"
)

func logger(c *goa.Context, next func()) {
  start := time.Now()

  fmt.Printf(
    "[%s] <-- %s %s\n",
    start.Format("2006-6-2 15:04:05"),
    c.Method,
    c.URL,
  )
  next()
  fmt.Printf(
    "[%s] --> %s %s %d %s\n",
    time.Now().Format("2006-6-2 15:04:05"),
    c.Method,
    c.URL,
    time.Since(start).Nanoseconds()/1e6,
    "ms",
  )
}

func main() {
  app := goa.New()
  router := router.New()

  router.GET("/", func(c *goa.Context) {
    c.String("Hello Goa!")
  })

  app.Use(logger)
  app.Use(router.Routes())
  log.Fatal(app.Listen(":3000"))
}
```

If you are unwilling to use goa-router, you can make a custom router middleware as you like.

## License

[MIT](https://github.com/goa-go/goa/blob/master/LICENSE)
