package protocol

import (
	"bytes"
	"encoding/binary"
	"errors"
	"reflect"
)

var (
	ErrorMismatchedStringLength = errors.New("fewer bytes available than string length")
)

type StringRule struct{}

func (sr StringRule) AppliesTo(typ reflect.Type) bool {
	return typ.Kind() == reflect.String
}

func (sr StringRule) Decode(r *bytes.Buffer) (interface{}, error) {
	size, err := binary.ReadUvarint(r)
	if err != nil {
		return nil, err
	}

	str := make([]byte, int(size))
	if read, err := r.Read(str); err != nil {
		return nil, err
	} else if read != int(size) {
		return nil, ErrorMismatchedStringLength
	}

	return string(str), nil

}
