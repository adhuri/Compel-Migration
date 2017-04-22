package main

import (
	"fmt"
	"os/exec"
	"time"

	"github.com/adhuri/Compel-Migration/protocol"
)

type CommandResult struct {
	Command   string
	TimeTaken time.Duration
	IsSuccess bool
}

func DumpMetadata(containerId, destinationIp, checkpointName, user string, chan1 chan CommandResult, commonChan chan CommandResult) {
	// Dump Metadata for a contianer
	startTime := time.Now()
	_, err := exec.Command("/home/"+user+"/scripts/DumpMetadata.sh", "-c", containerId, "-u", user, "-d", destinationIp, "-n", checkpointName).Output()
	if err != nil {
		fmt.Println("Dumping Contaier Metadata failed for container " + containerId)
		chan1 <- CommandResult{
			Command:   "Metadata Dump",
			TimeTaken: time.Nanosecond,
			IsSuccess: false,
		}
	}

	d1 := TimeTrack(startTime, "Dumping Container Metadata ")

	chan1 <- CommandResult{
		Command:   "Metadata Dump",
		TimeTaken: d1,
		IsSuccess: true,
	}

	// SCP Metadata to Destination
	startTime = time.Now()
	_, err = exec.Command("/home/"+user+"/scripts/MetadataSCP.sh", "-c", containerId, "-u", user, "-d", destinationIp, "-n", checkpointName).Output()
	if err != nil {
		fmt.Println("SCP Container Metadata failed " + containerId + " to Destination " + destinationIp)
		commonChan <- CommandResult{
			Command:   "Metadata Scp",
			TimeTaken: time.Nanosecond,
			IsSuccess: false,
		}
	}

	d2 := TimeTrack(startTime, "Transfering Container Metadata ")

	commonChan <- CommandResult{
		Command:   "Metadata Scp",
		TimeTaken: d2,
		IsSuccess: true,
	}

}

func ExecuteAndTransferCheckpoint(containerId, destinationIp, checkpointName, user string, chan2 chan CommandResult, commonChan chan CommandResult) {
	// Checkpoint a contianer
	startTime := time.Now()
	_, err := exec.Command("/home/"+user+"/scripts/ExecuteCheckpoint.sh", "-c", containerId, "-u", user, "-d", destinationIp, "-n", checkpointName).Output()
	if err != nil {
		fmt.Println("Checkpointing failed for container " + containerId)
		chan2 <- CommandResult{
			Command:   "Container Checkpoint",
			TimeTaken: time.Nanosecond,
			IsSuccess: false,
		}
	}

	d1 := TimeTrack(startTime, "Checkpointing Container ")

	chan2 <- CommandResult{
		Command:   "Container Checkpoint",
		TimeTaken: d1,
		IsSuccess: true,
	}

	// SCP checkpoint files
	startTime = time.Now()
	_, err = exec.Command("/home/"+user+"/scripts/CheckpointSCP.sh", "-c", containerId, "-u", user, "-d", destinationIp, "-n", checkpointName).Output()
	if err != nil {
		fmt.Println("SCP Checkpoint failed " + containerId + " to Destination " + destinationIp)
		commonChan <- CommandResult{
			Command:   "Checkpoint Transfer",
			TimeTaken: time.Nanosecond,
			IsSuccess: false,
		}
	}

	d2 := TimeTrack(startTime, "Transfering Container Checkpoint Files ")

	commonChan <- CommandResult{
		Command:   "Checkpoint Transfer",
		TimeTaken: d2,
		IsSuccess: true,
	}
}

func ExportAndTransferFileSystem(containerId, destinationIp, checkpointName, user string, commonChan chan CommandResult) {
	// Filesystem export a contianer
	startTime := time.Now()
	_, err := exec.Command("/home/"+user+"/scripts/ExportFilesystem.sh", "-c", containerId, "-u", user, "-d", destinationIp, "-n", checkpointName).Output()
	if err != nil {
		fmt.Println("Filesystem Export Failed for Container " + containerId)
		commonChan <- CommandResult{
			Command:   "Filesystem Export",
			TimeTaken: time.Nanosecond,
			IsSuccess: false,
		}
	}

	d1 := TimeTrack(startTime, "Exporting Container Filesystem ")

	commonChan <- CommandResult{
		Command:   "Filesystem Export",
		TimeTaken: d1,
		IsSuccess: true,
	}

	// SCP Filesystem files
	startTime = time.Now()
	_, err = exec.Command("/home/"+user+"/scripts/SCPFilesystem.sh", "-c", containerId, "-u", user, "-d", destinationIp, "-n", checkpointName).Output()
	if err != nil {
		fmt.Println("Filsystem SCP failed for container " + containerId + " to Destination " + destinationIp)
		commonChan <- CommandResult{
			Command:   "FileSystem Transfer",
			TimeTaken: time.Nanosecond,
			IsSuccess: false,
		}
	}

	d2 := TimeTrack(startTime, "Transfering Container Filesystem ")

	commonChan <- CommandResult{
		Command:   "FileSystem Transfer",
		TimeTaken: d2,
		IsSuccess: true,
	}

}

func RestoreRemoteContainer(containerId, destinationIp, checkpointName, user string, chan3 chan CommandResult) {
	// Remote Container Restoration
	startTime := time.Now()
	_, err := exec.Command("/home/"+user+"/scripts/RestoreRemote.sh", "-c", containerId, "-u", user, "-d", destinationIp, "-n", checkpointName).Output()
	if err != nil {
		fmt.Println("Restoration for container " + containerId + " on Destination " + destinationIp + " Failed.")
		chan3 <- CommandResult{
			Command:   "Container Restore",
			TimeTaken: time.Nanosecond,
			IsSuccess: false,
		}
	}

	d1 := TimeTrack(startTime, "Restoring Remote Container ")

	chan3 <- CommandResult{
		Command:   "Container Restore",
		TimeTaken: d1,
		IsSuccess: true,
	}
}

func CheckpointCleanup(containerId, destinationIp, checkpointName, user string) (CommandResult, error) {

	// Checkpoint Cleanup
	startTime := time.Now()
	_, err := exec.Command("/home/"+user+"/scripts/CheckpointCleanup.sh", "-c", containerId, "-u", user, "-d", destinationIp, "-n", checkpointName).Output()
	if err != nil {
		fmt.Println("Checkpoint Cleanup Failed")
		return CommandResult{
			Command:   "Checkpoint Cleanup",
			TimeTaken: time.Nanosecond,
			IsSuccess: false,
		}, err
	}
	d1 := TimeTrack(startTime, "Checkpoint Cleanup")
	return CommandResult{
		Command:   "Checkpoint Cleanup",
		TimeTaken: d1,
		IsSuccess: true,
	}, nil

}

func TimeTrack(start time.Time, name string) time.Duration {
	elapsed := time.Since(start)
	fmt.Println("        ", name, " : ", elapsed)
	return elapsed
}

func CheckpointAndRestore(containerId, destinationIp, checkpointName, user string, response *protocol.CheckpointResponse) {
	chan1 := make(chan CommandResult)
	chan2 := make(chan CommandResult)
	chan3 := make(chan CommandResult)
	commonChan := make(chan CommandResult, 4)

	// Start Metadata Dump and Transfer
	go DumpMetadata(containerId, destinationIp, checkpointName, user, chan1, commonChan)

	result := <-chan1
	response.StatusMap[result.Command] = protocol.Status{Duration: result.TimeTaken, IsSuccess: result.IsSuccess}
	if !result.IsSuccess {
		fmt.Println("Metadata Dump Failed")
		close(commonChan)
		close(chan2)
		close(chan3)
		return
	}

	// Start Container Checkpoint and Transfer
	go ExecuteAndTransferCheckpoint(containerId, destinationIp, checkpointName, user, chan2, commonChan)

	// Waiting for checkpointing to complete
	result = <-chan2
	response.StatusMap[result.Command] = protocol.Status{Duration: result.TimeTaken, IsSuccess: result.IsSuccess}
	if !result.IsSuccess {
		fmt.Println("Checkpoint Failed")
		close(commonChan)
		close(chan3)
		return
	}

	// Start FileSystem Export and Transfer
	go ExportAndTransferFileSystem(containerId, destinationIp, checkpointName, user, commonChan)

	// Waiting for all the checkpoint steps to complete
	completeStatus := true
	for i := 0; i < 4; i++ {
		result = <-commonChan
		response.StatusMap[result.Command] = protocol.Status{Duration: result.TimeTaken, IsSuccess: result.IsSuccess}
		completeStatus = completeStatus && result.IsSuccess
	}

	if !completeStatus {
		fmt.Println("One of the SCP Failed")
		close(chan3)
		return
	}

	// Start container restore
	go RestoreRemoteContainer(containerId, destinationIp, checkpointName, user, chan3)

	// wait for restore to complete
	result = <-chan3
	response.StatusMap[result.Command] = protocol.Status{Duration: result.TimeTaken, IsSuccess: result.IsSuccess}
	if !result.IsSuccess {
		fmt.Println("Restoration Failed")
		return
	}

	// Start cleanup
	result, err := CheckpointCleanup(containerId, destinationIp, checkpointName, user)
	if err != nil {
		return
	}
	response.StatusMap[result.Command] = protocol.Status{Duration: result.TimeTaken, IsSuccess: result.IsSuccess}

}
