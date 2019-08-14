package goa_test

import (
	"encoding/xml"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"

	"github.com/goa-go/goa"
	"github.com/goa-go/goa/router"
	json "github.com/json-iterator/go"
)

// xml test case
type Address struct {
	City, Country string
}
type Person struct {
	Id        int     `xml:"id,attr"`
	FirstName string  `xml:"name>first"`
	LastName  string  `xml:"name>last"`
	Age       int     `xml:"age"`
	Height    float32 `xml:"height,omitempty"`
	Married   bool
	Address
	Comment string `xml:",comment"`
}

func xmlHandler(c *goa.Context) {
	xml := &Person{Id: 26, FirstName: "Nicholas", LastName: "Cao", Age: 18}
	xml.Comment = " Nice man. "
	xml.Address = Address{"Have a guess", "CN"}
	c.XML(xml)
}

func jsonHandler(c *goa.Context) {
	c.JSON(goa.M{
		"string": "string",
		"int":    1,
		"json": goa.M{
			"key": "value",
		},
	})
}

func setStatus(c *goa.Context) {
	code := c.Param("code")
	int, err := strconv.Atoi(code)
	if err != nil {
		c.Status(400).String("plz input int")
	} else {
		c.Status(int).String("ok")
	}
}

func hello(c *goa.Context) {
	name := c.Query("name")
	c.String("hello " + name)
}

func postForm(c *goa.Context) {
	value := c.PostForm("key")
	c.String("key: " + value)
}

func initServer() *httptest.Server {
	app := goa.New()
	router := router.New()

	router.GET("/", func(c *goa.Context) {
		c.String("hello world")
	})
	router.GET("/xml", xmlHandler)
	router.GET("/json", jsonHandler)
	router.GET("/redirect", func(c *goa.Context) {
		c.Redirect(301, "/")
	})
	router.GET("/status/:code", setStatus)
	router.GET("/hello", hello)
	router.POST("/postForm", postForm)

	app.Use(router.Routes())

	// Before testing, must compose middlewares.
	app.ComposeMiddlewares()
	return httptest.NewServer(app)
}

func TestRequest(t *testing.T) {
	server := initServer()
	defer server.Close()

	resp, err := http.Get(server.URL)

	if err != nil {
		t.Error("request error")
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	if string(body) != "hello world" {
		t.Error("request error")
	}
}

func TestStatusCode(t *testing.T) {
	testStatusCode(t, 200)
	testStatusCode(t, 300)
	testStatusCode(t, 400)
	testStatusCode(t, 500)
}

func testStatusCode(t *testing.T, code int) {
	server := initServer()
	defer server.Close()

	resp, err := http.Get(server.URL + "/status/" + strconv.Itoa(code))
	if err != nil {
		t.Error("request /status error")
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	if string(body) != "ok" && resp.StatusCode == code {
		t.Error("status code error")
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

	var obj map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&obj)

	if obj["string"] == "string" && obj["int"] == 1.0 {
		if obj2, ok := obj["json"].(map[string]interface{}); ok {
			if obj2["key"] != "value" {
				t.Error("json error")
			}
		} else {
			t.Error("json error")
		}
	} else {
		t.Error("json error")
	}
}

func TestXml(t *testing.T) {
	server := initServer()
	defer server.Close()

	resp, err := http.Get(server.URL + "/xml")
	if err != nil {
		t.Error("request /xml error")
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	XML := &Person{Id: 26, FirstName: "Nicholas", LastName: "Cao", Age: 18}
	XML.Comment = " Nice man. "
	XML.Address = Address{"Have a guess", "CN"}

	b, _ := xml.Marshal(XML)

	if string(body) != string(b) {
		t.Error("xml error", string(b))
	}
}

func TestQuery(t *testing.T) {
	server := initServer()
	defer server.Close()

	resp, err := http.Get(server.URL + "/hello?name=nicholascao")
	if err != nil {
		t.Error("request error")
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	if string(body) != "hello nicholascao" {
		t.Error("request error")
	}
}

func TestPostForm(t *testing.T) {
	server := initServer()
	defer server.Close()

	resp, err := http.Post(server.URL+"/postForm", "application/x-www-form-urlencoded;", strings.NewReader("key=value"))
	if err != nil {
		t.Error("request error")
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)

	if string(body) != "key: value" {
		t.Error(string(body))
	}
}

func TestRedirect(t *testing.T) {
	server := initServer()
	defer server.Close()

	resp, err := http.Get(server.URL + "/redirect")
	if err != nil {
		t.Error("request error")
	}

	defer resp.Body.Close()
	b, _ := ioutil.ReadAll(resp.Body)
	if string(b) != "hello world" {
		t.Error("redirect error")
	}
}
