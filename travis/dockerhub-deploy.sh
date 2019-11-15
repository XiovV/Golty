#!/bin/bash
set -ev
if [ "${TRAVIS_PULL_REQUEST}" = "false" ]; then
    docker login -u $DOCKER_USER -p $DOCKER_PASSWORD
    docker build -t go-auto-yt .
    docker tag go-auto-yt xiovv/go-auto-yt:stable
    docker push xiovv/go-auto-yt:stable
fi