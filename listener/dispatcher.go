package listener

import (
	"sync"

	"github.com/ttaylorr/minecraft/player"
)

// A Dispatcher is a type used for distributing certain player-centric events to
// Listeners that are designated to listen to them.
type Dispatcher struct {
	lmu *sync.RWMutex

	// TODO: investigate cost/benefit perf by splitting into `pjl
	// []PlayerJoinListener` and `pll []PlayerLeaveListener`.  Will
	// definitely speed up execution, but adds another mutex and test
	listeners []listener
}

// NewDispatcher instantiates and returns a pointer to a new Dispatcher. The
// Dispatcher that is returned is initialized with no Listeners.
func NewDispatcher() *Dispatcher {
	return &Dispatcher{
		lmu: new(sync.RWMutex),
	}
}

// Register registers a listener within the *Dispatcher. Once that listener is
// registered, it will recieve new events.
//
// A call to Register() blocks execution until the Lock is granted, and the new
// listener can be appended.
func (d *Dispatcher) Register(l listener) {
	d.lmu.Lock()
	defer d.lmu.Unlock()

	d.listeners = append(d.listeners, l)
}

// PlayerJoin tells all PlayerJoinListeners that a player has joined by calling
// `func OnPlayerJoin` with the player as its argument.
func (d *Dispatcher) PlayerJoin(player *player.Player) {
	d.lmu.RLock()
	defer d.lmu.RUnlock()

	for _, l := range d.listeners {
		if pl, ok := l.(PlayerJoinListener); ok {
			pl.OnPlayerJoin(player)
		}
	}
}

// PlayerLeave tells all PlayerLeaveListeners that a player has left by calling
// `func OnPlayerLeave` with the player as its argument.
func (d *Dispatcher) PlayerLeave(player *player.Player) {
	d.lmu.RLock()
	defer d.lmu.RUnlock()

	for _, l := range d.listeners {
		if pl, ok := l.(PlayerLeaveListener); ok {
			pl.OnPlayerLeave(player)
		}
	}
}
