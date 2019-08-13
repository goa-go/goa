package parser

import (
	"io/ioutil"
	"net/http"
)

type String struct{}

func (p String) Parse(req *http.Request) (string, error) {
	b, err := ioutil.ReadAll(req.Body)
	return string(b), err
}
