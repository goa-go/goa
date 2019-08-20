package parser_test

import (
	"bytes"
	"encoding/xml"
	"errors"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/goa-go/goa"
	"github.com/goa-go/goa/parser"
	"github.com/goa-go/goa/router"
	json "github.com/json-iterator/go"
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

func initServer(t *testing.T) *httptest.Server {
	app := goa.New()
	router := router.New()

	router.POST("/json", func(c *goa.Context) {
		obj := Person{}
		c.ParseJSON(&obj)
		obj2 := Person{Id: 26, FirstName: "Nicholas", LastName: "Cao", Age: 18}
		obj2.Comment = " Nice man. "
		obj2.Address = Address{"Have a guess", "CN"}

		b1, _ := json.Marshal(obj)
		b2, _ := json.Marshal(obj2)
		if string(b1) != string(b2) {
			t.Error("parse json error")
		}
	})
	router.POST("/jsonfail", func(c *goa.Context) {
		obj := Person{}
		c.ParseJSON(&obj)
	})
	router.POST("/xml", func(c *goa.Context) {
		XML := Person{}
		c.ParseXML(&XML)
		XML2 := Person{Id: 26, FirstName: "Nicholas", LastName: "Cao", Age: 18}
		XML2.Comment = " Nice man. "
		XML2.Address = Address{"Have a guess", "CN"}

		b1, _ := xml.Marshal(XML)
		b2, _ := xml.Marshal(XML2)
		if string(b1) != string(b2) {
			t.Error("parse xml error")
		}
	})
	router.POST("/xmlfail", func(c *goa.Context) {
		xml := Person{}
		c.ParseXML(&xml)
	})
	router.POST("/string", func(c *goa.Context) {
		str := c.ParseString()
		if str != "string" {
			t.Error("parse string error")
		}
	})
	router.POST("/stringfail", func(c *goa.Context) {
		c.ParseString()
	})

	app.Use(router.Routes())

	// Before testing, must compose middlewares.
	app.ComposeMiddlewares()
	return httptest.NewServer(app)
}

func TestJSON(t *testing.T) {
	server := initServer(t)

	obj := Person{Id: 26, FirstName: "Nicholas", LastName: "Cao", Age: 18}
	obj.Comment = " Nice man. "
	obj.Address = Address{"Have a guess", "CN"}

	b, _ := json.Marshal(obj)

	resp, err := http.Post(server.URL+"/json", "application/json;", bytes.NewReader(b))

	if err != nil {
		t.Error("request error")
	}
	defer resp.Body.Close()
}

func TestJSONFail(t *testing.T) {
	server := initServer(t)

	resp, err := http.Post(server.URL+"/jsonfail", "application/json;", strings.NewReader("json"))

	if err != nil {
		t.Error("request error")
	}

	defer resp.Body.Close()
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Error("read error")
	}
	if string(b) != `readObjectStart: expect { or n, but found j, error found in #1 byte of ...|json|..., bigger context ...|json|...` {
		t.Error("json fail errir")
	}
}

func TestXML(t *testing.T) {
	server := initServer(t)

	XML := Person{Id: 26, FirstName: "Nicholas", LastName: "Cao", Age: 18}
	XML.Comment = " Nice man. "
	XML.Address = Address{"Have a guess", "CN"}

	b, _ := xml.Marshal(XML)

	resp, err := http.Post(server.URL+"/xml", "application/xml;", bytes.NewReader(b))

	if err != nil {
		t.Error("request error")
	}
	defer resp.Body.Close()
}

func TestXMLFail(t *testing.T) {
	server := initServer(t)
	resp, err := http.Post(server.URL+"/xmlfail", "application/xml;", strings.NewReader("xml"))

	if err != nil {
		t.Error("request error")
	}

	defer resp.Body.Close()
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Error("read error")
	}
	if string(b) != `EOF` {
		t.Error("xml fail error")
	}
}

func TestString(t *testing.T) {
	server := initServer(t)
	resp, err := http.Post(server.URL+"/string", "text/plain;", strings.NewReader("string"))

	if err != nil {
		t.Error("request error")
	}

	defer resp.Body.Close()
}

type reader struct{}

func (r reader) Read(b []byte) (int, error) {
	return 0, errors.New("error")
}

func TestStringFail(t *testing.T) {
	req, _ := http.NewRequest("POST", "/", strings.NewReader(""))
	req.Body = ioutil.NopCloser(&reader{})
	defer func() {
		e := recover()
		if e.(error).Error() != "error" {
			t.Error("string fail")
		}
	}()
	parser.String{}.Parse(req)
}
