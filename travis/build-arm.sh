#!/bin/bash
set -ev

PLATFORM=arm
DOCKERFILE_LOCATION="./docker/Dockerfile.arm"
DOCKER_USER="xiovv"
DOCKER_IMAGE="go-auto-yt"
DOCKER_TAG="stable"

if [ "${TRAVIS_PULL_REQUEST}" = "false" ]; then
    docker login -u $DOCKER_USER -p $DOCKER_PASSWORD

    buildctl build --frontend dockerfile.v0 \
        --frontend-opt platform=linux/${PLATFORM} \
        --frontend-opt filename=./${DOCKERFILE_LOCATION} \
        --exporter image \
        --exporter-opt name=docker.io/${DOCKER_USER}/${IMAGE}:${TAG}-${PLATFORM} \
        --exporter-opt push=true \
        --local dockerfile=. \
        --local context=.
fi