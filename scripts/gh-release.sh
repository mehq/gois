#!/usr/bin/env bash

set -e

version=$1
tag="v${version}"
release_notes=$(awk -v ver="[${version}]" ' /^## / { if (p) { exit }; if ($2 == ver) { p=1; next} } p' CHANGELOG.md | head -n -1)
release_name=$(awk -v var="## \[$version\]" '$0 ~ var {print $2 " " $3 " " $4}' CHANGELOG.md)
owner=mzbaulhaque
repo=gomage
create_payload="{\"name\":\"${release_name}\",\"tag_name\":\"${tag}\",\"body\":\"${release_notes}\"}"
release_id=$(curl -sS -u $owner:$GH_ACCESS_TOKEN -X POST -H "Accept: application/vnd.github.v3+json" https://api.github.com/repos/$owner/$repo/releases -d "$create_payload" | jq -r '.id')

# Upload release assets
for filename in dist/*.tar.gz; do
  GH_ASSET_UPLOAD_URL="https://uploads.github.com/repos/$owner/$repo/releases/$release_id/assets?name=$(basename $filename)"
  curl -o /dev/null -s -w "%{http_code}\n" -u $owner:$GH_ACCESS_TOKEN -X POST --data-binary @"$filename" -H "Accept: application/vnd.github.v3+json" -H "Content-Type: application/gzip" $GH_ASSET_UPLOAD_URL
done
