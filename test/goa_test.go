package goa_test

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"

	"github.com/goa-go/goa"
	"github.com/goa-go/goa/router"
)

// xml and json test case
type Address struct {
	City, Country string
}
type Person struct {
	Id        int     `xml:"id,attr"`
	FirstName string  `xml:"name>first"`
	LastName  string  `xml:"name>last"`
	Age       int     `xml:"age"`
	Height    float32 `xml:"height,omitempty" json:"height,omitempty"`
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
	obj := Person{Id: 26, FirstName: "Nicholas", LastName: "Cao", Age: 18}
	obj.Comment = " Nice man. "
	obj.Address = Address{"Have a guess", "CN"}
	c.JSON(obj)
}

func setStatus(c *goa.Context) {
	code := c.Param("code")
	code2 := c.Param("code2")
	if code2 != "" {
		panic("params error")
	}
	int, err := strconv.Atoi(code)
	if err != nil {
		c.Status = 400
		c.String("plz input int")
	} else {
		c.Status = int
		c.String("ok")
	}
}

func hello(c *goa.Context) {
	name := c.Query("name")
	name2 := c.Query("name2")
	c.String("hello " + name + "," + name2)
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
	router.GET("/html", func(c *goa.Context) {
		c.HTML("<p>html</p>")
	})
	router.GET("/xml", xmlHandler)
	router.GET("/json", jsonHandler)
	router.GET("/xmlerror", func(c *goa.Context) {
		c.XML([]byte{1, 2, 3})
	})
	router.GET("/redirect", func(c *goa.Context) {
		c.Redirect(302, "/")
	})
	router.GET("/keys", func(c *goa.Context) {
		c.Set("key", "value")
		v, _ := c.Get("key")
		c.String(v.(string))
	})
	router.GET("/keys2", func(c *goa.Context) {
		v2, _ := c.Get("key2")
		if v2 == nil {
			c.Error(goa.Error{
				Msg:    "key does not exist",
				Status: 500,
			})
		}
	})
	router.GET("/status/:code", setStatus)
	router.GET("/hello", hello)
	router.POST("/postForm", postForm)
	router.GET("/error", func(c *goa.Context) {
		c.Error(goa.Error{Msg: "msg"})
	})
	router.GET("/stringerror", func(c *goa.Context) {
		panic("error")
	})
	router.GET("/emptyerror", func(c *goa.Context) {
		panic("")
	})
	router.GET("/header", func(c *goa.Context) {
		c.SetHeader("Goa-Header", "test")
	})

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
		t.Error("html error")
	}
}

func TestStatusCode(t *testing.T) {
	testStatusCode(t, 200)
	testStatusCode(t, 300)
	testStatusCode(t, 400)
	testStatusCode(t, 500)
	testStatusCode(t, 99)
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

	if code < 100 || code > 999 {
		if string(body) != fmt.Sprintf("invalid status code: %d", code) {
			t.Error("error status code error")
		}
	} else if string(body) != "ok" || resp.StatusCode != code {
		t.Error("status code error: " + strconv.Itoa(code))
	}
}

func TestJSON(t *testing.T) {
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
		t.Error("json error")
	}
}

func TestXMLerror(t *testing.T) {
	server := initServer()
	defer server.Close()

	resp, err := http.Get(server.URL + "/xmlerror")

	if err != nil {
		t.Error("request error")
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	if string(body) != "xml: unsupported type: []uint8" {
		t.Error("xmlError error")
	}
}

func TestXML(t *testing.T) {
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

func TestKeys(t *testing.T) {
	server := initServer()
	defer server.Close()

	resp, err := http.Get(server.URL + "/keys")

	if err != nil {
		t.Error("keys error")
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	if string(body) != "value" {
		t.Error("keys error")
	}
}

func TestKeysFail(t *testing.T) {
	server := initServer()
	defer server.Close()

	resp, err := http.Get(server.URL + "/keys2")

	if err != nil {
		t.Error("keys error")
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	if string(body) != "key does not exist" && resp.StatusCode == 500 {
		t.Error("keys error")
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

	if string(body) != "hello nicholascao," {
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

func TestError(t *testing.T) {
	server := initServer()
	defer server.Close()

	resp, err := http.Get(server.URL + "/error")

	if err != nil {
		t.Error("request error")
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	if string(body) != "msg" && resp.StatusCode != 500 {
		t.Error("onerror error")
	}
}

func TestStringError(t *testing.T) {
	server := initServer()
	defer server.Close()

	resp, err := http.Get(server.URL + "/stringerror")

	if err != nil {
		t.Error("request error")
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	if string(body) != "error" && resp.StatusCode != 500 {
		t.Error("stringError error")
	}
}

func TestEmptyError(t *testing.T) {
	server := initServer()
	defer server.Close()

	resp, err := http.Get(server.URL + "/emptyerror")

	if err != nil {
		t.Error("request error")
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	if string(body) != "Internal Server Error" && resp.StatusCode != 500 {
		t.Error("emptyError error")
	}
}

func TestHeader(t *testing.T) {
	server := initServer()
	defer server.Close()

	resp, err := http.Get(server.URL + "/header")

	if err != nil {
		t.Error("request error")
	}
	defer resp.Body.Close()
	if resp.Header["Goa-Header"][0] != "test" {
		t.Error("header error")
	}
}

func TestListen(t *testing.T) {
	app := goa.New()
	var err error
	go func() {
		err = app.Listen(":3000")
	}()

	if err != nil {
		t.Error("listen error")
	}
}
