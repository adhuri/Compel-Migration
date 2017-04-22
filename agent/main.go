package main

import (
	"encoding/gob"
	"flag"
	"net"
	"os"
	"sync"

	"github.com/Sirupsen/logrus"
	"github.com/adhuri/Compel-Migration/protocol"
)

var (
	log *logrus.Logger
)

func init() {

	log = logrus.New()

	// Output logging to stdout
	log.Out = os.Stdout

	// Only log the info severity or above.
	log.Level = logrus.InfoLevel

	// Microseconds level logging
	customFormatter := new(logrus.TextFormatter)
	customFormatter.TimestampFormat = "2006-01-02 15:04:05.000000"
	customFormatter.FullTimestamp = true

	log.Formatter = customFormatter

}

func InitializeResponse(response *protocol.CheckpointResponse) {
	response.StatusMap["Metadata Dump"] = *protocol.NewStatus()
	response.StatusMap["Metadata Scp"] = *protocol.NewStatus()
	response.StatusMap["Container Checkpoint"] = *protocol.NewStatus()
	response.StatusMap["Checkpoint Transfer"] = *protocol.NewStatus()
	response.StatusMap["Filesystem Export"] = *protocol.NewStatus()
	response.StatusMap["FileSystem Transfer"] = *protocol.NewStatus()
	response.StatusMap["Container Restore"] = *protocol.NewStatus()
	response.StatusMap["Checkpoint Cleanup"] = *protocol.NewStatus()
}

func handleMigrationRequest(conn net.Conn, userName string) {
	migrationRequest := protocol.CheckpointRequest{}
	decoder := gob.NewDecoder(conn)
	err := decoder.Decode(&migrationRequest)
	//err := binary.Read(conn, binary.LittleEndian, &connectMessage)
	if err != nil {
		// If failure in parsing, close the connection and return
		log.Errorln("Bad Migration Request From Server" + err.Error())
		return
	}
	// If success, print the message received
	log.Infoln("Migration Request Received")
	log.Debugln("Migration Request Content : ", migrationRequest)
	migrationResponse := protocol.NewCheckpointResponse(migrationRequest)
	InitializeResponse(migrationResponse)

	containerId := migrationRequest.ContainerID
	checkpointName := migrationRequest.CheckpointName
	destinationIp := migrationRequest.DestinationAgentIP
	CheckpointAndRestore(containerId, destinationIp, checkpointName, userName, migrationResponse)
	// command := "./checkpoint.sh -c " + containerName + " -u " + userName + " -n " + checkpointName + " -d " + hostName
	// cmd := exec.Command("/bin/sh", "-c", command)

	// Create a ConnectAck Message

	// Send Connect Ack back to the client
	encoder := gob.NewEncoder(conn)
	err = encoder.Encode(migrationResponse)
	//err = binary.Write(conn, binary.LittleEndian, connectAck)
	if err != nil {
		// If failure in parsing, close the connection and return
		log.Errorln("Failure in Sending Migration Success Ack")
		return
	}
	log.Infoln("Migration Was Success")
	// close connection when done
	conn.Close()

}

func tcpListener(wg *sync.WaitGroup, userName string) {
	defer wg.Done()
	// Server listens on all interfaces for TCP connestion
	addr := ":" + "5052"
	log.Infoln("Migration Agent listening on TCP ", addr)
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalln("Agent Failed To Start ")
	}

	// Wait for clients to connect
	for {
		// Accept a connection and spin-off a goroutine
		conn, err := listener.Accept()
		if err != nil {
			// If error continue to wait for other clients to connect
			continue
		}
		log.Infoln(" Accepted Connection from Prediction Client ")
		go handleMigrationRequest(conn, userName)
	}
}

func main() {

	//migrationIp := flag.String("server", "127.0.0.1", "ip of the migration server")
	//migrationerverPort := flag.String("udpport", "7071", "udp port on the server")
	userName := flag.String("username", "ssakpal", "username on the server")
	flag.Parse()

	var wg sync.WaitGroup
	wg.Add(1)

	go tcpListener(&wg, *userName)

	wg.Wait()

}
