package goa

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func BenchmarkGoa(b *testing.B) {
	app := New()
	app.Use(func(c *Context, next func()) {
	})

	app.ComposeMiddlewares()
	run(b, app)
}

func BenchmarkGoaMiddleware(b *testing.B) {
	app := New()
	app.Use(func(c *Context, next func()) {
		next()
	})
	app.Use(func(c *Context, next func()) {
		next()
	})
	app.Use(func(c *Context, next func()) {
	})

	app.ComposeMiddlewares()
	run(b, app)
}

func BenchmarkGoaString(b *testing.B) {
	app := New()
	app.Use(func(c *Context, next func()) {
		c.String("string")
	})

	app.ComposeMiddlewares()
	run(b, app)
}

func BenchmarkGoaJSON(b *testing.B) {
	type obj struct {
		Key string `json:"key"`
	}

	app := New()
	app.Use(func(c *Context, next func()) {
		c.JSON(obj{"value"})
	})

	app.ComposeMiddlewares()
	run(b, app)
}

func run(b *testing.B, app *Goa) {
	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		panic(err)
	}

	w := httptest.NewRecorder()
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		app.ServeHTTP(w, req)
	}
}
