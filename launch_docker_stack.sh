#!/bin/bash
docker stack rm url_stack
sleep 10
docker build -t ibrahimfazili/urlshortner:latest .
docker image push ibrahimfazili/urlshortner:latest
docker build -t ibrahimfazili/urlwrite -f Dockerfile.writer .
docker image push ibrahimfazili/urlwrite:latest
docker stack deploy -c docker-compose.yml url_stack