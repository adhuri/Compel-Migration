#!/bin/bash

# getting arguments from the command line
if [ ! $# -eq 8 ]; then
	echo "USAGE: sudo ./CheckpointCleanup.sh -c CONTAINER_ID -d DESTINATION_IP -u USER -n CHECKPOINT_NAME
EXAMPLE: sudo ./CheckpointCleanup.sh -c hkj3434ljl43 -u ssakpal -n first -d 10.10.3.7"
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
		* )	echo >&2 "USAGE: sudo ./CheckpointCleanup.sh -c CONTAINER_ID -d DESTINATION_IP  -u USER -n CHECKPOINT_NAME
    EXAMPLE sudo ./CheckpointCleanup.sh -c hkj3434ljl43 -n first -d 10.10.3.7"
			exit 1
	esac
	shift
done


#Run Restore.sh on the remote machine to restore the container
SSH_RESTORE_COMMAND="ssh $USER@$DESTINATION_IP 'bash -s' -- < restore.sh \"-c\" \"$CONTAINER_ID\" \"-u\" \"$USER\" \"-n\" \"$CHECKPOINT_NAME\""
restore_timing_info=$(eval $SSH_RESTORE_COMMAND)

if [ $? != 0 ]; then
  echo "Restoration Failed"
  exit 1
fi

echo $restore_timing_info
