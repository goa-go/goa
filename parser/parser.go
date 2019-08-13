package parser

import (
	"net/http"
)

type Parser interface {
	// Parse request data.
	Parse(*http.Request) error
}
