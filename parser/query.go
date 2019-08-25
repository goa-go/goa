package parser

import (
	"net/http"
)

// Query is a json-parser instance.
type Query struct {
	Pointer interface{}
}

// Parse query-data.
func (p Query) Parse(req *http.Request) error {
	return mapByTag(p.Pointer, req.URL.Query(), "query")
}
