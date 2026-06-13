package output

import (
	"fmt"
	"io"
	"strings"
)

// Table is a minimal two-column key/value renderer used by `status`.
type Table struct {
	rows [][2]string
}

// Row adds a key/value pair.
func (t *Table) Row(k, v string) { t.rows = append(t.rows, [2]string{k, v}) }

// Render writes the aligned table to w.
func (t *Table) Render(w io.Writer) {
	width := 0
	for _, r := range t.rows {
		if l := len(r[0]); l > width {
			width = l
		}
	}
	for _, r := range t.rows {
		pad := strings.Repeat(" ", width-len(r[0]))
		fmt.Fprintf(w, "  %s%s : %s\n", r[0], pad, r[1])
	}
}
