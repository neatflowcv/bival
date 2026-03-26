#!/usr/bin/env bash

set -euo pipefail

usage() {
  echo "Usage: $0 <bucket-name>" >&2
}

if [[ $# -ne 1 ]]; then
  usage
  exit 1
fi

bucket_name="$1"
input_file="${bucket_name}.json"
sorted_file="${bucket_name}.sorted.json"
analyze_output_file="${bucket_name}.out"

run_step() {
  local label="$1"
  shift

  local started_at elapsed
  started_at=$(date +%s)

  echo "[START] $label"
  "$@"
  elapsed=$(( $(date +%s) - started_at ))
  echo "[DONE ] $label (${elapsed}s)"
}

total_started_at=$(date +%s)

echo "Bucket: $bucket_name"
echo "Input : $input_file"
echo "Sorted: $sorted_file"
echo "Output: $analyze_output_file"

run_step "sort $input_file -> $sorted_file" \
  go run ./cmd/bival sort "$input_file" "$sorted_file"

run_step "analyze $sorted_file -> $analyze_output_file" \
  bash -lc 'go run ./cmd/bival analyze "$1" >"$2" 2>&1' _ "$sorted_file" "$analyze_output_file"

echo "Lines : $(wc -l < "$analyze_output_file")"

total_elapsed=$(( $(date +%s) - total_started_at ))
echo "[DONE ] total (${total_elapsed}s)"
