package agent

// ResolveActive filters the registry down to the set of names listed as
// active. Unknown names are silently dropped — validation happens at the
// command layer where richer messaging is available.
func ResolveActive(all map[string]Definition, active []string) []Definition {
	out := make([]Definition, 0, len(active))
	for _, name := range active {
		if d, ok := all[name]; ok {
			out = append(out, d)
		}
	}
	return out
}
