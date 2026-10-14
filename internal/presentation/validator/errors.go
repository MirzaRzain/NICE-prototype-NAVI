package validator

import (
	"bytes"
	"fmt"
	"strings"
	"unicode"

	playground "github.com/go-playground/validator/v10"
)

func ValidationErrors(err error) (message map[string]interface{}) {
	ve, ok := err.(playground.ValidationErrors)
	if !ok {
		return
	}

	message = make(map[string]interface{})

	for _, e := range ve {
		jsonKey := strings.ToLower(e.Field())
		fieldName := convertCase(e.Field(), ' ')

		switch e.Tag() {
		case "datetime":
			message[jsonKey] = fmt.Sprintf("Field %s harus berupa tanggal & waktu", fieldName)
		case "email":
			message[jsonKey] = "Input harus berupa alamat email yang valid"
		case "max":
			message[jsonKey] = fmt.Sprintf("Field %s maksimal %s", fieldName, e.Param())
		case "min":
			message[jsonKey] = fmt.Sprintf("Field %s minimal %s", fieldName, e.Param())
		case "required":
			message[jsonKey] = fmt.Sprintf("Field %s tidak boleh kosong!", fieldName)
		case "oneof":
			message[jsonKey] = fmt.Sprintf("Field %s harus salah satu dari [%s]", fieldName, e.Param())
		case "eqfield":
			message[jsonKey] = fmt.Sprintf("Field %s harus sama dengan field %s", fieldName, e.Param())
		case "gt":
			message[jsonKey] = fmt.Sprintf("Field %s harus lebih besar dari %s", fieldName, e.Param())
		case "lt":
			message[jsonKey] = fmt.Sprintf("Field %s harus lebih kecil dari %s", fieldName, e.Param())
		case "gte":
			message[jsonKey] = fmt.Sprintf("Field %s harus lebih besar atau sama dengan %s", fieldName, e.Param())
		case "lte":
			message[jsonKey] = fmt.Sprintf("Field %s harus lebih kecil atau sama dengan %s", fieldName, e.Param())
		case "numeric":
			message[jsonKey] = fmt.Sprintf("Field %s harus berupa angka", fieldName)
		case "alphanum":
			message[jsonKey] = fmt.Sprintf("Field %s harus berupa alfanumerik", fieldName)
		case "url":
			message[jsonKey] = fmt.Sprintf("Field %s harus berupa URL yang valid", fieldName)
		case "uri":
			message[jsonKey] = fmt.Sprintf("Field %s harus berupa path atau URI yang valid", fieldName)
		default:
			message[jsonKey] = fmt.Sprintf("Field %s tidak valid (%s)", fieldName, e.Tag())
		}
	}

	return
}

func convertCase(t string, c rune) string {
	buf := &bytes.Buffer{}

	for i, r := range t {
		if i > 0 && unicode.IsUpper(r) {
			if t[i-1] != 'I' && r != 'D' {
				buf.WriteRune(c)
			}
		}

		buf.WriteRune(r)
	}

	return buf.String()
}
