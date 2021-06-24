#!/usr/bin/env bash

set -e

owner=mzbaulhaque
repository_name=gomage
tag=$(cat VERSION)
release_id=$(curl -sS -u $owner:$GH_ACCESS_TOKEN -H "Accept: application/vnd.github.v3+json" https://api.github.com/repos/$owner/$repository_name/releases/tags/$tag | jq -r '.id')

# Upload release assets
for filename in dist/*.tar.gz; do
  GH_ASSET_UPLOAD_URL="https://uploads.github.com/repos/$owner/$repository_name/releases/$release_id/assets?name=$(basename $filename)"
  curl -o /dev/null -s -w "%{http_code}\n" -u $owner:$GH_ACCESS_TOKEN -X POST --data-binary @"$filename" -H "Accept: application/vnd.github.v3+json" -H "Content-Type: application/gzip" $GH_ASSET_UPLOAD_URL
done
