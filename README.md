# switchic

> One config. Any AI coding assistant.

[![Go version](https://img.shields.io/badge/go-1.21%2B-blue)](https://go.dev/dl/)
[![License: MIT](https://img.shields.io/badge/license-MIT-green)](LICENSE)
[![Release](https://img.shields.io/github/v/release/Dandi-Pangestu/switchic)](https://github.com/Dandi-Pangestu/switchic/releases)

**switchic** is a CLI tool that manages the agentic context loaded into AI coding assistants. Instead of manually maintaining `CLAUDE.md`, Cursor rules, or Copilot instructions for each project, you define your agents, skills, and rules once — then let switchic generate the right files for whatever platform you're using.

```
switchic init                  # set up a project
switchic switch claude         # generate CLAUDE.md + .claude/agents/* + .claude/rules/*
switchic switch github-copilot # generate .github/copilot-instructions.md + agents + skills
switchic switch kiro           # generate AGENTS.md + .kiro/agents/* + .kiro/steering/*
switchic status                # see active components and token cost
```

---

## Table of Contents

- [Why switchic](#why-switchic)
- [Core Concepts](#core-concepts)
- [AI Platform Support](#ai-platform-support)
- [Install](#install)
- [Quickstart](#quickstart)
  - [Generated files: commit or ignore?](#generated-files-commit-or-ignore)
  - [Speed up setup with built-in skills](#speed-up-setup-with-built-in-skills)
  - [Multi-repo workspace](#multi-repo-workspace)
- [Commands](#commands)
- [Bundled Library](#bundled-library)
- [Coding Workflow](#coding-workflow)
- [User-Defined Assets](#user-defined-assets)
- [Config Reference](#config-reference)
- [Cost Optimization](#cost-optimization)
- [Architecture](#architecture)
- [Roadmap](#roadmap)
- [Contributing](#contributing)
- [License](#license)

---

## Why switchic

When you work with AI coding assistants, you configure them through context files: system prompts, sub-agent definitions, coding rules, and reusable skill guides. The problem is that every platform has a different format, and these files tend to grow unbounded — loading everything into every session even when most of it is irrelevant to the current task.

switchic solves two things:

1. **One source of truth** — define your agents, skills, and rules in a shared format. switchic generates the platform-specific files on demand.
2. **Context discipline** — turn components on or off per-project. Only what the current task needs gets loaded into the session, keeping token costs low.

---

## Core Concepts

switchic organizes agentic context into four building blocks. Understanding these will make every command self-explanatory.

### Agent

A **sub-agent** is a specialist persona that an AI assistant can delegate work to. Each agent has a name, a description, a set of instructions, and a list of tools it is allowed to use.

When you run `switchic switch claude`, each active agent is written to `.claude/agents/<name>.md` — the format Claude Code uses to load sub-agents into a session.

**Example use case:** the `code-reviewer` agent is only activated when you're doing review work. When you're doing a quick exploration, it stays off.

### Skill

A **skill** is a structured prompt guide — a set of instructions the AI follows when performing a specific task. Skills are not full agent personas; they are reference documents that get embedded into the main context file (e.g. `CLAUDE.md`).

**Example use case:** the `implementation-plan` skill tells the AI exactly how to format a structured implementation plan for a Jira ticket, including required sections and examples.

### Rule

A **rule** is a coding standard or constraint the AI must follow. Rules are injected into the platform context and apply globally across the session.

Rules support directories — enabling `backend` activates all rules inside that directory (e.g. `backend/api` and `backend/database`).

**Example use case:** enable `golang` for Go projects, `typescript` for frontend projects.

### Workflow

A **workflow** is a named preset that automatically activates a bundle of agents and skills suited to a task type. Instead of manually toggling individual components, you set a workflow and get a ready-to-run context.

**Example use case:** the `coding` workflow activates a hierarchical agent pipeline for ticket-to-code development — a `coding-orchestrator` that routes each request to only the capabilities it needs (fetch requirements, write a plan, implement, review, or manage sessions), instead of always running a fixed sequence. See [Coding Workflow](#coding-workflow) for the full capability list and example prompts.

---

## AI Platform Support

| Platform | Status | Generated files |
|---|---|---|
| **Claude** (Claude Code) | Ready | `CLAUDE.md`, `.claude/agents/*.md`, `.claude/rules/*.md`, `.claude/skills/*/SKILL.md` |
| **GitHub Copilot** | Ready | `AGENTS.md`, `.github/copilot-instructions.md`, `.github/agents/*.agent.md`, `.github/instructions/*.instructions.md`, `.github/skills/*/SKILL.md` |
| **Kiro** | Ready | `AGENTS.md`, `.kiro/steering/project.md`, `.kiro/agents/*.json`, `.kiro/steering/*.md`, `.kiro/skills/*/SKILL.md` |
| Cursor | Coming soon | — |
| Codex CLI | Coming soon | — |
| Windsurf | Coming soon | — |

The architecture is built around a platform adapter interface — adding a new platform means implementing one `Generate(ctx)` method without touching the rest of the tool.

---

## Install

Go 1.21+ is required ([download](https://go.dev/dl/)).

### System-wide (recommended)

```bash
git clone https://github.com/Dandi-Pangestu/switchic
cd switchic
make install        # installs to /usr/local/bin (may prompt for sudo)
switchic version
```

### User-local (no sudo)

```bash
make user-install   # installs to ~/.local/bin
```

If `~/.local/bin` is not on your `PATH`, the installer prints the exact line to add to your shell rc:

```bash
export PATH="$HOME/.local/bin:$PATH"
```

### One-shot script

```bash
./scripts/install.sh           # system install
./scripts/install.sh --user    # user-local install
```

### Custom prefix

```bash
PREFIX=/opt/switchic make install
```

### Uninstall

```bash
make uninstall
```

---

## Quickstart

### Single repo

```bash
cd path/to/your/repo

# 1. Initialize the project config
switchic init

# 2. Populate project context (optional — see "Speed up setup" below)

# 3. Generate platform files — pick your AI assistant
switchic switch claude          # for Claude Code
switchic switch github-copilot  # for GitHub Copilot
switchic switch kiro            # for Kiro

# 4. Check what's active and the token cost
switchic status

# 5. Tune the context — disable what you don't need right now
switchic remove agent code-reviewer
switchic add skill commit-msg

# 6. Regenerate after any mutation
switchic switch claude   # or: switchic switch kiro / switchic switch github-copilot
```

> **Conflict resolution:** If `CLAUDE.md` or `AGENTS.md` already exists and was not generated by switchic, the tool writes `CLAUDE.local.md` or `AGENTS.local.md` instead — leaving your file untouched. Use `--replace` to force-overwrite.

After `switchic switch claude`, your repo will have:

```
CLAUDE.md                       # main context Claude reads on launch
.claude/agents/<name>.md        # one file per active agent
.claude/rules/<name>.md         # one file per active rule
.claude/skills/<name>/SKILL.md  # one directory per active skill
```

> If `CLAUDE.md` is user-written, `CLAUDE.local.md` is generated instead — Claude Code loads both natively.

After `switchic switch github-copilot`, your repo will have:

```
AGENTS.md                                          # full project context (primary source of truth)
.github/copilot-instructions.md                    # pointer to AGENTS.md via @AGENTS.md
.github/agents/<name>.agent.md                     # one file per active agent
.github/instructions/<name>.instructions.md        # one file per active rule (path-specific)
.github/skills/<name>/SKILL.md                     # one directory per active skill
```

> If `AGENTS.md` is user-written, `AGENTS.local.md` is generated instead and `.github/copilot-instructions.md` is updated to reference it.

After `switchic switch kiro`, your repo will have:

```
AGENTS.md                          # main context file (AGENTS.md standard)
.kiro/steering/project.md          # always-included steering pointer to AGENTS.md
.kiro/agents/<name>.json           # one file per active agent (Kiro JSON format)
.kiro/steering/<name>.md           # one file per active rule (always-included steering)
.kiro/skills/<name>/SKILL.md       # one directory per active skill
```

> If `AGENTS.md` is user-written, `AGENTS.local.md` is generated instead and `.kiro/steering/project.md` is updated to reference it.

### Generated files: commit or ignore?

`.switchic/` is the source of truth — `CLAUDE.md`, `AGENTS.md`, and all other generated files are build artifacts. Two valid approaches:

**Commit generated files** (recommended for most teams): commit `.switchic/` alongside the generated files. Teammates without switchic installed get the full AI context immediately — Claude Code, Copilot, and Kiro pick up the files directly. Only whoever runs `switch` needs the tool installed.

**Gitignore generated files**: version-control only `.switchic/`. Every developer runs `switchic switch <platform>` after clone to regenerate. Git history stays clean and diffs contain only config changes, not generated output.

In both models, `.switchic/` is what you review in PRs, share across the team, and treat as the canonical source.

### Speed up setup with built-in skills

switchic ships three skills that replace manual config and migration work.
Enable them once, then invoke them from your AI assistant.

#### `generate-context` — auto-fill project context from the codebase

Instead of hand-writing `description`, `stack`, `commands`, `structure`, `conventions`,
`dos`, `donts`, and `docs` in `.switchic/config.yaml`, let the AI scan your repo and fill
them in:

```bash
switchic add skill generate-context
switchic switch claude      # or kiro / github-copilot
# then in your AI assistant: /generate-context
```

The skill reads your manifest files (go.mod, package.json, Makefile, README, etc.),
generates each config field with the correct YAML format, and writes them into
`.switchic/config.yaml` — preserving your existing `platform`, `workflows`, and component
lists untouched.

#### `agent-import` — migrate an agent from another platform

Convert an existing agent file (Claude Code `.md`, GitHub Copilot `.agent.md`, Cursor, etc.)
into a switchic YAML agent definition:

```bash
switchic add skill agent-import
switchic switch claude
# then in your AI assistant: /agent-import path/to/agent.md [output-path]
```

#### `skill-import` — migrate a skill from another platform

Convert an existing `SKILL.md` from any platform into a switchic YAML skill definition:

```bash
switchic add skill skill-import
switchic switch claude
# then in your AI assistant: /skill-import path/to/SKILL.md [output-path]
```

Once imported, drop the output YAML into `.switchic/agents/` or `.switchic/skills/` and
enable it with `switchic add`.

---

### Multi-repo workspace

Workspace mode lets you manage AI context across multiple repos from a single manifest. The generated context file includes a structured summary of all repos — names, roles, and notes — so the AI understands the full system without ingesting raw code from every repo.

```bash
mkdir my-workspace && cd my-workspace

# 1. Create the workspace manifest
switchic workspace init --name billing-workspace --notes "Billing platform — API, web, contracts"

# 2. Register repos
switchic workspace add ../billing-api      --role backend  --notes "REST API, owns invoicing"
switchic workspace add ../billing-web      --role frontend --notes "Customer-facing dashboard (React)"
switchic workspace add ../billing-contracts --role contracts --notes "Shared OpenAPI specs"

# 3. Check what's registered
switchic workspace list

# 4. Generate platform files — workspace context is included automatically
switchic switch claude   # or: kiro / github-copilot
```

**Workspace-level component overrides** — `add` and `remove` write to the workspace manifest when run from a workspace, overriding the per-repo config:

```bash
switchic add agent code-reviewer
switchic remove skill jira-ticket-description
```

**Custom context file per repo** — if a repo's context lives at a non-standard path, point to it explicitly:

```bash
switchic workspace add ../billing-api --context-file ../billing-api/docs/CLAUDE.md
```

**Remove a repo:**

```bash
switchic workspace remove billing-contracts
```

**Browsing out-of-tree repos** — repos registered with an absolute path (living outside the workspace directory) get a symlink under `repos/<name>`, so `ls repos/` surfaces every repo in one place regardless of where it actually lives on disk. Repos registered with a relative path (the sibling-directory pattern above) are already reachable in-tree and don't get one. Symlinks are created/removed automatically by `workspace add`/`remove`; run `switchic workspace link` to resync them by hand (e.g. after editing the manifest directly or cloning the workspace onto a new machine):

```bash
switchic workspace add /Users/me/code/billing-legacy --role backend
ls repos/   # billing-legacy -> /Users/me/code/billing-legacy

switchic workspace link
```

---

## Commands

| Command | Purpose |
|---|---|
| `init` | Create `.switchic/config.yaml` in the current repo |
| `switch <platform>` | Regenerate platform files; auto-creates `.local` variant if the main context file is user-written |
| `switch <platform> --replace` | Force-overwrite the main context file even if it is user-written |
| `status` | Show active platform, workflows, components, and token cost |
| `add agent <name>` | Enable an agent |
| `add skill <name>` | Enable a skill |
| `add rule <name>` | Enable a rule |
| `remove agent <name>` | Disable an agent |
| `remove skill <name>` | Disable a skill |
| `remove rule <name>` | Disable a rule |
| `workspace init` | Create `switchic.workspace.yaml` in the current directory |
| `workspace init --name <n> --notes <n>` | Create workspace with a custom name and description |
| `workspace add <path>` | Register a repo; accepts `--role`, `--notes`, `--context-file` |
| `workspace remove <name>` | Unregister a repo from the workspace |
| `workspace list` | List registered repos and warn about any missing on disk |
| `workspace link` | Create/refresh `repos/<name>` symlinks for out-of-tree repos |
| `version` | Print the binary version |

---

## Bundled Library

These agents, skills, rules, and workflows ship inside the binary. Use `switchic add` / `switchic remove` to toggle them, or override any of them with a local file (see [User-Defined Assets](#user-defined-assets)).

### Agents

| Name | Description |
|---|---|
| `coding-orchestrator` | Entry point for the end-to-end coding workflow — orchestrates session bootstrap, planning, implementation, and review |
| `code-implementer` | Executes an approved implementation plan by writing code and updating tests |
| `code-reviewer` | Reviews a completed implementation for correctness, quality, security, and test coverage |
| `implementation-plan-writer` | Analyzes a Jira ticket and PR context to produce a detailed, structured implementation plan |
| `session-bootstraper` | Provisions the working environment for a new workflow session (worktree, branch, session registry) |
| `session-manager` | Manages parallel agent sessions — lists active sessions and cleans up completed ones |
| `jira-requirements-fetcher` | Fetches comprehensive Jira ticket context and metadata given a ticket key |
| `pr-details-fetcher` | Fetches pull request or merge request details given a PR/MR identifier |
| `example` | Template agent demonstrating every available YAML field — use as a starting point for custom agents |

### Skills

| Name | Description |
|---|---|
| `implementation-plan` | Guide for producing structured implementation plans for Jira tickets |
| `jira-ticket-description` | Guide for producing Jira ticket descriptions with consistent sections |
| `code-review-documentation` | Guide for producing structured code review summaries |
| `generate-context` | Scan the project and auto-fill description, stack, commands, structure, conventions, dos, donts, and docs into `.switchic/config.yaml` |
| `agent-import` | Convert an agent definition from any AI coding platform into switchic format |
| `skill-import` | Convert a skill definition from any AI coding platform into switchic format |
| `worktree-setup` | Provision a git worktree at a sibling path of the repository root, copying gitignored local config (`.env*`, `config/master.key`, credentials keys) so tests run immediately |
| `worktree-cleanup` | Remove a git worktree created by the `worktree-setup` skill |
| `session-state` | Defines the per-ticket `state.yaml` schema and phase lifecycle the coding-orchestrator uses to resume interrupted sessions |
| `commit-msg` | Generate a conventional commit message from staged changes |
| `example` | Template skill demonstrating every available YAML field — use as a starting point for custom skills |

### Rules

Rules are referenced by their key path (without `.yaml`). Enabling a directory key activates all rules inside it — e.g. `backend` enables both `backend/api` and `backend/database`.

| Key | Description |
|---|---|
| `golang` | Go coding standards (based on Effective Go) |
| `typescript` | TypeScript coding standards (based on Google TypeScript Style Guide) |
| `backend/api` | REST API design standards (based on Azure API Design Guidelines) |
| `backend/database` | Database coding and querying best practices |

### Workflows

| Name | Description | Activates |
|---|---|---|
| `coding` | Drives software engineering work from a Jira ticket to reviewed code via a hierarchical, goal-routed orchestrator — run the full cycle, or just one capability (requirements, plan, implement, review, session management). See [Coding Workflow](#coding-workflow) | All coding agents + implementation-plan, code-review-documentation, session-state, worktree-setup, worktree-cleanup skills |

---

## Coding Workflow

The `coding` workflow's `coding-orchestrator` is a hierarchical coordinator, not a fixed
script. It routes every request to only the capabilities that request needs, delegates to
specialist sub-agents, validates each result against a contract before accepting it, and
tracks progress on disk so a session can be resumed instead of restarted.

`coding` is already active by default in `.switchic/config.yaml` (`workflows.active`) after
`switchic init`. If you removed it, add `coding` back to that list, then regenerate:

```bash
switchic switch claude   # or kiro / github-copilot
```

Then just talk to `@coding-orchestrator` (or let your platform route to it automatically)
— you don't need to know the internal agent names.

### Goals

Every request is classified into one of six goals. Only that goal's capabilities run.

| Goal | What it does | Example prompt |
|---|---|---|
| `requirements-only` | Fetches Jira ticket context (and PR/MR context if mentioned) | `@coding-orchestrator fetch the requirements for QC-1234` |
| `plan-only` | Fetches context if missing, then writes an implementation plan and stops at human approval | `@coding-orchestrator write an implementation plan for QC-1234` |
| `implement-only` | Requires an approved plan; bootstraps a session if needed, then implements | `@coding-orchestrator the plan for QC-1234 is approved, implement it` |
| `review-only` | Reviews the current implementation against the approved plan and ticket | `@coding-orchestrator review the implementation for QC-1234` |
| `full-cycle` | Runs everything above in order, looping implement ⇄ review until approved | `@coding-orchestrator drive QC-1234 from ticket to reviewed code` |
| `session-ops` | Lists, inspects, or cleans up sessions | `@coding-orchestrator list active sessions` · `@coding-orchestrator clean up QC-1234` |

More example prompts:

```
@coding-orchestrator what does QC-1234 ask for?
@coding-orchestrator revise the plan for QC-1234 — skip the migration, it's already done
@coding-orchestrator just run code-implementer against the approved plan for QC-1234
@coding-orchestrator review what was implemented for QC-1234
@coding-orchestrator please complete task QC-1234
@coding-orchestrator resume QC-1234
@coding-orchestrator status for QC-1234
@coding-orchestrator abandon QC-1234
```

Missing inputs are resolved dynamically, cheapest source first: **disk** (state file,
saved context, the plan) → **delegate** to the agent that produces it → **ask you**, and
only for things no agent can produce (the Jira key, PR/MR id, base branch, plan approval).

### Human gates

Two points always wait for you, regardless of goal:

- **Plan approval** — the orchestrator never lets `code-implementer` run against a plan
  still marked `draft`. Give feedback and it revises; approve and it proceeds.
- **Escalations** — infeasibility, unclear requirements, or security concerns raised by
  any sub-agent go straight to you instead of being silently worked around.

### Single-repo and workspace mode

In a single repo, session artifacts live at `<repo>/generated-docs/<JIRA_KEY>/`. In a
[workspace](#multi-repo-workspace), they live at the workspace root instead, and the
orchestrator asks which registered repo a ticket targets before touching code. A ticket
spanning several repos runs as one orchestrator session per repo (each with its own
plan and review, namespaced under the repo's name) — **workspace mode is single-repo-per-session
by design**, so run them sequentially in one conversation or concurrently in separate tabs.

### Parallel sessions and git guardrails

Starting a new ticket while another is active spins up a git worktree + feature branch
next to your repo (`<repo>-session-<JIRA_KEY>`) so sessions never collide, and registers it
in `sessions/registry.json`. The worktree setup copies your gitignored local config
(`.env`, `.env.test`, `config/master.key`, `config/credentials/*.key`) into the new
worktree, since `git worktree add` only checks out tracked files — this is what makes
`rspec` (or any test suite reading local env files) work immediately in a fresh worktree.

No agent in this workflow ever commits, pushes, merges, or opens a PR/MR:

- Implementation changes are left **uncommitted** in the working tree so you can review
  the diff and commit/push manually with your own message (the implementer's summary
  includes a suggested one).
- Session cleanup (`@coding-orchestrator clean up QC-1234`) only removes the worktree
  directory and the registry entry — never a merge or push — and refuses to run while
  uncommitted changes exist, so nothing is silently discarded.

### Cost awareness

Multi-agent workflows spend tokens on every delegation, so the workflow is tuned to spend
them where judgment is actually needed:

- Mechanical agents (Jira/PR fetchers, session bootstrap, session manager) run on a small,
  fast model with capped turns; planning and review run on a mid-tier model; the
  orchestrator and implementer inherit your session's model since that's where quality is
  worth paying for.
- Fetched Jira/PR context and diffs are saved to disk once and passed between agents by
  file path — never re-fetched or re-pasted into every delegation.
- Re-reviews after revision only re-check the previously flagged findings and the new
  diff, not the entire checklist again.
- For a single, trivial capability, you can skip the orchestrator layer entirely and call
  a sub-agent directly, e.g. `@jira-requirements-fetcher QC-1234`.

---

## User-Defined Assets

You can define your own agents, skills, rules, and workflows alongside the bundled ones.

Place YAML files under `.switchic/` in your project root:

```
.switchic/
├── agents/
│   └── my-reviewer.yaml
├── skills/
│   └── my-skill.yaml          # flat format
│   └── my-folder-skill/       # folder format (skill.yaml + prompt.md)
│       ├── skill.yaml
│       └── prompt.md
├── rules/
│   └── my-rule.yaml
└── workflows/
    └── my-workflow.yaml
```

**Resolution order:**

1. `.switchic/<kind>/` in your project root — checked first
2. Bundled defaults shipped with the binary — used as fallback

A local file whose `name` field (or filename) matches a bundled asset fully replaces it. Names that don't collide are simply added to the registry.

**YAML format** is identical to the bundled assets. The `example.yaml` files in each bundled directory document every available field — refer to `internal/assets/bundled/` in the source.

Once a file is in place, enable it like any bundled asset:

```bash
switchic add agent my-reviewer
switchic add skill my-skill
switchic switch claude
```

---

## Config Reference

### `.switchic/config.yaml` (per repo)

Created by `switchic init`. Edit manually or use `switchic add` / `switchic remove`.

```yaml
# ── Platform & workflow ──────────────────────────────────────────────────────
platform: claude           # target AI assistant: "claude" or "github-copilot"

workflows:
  active:
    - coding               # workflow presets to apply; their agents/skills are auto-activated

# ── Project identity ─────────────────────────────────────────────────────────
# These fields are written into the generated CLAUDE.md header so the AI
# understands what it is working on without reading the whole codebase.

name: broadcast-service    # project name; defaults to directory name if omitted
description: >
  Microservice responsible for broadcasting messages to email, SMS, and push
  notification channels. Provides a unified interface and ensures reliable
  delivery across platforms.

language: go               # primary language (used if stack is not set)

stack:                     # tech stack list (renders instead of language when set)
  - Go
  - PostgreSQL

# ── Developer commands ────────────────────────────────────────────────────────
# Injected into CLAUDE.md as a quick-reference table so the AI knows what
# commands to run without reading Makefile or package.json.

commands:
  build:
    run: make build
    description: Compile binary with version ldflags
  dev:
    run: go run .
    description: Start the dev server
  test:
    run: make test
    description: Run the full test suite

# ── Directory map ─────────────────────────────────────────────────────────────
# Key → short description. Tells the AI where to look before searching.

structure:
  cmd/:      CLI entry points — wiring only, no business logic
  config/:   Configuration files and templates
  internal/: Private packages and all business logic
  pkg/:      Reusable libraries and utilities

# ── Coding conventions ────────────────────────────────────────────────────────
# Non-default patterns only. Skip anything enforced by a linter or the
# language standard — the AI already knows those.

conventions:
  - Go modules for dependency management
  - Makefile for build and common tasks
  - Docker for containerization
  - GitHub Actions for CI/CD

# ── Do / Don't ────────────────────────────────────────────────────────────────
dos:
  - Add new providers by implementing the Provider interface in pkg/providers/
  - Write table-driven tests with t.Run() for handlers and providers

donts:
  - Don't put business logic in cmd/ — it's wiring only
  - Don't call external SDKs directly — always go through internal/providers/

# ── Reference docs ────────────────────────────────────────────────────────────
# Paths relative to the repo root. The "when" field scopes the @-mention to a
# specific trigger so the AI only loads the doc when it is relevant.

docs:
  - path: README.md
    when: project overview or setup instructions are needed

# ── Components ────────────────────────────────────────────────────────────────
agents:
  active: []             # extends the workflow preset; leave empty to rely on preset only

skills:
  active: []

rules:
  active:
    - golang
```

Fields omitted from `switchic init` output are optional — add them as your project grows.

### `switchic.workspace.yaml` (multi-repo)

Created by `switchic workspace init`. Place this file in a parent directory that sits alongside your repos.

```yaml
name: billing-workspace
notes: Billing platform — API, web dashboard, and shared contracts.
platform: claude           # target platform for all repos in this workspace

# Optional: workspace-level component overrides.
# When set, these take priority over any per-repo .switchic/config.yaml values.
agents:
  active:
    - code-reviewer

skills:
  active:
    - implementation-plan

rules:
  active:
    - golang

# Optional: reference docs shared across the whole workspace. Paths are
# relative to the workspace root. The "when" field scopes the @-mention to a
# specific trigger so the AI only loads the doc when it is relevant.
docs:
  - path: docs/architecture.md
    when: understanding how the repos fit together
  - path: README.md

repos:
  - name: billing-api
    path: ../billing-api   # relative to this workspace file
    role: backend          # optional label — appears in the generated context summary
    notes: REST API, owns user accounts and invoicing

  - name: billing-web
    path: ../billing-web
    role: frontend
    notes: Customer-facing dashboard (React)

  - name: billing-contracts
    path: ../billing-contracts
    role: contracts
    notes: Shared OpenAPI specs + generated client SDKs
    # context_file: path to this repo's context file relative to workspace root.
    # Use when the file lives at a non-standard path.
    context_file: ../billing-contracts/docs/AGENTS.md
```

See `examples/multi-repo/` for a fully-populated sample.

---

## Cost Optimization

Every token loaded into an AI session costs money and reduces quality (longer context = more noise). `switchic status` shows a rough byte and token estimate for the current config.

The goal is to keep only what the current task needs active.

```bash
# Check current context size
switchic status

# Strip components you don't need for this task
switchic remove agent code-reviewer
switchic remove skill jira-ticket-description
switchic remove rule typescript

# Regenerate
switchic switch claude
```

**Workflow presets help here** — instead of manually toggling agents and skills for different task types, define a workflow that bundles the right set. Switch workflows when you switch task types.

In workspace mode, the generated `CLAUDE.md` contains only a structured summary of repos (names, roles, notes) rather than any file contents — keeping multi-repo sessions affordable.

---

## Architecture

```
cmd/                  CLI entry points (Cobra) — kept thin
internal/app          LoadContext + Resolve — the orchestrator between layers
internal/config       Project config model, load/save, mutators
internal/workspace    Workspace manifest and repo registry
internal/platform     Platform adapter interface + Claude and GitHub Copilot adapters
internal/workflow     Workflow preset model and registry
internal/agent        Agent definitions and registry
internal/skill        Skill definitions and registry
internal/rules        Rule definitions and registry
internal/cost         Context-size estimator
internal/project      Single-repo and workspace bootstrap
internal/output       Presentation helpers (table, info, JSON)
internal/util         FS, paths, sentinel errors
internal/assets       embed.FS for all bundled YAML defaults
```

All bundled defaults (`configs/`, `workflows/`, `agents/`, `skills/`, `rules/`) live under `internal/assets/bundled/` and ship inside the binary via Go's `embed` package. Callers use `assets.FS()` which strips the `bundled/` prefix, so paths look like `"agents/planner.yaml"` throughout the codebase.

Generation is **deterministic** — the same config and filesystem state always produce the same output files. The main context file (`CLAUDE.md` / `AGENTS.md`) resolves to a `.local` variant when a user-written file is detected (no switchic banner), but this resolution is itself deterministic. Pointer files (`.github/copilot-instructions.md`, `.kiro/steering/project.md`) are always overwritten.

---

## Roadmap

### Coming soon

- **Cursor adapter** — generate `.cursorrules` and Cursor agent files
- **Codex CLI adapter** — generate `AGENTS.md`
- **Windsurf adapter** — generate Windsurf rules

### Planned

- Per-platform token estimates in `status`
- Low-cost / full-context workflow presets
- Richer workspace context (dependency maps, selected file globs per repo)

---

## Contributing

Contributions are welcome. Here's how to get started:

```bash
git clone https://github.com/Dandi-Pangestu/switchic
cd switchic
make build    # compile to ./bin/switchic
make test     # go test ./...
make vet      # go vet ./...
make fmt      # gofmt -w .
```

### Adding a bundled agent, skill, rule, or workflow

Drop a YAML file under `internal/assets/bundled/<kind>/` and rebuild — no embed list to update. Copy an `example.yaml` in the same directory as a starting point.

### Adding a platform adapter

1. Implement the `platform.Adapter` interface (`internal/platform/platform.go`)
2. Register it in `platform.Get` (`internal/platform/registry.go`)
3. Add a platform config YAML under `internal/assets/bundled/configs/platforms/`
4. Add `copilot:` (or equivalent) blocks to bundled agent, skill, and rule YAMLs

### PR checklist

- `make test` passes
- `make vet` passes
- New bundled assets include a meaningful `description` field
- New commands have a `--help` usage string

---

## License

MIT. See [LICENSE](LICENSE).
