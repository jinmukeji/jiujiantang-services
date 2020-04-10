#!/usr/bin/env bash
set -e
set -u
set -o pipefail

CUR=`dirname $0`

DB_FILE=${CUR}/GeoLite2-Country.mmdb.gz

echo "Updating GeoLite database..."
wget https://geolite.maxmind.com/download/geoip/database/GeoLite2-Country.mmdb.gz -O ${DB_FILE}.download
mv -f ${DB_FILE}.download ${DB_FILE}

echo "GeoLite database has been updated."
