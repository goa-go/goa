package responser

import (
	"encoding/xml"
	"net/http"
)

// XML is a xml-responser instance.
type XML struct {
	Data interface{}
}

// Respond xml-data.
func (r XML) Respond(w http.ResponseWriter) error {
	return xml.NewEncoder(w).Encode(r.Data)
}
