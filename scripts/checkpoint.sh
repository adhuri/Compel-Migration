#!/bin/bash

# getting inputs from the command line

#checkpoint a container
start=`date +%s.%N`
stuff
end=`date +%s.%N`
runtime=$((end-start))
echo $runtime
