package handshake

import "github.com/ttaylorr/minecraft/chat"

type MOTD struct {
	Version        string
	PlayersMax     int
	PlayersCurrent int
	Message        chat.TextComponent
}
