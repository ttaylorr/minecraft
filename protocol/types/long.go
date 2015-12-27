package types

import (
	"encoding/binary"
	"io"
)

type Long int64

func (_ Long) Decode(r io.Reader) (interface{}, error) {
	var l Long
	err := binary.Read(r, ByteOrder, &l)

	return l, err
}

func (l Long) Encode(w io.Writer) (int, error) {
	return 8, binary.Write(w, ByteOrder, l)
}
