#!/bin/bash
NAME=pxnginx_ubuntu18.04

docker rm -v -f $NAME
docker build -t $NAME -f Dockerfile .
docker run \
    -p 8081:80 \
    -it --name $NAME $NAME
