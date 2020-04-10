#!/usr/bin/env bash
set -e
set -u
set -o pipefail

aws ecr get-login-password | \
    docker login \
        --username AWS \
        --password-stdin 949191617935.dkr.ecr.cn-north-1.amazonaws.com.cn
