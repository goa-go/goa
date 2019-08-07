# Goa

[![CI](https://img.shields.io/travis/goa-go/goa.svg?style=flat)](https://travis-ci.org/goa-go/goa)
[![PR's Welcome](https://img.shields.io/badge/PRs-welcome-brightgreen.svg?style=flat)](https://github.com/goa-go/goa/pull/new)

Goa is under construction, if you are familiar with [koa](https://github.com/koajs/koa) or go and interested in this project, please join us.

## What is goa?
goa = go + [koa](https://github.com/koajs/koa)

## Installation

```bash
$ go get github.com/goa-go/goa
```

##  Hello Goa

```go
func main() {
  app := New()

  app.Use(func(c *goa.Context, next func()) {
    c.String("Hello Goa!")
    next()
  })
  app.Listen(":3000")
}
```

## Middleware

Goa is a web framework based on middleware. Here is an example of using goa-router and logger.
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

  fmt.Printf("[%s] <-- %s %s\n", start.Format("2006-6-2 15:04:05"), c.Method, c.URL)
  next()
  fmt.Printf("[%s] --> %s %s %d%s\n", time.Now().Format("2006-6-2 15:04:05"), c.Method, c.URL, time.Since(start).Nanoseconds()/1e6, "ms")
}

func json(c *goa.Context) {
  c.JSON(goa.M{
    "string": "string",
    "int":    1,
    "json": goa.M{
      "key": "value",
    },
  })
}

func main() {
  app := goa.New()
  router := router.New()

  router.GET("/", func(c *goa.Context) {
    c.String("hello world")
  })
  router.GET("/json", json)

  app.Use(logger)
  app.Use(router.Routes())
  app.Listen(":3000")
}
```

## License

[MIT](https://github.com/goa-go/goa/blob/master/LICENSE)
