#!/bin/bash

docker build -t ibrahimfazili/urlshortner:latest .
docker image push ibrahimfazili/urlshortner:latest
docker build -t ibrahimfazili/urlwrite -f Dockerfile.writer .
docker image push ibrahimfazili/urlwrite:latest