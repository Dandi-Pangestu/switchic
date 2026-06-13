#!/usr/bin/env bash
# worktree-cleanup.sh
# Removes a git worktree and deregisters the session from the session registry.
#
# Usage: worktree-cleanup.sh <JIRA_KEY> <REPO_ROOT>
#
# Arguments:
#   JIRA_KEY   - Jira ticket key of the session to remove, e.g. QC-99999
#   REPO_ROOT  - Absolute path to the primary repository root

set -euo pipefail

JIRA_KEY="${1:?JIRA_KEY is required}"
REPO_ROOT="${2:?REPO_ROOT is required}"

REGISTRY="${REPO_ROOT}/.github/sessions/registry.json"

if ! command -v jq &>/dev/null; then
  echo "ERROR: jq is required but not installed. Install via: brew install jq" >&2
  exit 1
fi

if [[ ! -f "$REGISTRY" ]]; then
  echo "ERROR: Registry not found at ${REGISTRY}" >&2
  exit 1
fi

WORKTREE_PATH=$(jq -r --arg key "$JIRA_KEY" \
  '.[] | select(.jira_key == $key) | .worktree_path' "$REGISTRY")

if [[ -z "$WORKTREE_PATH" ]]; then
  echo "ERROR: No session found for ${JIRA_KEY} in registry." >&2
  exit 1
fi

echo "Removing worktree at ${WORKTREE_PATH}..."
if [[ -d "$WORKTREE_PATH" ]]; then
  git -C "$REPO_ROOT" worktree remove "$WORKTREE_PATH" --force
  echo "Worktree removed."
else
  echo "Worktree directory not found on disk — pruning stale entry."
  git -C "$REPO_ROOT" worktree prune
fi

echo "Removing registry entry for ${JIRA_KEY}..."
jq --arg key "$JIRA_KEY" \
  '[.[] | select(.jira_key != $key)]' \
  "$REGISTRY" > "${REGISTRY}.tmp" \
  && mv "${REGISTRY}.tmp" "$REGISTRY"

echo "Session ${JIRA_KEY} cleaned up successfully."
