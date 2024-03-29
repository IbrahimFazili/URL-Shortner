#!/bin/bash
USAGE="Usage: $0 IP1 IP2"

if [ "$#" == "0" ]; then
    echo "$USAGE"
    exit 1
fi

MASTER="$1"
NODE="$2"
COMMAND="docker run --name cassandra-node -v /home/student/data/cassandra/:/var/lib/cassandra -d -e CASSANDRA_BROADCAST_ADDRESS=$NODE -p 7000:7000 -p 9042:9042 -e CASSANDRA_SEEDS=$MASTER cassandra"

sshpass -p hhhhiotwwg ssh student@$NODE "docker container stop cassandra-node; docker container rm cassandra-node; $COMMAND;"
while true;
do
    sleep 5
    STATUS=$(docker exec -it cassandra-node nodetool status | grep -e $NODE)
    STATUSUN=$(echo $STATUS | grep -e "UN")
    echo $STATUS
    [[ ! -z "$STATUSUN" ]] && break;
done;