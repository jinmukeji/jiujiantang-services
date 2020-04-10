#!/usr/bin/env bash
set -e
set -u
set -o pipefail

CUR=`dirname $0`

go run ${CUR} \
    --server_address=:9100 \
    --x_jwt_sign_in_key=jinmuhealth \
    --x_config_file=./data/config.yml
