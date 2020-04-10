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

info "Prune dangling images"
docker image prune -f

info "Begin to prune all existing docker images with 'latest' tag"

for b in ${BUILDS[@]}
do
    if [[ -n "$(docker images -q ${AWS_ECR_REPO}/${DOCKER_IMAGE_NAMESPACE}/$b:latest)" ]]; then
        docker rmi -f ${AWS_ECR_REPO}/${DOCKER_IMAGE_NAMESPACE}/$b:latest
    fi 
done

docker image ls -f "${DEFAULT_IMAGE_FILTER}"

info "Pruning all docker images is done."

