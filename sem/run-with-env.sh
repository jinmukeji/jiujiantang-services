#!/usr/bin/env bash
set -e
set -u
set -o pipefail

CUR=`dirname $0`

ENV_FILE=${CUR}/../build/local.svc-sem-gw.env
env $(cat ${ENV_FILE} | grep -v '^#'| xargs) go run ${CUR} \
   --server_address=:9091
