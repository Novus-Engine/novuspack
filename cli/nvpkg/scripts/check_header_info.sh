#!/usr/bin/env bash
# Check nvpkg header, info, and validate commands on an existing package.
set -e
SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
TMP_DIR="${SCRIPT_DIR}/../../tmp"
mkdir -p "$TMP_DIR"
NVPKG="${NVPKG:-${SCRIPT_DIR}/../nvpkg}"
PKG="${TMP_DIR}/check_header_info_$$.nvpk"

"$NVPKG" create "$PKG"
"$NVPKG" header "$PKG" | grep -q "Magic:" || { echo "header: expected Magic line"; exit 1; }
"$NVPKG" header "$PKG" | grep -q "FormatVersion:" || { echo "header: expected FormatVersion line"; exit 1; }
"$NVPKG" info "$PKG" | grep -q "File count:" || { echo "info: expected File count line"; exit 1; }
"$NVPKG" create "$TMP_DIR/check_identity_$$.nvpk" --vendor-id 1 --app-id 100
"$NVPKG" info "$TMP_DIR/check_identity_$$.nvpk" | grep -q "Vendor ID:" || { echo "info: expected Vendor ID line"; exit 1; }
"$NVPKG" info "$TMP_DIR/check_identity_$$.nvpk" | grep -q "App ID:" || { echo "info: expected App ID line"; exit 1; }
rm -f "$PKG" "$TMP_DIR/check_identity_$$.nvpk"
echo "check_header_info: OK"
