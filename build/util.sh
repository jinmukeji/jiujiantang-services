#!/usr/bin/env bash

# Useful Functions

function info() {
    local LIGHT_GREEN='\033[1;32m'
    local NC='\033[0m' # No Color

    printf "${LIGHT_GREEN}$1${NC}\n"
}

function warn() {
    local YELLOW='\033[33m'
    local NC='\033[0m' # No Color

    printf "${YELLOW}$1${NC}\n"
}

function error() {
    local RED='\033[0;31m'
    local NC='\033[0m' # No Color

    printf "${RED}$1${NC}\n"
}

function pushd () {
    builtin pushd "$@" > /dev/null
}

function popd () {
    builtin popd "$@" > /dev/null
}

# Check required command
function requireCommand () {
    type $1 >/dev/null 2>&1 || { error >&2 "[Aborting] Requires $1 but it's not installed."; exit 1; }
}
