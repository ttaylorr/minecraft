package util

import "encoding/binary"

var (
	maxVarintLength = 5
)

func Varint(n int64) []byte {
	buf := make([]byte, maxVarintLength)
	l := binary.PutVarint(buf, n)

	return buf[:l]
}

func Uvarint(n uint64) []byte {
	buf := make([]byte, maxVarintLength)
	l := binary.PutUvarint(buf, n)

	return buf[:l]
}
