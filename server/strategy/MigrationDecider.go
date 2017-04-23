package strategy

import (
	"github.com/Sirupsen/logrus"
	"github.com/adhuri/Compel-Migration/protocol"
	"github.com/adhuri/Compel-Migration/server/model"
)

// Returns if Migration is Needed
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

			movableStatus := server.CheckIfContainerIsMovable(container.ContainerId)

			newContainer := NewContainer(container.ContainerId, cpuAverage, memAverage, movableStatus)
			newBucket.ContainerDetails = append(newBucket.ContainerDetails, newContainer)
		}
		// Add newBucket to Buckets list
		Buckets = append(Buckets, newBucket)
	}

	PrintAllBuckets(Buckets, log)

	finalDecision, migrationDetails := migrationDecision(Buckets, server, log)

	return finalDecision, migrationDetails
}

// Returns Decision if true or false & Checkpoint Request for Buckets
func migrationDecision(buckets []*Bucket, server *model.Server, log *logrus.Logger) (bool, *protocol.CheckpointRequest) {

	switch len(buckets) {
	case 0:
		log.Warnln("No host running , migration not possible ")
	case 1:
		log.Warnln("Only one host running , migration not possible ")
	default:
		log.Infoln("Migration for Number of hosts : ", len(buckets), "started ....")
	}

	allMetrics := []string{"memory", "cpu"}
	for _, metric := range allMetrics {
		return stagedDecision(metric, buckets, server, log)
	}
	return false, &protocol.CheckpointRequest{}
}

func stagedDecision(metric string, buckets []*Bucket, server *model.Server, log *logrus.Logger) (bool, *protocol.CheckpointRequest) {

	decisionFlag, CheckpointRequest := metricDecision(metric, buckets, server, log)
	if decisionFlag {
		// Increment the counter for the False positive Checker
		server.IncrementFalsePositive(CheckpointRequest.ContainerID, metric)

		if !server.GetMigrationStatus() {
			// If Not Migrating any container
			log.Debugln(metric, "Decision : Server is not migrating - Considering request for migration")
			if !CheckIfFalsePositive(metric, CheckpointRequest.ContainerID, server, log) {
				// Check if set Threshold is crossed to avoid false positives
				log.Debugln(metric, "Decision : Configured Threshold for ", metric, " crossed ,not false positive - Considering request for migration")
				if !CheckIfMigrationTrashing(CheckpointRequest.ContainerID, server, log) {
					// If system is not thrashingDecision
					log.Debugln(metric, "Decision : Migration is not trashing - Considering request for migration")
					return true, CheckpointRequest
				} else {
					log.Debugln(metric, "Decision : Migration is trashing , Last migration done recently - Cannot Migrate")
				}
			} else {
				log.Warnln(metric, " Decision : Configured Threshold for memory not crossed, probably false positive - Cannot Migrate")

			}
		} else {
			log.Errorln("memory Decision : Previous Migration is in progress - Cannot Migrate to avoid CHAOS")
		}
	}
	return false, &protocol.CheckpointRequest{}
}
