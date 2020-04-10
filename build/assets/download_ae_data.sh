#!/usr/bin/env bash
set -e
set -u
set -o pipefail

CUR=$(dirname $0)
cd ${CUR}

# 需要配置 GitHub API Token
# 在 https://github.com/settings/tokens 中配置，并将值配置到 shell 的 rc 文件启动中
TOKEN=$GITHUB_RELEASE_GITHUB_API_TOKEN
VERSION=$(cat AE_VERSION) # tag name or the word "latest"
REPO="jinmukeji/ae"
GITHUB="https://api.github.com"

alias errcho='>&2 echo'

function gh_curl() {
    curl -H "Authorization: token $TOKEN" \
        -H "Accept: application/vnd.github.v3.raw" \
        $@
}

function download_label() {
    local LABEL=$1
    local parser
    local name_parser

    if [ "$VERSION" = "latest" ]; then
        # Github should return the latest release first.
        parser=".[0].assets | map(select(.label == \"$LABEL\"))[0].id"
        name_parser=".[0].assets | map(select(.label == \"$LABEL\"))[0].name"
    else
        parser=". | map(select(.tag_name == \"$VERSION\"))[0].assets | map(select(.label == \"$LABEL\"))[0].id"
        name_parser=". | map(select(.tag_name == \"$VERSION\"))[0].assets | map(select(.label == \"$LABEL\"))[0].name"
    fi

    local asset_id=$(echo $GH_RELEASES | jq -r "$parser")
    local name=$(echo $GH_RELEASES | jq -r "$name_parser")

    if [ "$asset_id" = "null" ]; then
        errcho "ERROR: version not found $VERSION"
        exit 1
    fi

    echo "Downloading $name ($asset_id)"

    wget -q --auth-no-challenge --header='Accept:application/octet-stream' \
        https://$TOKEN:@api.github.com/repos/$REPO/releases/assets/$asset_id \
        -O $out/$name
}


echo "AE Version: ${VERSION}"

# Output dir
out=${CUR}/ae_data_v2/${VERSION}
if [ ! -d "${out}" ]; then
    mkdir -p ${out}

    GH_RELEASES=$(gh_curl -s $GITHUB/repos/$REPO/releases)

    download_label "Lookups Dictionary Data"
    download_label "Lua Source Codes"
    download_label "Question Dictionary Data"
    download_label "Biz Config Data"
else
    echo "SKIPED downloading. ${out} already exists."
fi

cd ${out}

# 解压下载的文件
for FILE in ./*.gz
do
    if [ -f ${FILE} ]; then
    OUT_DIR=${FILE%.*.*}
    echo ${OUT_DIR}
    mkdir -p ${OUT_DIR}
    tar -xzf ${FILE} -C ${OUT_DIR}
    fi
done

