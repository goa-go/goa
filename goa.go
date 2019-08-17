package goa

import (
	"fmt"
	"net/http"
	"reflect"
	"regexp"

	"github.com/goa-go/goa/responser"
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

	Context           *Context
	middlewareHandler middlewareHandler
}

// New returns the initialized Goa instance.
func New() *Goa {
	return &Goa{}
}

func (app *Goa) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	app.Context = createContext(w, r)
	app.handleRequest(app.Context, app.middlewareHandler)
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
func (app *Goa) Listen(addr string) {
	app.ComposeMiddlewares()
	http.ListenAndServe(addr, app)
}

func compose(m Middlewares) middlewareHandler {
	return func(c *Context) {
		var dispatch func(int)
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

func (app *Goa) handleRequest(c *Context, fn middlewareHandler) {
	defer func() {
		if err := recover(); err != nil {
			app.onerror(err)
		}
	}()

	fn(c)
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
		c.SetHeader("Content-Type", c.Type)
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

func (app *Goa) onerror(err interface{}) {
	c := app.Context

	c.status(c.ErrorStatusCode)
	fmt.Fprint(c.ResponseWriter, err)
}
