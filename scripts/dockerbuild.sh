#!/usr/bin/env bash

set -e

ROOT=$(git rev-parse --show-toplevel)
REGISTRY="dockerfactory.rsglab.com"
IMAGE="$REGISTRY/rsg/stackdriver_exporter"
ROOT=$(git rev-parse --show-toplevel)
if git rev-parse --short HEAD &> /dev/null; then
    BUILD_TAG=$(git rev-parse --short HEAD)
else
    BUILD_TAG=$(date '+%Y-%m-%d-%H-%M')
fi

build() {
    echo "Building $IMAGE"
    docker build -t "$IMAGE" "$ROOT"
}

push() {
    echo "Pushing $IMAGE"
    docker push "$IMAGE"
}