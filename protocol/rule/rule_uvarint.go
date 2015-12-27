package rule

import (
	"bytes"
	"encoding/binary"
	"reflect"

	"github.com/ttaylorr/minecraft/util"
)

type UvarintRule struct{}

func (u UvarintRule) AppliesTo(typ reflect.Type) bool {
	return typ.Kind() == reflect.Uint64
}

func (u UvarintRule) Decode(r *bytes.Buffer) (interface{}, error) {
	return binary.ReadUvarint(r)
}

func (u UvarintRule) Encode(v interface{}) ([]byte, error) {
	u64, ok := v.(uint64)
	if !ok {
		return nil, ErrorMismatchedType("uint64", v)
	}

	return util.Uvarint(u64), nil
}
