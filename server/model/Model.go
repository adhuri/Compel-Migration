package model

import "sync"

type Server struct {
	sync.RWMutex
	previousMigrationMap map[string]int64 //timestamp when last migrated to avoid thrashing
	isMigrating          bool             // To avoid multiple containers migrating at same time CHAOS
}

func NewServer() *Server {
	return &Server{
		previousMigrationMap: make(map[string]int64),
		isMigrating:          false,
	}
}

func (server *Server) GetPreviousContainerMigrationTime(containerId string) int64 {
	server.RLock()
	defer server.RUnlock()
	timestamp, present := server.previousMigrationMap[containerId]
	if present {
		return timestamp
	}
	return 0
}

func (server *Server) GetMigrationStatus() bool {
	server.RLock()
	defer server.RUnlock()
	return server.isMigrating
}

func (server *Server) SetMigrationStatus(status bool) {
	server.Lock()
	defer server.Unlock()
	server.isMigrating = status
}

func (server *Server) SetPreviousContainerMigrationTime(containerId string, timestamp int64) {
	server.Lock()
	defer server.Unlock()
	server.previousMigrationMap[containerId] = timestamp
}
