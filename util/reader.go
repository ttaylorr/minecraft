package util

import "io"

type ByteReader struct {
	io.Reader
}

func (br ByteReader) ReadByte() (byte, error) {
	buf := make([]byte, 1)
	if _, err := br.Read(buf); err != nil {
		return 0, err
	}

	return buf[0], nil
}
