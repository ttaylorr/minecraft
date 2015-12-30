package handshake

import (
	"github.com/ttaylorr/minecraft/chat"
	"github.com/ttaylorr/minecraft/player"
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

func (h Handshaker) OnPlayerJoin(p *player.Player) {
	h.Handshake(p)
}

func (h *Handshaker) Handshake(player *player.Player) {
	for {
		p, err := player.Next()
		if err != nil {
			return
		}

		switch t := p.(type) {
		case packet.Handshake:
			player.SetState(h.GetState(t))
		case packet.StatusRequest:
			player.Write(h.GetStatus())
		case packet.StatusPing:
			player.Write(h.GetPong(t))

			return
		}
	}
}

func (h *Handshaker) GetState(hsk packet.Handshake) protocol.State {
	return protocol.State(uint8(hsk.NextState))
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
