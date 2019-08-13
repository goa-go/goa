package responser

import (
	"encoding/xml"
	"net/http"
)

type XML struct {
	Data interface{}
}

func (r XML) Respond(w http.ResponseWriter) error {
	return xml.NewEncoder(w).Encode(r.Data)
}
