#!/usr/bin/env bash

set -euo pipefail

usage() {
  echo "Usage: $0 <cephconf-dir>" >&2
}

if [[ $# -ne 1 ]]; then
  usage
  exit 1
fi

cephconf_dir_input="$1"

if [[ ! -d "$cephconf_dir_input" ]]; then
  echo "Directory not found: $cephconf_dir_input" >&2
  exit 1
fi

cephconf_dir="$(cd "$cephconf_dir_input" && pwd -P)"
container_name="$(basename "$cephconf_dir")"

podman run -d \
  --name="$container_name" \
  -v "$cephconf_dir:/etc/ceph" \
  -it quay.io/ceph/ceph:v18.2.7 \
  sleep infinity

echo "Started container: $container_name"
echo "Mounted: $cephconf_dir -> /etc/ceph"
