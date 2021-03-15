package version

import (
	"testing"

	"github.com/brevdev/brev-go-cli/internal/cmdcontext"
)

func TestBuildVersionString(t *testing.T) {
	contextStub := &cmdcontext.Context{}

    want := "unknown"
    got, err := buildVersionString(contextStub)

    if want != got || err != nil {
        t.Fatalf(`buildVersionString() = %q, %v, want match for %#q, nil`, got, err, want)
    }
}