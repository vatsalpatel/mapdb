package server

type Server interface {
	Start(port int) error
	Stop() error
}
