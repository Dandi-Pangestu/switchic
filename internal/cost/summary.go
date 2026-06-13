package cost

import (
	"fmt"
	"io"
)

// Print writes a one-block summary suitable for `status` output.
func (s Summary) Print(w io.Writer) {
	fmt.Fprintln(w, "Cost summary:")
	fmt.Fprintf(w, "  agents: %d   skills: %d   rules: %d   repos: %d\n",
		s.Agents, s.Skills, s.Rules, s.Repos)
	fmt.Fprintf(w, "  approx context: %d bytes (~%d tokens)\n",
		s.ApproxBytes, s.ApproxTokens)
}
