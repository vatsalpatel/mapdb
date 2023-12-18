package server

import (
	"fmt"
	"log"
	"net"
)

type TCPSyncServer struct {
	Port     int
	listener net.Listener
}

func NewTCPAsyncServer(port int) *TCPSyncServer {
	return &TCPSyncServer{
		Port: port,
	}
}

func (s *TCPSyncServer) handle(conn net.Conn) {
	buf := make([]byte, 1024)
	for {
		_, err := conn.Read(buf[:])
		log.Println("received", string(buf))
		if err != nil {
			conn.Close()
			log.Printf("Client disconnected: %v", conn.RemoteAddr())
			return
		}

		_, err = conn.Write(buf[:])
		if err != nil {
			conn.Close()
			log.Printf("Client disconnected: %v", conn.RemoteAddr())
			return
		}
	}
}

func (s *TCPSyncServer) Start() error {
	var err error
	s.listener, err = net.Listen("tcp", fmt.Sprintf(":%d", s.Port))
	if err != nil {
		return err
	}
	defer s.Stop()
	for {
		conn, err := s.listener.Accept()
		defer conn.Close()
		log.Println("client connected", conn.RemoteAddr())
		if err != nil {
			log.Printf("Error accepting connection: %v", err)
			continue
		}
		go s.handle(conn)
	}
}

func (s *TCPSyncServer) Stop() error {
	return s.listener.Close()
}
