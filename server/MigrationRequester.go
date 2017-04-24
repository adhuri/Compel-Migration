package main

import (
	"encoding/gob"
	"net"

	"github.com/Sirupsen/logrus"
	"github.com/adhuri/Compel-Migration/protocol"
	"github.com/adhuri/Compel-Migration/server/model"
)

func SendMigrationRequest(request *protocol.CheckpointRequest, server *model.Server, log *logrus.Logger) error {
	addr := request.SourceAgentIP + ":" + "5052"
	conn, err := net.Dial("tcp", addr)
	defer conn.Close()
	if err != nil {
		log.Warn("Server Not Alive")
		return err
	}

	encoder := gob.NewEncoder(conn)
	err = encoder.Encode(request)
	if err != nil {
		// If error occurs in sending a connect message to server then return
		log.Errorln("Failure While Sending Migration Request To Agent" + err.Error())
		return err
	}
	log.Infoln("Migration Request Successfully Sent - ", request)

	// read ack from the server
	agentReply := protocol.CheckpointResponse{}
	decoder := gob.NewDecoder(conn)
	err = decoder.Decode(&agentReply)
	if err != nil {
		// If error occurs while reading ACK from server then return
		log.Errorln("Bad Reply From Server " + err.Error())
		return err
	}
	// Print the ACK received from the server
	log.Infoln("Migration completed successfully")
	// If everything goes well return nil error
	return nil

}
