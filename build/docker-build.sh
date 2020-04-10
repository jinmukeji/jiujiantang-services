#!/usr/bin/env bash
set -e
set -u
set -o pipefail

CUR=`dirname $0`

# import util.sh
. ${CUR}/util.sh
. ${CUR}/const.sh

# Check all required commands
function checkAllRequiredCommands () {
    requireCommand go

    # https://stedolan.github.io/jq/
    requireCommand jq

    info "[OK] All required command tools are available."
}

warn "Checking tools"
checkAllRequiredCommands

# Build context info
BUILD_IMAGE_NAME=$1
BUILD_CONFIG_FILE="${CUR}/docker-build-config.json"
# check config
CHECK_BUILD_NAME=$(jq --raw-output '.["docker-builds"][] | select(.image_name == "'$BUILD_IMAGE_NAME'")' $BUILD_CONFIG_FILE)
if [ -z "${CHECK_BUILD_NAME}" ]
then
    error "[ERROR] $BUILD_IMAGE_NAME is not found in config file $BUILD_CONFIG_FILE"
    exit 1
fi

BUILD_DOCKER_FILE=${CUR}/$(jq --raw-output '.["docker-builds"][] | select(.image_name == "'$BUILD_IMAGE_NAME'") | .dockerfile' $BUILD_CONFIG_FILE)
BUILD_WORKING_DIR=${CUR}/$(jq --raw-output '.["docker-builds"][] | select(.image_name == "'$BUILD_IMAGE_NAME'") | .working_dir' $BUILD_CONFIG_FILE)

echo
warn "Building ${BUILD_IMAGE_NAME}"
info "Dockerfile:              ${BUILD_DOCKER_FILE}"
info "Working_dir:             ${BUILD_WORKING_DIR}"
docker build -t ${AWS_ECR_REPO}/${DOCKER_IMAGE_NAMESPACE}/${BUILD_IMAGE_NAME} --file ${BUILD_DOCKER_FILE} ${BUILD_WORKING_DIR}
if [ "$(uname)" = "Darwin" ]; then
    # Build on macOS system
    IMAGE_INFORMATION=`${CUR}/out/${BUILD_IMAGE_NAME}_darwin_amd64 --version`
else
    # Build on Linux system
    IMAGE_INFORMATION=`${CUR}/out/${BUILD_IMAGE_NAME}_linux_amd64 --version`
fi

IMAGE_VERSION=`echo  ${IMAGE_INFORMATION}| cut -f 1 -d " "`
docker tag ${AWS_ECR_REPO}/${DOCKER_IMAGE_NAMESPACE}/${BUILD_IMAGE_NAME}:latest ${AWS_ECR_REPO}/${DOCKER_IMAGE_NAMESPACE}/${BUILD_IMAGE_NAME}:${IMAGE_VERSION}
