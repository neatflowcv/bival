#!/usr/bin/env bash

set -euo pipefail

usage() {
  echo "Usage: $0 <container-name> <bucket-name>" >&2
}

if [[ $# -ne 2 ]]; then
  usage
  exit 1
fi

container_name="$1"
bucket_name="$2"
output_file="${bucket_name}.stats.json"

podman exec -i "$container_name" \
  radosgw-admin bucket stats \
  --bucket="$bucket_name" \
  >"$output_file"

echo "Saved bucket stats to $output_file"
