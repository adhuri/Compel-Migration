#!/bin/bash

REPO_NAME="Compel-Migration"
SERVER_NAME="server"
#AGENT_NAME="compel-monitoring-agent"

echo "Building $REPO_NAME "

if go build -o $GOPATH/bin/$REPO_NAME -i github.com/adhuri/$REPO_NAME/$SERVER_NAME ;then
echo "+Successful"
else echo "-Failed"
fi


#echo "Building $AGENT_NAME"

#if go build -o $GOPATH/bin/$AGENT_NAME -i github.com/adhuri/$REPO_NAME/$AGENT_NAME ; then
#echo "+Successful"
#else echo "-Failed"
#fi
