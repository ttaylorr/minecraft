package protocol

import (
	"bytes"
	"encoding/binary"
	"reflect"
)

type VarintRule struct{}

func (v VarintRule) AppliesTo(typ reflect.Type) bool {
	return typ.Kind() == reflect.Int64
}

func (v VarintRule) Decode(r *bytes.Buffer) (interface{}, error) {
	return binary.ReadVarint(r)
}
