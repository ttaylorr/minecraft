package types

import (
	"bytes"
	"encoding/binary"
)

func UnsignedShortDecoder(r *bytes.Buffer) (interface{}, error) {
	buf := make([]byte, 2)
	if _, err := r.Read(buf); err != nil {
		return nil, err
	}

	return binary.BigEndian.Uint16(buf), nil
}
