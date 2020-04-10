#!/usr/bin/env bash
set -e

CUR=$(dirname $0)
echo $CUR
cd ${CUR}/.. && golangci-lint run
