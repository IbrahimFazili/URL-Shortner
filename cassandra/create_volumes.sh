#!/bin/bash
USAGE="Usage: $0 IP1 IP2 IP3 ..."

if [ "$#" == "0" ]; then
        echo "$USAGE"
        exit 1
fi
while (( "$#" )); do
	sshpass -p hhhhiotwwg ssh student@$1 "mkdir -p data/cassandra/; mkdir -p data/redis/; mkdir -p data/url-shortner/; mkdir -p data/url-writer/"
	echo "created /data folder with cassandra, redis and url-shortner subfolders folder in $1"
	shift
done
