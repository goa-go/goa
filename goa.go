package goa

import (
	"fmt"
	"log"
	"net/http"
	"reflect"
	"regexp"
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
	if !c.redirected {
		app.handleResponse(c)
	}
}

func (app *Goa) handleResponse(c *Context) {
	body := c.Body

	// handle body
	if body != nil {
		if str, ok := body.(string); ok {
			if match, _ := regexp.MatchString(`^\s*<`, str); match {
				c.Type = "text/html; charset=utf-8"
			} else {
				c.Type = "text/plain; charset=utf-8"
			}
			c.responser = responser.String{Data: str}
		} else if reflect.TypeOf(body).Kind() == reflect.Struct || reflect.TypeOf(body).Kind() == reflect.Map {
			c.Type = "application/json; charset=utf-8"
			c.responser = responser.JSON{Data: body}
		}
	}

	// Content-Type
	if c.Type != "" {
		c.writeContentType(c.Type)
	}

	// Status code
	if c.Status == 0 {
		if c.responser == nil && body == nil {
			c.Status = 404
		} else {
			c.Status = 200
		}
	}
	c.status(c.Status)

	// Response
	if c.responser != nil {
		c.respond(c.responser)
		return
	}
}

func (app *Goa) onerror(c *Context, err interface{}) {
	var errResponse interface{}

	if e, ok := err.(Error); ok {
		c.errorStatusCode = e.Status
		errResponse = e.Msg
	} else if e, ok := err.(error); ok {
		log.Printf("[ERROR] %+v", errors.WithStack(e))
		errResponse = e.Error()
	} else {
		log.Print("[ERROR] ", err)
		errResponse = err
	}

	c.Type = "text/plain; charset=utf-8"
	c.writeContentType(c.Type)
	c.Status = c.errorStatusCode
	c.status(c.Status)
	if errResponse != nil && errResponse != "" {
		fmt.Fprint(c.ResponseWriter, errResponse)
	} else {
		fmt.Fprint(c.ResponseWriter, "Internal Server Error")
	}
}
