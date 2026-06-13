# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Build, test, run

```bash
make build            # compile to ./bin/switchic (with version ldflags)
make install          # build + install to /usr/local/bin (may sudo)
make user-install     # build + install to ~/.local/bin (no sudo)
make uninstall        # remove from both locations
make test             # go test ./...
make vet              # go vet ./...
make fmt              # gofmt -w .
make clean            # rm -rf ./bin

# fast inner loop (no Makefile, no version ldflags):
go build ./...
go vet ./...

# run a single test once they exist:
go test -run TestName ./internal/<pkg>/...
```

The Makefile injects `VERSION` via `-ldflags -X github.com/Dandi-Pangestu/switchic/cmd.Version=...`. Plain `go build` produces a binary that reports `dev`.

Module path: `github.com/Dandi-Pangestu/switchic`. Go 1.21+ required (uses `slices` and `embed` with `all:` directives).

## Architecture

### Layered flow

Reading the code top-down:

1. **`cmd/`** — Cobra command handlers, kept thin. Each command reads CLI args, calls into `internal/app`, formats output via `internal/output`.
2. **`internal/app`** — the orchestrator. `LoadContext(cwd)` discovers whether the user is in a workspace, a project, both, or neither. `Resolve(...)` merges configs and produces a `Resolved` struct that the platform adapter consumes via `Resolved.ToContext()`.
3. **`internal/platform`** — the abstraction boundary. `Adapter.Generate(ctx)` is the single extension point for new platforms (Cursor, Copilot, etc.). `platform.Get(name)` is the registry; today it only returns `Claude{}`.
4. **`internal/{agent,skill,rules,workflow}`** — each package owns its YAML model, a registry that reads from the embedded asset FS, and a resolver that filters by active names.
5. **`internal/{config,workspace,project}`** — config models, load/save, mutators, single-repo + workspace bootstrap.

### Embedded assets pattern (important)

All bundled defaults live under `internal/assets/bundled/` (`configs/`, `workflows/`, `agents/`, `skills/`, `rules/`). They ship inside the binary via `//go:embed all:bundled` in `internal/assets/assets.go`.

Callers must use `assets.FS()` — which returns `fs.Sub(raw, "bundled")` — so they see paths like `"agents/planner.yaml"`, not `"bundled/agents/planner.yaml"`. The `bundled/` prefix is stripped at the boundary.

When adding a new bundled YAML, drop it under `internal/assets/bundled/<kind>/` and rebuild — no embed list to update.

### Two-tier config with workspace precedence

There are two config files:
- `.switchic/config.yaml` — per-repo, written by `init`
- `switchic.workspace.yaml` — per-workspace, written by `workspace init`

`app.Resolve` applies workspace overrides on top of the project config: when workspace mode is active AND the workspace has a non-empty list for agents/skills/rules/workflow/platform, the workspace value wins. The `workspace init` flow seeds these lists from `config.Defaults()` so the workspace is usable immediately.

Commands that mutate config (`add`, `remove`) write to the workspace manifest when in workspace mode, otherwise the project config. The helper that handles this is `cmd/mutate.go:runMutation` — all six `add|remove agent|skill|rule` commands are thin wrappers around it. Don't duplicate that logic; extend `runMutation` and `listFor`.

### Deterministic generation

`platform.Claude.Generate` overwrites `CLAUDE.md` and `.claude/{agents,rules}/*.md` on every `switch claude`. The generated `CLAUDE.md` itself contains a banner saying "do not edit by hand" — the tool blows away local edits. Manual changes belong in the bundled YAMLs (or, eventually, in user-supplied agent/skill/rule overrides).

Same inputs → same outputs. This is intentional for golden testing and for safe re-runs.

### Optional workflow stages

`workflow.Stage.Optional == true` means: drop the stage if its agent is not in the active list. Required stages stay in the plan even when their agent is disabled, so `cmd/run` can flag them as broken with the `[x]` marker. The logic is in `internal/workflow/stage.go`.

### Mode discovery

`util.FindWorkspaceRoot` walks up from cwd looking for `switchic.workspace.yaml`. Project config is read from cwd's `.switchic/config.yaml` directly — it does **not** walk up. This means: running `switchic status` from a subdirectory of a workspace will find the workspace but not the project config of the current sub-repo. That is by design — workspace context drives the platform output, not the inner repo's local config.

## Conventions used here

- Errors funnel through `internal/util`'s sentinels (`ErrNotFound`, `ErrAlreadyExists`, `ErrInvalidConfig`, `ErrUnknownPlatform`) and `util.Wrap` for context. Callers branch with `errors.Is`.
- Atomic writes only: `util.WriteFile` writes to `<path>.tmp` then renames. Don't introduce direct `os.WriteFile` for config or generated files.
- Sorted, deduplicated lists in active sets: `config.AddTo` sorts on insert.
- Diagnostics about "undefined symbol X" that appear right after creating a file but before its sibling in the same package are stale — same-package symbols resolve once both files exist.

## Reference

- `README.md` — user-facing quickstart, command table, install paths.
- `implementation-plan_Version4.md` — original spec; MVP is Claude-only by design. Don't extend the platform registry without adding the corresponding adapter.
