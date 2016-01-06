package server

import (
	"crypto/rsa"
	"fmt"
	"net"
	"reflect"

	"github.com/ttaylorr/minecraft/listener"
	"github.com/ttaylorr/minecraft/player"
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

	priv, err := NewPrivateKey()
	if err != nil {
		return nil, err
	}

	auth, err := auth.New(priv)
	if err != nil {
		return nil, err
	}

	return &Server{
		Handshaker:    handshake.New(),
		Authenticator: auth,

		dispatch:   listener.NewDispatcher(),
		privateKey: priv,
		conn:       conn,
	}, nil
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

func (s *Server) acceptConnections(errs chan error) {
	for {
		// TODO: move to goroutine
		client, err := s.nextConnection()
		if err != nil {
			errs <- err
			continue
		}

		player, err := s.handshakeOrAuth(client)
		if err != nil {
			errs <- err
			continue
		} else if player != nil {
			s.dispatch.PlayerJoin(player)
		}
	}
}

func (s *Server) nextConnection() (*protocol.Connection, error) {
	ln, err := s.conn.Accept()
	if err != nil {
		return nil, err
	}

	return protocol.NewConnection(ln), nil
}

func (s *Server) handshakeOrAuth(client *protocol.Connection) (*player.Player, error) {
	p, err := client.Next()
	if err != nil {
		return nil, err
	}

	hsk, ok := p.(packet.Handshake)
	if !ok {
		return nil, fmt.Errorf("expected packet.Handshake, received %s", reflect.TypeOf(p))
	}

	state := protocol.State(uint8(hsk.NextState))
	client.SetState(state)

	switch state {
	case protocol.HandshakeState:
		s.Handshaker.Handshake(client)
		return nil, nil
	case protocol.LoginState:
		player, err := s.Authenticator.Login(client)
		if err != nil {
			return nil, fmt.Errorf("unable to login player: %s", err)
		}

		return player, nil
	default:
		return nil, fmt.Errorf("neither handshake nor auth state received")
	}
}
