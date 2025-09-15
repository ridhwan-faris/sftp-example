package helper

import (
	"encoding/json"
	"fmt"
	"net/http"
)

var (
	invalidFormat = "Field '%s' harus berupa %s"
)

func Bind(r *http.Request, payload interface{}) error {
	err := json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		errDecode, isDecodeErr := err.(*json.UnmarshalTypeError)
		if isDecodeErr {
			return fmt.Errorf(invalidFormat, errDecode.Field, errDecode.Type.String())
		}
	}
	return err
}
