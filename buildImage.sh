#!/usr/bin/env sh

REPO=""
IMG="go-http-replier"
TAG=`date +%s`
BUILDID=$1

if [ ! -z ${BUILDID} ]; then
    TAG=${BUILDID}
fi

if ! docker build --tag ${REPO}${IMG}:${TAG} . ; then
    echo "Failed to build image!"
    exit 1
fi
echo "Image: ${REPO}${IMG}:${TAG} has been built."
