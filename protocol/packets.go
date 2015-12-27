package protocol

import (
	"reflect"

	"github.com/ttaylorr/minecraft/protocol/packet"
)

var (
	Packets = map[packet.Direction]map[State]map[int]reflect.Type{
		packet.DirectionServerbound: map[State]map[int]reflect.Type{
			EmptyState: map[int]reflect.Type{
				0x00: reflect.TypeOf(packet.Handshake{}),
			},
			StatusState: map[int]reflect.Type{
				0x00: reflect.TypeOf(packet.StatusRequest{}),
				0x01: reflect.TypeOf(packet.StatusPing{}),
			},
		},
		packet.DirectionClientbound: map[State]map[int]reflect.Type{
			StatusState: map[int]reflect.Type{
				0x00: reflect.TypeOf(packet.StatusResponse{}),
				0x01: reflect.TypeOf(packet.StatusPong{}),
			},
		},
	}
)

func GetPacket(d packet.Direction, s State, id int) reflect.Type {
	return Packets[d][s][id]
}
