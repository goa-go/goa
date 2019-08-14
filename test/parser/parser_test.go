package parser_test

import (
	"bytes"
	"encoding/xml"
	"net/http"
	"testing"

	"github.com/goa-go/goa/parser"
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

func TestParseJSON(t *testing.T) {
	obj := Person{Id: 26, FirstName: "Nicholas", LastName: "Cao", Age: 18}
	obj.Comment = " Nice man. "
	obj.Address = Address{"Have a guess", "CN"}

	b, _ := json.Marshal(obj)

	req, _ := http.NewRequest("POST", "/", bytes.NewBuffer(b))
	obj2 := Person{}
	jsonParser := parser.JSON{
		Pointer: &obj2,
	}
	err := jsonParser.Parse(req)
	if err != nil {
		t.Error("parse json error")
	}
	b2, _ := json.Marshal(obj2)

	if string(b) != string(b2) {
		t.Error("parse json error")
	}
}

func TestParseXML(t *testing.T) {
	XML := Person{Id: 26, FirstName: "Nicholas", LastName: "Cao", Age: 18}
	XML.Comment = " Nice man. "
	XML.Address = Address{"Have a guess", "CN"}

	b, _ := xml.Marshal(XML)

	req, _ := http.NewRequest("POST", "/", bytes.NewBuffer(b))
	XML2 := Person{}
	xmlParser := parser.XML{
		Pointer: &XML2,
	}
	err := xmlParser.Parse(req)
	if err != nil {
		t.Error("parse xml error")
	}
	b2, _ := xml.Marshal(XML2)
	if string(b) != string(b2) {
		t.Error("parse xml error")
	}
}

func TestParseString(t *testing.T) {
	req, _ := http.NewRequest("POST", "/", bytes.NewBufferString("string"))
	str := parser.String{}
	s, err := str.Parse(req)
	if err != nil {
		t.Error("parse string error")
	}
	if s != "string" {
		t.Error("parse string error")
	}
}
