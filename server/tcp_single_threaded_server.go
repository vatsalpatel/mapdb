package server

import (
	"fmt"
	"log"
	"net"
	"time"

	"github.com/vatsalpatel/mapdb/core"
)

type TCPSingleThreadedServer struct {
	core.IEngine
	Port     int
	listener net.Listener
	channel  chan channelItem
	shutdown chan struct{}
}

type channelItem struct {
	conn net.Conn
	data []byte
}

func NewTCPSingleThreadedServer(port int, engine core.IEngine) *TCPSingleThreadedServer {
	return &TCPSingleThreadedServer{
		IEngine: engine,
		Port:    port,
	}
}

func (s *TCPSingleThreadedServer) Start() error {
	s.channel = make(chan channelItem, 1000)
	s.shutdown = make(chan struct{})
	var err error
	s.listener, err = net.Listen("tcp", fmt.Sprintf(":%d", s.Port))
	log.Println("single threaded tcp server started on port", s.Port)
	if err != nil {
		return err
	}

	go s.worker()
	for {
		conn, err := s.listener.Accept()
		if err != nil {
			continue
		}
		go waitAndRead(s.channel, conn)
	}
}

func (s *TCPSingleThreadedServer) Stop() error {
	s.shutdown <- struct{}{}
	close(s.channel)

	err := s.listener.Close()
	if err != nil {
		return err
	}
	<-time.After(time.Second)
	return s.IEngine.Shutdown()
}

func (s *TCPSingleThreadedServer) handle(item channelItem) {
	resp := s.IEngine.Handle(item.data)

	_, err := item.conn.Write(resp)
	if err != nil {
		item.conn.Close()
	}
}

func (s *TCPSingleThreadedServer) worker() {
	for {
		select {
		case <-s.shutdown:
			return
		case item := <-s.channel:
			s.handle(item)
		}
	}
}

func waitAndRead(channel chan channelItem, conn net.Conn) {
	buf := make([]byte, 1024)
	for {
		read, err := conn.Read(buf[:])
		if err != nil {
			return
		}
		if read == 0 {
			continue
		}
		channel <- channelItem{conn: conn, data: buf}
	}
}
