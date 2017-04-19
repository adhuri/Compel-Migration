#!/bin/bash

REPO_NAME="Compel-Migration"
SERVER_NAME="server"
AGENT_NAME="agent"

echo "Building compel-migration-$SERVER_NAME "

if go build -o $GOPATH/bin/$SERVER_NAME -i github.com/adhuri/$REPO_NAME/$SERVER_NAME ;then
mv $GOPATH/bin/$SERVER_NAME $GOPATH/bin/Compel-Migration-$SERVER_NAME
echo "+Successful"
else echo "-Failed"
fi


echo "Building compel-migration-$AGENT_NAME"

if go build -o $GOPATH/bin/$AGENT_NAME -i github.com/adhuri/$REPO_NAME/$AGENT_NAME ; then
mv $GOPATH/bin/$AGENT_NAME $GOPATH/bin/Compel-Migration-$AGENT_NAME
echo "+Successful"
else echo "-Failed"
fi
