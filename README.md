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
- [Commands](#commands)
- [Bundled Library](#bundled-library)
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

**Example use case:** the `coding` workflow activates the full agent pipeline for ticket-to-PR development: context fetcher, planner, implementer, reviewer, and session manager — all wired together.

---

## AI Platform Support

| Platform | Status | Generated files |
|---|---|---|
| **Claude** (Claude Code) | Ready | `CLAUDE.md`, `.claude/agents/*.md`, `.claude/rules/*.md`, `.claude/skills/*/SKILL.md` |
| **GitHub Copilot** | Ready | `AGENTS.md`, `.github/copilot-instructions.md`, `.github/agents/*.agent.md`, `.github/instructions/*.instructions.md`, `.github/skills/*/SKILL.md` |
| **Kiro** (Kiro CLI) | Ready | `AGENTS.md`, `.kiro/steering/project.md`, `.kiro/agents/*.json`, `.kiro/steering/*.md`, `.kiro/skills/*/SKILL.md` |
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

# 2. Generate platform files — pick your AI assistant
switchic switch claude          # for Claude Code
switchic switch github-copilot  # for GitHub Copilot
switchic switch kiro            # for Kiro

# 3. Check what's active and the token cost
switchic status

# 4. Tune the context — disable what you don't need right now
switchic remove agent code-reviewer
switchic add skill commit-msg

# 5. Regenerate after any mutation
switchic switch claude   # or: switchic switch kiro / switchic switch github-copilot
```

After `switchic switch claude`, your repo will have:

```
CLAUDE.md                       # main context Claude reads on launch
.claude/agents/<name>.md        # one file per active agent
.claude/rules/<name>.md         # one file per active rule
.claude/skills/<name>/SKILL.md  # one directory per active skill
```

After `switchic switch github-copilot`, your repo will have:

```
AGENTS.md                                          # full project context (primary source of truth)
.github/copilot-instructions.md                    # pointer to AGENTS.md via @AGENTS.md
.github/agents/<name>.agent.md                     # one file per active agent
.github/instructions/<name>.instructions.md        # one file per active rule (path-specific)
.github/skills/<name>/SKILL.md                     # one directory per active skill
```

After `switchic switch kiro`, your repo will have:

```
AGENTS.md                          # main context file (AGENTS.md standard)
.kiro/steering/project.md          # always-included steering pointer to AGENTS.md
.kiro/agents/<name>.json           # one file per active agent (Kiro JSON format)
.kiro/steering/<name>.md           # one file per active rule (always-included steering)
.kiro/skills/<name>/SKILL.md       # one directory per active skill
```

### Multi-repo workspace *(coming soon)*

> Workspace commands are not yet released. Running any `switchic workspace` subcommand will print an error and exit.

Once released, the workflow will look like this:

```bash
mkdir my-workspace && cd my-workspace

switchic workspace init --name billing-workspace
switchic workspace add ../billing-api      --role backend
switchic workspace add ../billing-web      --role frontend
switchic workspace add ../billing-contracts --role contracts
switchic workspace list

switchic switch claude   # CLAUDE.md now includes a workspace-level summary
```

In workspace mode, the generated `CLAUDE.md` includes a structured summary of all repos (names, roles, notes) so the AI understands the full system without ingesting raw code from every repo.

---

## Commands

| Command | Purpose |
|---|---|
| `init` | Create `.switchic/config.yaml` in the current repo |
| `switch <platform>` | Regenerate platform files from the current config |
| `status` | Show active platform, workflows, components, and token cost |
| `add agent <name>` | Enable an agent |
| `add skill <name>` | Enable a skill |
| `add rule <name>` | Enable a rule |
| `remove agent <name>` | Disable an agent |
| `remove skill <name>` | Disable a skill |
| `remove rule <name>` | Disable a rule |
| `workspace init` | *(coming soon)* Create `switchic.workspace.yaml` |
| `workspace add <path>` | *(coming soon)* Register a repo in the workspace |
| `workspace remove <name>` | *(coming soon)* Unregister a repo |
| `workspace list` | *(coming soon)* List registered repos |
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
| `agent-import` | Convert an agent definition from any AI coding platform into switchic format |
| `skill-import` | Convert a skill definition from any AI coding platform into switchic format |
| `worktree-setup` | Provision a git worktree at a sibling path of the repository root |
| `worktree-cleanup` | Remove a git worktree created by the `worktree-setup` skill |
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
| `coding` | Drives a ticket from Jira to merged code in four phases: context gathering → planning → implementation loop → session cleanup | All coding agents + implementation-plan, code-review-documentation, worktree-setup, worktree-cleanup skills |

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

### `switchic.workspace.yaml` (multi-repo) *(coming soon)*

Created by `switchic workspace init`.

```yaml
name: billing-workspace
platform: claude
repos:
  - name: billing-api
    path: ../billing-api
    role: backend
  - name: billing-web
    path: ../billing-web
    role: frontend
```

See `examples/` for fully-populated samples.

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

In workspace mode *(coming soon)*, the generated `CLAUDE.md` will contain only a structured summary of repos (names, roles, notes) rather than any file contents — keeping multi-repo sessions affordable.

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

Generation is **deterministic** — the same config always produces the same output files. This makes generated files safe to commit and review in diffs.

---

## Roadmap

### Coming soon

- **Workspace support** — multi-repo context (`switchic workspace init/add/remove/list`)
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
