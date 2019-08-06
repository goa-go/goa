// Copyright 2013 Julien Schmidt. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be found
// in the LICENSE file.

package router

import (
	"net/http"

	"github.com/goa-go/goa"
)

type Handler func(*goa.Context)

// Router is a http.Handler which can be used to dispatch requests to different
// handler functions via configurable routes
type Router struct {
	trees map[string]*node

	// Enables automatic redirection if the current route can't be matched but a
	// handler for the path with (without) the trailing slash exists.
	// For example if /foo/ is requested but a route only exists for /foo, the
	// client is redirected to /foo with http status code 301 for GET requests
	// and 307 for all other request methods.
	RedirectTrailingSlash bool

	// If enabled, the router tries to fix the current request path, if no
	// handler is registered for it.
	// First superfluous path elements like ../ or // are removed.
	// Afterwards the router does a case-insensitive lookup of the cleaned path.
	// If a handler can be found for this route, the router makes a redirection
	// to the corrected path with status code 301 for GET requests and 307 for
	// all other request methods.
	// For example /FOO and /..//Foo could be redirected to /foo.
	// RedirectTrailingSlash is independent of this option.
	RedirectFixedPath bool

	// If enabled, the router checks if another method is allowed for the
	// current route, if the current request can not be routed.
	// If this is the case, the request is answered with 'Method Not Allowed'
	// and HTTP status code 405.
	// If no other Method is allowed, the request is delegated to the NotFound
	// handler.
	HandleMethodNotAllowed bool

	// If enabled, the router automatically replies to OPTIONS requests.
	// Custom OPTIONS handlers take priority over automatic replies.
	HandleOPTIONS bool

	// Configurable http.Handler which is called when no matching route is
	// found. If it is not set, http.NotFound is used.
	NotFound http.Handler

	// Configurable http.Handler which is called when a request
	// cannot be routed and HandleMethodNotAllowed is true.
	// If it is not set, http.Error with http.StatusMethodNotAllowed is used.
	// The "Allow" header with allowed request methods is set before the handler
	// is called.
	MethodNotAllowed http.Handler

	// Function to handler panics recovered from http handlers.
	// It should be used to generate a error page and return the http error code
	// 500 (Internal Server Error).
	// The handler can be used to keep your server from crashing because of
	// unrecovered panics.
	PanicHandler func(http.ResponseWriter, *http.Request, interface{})
}

// Make sure the Router conforms with the http.Handler interface
// var _ http.Handler = New()

// New returns a new initialized Router.
// Path auto-correction, including trailing slashes, is enabled by default.
func New() *Router {
	return &Router{
		RedirectTrailingSlash:  true,
		RedirectFixedPath:      true,
		HandleMethodNotAllowed: true,
		HandleOPTIONS:          true,
	}
}

// GET is a shortcut for router.Handler("GET", path, handler)
func (r *Router) GET(path string, handler Handler) {
	r.Handle("GET", path, handler)
}

// HEAD is a shortcut for router.Handler("HEAD", path, handler)
func (r *Router) HEAD(path string, handler Handler) {
	r.Handle("HEAD", path, handler)
}

// OPTIONS is a shortcut for router.Handler("OPTIONS", path, handler)
func (r *Router) OPTIONS(path string, handler Handler) {
	r.Handle("OPTIONS", path, handler)
}

// POST is a shortcut for router.Handler("POST", path, handler)
func (r *Router) POST(path string, handler Handler) {
	r.Handle("POST", path, handler)
}

// PUT is a shortcut for router.Handler("PUT", path, handler)
func (r *Router) PUT(path string, handler Handler) {
	r.Handle("PUT", path, handler)
}

// PATCH is a shortcut for router.Handler("PATCH", path, handler)
func (r *Router) PATCH(path string, handler Handler) {
	r.Handle("PATCH", path, handler)
}

// DELETE is a shortcut for router.Handler("DELETE", path, handler)
func (r *Router) DELETE(path string, handler Handler) {
	r.Handle("DELETE", path, handler)
}

// Handle registers a new request handle with the given path and method.
//
// For GET, POST, PUT, PATCH and DELETE requests the respective shortcut
// functions can be used.
//
// This function is intended for bulk loading and to allow the usage of less
// frequently used, non-standardized or custom methods (e.g. for internal
// communication with a proxy).
func (r *Router) Handle(method, path string, handle Handler) {
	if path[0] != '/' {
		panic("path must begin with '/' in path '" + path + "'")
	}

	if r.trees == nil {
		r.trees = make(map[string]*node)
	}

	root := r.trees[method]
	if root == nil {
		root = new(node)
		r.trees[method] = root
	}

	root.addRoute(path, handle)
}

func (r *Router) recv(w http.ResponseWriter, req *http.Request) {
	if rcv := recover(); rcv != nil {
		r.PanicHandler(w, req, rcv)
	}
}

func (r *Router) allowed(path, reqMethod string) (allow string) {
	if path == "*" { // server-wide
		for method := range r.trees {
			if method == "OPTIONS" {
				continue
			}

			// add request method to list of allowed methods
			if len(allow) == 0 {
				allow = method
			} else {
				allow += ", " + method
			}
		}
	} else { // specific path
		for method := range r.trees {
			// Skip the requested method - we already tried this one
			if method == reqMethod || method == "OPTIONS" {
				continue
			}

			handle, _, _ := r.trees[method].getValue(path)
			if handle != nil {
				// add request method to list of allowed methods
				if len(allow) == 0 {
					allow = method
				} else {
					allow += ", " + method
				}
			}
		}
	}
	if len(allow) > 0 {
		allow += ", OPTIONS"
	}
	return
}

// ServeHTTP makes the router implement the http.Handler interface.
func (r *Router) Routes() goa.Middleware {
	return func(c *goa.Context, next func()) {
		if r.PanicHandler != nil {
			defer r.recv(c.ResponseWriter, c.Request)
		}

		path := c.Path

		if root := r.trees[c.Method]; root != nil {
			if Handler, ps, tsr := root.getValue(path); Handler != nil {
				c.Params = ps
				Handler(c)
				return
			} else if c.Method != "CONNECT" && path != "/" {
				code := 301 // Permanent redirect, request with GET method
				if c.Method != "GET" {
					// Temporary redirect, request with same method
					// As of Go 1.3, Go does not support status code 308.
					code = 307
				}

				if tsr && r.RedirectTrailingSlash {
					if len(path) > 1 && path[len(path)-1] == '/' {
						c.Path = path[:len(path)-1]
					} else {
						c.Path = path + "/"
					}
					c.Redirect(code, c.URL.String())
					return
				}

				// Try to fix the request path
				if r.RedirectFixedPath {
					fixedPath, found := root.findCaseInsensitivePath(
						CleanPath(path),
						r.RedirectTrailingSlash,
					)
					if found {
						c.Path = string(fixedPath)
						c.Redirect(code, c.URL.String())
						return
					}
				}
			}
		}

		if c.Method == "OPTIONS" && r.HandleOPTIONS {
			// Handle OPTIONS requests
			if allow := r.allowed(path, c.Method); len(allow) > 0 {
				c.ResponseWriter.Header().Set("Allow", allow)
				return
			}
		} else {
			// Handle 405
			if r.HandleMethodNotAllowed {
				if allow := r.allowed(path, c.Method); len(allow) > 0 {
					c.ResponseWriter.Header().Set("Allow", allow)
					if r.MethodNotAllowed != nil {
						r.MethodNotAllowed.ServeHTTP(c.ResponseWriter, c.Request)
					} else {
						http.Error(c.ResponseWriter,
							http.StatusText(http.StatusMethodNotAllowed),
							http.StatusMethodNotAllowed,
						)
					}
					return
				}
			}
		}

		// Handle 404
		if r.NotFound != nil {
			r.NotFound.ServeHTTP(c.ResponseWriter, c.Request)
		} else {
			http.NotFound(c.ResponseWriter, c.Request)
		}

		next()
	}
}
