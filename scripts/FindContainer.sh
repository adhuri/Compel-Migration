#!/bin/bash

# getting arguments from the command line
if [ ! $# -eq 2 ]; then
	echo "USAGE: sudo ./FindContainer.sh -c CONTAINER_ID
EXAMPLE: sudo ./FindContainer.sh -c hkj3434ljl43"
	exit 1
fi

while [ $# -gt 0 ]; do
	case $1 in
		-c )	shift
			CONTAINER_ID=$1
			;;
		* )	echo >&2 "USAGE: sudo ./FindContainer.sh -c CONTAINER_ID
    EXAMPLE sudo ./FindContainer.sh -c hkj3434ljl43"
			exit 1
	esac
	shift
done

if [ "$(docker ps | grep $CONTAINER_ID)" ]; then
  echo "Container Present"
  exit 0
fi

exit 1
