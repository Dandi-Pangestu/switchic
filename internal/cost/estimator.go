// Package cost provides a lightweight estimator for the size of the context
// switchic would produce for a target platform. It is intentionally
// approximate — the goal is to give users a relative signal, not an exact
// token count.
package cost

import (
	"github.com/Dandi-Pangestu/switchic/internal/agent"
	"github.com/Dandi-Pangestu/switchic/internal/rules"
	"github.com/Dandi-Pangestu/switchic/internal/skill"
)

// Summary is a rough size accounting of active components.
type Summary struct {
	Agents      int
	Skills      int
	Rules       int
	Repos       int
	ApproxBytes int
	ApproxTokens int
}

// Estimate counts characters across active definitions and divides by 4 to
// approximate tokens — close enough for "is my context bloated?" signals.
func Estimate(
	allAgents map[string]agent.Definition, activeAgents []string,
	allSkills map[string]skill.Definition, activeSkills []string,
	allRules map[string]rules.Definition, activeRules []string,
	repoCount int,
) Summary {
	bytes := 0
	for _, n := range activeAgents {
		if d, ok := allAgents[n]; ok {
			bytes += len(d.Description) + len(d.Instructions)
			for _, s := range d.RequiredSkills {
				bytes += len(s)
			}
		}
	}
	for _, n := range activeSkills {
		if d, ok := allSkills[n]; ok {
			bytes += len(d.Description) + len(d.Prompt)
		}
	}
	for _, n := range activeRules {
		if d, ok := allRules[n]; ok {
			bytes += len(d.Description) + len(d.Content)
		}
	}
	return Summary{
		Agents:       len(activeAgents),
		Skills:       len(activeSkills),
		Rules:        len(activeRules),
		Repos:        repoCount,
		ApproxBytes:  bytes,
		ApproxTokens: bytes / 4,
	}
}
