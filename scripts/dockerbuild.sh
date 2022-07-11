#!/usr/bin/env bash

set -e

ROOT=$(git rev-parse --show-toplevel)
REGISTRY="dockerfactory.rsglab.com"
IMAGE_NAME="$REGISTRY/rsg/stackdriver_exporter"
ROOT=$(git rev-parse --show-toplevel)
if git rev-parse --short HEAD &> /dev/null; then
    BUILD_TAG=$(git rev-parse --short HEAD)
else
    BUILD_TAG=$(date '+%Y-%m-%d-%H-%M')
fi

IMAGE="$IMAGE_NAME:$BUILD_TAG"

build() {
    echo "Building $IMAGE"
    docker build -t "$IMAGE" "$ROOT"
}

push() {
    echo "Pushing $IMAGE"
    docker push "$IMAGE"
}

if [ -z "$1" ]; then
    echo "Usage: $0 <function>"
    echo "Available Functions: build, push"
fi

"$1"