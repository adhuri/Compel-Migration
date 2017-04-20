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
