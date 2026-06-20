package platform

import "github.com/Dandi-Pangestu/switchic/internal/util"

// Get returns the adapter for the named platform.
// Unknown names produce ErrUnknownPlatform so callers can surface a helpful message.
func Get(name string) (Adapter, error) {
	switch name {
	case "claude":
		return Claude{}, nil
	case "github-copilot":
		return Copilot{}, nil
	case "kiro":
		return Kiro{}, nil
	default:
		return nil, util.Wrap(util.ErrUnknownPlatform, "platform %q", name)
	}
}

// Available returns the list of platform names this build supports.
func Available() []string {
	return []string{"claude", "github-copilot", "kiro"}
}
