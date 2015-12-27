package rule

import (
	"bytes"
	"encoding/binary"
	"reflect"

	"github.com/ttaylorr/minecraft/util"
)

type UvarintRule struct{}

func (u UvarintRule) AppliesTo(typ reflect.Type) bool {
	return typ.Kind() == reflect.Uint32
}

func (u UvarintRule) Decode(r *bytes.Buffer) (interface{}, error) {
	i, err := binary.ReadUvarint(r)
	return uint32(i), err
}

func (u UvarintRule) Encode(v interface{}) ([]byte, error) {
	u32, ok := v.(uint32)
	if !ok {
		return nil, ErrorMismatchedType("uint32", v)
	}

	return util.Uvarint(u32), nil
}
