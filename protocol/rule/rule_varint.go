package rule

import (
	"bytes"
	"encoding/binary"
	"reflect"

	"github.com/ttaylorr/minecraft/util"
)

type VarintRule struct{}

func (v VarintRule) AppliesTo(typ reflect.Type) bool {
	return typ.Kind() == reflect.Int64
}

func (v VarintRule) Decode(r *bytes.Buffer) (interface{}, error) {
	return binary.ReadVarint(r)
}

func (u VarintRule) Encode(v interface{}) ([]byte, error) {
	i64, ok := v.(int64)
	if !ok {
		return nil, ErrorMismatchedType("int64", v)
	}

	return util.Varint(i64), nil
}
