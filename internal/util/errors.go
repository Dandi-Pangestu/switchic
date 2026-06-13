package util

import (
	"errors"
	"fmt"
)

// Sentinels returned by domain code so callers can branch on them with errors.Is.
var (
	ErrNotFound       = errors.New("not found")
	ErrAlreadyExists  = errors.New("already exists")
	ErrInvalidConfig  = errors.New("invalid config")
	ErrUnknownPlatform = errors.New("unknown platform")
)

// Wrap returns nil if err is nil; otherwise it adds context with %w semantics.
func Wrap(err error, format string, args ...any) error {
	if err == nil {
		return nil
	}
	return fmt.Errorf(format+": %w", append(args, err)...)
}
