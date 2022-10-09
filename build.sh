#!/bin/zsh

CONTAINER_IMAGE_OS="linux"
CONTAINER_IMAGE_ARCH="amd64"
DOCKER_FILE="Dockerfile"
IMG_NAME="mlog"

~/go/bin/gox -osarch="$CONTAINER_IMAGE_OS/$CONTAINER_IMAGE_ARCH"
docker build -f $DOCKER_FILE -t $IMG_NAME .