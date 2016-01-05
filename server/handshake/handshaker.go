package handshake

import (
	"github.com/ttaylorr/minecraft/chat"
	"github.com/ttaylorr/minecraft/protocol"
	"github.com/ttaylorr/minecraft/protocol/packet"
)

type Handshaker struct {
	motd *MOTD
}

func New() *Handshaker {
	motd := &MOTD{
		Version: "1.8.8",
		Message: chat.TextComponent{
			"Hello, world!",
			chat.Component{
				Bold:  true,
				Color: chat.ColorRed,
			},
		},
	}

	return &Handshaker{motd}
}

func (h *Handshaker) Handshake(conn *protocol.Connection) {
	for {
		p, err := conn.Next()
		if err != nil {
			return
		}

		switch t := p.(type) {
		case packet.StatusRequest:
			conn.Write(h.GetStatus())
		case packet.StatusPing:
			conn.Write(h.GetPong(t))

			return
		}
	}
}

func (h *Handshaker) GetPong(ping packet.StatusPing) packet.StatusPong {
	return packet.StatusPong{
		Payload: ping.Payload,
	}
}

func (h *Handshaker) GetStatus() packet.StatusResponse {
	r := packet.StatusResponse{}

	r.Status.Version.Name = h.motd.Version
	r.Status.Version.Protocol = 47
	r.Status.Players.Max = h.motd.PlayersMax
	r.Status.Players.Online = h.motd.PlayersCurrent
	r.Status.Description = h.motd.Message

	return r
}
