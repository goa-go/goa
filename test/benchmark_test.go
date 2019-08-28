package goa_test

import (
	"net/http"
	"testing"

	"github.com/goa-go/goa"
)

func BenchmarkGoa(b *testing.B) {
	app := goa.New()
	app.Use(func(c *goa.Context, next func()) {
	})

	app.ComposeMiddlewares()
	run(b, app)
}

func BenchmarkGoaMiddleware(b *testing.B) {
	app := goa.New()
	app.Use(func(c *goa.Context, next func()) {
		next()
	})
	app.Use(func(c *goa.Context, next func()) {
		next()
	})
	app.Use(func(c *goa.Context, next func()) {
	})

	app.ComposeMiddlewares()
	run(b, app)
}

func BenchmarkGoaString(b *testing.B) {
	app := goa.New()
	app.Use(func(c *goa.Context, next func()) {
		c.String("string")
	})

	app.ComposeMiddlewares()
	run(b, app)
}

type obj struct {
	key string
}

func BenchmarkGoaJSON(b *testing.B) {
	app := goa.New()
	app.Use(func(c *goa.Context, next func()) {
		c.JSON(obj{"value"})
	})

	app.ComposeMiddlewares()
	run(b, app)
}

type writer struct {
	headers http.Header
}

func (m *writer) Header() (h http.Header) {
	return m.headers
}

func (m *writer) Write(p []byte) (n int, err error) {
	return len(p), nil
}

func (m *writer) WriteString(s string) (n int, err error) {
	return len(s), nil
}

func (m *writer) WriteHeader(int) {}

func run(b *testing.B, app *goa.Goa) {
	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		panic(err)
	}

	w := &writer{
		http.Header{},
	}
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		app.ServeHTTP(w, req)
	}
}
