package responser

import (
	"encoding/json"
	"net/http"
)

// JSON is a json-responser instance.
type JSON struct {
	Data interface{}
}

// Respond json-data.
func (r JSON) Respond(w http.ResponseWriter) error {
	return json.NewEncoder(w).Encode(&r.Data)
}
