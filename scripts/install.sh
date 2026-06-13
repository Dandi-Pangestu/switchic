#!/usr/bin/env bash
# One-shot installer for switchic-flow.
#
# Usage:
#   ./scripts/install.sh             # install to /usr/local/bin (may sudo)
#   ./scripts/install.sh --user      # install to ~/.local/bin (no sudo)
#   PREFIX=/opt/switchic ./scripts/install.sh
#
# Requires: go >= 1.21
set -euo pipefail

MODE="system"
for arg in "$@"; do
  case "$arg" in
    --user) MODE="user" ;;
    -h|--help)
      grep '^#' "$0" | sed 's/^# \{0,1\}//'
      exit 0
      ;;
    *) echo "unknown argument: $arg" >&2; exit 2 ;;
  esac
done

SCRIPT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
REPO_ROOT="$( cd "$SCRIPT_DIR/.." && pwd )"

cd "$REPO_ROOT"

if ! command -v go >/dev/null 2>&1; then
  echo "Error: 'go' is not installed or not on PATH." >&2
  echo "Install Go from https://go.dev/dl/ and re-run." >&2
  exit 1
fi

if [ "$MODE" = "user" ]; then
  exec make user-install
else
  exec make install
fi
