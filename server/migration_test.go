package main

import (
	"encoding/gob"
	"net"
	"os"
	"testing"

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

func TestMigrationCode(t *testing.T) {
	initLog()

	predictionData := dummyData()
	tlog.Infoln("For General Scenario")
	simulateTCPConnection(predictionData)
}

func TestScenario1(t *testing.T) {
	initLog()
	//Scenario 1 - Two containers on same machine
	tlog.Infoln(" NOTICE============> Make sure you have mysql container as immovable while starting server")
	tlog.Infoln("Run - Compel-Migration-server  -migrate=true -immovable mysql")
	initLog()
	predictionData := twoHostsScenario1(t)
	tlog.Infoln("Scenario 1 - Two containers on same machine")
	simulateTCPConnection(predictionData)
	tlog.Warnln("No Migration Should happen")
}

func TestScenario2(t *testing.T) {
	initLog()
	//Scenario 2 - Two containers on different machine
	tlog.Infoln(" ============= NOTICE============> Make sure you have mysql container as immovable while starting server")
	tlog.Infoln("Run - Compel-Migration-server  -migrate=true -immovable mysql")
	initLog()
	predictionData := twoHostsScenario2(t)
	tlog.Infoln("Scenario 2 - Two containers on different machine")
	simulateTCPConnection(predictionData)
	tlog.Warnln("rubis should migrate to Machine 1")
}

func TestScenario3(t *testing.T) {
	initLog()
	//Scenario 1 - Two containers on different machine
	tlog.Infoln(" ============= NOTICE============> Make sure you have mysql container as immovable while starting server")
	tlog.Infoln("Run - Compel-Migration-server  -migrate=true -immovable mysql")
	initLog()
	predictionData := twoHostsScenario3(t)
	tlog.Infoln("Scenario 3 - Two containers on different machine ")
	simulateTCPConnection(predictionData)
	tlog.Warnln("rubis should migrate to Machine 2")
}

func TestScenario4(t *testing.T) {
	initLog()
	//Scenario 1 - Two containers on different machine
	tlog.Infoln(" ============= NOTICE============> Make sure you have mysql container as immovable while starting server")
	tlog.Infoln("Run - Compel-Migration-server  -migrate=true -immovable mysql")
	initLog()
	predictionData := twoHostsScenario4(t)
	tlog.Infoln("Scenario 4 - Two containers on same machine - overload")
	simulateTCPConnection(predictionData)
	tlog.Warnln("rubis should migrate to Machine 2")
}

func TestScenario5(t *testing.T) {
	initLog()
	//Scenario 1 - Two containers on different machine
	tlog.Infoln(" ============= NOTICE============> Make sure you have mysql container as immovable while starting server")
	tlog.Infoln("Run - Compel-Migration-server  -migrate=true -immovable mysql")
	initLog()
	predictionData := twoHostsScenario5(t)
	tlog.Infoln("Scenario 5 - Two containers on different machine - more resources rubis consumed on machine2")
	simulateTCPConnection(predictionData)
	tlog.Warnln("No migration")
}

func sendPredictionData(conn net.Conn, data *protocol.PredictionData) error {

	encoder := gob.NewEncoder(conn)
	err := encoder.Encode(data)
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

func simulateTCPConnection(data *protocol.PredictionData) {
	addr := "127.0.0.1" + ":" + "5051"
	log.Info("Connect test client to server ", addr)
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		// Before trying to reconnect to the server wait for 3 seconds
		tlog.Warn("Server Not Alive")
	} else {
		// If connection successful send a connect message
		err = sendPredictionData(conn, data)
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
