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

// Context is used to receive requests and respond to requests.
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

	// Error Status code.
	errorStatusCode int

	// Content-Type
	Type string

	redirected bool

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
		Request:         r,
		ResponseWriter:  w,
		Method:          r.Method,
		URL:             r.URL,
		Path:            r.URL.Path,
		Header:          r.Header,
		errorStatusCode: 500,
	}
}

// Set value.
func (c *Context) Set(key string, value interface{}) {
	if c.Keys == nil {
		c.Keys = make(map[string]interface{})
	}
	c.Keys[key] = value
}

// Get value, return (value, exists).
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

// Parse handles parser.
func (c *Context) Parse(p parser.Parser) {
	if err := p.Parse(c.Request); err != nil {
		panic(err)
	}
}

// ParseJSON parses json-data, require a pointer.
func (c *Context) ParseJSON(pointer interface{}) {
	c.Parse(parser.JSON{Pointer: pointer})
}

// ParseXML parses xml-data, require a pointer.
func (c *Context) ParseXML(pointer interface{}) {
	c.Parse(parser.XML{Pointer: pointer})
}

// ParseString returns string-data
func (c *Context) ParseString() string {
	str := parser.String{}.Parse(c.Request)

	return str
}

/* handle response */

// status sets the HTTP response code.
func (c *Context) status(code int) {
	if code < 100 || code > 999 {
		panic(fmt.Errorf("invalid status code: %d", code))
	}
	c.ResponseWriter.WriteHeader(code)
}

// M is a convenient alias for a map[string]interface{} map.
// Use is as c.JSON(&goa.M{...})
type M map[string]interface{}

func (c *Context) respond(r responser.Responser) {
	if err := r.Respond(c.ResponseWriter); err != nil {
		panic(err)
	}
}

// JSON responds json-data.
func (c *Context) JSON(json interface{}) {
	c.Type = "application/json; charset=utf-8"
	c.responser = responser.JSON{Data: json}
}

// XML responds xml-data.
func (c *Context) XML(xml interface{}) {
	c.Type = "application/xml; charset=utf-8"
	c.responser = responser.XML{Data: xml}
}

// String responds string-data.
func (c *Context) String(str string) {
	c.Type = "text/plain; charset=utf-8"
	c.responser = responser.String{Data: str}
}

// HTML responds html.
func (c *Context) HTML(str string) {
	c.Type = "text/html; charset=utf-8"
	c.responser = responser.String{Data: str}
}

// Redirect replies to the request with a redirect to url and a status code.
func (c *Context) Redirect(code int, url string) {
	c.redirected = true
	http.Redirect(c.ResponseWriter, c.Request, url, code)
}

// SetHeader sets http response header.
// It should be called before Status and Respond.
func (c *Context) SetHeader(key string, value string) {
	c.ResponseWriter.Header().Set(key, value)
}

// Error is used like c.Error(goa.Error{...}).
// It will create a http-error.
type Error struct {
	Msg    string
	Status int
}

// Error throw a http-error.
func (c *Context) Error(err Error) {
	if err.Status == 0 {
		err.Status = 500
	}

	panic(err)
}
