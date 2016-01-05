package listener

import "github.com/ttaylorr/minecraft/player"

type (
	Listener interface {
		IsListener() bool
	}

	PlayerJoinListener interface {
		OnPlayerJoin(*player.Player)
		Listener
	}

	PlayerLeaveListener interface {
		OnPlayerLeave(*player.Player)
		Listener
	}
)
