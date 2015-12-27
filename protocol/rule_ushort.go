package protocol

import (
	"bytes"
	"encoding/binary"
	"reflect"
)

var (
	ByteOrder = binary.BigEndian
)

type UshortRule struct{}

func (u UshortRule) AppliesTo(typ reflect.Type) bool {
	return typ.Kind() == reflect.Uint16
}

func (u UshortRule) Decode(r *bytes.Buffer) (interface{}, error) {
	buf := u.Buffer()
	if _, err := r.Read(buf); err != nil {
		return nil, err
	}

	return ByteOrder.Uint16(buf), nil
}

func (u UshortRule) Encode(v interface{}) ([]byte, error) {
	u16, ok := v.(uint16)
	if !ok {
		return nil, ErrorMismatchedType("uint16", v)
	}

	buf := u.Buffer()
	ByteOrder.PutUint16(buf, u16)

	return buf, nil
}

func (u UshortRule) Buffer() []byte {
	return make([]byte, 2)
}
