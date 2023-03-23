#!/bin/bash
USAGE="Usage: $0 IP1 IP2 IP3 ..."

if [ "$#" == "0" ]; then
        echo "$USAGE"
        exit 1
fi

while (( "$#" )); do
        sshpass -p hhhhiotwwg ssh student@$1 "docker swarm leave --force"
        hostName=`sshpass -p hhhhiotwwg ssh student@$1 hostname`
        sleep 5
        docker node rm $hostName
        echo "Removed $1 from swarm"
        shift
done