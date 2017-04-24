package main

import (
	"testing"
	"time"

	"github.com/adhuri/Compel-Migration/protocol"
)

func dummyData() *protocol.PredictionData {
	cpuPredictions1 := []float32{1.2, 1.3, 1.14, 1.5}
	memoryPredictions1 := []float32{1.2, 1.3, 1.14, 1.5}

	cpuPredictions2 := []float32{2.3, 3.3, 4.14, 4.5}
	memoryPredictions2 := []float32{7.2, 8.3, 9.14, 5.5}

	cpuPredictions3 := []float32{3.2, 4.3, 2.14, 11.5}
	memoryPredictions3 := []float32{4.2, 4.3, 4.14, 9.5}

	// Agent IP 1
	agentIP1 := "127.0.0.2"

	containerInfo1 := protocol.NewContainerInfo("container31", cpuPredictions1, memoryPredictions1)
	containerInfo2 := protocol.NewContainerInfo("container32", cpuPredictions2, memoryPredictions2)

	containerData1 := []protocol.ContainerInfo{*containerInfo1, *containerInfo2}

	clientInfo1 := protocol.NewClientInfo(agentIP1, containerData1)

	// Agent IP 2

	agentIP2 := "127.0.0.1"

	containerInfo3 := protocol.NewContainerInfo("container41", cpuPredictions3, memoryPredictions3)

	containerData2 := []protocol.ContainerInfo{*containerInfo3}

	clientInfo2 := protocol.NewClientInfo(agentIP2, containerData2)

	// Add both Agents to ClientInfo

	clientData1 := []protocol.ClientInfo{*clientInfo1, *clientInfo2}

	predictionData := protocol.NewPredictionData(time.Now().Unix(), clientData1)

	return predictionData
}

func twoHostsScenario1(t *testing.T) *protocol.PredictionData {

	//Scenario 1 - Two containers on same machine
	cpuPredictions1 := []float32{40.0, 40.0}
	memoryPredictions1 := []float32{40.0, 40.0}

	cpuPredictions2 := []float32{50.0, 50.0}
	memoryPredictions2 := []float32{50.0, 50.0}

	// Agent IP 1
	agentIP1 := "10.10.3.1"

	containerInfo1 := protocol.NewContainerInfo("mysql", cpuPredictions1, memoryPredictions1)
	containerInfo2 := protocol.NewContainerInfo("rubis", cpuPredictions2, memoryPredictions2)

	containerData1 := []protocol.ContainerInfo{*containerInfo1, *containerInfo2}

	clientInfo1 := protocol.NewClientInfo(agentIP1, containerData1)

	// Agent IP 2

	agentIP2 := "10.10.4.1"
	containerData2 := []protocol.ContainerInfo{}
	clientInfo2 := protocol.NewClientInfo(agentIP2, containerData2)

	// Add both Agents to ClientInfo

	clientData1 := []protocol.ClientInfo{*clientInfo1, *clientInfo2}

	predictionData := protocol.NewPredictionData(time.Now().Unix(), clientData1)

	return predictionData
}

func twoHostsScenario2(t *testing.T) *protocol.PredictionData {

	//Scenario 2
	cpuPredictions1 := []float32{50.0, 50.0}
	memoryPredictions1 := []float32{50.0, 50.0}

	cpuPredictions2 := []float32{40.0, 40.0}
	memoryPredictions2 := []float32{40.0, 40.0}

	// Agent IP 1
	agentIP1 := "10.10.3.1"

	containerInfo1 := protocol.NewContainerInfo("mysql", cpuPredictions1, memoryPredictions1)
	containerInfo2 := protocol.NewContainerInfo("rubis", cpuPredictions2, memoryPredictions2)

	containerData1 := []protocol.ContainerInfo{*containerInfo1}

	clientInfo1 := protocol.NewClientInfo(agentIP1, containerData1)

	// Agent IP 2

	agentIP2 := "10.10.4.1"
	containerData2 := []protocol.ContainerInfo{*containerInfo2}
	clientInfo2 := protocol.NewClientInfo(agentIP2, containerData2)

	// Add both Agents to ClientInfo

	clientData1 := []protocol.ClientInfo{*clientInfo1, *clientInfo2}

	predictionData := protocol.NewPredictionData(time.Now().Unix(), clientData1)

	return predictionData
}

func twoHostsScenario3(t *testing.T) *protocol.PredictionData {

	//Scenario 3 mysql on 2nd machine , rubis on  1st
	cpuPredictions1 := []float32{40.0, 40.0}
	memoryPredictions1 := []float32{40.0, 40.0}

	cpuPredictions2 := []float32{50.0, 50.0}
	memoryPredictions2 := []float32{50.0, 50.0}

	// Agent IP 1
	agentIP1 := "10.10.3.1"

	containerInfo1 := protocol.NewContainerInfo("rubis", cpuPredictions1, memoryPredictions1)
	containerData1 := []protocol.ContainerInfo{*containerInfo1}
	clientInfo1 := protocol.NewClientInfo(agentIP1, containerData1)

	// Agent IP 2

	agentIP2 := "10.10.4.1"

	containerInfo2 := protocol.NewContainerInfo("mysql", cpuPredictions2, memoryPredictions2)
	containerData2 := []protocol.ContainerInfo{*containerInfo2}
	clientInfo2 := protocol.NewClientInfo(agentIP2, containerData2)

	// Add both Agents to ClientInfo

	clientData1 := []protocol.ClientInfo{*clientInfo1, *clientInfo2}

	predictionData := protocol.NewPredictionData(time.Now().Unix(), clientData1)

	return predictionData
}

func twoHostsScenario4(t *testing.T) *protocol.PredictionData {

	//Scenario 4 both on same machine 1, overload

	cpuPredictions1 := []float32{70.0, 70.0}
	memoryPredictions1 := []float32{70.0, 70.0}

	cpuPredictions2 := []float32{60.0, 60.0}
	memoryPredictions2 := []float32{60.0, 60.0}

	// Agent IP 1
	agentIP1 := "10.10.3.1"

	containerInfo1 := protocol.NewContainerInfo("rubis", cpuPredictions1, memoryPredictions1)
	containerInfo2 := protocol.NewContainerInfo("mysql", cpuPredictions2, memoryPredictions2)

	containerData1 := []protocol.ContainerInfo{*containerInfo1, *containerInfo2}

	clientInfo1 := protocol.NewClientInfo(agentIP1, containerData1)

	// Agent IP 2

	agentIP2 := "10.10.4.1"

	containerData2 := []protocol.ContainerInfo{}
	clientInfo2 := protocol.NewClientInfo(agentIP2, containerData2)

	// Add both Agents to ClientInfo

	clientData1 := []protocol.ClientInfo{*clientInfo1, *clientInfo2}

	predictionData := protocol.NewPredictionData(time.Now().Unix(), clientData1)

	return predictionData
}

func twoHostsScenario5(t *testing.T) *protocol.PredictionData {

	//Scenario 5 Edge case No migration

	cpuPredictions1 := []float32{70.0, 70.0}
	memoryPredictions1 := []float32{70.0, 70.0}

	cpuPredictions2 := []float32{40.0, 40.0}
	memoryPredictions2 := []float32{40.0, 40.0}

	// Agent IP 1
	agentIP1 := "10.10.3.1"

	containerInfo1 := protocol.NewContainerInfo("rubis", cpuPredictions1, memoryPredictions1)

	containerData1 := []protocol.ContainerInfo{*containerInfo1}

	clientInfo1 := protocol.NewClientInfo(agentIP1, containerData1)

	// Agent IP 2

	agentIP2 := "10.10.4.1"
	containerInfo2 := protocol.NewContainerInfo("mysql", cpuPredictions2, memoryPredictions2)
	containerData2 := []protocol.ContainerInfo{*containerInfo2}
	clientInfo2 := protocol.NewClientInfo(agentIP2, containerData2)

	// Add both Agents to ClientInfo

	clientData1 := []protocol.ClientInfo{*clientInfo1, *clientInfo2}

	predictionData := protocol.NewPredictionData(time.Now().Unix(), clientData1)

	return predictionData
}
