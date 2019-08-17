package responser

import (
	"fmt"
	"net/http"
)

// String is a string-responser instance.
type String struct {
	Data string
}

// Respond string-data.(text/html)
func (r String) Respond(w http.ResponseWriter) error {
	_, err := fmt.Fprint(w, r.Data)
	return err
}
