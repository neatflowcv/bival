#!/usr/bin/env bash

set -euo pipefail

usage() {
  echo "Usage: $0 <container-name> <bucket-name> <object-name>" >&2
}

if [[ $# -ne 3 ]]; then
  usage
  exit 1
fi

container_name="$1"
bucket_name="$2"
object_name="$3"

if podman exec -i "$container_name" \
  radosgw-admin object stat \
  --bucket="$bucket_name" \
  --object="$object_name" \
  >/dev/null 2>&1; then
  echo "Object exists: bucket=$bucket_name object=$object_name"
  exit 0
fi

echo "Object not found: bucket=$bucket_name object=$object_name" >&2
exit 1
