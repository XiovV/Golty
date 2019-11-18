#!/bin/bash
set -ev
docker build --build-arg TRAVIS_PULL_REQUEST -t go-auto-yt .
if [ "${TRAVIS_PULL_REQUEST}" = "false" ]; then
    docker login -u $DOCKER_USER -p $DOCKER_PASSWORD
    docker tag go-auto-yt xiovv/go-auto-yt:stable
    docker push xiovv/go-auto-yt:stable
fi