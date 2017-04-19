package main

import (
	"encoding/gob"
	"net"
	"os"
	"sync"

	"github.com/Sirupsen/logrus"

	migrationProtocol "github.com/adhuri/Compel-Migration/protocol"
	"github.com/adhuri/Compel-Monitoring/compel-monitoring-server/model"
)

// on receiving tcp packet
// process it
// if decision to migrate send a migration request

var (
	log    *logrus.Logger
	server *model.Server
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

func handlePredictionDataMessage(conn net.Conn) {
	// Read the ConnectRequest
	predictionDataMessage := migrationProtocol.PredictionData{}
	decoder := gob.NewDecoder(conn)
	err := decoder.Decode(&predictionDataMessage)
	//err := binary.Read(conn, binary.LittleEndian, &connectMessage)
	if err != nil {
		// If failure in parsing, close the connection and return
		log.Errorln("Bad Prediction Data Message From Client" + err.Error())
		return
	} else {
		// If success, print the message received
		log.Infoln("Prediction Data Received")
		log.Debugln("Prediction Data Content : ", predictionDataMessage)
	}

	// Create a ConnectAck Message
	predictionAck := migrationProtocol.NewPredictionDataResponse(predictionDataMessage.Timestamp, true)

	// Send Connect Ack back to the client
	encoder := gob.NewEncoder(conn)
	err = encoder.Encode(predictionAck)
	//err = binary.Write(conn, binary.LittleEndian, connectAck)
	if err != nil {
		// If failure in parsing, close the connection and return
		log.Errorln("Prediction Data Ack")
		return
	}
	log.Infoln("Prediction Data Ack Sent for Request Id " + string(predictionDataMessage.Timestamp))
	// close connection when done
	conn.Close()

}

func tcpListener(wg *sync.WaitGroup) {
	defer wg.Done()
	// Server listens on all interfaces for TCP connestion
	addr := ":" + "5051"

	listener, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalln("Server Failed To Start ")
	}

	// Wait for clients to connect
	for {
		// Accept a connection and spin-off a goroutine
		conn, err := listener.Accept()
		if err != nil {
			// If error continue to wait for other clients to connect
			continue
		}
		go handlePredictionDataMessage(conn)
	}
}

func main() {
	// tcp listener

	var wg sync.WaitGroup
	wg.Add(1)

	go tcpListener(&wg)

	wg.Wait()

}
