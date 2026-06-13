package cmd

import (
	"maps"
	"strconv"
	"strings"

	"github.com/spf13/cobra"

	"github.com/Dandi-Pangestu/switchic/internal/agent"
	"github.com/Dandi-Pangestu/switchic/internal/app"
	"github.com/Dandi-Pangestu/switchic/internal/assets"
	"github.com/Dandi-Pangestu/switchic/internal/cost"
	"github.com/Dandi-Pangestu/switchic/internal/output"
	"github.com/Dandi-Pangestu/switchic/internal/rules"
	"github.com/Dandi-Pangestu/switchic/internal/skill"
	"github.com/Dandi-Pangestu/switchic/internal/util"
)

var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Show the current platform, workflow, and active components",
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx, err := app.LoadContext(util.Cwd())
		if err != nil {
			return err
		}
		if !ctx.IsWorkspace && !ctx.HasProject() {
			output.Info(cmd.OutOrStdout(), "No switchic config found in %s.", ctx.WorkingDir)
			output.Info(cmd.OutOrStdout(), "Run `switchic init` or `switchic workspace init` to get started.")
			return nil
		}

		resolved, err := app.Resolve(ctx.PrimaryRoot(), ctx.Project, ctx.Workspace, ctx.IsWorkspace)
		if err != nil {
			return err
		}

		t := output.Table{}
		t.Row("mode", modeLabel(ctx))
		t.Row("platform", resolved.Platform)
		wfNames := make([]string, len(resolved.Workflows))
		for i, w := range resolved.Workflows {
			wfNames[i] = w.Name
		}
		t.Row("workflows", strings.Join(wfNames, ", "))
		if ctx.Project.Language != "" {
			t.Row("language", ctx.Project.Language)
		}
		t.Row("agents", namesOfAgents(resolved.Agents))
		t.Row("skills", namesOfSkills(resolved.Skills))
		t.Row("rules", namesOfRules(resolved.Rules))
		if ctx.IsWorkspace {
			repos := make([]string, 0, len(ctx.Workspace.Repos))
			for _, r := range ctx.Workspace.Repos {
				repos = append(repos, r.Name)
			}
			t.Row("repos", strings.Join(repos, ", "))
			t.Row("repo count", strconv.Itoa(len(ctx.Workspace.Repos)))
		}
		t.Render(cmd.OutOrStdout())

		// Cost summary needs the full merged registries (bundled + local).
		allAgents, _ := agent.LoadAll()
		allSkills, _ := skill.LoadAll()
		allRules, _ := rules.LoadAll()
		if localFSys := assets.LocalFS(ctx.PrimaryRoot()); localFSys != nil {
			if local, err := agent.LoadAllFrom(localFSys); err == nil {
				maps.Copy(allAgents, local)
			}
			if local, err := skill.LoadAllFrom(localFSys); err == nil {
				maps.Copy(allSkills, local)
			}
			if local, err := rules.LoadAllFrom(localFSys); err == nil {
				maps.Copy(allRules, local)
			}
		}
		summary := cost.Estimate(
			allAgents, namesList(resolved.Agents, agentName),
			allSkills, namesList(resolved.Skills, skillName),
			allRules, namesList(resolved.Rules, ruleName),
			len(ctx.Workspace.Repos),
		)
		output.Section(cmd.OutOrStdout(), "Cost")
		summary.Print(cmd.OutOrStdout())
		return nil
	},
}

func modeLabel(c app.Context) string {
	if c.IsWorkspace && c.HasProject() {
		return "workspace + project"
	}
	if c.IsWorkspace {
		return "workspace"
	}
	return "project"
}

func namesList[T any](defs []T, f func(T) string) []string {
	out := make([]string, len(defs))
	for i, d := range defs {
		out[i] = f(d)
	}
	return out
}

func agentName(a agent.Definition) string { return a.Name }
func skillName(s skill.Definition) string { return s.Name }
func ruleName(r rules.Definition) string  { return r.Name }

func namesOfAgents(d []agent.Definition) string { return strings.Join(namesList(d, agentName), ", ") }
func namesOfSkills(d []skill.Definition) string { return strings.Join(namesList(d, skillName), ", ") }
func namesOfRules(d []rules.Definition) string  { return strings.Join(namesList(d, ruleName), ", ") }

func init() { rootCmd.AddCommand(statusCmd) }
