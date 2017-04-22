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
	isSuccess bool
}

type Status struct {
	IsSuccess bool
	Duration  time.Duration
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
		isSuccess: false,
	}
}

func NewStatus() *Status {

	return &Status{
		IsSuccess: false,
		Duration:  time.Nanosecond,
	}
}
