#!/usr/bin/env bash

set -euo pipefail

usage() {
  echo "Usage: $0 <input-file>" >&2
}

if [[ $# -ne 1 ]]; then
  usage
  exit 1
fi

input_file="$1"

if [[ ! -f "$input_file" ]]; then
  echo "Input file not found: $input_file" >&2
  exit 1
fi

input_dir=$(dirname "$input_file")
input_base=$(basename "$input_file")
input_stem="${input_base%.*}"

sorted_file="${input_dir}/${input_stem}.sorted.json"
analyze_output_file="${input_dir}/${input_stem}.out"

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
