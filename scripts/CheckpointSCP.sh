#!/bin/bash

# getting arguments from the command line
if [ ! $# -eq 8 ]; then
	echo "USAGE: sudo ./CheckpointSCP.sh -c CONTAINER_ID -d DESTINATION_IP -u USER -n CHECKPOINT_NAME
EXAMPLE: sudo ./CheckpointSCP.sh -c hkj3434ljl43 -u ssakpal -n first -d 10.10.3.7"
	exit 1
fi

while [ $# -gt 0 ]; do
	case $1 in
		-c )	shift
			CONTAINER_ID=$1
			;;
		-d )	shift
			DESTINATION_IP=$1
			;;
    -u )	shift
  		USER=$1
  		;;
		-n )	shift
			CHECKPOINT_NAME=$1
			;;
		* )	echo >&2 "USAGE: sudo ./CheckpointSCP.sh -c CONTAINER_ID -d DESTINATION_IP  -u USER -n CHECKPOINT_NAME
    EXAMPLE sudo ./CheckpointSCP.sh -c hkj3434ljl43 -n first -d 10.10.3.7"
			exit 1
	esac
	shift
done


#SCP checkpoint file to the DESTINATION_IP
DIRECTORY="/home/$USER/checkpoint"
CHECKPOINT_LOCATION="$DIRECTORY/$CHECKPOINT_NAME"
SCP_LOCATION="$USER@$DESTINATION_IP:$DIRECTORY/"
sudo rsync -r $CHECKPOINT_LOCATION $SCP_LOCATION
if [ $? != 0 ]; then
  echo "Checkpoint SCP Failed"
  exit 1
fi
