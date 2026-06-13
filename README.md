# switchic

A Go CLI that switches a repo or workspace between AI coding platforms,
manages the agentic context that gets loaded into each one, and keeps that
context lean.

The MVP supports **Claude** as the first-class platform. The same architecture
will support Cursor, Copilot, Codex CLI, Windsurf, and others in later phases.

## What it does

- **Switch platforms** — `switchic switch claude` regenerates the
  platform-specific files (e.g. `CLAUDE.md`, `.claude/agents/*`) from one
  shared source of truth.
- **Multi-repo workspaces** *(coming soon)* — register related repos in
  `switchic.workspace.yaml` and generate a single Claude context that
  describes the whole system.
- **Cost-aware components** — turn agents, skills, and rules on or off so
  only what the current task needs is loaded into the session.
- **Workflow presets** — a workflow bundles the agents and skills for a
  task type (e.g. `coding`). Set one or more in config and every agent
  exports automatically — no manual list needed.
- **User-defined assets** — drop your own agents, skills, rules, or workflows
  into `.switchic/agents/`, `.switchic/skills/`, `.switchic/rules/`, or
  `.switchic/workflows/` and they are picked up automatically. A local asset
  with the same name as a bundled one replaces it.

## Install

You need Go 1.21+ installed (https://go.dev/dl/).

### System-wide install (recommended)

```bash
git clone <this repo>
cd switchic
make install        # installs to /usr/local/bin (may prompt for sudo)
switchic version
```

### User-local install (no sudo)

```bash
make user-install   # installs to ~/.local/bin
```

If the chosen directory is not on your `PATH`, the installer prints the
exact line to add to your shell rc:

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

## Quickstart — single repo

```bash
cd path/to/your/repo
switchic init                 # creates .switchic/config.yaml
switchic switch claude        # writes CLAUDE.md + .claude/agents + .claude/rules
switchic status               # see active components and cost summary
switchic remove agent verifier
switchic add skill summarize
switchic switch claude        # regenerate after mutations
```

## Quickstart — multi-repo workspace *(coming soon)*

> **Note:** Workspace commands are not yet available. Running any
> `switchic workspace` subcommand will print an error and exit.

Once released, the workflow will look like:

```bash
mkdir my-workspace && cd my-workspace
switchic workspace init --name billing-workspace
switchic workspace add ../billing-api      --role backend
switchic workspace add ../billing-web      --role frontend
switchic workspace add ../billing-contracts --role contracts
switchic workspace list
switchic switch claude        # CLAUDE.md now contains a workspace summary
```

In workspace mode, the generated `CLAUDE.md` will include a section that lists
every repo with its role and notes, so the assistant understands cross-repo
context without ingesting raw code.

## Commands

| Command | Purpose |
|---|---|
| `init` | Create `.switchic/config.yaml` in the current repo |
| `switch <platform>` | Update active platform and regenerate its files |
| `status` | Show platform, workflows, active components, cost summary |
| `workspace init` | *(coming soon)* Create `switchic.workspace.yaml` |
| `workspace add <path>` | *(coming soon)* Register a repo in the workspace |
| `workspace remove <name>` | *(coming soon)* Unregister a repo |
| `workspace list` | *(coming soon)* List repos (flags missing paths) |
| `add agent\|skill\|rule <name>` | Enable a component |
| `remove agent\|skill\|rule <name>` | Disable a component |
| `version` | Print the binary version |

## Bundled assets

These are the agents, skills, rules, and workflows that ship with the binary.
Use them with `switchic add` / `switchic remove`, or override them by placing
a file with the same name under `.switchic/` (see [User-defined assets](#user-defined-assets)).

### Agents

| Name | Description |
|---|---|
| `coding-orchestrator` | Entry point for the end-to-end coding workflow — orchestrates session bootstrap, planning, implementation, and review |
| `code-implementer` | Executes an approved implementation plan by writing code and updating tests |
| `code-reviewer` | Reviews a completed implementation for correctness, quality, security, and test coverage |
| `implementation-plan-writer` | Analyzes a Jira ticket and PR context to produce a detailed, structured implementation plan |
| `session-bootstraper` | Provisions the working environment for a new workflow session |
| `session-manager` | Manages parallel agent sessions — lists active sessions and cleans up completed ones |
| `jira-requirements-fetcher` | Fetches comprehensive Jira ticket context and metadata given a ticket key |
| `pr-details-fetcher` | Fetches pull request or merge request details given a PR/MR identifier |
| `example` | Demonstrates every available agent YAML field — use as a starting-point template |

### Skills

| Name | Description |
|---|---|
| `implementation-plan` | Guide for producing structured implementation plans for Jira tickets |
| `jira-ticket-description` | Guide for producing Jira ticket descriptions with consistent sections |
| `code-review-documentation` | Guide for producing structured code review summaries and documentation |
| `agent-import` | Convert an agent definition from any AI coding platform into switchic format |
| `skill-import` | Convert a skill definition from any AI coding platform into switchic format |
| `worktree-setup` | Provision a git worktree at a sibling path of the repository root |
| `worktree-cleanup` | Remove a git worktree created by the `worktree-setup` skill |
| `commit-msg` | Generate a conventional commit message from staged changes |
| `example` | Demonstrates every available skill YAML field — use as a starting-point template |

### Rules

Rules are referenced by their key path (without `.yaml`). Directory names
expand to all rules inside them — e.g. `backend` enables both `backend/api`
and `backend/database`.

| Key | Description |
|---|---|
| `golang` | Go coding standards (based on Effective Go) |
| `typescript` | TypeScript coding standards (based on Google TS Style Guide) |
| `backend/api` | REST API design standards (based on Azure API Design Guidelines) |
| `backend/database` | Database coding and querying best practices |

### Workflows

| Name | Description |
|---|---|
| `coding` | Drives a ticket from Jira to merged code in four phases: session bootstrap → planning → implementation → review |

---

## User-defined assets

You can define your own agents, skills, rules, and workflows by placing YAML
files under `.switchic/` in your project (or workspace root):

```
.switchic/
├── agents/
│   └── my-reviewer.yaml      # custom agent
├── skills/
│   └── my-skill.yaml         # custom skill (flat format)
│   └── my-folder-skill/      # custom skill (folder format)
│       ├── skill.yaml
│       └── prompt.md
├── rules/
│   └── my-rule.yaml          # custom rule (subdirectories supported)
└── workflows/
    └── my-workflow.yaml      # custom workflow
```

**Resolution order:**

1. `.switchic/<kind>/` in the project (or workspace) root — checked first
2. Bundled defaults shipped with the binary — used as fallback

A local file whose `name` field (or filename) matches a bundled asset fully
replaces that bundled asset. Names that don't collide with anything bundled
are simply added to the registry.

**YAML format** is identical to the bundled assets — refer to
`internal/assets/bundled/` in the source for examples.

Once the file is in place, enable the asset just like any bundled one:

```bash
switchic add agent my-reviewer
switchic add skill my-skill
switchic switch claude
```

## Config files

### `.switchic/config.yaml` (per repo)

```yaml
platform: claude
workflows:
  active: [coding]       # add more workflow names to merge their presets
language: go
agents:
  active: []             # extends the workflow preset; leave empty to use preset only
skills:
  active: []
rules:
  active: [global, golang]
```

### `switchic.workspace.yaml` (multi-repo) *(coming soon)*

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

## Generated outputs (Claude)

After `switchic switch claude`:

```
CLAUDE.md                       # main context Claude reads on launch
.claude/agents/<name>.md        # one file per active agent
.claude/rules/<name>.md         # one file per active rule
```

Generation is **deterministic** — the same config always produces the same
files, which makes the outputs safe to commit and review in diffs.

## Cost optimization

`status` prints a rough byte and token estimate of the context that would
be loaded. To shrink it:

```bash
switchic remove agent reviewer
switchic remove skill summarize
switchic remove rule docs
switchic switch claude
```

In workspace mode *(coming soon)* the context will only include a structured
summary of the repos (names, roles, notes), not their contents — keeping
multi-repo sessions affordable.

## Architecture

```
cmd/                CLI (Cobra) — thin command handlers
internal/app        Bootstrap and resolver between layers
internal/config     Project config model + load/save + mutators
internal/workspace  Workspace manifest + repo registry
internal/platform   Platform adapter interface + Claude adapter
internal/workflow   Workflow preset model + registry
internal/agent      Agent definitions + registry
internal/skill      Skill definitions + registry
internal/rules      Rule definitions + registry
internal/cost       Context-size estimator
internal/project    Single-repo + workspace bootstrap
internal/output     Presentation helpers (table, info, JSON)
internal/util       FS, paths, sentinel errors
internal/assets     embed.FS of all bundled YAML defaults
```

All bundled defaults (`configs/`, `workflows/`, `agents/`, `skills/`,
`rules/`) live under `internal/assets/bundled/` and ship inside the
binary via Go's `embed` package.

User-defined assets are loaded from `.switchic/{agents,skills,rules,workflows}/`
at runtime and merged on top of the bundled set before resolution. Local
definitions take precedence — same name means local wins.

## Roadmap

After the MVP stabilizes:

- additional platform adapters (Cursor, Copilot, Codex CLI, Windsurf)
- richer workspace context (dependency maps, selected file globs per repo)
- per-platform token estimates
- low-cost / full-context presets

## License

MIT (see LICENSE if present, or add one before publishing).
