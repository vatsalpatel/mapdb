package server

import (
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"sync"
	"time"

	"github.com/vatsalpatel/mapdb/core"
)

type TCPSyncServer struct {
	core.IEngine
	Port          int
	listener      net.Listener
	connections   map[net.Addr]net.Conn
	connectionsMu sync.RWMutex
}

func NewTCPSyncServer(port int, engine core.IEngine) *TCPSyncServer {
	return &TCPSyncServer{
		IEngine: engine,
		Port:    port,
	}
}

func (s *TCPSyncServer) Start() error {
	s.connections = make(map[net.Addr]net.Conn)

	var err error
	s.listener, err = net.Listen("tcp", fmt.Sprintf(":%d", s.Port))
	log.Println("sync tcp server started on port", s.Port)
	if err != nil {
		return err
	}
	defer s.Stop()

	go s.worker()
	for {
		conn, err := s.listener.Accept()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		s.connectionsMu.Lock()
		s.connections[conn.RemoteAddr()] = conn
		s.connectionsMu.Unlock()
	}
}

func (s *TCPSyncServer) Stop() error {
	return s.listener.Close()
}

func (s *TCPSyncServer) handle(conn net.Conn) error {
	buf := make([]byte, 1024)
	conn.SetReadDeadline(time.Now().Add(100 * time.Nanosecond))
	read, err := conn.Read(buf[:])
	if err != nil {
		if err == io.EOF {
			return err
		}
		return nil
	}

	if read == 0 {
		return nil
	}

	resp := s.IEngine.Handle(buf)

	_, err = conn.Write(resp)
	if err != nil {
		conn.Close()
	}
	return nil
}

func (s *TCPSyncServer) worker() {
	for {
		s.connectionsMu.RLock()
		for _, conn := range s.connections {
			err := s.handle(conn)
			if err != nil {
				conn.Close()
				delete(s.connections, conn.RemoteAddr())
			}
		}
		s.connectionsMu.RUnlock()
	}
}
