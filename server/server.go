package server

import (
	"crypto/rand"
	"crypto/rsa"
	"fmt"
	"net"

	"github.com/ttaylorr/minecraft/listener"
	"github.com/ttaylorr/minecraft/protocol"
	"github.com/ttaylorr/minecraft/protocol/packet"
	"github.com/ttaylorr/minecraft/server/auth"
	"github.com/ttaylorr/minecraft/server/handshake"
)

type Server struct {
	Handshaker    *handshake.Handshaker
	Authenticator *auth.Authenticator

	dispatch   *listener.Dispatcher
	privateKey *rsa.PrivateKey
	conn       net.Listener
}

func New(port int) (*Server, error) {
	conn, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return nil, err
	}

	s := &Server{
		Handshaker: handshake.New(),

		dispatch: listener.NewDispatcher(),
		conn:     conn,
	}

	s.GeneratePrivateKey()
	s.Authenticator = auth.New(s.privateKey)

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

		p, err := client.Next()
		if err != nil {
			errs <- err
			continue
		}

		if hsk, ok := p.(packet.Handshake); ok {
			state := protocol.State(uint8(hsk.NextState))
			client.SetState(state)

			switch state {
			case protocol.StatusState:
				s.Handshaker.Handshake(client)
			case protocol.LoginState:
				player, err := s.Authenticator.Login(client)
				if err != nil {
					fmt.Println("unable to login player: %v", err)
				}
				s.dispatch.PlayerJoin(player)
			}
		}
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
