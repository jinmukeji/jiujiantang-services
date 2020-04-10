#!/usr/bin/env bash
set -e
set -u
set -o pipefail

CUR=`dirname $0`

ENV_FILE=${CUR}/../build/local.svc-jinmuid.env
# 喜马把脉技术部公众号调试
env $(cat ${ENV_FILE} | grep -v '^#'| xargs) go run ${CUR} \
    --server_address=:9100 \
    --x_jwt_sign_in_key=jinmuhealth \
