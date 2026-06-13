package output

import (
	"encoding/json"
	"io"
)

// JSON writes v as indented JSON. Used by --json flags on commands.
func JSON(w io.Writer, v any) error {
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(v)
}
