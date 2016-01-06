package player

import (
	"github.com/ttaylorr/minecraft/protocol"
	"github.com/ttaylorr/minecraft/protocol/packet"
)

type Player struct {
	Username string
	UUID     string

	conn *protocol.Connection
}

func New(conn *protocol.Connection) *Player {
	return &Player{
		conn: conn,
	}
}

func (p *Player) SetState(state protocol.State) {
	p.conn.SetState(state)
}

func (p *Player) Next() (interface{}, error) {
	return p.conn.Next()
}

func (p *Player) Write(h packet.Holder) (int, error) {
	return p.conn.Write(h)
}
