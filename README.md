# Goa

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

  app.Use(func(w http.ResponseWriter, r *http.Request, next func()) {
    w.Write([]byte("Hello Goa!"))
    next()
  })
  app.Listen(":3000")
}
```

## Middleware

Goa is a web framework based on middleware. Here is an example.
```go
import (
  "fmt"
  "net/http"
  "time"

  "github.com/goa-go/goa"
  "github.com/julienschmidt/httprouter"
)

func logger(w http.ResponseWriter, r *http.Request, next func()) {
  start := time.Now()
  next()
  fmt.Println(r.Method, r.URL, time.Since(start).Nanoseconds()/1e6, "ms")
}

func main() {
  app := goa.New()
  router := goa.NewRouter()

  router.GET("/", func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
    fmt.Fprintln(w, "Hello Goa!")
  })
  app.Use(logger)
  app.Use(router.Routes())
  app.Listen(":3000")
}
```

## License

[MIT](https://github.com/goa-go/goa/blob/master/LICENSE)
