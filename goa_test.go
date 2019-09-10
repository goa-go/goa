package goa

import (
	"errors"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func testServer(m Middleware) *httptest.Server {
	app := New()

	app.Use(m)

	// Before testing, must compose middlewares.
	app.ComposeMiddlewares()
	return httptest.NewServer(app)
}

func TestSimpleRequest(t *testing.T) {
	ts := testServer(func(c *Context, next func()) {
		c.String("Hello Goa!")
	})
	defer ts.Close()
	resp, err := http.Get(ts.URL)
	assert.Nil(t, err)

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, "Hello Goa!", string(body))
}

func TestShouldNotHandleRespond(t *testing.T) {
	ts := testServer(func(c *Context, next func()) {
		c.Handled = true
		c.Status(500)
		c.String("error")
	})
	defer ts.Close()
	resp, err := http.Get(ts.URL)
	assert.Nil(t, err)

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, "", string(body))
}

func TestShouldNotHandleRespondWithRedirect(t *testing.T) {
	ts := testServer(func(c *Context, next func()) {
		if c.Path != "/path" {
			c.Redirect(301, "/path")
			c.Status(500)
			c.String("error")
		}
		c.String("redirected")
	})
	defer ts.Close()
	resp, err := http.Get(ts.URL)
	assert.Nil(t, err)

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, "redirected", string(body))
}

/* test *goa.onerror */
func TestGoaError(t *testing.T) {
	ts := testServer(func(c *Context, next func()) {
		c.Error(404, http.StatusText(404))
	})
	defer ts.Close()
	resp, err := http.Get(ts.URL)
	assert.Nil(t, err)

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	assert.Nil(t, err)
	assert.Equal(t, 404, resp.StatusCode)
	assert.Equal(t, http.StatusText(404), string(body))
}

func TestError(t *testing.T) {
	ts := testServer(func(c *Context, next func()) {
		panic(errors.New("error"))
	})
	defer ts.Close()
	resp, err := http.Get(ts.URL)
	assert.Nil(t, err)

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	assert.Nil(t, err)
	assert.Equal(t, 500, resp.StatusCode)
	assert.Equal(t, "error", string(body))
}

func TestStringError(t *testing.T) {
	ts := testServer(func(c *Context, next func()) {
		panic("error")
	})
	defer ts.Close()
	resp, err := http.Get(ts.URL)
	assert.Nil(t, err)

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	assert.Nil(t, err)
	assert.Equal(t, 500, resp.StatusCode)
	assert.Equal(t, "error", string(body))
}

func TestIntError(t *testing.T) {
	ts := testServer(func(c *Context, next func()) {
		panic(1)
	})
	defer ts.Close()
	resp, err := http.Get(ts.URL)
	assert.Nil(t, err)

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	assert.Nil(t, err)
	assert.Equal(t, 500, resp.StatusCode)
	assert.Equal(t, http.StatusText(500), string(body))
}

func TestRespondError(t *testing.T) {
	ts := testServer(func(c *Context, next func()) {
		c.XML([]byte{1, 2, 3})
	})
	defer ts.Close()
	resp, err := http.Get(ts.URL)
	assert.Nil(t, err)

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	assert.Nil(t, err)
	assert.Equal(t, http.StatusText(500), string(body))
}

func TestListen(t *testing.T) {
	var err error
	app := New()

	go func() {
		err = app.Listen(":3000")
	}()
	assert.Nil(t, err)
}
