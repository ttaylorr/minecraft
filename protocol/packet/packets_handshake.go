package packet

import "github.com/ttaylorr/minecraft/protocol/types"

type Handshake struct {
	ProtocolVersion types.UVarint
	ServerAddress   types.String
	ServerPort      types.UShort
	NextState       types.UVarint
}

func (h Handshake) ID() int {
	return 0x00
}
