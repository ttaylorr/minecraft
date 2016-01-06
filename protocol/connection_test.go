package protocol

import (
	"bytes"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/ttaylorr/minecraft/protocol/packet"
)

func TestConnectionConstruction(t *testing.T) {
	b := new(bytes.Buffer)
	conn := NewConnection(b)

	assert.IsType(t, conn, &Connection{})
}

func TestSettingState(t *testing.T) {
	b := new(bytes.Buffer)
	conn := NewConnection(b)

	conn.SetState(StatusState)

	assert.Equal(t, conn.d.State, StatusState)
}

func TestPacketReading(t *testing.T) {
	p := []byte{
		0x0f, 0x00, 0x2f, 0x09, 0x6c, 0x6f, 0x63, 0x61,
		0x6c, 0x68, 0x64, 0x73, 0x74, 0x63, 0xdd, 0x01,
	}

	r := bytes.NewBuffer(p)
	c := NewConnection(r)

	next, err := c.packet()

	assert.Nil(t, err)
	assert.Equal(t, next.ID, 0)
	assert.Equal(t, next.Direction, packet.DirectionServerbound)
	assert.Equal(t, next.Data, p[2:])
}

func TestInvalidPacketsGetErrors(t *testing.T) {
	r := bytes.NewBuffer([]byte{
		0x0f, 0x00, 0x2f, 0x09, 0x6c, 0x6f, 0x63, 0x61,
	})
	c := NewConnection(r)

	next, err := c.packet()

	assert.Nil(t, next)
	assert.Equal(t, errors.New("unexpected EOF"), err)
}

func TestInvalidPacketsDontGetDecoded(t *testing.T) {
	p := []byte{
		0x0f, 0x00, // packet payload is too short (0 bytes)
	}

	r := bytes.NewBuffer(p)
	c := NewConnection(r)

	v, err := c.Next()

	assert.Nil(t, v)
	assert.NotNil(t, err)
}
