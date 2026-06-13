// Package output centralizes CLI presentation helpers so commands stay slim.
package output

import (
	"fmt"
	"io"
)

// Info prints a single-line informational message.
func Info(w io.Writer, format string, args ...any) {
	fmt.Fprintf(w, format+"\n", args...)
}

// Section prints a header above a block of detail.
func Section(w io.Writer, title string) {
	fmt.Fprintf(w, "\n== %s ==\n", title)
}

// Bullet prints an indented bullet line.
func Bullet(w io.Writer, format string, args ...any) {
	fmt.Fprintf(w, "  - "+format+"\n", args...)
}
