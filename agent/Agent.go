package main

import "sync"

type Agent struct {
	sync.RWMutex
	isMigrationGoingOn bool
}

func NewAgent() *Agent {
	return &Agent{
		isMigrationGoingOn: false,
	}
}

func (agent *Agent) GetMigrationGoingStatus() bool {
	agent.RLock()
	defer agent.RUnlock()

	return agent.isMigrationGoingOn

}

func (agent *Agent) SetMigrationGoingStatus(status bool) {
	agent.Lock()
	defer agent.Unlock()

	agent.isMigrationGoingOn = status

}
