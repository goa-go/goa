package parser

import (
	"net/http"
)

// Parser is a handler for http-request.
type Parser interface {
	// Parse request data.
	Parse(*http.Request) error
}
