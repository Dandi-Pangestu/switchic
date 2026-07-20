// Package app wires the domain packages together for the CLI commands. It
// hides registry loading and resolution behind a small surface so that cmd/
// stays focused on flag handling and presentation.
package app

import (
	"maps"
	"sort"

	"github.com/Dandi-Pangestu/switchic/internal/agent"
	"github.com/Dandi-Pangestu/switchic/internal/assets"
	"github.com/Dandi-Pangestu/switchic/internal/config"
	"github.com/Dandi-Pangestu/switchic/internal/platform"
	"github.com/Dandi-Pangestu/switchic/internal/rules"
	"github.com/Dandi-Pangestu/switchic/internal/skill"
	"github.com/Dandi-Pangestu/switchic/internal/util"
	"github.com/Dandi-Pangestu/switchic/internal/workflow"
	"github.com/Dandi-Pangestu/switchic/internal/workspace"
)

// Resolved is the fully merged view of project + workspace + bundled registry
// that the platform adapter consumes.
type Resolved struct {
	Platform  string
	Workflows []workflow.Workflow
	Agents    []agent.Definition
	Skills    []skill.Definition
	Rules     []rules.Definition
	Root      string
	Project   config.Project
	Workspace workspace.Manifest
	IsWS      bool
}

// Resolve loads all registries, merges project/workspace overrides, and
// filters down to active components. The caller chooses which is "primary"
// by passing isWorkspace.
//
// User-defined assets under root/.switchic/{agents,skills,rules,workflows}/
// are loaded first and take precedence over bundled defaults when names collide.
func Resolve(root string, proj config.Project, ws workspace.Manifest, isWorkspace bool) (Resolved, error) {
	// Pick the effective active lists. Workspace overrides project lists
	// only when workspace mode is active and the workspace lists are set.
	active := proj
	if isWorkspace {
		if len(ws.Workflows.Active) > 0 {
			active.Workflows = ws.Workflows
		}
		if len(ws.Agents.Active) > 0 {
			active.Agents = ws.Agents
		}
		if len(ws.Skills.Active) > 0 {
			active.Skills = ws.Skills
		}
		if len(ws.Rules.Active) > 0 {
			active.Rules = ws.Rules
		}
		if ws.Platform != "" {
			active.Platform = ws.Platform
		}
		if len(ws.Docs) > 0 {
			active.Docs = ws.Docs
		}
	}
	if active.Platform == "" {
		active.Platform = "claude"
	}
	if len(active.Workflows.Active) == 0 {
		active.Workflows.Active = []string{"coding"}
	}

	// Load bundled registries.
	allAgents, err := agent.LoadAll()
	if err != nil {
		return Resolved{}, err
	}
	allSkills, err := skill.LoadAll()
	if err != nil {
		return Resolved{}, err
	}
	allRules, err := rules.LoadAll()
	if err != nil {
		return Resolved{}, err
	}
	allWorkflows, err := workflow.LoadAll()
	if err != nil {
		return Resolved{}, err
	}

	// Overlay user-local assets from root/.switchic/ — local definitions win
	// on name collision, so users can fully replace a bundled asset.
	if localFSys := assets.LocalFS(root); localFSys != nil {
		if local, err := agent.LoadAllFrom(localFSys); err != nil {
			return Resolved{}, util.Wrap(err, "load local agents")
		} else {
			maps.Copy(allAgents, local)
		}
		if local, err := skill.LoadAllFrom(localFSys); err != nil {
			return Resolved{}, util.Wrap(err, "load local skills")
		} else {
			maps.Copy(allSkills, local)
		}
		if local, err := rules.LoadAllFrom(localFSys); err != nil {
			return Resolved{}, util.Wrap(err, "load local rules")
		} else {
			maps.Copy(allRules, local)
		}
		if local, err := workflow.LoadAllFrom(localFSys); err != nil {
			return Resolved{}, util.Wrap(err, "load local workflows")
		} else {
			maps.Copy(allWorkflows, local)
		}
	}

	// Load all active workflows and collect their preset agents/skills.
	var workflows []workflow.Workflow
	var wfAgents, wfSkills []string
	for _, name := range active.Workflows.Active {
		w, ok := allWorkflows[name]
		if !ok {
			return Resolved{}, util.Wrap(util.ErrNotFound, "workflow %q", name)
		}
		workflows = append(workflows, w)
		wfAgents = append(wfAgents, w.Agents...)
		wfSkills = append(wfSkills, w.Skills...)
	}

	// Workflow presets are defaults; config lists extend them.
	agentNames := mergeUnique(wfAgents, active.Agents.Active)
	activeAgents := agent.ResolveActive(allAgents, agentNames)

	// Union workflow skills + config skills + required_skills from active agents.
	skillSet := make(map[string]struct{})
	for _, name := range mergeUnique(wfSkills, active.Skills.Active) {
		skillSet[name] = struct{}{}
	}
	for _, a := range activeAgents {
		for _, s := range a.RequiredSkills {
			skillSet[s] = struct{}{}
		}
	}
	mergedSkills := make([]string, 0, len(skillSet))
	for name := range skillSet {
		mergedSkills = append(mergedSkills, name)
	}
	sort.Strings(mergedSkills)

	return Resolved{
		Platform:  active.Platform,
		Workflows: workflows,
		Agents:    activeAgents,
		Skills:    skill.ResolveActive(allSkills, mergedSkills),
		Rules:     rules.ResolveActive(allRules, active.Rules.Active),
		Root:      root,
		Project:   active,
		Workspace: ws,
		IsWS:      isWorkspace,
	}, nil
}

// mergeUnique returns a deduplicated slice of base + extra, preserving order.
func mergeUnique(base, extra []string) []string {
	seen := make(map[string]struct{}, len(base)+len(extra))
	out := make([]string, 0, len(base)+len(extra))
	for _, s := range append(base, extra...) {
		if _, ok := seen[s]; !ok {
			seen[s] = struct{}{}
			out = append(out, s)
		}
	}
	return out
}

// ToContext maps a Resolved into the platform.Context the adapter expects.
func (r Resolved) ToContext() platform.Context {
	return platform.Context{
		Root:        r.Root,
		Project:     r.Project,
		Workspace:   r.Workspace,
		IsWorkspace: r.IsWS,
		Workflows:   r.Workflows,
		Agents:      r.Agents,
		Skills:      r.Skills,
		Rules:       r.Rules,
	}
}
