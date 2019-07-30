package goa

import (
	"net/http"
	"strings"

	"github.com/julienschmidt/httprouter"
)

type Router struct {
	methods    []string
	httprouter *httprouter.Router
}

func NewRouter() *Router {
	r := &Router{
		[]string{"GET", "POST", "PUT", "DELETE", "PATCH", "OPTIONS", "HEAD"},
		nil,
	}
	r.httprouter = httprouter.New()
	return r
}

func (router *Router) Register(method string, path string, handle httprouter.Handle) {
	isMethod := false
	for _, v := range router.methods {
		if strings.EqualFold(method, v) {
			isMethod = true
		}
	}

	if !isMethod {
		return
	} else {
		router.httprouter.Handle(method, path, handle)
	}
}
func (router *Router) GET(path string, handle httprouter.Handle) {
	router.Register("GET", path, handle)
}
func (r *Router) POST(path string, handler httprouter.Handle) {
	r.Register("POST", path, handler)
}
func (r *Router) Put(path string, handler httprouter.Handle) {
	r.Register("PUT", path, handler)
}
func (r *Router) Delete(path string, handler httprouter.Handle) {
	r.Register("DELETE", path, handler)
}
func (r *Router) Patch(path string, handler httprouter.Handle) {
	r.Register("PATCH", path, handler)
}
func (r *Router) Options(path string, handler httprouter.Handle) {
	r.Register("OPTIONS", path, handler)
}
func (r *Router) Head(path string, handler httprouter.Handle) {
	r.Register("HEAD", path, handler)
}

func (router *Router) Routes() Middleware {
	return func(w http.ResponseWriter, r *http.Request, next func()) {
		router.httprouter.ServeHTTP(w, r)
		next()
	}
}
