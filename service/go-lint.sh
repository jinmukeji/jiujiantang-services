#!/usr/bin/env bash
set -e
CUR=`dirname $0`

# Run gometalinter with configuration file.
# Refer to https://github.com/alecthomas/gometalinter#configuration-file
gometalinter \
  --config ${CUR}/.gometalinter.json \
  ${CUR}/...
