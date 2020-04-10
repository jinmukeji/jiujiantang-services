#!/usr/bin/env bash
set -e
set -u
set -o pipefail

CUR=`dirname $0`

# import util.sh
. ${CUR}/util.sh
. ${CUR}/const.sh

BUILD_CONFIG_FILE="${CUR}/docker-build-config.json"

BUILDS=$(jq --raw-output '.["docker-builds"][].image_name' ${BUILD_CONFIG_FILE})

#copy files to be used to out
${CUR}/prepare-out-resources.sh

# Prune old images at first
${CUR}/prune-docker-images.sh

info "Begin to docker build all"

for b in ${BUILDS[@]}
do
    $CUR/docker-build.sh $b
    echo
    echo -------------------------------------------------------------------
    echo
done

echo
info "Prune dangling images"
docker image prune -f

docker image ls -f "${DEFAULT_IMAGE_FILTER}"
echo
echo -------------------------------------------------------------------
echo

info "Build all is done."
