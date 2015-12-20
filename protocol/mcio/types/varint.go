package types

import (
	"bytes"
	"encoding/binary"
)

type Varint struct{}

func (v Varint) Decode(r *bytes.Buffer) (interface{}, error) {
	return binary.ReadVarint(r)
}
