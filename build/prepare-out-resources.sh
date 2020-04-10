#!/usr/bin/env bash
set -e
set -u
set -o pipefail

CUR=$(dirname $0)
OUT_DIR=${CUR}/out
ASSETS_DIR=${CUR}/assets

. ${CUR}/util.sh

info "Preparing output resources..."
if [ ! -d "${OUT_DIR}" ]; then
    mkdir -p ${OUT_DIR}
fi

# AE rule and data files
info "--> preparing ae v1 data..."
if [ -d "${OUT_DIR}/ae_data" ]; then
    rm -rf ${OUT_DIR}/ae_data
fi
mkdir -p ${OUT_DIR}/ae_data
tar xzf ${ASSETS_DIR}/ae_data-v1.0.1.tar.gz -C ${OUT_DIR}/ae_data

info "--> preparing sys config data..."
if [ -d "${OUT_DIR}/sys" ]; then
    rm -rf ${OUT_DIR}/sys
fi
mkdir -p ${OUT_DIR}/sys_data
cp ${ASSETS_DIR}/sys_data/config.yml ${OUT_DIR}/sys_data


info "--> preparing ae v2 data..."

AE_V2_VERSION=$(cat ${CUR}/assets/AE_VERSION)

AE_V2_DIR="assets/ae_data_v2/${AE_V2_VERSION}"

# 创建存储 v2 版本的 ae_data
if [ -d "${OUT_DIR}/ae_data_v2" ]; then
    rm -rf ${OUT_DIR}/ae_data_v2
fi
mkdir -p ${OUT_DIR}/ae_data_v2

AE_DATA_V2=${OUT_DIR}/ae_data_v2

cp -r ${AE_V2_DIR}/lua_src-${AE_V2_VERSION} ${AE_DATA_V2}/lua_src

cp -r ${AE_V2_DIR}/lookups-${AE_V2_VERSION} ${AE_DATA_V2}/lookups

cp -r ${AE_V2_DIR}/question-${AE_V2_VERSION} ${AE_DATA_V2}/question

cp ${AE_V2_DIR}/biz_conf-${AE_V2_VERSION}/presets.yaml ${AE_DATA_V2}/presets.yaml

echo
info "Finished to prepare output resources."
