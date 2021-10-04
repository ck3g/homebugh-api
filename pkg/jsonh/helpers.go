package jsonh

import (
	"encoding/json"
	"reflect"
)

func Equal(j1, j2 []byte) bool {
	if !json.Valid(j1) || !json.Valid(j2) {
		return false
	}

	var json1 interface{}
	var json2 interface{}

	json.Unmarshal(j1, &json1)
	json.Unmarshal(j2, &json2)

	return reflect.DeepEqual(json1, json2)
}
