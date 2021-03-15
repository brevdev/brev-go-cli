package cmdcontext

import (
	"fmt"
	"os"
	"testing"
)

func TestInitVerbose(t *testing.T) {
	context := &Context{}

	// verbose = true
	context.Init(true)

	want := os.Stdout
	got := context.Out

	if want != got {
		t.Errorf(`context.Out = %v, want %v`, got, want)
	}
}

func TestInitDefault(t *testing.T) {
	context := &Context{}

	// verbose = false
	context.Init(false)

	want := NoopWriter{}
	got := context.Out

	if !typeIsEqual(want, got) {
		t.Errorf(`context.Out = %v, want %v`, got, want)
	}
}

func typeIsEqual(a, b interface{}) bool {
    return fmt.Sprintf("%T", a) == fmt.Sprintf("%T", b)
}