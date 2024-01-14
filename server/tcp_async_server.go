package server

import (
	"fmt"
	"log"
	"net"

	"github.com/vatsalpatel/mapdb/core"
)

type TCPAsyncServer struct {
	core.IEngine
	Port     int
	listener net.Listener
}

func NewTCPAsyncServer(port int, engine core.IEngine) *TCPAsyncServer {
	return &TCPAsyncServer{
		IEngine: engine,
		Port:    port,
	}
}

func (s *TCPAsyncServer) Start() error {
	var err error
	s.listener, err = net.Listen("tcp", fmt.Sprintf(":%d", s.Port))
	log.Println("async tcp server started on port", s.Port)
	if err != nil {
		return err
	}
	for {
		conn, err := s.listener.Accept()
		if err != nil {
			log.Printf("Error accepting connection: %v", err)
			continue
		}
		go s.handle(conn)
	}
}

func (s *TCPAsyncServer) Stop() error {
	return s.listener.Close()
}

func (s *TCPAsyncServer) handle(conn net.Conn) {
	buf := make([]byte, 1024)
	for {
		_, err := conn.Read(buf[:])
		if err != nil {
			conn.Close()
			return
		}

		resp := s.IEngine.Handle(buf)

		_, err = conn.Write(resp)
		if err != nil {
			conn.Close()
			return
		}
	}
}
