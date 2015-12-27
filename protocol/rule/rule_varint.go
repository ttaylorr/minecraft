package rule

import (
	"bytes"
	"encoding/binary"
	"reflect"

	"github.com/ttaylorr/minecraft/util"
)

type VarintRule struct{}

func (v VarintRule) AppliesTo(typ reflect.Type) bool {
	return typ.Kind() == reflect.Int32
}

func (v VarintRule) Decode(r *bytes.Buffer) (interface{}, error) {
	i, err := binary.ReadVarint(r)
	return int32(i), err
}

func (u VarintRule) Encode(v interface{}) ([]byte, error) {
	i32, ok := v.(int32)
	if !ok {
		return nil, ErrorMismatchedType("int32", v)
	}

	return util.Varint(i32), nil
}
