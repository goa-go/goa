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

	status         int
	explicitStatus bool

	// Content-Type
	ct string

	// whether handled response by c.ResponseWriter
	Handled    bool
	redirected bool

	index int8
	len   int8
	app   *Goa

	responser responser.Responser
}

func (c *Context) init(w http.ResponseWriter, r *http.Request) {
	c.Request = r
	c.ResponseWriter = w
	c.Method = r.Method
	c.URL = r.URL
	c.Path = r.URL.Path
	c.Header = r.Header
	c.status = http.StatusNotFound
	c.explicitStatus = false

	c.Params = nil
	c.Keys = nil
	c.queryMap = nil
	c.ct = ""
	c.Handled = false
	c.redirected = false
	c.responser = nil
	c.index = 0
	c.len = int8(len(c.app.middlewares))
}

// Next implements the next middleware.
// For example,
// app.Use(func(c *goa.Context) {
//   //do sth
//   c.Next()
//   //do sth
// })
func (c *Context) Next() {
	if c.index >= c.len-1 {
		return
	}
	c.index++
	c.app.middlewares[c.index](c)
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

func (c *Context) parse(p parser.Parser) error {
	return p.Parse(c.Request)
}

// ParseJSON parses json-data, require a pointer.
func (c *Context) ParseJSON(pointer interface{}) error {
	return c.parse(parser.JSON{Pointer: pointer})
}

// ParseXML parses xml-data, require a pointer.
func (c *Context) ParseXML(pointer interface{}) error {
	return c.parse(parser.XML{Pointer: pointer})
}

// ParseString returns string-data
func (c *Context) ParseString() (string, error) {
	return parser.String{}.Parse(c.Request)
}

// ParseQuery can parse query, require a pointer.
// Just like json, it also needs a "query" tag. Here is a example.
//
// type Person struct {
// 	Name string `query:"name"`
// 	Age  int    `query:"age"`
// }
//
// p := &Person{}
// c.ParseQuery(p)
func (c *Context) ParseQuery(pointer interface{}) error {
	return c.parse(parser.Query{Pointer: pointer})
}

// ParseForm can parse form-data and x-www-form-urlencoded,
// the latter is not available when the request method is get,
// require a pointer.
// Just like json, it also needs a "form" tag. Here is a example.
//
// type Person struct {
// 	Name string `form:"name"`
// 	Age  int    `form:"age"`
// }
//
// p := &Person{}
// c.ParseForm(p)
func (c *Context) ParseForm(pointer interface{}) error {
	return c.parse(parser.Form{Pointer: pointer})
}

// Cookie returns the named cookie provided in the request 
// or ErrNoCookie if not found.
func (c *Context) Cookie(name string) (string, error) {
	cookie, err := c.Request.Cookie(name)
	if err != nil {
		return "", err
	}
	val, _ := url.QueryUnescape(cookie.Value)
	return val, nil
}

/* handle response */

// Status sets the HTTP response code.
func (c *Context) Status(code int) {
	if code < 100 || code > 999 {
		panic(fmt.Errorf("invalid status code: %d", code))
	}
	c.explicitStatus = true
	c.status = code
}

// GetStatus returns c.status.
func (c *Context) GetStatus() int {
	return c.status
}

// M is a convenient alias for a map[string]interface{} map.
// Use is as c.JSON(&goa.M{...})
type M map[string]interface{}

func (c *Context) respond(r responser.Responser) error {
	return r.Respond(c.ResponseWriter)
}

// JSON responds json-data.
func (c *Context) JSON(json interface{}) {
	if !c.explicitStatus {
		c.Status(http.StatusOK)
	}

	c.ct = "application/json; charset=utf-8"
	c.responser = responser.JSON{Data: json}
}

// XML responds xml-data.
func (c *Context) XML(xml interface{}) {
	if !c.explicitStatus {
		c.Status(http.StatusOK)
	}

	c.ct = "application/xml; charset=utf-8"
	c.responser = responser.XML{Data: xml}
}

// String responds string-data.
func (c *Context) String(str string) {
	if !c.explicitStatus {
		c.Status(http.StatusOK)
	}

	c.ct = "text/plain; charset=utf-8"
	c.responser = responser.String{Data: str}
}

// HTML responds html.
func (c *Context) HTML(html string) {
	if !c.explicitStatus {
		c.Status(http.StatusOK)
	}

	c.ct = "text/html; charset=utf-8"
	c.responser = responser.String{Data: html}
}

// Redirect replies to the request with a redirect to url and a status code.
func (c *Context) Redirect(code int, url string) {
	if code < http.StatusMultipleChoices || code > http.StatusPermanentRedirect {
		panic(fmt.Errorf("cannot redirect with status code %d", code))
	}
	c.redirected = true
	c.Status(code)
	http.Redirect(c.ResponseWriter, c.Request, url, code)
}

// SetHeader sets http response header.
// It should be called before Status and Respond.
func (c *Context) SetHeader(key string, value string) {
	c.ResponseWriter.Header().Set(key, value)
}

func (c *Context) writeContentType(value string) {
	header := c.ResponseWriter.Header()
	if val := header["Content-Type"]; len(val) == 0 {
		header["Content-Type"] = []string{value}
	}
}

// SetCookie adds a Set-Cookie header to the ResponseWriter's headers.
func (c *Context) SetCookie(cookie *http.Cookie) {
	http.SetCookie(c.ResponseWriter, cookie)
}

// Error is used like c.Error(goa.Error{...}).
// It will create a http-error.
type Error struct {
	Code int
	Msg  string
}

// Error throw a http-error, it would be catched by goa.
func (c *Context) Error(code int, msg string) {
	panic(Error{
		code,
		msg,
	})
}
