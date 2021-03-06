package parser

import (
	"encoding/json"
	"net/http"
)

// JSON is a json-parser instance.
type JSON struct {
	Pointer interface{}
}

// Parse json-data.
func (p JSON) Parse(req *http.Request) error {
	return json.NewDecoder(req.Body).Decode(p.Pointer)
}
