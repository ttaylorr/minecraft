package types

import (
	"bytes"
	"encoding/binary"
)

func UVarintDecoder(r *bytes.Buffer) (interface{}, error) {
	return binary.ReadUvarint(r)
}
