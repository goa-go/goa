package responser

import (
	"net/http"
)

// Responser is a handler for http-response.
type Responser interface {
	// Respond writes data with custom ContentType.
	Respond(http.ResponseWriter) error
}
