package protocol

import (
	"reflect"

	"github.com/ttaylorr/minecraft/protocol/packet"
)

var (
	Packets = map[packet.Direction]map[State]map[int]reflect.Type{
		packet.DirectionServerbound: map[State]map[int]reflect.Type{
			HandshakeState: map[int]reflect.Type{
				0x00: reflect.TypeOf(packet.Handshake{}),
			},
			StatusState: map[int]reflect.Type{
				0x00: reflect.TypeOf(packet.StatusRequest{}),
				0x01: reflect.TypeOf(packet.StatusPing{}),
			},
			LoginState: map[int]reflect.Type{
				0x00: reflect.TypeOf(packet.LoginStart{}),
				0x01: reflect.TypeOf(packet.LoginEncryptionResponse{}),
				0x03: reflect.TypeOf(packet.LoginSetCompression{}),
			},
		},
		packet.DirectionClientbound: map[State]map[int]reflect.Type{
			StatusState: map[int]reflect.Type{
				0x00: reflect.TypeOf(packet.StatusResponse{}),
				0x01: reflect.TypeOf(packet.StatusPong{}),
			},
			LoginState: map[int]reflect.Type{
				0x01: reflect.TypeOf(packet.LoginEncryptionRequest{}),
				0x02: reflect.TypeOf(packet.LoginSuccess{}),
				0x03: reflect.TypeOf(packet.LoginSetCompression{}),
			},
		},
	}
)

func GetPacket(d packet.Direction, s State, id int) reflect.Type {
	return Packets[d][s][id]
}
