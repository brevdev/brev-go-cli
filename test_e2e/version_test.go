package test_e2e

import (
	"regexp"
	"testing"
)

func TestVersion(t *testing.T) {
	r := regexp.MustCompile(`
Current version: (\d.\d.\d)

You're up to date!

`)
	out, _ := Brev("version")
	if !r.MatchString(out) {
		t.Errorf("brev version => %s; expected => %s", out, r.String())
	}
}
