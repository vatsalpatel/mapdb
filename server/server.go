package server

import "github.com/vatsalpatel/mapdb/core"

type IServer interface {
	core.IEngine
	Start() error
	Stop() error
}
