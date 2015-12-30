package listener

import "github.com/ttaylorr/minecraft/player"

type (
	PlayerJoinListener interface {
		OnPlayerJoin(*player.Player)
	}

	PlayerLeaveListener interface {
		OnPlayerLeave(*player.Player)
	}

	PlayerListener interface {
		PlayerJoinListener
		PlayerLeaveListener
	}
)
