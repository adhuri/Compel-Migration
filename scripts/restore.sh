#!/bin/bash

# getting arguments from the command line
if [ ! $# -eq 6 ]; then
	echo "USAGE: sudo ./checkpoint.sh -c CONTAINER_ID -u USER -n CHECKPOINT_NAME
EXAMPLE: sudo ./checkpoint.sh -c hkj3434ljl43 -u ssakpal -n first"
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
		* )	echo >&2 "USAGE: sudo ./checkpoint.sh -c CONTAINER_ID  -u USER -n CHECKPOINT_NAME
    EXAMPLE sudo ./checkpoint.sh -c hkj3434ljl43 -n first"
			exit 1
	esac
	shift
done

ARGS=""
PORT_MAPPING=""
CONTAINER_NAME=""
filename="/home/$USER/${CONTAINER_ID}_metadata.conf"
while read -r line
do
    value="$line"
    if [[ $value == *"ENV"* ]] || [[ $value == *"CMD"* ]] || [[ $value == *"EXPOSE"* ]] || [[ $value == *"ENTRYPOINT"* ]]; then
        #echo "It's there! : $value"
        ARGS="$ARGS --change \"$value\""
    else
      if [[ $value == *"-p"* ]]; then
        echo "PORT arguments : $value"
        PORT_MAPPING=$value
      else
        echo "CONTAINER NAME : $value"
        CONTAINER_NAME=$value
      fi
    fi
    #echo "Name read from file - $name"
done < "$filename"

echo "$ARGS"
echo "$CONTAINER_NAME"
echo "$PORT_MAPPING "

TAR_NAME="/home/$USER/$CHECKPOINT_NAME.tar"
echo $TAR_NAME
DOCKER_IMPORT_COMMAND="docker import $ARGS $TAR_NAME"
#eval $IMAGE
IMAGE=$(eval $DOCKER_IMPORT_COMMAND)
#c=${!IMAGE}
echo "IMAGE NAME: $IMAGE"
