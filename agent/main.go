package main

import (
	"encoding/gob"
	"flag"
	"fmt"
	"net"
	"os"
	"os/exec"
	"sync"
	"time"

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

func containerPresent(userName, containerId string) bool {
	_, err := exec.Command("/home/"+userName+"/scripts/FindContainer.sh", "-c", containerId).Output()
	if err != nil {
		return false
	}
	return true
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
	response.StatusMap["Stop Load-Balancer"] = *protocol.NewStatus()
	response.StatusMap["Start Load-Balancer"] = *protocol.NewStatus()
}

func handleMigrationRequest(conn net.Conn, userName string, active bool, agent *Agent) {
	//defer conn.Close()

	migrationRequest := protocol.CheckpointRequest{}

	decoder := gob.NewDecoder(conn)
	err := decoder.Decode(&migrationRequest)
	//err := binary.Read(conn, binary.LittleEndian, &connectMessage)
	if err != nil {
		// If failure in parsing, close the connection and return
		log.Errorln("Bad Migration Request From Server" + err.Error())
		return
	}

	containerId := migrationRequest.ContainerID
	if !containerPresent(userName, containerId) {
		log.Errorln("Container " + containerId + " Not Present.")
		return
	}

	if agent.GetMigrationGoingStatus() {
		log.Errorln("Migration Going On: Cannont Migrate at this time")
		return
	}

	agent.SetMigrationGoingStatus(true)

	// If success, print the message received
	log.Infoln("Migration Request Received")
	log.Debugln("Migration Request Content : ", migrationRequest)
	migrationResponse := protocol.NewCheckpointResponse(migrationRequest)
	InitializeResponse(migrationResponse)

	flag := false

	startTime := time.Now()
	if active {
		flag = true
		fmt.Println("")
		log.Infoln("						MIGRATION STATS")
		checkpointName := migrationRequest.CheckpointName
		destinationIp := migrationRequest.DestinationAgentIP
		CheckpointAndRestore(containerId, destinationIp, checkpointName, userName, migrationResponse)
	}
	elapsed := time.Since(startTime)

	agent.SetMigrationGoingStatus(false)

	if flag {
		migrationResponse.TotalDuration = elapsed
		PrintResponse(migrationResponse)

	} else {
		log.Infoln("Migration is set to OFF")
	}

	encoder := gob.NewEncoder(conn)
	err = encoder.Encode(migrationResponse)
	//err = binary.Write(conn, binary.LittleEndian, connectAck)
	if err != nil {
		// If failure in parsing, close the connection and return
		log.Errorln("Failure in Sending Checkpoint Response")
		return
	}
	log.Infoln("Checkpoint Response Successfully Sent to the Migration Server")
	// close connection when done
	conn.Close()

}

func tcpListener(wg *sync.WaitGroup, userName string, active bool, agent *Agent) {
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
		go handleMigrationRequest(conn, userName, active, agent)
	}
}

func PrintResponse(response *protocol.CheckpointResponse) {

	arr := []string{"Metadata Dump",
		"Metadata Scp",
		"Container Checkpoint",
		"Checkpoint Transfer",
		"Filesystem Export",
		"FileSystem Transfer",
		"Container Restore",
		"Checkpoint Cleanup",
	}

	for _, k := range arr {
		v := response.StatusMap[k]
		d := v.Duration / time.Millisecond
		fmt.Println("        Status : ", v.IsSuccess, "\t\tDuration : "+d.String()+"\t\tActivity : "+k)
	}
	fmt.Println("        OVERALL STATUS :\t", response.IsSuccess)
	fmt.Println("        OVERALL TIME TAKEN :\t", response.TotalDuration.String())

}

func main() {

	//migrationIp := flag.String("server", "127.0.0.1", "ip of the migration server")
	//migrationerverPort := flag.String("udpport", "7071", "udp port on the server")
	userName := flag.String("username", "ssakpal", "username on the server")
	active := flag.Bool("active", false, "username on the server")
	flag.Parse()

	if !*active {
		log.Warnln("MIGRATION IS SWITCHED OFF --(to active set -active=true ) ")
		log.Warnln("MIGRATION IS SWITCHED OFF --(to active set -active=true )")
	}
	agent := NewAgent()

	var wg sync.WaitGroup
	wg.Add(1)

	go tcpListener(&wg, *userName, *active, agent)

	wg.Wait()

}
