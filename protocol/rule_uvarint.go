package protocol

import (
	"bytes"
	"encoding/binary"
	"reflect"
)

type UvarintRule struct{}

func (u UvarintRule) AppliesTo(typ reflect.Type) bool {
	return typ.Kind() == reflect.Uint64
}

func (u UvarintRule) Decode(r *bytes.Buffer) (interface{}, error) {
	return binary.ReadUvarint(r)
}
