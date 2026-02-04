#!/usr/bin/env bash
# Check nvpkg create command: create empty package, then info.
set -e
SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
TMP_DIR="${SCRIPT_DIR}/../../tmp"
mkdir -p "$TMP_DIR"
NVPKG="${NVPKG:-${SCRIPT_DIR}/../nvpkg}"
PKG="${TMP_DIR}/check_create_$$.nvpk"

"$NVPKG" create "$PKG"
test -f "$PKG" || { echo "create: package file not created"; exit 1; }
"$NVPKG" info "$PKG" | head -1 | grep -q "Path:" || { echo "info: expected Path line"; exit 1; }
rm -f "$PKG"
echo "check_create: OK"
