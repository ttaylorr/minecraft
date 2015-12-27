package main

import (
	"io"
	"log"
	"math/rand"

	"github.com/ttaylorr/minecraft/chat"
	"github.com/ttaylorr/minecraft/protocol"
	"github.com/ttaylorr/minecraft/protocol/packet"
	"github.com/ttaylorr/minecraft/server"
)

func main() {
	server, err := server.New(25565)
	if err != nil {
		panic(err)
	}

	if err := server.Start(); err != nil {
		panic(err)
	}

	log.Fatalf("Shutting down gracefully...")
}

func handleConnection(c *protocol.Connection) {
	for {
		p, err := c.Next()
		if err == io.EOF {
			return
		}

		switch t := p.(type) {
		case packet.Handshake:
			state := protocol.State(uint8(t.NextState))
			c.SetState(state)
		case packet.StatusRequest:
			resp := packet.StatusResponse{}
			resp.Status.Version.Name = "1.8.8"
			resp.Status.Version.Protocol = 47
			resp.Status.Players.Max = rand.Intn(100)
			resp.Status.Players.Online = rand.Intn(101)
			resp.Status.Description = chat.TextComponent{
				"Hello from Golang!",

				chat.Component{
					Bold:  true,
					Color: chat.ColorRed,
				},
			}

			c.Write(resp)
		case packet.StatusPing:
			pong := packet.StatusPong{}
			pong.Payload = t.Payload

			c.Write(pong)
			return
		}
	}
}
