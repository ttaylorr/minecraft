package server

import (
	"crypto/rand"
	"crypto/rsa"
	"fmt"
	"net"

	"github.com/ttaylorr/minecraft/player"
	"github.com/ttaylorr/minecraft/protocol"
	"github.com/ttaylorr/minecraft/server/handshake"
)

type Server struct {
	Handshaker *handshake.Handshaker

	privateKey *rsa.PrivateKey

	conn net.Listener
}

func New(port int) (*Server, error) {
	conn, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return nil, err
	}

	s := &Server{
		Handshaker: handshake.New(),

		conn: conn,
	}

	s.GeneratePrivateKey()

	return s, nil
}

func (s *Server) acceptConnections(errs chan error) {
	for {
		ln, err := s.conn.Accept()
		if err != nil {
			errs <- err
			continue
		}

		client := protocol.NewConnection(ln)
		player := player.New(client)

		s.Handshaker.OnPlayerJoin(player)
	}
}

func (s *Server) Start() error {
	errs := make(chan error)
	defer close(errs)

	go s.acceptConnections(errs)

	for {
		select {
		case err := <-errs:
			return err
		}
	}

	return nil
}

func (s *Server) GeneratePrivateKey() error {
	priv, err := rsa.GenerateKey(rand.Reader, PrivateKeySize)
	if err != nil {
		return err
	}

	s.privateKey = priv
	s.privateKey.Precompute()

	return nil
}
