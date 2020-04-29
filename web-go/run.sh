#!/usr/bin/env bash
set -e
set -u
set -o pipefail

CUR=`dirname $0`

go run ${CUR} \
    --registry=mdns \
    --server_address=:9011 \
    --x_config_file=./data/resource.yml
