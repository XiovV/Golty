#!/bin/bash
set -ev

# Set Build Options
PLATFORM=arm
DOCKERFILE_LOCATION="./docker/Dockerfile.arm"
DOCKER_IMAGE="go-auto-yt"
DOCKER_TAG="stable"

# If This Isn't A PR, Push to Dockerhub
if [ "${TRAVIS_PULL_REQUEST}" = "false" ]; then
    docker login -u $DOCKER_USER -p $DOCKER_PASSWORD

    buildctl build --frontend dockerfile.v0 \
        --progress=plain \
        --opt platform=linux/${PLATFORM} \
        --opt filename=${DOCKERFILE_LOCATION} \
        --opt build-arg:TRAVIS_PULL_REQUEST=${TRAVIS_PULL_REQUEST} \
        --output type=image,name=docker.io/${DOCKER_USER}/${DOCKER_IMAGE}:${DOCKER_TAG}-${PLATFORM},push=true \
        --local dockerfile=. \
        --local context=.
else
    # If This is a PR, Build to Check for Errors
    buildctl build --frontend dockerfile.v0 \
        --progress=plain \
        --opt platform=linux/${PLATFORM} \
        --opt filename=${DOCKERFILE_LOCATION} \
        --opt build-arg:TRAVIS_PULL_REQUEST=false \
        --output type=docker,name=${DOCKER_IMAGE}:${DOCKER_TAG}-${PLATFORM} \
        --local dockerfile=. \
        --local context=. \
        | docker load
fi