package protocol_test

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/ttaylorr/minecraft/protocol"
	"github.com/ttaylorr/minecraft/protocol/packet"
)

func TestDecoderInitialization(t *testing.T) {
	var d *protocol.Dealer

	d = protocol.NewDealer()
	assert.IsType(t, d, &protocol.Dealer{})
	assert.Empty(t, d.Rules)

	d = protocol.NewDealer(protocol.StringRule{})
	assert.IsType(t, d, &protocol.Dealer{})
	assert.Len(t, d.Rules, 1)
	assert.Equal(t, d.Rules[0], protocol.StringRule{})
}

func TestDefaultDecoderInitialization(t *testing.T) {
	d := protocol.DefaultDealer()
	assert.IsType(t, d, &protocol.Dealer{})
	assert.Len(t, d.Rules, 4)
}

func TestPacketDecoding(t *testing.T) {
	d := protocol.DefaultDealer()
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

	assert.Equal(t, hsk.ProtocolVersion, uint64(47))
	assert.Equal(t, hsk.ServerAddress, "localhdst")
	assert.Equal(t, hsk.ServerPort, uint16(25565))
	assert.Equal(t, hsk.NextState, uint64(1))
}

func TestPacketEncoding(t *testing.T) {
	handshake := packet.Handshake{
		ProtocolVersion: uint64(47),
		ServerAddress:   "localhdst",
		ServerPort:      uint16(25565),
		NextState:       uint64(1),
	}

	d := protocol.DefaultDealer()
	e, err := d.Encode(handshake)

	assert.Nil(t, err)
	assert.Equal(t, e, []byte{
		0x0f, 0x00, 0x2f, 0x09, 0x6c, 0x6f, 0x63, 0x61,
		0x6c, 0x68, 0x64, 0x73, 0x74, 0x63, 0xdd, 0x01,
	})
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

func TestFindingSingleMatchingRule(t *testing.T) {
	rule := &protocol.StringRule{}
	d := protocol.NewDealer(rule)

	typ := reflect.TypeOf(struct {
		Field string
	}{}).Field(0).Type

	assert.Equal(t, d.GetRule(typ), rule)
}

func TestFindingNoMatchingRules(t *testing.T) {
	d := protocol.NewDealer()
	typ := reflect.TypeOf("")

	assert.Nil(t, d.GetRule(typ))
}
