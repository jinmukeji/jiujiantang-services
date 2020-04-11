#!/usr/bin/env bash
set -e

# Required command-line tools
#   openapi-generator
#   redoc-cli
#   yq -- brew install python-yq

CUR=`dirname $0`
# TODO: 区分 Spec 生成 SDK
SPEC_FILE="$CUR/api-v2-oas3.yaml"
CONFIG_DIR="$CUR/openapi-generator-config"
TEMPLATE_DIR="$CUR/templates"
API_VERSION=`yq --raw-output '.info.version' $SPEC_FILE `
OUT_DIR="$CUR/sdk_out/$API_VERSION"
GROUP_ID="com.himalife.api"

function info() {
    echo -e "\033[33m$1\033[0m"
}

function generate()
{
    local LANG="$1"
    local MSG="$2"
    local CONFIG_FILE="$CONFIG_DIR/$LANG.json"
    local LANG_TEMPLATE="$TEMPLATE_DIR/$LANG"

    info "$MSG"

    if [ -f $CONFIG_FILE ]; then
        if [ -d $LANG_TEMPLATE ]; then
            openapi-generator generate \
                --template-dir $LANG_TEMPLATE \
                --group-id $GROUP_ID \
                --artifact-version $API_VERSION \
                -i $SPEC_FILE \
                -c $CONFIG_FILE \
                -g $LANG \
                -o $OUT_DIR/$LANG
        else
            openapi-generator generate \
                --group-id $GROUP_ID \
                --artifact-version $API_VERSION \
                -i $SPEC_FILE \
                -c $CONFIG_FILE \
                -g $LANG \
                -o $OUT_DIR/$LANG
        fi
    else
        if [ -d $LANG_TEMPLATE ]; then
            openapi-generator generate \
                --template-dir $LANG_TEMPLATE \
                --group-id $GROUP_ID \
                --artifact-version $API_VERSION \
                -i $SPEC_FILE \
                -g $LANG \
                -o $OUT_DIR/$LANG
        else
            openapi-generator generate \
                --group-id $GROUP_ID \
                --artifact-version $API_VERSION \
                -i $SPEC_FILE \
                -g $LANG \
                -o $OUT_DIR/$LANG
        fi
    fi

    echo
}

function compress()
{
    local SDK_OUT_FILE="$CUR/sdk_out/sdk-all-$API_VERSION.tar.gz"
    tar -C $OUT_DIR -czf $SDK_OUT_FILE .
}

function generateRestDoc()
{
    local DEST=$OUT_DIR/rest-doc

    if [ ! -d "$DEST" ]; then
     mkdir -p $DEST
    fi

    info "Generating ReDoc HTML documentation..."
    redoc-cli bundle $SPEC_FILE \
        -o $DEST/index.html \
        --title "喜马把脉平台 API V2" \
        --options.theme.colors.main=#00b2a5 \
        --options.pathInMiddlePanel
    echo
}

# Check output dir existing
if [ -d "$OUT_DIR" ]; then
  info "The target out directry exists."
  info "Removing directory: $OUT_DIR"
  rm -rd "$OUT_DIR"
  echo
fi
mkdir -p "$OUT_DIR"

###########################
# Generate documents
###########################

# Static HTML SDK doc
generateRestDoc
generate html "Generating static HTML documentation..."
generate html2 "Generating static HTML2 documentation..."
generate openapi "Generating OpenAPI JSON spec..."
generate openapi-yaml "Generating OpenAPI YAML spec..."

###########################
# Generate client SDKs
###########################

# Android Client SDK
generate android "Generating Android client SDK..."
# Build .jar package
pushd $OUT_DIR/android
mvn package
popd

# # Objective-C Client SDK
generate objc "Generating Objective-C client SDK..."

# # Swift 4 Client SDK
# generate swift4 "Generating Swift 4 client SDK..."

# # Kotlin Client SDK
# generate kotlin "Generating Kotlin client SDK..."

# # Java Client SDK
# generate java "Generating Java client SDK..."

# # Go Client SDK
# generate go "Generating Go client SDK..."

# Python Client SDK
generate python "Generating Python client SDK..."

echo
info "Finished to generate SDK codes."

# Compress SDK dir
echo
info "Compressing SDK dir..."
compress

# Done
echo
info "All is done."
