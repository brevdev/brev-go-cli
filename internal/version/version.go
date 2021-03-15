package version

import (
	"github.com/brevdev/brev-go-cli/internal/cmdcontext"
	"github.com/brevdev/brev-go-cli/internal/config"
)

// buildVersionString returns the full output for the 'version' command. The result of
// this function may be immediately printed to standard out.
func buildVersionString(context *cmdcontext.Context) (string, error) {
	return config.GetVersion(), nil
}
