#!/usr/bin/env bash
set -e
set -u
set -o pipefail

CUR=`dirname $0`

# import util.sh
. ${CUR}/util.sh
. ${CUR}/const.sh

DOCKER_IMAGES=$(docker image ls -f "${DEFAULT_IMAGE_FILTER}" | grep "${AWS_ECR_REPO}/${DOCKER_IMAGE_NAMESPACE}" | awk '{print $1":"$2}')

for img in ${DOCKER_IMAGES[@]}
do
    warn "Pushing $img"
    docker push $img
    echo -----------------------------------------------------------------------------------------------------------------
    echo
done

warn "Pushing all docker images is done."
