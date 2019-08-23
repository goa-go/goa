package parser_test

import (
	"bytes"
	"encoding/xml"
	"errors"
	"net/http"
	"net/http/httptest"
	"net/url"
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

type Form struct {
	Int        int   `form:"int"`
	Int8       int8  `form:"int8"`
	Int16      int16 `form:"int16"`
	Int32      int32 `form:"int32"`
	Int64      int64 `form:"int64"`
	Uint       uint
	Uint8      uint8
	Uint16     uint16
	Uint32     uint32
	Uint64     uint64
	Float32    float32
	Float64    float64
	String     string
	Bool       bool
	IntArray   [2]int
	StrSlice   []string
	FloatSlice []float32
	Ignore     string `form:"-"`
	str        string
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
		e := c.ParseJSON(&obj)
		if e.Error() != `readObjectStart: expect { or n, but found j, error found in #1 byte of ...|json|..., bigger context ...|json|...` {
			t.Error("json fail error")
		}
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
		e := c.ParseXML(&xml)
		if e.Error() != `EOF` {
			t.Error("xml fail error")
		}
	})

	router.POST("/string", func(c *goa.Context) {
		str, _ := c.ParseString()
		if str != "string" {
			t.Error("parse string error")
		}
	})
	router.POST("/form", func(c *goa.Context) {
		f := Form{}
		c.ParseForm(&f)
		if f.Int != 1 || f.Int8 != 1 || f.Int16 != 1 || f.Int32 != 1 || f.Int64 != 1 {
			t.Error("parse form-int error")
		}
		if f.Uint != 1 || f.Uint8 != 1 || f.Uint16 != 1 || f.Uint32 != 1 || f.Uint64 != 1 {
			t.Error("parse form-uint error")
		}
		if f.Float32 != 1.0 || f.Float64 != 1.0 {
			t.Error("parse form-float error")
		}
		if f.String != "string" {
			t.Error("parse form-string error")
		}
		if !f.Bool {
			t.Error("parse form-bool error")
		}
		if f.IntArray != [2]int{1, 2} {
			t.Error("parse form-array error")
		}

		if len(f.StrSlice) == 0 {
			t.Error("parse form-slice error")
		} else {
			for _, v := range f.StrSlice {
				if v != "a" {
					t.Error("parse form-slice error")
				}
			}
		}

		if f.Ignore != "" {
			t.Error("parse form-ignore error")
		}
	})

	router.POST("/formfail", func(c *goa.Context) {
		f := Form{}
		err := c.ParseForm(f)
		if err == nil {
			t.Error("parse form, require a ptr error")
		}
		err = c.ParseForm(&f)
		if err == nil {
			t.Error("parse form-array error")
		}
	})

	router.POST("/form-uint-fail", func(c *goa.Context) {
		f := Form{}
		err := c.ParseForm(&f)
		if err == nil {
			t.Error("parse form-uint-fail error")
		}
	})

	router.POST("/form-int-array-fail", func(c *goa.Context) {
		f := Form{}
		err := c.ParseForm(&f)
		if err == nil {
			t.Error("parse form-int-array-fail error", f)
		}
	})

	router.POST("/form-float-slice-fail", func(c *goa.Context) {
		f := Form{}
		err := c.ParseForm(&f)
		if err == nil {
			t.Error("parse form-float-slice-fail error")
		}
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
}

func TestString(t *testing.T) {
	server := initServer(t)
	resp, err := http.Post(server.URL+"/string", "text/plain;", strings.NewReader("string"))

	if err != nil {
		t.Error("request error")
	}

	defer resp.Body.Close()
}

func TestForm(t *testing.T) {
	server := initServer(t)
	Data := map[string]string{
		"int":      "1",
		"int8":     "1",
		"int16":    "1",
		"int32":    "1",
		"int64":    "1",
		"Uint":     "1",
		"Uint8":    "1",
		"Uint16":   "1",
		"Uint32":   "1",
		"Uint64":   "1",
		"Float32":  "1.0",
		"Float64":  "1.0",
		"String":   "string",
		"Bool":     "true",
		"IntArray": "1",
		"StrSlice": "a",
		"Ingore":   "ignore",
		"str":      "str",
	}
	DataUrlVal := url.Values{}
	for key, val := range Data {
		DataUrlVal.Add(key, val)
	}
	DataUrlVal.Add("IntArray", "2")
	DataUrlVal.Add("StrSlice", "a")
	resp, err := http.Post(server.URL+"/form", "application/x-www-form-urlencoded;", strings.NewReader(DataUrlVal.Encode()))

	if err != nil {
		t.Error("request error")
	}

	defer resp.Body.Close()
}

func TestFormFail(t *testing.T) {
	server := initServer(t)
	DataUrlVal := url.Values{}

	DataUrlVal.Add("IntArray", "1")
	DataUrlVal.Add("IntArray", "1")
	DataUrlVal.Add("IntArray", "1")
	resp, err := http.Post(server.URL+"/formfail", "application/x-www-form-urlencoded;", strings.NewReader(DataUrlVal.Encode()))

	if err != nil {
		t.Error("request error")
	}

	defer resp.Body.Close()
}

func TestFormFail2(t *testing.T) {
	server := initServer(t)
	DataUrlVal := url.Values{}

	DataUrlVal.Add("Uint", "a")
	resp, err := http.Post(server.URL+"/form-uint-fail", "application/x-www-form-urlencoded;", strings.NewReader(DataUrlVal.Encode()))

	if err != nil {
		t.Error("request error")
	}

	defer resp.Body.Close()
}

func TestFormFail3(t *testing.T) {
	server := initServer(t)
	DataUrlVal := url.Values{}

	DataUrlVal.Add("IntArray", "a")
	DataUrlVal.Add("IntArray", "a")
	resp, err := http.Post(server.URL+"/form-int-array-fail", "application/x-www-form-urlencoded;", strings.NewReader(DataUrlVal.Encode()))

	if err != nil {
		t.Error("request error")
	}

	defer resp.Body.Close()
}

func TestFormFail4(t *testing.T) {
	server := initServer(t)
	DataUrlVal := url.Values{}

	DataUrlVal.Add("FloatSlice", "a")
	DataUrlVal.Add("FloatSlice", "a")
	resp, err := http.Post(server.URL+"/form-float-slice-fail", "application/x-www-form-urlencoded;", strings.NewReader(DataUrlVal.Encode()))

	if err != nil {
		t.Error("request error")
	}

	defer resp.Body.Close()
}

type Reader struct{}

func (r *Reader) Read(p []byte) (int, error) {
	return 0, errors.New("error")
}

func TestFormFail5(t *testing.T) {
	req, _ := http.NewRequest("POST", "/", &Reader{})
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded;")
	err := parser.Form{Pointer: &Reader{}}.Parse(req)
	if err == nil {
		t.Error("fail to let req.ParseForm err")
	}
}

func TestFormFail6(t *testing.T) {
	req, _ := http.NewRequest("POST", "/", strings.NewReader(""))
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded;")
	req.MultipartReader()
	err := parser.Form{Pointer: &Reader{}}.Parse(req)
	if err == nil {
		t.Error("fail to let req.ParseMultipartForm err")
	}
}
