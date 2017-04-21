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
  echo "Checkpoint directory created"
fi

#capture metadata of the container
start=`date +%s%3N`
time python read_container_metadata.py -n $CONTAINER_ID -u $USER
end=`date +%s%3N`
runtime=$((end-start))
echo "Execution Time : $runtime milliseconds"

#checkpoint a container
start=`date +%s%3N`
time docker checkpoint create --checkpoint-dir $DIRECTORY $CONTAINER_ID $CHECKPOINT_NAME
end=`date +%s%3N`
runtime=$((end-start))
echo "Execution Time : $runtime milliseconds"

#exporting the file system
start=`date +%s%3N`
TAR_NAME="/home/$USER/$CHECKPOINT_NAME.tar"
time docker export $CONTAINER_ID > $TAR_NAME
end=`date +%s%3N`
runtime=$((end-start))
echo "Execution Time : $runtime milliseconds"

#SCP file system to DESTINATION_IP
start=`date +%s%3N`
SCP_LOCATION="$USER@$DESTINATION_IP:/home/$USER"
time scp -r $TAR_NAME $SCP_LOCATION
end=`date +%s%3N`
runtime=$((end-start))
echo "Execution Time : $runtime milliseconds"

#SCP checkpoint file to the DESTINATION_IP
start=`date +%s%3N`
CHECKPOINT_LOCATION="$DIRECTORY/$CHECKPOINT_NAME/"
SCP_LOCATION="$USER@$DESTINATION_IP:$DIRECTORY"
time scp -r $CHECKPOINT_LOCATION $SCP_LOCATION
end=`date +%s%3N`
runtime=$((end-start))
echo "Execution Time : $runtime milliseconds"

#SCP Metadata file
start=`date +%s%3N`
METADATA_LOCATION="/home/$USER/${CONTAINER_ID}_metadata.conf"
SCP_LOCATION="$USER@$DESTINATION_IP:/home/$USER"
time scp -r $METADATA_LOCATION $SCP_LOCATION
end=`date +%s%3N`
runtime=$((end-start))
echo "Execution Time : $runtime milliseconds"


#Run Restore.sh on the remote machine to restore the container
SSH_RESTORE_COMMAND="ssh root@$DESTINATION_IP 'bash -s' -- < restore.sh \"-c\" \"$CONTAINER_ID\" \"-u\" \"$USER\" \"-n\" \"$CHECKPOINT_NAME\""
restore_timing_info=$(eval $SSH_RESTORE_COMMAND)
echo $restore_timing_info


#clean up
start=`date +%s%3N`
docker rm $CONTAINER_ID
rm $TAR_NAME
end=`date +%s%3N`
runtime=$((end-start))
echo "Execution Time : $runtime milliseconds"
