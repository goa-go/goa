package parser

import (
	"io/ioutil"
	"net/http"

	"github.com/goa-go/goa/util"
)

// String is a json-parser instance.
type String struct{}

// Parse string-data.
func (p String) Parse(req *http.Request) (string, error) {
	b, err := ioutil.ReadAll(req.Body)
	return util.Bytes2str(b), err
}
