#!/usr/bin/env bash
# Apply every SQL migration in lexical order.
# Usage: scripts/apply-migrations.sh <database-url>
set -euo pipefail

url="${1:-${GOWRITE_TEST_DATABASE_URL:-}}"
if [[ -z "$url" ]]; then
  echo "usage: $0 <database-url> (or set GOWRITE_TEST_DATABASE_URL)" >&2
  exit 2
fi

root="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
migrations_dir="$root/impl/db/migrations"

files=$(find "$migrations_dir" -maxdepth 1 -type f -name '*.sql' -print | LC_ALL=C sort)
if [[ -z "$files" ]]; then
  echo "no migrations found in $migrations_dir" >&2
  exit 1
fi

while IFS= read -r file; do
  echo "applying $(basename "$file")"
  psql "$url" -v ON_ERROR_STOP=1 -q -f "$file"
done <<< "$files"
