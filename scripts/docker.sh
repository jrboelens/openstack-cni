#!/bin/bash
set -e

function usage() {
    echo "usage $0 <build|push> <values file>"
    exit 1
}

COMMAND=$1
if [ "$COMMAND" == "" ]; then
    usage
fi

VALUES_FILE=$2
if [ "$VALUES_FILE" == "" ]; then
    usage
fi

REPO=$(python -c 'import sys,yaml;print(yaml.safe_load(sys.stdin.read())["image"]["push-repository"])' < $VALUES_FILE)
TAG=$(python -c 'import sys,yaml;print(yaml.safe_load(sys.stdin.read())["image"]["tag"])' < $VALUES_FILE)


if [ $COMMAND = "build" ]; then
    echo "${COMMAND}ing ${REPO}:${TAG}"
    docker build -t $REPO:$TAG ./
elif [ $COMMAND = "push" ]; then
    echo "${COMMAND}ing ${REPO}:${TAG}"
    docker push $REPO:$TAG
else
    echo "invalid command $COMMAND"
    usage
fi