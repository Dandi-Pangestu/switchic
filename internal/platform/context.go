package platform

import (
	"fmt"
	"path/filepath"
	"sort"
	"strings"
)

// buildContextBody renders the platform-agnostic sections of the main context
// file: header (project/workspace identity), commands, structure, conventions,
// do/don't, and references. The trailing separator and platform footer are left
// to the caller.
func buildContextBody(ctx Context) string {
	var b strings.Builder

	// ── Header: project / workspace identity ─────────────────────────────────
	if ctx.IsWorkspace {
		fmt.Fprintf(&b, "# %s\n\n", ctx.Workspace.Name)
		if ctx.Workspace.Notes != "" {
			fmt.Fprintf(&b, "%s\n\n", ctx.Workspace.Notes)
		}
		b.WriteString("## Repos\n\n")
		for _, r := range ctx.Workspace.Repos {
			role := r.Role
			if role == "" {
				role = "repo"
			}
			fmt.Fprintf(&b, "- **%s** (`%s`) — %s", r.Name, r.Path, role)
			if r.Notes != "" {
				fmt.Fprintf(&b, ": %s", r.Notes)
			}
			b.WriteString("\n")
		}
		b.WriteString("\n")
	} else {
		name := ctx.Project.Name
		if name == "" {
			name = filepath.Base(ctx.Root)
		}
		fmt.Fprintf(&b, "# %s\n\n", name)

		if ctx.Project.Description != "" {
			fmt.Fprintf(&b, "%s\n\n", ctx.Project.Description)
		}

		if len(ctx.Project.Stack) > 0 {
			fmt.Fprintf(&b, "**Stack:** %s\n\n", strings.Join(ctx.Project.Stack, " · "))
		} else if ctx.Project.Language != "" && ctx.Project.Language != "auto" {
			fmt.Fprintf(&b, "**Language:** %s\n\n", ctx.Project.Language)
		}
	}

	// ── Commands ──────────────────────────────────────────────────────────────
	if len(ctx.Project.Commands) > 0 {
		b.WriteString("## Commands\n\n")
		b.WriteString("| Label | Command | Description |\n")
		b.WriteString("|-------|---------|-------------|\n")
		keys := make([]string, 0, len(ctx.Project.Commands))
		for k := range ctx.Project.Commands {
			keys = append(keys, k)
		}
		sort.Strings(keys)
		for _, k := range keys {
			cmd := ctx.Project.Commands[k]
			fmt.Fprintf(&b, "| `%s` | `%s` | %s |\n", k, cmd.Run, cmd.Description)
		}
		b.WriteString("\n")
	}

	// ── Directory structure ───────────────────────────────────────────────────
	if len(ctx.Project.Structure) > 0 {
		b.WriteString("## Structure\n\n")
		keys := make([]string, 0, len(ctx.Project.Structure))
		for k := range ctx.Project.Structure {
			keys = append(keys, k)
		}
		sort.Strings(keys)
		for _, k := range keys {
			fmt.Fprintf(&b, "- `%s` — %s\n", k, ctx.Project.Structure[k])
		}
		b.WriteString("\n")
	}

	// ── Conventions ───────────────────────────────────────────────────────────
	if len(ctx.Project.Conventions) > 0 {
		b.WriteString("## Conventions\n\n")
		b.WriteString("> Non-default patterns only — skip anything enforced by a linter or standard for the language.\n\n")
		for _, c := range ctx.Project.Conventions {
			fmt.Fprintf(&b, "- %s\n", c)
		}
		b.WriteString("\n")
	}

	// ── Do / Don't ───────────────────────────────────────────────────────────
	if len(ctx.Project.Dos) > 0 || len(ctx.Project.Donts) > 0 {
		b.WriteString("## Do / Don't\n\n")
		if len(ctx.Project.Dos) > 0 {
			b.WriteString("### Do\n\n")
			for _, item := range ctx.Project.Dos {
				fmt.Fprintf(&b, "- %s\n", item)
			}
			b.WriteString("\n")
		}
		if len(ctx.Project.Donts) > 0 {
			b.WriteString("### Don't\n\n")
			for _, item := range ctx.Project.Donts {
				fmt.Fprintf(&b, "- %s\n", item)
			}
			b.WriteString("\n")
		}
	}

	// ── Reference docs ────────────────────────────────────────────────────────
	if len(ctx.Project.Docs) > 0 {
		b.WriteString("## References\n\n")
		for _, doc := range ctx.Project.Docs {
			if doc.When != "" {
				fmt.Fprintf(&b, "- When %s, read `@%s`.\n", doc.When, doc.Path)
			} else {
				fmt.Fprintf(&b, "- @%s\n", doc.Path)
			}
		}
		b.WriteString("\n")
	}

	return b.String()
}
