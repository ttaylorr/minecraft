package types

import (
	"encoding/binary"
	"io"

	"github.com/ttaylorr/minecraft/util"
)

type ByteArray []byte

func (_ ByteArray) Decode(r io.Reader) (interface{}, error) {
	l, err := binary.ReadUvarint(util.ByteReader{r})
	if err != nil {
		return nil, err
	}

	buf := make([]byte, int(l))
	_, err = r.Read(buf)
	if err != nil {
		return nil, err
	}

	return buf, nil
}

func (b ByteArray) Encode(w io.Writer) (n int, err error) {
	n, err = w.Write(util.Uvarint(uint32(len(b))))
	if err != nil {
		return
	}

	_, err = w.Write(b)
	n += len(b)

	return
}
