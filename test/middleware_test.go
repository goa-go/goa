package goa_test

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/goa-go/goa"
)

func TestMiddleware(t *testing.T) {
	app := goa.New()
	calls := []int{}

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

	go app.Listen(":3000")

	http.Get("http://localhost:3000")

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
