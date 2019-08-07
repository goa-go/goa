package main_test

import (
	"encoding/json"
	"encoding/xml"
	"io/ioutil"
	"net/http"
	"strconv"
	"testing"
)

func TestRequest(t *testing.T) {
	resp, err := http.Get("http://localhost:3000")
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
	resp, err := http.Get("http://localhost:3000/status/" + strconv.Itoa(code))
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
	resp, err := http.Get("http://localhost:3000/json")
	if err != nil {
		t.Error("request /json error")
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	b, _ := json.MarshalIndent(map[string]interface{}{
		"string": "string",
		"int":    1,
		"json": map[string]interface{}{
			"key": "value",
		},
	}, "", "  ")

	if string(body) != string(b) {
		t.Error("json error")
	}
}

func TestXml(t *testing.T) {
	resp, err := http.Get("http://localhost:3000/xml")
	if err != nil {
		t.Error("request /xml error")
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	type Address struct {
		City, State string
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
	XML := &Person{Id: 13, FirstName: "John", LastName: "Doe", Age: 42}
	XML.Comment = " Need more details. "
	XML.Address = Address{"Hanga Roa", "Easter Island"}

	b, _ := xml.MarshalIndent(XML, "", "  ")

	if string(body) != string(b) {
		t.Error("xml error")
	}
}