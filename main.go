package main

import (
	"flag"

	"github.com/vatsalpatel/mapdb/core"
	"github.com/vatsalpatel/mapdb/server"
	"github.com/vatsalpatel/mapdb/store"
)

func setupFlags(port *int, serverType *int) {
	flag.IntVar(port, "port", 6379, "Port to listen on")
	flag.IntVar(serverType, "server-type", 0, "Server type to run")
	flag.Parse()
}

func main() {
	var port, serverType int
	setupFlags(&port, &serverType)

	persistentStorage := store.NewFileStore("dump.rdb")

	var s server.IServer
	switch serverType {
	case 0:
		engine := core.NewEngine(store.NewMemory[*core.Item](), persistentStorage)
		s = server.NewTCPSyncServer(port, engine)
	case 1:
		engine := core.NewEngine(store.NewMemory[*core.Item](), persistentStorage)
		s = server.NewTCPSingleThreadedServer(port, engine)
	case 2:
		engine := core.NewEngine(store.NewThreadSafeMemory[*core.Item](), persistentStorage)
		s = server.NewTCPAsyncServer(port, engine)
	}

	s.Start()
	defer s.Stop()
}
