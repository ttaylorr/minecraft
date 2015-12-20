package protocol

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/ttaylorr/minecraft/protocol/packet"
)

func TestConnectionConstruction(t *testing.T) {
	b := new(bytes.Buffer)
	conn := NewConnection(b)

	assert.NotNil(t, conn, "expected a *Connection object, got nil")
}

func TestPacketReading(t *testing.T) {
	p := []byte{
		0x0f, 0x00, 0x2f, 0x09, 0x6c, 0x6f, 0x63, 0x61,
		0x6c, 0x68, 0x64, 0x73, 0x74, 0x63, 0xdd, 0x01,
	}

	r := bytes.NewReader(p)
	c := NewConnection(r)

	next, err := c.packet()

	assert.Nil(t, err)
	assert.Equal(t, next.ID, 0)
	assert.Equal(t, next.Direction, packet.DirectionServerbound)
	assert.Equal(t, next.Data, p[2:])
}
