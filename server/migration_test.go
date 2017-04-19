package main

import (
	"encoding/gob"
	"net"
	"os"
	"testing"
	"time"

	"github.com/Sirupsen/logrus"

	"github.com/adhuri/Compel-Migration/protocol"
)

var (
	tlog *logrus.Logger
)

func initLog() {

	tlog = logrus.New()

	// Output logging to stdout
	tlog.Out = os.Stdout

	// Only log the info severity or above.
	tlog.Level = logrus.DebugLevel

	// Microseconds level logging
	customFormatter := new(logrus.TextFormatter)
	customFormatter.TimestampFormat = "2006-01-02 15:04:05.000000"
	customFormatter.FullTimestamp = true

	tlog.Formatter = customFormatter

}

type PredictionData struct {
	Timestamp  int64
	ClientData []ClientInfo
}

type ClientInfo struct {
	AgentIp       string
	ContainerData []ContainerInfo
}

type ContainerInfo struct {
	ContainerId string
	CPU         []float64
	Memory      []float64
}

func dummyData() *protocol.PredictionData {
	cpuPredictions1 := []float64{1.2, 1.3, 1.14, 1.5}
	memoryPredictions1 := []float64{1.2, 1.3, 1.14, 1.5}

	cpuPredictions2 := []float64{2.3, 3.3, 4.14, 4.5}
	memoryPredictions2 := []float64{7.2, 8.3, 9.14, 5.5}

	cpuPredictions3 := []float64{3.2, 4.3, 2.14, 11.5}
	memoryPredictions3 := []float64{4.2, 4.3, 4.14, 9.5}

	// Agent IP 1
	agentIP1 := "10.10.3.1"

	containerInfo1 := protocol.NewContainerInfo("container31", cpuPredictions1, memoryPredictions1)
	containerInfo2 := protocol.NewContainerInfo("container32", cpuPredictions2, memoryPredictions2)

	containerData1 := []protocol.ContainerInfo{*containerInfo1, *containerInfo2}

	clientInfo1 := protocol.NewClientInfo(agentIP1, containerData1)

	// Agent IP 2

	agentIP2 := "10.10.4.1"

	containerInfo3 := protocol.NewContainerInfo("container41", cpuPredictions3, memoryPredictions3)

	containerData2 := []protocol.ContainerInfo{*containerInfo3}

	clientInfo2 := protocol.NewClientInfo(agentIP2, containerData2)

	// Add both Agents to ClientInfo

	clientData1 := []protocol.ClientInfo{*clientInfo1, *clientInfo2}

	predictionData := protocol.NewPredictionData(time.Now().Unix(), clientData1)

	return predictionData
}

func sendPredictionData(conn net.Conn) error {

	predictionData := dummyData()

	//log.Infoln("Sample Prediction data message ", predictionData)

	encoder := gob.NewEncoder(conn)
	err := encoder.Encode(predictionData)
	if err != nil {
		// If error occurs in sending a connect message to server then return
		tlog.Errorln("Failure While Sending Data To Server " + err.Error())
		return err
	}
	tlog.Infoln("Prediction Data Message Successfully Sent")

	// read ack from the server
	serverReply := protocol.PredictionDataResponse{}
	decoder := gob.NewDecoder(conn)
	err = decoder.Decode(&serverReply)
	if err != nil {
		// If error occurs while reading ACK from server then return
		tlog.Errorln("Bad Reply From Server " + err.Error())
		return err

	} else {
		// Print the ACK received from the server
		tlog.Infoln("Prediction Data Message ACK Received")
	}
	// If everything goes well return nil error
	return nil
}

func TestMigrationCode(t *testing.T) {
	initLog()
	addr := "127.0.0.1" + ":" + "5051"
	log.Info("Connect test client to server ", addr)
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		// Before trying to reconnect to the server wait for 3 seconds
		tlog.Warn("Server Not Alive")
	} else {
		// If connection successful send a connect message
		err = sendPredictionData(conn)
		if err != nil {
			// Connect Protocol failed midway; Retry
			tlog.Error("Sending of Data Failed. Try Reconnecting to Migration Server")
			defer conn.Close()
		} else {
			// Client Registration successful
			tlog.Infoln("Connected to Server")
		}
	}

}
