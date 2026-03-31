#!/usr/bin/env bash

set -euo pipefail

usage() {
  echo "Usage: $0 <container-name> <bucket-name> [shard-id]" >&2
}

if [[ $# -lt 2 || $# -gt 3 ]]; then
  usage
  exit 1
fi

container_name="$1"
bucket_name="$2"
shard_id="${3:-}"
output_file="${bucket_name}.json"

radosgw_admin_args=(
  radosgw-admin bi list
  --bucket="$bucket_name"
)

if [[ -n "$shard_id" ]]; then
  radosgw_admin_args+=(--shard-id "$shard_id")
  output_file="${bucket_name}.${shard_id}.json"
fi

podman exec -i "$container_name" \
  "${radosgw_admin_args[@]}" \
  >"$output_file"

echo "Saved bi list to $output_file"
