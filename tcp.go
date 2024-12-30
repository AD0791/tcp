package main

import (
	"fmt"
	"net"
)

type Server struct {
	listenAddr string
	listener   net.Listener
	quitCh     chan struct{}
}

func NewTCPServer(address string) *Server {
	return &Server{
		listenAddr: address,
	}
}

func (s *Server) Start() error {
	ln, err := net.Listen("tcp", s.listenAddr)
	if err != nil {
		return err
	}
	defer ln.Close()
	s.listener = ln

	<-s.quitCh
	return nil
}

func main() {
	fmt.Println("tcp server")
}
