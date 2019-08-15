// goa
package goa

import (
	"fmt"
	"net/http"
	"reflect"
	"regexp"

	"github.com/goa-go/goa/responser"
)

type Middleware func(*Context, func())
type Middlewares []Middleware

type middlewareHandler func(*Context)

type Goa struct {
	middlewares Middlewares

	Context           *Context
	middlewareHandler middlewareHandler
}

// Init goa.
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

// Compose middleware,
// httptest is only available after ComposeMiddlewares is called.
func (app *Goa) ComposeMiddlewares() {
	app.middlewareHandler = compose(app.middlewares)
}

// Start server with addr.
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
	app.handleResponse(c)
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
	c.status()

	// Response
	if c.responser != nil {
		c.Respond(c.responser)
		return
	}
}

func (app *Goa) onerror(err interface{}) {
	c := app.Context

	if c.Status != 0 {
		c.status()
	} else {
		c.Status = 500
		c.status()
	}

	fmt.Fprint(c.ResponseWriter, err)

	fmt.Println(err)
}
