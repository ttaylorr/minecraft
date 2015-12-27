package rule

import (
	"bytes"
	"encoding/json"
	"errors"
	"reflect"
)

type JsonRule struct {
	s StringRule
}

func (j JsonRule) AppliesTo(typ reflect.Type) bool {
	return typ.Kind() == reflect.Struct
}

func (j JsonRule) Encode(v interface{}) ([]byte, error) {
	data, err := json.Marshal(v)
	if err != nil {
		return nil, err
	}

	return j.s.Encode(string(data))
}

func (j JsonRule) Decode(buf *bytes.Buffer) (interface{}, error) {
	return nil, errors.New("not yet implemented")
}
