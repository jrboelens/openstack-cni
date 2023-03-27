#!/bin/bash
set -e

function usage() {
    echo "usage $0 <values file>"
    exit 1
}

VALUES_FILE=$1
if [ "$VALUES_FILE" == "" ]; then
    usage
fi

helm upgrade openstack-cni helm --install -f $VALUES_FILE