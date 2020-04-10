#!/usr/bin/env bash

# Examples:
#   $ ./go-build-all.sh
#   $ GOOS=linux ./go-build-all.sh

set -e
set -u
set -o pipefail

CUR=`dirname $0`

# import util.sh
. ${CUR}/util.sh

BUILD_CONFIG_FILE="${CUR}/go-build-config.json"

BUILDS=$(jq --raw-output '.["go-builds"][].name' ${BUILD_CONFIG_FILE})

info "Begin to go build all"
echo

for b in ${BUILDS[@]}
do
    # Build for Linux
    GOOS=linux $CUR/go-build.sh $b

    # Build for macOS
    if [ "$(uname)" = "Darwin" ]; then
        GOOS=darwin $CUR/go-build.sh $b
    fi

    echo
    echo -------------------------------------------------------------------
    echo
done

info "Build all is done."
