package parser

import (
	"io/ioutil"
	"net/http"
)

// String is a json-parser instance.
type String struct{}

// Parse string-data.
func (p String) Parse(req *http.Request) string {
	b, err := ioutil.ReadAll(req.Body)
	if err != nil {
		panic(err)
	}
	return string(b)
}
