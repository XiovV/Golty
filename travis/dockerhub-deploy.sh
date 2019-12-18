#!/bin/bash
set -ev
export DOCKER_CLI_EXPERIMENTAL=enabled

if [ "${TRAVIS_PULL_REQUEST}" = "false" ]; then
    docker login -u $DOCKER_USER -p $DOCKER_PASSWORD
    docker tag go-auto-yt xiovv/go-auto-yt:stable

    docker manifest create xiovv/go-auto-yt:stable \
            xiovv/go-auto-yt:stable-amd64 \
            xiovv/go-auto-yt:stable-arm \
            xiovv/go-auto-yt:stable-arm64

    docker manifest annotate xiovv/go-auto-yt:stable xiovv/go-auto-yt:stable-amd64 --arch amd64
    docker manifest annotate xiovv/go-auto-yt:stable xiovv/go-auto-yt:stable-arm --arch arm
    docker manifest annotate xiovv/go-auto-yt:stable xiovv/go-auto-yt:stable-arm64 --arch arm64

    docker manifest push xiovv/go-auto-yt:stable
else
    docker images
fi