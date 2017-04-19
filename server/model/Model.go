package model

import "sync"

type Server struct {
	sync.RWMutex
	previousMigrationMap map[string]int64 //timestamp when last migrated to avoid thrashing
	migrationStatus      map[string]bool  // To avoid duplicate migration requests
}

func NewServer() *Server {
	return &Server{
		previousMigrationMap: make(map[string]int64),
		migrationStatus:      make(map[string]bool),
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

func (server *Server) GetContainerMigrationStatus(containerId string) bool {
	server.RLock()
	defer server.RUnlock()
	status, present := server.migrationStatus[containerId]
	if present {
		return status
	}
	return false
}

func (server *Server) SetContainerMigrationStatus(containerId string, status bool) {
	server.Lock()
	defer server.Unlock()
	server.migrationStatus[containerId] = status
}

func (server *Server) SetPreviousContainerMigrationTime(containerId string, timestamp int64) {
	server.Lock()
	defer server.Unlock()
	server.previousMigrationMap[containerId] = timestamp
}
