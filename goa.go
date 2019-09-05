package goa

import (
	"log"
	"net/http"
	"sync"

	"github.com/goa-go/goa/responser"
	"github.com/pkg/errors"
)

// Middleware is based part of goa,
// any processing will take place here.
// should be used liked app.Use(middleware).
type Middleware func(*Context, func())

// Middlewares is []Middleware.
type Middlewares []Middleware

type middlewareHandler func(*Context)

// Goa is the framework's instance.
type Goa struct {
	middlewares Middlewares

	pool              sync.Pool
	middlewareHandler middlewareHandler
}

// New returns the initialized Goa instance.
func New() *Goa {
	app := &Goa{}
	app.pool.New = func() interface{} {
		return &Context{}
	}
	return app
}

// ServeHTTP makes the router implement the http.Handler interface.
func (app *Goa) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	c := app.pool.Get().(*Context)
	c.init(w, r)

	app.handleRequest(c)

	app.pool.Put(c)
}

// Use a middleware.
func (app *Goa) Use(m Middleware) {
	app.middlewares = append(app.middlewares, m)
}

// ComposeMiddlewares composes middleware,
// it doesn't need to be called manually except in testing,
// but httptest is only available after ComposeMiddlewares is called.
func (app *Goa) ComposeMiddlewares() {
	app.middlewareHandler = compose(app.middlewares)
}

// Listen starts server with the addr.
func (app *Goa) Listen(addr string) error {
	app.ComposeMiddlewares()
	return http.ListenAndServe(addr, app)
}

var dispatch func(i int)

func compose(m Middlewares) middlewareHandler {
	return func(c *Context) {
		dispatch = func(i int) {
			if i == len(m) {
				return
			}
			fn := m[i]
			fn(c, func() {
				dispatch(i + 1)
			})
		}
		dispatch(0)
	}
}

func (app *Goa) handleRequest(c *Context) {
	defer func() {
		if err := recover(); err != nil {
			app.onerror(c, err)
		}
	}()

	app.middlewareHandler(c)
	if !c.redirected && !c.Handled {
		app.handleResponse(c)
	}
}

func (app *Goa) handleResponse(c *Context) {

	// Content-Type
	if c.Type != "" {
		c.writeContentType(c.Type)
	}

	// Status code
	c.ResponseWriter.WriteHeader(c.status)

	// Response
	if c.responser == nil {
		c.String(http.StatusText(c.status))
	}

	if err := c.respond(c.responser); err != nil {
		log.Printf("[ERROR] %+v", errors.WithStack(err))
		c.respond(responser.String{Data: http.StatusText(http.StatusInternalServerError)})
	}
}

func (app *Goa) onerror(c *Context, err interface{}) {
	code := http.StatusInternalServerError
	msg := http.StatusText(http.StatusInternalServerError)

	if e, ok := err.(Error); ok {
		code = e.Code
		msg = e.Msg
	} else if e, ok := err.(error); ok {
		log.Printf("[ERROR] %+v", errors.WithStack(e))
		msg = e.Error()
	} else if str, ok := err.(string); ok {
		log.Print("[ERROR] ", str)
		msg = str
	} else {
		log.Print("[ERROR] ", err)
	}

	c.Type = "text/plain; charset=utf-8"
	c.writeContentType(c.Type)
	c.SetHeader("X-Content-Type-Options", "nosniff")

	c.ResponseWriter.WriteHeader(code)
	c.respond(responser.String{Data: msg})
}
