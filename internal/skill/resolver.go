package skill

import "github.com/Dandi-Pangestu/switchic/internal/util"

// ResolveActive filters the registry to the set of active names.
func ResolveActive(all map[string]Definition, active []string) []Definition {
	out := make([]Definition, 0, len(active))
	for _, name := range active {
		if d, ok := all[name]; ok {
			out = append(out, d)
		}
	}
	return out
}

// Validate ensures every active name exists in the registry.
func Validate(all map[string]Definition, active []string) []error {
	var errs []error
	for _, name := range active {
		if _, ok := all[name]; !ok {
			errs = append(errs, util.Wrap(util.ErrNotFound, "skill %q", name))
		}
	}
	return errs
}
