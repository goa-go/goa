package encode

import (
	"encoding/json"
	"encoding/xml"
)

func JSON(v interface{}) []byte {
	b, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		panic(err)
	}
	return b
}

func String(v string) []byte {
	return []byte(v)
}

func XML(v interface{}) []byte {
	b, err := xml.MarshalIndent(v, "", "  ")
	if err != nil {
		panic(err)
	}
	return b
}
