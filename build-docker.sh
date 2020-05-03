#!/usr/bin/env bash

TAG_PREFIX="cmsc483-project"
LB_TAG="load-balancer"
LB_DOCKERFILE="load-balancer.dockerfile"
WS_TAG="simple-http-server"
WS_DOCKERFILE="simple-http-server.dockerfile"

echo "Building Images!"
docker build -f ${LB_DOCKERFILE} -t ${TAG_PREFIX}/${LB_TAG} . && 
docker build -f ${WS_DOCKERFILE} -t ${TAG_PREFIX}/${WS_TAG} . && 
echo "Done building images! Tagged ${TAG_PREFIX}/{${LB_TAG},${WS_TAG}}"
