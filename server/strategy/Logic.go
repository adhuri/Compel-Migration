package strategy

import (
	"github.com/Sirupsen/logrus"
	"github.com/adhuri/Compel-Migration/protocol"
	"github.com/adhuri/Compel-Migration/server/model"
)

func MigrationNeeded(predictionData *protocol.PredictionData, server *model.Server, log *logrus.Logger) (bool, *protocol.CheckpointRequest) {
	// complete this function with the migration logic

	log.Debugln("Checking if Migration Needed")
	Buckets := []*Bucket{} // Buckets data structure
	log.Debugln("Buckets Data structure ", Buckets)

	for clientIndex, client := range predictionData.ClientData {
		log.Debug("Creating Bucket for Agent ", clientIndex, " , AgentIP ", client.AgentIp)
		newBucket := NewBucket(client.AgentIp)

		for containerIndex, container := range client.ContainerData {
			log.Debug("Creating Container for Container ", containerIndex, " , ContainerID ", container.ContainerId)
			// Calculate Value - using max or average strategy
			cpuAverage := calculateValue(container.CPU)
			memAverage := calculateValue(container.Memory)

			newContainer := NewContainer(container.ContainerId, cpuAverage, memAverage)
			newBucket.ContainerDetails = append(newBucket.ContainerDetails, newContainer)
		}
		// Add newBucket to Buckets list
		Buckets = append(Buckets, newBucket)
	}

	PrintAllBuckets(Buckets, log)
	return true, &protocol.CheckpointRequest{}
}
