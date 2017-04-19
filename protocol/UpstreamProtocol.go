package protocol

type PredictionData struct {
	Timestamp  int64
	ClientData []ClientInfo
}

type ClientInfo struct {
	AgentIp       string
	ContainerData []ContainerInfo
}

type ContainerInfo struct {
	ContainerId string
	CPU         []float64
	Memory      []float64
}

type PredictionDataResponse struct {
	Timestamp int64
	IsSucess  bool
}

func NewPredictionDataResponse(timestamp int64, status bool) *PredictionDataResponse {
	return &PredictionDataResponse{Timestamp: timestamp,
		IsSucess: status,
	}
}

func NewPredictionData(timestamp int64, clientInfo []ClientInfo) *PredictionData {
	return &PredictionData{
		Timestamp:  timestamp,
		ClientData: clientInfo,
	}
}

func NewClientInfo(agentIp string, containerInfo []ContainerInfo) *ClientInfo {
	return &ClientInfo{
		AgentIp:       agentIp,
		ContainerData: containerInfo,
	}
}

func NewContainerInfo(containerId string, cpuPredictions, memoryPredictions []float64) *ContainerInfo {
	return &ContainerInfo{
		ContainerId: containerId,
		CPU:         cpuPredictions,
		Memory:      memoryPredictions,
	}

}
