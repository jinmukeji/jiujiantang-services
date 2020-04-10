#!/usr/bin/env bash
set -e
set -u
set -o pipefail

CUR=`dirname $0`

. ${CUR}/const.sh

docker-compose \
    --file ${CUR}/docker-compose.yml \
    --project-name ${DOCKER_COMPOSE_PROJECT_NAME} down

docker-compose \
    --file ${CUR}/docker-compose.yml \
    --project-name ${DOCKER_COMPOSE_PROJECT_NAME} rm

