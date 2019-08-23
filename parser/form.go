package parser

import (
	"fmt"
	"net/http"
	"reflect"
	"strconv"
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
	value := reflect.ValueOf(p.Pointer)
	if value.Kind() != reflect.Ptr {
		return fmt.Errorf("Expected a pointer, but got a %s", value.Kind())
	}

	if err := req.ParseForm(); err != nil {
		return err
	}
	if err := req.ParseMultipartForm(defaultMaxMemory); err != nil {
		if err != http.ErrNotMultipart {
			return err
		}
	}

	value = value.Elem()
	t := value.Type()
	for i := 0; i < value.NumField(); i++ {
		field := value.Field(i)
		tfield := t.Field(i)
		if tfield.PkgPath != "" && !tfield.Anonymous { // unexported
			continue
		}
		key, ignore := getKey(tfield)
		if ignore {
			// tag is "-"
			continue
		}

		if strArray, ok := req.Form[key]; ok {
			switch field.Kind() {
			case reflect.Slice:
				if err := setSlice(field, strArray); err != nil {
					return err
				}
			case reflect.Array:
				if len(strArray) != field.Len() {
					return fmt.Errorf("%q is not valid value for %s", strArray, field.Type().String())
				}
				if err := setArray(field, strArray); err != nil {
					return err
				}
			default:
				if err := setValue(field, strArray[0]); err != nil {
					return err
				}
			}
		} else {
			continue
		}
	}
	return nil
}

// getKey returns key and wether to ignore.
func getKey(field reflect.StructField) (string, bool) {
	name := field.Name
	tag := field.Tag.Get("form")

	if tag == "-" {
		return "", true
	}
	if tag == "" {
		return name, false
	}
	return tag, false
}

func setValue(field reflect.Value, val string) error {
	switch field.Kind() {
	case reflect.Int:
		return setIntField(field, val, 0)
	case reflect.Int8:
		return setIntField(field, val, 8)
	case reflect.Int16:
		return setIntField(field, val, 16)
	case reflect.Int32:
		return setIntField(field, val, 32)
	case reflect.Int64:
		return setIntField(field, val, 64)
	case reflect.Uint:
		return setUintField(field, val, 0)
	case reflect.Uint8:
		return setUintField(field, val, 8)
	case reflect.Uint16:
		return setUintField(field, val, 16)
	case reflect.Uint32:
		return setUintField(field, val, 32)
	case reflect.Uint64:
		return setUintField(field, val, 64)
	case reflect.Float32:
		return setFloatField(field, val, 32)
	case reflect.Float64:
		return setFloatField(field, val, 64)
	case reflect.Bool:
		return setBoolField(field, val)
	case reflect.String:
		field.SetString(val)
	}
	return nil
}

func setIntField(field reflect.Value, val string, bitSize int) error {
	intVal, err := strconv.ParseInt(val, 10, bitSize)
	if err == nil {
		field.SetInt(intVal)
	}
	return err
}

func setUintField(field reflect.Value, val string, bitSize int) error {
	uintVal, err := strconv.ParseUint(val, 10, bitSize)
	if err == nil {
		field.SetUint(uintVal)
	}
	return err
}

func setFloatField(field reflect.Value, val string, bitSize int) error {
	floatVal, err := strconv.ParseFloat(val, bitSize)
	if err == nil {
		field.SetFloat(floatVal)
	}
	return err
}

func setBoolField(field reflect.Value, val string) error {
	boolVal, err := strconv.ParseBool(val)
	if err == nil {
		field.SetBool(boolVal)
	}
	return err
}

func setArray(field reflect.Value, strArray []string) error {
	for i, str := range strArray {
		err := setValue(field.Index(i), str)
		if err != nil {
			return err
		}
	}
	return nil
}

func setSlice(field reflect.Value, strArray []string) error {
	slice := reflect.MakeSlice(field.Type(), len(strArray), len(strArray))
	err := setArray(slice, strArray)
	if err != nil {
		return err
	}
	field.Set(slice)
	return nil
}
