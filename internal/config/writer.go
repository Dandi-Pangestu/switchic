package config

import (
	"fmt"
	"strings"

	"github.com/Dandi-Pangestu/switchic/internal/util"

	"gopkg.in/yaml.v3"
)

// Save serializes p as YAML and writes it to path atomically.
func Save(path string, p Project) error {
	data, err := yaml.Marshal(p)
	if err != nil {
		return util.Wrap(err, "marshal project config")
	}
	return util.WriteFile(path, data)
}

// WriteInitial writes a human-friendly config.yaml with commented placeholders
// for all optional fields. Used by `init` so users see every knob upfront.
func WriteInitial(path string, p Project) error {
	lang := p.Language
	if lang == "" {
		lang = "auto"
	}

	fmtList := func(items []string) string {
		lines := make([]string, len(items))
		for i, item := range items {
			lines[i] = "    - " + item
		}
		return strings.Join(lines, "\n")
	}

	content := fmt.Sprintf(`# switchic project configuration
# Edit this file to customise the generated CLAUDE.md.
# Re-run `+"`switchic switch claude`"+` after any change.

platform: %s
workflows:
  active:
%s
language: %s

# ── Project identity ──────────────────────────────────────────────────────────
# These fields drive the header of the generated CLAUDE.md.
# name: My Project
# description: One or two sentences describing what this project does and why.
# stack:
#   - Go
#   - PostgreSQL

# ── Essential commands ────────────────────────────────────────────────────────
# Rendered as a Commands table (Label | Command | Description).
# Include only the commands developers run daily.
# commands:
#   build:
#     run: make build
#     description: Compile binary with version ldflags
#   test:
#     run: make test
#     description: Run the full test suite
#   dev:
#     run: go run .
#     description: Start the dev server

# ── Directory structure ───────────────────────────────────────────────────────
# Non-obvious paths only — skip anything self-evident from the name.
# structure:
#   cmd/:      CLI entry points
#   internal/: private packages

# ── Conventions ───────────────────────────────────────────────────────────────
# Non-default patterns only — skip anything enforced by a linter or obvious
# from the language standard.
# conventions:
#   - Use Zustand for state management, never Redux
#   - All API responses use { success, data, error } shape
#   - Database migrations go in src/database/migrations/

# ── Do / Don't ────────────────────────────────────────────────────────────────
# Explicit guardrails so Claude doesn't generate code that breaks your patterns.
# dos:
#   - Add new providers by implementing the Provider interface in pkg/providers/
#   - Write table-driven tests with t.Run() for handlers and providers
# donts:
#   - Don't put business logic in cmd/ — it's wiring only
#   - Don't call external SDKs directly — always go through internal/providers/

# ── Reference docs ────────────────────────────────────────────────────────────
# Files Claude should read when a specific context arises.
# docs:
#   - path: docs/api.md
#     when: working on API endpoints
#   - path: docs/db-schema.md
#     when: changing database models

agents:
  active:
%s

skills:
  active:
%s

rules:
  active:
%s
`, p.Platform, fmtList(p.Workflows.Active), lang, fmtList(p.Agents.Active), fmtList(p.Skills.Active), fmtList(p.Rules.Active))

	return util.WriteFile(path, []byte(content))
}
