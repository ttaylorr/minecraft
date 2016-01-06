package listener

import "github.com/ttaylorr/minecraft/player"

type (
	// listener is the "base-type" for all listeners. It exists soley to
	// provide easy-calling for Dispatcher#Register
	listener interface {
		IsListener() bool
	}

	// PlayerJoinListener is a listener that handles players joining the
	// server. `func OnPlayerJoin` will be called once a player has logged
	// into a server.
	PlayerJoinListener interface {
		OnPlayerJoin(*player.Player)
		listener
	}

	// PlayerLeaveListener is a listener that handles players joining the
	// server. `func OnPlayerLeave` will be called before a player has logged
	// off of a server.
	PlayerLeaveListener interface {
		OnPlayerLeave(*player.Player)
		listener
	}
)
