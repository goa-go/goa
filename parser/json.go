package parser

import (
	"net/http"

	json "github.com/json-iterator/go"
)

type JSON struct {
	Pointer interface{}
}

func (p JSON) Parse(req *http.Request) error {
	return json.NewDecoder(req.Body).Decode(p.Pointer)
}
