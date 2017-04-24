package model

import "sync"

type Server struct {
	sync.RWMutex
	previousMigrationMap             map[string]int64 //timestamp when last migrated to avoid thrashing
	previousSystemMigrationTimeStamp int64            //To avoid System level thrashing
	isMigrating                      bool             // To avoid multiple containers migrating at same time CHAOS
	immovableContainersList          []string
	thrashingThreshold               int64
	cpufp                            int64
	memfp                            int64
	falsePositiveMap                 map[string]*FalsePositive // Check false positives before accepting decision of migrating
}

// To handle False Positives
type FalsePositive struct {
	cpuCount int64
	memCount int64
}

func NewFalsePositive() *FalsePositive {
	return &FalsePositive{
		cpuCount: 0,
		memCount: 0,
	}
}

func NewServer(immovableContainersList []string, threshold int64, cpufp int64, memfp int64) *Server {
	return &Server{
		previousMigrationMap:             make(map[string]int64),
		isMigrating:                      false,
		immovableContainersList:          immovableContainersList,
		thrashingThreshold:               threshold,
		cpufp:                            cpufp,
		memfp:                            memfp,
		falsePositiveMap:                 make(map[string]*FalsePositive),
		previousSystemMigrationTimeStamp: 0,
	}
}

func (server *Server) GetPreviousSystemMigrationTime() int64 {
	server.RLock()
	defer server.RUnlock()
	return server.previousSystemMigrationTimeStamp
}

func (server *Server) SetPreviousSystemMigrationTime(timestamp int64) {
	server.Lock()
	defer server.Unlock()
	server.previousSystemMigrationTimeStamp = timestamp
}

// For system level thrashing
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

// For container level thrashing
func (server *Server) SetPreviousContainerMigrationTime(containerId string, timestamp int64) {
	server.Lock()
	defer server.Unlock()
	server.previousMigrationMap[containerId] = timestamp
}

func (server *Server) CheckIfContainerIsMovable(containerId string) bool {
	server.RLock()
	defer server.RUnlock()

	for _, immovableContainer := range server.immovableContainersList {
		if immovableContainer == containerId {
			return false
		}
	}

	return true
}

func (server *Server) GetThrashingThreshold() int64 {
	server.RLock()
	defer server.RUnlock()
	return server.thrashingThreshold
}

func (server *Server) ResetFalsePositive(containerID string) {
	server.Lock()
	defer server.Unlock()

	server.falsePositiveMap[containerID] = NewFalsePositive()
}

func (server *Server) IncrementFalsePositive(containerID string, metric string) {
	server.Lock()
	defer server.Unlock()

	// Check to see if it is present
	if present := server.falsePositiveMap[containerID]; present == nil {
		server.falsePositiveMap[containerID] = NewFalsePositive()
	}

	switch metric {
	case "cpu":
		server.falsePositiveMap[containerID].cpuCount++
	case "memory":
		server.falsePositiveMap[containerID].memCount++

	}

}

func (server *Server) GetFalsePositiveMap(containerID string, metric string) int64 {
	server.RLock()
	defer server.RUnlock()

	if present := server.falsePositiveMap[containerID]; present == nil {
		server.falsePositiveMap[containerID] = NewFalsePositive()
	}

	switch metric {
	case "cpu":
		return server.falsePositiveMap[containerID].cpuCount
	case "memory":
		return server.falsePositiveMap[containerID].memCount

	}
	return 0
}

func (server *Server) GetFalsePositiveThreshold(containerID string, metric string) int64 {
	server.RLock()
	defer server.RUnlock()

	switch metric {
	case "cpu":
		return server.cpufp
	case "memory":
		return server.memfp

	}
	return 0
}
