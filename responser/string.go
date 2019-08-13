package responser

import (
	"fmt"
	"net/http"
)

type String struct {
	Data string
}

func (r String) Respond(w http.ResponseWriter) error {
	_, err := fmt.Fprint(w, r.Data)
	return err
}
