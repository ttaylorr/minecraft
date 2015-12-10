package protocol

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConnectionConstruction(t *testing.T) {
	b := new(bytes.Buffer)
	conn := NewConnection(b)

	assert.NotNil(t, conn, "expected a *Connection object, got nil")
	assert.Equal(t, conn.r, b, "expected Connection to have matching reader")
}

func TestPacketReading(t *testing.T) {
	p := []byte{
		0x0f, 0x00, 0x2f, 0x09, 0x6c, 0x6f, 0x63, 0x61,
		0x6c, 0x68, 0x64, 0x73, 0x74, 0x63, 0xdd, 0x01,
	}

	r := bytes.NewReader(p)
	c := NewConnection(r)

	packet, err := c.Next()

	assert.Nil(t, err)
	assert.Equal(t, packet.ID, 0)
	assert.Equal(t, packet.Direction, DirectionServerbound)
	assert.Equal(t, packet.Data, p[2:])
}
