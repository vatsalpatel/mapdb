package server

import (
	"fmt"
	"log"
	"net"
	"os"

	"github.com/vatsalpatel/mapdb/core"
)

type TCPSingleThreadedServer struct {
	core.IEngine
	Port     int
	listener net.Listener
	channel  chan channelItem
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
	var err error
	s.listener, err = net.Listen("tcp", fmt.Sprintf(":%d", s.Port))
	log.Println("single threaded tcp server started on port", s.Port)
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
		go waitAndRead(s.channel, conn)
	}
}

func (s *TCPSingleThreadedServer) Stop() error {
	return s.listener.Close()
}

func (s *TCPSingleThreadedServer) handle(item channelItem) {
	resp := s.IEngine.Handle(item.data)

	_, err := item.conn.Write(resp)
	if err != nil {
		item.conn.Close()
	}
}

func (s *TCPSingleThreadedServer) worker() {
	for item := range s.channel {
		s.handle(item)
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
