package main

import (
	"fmt"
	"os/exec"
)

func DumpMetadata(containerId, destinationIp, checkpointName, user string) error {
	// Dump Metadata for a contianer
	_, err := exec.Command("/home/"+user+"/DumpMetadata.sh", "-c", containerId, "-u", user, "-d", destinationIp, "-n", containerId).Output()
	if err != nil {
		fmt.Printf("Checkpointing failed for container " + containerId)
		return err
	}

	// SCP Metadata to Destination
	_, err = exec.Command("/home/"+user+"/MetadataSCP.sh", "-c", containerId, "-u", user, "-d", destinationIp, "-n", containerId).Output()
	if err != nil {
		fmt.Printf("SCP Checkpoint failed " + containerId + " to Destination " + destinationIp)
		return err
	}

	return nil

}

func ExecuteTransferCheckpoint(containerId, destinationIp, checkpointName, user string) error {
	// Checkpoint a contianer
	_, err := exec.Command("/home/"+user+"/ExecuteCheckpoint.sh", "-c", containerId, "-u", user, "-d", destinationIp, "-n", containerId).Output()
	if err != nil {
		fmt.Printf("Checkpointing failed for container " + containerId)
		return err
	}

	// SCP checkpoint files
	_, err = exec.Command("/home/"+user+"/CheckpointSCP.sh", "-c", containerId, "-u", user, "-d", destinationIp, "-n", containerId).Output()
	if err != nil {
		fmt.Printf("SCP Checkpoint failed " + containerId + " to Destination " + destinationIp)
		return err
	}

	return nil
}

func ExportTransferFileSystem(containerId, destinationIp, checkpointName, user string) error {
	// Filesystem export a contianer
	_, err := exec.Command("/home/"+user+"/ExportFilesystem.sh", "-c", containerId, "-u", user, "-d", destinationIp, "-n", containerId).Output()
	if err != nil {
		fmt.Printf("Filesystem Export Failed for Container " + containerId)
		return err
	}

	// SCP Filesystem files
	_, err = exec.Command("/home/"+user+"/SCPFilesystem.sh", "-c", containerId, "-u", user, "-d", destinationIp, "-n", containerId).Output()
	if err != nil {
		fmt.Printf("Filsystem SCP failed for container " + containerId + " to Destination " + destinationIp)
		return err
	}

	return nil

}

func RestoreContainer() {

}

func CheckpointCleanup() {

}

func Checkpoint(containerId, destinationIp, checkpointName, user string) {

}
