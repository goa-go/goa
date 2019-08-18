// test c.Body = ...
package body_test

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/goa-go/goa"
	"github.com/goa-go/goa/router"
	json "github.com/json-iterator/go"
)

// json test case
type Address struct {
	City, Country string
}
type Person struct {
	Id        int
	FirstName string
	LastName  string
	Age       int
	Height    float32 `json:"height,omitempty"`
	Married   bool
	Address
	Comment string
}

func jsonHandler(c *goa.Context) {
	obj := Person{Id: 26, FirstName: "Nicholas", LastName: "Cao", Age: 18}
	obj.Comment = " Nice man. "
	obj.Address = Address{"Have a guess", "CN"}

	c.Body = obj
}

func initServer() *httptest.Server {
	app := goa.New()
	router := router.New()

	router.GET("/string", func(c *goa.Context) {
		c.Body = "string"
	})

	router.GET("/html", func(c *goa.Context) {
		c.Body = "<p>html</p>"
	})
	router.GET("/json", jsonHandler)

	app.Use(router.Routes())

	// Before testing, must compose middlewares.
	app.ComposeMiddlewares()
	return httptest.NewServer(app)
}

func TestString(t *testing.T) {
	server := initServer()
	defer server.Close()

	resp, err := http.Get(server.URL + "/string")

	if err != nil {
		t.Error("request error")
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	if string(body) != "string" {
		t.Error("string-body error")
	}
}

func TestHTML(t *testing.T) {
	server := initServer()
	defer server.Close()

	resp, err := http.Get(server.URL + "/html")

	if err != nil {
		t.Error("request error")
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	if string(body) != "<p>html</p>" || resp.Header["Content-Type"][0] != "text/html; charset=utf-8" {
		t.Error("html-body error")
	}
}

func TestJson(t *testing.T) {
	server := initServer()
	defer server.Close()

	resp, err := http.Get(server.URL + "/json")
	if err != nil {
		t.Error("request /json error")
	}
	defer resp.Body.Close()

	obj := Person{Id: 26, FirstName: "Nicholas", LastName: "Cao", Age: 18}
	obj.Comment = " Nice man. "
	obj.Address = Address{"Have a guess", "CN"}

	obj2 := Person{}
	json.NewDecoder(resp.Body).Decode(&obj2)

	b1, _ := json.Marshal(obj)
	b2, _ := json.Marshal(obj2)
	if string(b1) != string(b2) {
		t.Error("json-body error")
	}
}
