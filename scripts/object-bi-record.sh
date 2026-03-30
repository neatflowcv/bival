#!/usr/bin/env bash

set -euo pipefail

usage() {
  echo "Usage: $0 <container-name> <bucket-name> <object-name>" >&2
}

fail() {
  echo "Error: $*" >&2
  exit 1
}

require_cmd() {
  local cmd="$1"

  if ! command -v "$cmd" >/dev/null 2>&1; then
    fail "required command not found: $cmd"
  fi
}

run_radosgw_admin() {
  podman exec -i "$container_name" radosgw-admin "$@"
}

run_json_query() {
  local label="$1"
  shift

  local output
  if ! output=$(run_radosgw_admin "$@" 2>&1); then
    fail "$label failed: $output"
  fi

  printf '%s\n' "$output"
}

extract_json_field() {
  local label="$1"
  local filter="$2"
  local json_input="$3"

  local value
  if ! value=$(printf '%s\n' "$json_input" | jq -er "$filter" 2>/dev/null); then
    fail "missing or invalid field for $label"
  fi

  printf '%s\n' "$value"
}

if [[ $# -ne 3 ]]; then
  usage
  exit 1
fi

container_name="$1"
bucket_name="$2"
object_name="$3"

require_cmd podman
require_cmd jq
require_cmd python3

bucket_stats_json=$(run_json_query "bucket stats" bucket stats --bucket="$bucket_name")
bucket_layout_json=$(run_json_query "bucket layout" bucket layout --bucket="$bucket_name")
zone_json=$(run_json_query "zone get" zone get)

bucket_id=$(extract_json_field "bucket id" '.id | strings | select(length > 0)' "$bucket_stats_json")
num_shards=$(extract_json_field "num_shards" '.num_shards | tonumber' "$bucket_stats_json")
gen=$(extract_json_field "bucket layout generation" '.layout.current_index.gen | tonumber' "$bucket_layout_json")
index_pool=$(extract_json_field "index pool" '.placement_pools[0].val.index_pool | strings | select(length > 0)' "$zone_json")

shard_json=$(run_json_query "bucket object shard" bucket object shard --object="$object_name" --num-shards "$num_shards")
shard=$(extract_json_field "shard" '.shard | tonumber' "$shard_json")

bi_list_file=$(mktemp)
cleanup() {
  rm -f "$bi_list_file"
}
trap cleanup EXIT

if ! run_radosgw_admin bi list --bucket="$bucket_name" --object="$object_name" --shard-id="$shard" >"$bi_list_file"; then
  fail "bi list failed for bucket=$bucket_name object=$object_name shard=$shard"
fi

if [[ "$gen" == "0" ]]; then
  index_object=".dir.${bucket_id}.${shard}"
else
  index_object=".dir.${bucket_id}.${gen}.${shard}"
fi

echo "Index Info"
echo "container_name=$container_name"
echo "bucket=$bucket_name"
echo "object=$object_name"
echo "bucket_id=$bucket_id"
echo "gen=$gen"
echo "shard=$shard"
echo "INDEX_POOL=\"$index_pool\""
echo "INDEX_OBJECT=\"$index_object\""
echo
echo "Record Summary"

record_summary="$(
  python3 - "$bi_list_file" <<'PY'
import sys


class Parser:
    def __init__(self, data: bytes) -> None:
        self.data = data
        self.pos = 0

    def error(self, message: str) -> None:
        raise ValueError(f"{message} at byte {self.pos}")

    def peek(self) -> int | None:
        if self.pos >= len(self.data):
            return None
        return self.data[self.pos]

    def consume(self, expected: int | None = None) -> int:
        if self.pos >= len(self.data):
            self.error("unexpected end of input")
        value = self.data[self.pos]
        if expected is not None and value != expected:
            self.error(f"expected byte {expected:#x}, found {value:#x}")
        self.pos += 1
        return value

    def skip_ws(self) -> None:
        while self.pos < len(self.data) and self.data[self.pos] in b" \t\r\n":
            self.pos += 1

    def parse_string_bytes(self) -> bytes:
        self.consume(ord('"'))
        out = bytearray()

        while True:
            if self.pos >= len(self.data):
                self.error("unterminated string")

            ch = self.consume()
            if ch == ord('"'):
                return bytes(out)

            if ch != ord("\\"):
                out.append(ch)
                continue

            esc = self.consume()
            if esc in (ord('"'), ord("\\"), ord("/")):
                out.append(esc)
                continue
            if esc == ord("b"):
                out.append(0x08)
                continue
            if esc == ord("f"):
                out.append(0x0C)
                continue
            if esc == ord("n"):
                out.append(0x0A)
                continue
            if esc == ord("r"):
                out.append(0x0D)
                continue
            if esc == ord("t"):
                out.append(0x09)
                continue
            if esc == ord("u"):
                out.extend(self.parse_unicode_escape())
                continue

            self.error(f"unsupported escape sequence \\{chr(esc)}")

    def parse_unicode_escape(self) -> bytes:
        codepoint = self.parse_hex_quad()

        if 0xD800 <= codepoint <= 0xDBFF:
            if self.pos + 6 > len(self.data) or self.data[self.pos:self.pos + 2] != b"\\u":
                self.error("missing low surrogate")
            self.pos += 2
            low = self.parse_hex_quad()
            if not 0xDC00 <= low <= 0xDFFF:
                self.error("invalid low surrogate")
            codepoint = 0x10000 + ((codepoint - 0xD800) << 10) + (low - 0xDC00)
        elif 0xDC00 <= codepoint <= 0xDFFF:
            self.error("unexpected low surrogate")

        return chr(codepoint).encode("utf-8")

    def parse_hex_quad(self) -> int:
        if self.pos + 4 > len(self.data):
          self.error("truncated unicode escape")
        raw = self.data[self.pos:self.pos + 4]
        self.pos += 4
        try:
            return int(raw.decode("ascii"), 16)
        except ValueError as exc:
            raise ValueError(f"invalid unicode escape at byte {self.pos - 4}") from exc

    def skip_primitive(self) -> None:
        start = self.pos
        while self.pos < len(self.data) and self.data[self.pos] not in b" \t\r\n,]}":
            self.pos += 1
        if self.pos == start:
            self.error("expected primitive")

    def skip_array(self) -> None:
        self.consume(ord("["))
        self.skip_ws()
        if self.peek() == ord("]"):
            self.consume(ord("]"))
            return

        while True:
            self.skip_value()
            self.skip_ws()
            ch = self.peek()
            if ch == ord(","):
                self.consume(ord(","))
                self.skip_ws()
                continue
            if ch == ord("]"):
                self.consume(ord("]"))
                return
            self.error("expected ',' or ']' in array")

    def skip_object(self) -> None:
        self.consume(ord("{"))
        self.skip_ws()
        if self.peek() == ord("}"):
            self.consume(ord("}"))
            return

        while True:
            self.parse_string_bytes()
            self.skip_ws()
            self.consume(ord(":"))
            self.skip_ws()
            self.skip_value()
            self.skip_ws()
            ch = self.peek()
            if ch == ord(","):
                self.consume(ord(","))
                self.skip_ws()
                continue
            if ch == ord("}"):
                self.consume(ord("}"))
                return
            self.error("expected ',' or '}' in object")

    def skip_value(self) -> None:
        self.skip_ws()
        ch = self.peek()
        if ch == ord('"'):
            self.parse_string_bytes()
            return
        if ch == ord("{"):
            self.skip_object()
            return
        if ch == ord("["):
            self.skip_array()
            return
        self.skip_primitive()

    def parse_record_object(self) -> tuple[str | None, bytes | None]:
        record_type = None
        idx_bytes = None

        self.consume(ord("{"))
        self.skip_ws()
        if self.peek() == ord("}"):
            self.consume(ord("}"))
            return record_type, idx_bytes

        while True:
            key_bytes = self.parse_string_bytes()
            try:
                key = key_bytes.decode("ascii")
            except UnicodeDecodeError as exc:
                raise ValueError(f"non-ascii object key at byte {self.pos}") from exc

            self.skip_ws()
            self.consume(ord(":"))
            self.skip_ws()

            if key == "type":
                value_bytes = self.parse_string_bytes()
                record_type = value_bytes.decode("ascii", "replace")
            elif key == "idx":
                idx_bytes = self.parse_string_bytes()
            else:
                self.skip_value()

            self.skip_ws()
            ch = self.peek()
            if ch == ord(","):
                self.consume(ord(","))
                self.skip_ws()
                continue
            if ch == ord("}"):
                self.consume(ord("}"))
                return record_type, idx_bytes
            self.error("expected ',' or '}' in record object")

    def parse_records(self) -> list[tuple[str, bytes]]:
        records: list[tuple[str, bytes]] = []

        self.skip_ws()
        self.consume(ord("["))
        self.skip_ws()
        if self.peek() == ord("]"):
            self.consume(ord("]"))
            self.skip_ws()
            if self.pos != len(self.data):
                self.error("unexpected trailing data")
            return records

        while True:
            record_type, idx_bytes = self.parse_record_object()
            if record_type is None or idx_bytes is None:
                self.error("record is missing type or idx")
            records.append((record_type, idx_bytes))

            self.skip_ws()
            ch = self.peek()
            if ch == ord(","):
                self.consume(ord(","))
                self.skip_ws()
                continue
            if ch == ord("]"):
                self.consume(ord("]"))
                self.skip_ws()
                if self.pos != len(self.data):
                    self.error("unexpected trailing data")
                return records
            self.error("expected ',' or ']' after record")


def escape_bytes(raw: bytes) -> str:
    parts = []
    for byte in raw:
        if 0x20 <= byte <= 0x7E:
            parts.append(chr(byte))
        else:
            parts.append(f"\\x{byte:02x}")
    return "".join(parts)


with open(sys.argv[1], "rb") as f:
    data = f.read()

parser = Parser(data)
records = parser.parse_records()
for record_type, idx_bytes in records:
    print(f'IDX="{escape_bytes(idx_bytes)}"')
PY
)" || fail "failed to parse bi list output"

if [[ -z "$record_summary" ]]; then
  echo "(no records)"
else
  printf '%s\n' "$record_summary"
fi

echo
cat "$bi_list_file"
