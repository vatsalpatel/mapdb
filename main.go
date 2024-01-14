package main

import (
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"

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
		engine := core.NewEngine(store.NewThreadSafeMemory[*core.Item](), persistentStorage)
		s = server.NewTCPAsyncServer(port, engine)
	case 1:
		engine := core.NewEngine(store.NewMemory[*core.Item](), persistentStorage)
		s = server.NewTCPSingleThreadedServer(port, engine)
	case 2:
		engine := core.NewEngine(store.NewMemory[*core.Item](), persistentStorage)
		s = server.NewTCPSyncServer(port, engine)
	}

	go func() {
		if err := s.Start(); err != nil {
			log.Fatal(err)
		}
	}()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)
	<-sigChan
	log.Println("Shutting down...")
	if err := s.Stop(); err != nil {
		log.Fatal(err)
	}
	log.Println("Shutdown complete")
}
