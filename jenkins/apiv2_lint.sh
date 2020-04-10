#!/usr/bin/env bash
set -e

CUR=$(dirname $0)/../apitest 
echo $CUR
cd ${CUR} && pylint */**
