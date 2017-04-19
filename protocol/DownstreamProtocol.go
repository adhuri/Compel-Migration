package protocol

type CheckpointRequest struct {
	SourceAgentIP      string
	ContainerID        string
	DestinationAgentIP string
	CheckpointName     string
}

type CheckpointResponse struct {
	Request   CheckpointRequest
	IsSuccess bool
}

func NewCheckpointRequest(sourceIp, containerId, destinationIp, checkpointName string) *CheckpointRequest {
	return &CheckpointRequest{
		SourceAgentIP:      sourceIp,
		ContainerID:        containerId,
		DestinationAgentIP: destinationIp,
		CheckpointName:     checkpointName,
	}
}

func NewCheckpointResponse(request CheckpointRequest, status bool) *CheckpointResponse {
	return &CheckpointResponse{
		Request:   request,
		IsSuccess: status,
	}
}
