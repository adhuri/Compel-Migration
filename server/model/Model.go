package model

import "sync"

type Server struct {
	sync.RWMutex
}

func NewServer() *Server {
	return &Server{}
}
