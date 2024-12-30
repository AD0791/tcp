package main

import (
	"fmt"
	"log"
	"net"
)

type Message struct {
	from    string
	payload []byte
}

type Server struct {
	listenAddr string
	ln         net.Listener
	quitch     chan struct{}
	msgch      chan Message
}

func NewTCPServer(address string) *Server {
	return &Server{
		listenAddr: address,
		quitch:     make(chan struct{}),
		msgch:      make(chan Message, 10),
	}
}

func (s *Server) Start() error {
	ln, err := net.Listen("tcp", s.listenAddr)
	if err != nil {
		return err
	}
	defer ln.Close()
	s.ln = ln

	go s.acceptLoop()

	<-s.quitch
	close(s.msgch)
	return nil
}

func (s *Server) acceptLoop() {
	for {
		conn, err := s.ln.Accept()
		if err != nil {
			fmt.Println("acccept error: ", err)
			continue
		}
		go s.readLoop(conn)
	}
}

func (s *Server) readLoop(conn net.Conn) {
	defer conn.Close()
	buff := make([]byte, 2048)
	for {
		n, err := conn.Read(buff)
		if err != nil {
			fmt.Println("read error: ", err)
			continue
		}
		s.msgch <- Message{
			from:    conn.RemoteAddr().String(),
			payload: buff[:n],
		}
	}
}

func (s *Server) Stop() {
	close(s.quitch)
}

func main() {
	foundation_server := NewTCPServer(":3000")
	defer foundation_server.Stop()
	log.Println("The tcp server has started in the port :3000")
	go func() {
		for msg := range foundation_server.msgch {
			log.Printf("the message from the tcp connection: (%s):%s ", msg.from, string(msg.payload))
		}
	}()
	log.Fatal(foundation_server.Start())
}
