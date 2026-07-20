#!/usr/bin/env bash
# worktree-cleanup.sh
# Removes a git worktree and deregisters the session from the session registry.
#
# Usage: worktree-cleanup.sh <JIRA_KEY> <REPO_ROOT> [SESSION_ROOT] [force]
#
# Arguments:
#   JIRA_KEY     - Jira ticket key of the session to remove, e.g. QC-99999
#   REPO_ROOT    - Absolute path to the primary repository root (fallback when the
#                  registry entry has no repo_root field)
#   SESSION_ROOT - Where session artifacts live (default: REPO_ROOT). In workspace
#                  mode, pass the workspace root.
#   force        - Pass the literal string "force" to remove the worktree even if it
#                  has uncommitted changes (they will be discarded).
#
# This script only removes the worktree directory and the registry entry. It never
# merges, pushes, or deletes branches — the feature branch and its commits stay in
# the local repository.

set -euo pipefail

JIRA_KEY="${1:?JIRA_KEY is required}"
REPO_ROOT="${2:?REPO_ROOT is required}"
SESSION_ROOT="${3:-$REPO_ROOT}"
FORCE="${4:-}"

REGISTRY="${SESSION_ROOT}/sessions/registry.json"

if ! command -v jq &>/dev/null; then
  echo "ERROR: jq is required but not installed. Install via: brew install jq" >&2
  exit 1
fi

if [[ ! -f "$REGISTRY" ]]; then
  echo "ERROR: Registry not found at ${REGISTRY}" >&2
  exit 1
fi

WORKTREE_PATH=$(jq -r --arg key "$JIRA_KEY" \
  '[.[] | select(.jira_key == $key) | .worktree_path][0] // empty' "$REGISTRY")

if [[ -z "$WORKTREE_PATH" ]]; then
  echo "ERROR: No session found for ${JIRA_KEY} in registry." >&2
  exit 1
fi

# The repo the worktree belongs to: prefer the registry entry, fall back to the arg.
ENTRY_REPO_ROOT=$(jq -r --arg key "$JIRA_KEY" \
  '[.[] | select(.jira_key == $key) | .repo_root][0] // empty' "$REGISTRY")
GIT_REPO_ROOT="${ENTRY_REPO_ROOT:-$REPO_ROOT}"

ABSOLUTE_GIT_REPO_ROOT="$(cd "$GIT_REPO_ROOT" && pwd -P)"

if [[ "$(cd "$WORKTREE_PATH" 2>/dev/null && pwd -P || echo "$WORKTREE_PATH")" == "$ABSOLUTE_GIT_REPO_ROOT" ]]; then
  # Primary session: the "worktree" is the repository root itself. Never remove it —
  # only deregister the session.
  echo "Session ${JIRA_KEY} is the primary session (worktree is the repo root). Skipping worktree removal."
else
  echo "Removing worktree at ${WORKTREE_PATH}..."
  if [[ -d "$WORKTREE_PATH" ]]; then
    if [[ "$FORCE" != "force" ]] && [[ -n "$(git -C "$WORKTREE_PATH" status --porcelain 2>/dev/null)" ]]; then
      echo "ERROR: Worktree at ${WORKTREE_PATH} has uncommitted changes." >&2
      echo "Commit them first, or re-run with 'force' as the 4th argument to discard them." >&2
      exit 1
    fi
    BRANCH=$(git -C "$WORKTREE_PATH" rev-parse --abbrev-ref HEAD 2>/dev/null || echo "")
    git -C "$GIT_REPO_ROOT" worktree remove "$WORKTREE_PATH" --force
    echo "Worktree removed."
    if [[ -n "$BRANCH" ]]; then
      echo "Branch ${BRANCH} is preserved locally (nothing was merged or pushed)."
    fi
  else
    echo "Worktree directory not found on disk — pruning stale entry."
    git -C "$GIT_REPO_ROOT" worktree prune
  fi
fi

echo "Removing registry entry for ${JIRA_KEY}..."
jq --arg key "$JIRA_KEY" \
  '[.[] | select(.jira_key != $key)]' \
  "$REGISTRY" > "${REGISTRY}.tmp" \
  && mv "${REGISTRY}.tmp" "$REGISTRY"

echo "Session ${JIRA_KEY} cleaned up successfully."
