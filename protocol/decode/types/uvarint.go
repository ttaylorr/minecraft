package types

import (
	"bytes"
	"encoding/binary"
)

type Uvarint struct{}

func (u Uvarint) Decode(r *bytes.Buffer) (interface{}, error) {
	return binary.ReadUvarint(r)
}
