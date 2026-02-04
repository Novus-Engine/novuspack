#!/usr/bin/env bash
# Run all nvpkg functionality checks. Uses tmp at repo root.
set -e
SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
NVPKG="${NVPKG:-${SCRIPT_DIR}/../nvpkg}"
if ! test -x "$NVPKG"; then
	echo "Build nvpkg first: cd $(dirname "$SCRIPT_DIR") && go build -o nvpkg ."
	exit 1
fi
"$SCRIPT_DIR/check_create.sh"
"$SCRIPT_DIR/check_header_info.sh"
"$SCRIPT_DIR/check_validate.sh"
"$SCRIPT_DIR/check_add_list_read.sh"
"$SCRIPT_DIR/check_remove.sh"
echo "All checks passed."
