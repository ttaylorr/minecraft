package util_test

import (
	"bytes"
	"io"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/ttaylorr/minecraft/util"
)

func TestByteReaderReadsOneByte(t *testing.T) {
	source := bytes.NewReader([]byte{1, 2, 3})
	br := util.ByteReader{source}

	first, err := br.ReadByte()
	assert.Nil(t, err)
	assert.Equal(t, first, byte(1))

	second, err := br.ReadByte()
	assert.Nil(t, err)
	assert.Equal(t, second, byte(2))

	third, err := br.ReadByte()
	assert.Nil(t, err)
	assert.Equal(t, third, byte(3))
}

func TestEmptyByteReaderEOFs(t *testing.T) {
	source := bytes.NewReader([]byte{})
	br := util.ByteReader{source}

	_, err := br.ReadByte()
	assert.Equal(t, io.EOF, err)
}
