package packet

import "reflect"

var (
	Packets = map[int]reflect.Type{
		0x00: reflect.TypeOf(Handshake{}),
	}
)
