package strategy

import (
	"github.com/Sirupsen/logrus"
	"github.com/adhuri/Compel-Migration/protocol"
	"github.com/adhuri/Compel-Migration/server/model"
)

func metricDecision(metric string, buckets []*Bucket, server *model.Server, log *logrus.Logger) (bool, *protocol.CheckpointRequest) {

	return false, &protocol.CheckpointRequest{}
}

func CheckIfFalsePositive(containerID string, server *model.Server) bool {
	// Fetch Counter  from server object
	return false
}

func CheckIfMigrationTrashing(containerID string, server *model.Server) bool {
	// Fetch Timestamp  from server object
	return false
}
