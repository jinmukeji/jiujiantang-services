#!/usr/bin/env bash
set -e
set -u
set -o pipefail

CUR=`dirname $0`

# 金姆技术部公众号调试
go run ${CUR} \
    --server_address=:9100
