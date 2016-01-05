package listener

import (
	"sync"

	"github.com/ttaylorr/minecraft/player"
)

type Dispatcher struct {
	lmu       *sync.Mutex
	listeners []Listener
}

func NewDispatcher() *Dispatcher {
	return &Dispatcher{}
}

func (d *Dispatcher) Register(l Listener) {
	d.lmu.Lock()
	defer d.lmu.Unlock()

	d.listeners = append(d.listeners, l)
}

func (d *Dispatcher) PlayerJoin(player *player.Player) {
	for _, l := range d.listeners {
		if pl, ok := l.(PlayerJoinListener); ok {
			pl.OnPlayerJoin(player)
		}
	}
}

func (d *Dispatcher) PlayerLeave(player *player.Player) {
	for _, l := range d.listeners {
		if pl, ok := l.(PlayerLeaveListener); ok {
			pl.OnPlayerLeave(player)
		}
	}
}
