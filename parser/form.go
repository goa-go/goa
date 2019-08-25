package parser

import (
	"net/http"
)

// Form is a form-parser instance.
type Form struct {
	Pointer interface{}
}

const (
	defaultMaxMemory = 32 << 20 // 32 MB
)

// Parse form-data.
func (p Form) Parse(req *http.Request) error {
	if err := req.ParseForm(); err != nil {
		return err
	}
	if err := req.ParseMultipartForm(defaultMaxMemory); err != nil {
		if err != http.ErrNotMultipart {
			return err
		}
	}

	return mapByTag(p.Pointer, req.Form, "form")
}
