package cmd

import (
	"errors"
	"fmt"
	"maps"

	"github.com/spf13/cobra"

	"github.com/Dandi-Pangestu/switchic/internal/agent"
	"github.com/Dandi-Pangestu/switchic/internal/app"
	"github.com/Dandi-Pangestu/switchic/internal/assets"
	"github.com/Dandi-Pangestu/switchic/internal/config"
	"github.com/Dandi-Pangestu/switchic/internal/output"
	"github.com/Dandi-Pangestu/switchic/internal/rules"
	"github.com/Dandi-Pangestu/switchic/internal/skill"
	"github.com/Dandi-Pangestu/switchic/internal/util"
	"github.com/Dandi-Pangestu/switchic/internal/workspace"
)

// componentKind identifies which list (agents/skills/rules) a mutation targets.
type componentKind int

const (
	kindAgent componentKind = iota
	kindSkill
	kindRule
)

func (k componentKind) label() string {
	switch k {
	case kindAgent:
		return "agent"
	case kindSkill:
		return "skill"
	default:
		return "rule"
	}
}

// validName checks whether name is known in either the bundled registry or the
// user-local .switchic/ directory for the given project root.
func validName(k componentKind, name string, root string) (bool, error) {
	// Build the merged map: bundled first, then local overrides.
	switch k {
	case kindAgent:
		m, err := agent.LoadAll()
		if err != nil {
			return false, err
		}
		if localFSys := assets.LocalFS(root); localFSys != nil {
			if local, err := agent.LoadAllFrom(localFSys); err == nil {
				maps.Copy(m, local)
			}
		}
		_, ok := m[name]
		return ok, nil
	case kindSkill:
		m, err := skill.LoadAll()
		if err != nil {
			return false, err
		}
		if localFSys := assets.LocalFS(root); localFSys != nil {
			if local, err := skill.LoadAllFrom(localFSys); err == nil {
				maps.Copy(m, local)
			}
		}
		_, ok := m[name]
		return ok, nil
	default:
		m, err := rules.LoadAll()
		if err != nil {
			return false, err
		}
		if localFSys := assets.LocalFS(root); localFSys != nil {
			if local, err := rules.LoadAllFrom(localFSys); err == nil {
				maps.Copy(m, local)
			}
		}
		_, ok := m[name]
		return ok, nil
	}
}

// listFor returns a pointer to the active slice on the right config object,
// preferring the workspace manifest when workspace mode is active and the
// workspace already has that list defined.
func listFor(k componentKind, ctx *app.Context) *[]string {
	if ctx.IsWorkspace {
		switch k {
		case kindAgent:
			return &ctx.Workspace.Agents.Active
		case kindSkill:
			return &ctx.Workspace.Skills.Active
		default:
			return &ctx.Workspace.Rules.Active
		}
	}
	switch k {
	case kindAgent:
		return &ctx.Project.Agents.Active
	case kindSkill:
		return &ctx.Project.Skills.Active
	default:
		return &ctx.Project.Rules.Active
	}
}

// runMutation performs an add or remove for a given component kind and
// persists the affected config file.
func runMutation(cmd *cobra.Command, k componentKind, name string, add bool) error {
	ctx, err := app.LoadContext(util.Cwd())
	if err != nil {
		return err
	}
	if !ctx.IsWorkspace && !ctx.HasProject() {
		return errors.New("no switchic config found — run `switchic init` first")
	}

	if add {
		ok, err := validName(k, name, ctx.PrimaryRoot())
		if err != nil {
			return err
		}
		if !ok {
			return fmt.Errorf("unknown %s %q (not found in bundled or local .switchic/)", k.label(), name)
		}
	}

	list := listFor(k, &ctx)
	changed := false
	if add {
		changed = config.AddTo(list, name)
	} else {
		changed = config.RemoveFrom(list, name)
	}
	if !changed {
		if add {
			output.Info(cmd.OutOrStdout(), "%s %q was already active", k.label(), name)
		} else {
			output.Info(cmd.OutOrStdout(), "%s %q was not active", k.label(), name)
		}
		return nil
	}

	if ctx.IsWorkspace {
		if err := workspace.Save(util.WorkspacePath(ctx.WorkspaceRoot), ctx.Workspace); err != nil {
			return err
		}
	} else {
		if err := config.Save(util.ProjectConfigPath(ctx.WorkingDir), ctx.Project); err != nil {
			return err
		}
	}

	verb := "Enabled"
	if !add {
		verb = "Disabled"
	}
	output.Info(cmd.OutOrStdout(), "%s %s %q. Run `switchic switch claude` to regenerate.", verb, k.label(), name)
	return nil
}
