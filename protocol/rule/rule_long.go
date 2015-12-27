package rule

import (
	"bytes"
	"encoding/binary"
	"reflect"
)

type LongRule struct{}

func (l LongRule) AppliesTo(typ reflect.Type) bool {
	return typ.Kind() == reflect.Int64
}

func (l LongRule) Decode(r *bytes.Buffer) (interface{}, error) {
	var n int64
	if err := binary.Read(r, ByteOrder, &n); err != nil {
		return nil, err
	}

	return n, nil
}

func (l LongRule) Encode(v interface{}) ([]byte, error) {
	i64, ok := v.(int64)
	if !ok {
		return nil, ErrorMismatchedType("int64", v)
	}

	buf := new(bytes.Buffer)
	if err := binary.Write(buf, ByteOrder, i64); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}
