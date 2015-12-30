package server

import (
	"fmt"
	"net"
)

type Server struct {
	conn *net.Conn
}

func New(port int) (*Server, error) {
	conn, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%d", port))
	if err != nil {
		return nil, err
	}

	return &Server{
		conn: conn,
	}
}
