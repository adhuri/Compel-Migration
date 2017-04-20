package strategy

import (
	"strconv"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/adhuri/Compel-Migration/protocol"
	"github.com/adhuri/Compel-Migration/server/model"
)

func CheckIfFalsePositive(containerID string, server *model.Server) bool {
	// Fetch Counter  from server object
	return false
}

func CheckIfMigrationTrashing(containerID string, server *model.Server, log *logrus.Logger) bool {
	// Fetch Timestamp  from server object
	//Unix Time stamps are number of seconds

	thresholdThrashing := server.GetThrashingThreshold() //5 minutes interval between thrashing
	secondsElapsed := time.Now().Unix() - server.GetPreviousContainerMigrationTime(containerID)
	if secondsElapsed < thresholdThrashing {
		log.Warnln("Thrashing might occur for container ", containerID, " : Migration for this container was done - ", (float32(secondsElapsed / 60.0)), " minutes ago")
		return true
	}

	return false
}

func metricDecision(metric string, buckets []*Bucket, server *model.Server, log *logrus.Logger) (bool, *protocol.CheckpointRequest) {
	// We handle Two main cases when
	// a ) some agents have -ve free memory due to prediction
	// b )all agents have  + ve free due to predicition
	unixTimestamp := time.Now().Unix()
	timestamp := strconv.FormatInt(unixTimestamp, 10)
	sortBucketsAsc(buckets, metric) // In place sort

	//For all values in buckets i to j
	for bucketIndexi, bucketi := range buckets {
		// Case a

		if bucketi.GetValue(metric) < 0 {
			// For all positive values in buckets k to j where i<k<j
			for _, bucketk := range buckets {

				if bucketk.GetValue(metric) >= 0 {
					// Sort Containers descending on i
					sortContainersDesc(bucketi.ContainerDetails, metric)

					// For each container on i
					for _, containeri := range bucketi.ContainerDetails {
						// Check if container static or if we cannot move
						if !containeri.movableContainer || (bucketk.GetValue(metric)-containeri.GetValue(metric)) < 0 {
							// Cannot Move this container container i
							log.Debugln(metric, " Decision : Immovable Container ", containeri.ContainerID, " cannot be moved since it is configured immovable")

							continue
						} else {
							// Move containiner Ci to agent k
							return true, &protocol.CheckpointRequest{
								SourceAgentIP:      bucketi.AgentIP,
								ContainerID:        containeri.ContainerID,
								CheckpointName:     timestamp,
								DestinationAgentIP: bucketk.AgentIP,
							}
						}
					} // end of containers

				} // end of positive checker

			} // end of k
			// No potential Candidate found ( all could be -ves)
			return false, &protocol.CheckpointRequest{}
		}

		// Case b
		bucketIndexj := len(buckets) - 1 // Switch at end
		// Till i and j are not same
		for bucketIndexi != bucketIndexj {
			// sort in place
			sortContainersDesc(buckets[bucketIndexj].ContainerDetails, metric)
			// Sorted  containers in bucket j
			for _, containerj := range buckets[bucketIndexj].ContainerDetails {
				if !containerj.movableContainer || (bucketi.GetValue(metric)-containerj.GetValue(metric)) < 0 {
					// Container j cannot be moved
					log.Debugln(metric, " Decision : Immovable Container ", containerj.ContainerID, " cannot be moved since it is configured immovable")
					continue
				} else {
					// Move containiner Cj to agent i
					return true, &protocol.CheckpointRequest{
						SourceAgentIP:      buckets[bucketIndexj].AgentIP,
						ContainerID:        containerj.ContainerID,
						CheckpointName:     timestamp,
						DestinationAgentIP: bucketi.AgentIP,
					}

				}
				// If no potential candidate in j move to previous container

			}
			bucketIndexj--

		}
		// No potential Candidate found
		return false, &protocol.CheckpointRequest{}

	}

	return true, &protocol.CheckpointRequest{}
}
