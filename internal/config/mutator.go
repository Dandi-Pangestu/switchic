package config

import (
	"slices"
	"sort"
)

// AddTo appends name into the slice, deduping. Returns true if added.
func AddTo(list *[]string, name string) bool {
	if slices.Contains(*list, name) {
		return false
	}
	*list = append(*list, name)
	sort.Strings(*list)
	return true
}

// RemoveFrom removes name from the slice. Returns true if removed.
func RemoveFrom(list *[]string, name string) bool {
	out := (*list)[:0]
	removed := false
	for _, v := range *list {
		if v == name {
			removed = true
			continue
		}
		out = append(out, v)
	}
	*list = out
	return removed
}
