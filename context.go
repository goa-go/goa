package goa

import (
	"fmt"
	"net/http"
	"net/url"

	"github.com/goa-go/goa/parser"
	"github.com/goa-go/goa/responser"
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
	Header http.Header

	queryMap url.Values
	Params   Params
	Keys     map[string]interface{}

	// Response status code.
	Status int

	// Content-Type
	Type string

	// Body will be wrote in response,
	// use it just like `c.Body = ...`.
	// Only string and struct will be supported.
	// If body is struct, will parse it as json.
	// If u want to respond xml or other type data, u can `c.XML(...)` or encode it into string.
	Body interface{}

	responser responser.Responser
}

func createContext(w http.ResponseWriter, r *http.Request) *Context {

	return &Context{
		Request:        r,
		ResponseWriter: w,
		Method:         r.Method,
		URL:            r.URL,
		Path:           r.URL.Path,
		Header:         r.Header,
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

/* handle request */

// Query returns the keyed url query value or ""
func (c *Context) Query(key string) string {
	query, _ := c.GetQuery(key)
	return query
}

// GetQuery returns the keyed url query value and isExit
// if it exists, return (value, true)
// otherwise it returns ("", false)
func (c *Context) GetQuery(key string) (string, bool) {
	if querys, ok := c.GetQueryArray(key); ok {
		return querys[0], true
	}
	return "", false
}

// GetQueryArray returns a slice of value for a given query key.
// And returns whether at least one value exists for the given key.
func (c *Context) GetQueryArray(key string) ([]string, bool) {
	c.initQuery()
	if querys, ok := c.queryMap[key]; ok && len(querys) > 0 {
		return querys, true
	}
	return []string{}, false
}

func (c *Context) initQuery() {
	if c.queryMap == nil {
		c.queryMap = make(url.Values)
		c.queryMap, _ = url.ParseQuery(c.Request.URL.RawQuery)
	}
}

// PostForm returns the value from a POST form or "".
func (c *Context) PostForm(key string) string {
	return c.Request.PostFormValue(key)
}

// Param returns the value of the URL param or "".
// When using goa-router, it works.
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

// handle parser
func (c *Context) Parse(p parser.Parser) {
	if err := p.Parse(c.Request); err != nil {
		panic(err)
	}
}

// parse json-data, require a pointer
func (c *Context) ParseJSON(pointer interface{}) {
	c.Parse(parser.JSON{Pointer: pointer})
}

// parse xml-data, require a pointer
func (c *Context) ParseXML(pointer interface{}) {
	c.Parse(parser.XML{Pointer: pointer})
}

// ParseString returns string-data
func (c *Context) ParseString() string {
	str, err := parser.String{}.Parse(c.Request)

	if err != nil {
		panic(err)
	}

	return str
}

/* handle response */

// status sets the HTTP response code.
func (c *Context) status() *Context {
	code := c.Status
	if code < 100 || code > 999 {
		panic(fmt.Sprintf("invalid status code: %d", code))
	}
	c.ResponseWriter.WriteHeader(code)
	return c
}

// M is a convenient alias for a map[string]interface{} map.
// Use is as c.JSON(&goa.M{...})
type M map[string]interface{}

func (c *Context) Respond(r responser.Responser) {
	if err := r.Respond(c.ResponseWriter); err != nil {
		panic(err)
	}
}

// respond json-data
func (c *Context) JSON(json interface{}) {
	// writeContentType(c.ResponseWriter, []string{"application/json; charset=utf-8"})
	c.Type = "application/json; charset=utf-8"
	// c.Respond(responser.JSON{Data: json})
	c.responser = responser.JSON{Data: json}
}

// respond xml-data
func (c *Context) XML(xml interface{}) {
	// writeContentType(c.ResponseWriter, []string{"application/xml; charset=utf-8"})
	c.Type = "application/xml; charset=utf-8"
	// c.Respond(responser.XML{Data: xml})
	c.responser = responser.XML{Data: xml}
}

// respond string-data
func (c *Context) String(str string) {
	// writeContentType(c.ResponseWriter, []string{"text/plain; charset=utf-8"})
	c.Type = "text/plain; charset=utf-8"
	// c.Respond(responser.String{Data: str})
	c.responser = responser.String{Data: str}
}

// respond html
func (c *Context) HTML(str string) {
	c.Type = "text/html; charset=utf-8"
	c.responser = responser.String{Data: str}
}

// redirect
func (c *Context) Redirect(code int, url string) {
	http.Redirect(c.ResponseWriter, c.Request, url, code)
}

// Set http response header.
// It should be called before Status and Respond.
func (c *Context) SetHeader(key string, value string) {
	c.ResponseWriter.Header().Set(key, value)
}

func writeContentType(w http.ResponseWriter, contentType []string) {
	header := w.Header()
	if val := header["Content-Type"]; len(val) == 0 {
		header["Content-Type"] = contentType
	}
}

func (c *Context) Error(msg string) {
	panic(msg)
}
