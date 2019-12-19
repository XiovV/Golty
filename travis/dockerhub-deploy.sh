#!/bin/bash
set -ev
export DOCKER_CLI_EXPERIMENTAL=enabled

if [ "${TRAVIS_PULL_REQUEST}" = "false" ]; then
    docker login -u $DOCKER_USER -p $DOCKER_PASSWORD

    docker manifest create ${DOCKER_USER}/go-auto-yt:stable \
            ${DOCKER_USER}/go-auto-yt:stable-amd64 \
            ${DOCKER_USER}/go-auto-yt:stable-arm \
            ${DOCKER_USER}/go-auto-yt:stable-arm64

    docker manifest annotate ${DOCKER_USER}/go-auto-yt:stable ${DOCKER_USER}/go-auto-yt:stable-amd64 --arch amd64
    docker manifest annotate ${DOCKER_USER}/go-auto-yt:stable ${DOCKER_USER}/go-auto-yt:stable-arm --arch arm
    docker manifest annotate ${DOCKER_USER}/go-auto-yt:stable ${DOCKER_USER}/go-auto-yt:stable-arm64 --arch arm64

    docker manifest push ${DOCKER_USER}/go-auto-yt:stable
else
    docker images
fi