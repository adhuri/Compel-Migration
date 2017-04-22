package main

import (
	"encoding/gob"
	"fmt"
	"net"
	"testing"

	"github.com/adhuri/Compel-Migration/protocol"
)

func TestMigrationCode(t *testing.T) {
	request := protocol.CheckpointRequest{
		SourceAgentIP:      "172.31.18.122",
		ContainerID:        "e841db3b15c5",
		DestinationAgentIP: "172.31.26.199",
		CheckpointName:     "third",
	}
	simulateTCPConnection(&request)
}

func sendCheckpoinRequest(conn net.Conn, request *protocol.CheckpointRequest) error {

	encoder := gob.NewEncoder(conn)
	err := encoder.Encode(request)
	if err != nil {
		// If error occurs in sending a connect message to server then return
		fmt.Println("Failure While Sending Data To Server " + err.Error())
		return err
	}
	fmt.Println("Checkpoint Request Successfully Sent")

	// read ack from the server
	serverReply := protocol.CheckpointResponse{}
	decoder := gob.NewDecoder(conn)
	err = decoder.Decode(&serverReply)
	if err != nil {
		// If error occurs while reading ACK from server then return
		fmt.Println("Bad Reply From Server " + err.Error())
		return err

	} else {
		// Print the ACK received from the server
		fmt.Println("Checkpoint Response Message Received")
	}

	// If everything goes well return nil error
	return nil
}

func simulateTCPConnection(request *protocol.CheckpointRequest) {
	addr := "52.35.219.45" + ":" + "5052"
	log.Info("Connect test client to server ", addr)
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		// Before trying to reconnect to the server wait for 3 seconds
		fmt.Println("Server Not Alive")
	} else {
		// If connection successful send a connect message
		err = sendCheckpoinRequest(conn, request)
		if err != nil {
			// Connect Protocol failed midway; Retry
			fmt.Println("Sending of Data Failed. Try Reconnecting to Migration Server")
			defer conn.Close()
		} else {
			// Client Registration successful
			fmt.Println("Connected to Server")
		}
	}
}
