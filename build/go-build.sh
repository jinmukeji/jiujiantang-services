#!/usr/bin/env bash

# Examples:
#   $ ./go-build.sh svc-biz-core
#   $ GOOS=linux ./go-build.sh svc-biz-core

set -e
set -u
set -o pipefail

CUR=`dirname $0`
REPO_ROOT=$(realpath ${CUR}/..)

# import util.sh
. ${CUR}/util.sh

# Check all  required commands
function checkAllRequiredCommands () {
    requireCommand go

    # https://stedolan.github.io/jq/
    requireCommand jq

    info "[OK] All required command tools are available."
}

warn "Checking tools"
checkAllRequiredCommands

BUILD_CONFIG_FILE="${CUR}/go-build-config.json"
BUILD_OUT_DIR=${CUR}/out

if [ ! -d "${BUILD_OUT_DIR}" ]; then
  mkdir -p ${BUILD_OUT_DIR}
fi

# Build context info
BUILD_NAME=$1
# check config
CHECK_BUILD_NAME=$(jq --raw-output '.["go-builds"][] | select(.name == "'$BUILD_NAME'")' $BUILD_CONFIG_FILE)

if [ -z "${CHECK_BUILD_NAME}" ]
then
    error "[ERROR] $BUILD_NAME is not found in config file $BUILD_CONFIG_FILE"
    exit 1
fi


BUILD_PRODUCT_VERSION=$(jq --raw-output '.["go-builds"][] | select(.name == "'$BUILD_NAME'") | .version' $BUILD_CONFIG_FILE)
BUILD_PACKAGE_NAME=$(jq --raw-output '.["go-builds"][] | select(.name == "'$BUILD_NAME'") | .package' $BUILD_CONFIG_FILE)
BUILD_ENABLED_CGO=$(jq --raw-output '.["go-builds"][] | select(.name == "'$BUILD_NAME'") | (if .enabled_cgo then "1" else "0" end)' $BUILD_CONFIG_FILE)
BUILD_FILES=${CUR}/$(jq --raw-output '.["go-builds"][] | select(.name == "'$BUILD_NAME'") | .build_files' $BUILD_CONFIG_FILE)
GIT_SHA=`git rev-parse HEAD`


if [[ -v "GIT_BRANCH" ]]; then
  GIT_BRANCH=${GIT_BRANCH}
else
  GIT_BRANCH=`git symbolic-ref --short HEAD`
fi


GO_VERSION="`go version`"
BUILD_VERSION=`git describe --always --long --dirty`
BUILD_TIME=`date +%FT%T%z`
BUILD_OUT_SUFFIX=$(go env GOOS)_$(go env GOARCH)
BUILD_OUT_NAME=${BUILD_NAME}_${BUILD_OUT_SUFFIX}

echo
warn "Building ${BUILD_NAME} for ${BUILD_OUT_SUFFIX}"
info "Package:              ${BUILD_PACKAGE_NAME}"
info "Product Version:      ${BUILD_PRODUCT_VERSION}"
info "Git SHA:              ${GIT_SHA}"
info "Git Branch:           ${GIT_BRANCH}"
info "Build Version:        ${BUILD_VERSION}"
info "Build Time:           ${BUILD_TIME}"
info "Go Version:           ${GO_VERSION}"
info "GOOS:                 $(go env GOOS)"
info "GOARCH:               $(go env GOARCH)"
info "Enabled CGO:          ${BUILD_ENABLED_CGO}"

CGO_ENABLED=${BUILD_ENABLED_CGO} go build -v \
    -gcflags "all=-trimpath=${REPO_ROOT}" \
    -asmflags "all=-trimpath=${REPO_ROOT}" \
    -ldflags "\
    -s \
    -X \"${BUILD_PACKAGE_NAME}/config.ProductVersion=${BUILD_PRODUCT_VERSION}\" \
    -X \"${BUILD_PACKAGE_NAME}/config.GitSHA=${GIT_SHA}\" \
    -X \"${BUILD_PACKAGE_NAME}/config.GitBranch=${GIT_BRANCH}\" \
    -X \"${BUILD_PACKAGE_NAME}/config.GoVersion=${GO_VERSION}\" \
    -X \"${BUILD_PACKAGE_NAME}/config.BuildVersion=${BUILD_VERSION}\" \
    -X \"${BUILD_PACKAGE_NAME}/config.BuildTime=${BUILD_TIME}\" " \
    -o ${BUILD_OUT_DIR}/${BUILD_OUT_NAME} ${BUILD_FILES}

echo
warn "[DONE] Build outputs to: ${BUILD_OUT_DIR}/${BUILD_OUT_NAME}"
