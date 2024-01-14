package server

import (
	"fmt"
	"io"
	"log"
	"net"
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
	shutdown      chan struct{}
}

func NewTCPSyncServer(port int, engine core.IEngine) *TCPSyncServer {
	return &TCPSyncServer{
		IEngine: engine,
		Port:    port,
	}
}

func (s *TCPSyncServer) Start() error {
	s.connections = make(map[net.Addr]net.Conn)
	s.shutdown = make(chan struct{})

	var err error
	s.listener, err = net.Listen("tcp", fmt.Sprintf(":%d", s.Port))
	log.Println("sync tcp server started on port", s.Port)
	if err != nil {
		return err
	}

	go s.worker()
	for {
		conn, err := s.listener.Accept()
		if err != nil {
			continue
		}
		s.connectionsMu.Lock()
		s.connections[conn.RemoteAddr()] = conn
		s.connectionsMu.Unlock()
	}
}

func (s *TCPSyncServer) Stop() error {
	s.connectionsMu.Lock()
	for addr, conn := range s.connections {
		conn.Close()
		delete(s.connections, addr)
	}
	s.connectionsMu.Unlock()

	s.shutdown <- struct{}{}
	close(s.shutdown)

	err := s.listener.Close()
	if err != nil {
		return err
	}
	<-time.After(time.Second)
	return s.IEngine.Shutdown()
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
		select {
		case <-s.shutdown:
			return
		default:
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
}
