package main

import (
	"encoding/gob"
	"fmt"
	"net"
	"os"
	"strconv"
	"sync"

	"github.com/Sirupsen/logrus"

	migrationProtocol "github.com/adhuri/Compel-Migration/protocol"
	model "github.com/adhuri/Compel-Migration/server/model"
	strategy "github.com/adhuri/Compel-Migration/server/strategy"
)

var (
	log    *logrus.Logger
	server *model.Server
)

func init() {

	log = logrus.New()

	// Output logging to stdout
	log.Out = os.Stdout

	// Only log the info severity or above.
	log.Level = logrus.DebugLevel
	// Microseconds level logging
	customFormatter := new(logrus.TextFormatter)
	customFormatter.TimestampFormat = "2006-01-02 15:04:05.000000"
	customFormatter.FullTimestamp = true

	log.Formatter = customFormatter

}

func main() {
	// tcp listener

	server = model.NewServer()

	var wg sync.WaitGroup
	wg.Add(1)

	go tcpListener(&wg, server)

	wg.Wait()

}

func handlePredictionDataMessage(conn net.Conn, server *model.Server) {
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

	//fmt.Println("Yayy ----------------------", predictionDataMessage)
	ts := strconv.FormatInt(predictionDataMessage.Timestamp, 10)
	fmt.Println("Timestamp", ts)
	log.Infoln("Prediction Data Ack Sent for Request Id " + ts)
	// close connection when done
	conn.Close()

	// migration decision
	migrationNeeded, migrationInfo := strategy.MigrationNeeded(&predictionDataMessage, server)

	// send migration request if decided to migrate
	if migrationNeeded {
		err = SendMigrationRequest(migrationInfo, server, log)
		if err != nil {
			log.Infoln("Migration Was Failure")
			return
		}
		log.Infoln("Migration Was Success")
	} else {
		log.Infoln("Migration Was Not Needed")
	}
}

func tcpListener(wg *sync.WaitGroup, server *model.Server) {
	defer wg.Done()
	// Server listens on all interfaces for TCP connestion
	addr := ":" + "5051"
	log.Infoln("Migration Server listening on TCP ", addr)
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
		log.Infoln(" Accepted Connection from Prediction Client ")
		go handlePredictionDataMessage(conn, server)
	}
}
