package types

import (
	"bytes"
	"encoding/binary"
)

func VarintDecoder(r *bytes.Buffer) (interface{}, error) {
	return binary.ReadVarint(r)
}
