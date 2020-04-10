#!/usr/bin/env bash
set -e

CUR=$(dirname $0)/../api-jinmuid/webapi
echo $CUR
cd ${CUR} && pylint */**
