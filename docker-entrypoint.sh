#!/usr/bin/env sh

set -e
envsubst < /app/config.tmpl > /app/config.json
jq '.' /app/config.json > /dev/null
exec ./getir-case