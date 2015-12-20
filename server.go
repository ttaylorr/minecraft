package main

import (
	"fmt"
	"io"
	"net"

	"github.com/ttaylorr/minecraft/protocol"
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

		fmt.Println(p)
	}
}
