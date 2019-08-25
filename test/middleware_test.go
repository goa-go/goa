package goa_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/goa-go/goa"
)

var calls = []int{}

func testMiddlewareServer() *httptest.Server {
	app := goa.New()

	app.Use(func(c *goa.Context, next func()) {
		calls = append(calls, 1)
		next()
		calls = append(calls, 6)
	})

	app.Use(func(c *goa.Context, next func()) {
		calls = append(calls, 2)
		next()
		calls = append(calls, 5)
	})

	app.Use(func(c *goa.Context, next func()) {
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

	if len(calls) != 6 {
		fmt.Println(calls)
		t.Error("middleware error")
	}

	for i, v := range calls {
		if v != i+1 {
			t.Error("middleware error")
		}
	}
}
