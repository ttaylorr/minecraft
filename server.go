package main

import (
	"io"
	"math/rand"
	"net"

	"github.com/ttaylorr/minecraft/protocol"
	"github.com/ttaylorr/minecraft/protocol/packet"
)

func main() {
	conn, _ := net.Listen("tcp", "0.0.0.0:25565")
	for {
		client, _ := conn.Accept()
		go handleConnection(protocol.NewConnection(client))
	}
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
			resp.Status.Description.Text = "Hello from Golang!"

			c.Write(resp)
		case packet.StatusPing:
			pong := packet.StatusPong{}
			pong.Payload = t.Payload

			c.Write(pong)
		}
	}
}
