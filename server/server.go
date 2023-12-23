package server

import "github.com/vatsalpatel/radish/core"

type IServer interface {
	core.IEngine
	Start(port int) error
	Stop() error
}
