#!/bin/bash
USAGE="Usage: $0 IP1 IP2 IP3 ..."

if [ "$#" == "0" ]; then
        echo "$USAGE"
        exit 1
fi

while (( "$#" )); do
	# ./tear_volumes.sh "$1"
        sshpass -p hhhhiotwwg ssh student@$1 "docker container stop cassandra-node; docker container rm cassandra-node;"
        echo "Stopped cassandra on $1"
        shift
done
