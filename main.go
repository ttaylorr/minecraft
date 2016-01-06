package main

import "github.com/ttaylorr/minecraft/server"

func main() {
	server, err := server.New(25565)
	if err != nil {
		panic(err)
	}

	server.Start()
}
