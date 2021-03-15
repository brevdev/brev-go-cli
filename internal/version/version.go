package version

import (
	"github.com/brevdev/brev-go-cli/internal/cmdcontext"
)

// buildVersionString returns the full output for the 'version' command. The result of
// this function may be immediately printed to standard out.
func buildVersionString(context *cmdcontext.Context) (string, error) {
	return "1.0", nil
}
