package protocol

import (
	"bytes"
	"encoding/binary"
	"reflect"
)

type UshortRule struct{}

func (u UshortRule) AppliesTo(typ reflect.Type) bool {
	return typ.Kind() == reflect.Uint16
}

func (u UshortRule) Decode(r *bytes.Buffer) (interface{}, error) {
	buf := make([]byte, 2)
	if _, err := r.Read(buf); err != nil {
		return nil, err
	}

	return binary.BigEndian.Uint16(buf), nil
}
