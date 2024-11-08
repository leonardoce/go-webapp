#!/usr/bin/bash

# This script deploy the example web application in the local Kind
# cluster

export KO_DOCKER_REPO=kind.local
export KIND_CLUSTER_NAME=$(kind get clusters)

cd "$(dirname "$0")"/..
ko apply -f k8s
