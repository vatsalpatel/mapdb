package main

import (
	"flag"

	"github.com/vatsalpatel/radish/server"
)

func setupFlags(port *int) {
	flag.IntVar(port, "port", 7379, "Port to listen on")
	flag.Parse()
}

func main() {
	var port int
	setupFlags(&port)
	s := server.NewTCPSyncServer(port)
	s.Start()
	defer s.Stop()
}
