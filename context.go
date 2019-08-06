package goa

import (
	"net/http"
	"net/url"

	"github.com/goa-go/goa/encode"
)

// Param is a single URL parameter, consisting of a key and a value.
type Param struct {
	Key   string
	Value string
}

// Params is a Param-slice, as returned by the router.
// The slice is ordered, the first URL parameter is also the first slice value.
// It is therefore safe to read values by the index.
type Params []Param

type Context struct {
	Request        *http.Request
	ResponseWriter http.ResponseWriter

	Method string
	URL    *url.URL
	Path   string

	Params Params
	Keys   map[string]interface{}
}

func createContext(w http.ResponseWriter, r *http.Request) *Context {
	return &Context{
		Request:        r,
		ResponseWriter: w,
		Method:         r.Method,
		URL:            r.URL,
		Path:           r.URL.Path,
	}
}

// Context set value.
func (c *Context) Set(key string, value interface{}) {
	if c.Keys == nil {
		c.Keys = make(map[string]interface{})
	}
	c.Keys[key] = value
}

// Context get value, return (value, exists).
func (c *Context) Get(key string) (value interface{}, exists bool) {
	value, exists = c.Keys[key]
	return
}

func (c *Context) Param(key string) string {
	return c.Params.Get(key)
}

// Get returns the value of the first Param which key matches the given name.
// If no matching Param is found, an empty string is returned.
func (ps Params) Get(name string) string {
	for _, param := range ps {
		if param.Key == name {
			return param.Value
		}
	}
	return ""
}

// Status sets the HTTP response code.
// And return context, so c.Status(200).JSON(...) is supported.
func (c *Context) Status(code int) *Context {

	c.ResponseWriter.WriteHeader(code)

	return c
}

func (c *Context) Render(data []byte) {
	_, err := c.ResponseWriter.Write(data)
	if err != nil {
		panic(err)
	}
}

// M is a convenient alias for a map[string]interface{} map.
// Use is as c.JSON(&goa.M{...})
type M map[string]interface{}

func (c *Context) JSON(json interface{}) {
	writeContentType(c.ResponseWriter, []string{"application/json; charset=utf-8"})
	c.Render(encode.JSON(json))
}

func (c *Context) String(str string) {
	writeContentType(c.ResponseWriter, []string{"text/plain; charset=utf-8"})
	c.Render(encode.String(str))
}

func (c *Context) XML(xml interface{}) {
	writeContentType(c.ResponseWriter, []string{"application/xml; charset=utf-8"})
	c.Render(encode.XML(xml))
}

func (c *Context) Redirect(code int, location string) {
	http.Redirect(c.ResponseWriter, c.Request, location, code)
}

func writeContentType(w http.ResponseWriter, contentType []string) {
	w.Header()["Content-Type"] = contentType
}
