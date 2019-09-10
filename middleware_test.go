package goa

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

var calls = []int{}

func testMiddlewareServer() *httptest.Server {
	app := New()

	app.Use(func(c *Context, next func()) {
		calls = append(calls, 1)
		next()
		calls = append(calls, 6)
	})

	app.Use(func(c *Context, next func()) {
		calls = append(calls, 2)
		next()
		calls = append(calls, 5)
	})

	app.Use(func(c *Context, next func()) {
		calls = append(calls, 3)
		next()
		calls = append(calls, 4)
	})

	// Before testing, must compose middlewares.
	app.ComposeMiddlewares()
	return httptest.NewServer(app)
}

func TestMiddleware(t *testing.T) {
	server := testMiddlewareServer()
	defer server.Close()

	http.Get(server.URL)

	assert.Equal(t, calls, []int{1, 2, 3, 4, 5, 6})
}
