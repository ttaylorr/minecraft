package packet

import (
	"github.com/ttaylorr/minecraft/chat"
	"github.com/ttaylorr/minecraft/protocol/types"
)

type (
	StatusRequest struct{}

	StatusResponse struct {
		Status struct {
			Version struct {
				Name     string `json:"name"`
				Protocol int    `json:"protocol"`
			} `json:"version"`

			Players struct {
				Max    int `json:"max"`
				Online int `json:"online"`
			} `json:"players"`

			Description chat.TextComponent `json:"description"`
		}
	}

	StatusPing struct {
		Payload types.Long
	}

	StatusPong struct {
		Payload types.Long
	}
)

func (r StatusRequest) ID() int  { return 0x00 }
func (r StatusResponse) ID() int { return 0x00 }
func (p StatusPing) ID() int     { return 0x01 }
func (p StatusPong) ID() int     { return 0x01 }
