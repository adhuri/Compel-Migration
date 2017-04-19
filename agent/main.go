package main

import (
	"encoding/gob"
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

func handleMigrationRequest(conn net.Conn) {
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

	// Create a ConnectAck Message
	migrationAck := protocol.NewCheckpointResponse(migrationRequest, true)

	// Send Connect Ack back to the client
	encoder := gob.NewEncoder(conn)
	err = encoder.Encode(migrationAck)
	//err = binary.Write(conn, binary.LittleEndian, connectAck)
	if err != nil {
		// If failure in parsing, close the connection and return
		log.Errorln("Prediction Data Ack")
		return
	}
	// close connection when done
	conn.Close()

}

func tcpListener(wg *sync.WaitGroup) {
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
		go handleMigrationRequest(conn)
	}
}

func main() {

	var wg sync.WaitGroup
	wg.Add(1)

	go tcpListener(&wg)

	wg.Wait()

}
