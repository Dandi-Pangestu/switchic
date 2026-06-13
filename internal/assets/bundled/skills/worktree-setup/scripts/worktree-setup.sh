#!/usr/bin/env bash
# worktree-setup.sh
# Provisions a new git worktree for a parallel agent session.
#
# Usage: worktree-setup.sh <JIRA_KEY> <JIRA_TITLE_SLUG> <REPO_ROOT>
#
# Arguments:
#   JIRA_KEY        - Jira ticket key, e.g. QC-99999
#   JIRA_TITLE_SLUG - Lowercase hyphenated title slug, max 6 words, e.g. fix-nil-message-send
#   REPO_ROOT       - Absolute path to the repository root
#
# Stdout (parsed by the invoking agent):
#   WORKTREE_PATH=<absolute path>
#   BRANCH_NAME=<branch name>

set -euo pipefail

JIRA_KEY="${1:?JIRA_KEY is required}"
JIRA_TITLE_SLUG="${2:?JIRA_TITLE_SLUG is required}"
REPO_ROOT="${3:?REPO_ROOT is required}"

BRANCH_NAME="feature/${JIRA_KEY}-${JIRA_TITLE_SLUG}"
WORKTREE_PATH="${REPO_ROOT}/../hub_core-session-${JIRA_KEY}"
REGISTRY="${REPO_ROOT}/.github/sessions/registry.json"
STARTED_AT=$(date -u +"%Y-%m-%dT%H:%M:%SZ")

if ! command -v jq &>/dev/null; then
  echo "ERROR: jq is required but not installed. Install via: brew install jq" >&2
  exit 1
fi

echo "Pulling latest production..."
git -C "$REPO_ROOT" pull origin production

echo "Creating worktree at ${WORKTREE_PATH} on branch ${BRANCH_NAME}..."
git -C "$REPO_ROOT" worktree add -b "$BRANCH_NAME" "$WORKTREE_PATH" production

ABSOLUTE_WORKTREE_PATH="$(cd "$WORKTREE_PATH" && pwd)"

echo "Registering session in registry.json..."
NEW_ENTRY=$(jq -n \
  --arg jira_key "$JIRA_KEY" \
  --arg branch "$BRANCH_NAME" \
  --arg worktree_path "$ABSOLUTE_WORKTREE_PATH" \
  --arg started_at "$STARTED_AT" \
  '{jira_key: $jira_key, branch: $branch, worktree_path: $worktree_path, status: "active", started_at: $started_at}')

jq --argjson entry "$NEW_ENTRY" '. += [$entry]' "$REGISTRY" > "${REGISTRY}.tmp" \
  && mv "${REGISTRY}.tmp" "$REGISTRY"

echo "WORKTREE_PATH=${ABSOLUTE_WORKTREE_PATH}"
echo "BRANCH_NAME=${BRANCH_NAME}"
