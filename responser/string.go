package responser

import (
	"net/http"

	"github.com/goa-go/goa/util"
)

// String is a string-responser instance.
type String struct {
	Data string
}

// Respond string-data.(text/html)
func (r String) Respond(w http.ResponseWriter) error {
	_, err := w.Write(util.Str2bytes(r.Data))
	return err
}
