#!/usr/bin/env bash
set -e
set -u
set -o pipefail

CUR=`dirname $0`

. ${CUR}/build/util.sh

cd $CUR

function update_all() {
    echo
    info "[1/2] Update all go packages..."
    export GO111MODULE=on
    go env -w GOPROXY=https://goproxy.io,direct
    go env -w GOPRIVATE="github.com/jinmukeji/*"

    go mod tidy
    go get -u all

    echo
    info "[2/2] Tidy go mod packages..."
    go mod tidy
}

update_all

echo
info "Done"
