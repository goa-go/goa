package responser

import (
	"net/http"

	json "github.com/json-iterator/go"
)

type JSON struct {
	Data interface{}
}

func (r JSON) Respond(w http.ResponseWriter) error {
	return json.NewEncoder(w).Encode(&r.Data)
}
