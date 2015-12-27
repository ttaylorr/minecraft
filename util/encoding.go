package util

import "encoding/binary"

var (
	maxVarintLength = 5
)

func Varint(n int32) []byte {
	buf := make([]byte, maxVarintLength)
	l := binary.PutVarint(buf, int64(n))

	return buf[:l]
}

func Uvarint(n uint32) []byte {
	buf := make([]byte, maxVarintLength)
	l := binary.PutUvarint(buf, uint64(n))

	return buf[:l]
}
