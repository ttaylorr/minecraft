package decode_test

import (
	"bytes"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/ttaylorr/minecraft/protocol"
	"github.com/ttaylorr/minecraft/protocol/decode"
	"github.com/ttaylorr/minecraft/protocol/packet"
)

func TestGettingFieldTypes(t *testing.T) {
	v := reflect.TypeOf(struct {
		SomeField string `type:"string"`
	}{})
	field := v.Field(0)

	assert.Equal(t, decode.FieldType(field), "string")
}

func TestGettingFieldDecoder(t *testing.T) {
	v := reflect.TypeOf(struct {
		SomeField string `type:"string"`
	}{})
	field := v.Field(0)

	assert.Equal(t, decode.Types["string"], decode.GetFieldType(field))
}

func TestSettingFieldValue(t *testing.T) {
	v := reflect.ValueOf(&struct {
		SomeField string `type:"string"`
	}{}).Elem()

	assert.Equal(t, v.Field(0).Interface(), "")
	decode.SetFieldValue(v, 0, "some string")
	assert.Equal(t, v.Field(0).Interface(), "some string")
}

func TestPacketDecoding(t *testing.T) {
	p := []byte{
		0x0f, 0x00, 0x2f, 0x09, 0x6c, 0x6f, 0x63, 0x61,
		0x6c, 0x68, 0x64, 0x73, 0x74, 0x63, 0xdd, 0x01,
	}

	r := bytes.NewReader(p)
	c := protocol.NewConnection(r)

	next, err := c.Next()

	assert.Nil(t, err)
	hsk, ok := next.(packet.Handshake)

	assert.True(t, ok)

	assert.Equal(t, hsk.ProtocolVersion, uint64(47))
	assert.Equal(t, hsk.ServerAddress, "localhdst")
	assert.Equal(t, hsk.ServerPort, uint16(25565))
	assert.Equal(t, hsk.NextState, uint64(1))
}
