#!/usr/bin/env bash
# worktree-setup.sh
# Provisions a new git worktree for a parallel agent session.
#
# Usage: worktree-setup.sh <JIRA_KEY> <JIRA_TITLE_SLUG> <REPO_ROOT> [BASE_BRANCH] [SESSION_ROOT]
#
# Arguments:
#   JIRA_KEY        - Jira ticket key, e.g. QC-99999
#   JIRA_TITLE_SLUG - Lowercase hyphenated title slug, max 6 words, e.g. fix-nil-message-send
#   REPO_ROOT       - Absolute path to the target repository root
#   BASE_BRANCH     - Branch to base the feature branch on (default: production)
#   SESSION_ROOT    - Where session artifacts live (default: REPO_ROOT).
#                     In workspace mode, pass the workspace root so all sessions
#                     register in one registry.
#
# Stdout (parsed by the invoking agent):
#   WORKTREE_PATH=<absolute path>
#   BRANCH_NAME=<branch name>
#   BASE_BRANCH=<base branch used>

set -euo pipefail

JIRA_KEY="${1:?JIRA_KEY is required}"
JIRA_TITLE_SLUG="${2:?JIRA_TITLE_SLUG is required}"
REPO_ROOT="${3:?REPO_ROOT is required}"
BASE_BRANCH="${4:-production}"
SESSION_ROOT="${5:-$REPO_ROOT}"

# Resolve symlinks (workspace repos may be symlinked under repos/) so the
# worktree is created next to the real repository, not inside the workspace.
REPO_ROOT="$(cd "$REPO_ROOT" && pwd -P)"
SESSION_ROOT="$(cd "$SESSION_ROOT" && pwd)"

REPO_NAME="$(basename "$REPO_ROOT")"
BRANCH_NAME="feature/${JIRA_KEY}-${JIRA_TITLE_SLUG}"
WORKTREE_PATH="${REPO_ROOT}/../${REPO_NAME}-session-${JIRA_KEY}"
REGISTRY="${SESSION_ROOT}/sessions/registry.json"
STARTED_AT=$(date -u +"%Y-%m-%dT%H:%M:%SZ")

if ! command -v jq &>/dev/null; then
  echo "ERROR: jq is required but not installed. Install via: brew install jq" >&2
  exit 1
fi

if [[ ! -f "$REGISTRY" ]]; then
  mkdir -p "$(dirname "$REGISTRY")"
  echo "[]" > "$REGISTRY"
fi

echo "Pulling latest ${BASE_BRANCH}..."
git -C "$REPO_ROOT" fetch origin "$BASE_BRANCH"

echo "Creating worktree at ${WORKTREE_PATH} on branch ${BRANCH_NAME}..."
git -C "$REPO_ROOT" worktree add -b "$BRANCH_NAME" "$WORKTREE_PATH" "origin/${BASE_BRANCH}"

ABSOLUTE_WORKTREE_PATH="$(cd "$WORKTREE_PATH" && pwd)"

echo "Registering session in registry.json..."
NEW_ENTRY=$(jq -n \
  --arg jira_key "$JIRA_KEY" \
  --arg branch "$BRANCH_NAME" \
  --arg worktree_path "$ABSOLUTE_WORKTREE_PATH" \
  --arg repo_root "$REPO_ROOT" \
  --arg base_branch "$BASE_BRANCH" \
  --arg started_at "$STARTED_AT" \
  '{jira_key: $jira_key, branch: $branch, worktree_path: $worktree_path, repo_root: $repo_root, base_branch: $base_branch, status: "active", started_at: $started_at}')

jq --argjson entry "$NEW_ENTRY" '. += [$entry]' "$REGISTRY" > "${REGISTRY}.tmp" \
  && mv "${REGISTRY}.tmp" "$REGISTRY"

echo "WORKTREE_PATH=${ABSOLUTE_WORKTREE_PATH}"
echo "BRANCH_NAME=${BRANCH_NAME}"
echo "BASE_BRANCH=${BASE_BRANCH}"
