package responser

import (
	"encoding/json"
	"encoding/xml"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

// test case
type address struct {
	City, Country string
}
type person struct {
	ID        int     `xml:"id,attr"`
	FirstName string  `xml:"name>first"`
	LastName  string  `xml:"name>last"`
	Age       int     `xml:"age"`
	Height    float32 `xml:"height,omitempty" json:"height,omitempty"`
	Married   bool    `xml:"-" json:"-"`
	Address   address
	Comment   string `xml:",comment"`
}

func TestRespondJSON(t *testing.T) {
	p := person{ID: 26, FirstName: "Nicholas", LastName: "Cao", Age: 18}
	p.Comment = " Nice man. "
	p.Address = address{"Have a guess", "CN"}
	w := httptest.NewRecorder()
	err := JSON{Data: p}.Respond(w)

	bytes, _ := json.Marshal(p)

	assert.Nil(t, err)
	assert.Equal(t, string(bytes)+"\n", w.Body.String())
}

// func TestRespondJSONFailed(t *testing.T) {
// 	w := httptest.NewRecorder()
// 	err := JSON{Data: []byte{1, 2, 3}}.Respond(w)

// 	if err == nil {
// 		t.Fatal()
// 	}
// }

func TestRespondXML(t *testing.T) {
	p := person{ID: 26, FirstName: "Nicholas", LastName: "Cao", Age: 18}
	p.Comment = " Nice man. "
	p.Address = address{"Have a guess", "CN"}
	w := httptest.NewRecorder()
	err := XML{Data: p}.Respond(w)

	bytes, _ := xml.Marshal(p)

	assert.Nil(t, err)
	assert.Equal(t, string(bytes), w.Body.String())
}

func TestRespondXMLFailed(t *testing.T) {
	w := httptest.NewRecorder()
	assert.Error(t, XML{Data: []byte{1, 2, 3}}.Respond(w))
}

func TestRespondString(t *testing.T) {
	w := httptest.NewRecorder()
	err := String{Data: "string"}.Respond(w)

	assert.Nil(t, err)
	assert.Equal(t, "string", w.Body.String())
}

// func TestRespondStringFailed(t *testing.T) {
// 	w := httptest.NewRecorder()
// 	err := String{Data: "string"}.Respond(w)

// 	body, _ := ioutil.ReadAll(w.Result().Body)

// 	if err != nil || string(body) != "string" {
// 		t.Errorf("respond string failed: %v", string(body))
// 	}
// }
