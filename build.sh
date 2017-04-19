#!/bin/bash

REPO_NAME="Compel-Migration"
SERVER_NAME="server"
AGENT_NAME="agent"

echo "Building $REPO_NAME "

if go build -o $GOPATH/bin/$SERVER_NAME -i github.com/adhuri/$REPO_NAME/$SERVER_NAME ;then
mv $GOPATH/bin/$SERVER_NAME $GOPATH/bin/migration_server
echo "+Successful"
else echo "-Failed"
fi


echo "Building $AGENT_NAME"

if go build -o $GOPATH/bin/$AGENT_NAME -i github.com/adhuri/$REPO_NAME/$AGENT_NAME ; then
mv $GOPATH/bin/$AGENT_NAME $GOPATH/bin/migration_agent
echo "+Successful"
else echo "-Failed"
fi
