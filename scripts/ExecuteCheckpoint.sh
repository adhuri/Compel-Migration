#!/bin/bash

# getting arguments from the command line
if [ ! $# -eq 8 ]; then
	echo "USAGE: sudo ./checkpoint.sh -c CONTAINER_ID -d DESTINATION_IP -u USER -n CHECKPOINT_NAME
EXAMPLE: sudo ./checkpoint.sh -c hkj3434ljl43 -u ssakpal -n first -d 10.10.3.7"
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
		* )	echo >&2 "USAGE: sudo ./checkpoint.sh -c CONTAINER_ID -d DESTINATION_IP  -u USER -n CHECKPOINT_NAME
    EXAMPLE sudo ./checkpoint.sh -c hkj3434ljl43 -n first -d 10.10.3.7"
			exit 1
	esac
	shift
done


#creating a checkpoint directory if not present
DIRECTORY="/home/$USER/checkpoint"
if [ ! -d "$DIRECTORY" ]; then
  mkdir $DIRECTORY
	sudo chown $USER $DIRECTORY
  echo "Checkpoint directory created"
fi


#checkpoint a container
sudo docker checkpoint create --checkpoint-dir $DIRECTORY $CONTAINER_ID $CHECKPOINT_NAME
if [ $? != 0 ]; then
  echo "Checkpoint Execution Failed"
  exit 1
fi
