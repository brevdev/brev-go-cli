package test_e2e

import (
	"bytes"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path"
	"strings"
	"testing"
)

var bin = flag.String("brev-bin", "brev", "brev binary path")
var abs = ""

func TestMain(m *testing.M) {
	beforeAll()
	code := m.Run()
	afterAll()
	os.Exit(code)
}

func beforeAll() {

}

func afterAll() {

}

// Brev executes a brev command, e.g.: Brev("endpoint run --v")
func Brev(cmdArgs string) (string, error) {
	var out bytes.Buffer
	cmd := exec.Command(getAbsBin(), strings.Fields(cmdArgs)...)
	cmd.Stdout = &out

	err := cmd.Run()

	return out.String(), err
}

func getAbsBin() string {
	if abs != "" {
		return abs
	}
	if path.IsAbs(*bin) {
		abs = *bin
	}
	wd, err := os.Getwd()
	if err != nil {
		panic(fmt.Sprintf("Getwd failed: %s", err))
	}
	abs := path.Clean(path.Join(wd, "..", *bin))
	return abs
}
