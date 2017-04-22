package protocol

import "time"

type CheckpointRequest struct {
	SourceAgentIP      string
	ContainerID        string
	DestinationAgentIP string
	CheckpointName     string
}

type CheckpointResponse struct {
	Request   CheckpointRequest
	StatusMap map[string]Status
}

type Status struct {
	isSucess bool
	duration time.Duration
}

func NewCheckpointRequest(sourceIp, containerId, destinationIp, checkpointName string) *CheckpointRequest {
	return &CheckpointRequest{
		SourceAgentIP:      sourceIp,
		ContainerID:        containerId,
		DestinationAgentIP: destinationIp,
		CheckpointName:     checkpointName,
	}
}

func NewCheckpointResponse(request CheckpointRequest) *CheckpointResponse {
	return &CheckpointResponse{
		Request:   request,
		StatusMap: make(map[string]Status),
	}
}

func NewStatus() *Status {

	return &Status{
		isSucess: false,
		duration: time.Nanosecond,
	}
}
