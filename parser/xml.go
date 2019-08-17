package parser

import (
	"encoding/xml"
	"net/http"
)

// XML is a json-parser instance.
type XML struct {
	Pointer interface{}
}

// Parse xml-data.
func (p XML) Parse(req *http.Request) error {
	return xml.NewDecoder(req.Body).Decode(p.Pointer)
}
