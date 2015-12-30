package main

import (
	"sync"

	"github.com/ttaylorr/minecraft/server"
)

func main() {
	wg := new(sync.WaitGroup)
	wg.Add(1)

	server, err := server.New(wg, 25565)
	if err != nil {
		panic(err)
	}
	server.Start()

	wg.Wait()
}
