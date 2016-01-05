package protocol_test

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/ttaylorr/minecraft/protocol"
	"github.com/ttaylorr/minecraft/protocol/packet"
	"github.com/ttaylorr/minecraft/protocol/types"
)

func TestDecoderInitialization(t *testing.T) {
	d := protocol.NewDealer()
	assert.IsType(t, d, &protocol.Dealer{})
}

func TestErrorForUnknownPacketType(t *testing.T) {
	d := protocol.NewDealer()
	p := &packet.Packet{
		ID: 0x99,
	}

	v, err := d.Decode(p)
	assert.Nil(t, v)
	assert.Equal(t, err, protocol.UnknownPacketError)
}

func TestPacketDecoding(t *testing.T) {
	d := protocol.NewDealer()
	p := &packet.Packet{
		ID: 0x00,
		Data: []byte{
			0x2f, 0x09, 0x6c, 0x6f, 0x63, 0x61, 0x6c, 0x68,
			0x64, 0x73, 0x74, 0x63, 0xdd, 0x01,
		},
	}

	v, err := d.Decode(p)

	assert.Nil(t, err)
	hsk, _ := v.(packet.Handshake)

	assert.IsType(t, hsk, packet.Handshake{})

	assert.Equal(t, types.UVarint(47), hsk.ProtocolVersion)
	assert.Equal(t, types.String("localhdst"), hsk.ServerAddress)
	assert.Equal(t, types.UShort(25565), hsk.ServerPort)
	assert.Equal(t, types.UVarint(1), hsk.NextState)
}

func TestPacketEncoding(t *testing.T) {
	handshake := packet.Handshake{
		ProtocolVersion: 47,
		ServerAddress:   "localhdst",
		ServerPort:      25565,
		NextState:       1,
	}

	d := protocol.NewDealer()
	e, err := d.Encode(handshake)

	assert.Nil(t, err)
	assert.Equal(t, []byte{
		0x0f, 0x00, 0x2f, 0x09, 0x6c, 0x6f, 0x63, 0x61,
		0x6c, 0x68, 0x64, 0x73, 0x74, 0x63, 0xdd, 0x01,
	}, e)
}

func TestFindingHolderWithValidID(t *testing.T) {
	d := protocol.NewDealer()
	pack := &packet.Packet{ID: 0x00}

	holder := d.GetHolderType(pack)
	assert.Equal(t, holder, reflect.TypeOf(packet.Handshake{}))
}

func TestFindingHolderWithInvalidID(t *testing.T) {
	d := protocol.NewDealer()
	pack := &packet.Packet{ID: -1}

	assert.Nil(t, d.GetHolderType(pack))
}
