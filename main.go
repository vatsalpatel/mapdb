package main

import (
	"flag"

	"github.com/vatsalpatel/radish/core"
	"github.com/vatsalpatel/radish/server"
	"github.com/vatsalpatel/radish/store"
)

func setupFlags(port *int) {
	flag.IntVar(port, "port", 6379, "Port to listen on")
	flag.Parse()
}

func main() {
	var port int
	setupFlags(&port)

	storage := store.NewMemory()
	engine := core.NewEngine(storage)
	s := server.NewTCPAsyncServer(port, engine)
	s.Start()
	defer s.Stop()
}
