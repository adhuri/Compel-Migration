package main

import (
	"fmt"
	"os/exec"
	"time"
)

func DumpMetadata(containerId, destinationIp, checkpointName, user string) error {
	// Dump Metadata for a contianer
	startTime := time.Now()
	_, err := exec.Command("/home/"+user+"/DumpMetadata.sh", "-c", containerId, "-u", user, "-d", destinationIp, "-n", containerId).Output()
	if err != nil {
		fmt.Printf("Dumping Contaier Metadata failed for container " + containerId)
		return err
	}
	d1 := TimeTrack(startTime, "Dumping Container Metadata ")

	// SCP Metadata to Destination
	startTime = time.Now()
	_, err = exec.Command("/home/"+user+"/MetadataSCP.sh", "-c", containerId, "-u", user, "-d", destinationIp, "-n", containerId).Output()
	if err != nil {
		fmt.Printf("SCP Container Metadata failed " + containerId + " to Destination " + destinationIp)
		return err
	}
	d2 := TimeTrack(startTime, "Transfering Container Metadata ")
	return nil

}

func ExecuteAndTransferCheckpoint(containerId, destinationIp, checkpointName, user string) error {
	// Checkpoint a contianer
	startTime := time.Now()
	_, err := exec.Command("/home/"+user+"/ExecuteCheckpoint.sh", "-c", containerId, "-u", user, "-d", destinationIp, "-n", containerId).Output()
	if err != nil {
		fmt.Printf("Checkpointing failed for container " + containerId)
		return err
	}
	d1 := TimeTrack(startTime, "Checkpointing Container ")

	// SCP checkpoint files
	startTime = time.Now()
	_, err = exec.Command("/home/"+user+"/CheckpointSCP.sh", "-c", containerId, "-u", user, "-d", destinationIp, "-n", containerId).Output()
	if err != nil {
		fmt.Printf("SCP Checkpoint failed " + containerId + " to Destination " + destinationIp)
		return err
	}
	d2 := TimeTrack(startTime, "Transfering Container Checkpoint Files ")
	return nil
}

func ExportAndTransferFileSystem(containerId, destinationIp, checkpointName, user string) error {
	// Filesystem export a contianer
	startTime := time.Now()
	_, err := exec.Command("/home/"+user+"/ExportFilesystem.sh", "-c", containerId, "-u", user, "-d", destinationIp, "-n", containerId).Output()
	if err != nil {
		fmt.Printf("Filesystem Export Failed for Container " + containerId)
		return err
	}
	d1 := TimeTrack(startTime, "Exporting Container Filesystem ")

	// SCP Filesystem files
	startTime = time.Now()
	_, err = exec.Command("/home/"+user+"/SCPFilesystem.sh", "-c", containerId, "-u", user, "-d", destinationIp, "-n", containerId).Output()
	if err != nil {
		fmt.Printf("Filsystem SCP failed for container " + containerId + " to Destination " + destinationIp)
		return err
	}
	d2 := TimeTrack(startTime, "Transfering Container Filesystem ")
	return nil

}

func RestoreRemoteContainer(containerId, destinationIp, checkpointName, user string) error {
	// Remote Container Restoration
	startTime := time.Now()
	_, err := exec.Command("/home/"+user+"/RestoreRemote.sh", "-c", containerId, "-u", user, "-d", destinationIp, "-n", containerId).Output()
	if err != nil {
		fmt.Printf("Restoration for container " + containerId + " on Destination " + destinationIp + " Failed.")
		return err
	}
	d1 := TimeTrack(startTime, "Restoring Remote Container ")
	return nil
}

func CheckpointCleanup(containerId, destinationIp, checkpointName, user string) error {

	// Checkpoint Cleanup
	_, err := exec.Command("/home/"+user+"/CheckpointCleanup.sh", "-c", containerId, "-u", user, "-d", destinationIp, "-n", containerId).Output()
	if err != nil {
		fmt.Printf("Checkpoint Cleanup Failed")
		return err
	}
	return nil

}

func TimeTrack(start time.Time, name string) time.Duration {
	elapsed := time.Since(start)
	fmt.Printf(name, " took ", elapsed, "\n")
	return elapsed
}

func CheckpointAndRestore(containerId, destinationIp, checkpointName, user string) {

	go DumpMetadata(containerId, destinationIp, checkpointName, user)
	go ExecuteAndTransferCheckpoint(containerId, destinationIp, checkpointName, user)
	go ExportAndTransferFileSystem(containerId, destinationIp, checkpointName, user)
	go RestoreRemoteContainer(containerId, destinationIp, checkpointName, user)
	CheckpointCleanup(containerId, destinationIp, checkpointName, user)

}
