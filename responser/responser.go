package responser

import (
	"net/http"
)

type Responser interface {
	// Respond writes data with custom ContentType.
	Respond(http.ResponseWriter) error
}
