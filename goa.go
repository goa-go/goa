// goa
package goa

import (
	"net/http"
)

const Version = "1.0.0"

type Middleware func(*Context, func())
type Middlewares []Middleware

type HandleRequest func(*Context)

type Goa struct {
	middlewares Middlewares

	Context       *Context
	handleRequest HandleRequest
}

func (app *Goa) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	app.Context = createContext(w, r)
	app.handleRequest(app.Context)
}

func (app *Goa) Use(m Middleware) {
	app.middlewares = append(app.middlewares, m)
}

func (app *Goa) Listen(addr string) {
	app.handleRequest = compose(app.middlewares)
	http.ListenAndServe(addr, app)
}

func compose(m Middlewares) HandleRequest {
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

func New() *Goa {
	return &Goa{}
}
