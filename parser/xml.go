package parser

import (
	"encoding/xml"
	"net/http"
)

type XML struct {
	Pointer interface{}
}

func (p XML) Parse(req *http.Request) error {
	return xml.NewDecoder(req.Body).Decode(p.Pointer)
}
