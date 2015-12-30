package server

import (
	"fmt"
	"net"
	"sync"

	"github.com/ttaylorr/minecraft/player"
	"github.com/ttaylorr/minecraft/protocol"
	"github.com/ttaylorr/minecraft/server/handshake"
)

type Server struct {
	Handshaker *handshake.Handshaker

	wg   *sync.WaitGroup
	conn net.Listener
}

func New(wg *sync.WaitGroup, port int) (*Server, error) {
	conn, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%d", port))
	if err != nil {
		return nil, err
	}

	s := &Server{
		Handshaker: handshake.New(),

		wg:   wg,
		conn: conn,
	}

	return s, nil
}

func (s *Server) Start() {
	go s.acceptConnections()
}

func (s *Server) acceptConnections() {
	for {
		ln, err := s.conn.Accept()
		if err != nil {
			continue
		}

		client := protocol.NewConnection(ln)
		player := player.New(client)

		s.Handshaker.OnPlayerJoin(player)
	}

	s.wg.Done()
}
