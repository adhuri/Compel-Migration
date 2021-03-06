#!/bin/bash

# getting arguments from the command line
if [ ! $# -eq 6 ]; then
	echo "USAGE: sudo ./restore.sh -c CONTAINER_ID -u USER -n CHECKPOINT_NAME
EXAMPLE: sudo ./restore.sh -c hkj3434ljl43 -u ssakpal -n first"
	exit 1
fi

while [ $# -gt 0 ]; do
	case $1 in
		-c )	shift
			CONTAINER_ID=$1
			;;
    -u )	shift
  		USER=$1
  		;;
		-n )	shift
			CHECKPOINT_NAME=$1
			;;
		* )	echo >&2 "USAGE: sudo ./restore.sh -c CONTAINER_ID  -u USER -n CHECKPOINT_NAME
    EXAMPLE sudo ./restore.sh -c hkj3434ljl43 -n first"
			exit 1
	esac
	shift
done


# Initializing Variables to use in future
ARGS=""
PORT_MAPPING=""
CONTAINER_NAME=""
filename="/home/$USER/${CONTAINER_ID}_metadata.conf"


# Read Metadata file line by line
while read -r line
do
    value="$line"
    if [[ $value == *"ENV"* ]] || [[ $value == *"CMD"* ]] || [[ $value == *"EXPOSE"* ]] || [[ $value == *"ENTRYPOINT"* ]]; then
        # METADATA used to import the image
        ARGS="$ARGS --change \"$value\""
    else
      if [[ $value == *"-p"* ]]; then
        # PORT_MAPPING
        PORT_MAPPING=$value
      else
        # CONTAINER_NAME
        CONTAINER_NAME=$value
      fi
    fi
done < "$filename"


# Importing the tar line by line
start=`date +%s%3N`
TAR_NAME="/home/$USER/$CHECKPOINT_NAME.tar"
DOCKER_IMPORT_COMMAND="sudo docker import $ARGS $TAR_NAME"
IMAGE=$(eval $DOCKER_IMPORT_COMMAND)
echo $?
end=`date +%s%3N`
runtime=$((end-start))
echo "Docker Image Importing took : $runtime milliseconds"


# Trimming the image name because output of previous command can't be directly used
IMAGE_NAME=${IMAGE##*:}


# Create a new Docker container by the same name as the original container
start=`date +%s%3N`
DOCKER_CREATE_COMMAND="sudo docker create --name $CONTAINER_NAME --network=my-net $PORT_MAPPING $IMAGE_NAME"
new_container_id=$(eval $DOCKER_CREATE_COMMAND)
echo $?
end=`date +%s%3N`
runtime=$((end-start))
echo "Docker Container Creation took : $runtime milliseconds"


# Restore the container
start=`date +%s%3N`
DIRECTORY="/home/$USER/checkpoint"
DOCKER_RESTORE_COMMAND="sudo docker start --checkpoint $CHECKPOINT_NAME --checkpoint-dir=\"$DIRECTORY\" $CONTAINER_NAME"
eval $DOCKER_RESTORE_COMMAND
echo $?
end=`date +%s%3N`
runtime=$((end-start))
echo "Docker Container Restoration took : $runtime milliseconds"


#clean up
# start=`date +%s%3N`
# CHECKPOINT_FILES="$DIRECTORY/$CHECKPOINT_NAME"
# docker rm $CONTAINER_ID
# rm TAR_NAME
# rm -rf CHECKPOINT_FILES
# end=`date +%s%3N`
# runtime=$((end-start))
# echo "Clean Up Execution Time : $runtime milliseconds"
