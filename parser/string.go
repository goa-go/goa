package parser

import (
	"io/ioutil"
	"net/http"
)

// String is a json-parser instance.
type String struct{}

// Parse string-data.
func (p String) Parse(req *http.Request) (string, error) {
	b, err := ioutil.ReadAll(req.Body)
	return string(b), err
}
