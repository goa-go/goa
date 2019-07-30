package goa

import (
	"net/http"
)

type Middleware func(http.ResponseWriter, *http.Request, func())
type Middlewares []Middleware

type HandleRequest func(http.ResponseWriter, *http.Request)

type Goa struct {
	middlewares   Middlewares
	handleRequest HandleRequest
}

func (app *Goa) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	app.handleRequest(w, r)
}

func (app *Goa) Use(m Middleware) {
	app.middlewares = append(app.middlewares, m)
}

func (app *Goa) Listen(addr string) {
	app.handleRequest = compose(app.middlewares)
	http.ListenAndServe(addr, app)
}

func compose(m Middlewares) HandleRequest {
	return func(w http.ResponseWriter, r *http.Request) {
		var dispatch func(i int)
		dispatch = func(i int) {
			if i == len(m) {
				return
			}
			fn := m[i]
			fn(w, r, func() {
				dispatch(i + 1)
			})
		}

		dispatch(0)
	}
}

func New() *Goa {
	return &Goa{}
}
