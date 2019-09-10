package parser

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"errors"
	"net/http"
	"net/url"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

// test case
type address struct {
	City, Country string
}
type person struct {
	ID        int     `xml:"id,attr" query:"id"`
	FirstName string  `xml:"name>first" query:"firstName"`
	LastName  string  `xml:"name>last" query:"lastName"`
	Age       int     `xml:"age" query:"age"`
	Height    float32 `xml:"height,omitempty" json:"height,omitempty"`
	Married   bool    `xml:"-" json:"-"`
	Address   address
	Comment   string `xml:",comment"`
}

func getRequest(body []byte) *http.Request {
	req, _ := http.NewRequest("POST", "/", bytes.NewReader(body))
	return req
}

func TestParseJSON(t *testing.T) {
	p := person{ID: 26, FirstName: "Nicholas", LastName: "Cao", Age: 18}
	p.Comment = " Nice man. "
	p.Address = address{"Have a guess", "CN"}

	b, _ := json.Marshal(p)
	req := getRequest(b)
	ptr := &person{}

	err := JSON{Pointer: ptr}.Parse(req)
	assert.Nil(t, err)
	assert.Equal(t, p, *ptr)
}

func TestParseJSONFailed(t *testing.T) {
	req := getRequest([]byte(" "))
	ptr := &person{}

	assert.Error(t, JSON{Pointer: ptr}.Parse(req))
}

func TestParseXML(t *testing.T) {
	p := person{ID: 26, FirstName: "Nicholas", LastName: "Cao", Age: 18}
	p.Comment = " Nice man. "
	p.Address = address{"Have a guess", "CN"}

	b, _ := xml.Marshal(p)
	req := getRequest(b)
	ptr := &person{}

	err := XML{Pointer: ptr}.Parse(req)
	assert.Nil(t, err)
	assert.Equal(t, p, *ptr)
}

func TestParseXMLFailed(t *testing.T) {
	req := getRequest([]byte(" "))
	ptr := &person{}

	assert.Error(t, JSON{Pointer: ptr}.Parse(req))
}

func TestParseString(t *testing.T) {
	req := getRequest([]byte("string"))
	str, err := String{}.Parse(req)

	assert.Nil(t, err)
	assert.Equal(t, "string", str)
}

// func TestParseStringFailed(t *testing.T) {
// 	req := getRequest([]byte("string"))
// 	req.Body.Close()
// 	str, err := String{}.Parse(req)

// 	if err != nil || str != "string" {
// 		t.Errorf("parse string failed: %s", str)
// 	}
// }

type form struct {
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

func TestParseForm(t *testing.T) {
	data := map[string]string{
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
	dataURLVal := url.Values{}
	for key, val := range data {
		dataURLVal.Add(key, val)
	}
	dataURLVal.Add("IntArray", "2")
	dataURLVal.Add("StrSlice", "b")
	req, _ := http.NewRequest("POST", "/", strings.NewReader(dataURLVal.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded;")

	f := &form{}
	err := Form{Pointer: f}.Parse(req)

	assert := assert.New(t)
	assert.Nil(err)

	assert.Equal(1, f.Int)
	assert.Equal(int8(1), f.Int8)
	assert.Equal(int16(1), f.Int16)
	assert.Equal(int32(1), f.Int32)
	assert.Equal(int64(1), f.Int64)

	assert.Equal(uint(1), f.Uint)
	assert.Equal(uint8(1), f.Uint8)
	assert.Equal(uint16(1), f.Uint16)
	assert.Equal(uint32(1), f.Uint32)
	assert.Equal(uint64(1), f.Uint64)

	assert.Equal(float32(1.0), f.Float32)
	assert.Equal(float64(1.0), f.Float64)

	assert.Equal("string", f.String)

	assert.True(f.Bool)

	assert.Equal([2]int{1, 2}, f.IntArray)

	assert.Equal([]string{"a", "b"}, f.StrSlice)

	assert.Equal("", f.Ignore)
}

type Reader struct{}

func (r *Reader) Read(p []byte) (int, error) {
	return 0, errors.New("error")
}

func TestReqParseFormFailed(t *testing.T) {
	req, _ := http.NewRequest("POST", "/", &Reader{})
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded;")

	assert.Error(t, Form{Pointer: &Reader{}}.Parse(req))
}

func TestReqParseMultipartFormFailed(t *testing.T) {
	req, _ := http.NewRequest("POST", "/", strings.NewReader(""))
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded;")
	req.MultipartReader()

	assert.Error(t, Form{Pointer: &Reader{}}.Parse(req))
}

func TestParseQuerry(t *testing.T) {
	p := &person{}
	req, _ := http.NewRequest("GET", "/?id=1&firstName=Nicholas&lastName=Cao&age=18", nil)
	err := Query{Pointer: p}.Parse(req)
	assert := assert.New(t)

	assert.Nil(err)
	assert.Equal(1, p.ID, 1)
	assert.Equal("Nicholas", p.FirstName)
	assert.Equal("Cao", p.LastName)
	assert.Equal(18, p.Age)
}

/* mapByTagFailed */

func TestMapByTagNoPtr(t *testing.T) {
	assert.Error(t, mapByTag(form{}, nil, ""))
}

func TestMapTypeError(t *testing.T) {
	dataURLVal := url.Values{}
	dataURLVal.Add("Uint", "a")

	req, _ := http.NewRequest("POST", "/", strings.NewReader(dataURLVal.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded;")

	req.ParseForm()
	req.ParseMultipartForm(defaultMaxMemory)

	f := &form{}
	assert.Error(t, mapByTag(f, req.Form, "form"))
}

func TestMapArrayFailed(t *testing.T) {
	dataURLVal := url.Values{}

	dataURLVal.Add("IntArray", "1")
	dataURLVal.Add("IntArray", "1")
	dataURLVal.Add("IntArray", "1")
	req, _ := http.NewRequest("POST", "/", strings.NewReader(dataURLVal.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded;")

	req.ParseForm()
	req.ParseMultipartForm(defaultMaxMemory)

	f := &form{}
	assert.Error(t, mapByTag(f, req.Form, "form"))
}

func TestMapIntArrayFailed(t *testing.T) {
	dataURLVal := url.Values{}
	dataURLVal.Add("IntArray", "a")
	dataURLVal.Add("IntArray", "a")

	req, _ := http.NewRequest("POST", "/", strings.NewReader(dataURLVal.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded;")

	req.ParseForm()
	req.ParseMultipartForm(defaultMaxMemory)

	f := &form{}
	assert.Error(t, mapByTag(f, req.Form, "form"))
}

func TestMapSliceFailed(t *testing.T) {
	dataURLVal := url.Values{}

	dataURLVal.Add("FloatSlice", "a")
	req, _ := http.NewRequest("POST", "/", strings.NewReader(dataURLVal.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded;")

	req.ParseForm()
	req.ParseMultipartForm(defaultMaxMemory)

	f := &form{}
	assert.Error(t, mapByTag(f, req.Form, "form"))
}
