#!/usr/bin/env bash
# Check nvpkg validate on an existing package.
set -e
SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
TMP_DIR="${SCRIPT_DIR}/../../tmp"
mkdir -p "$TMP_DIR"
NVPKG="${NVPKG:-${SCRIPT_DIR}/../nvpkg}"
PKG="${TMP_DIR}/check_validate_$$.nvpk"

"$NVPKG" create "$PKG"
"$NVPKG" validate "$PKG" | grep -q "OK" || { echo "validate: expected OK line"; exit 1; }
if "$NVPKG" validate /nonexistent/pkg.nvpk 2>/dev/null; then
	echo "validate nonexistent: should fail"
	exit 1
fi
rm -f "$PKG"
echo "check_validate: OK"
