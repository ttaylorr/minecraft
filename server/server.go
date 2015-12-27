package server

import (
	"crypto/rand"
	"crypto/rsa"
	"fmt"
	"net"
)

type Server struct {
	conn       net.Listener
	privateKey *rsa.PrivateKey
}

func New(port int) (*Server, error) {
	conn, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return nil, err
	}

	s := &Server{
		conn: conn,
	}

	if err = s.GeneratePrivateKey(); err != nil {
		return nil, err
	}

	return s, nil
}

func (s *Server) Start() error {
	errs := make(chan error)
	defer close(errs)

	go s.AcceptConnections(errs)

	for {
		select {
		case err := <-errs:
			return err
		}
	}

	return nil
}

func (s *Server) AcceptConnections(errs chan error) {
	for {
		_, err := s.conn.Accept()
		if err != nil {
			errs <- err
			return
		}
	}
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
