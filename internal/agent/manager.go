package agent

import "github.com/Dandi-Pangestu/switchic/internal/util"

// Validate ensures every active name exists in the registry. Unknown names
// are returned as a slice of errors-aware messages.
func Validate(all map[string]Definition, active []string) []error {
	var errs []error
	for _, name := range active {
		if _, ok := all[name]; !ok {
			errs = append(errs, util.Wrap(util.ErrNotFound, "agent %q", name))
		}
	}
	return errs
}
