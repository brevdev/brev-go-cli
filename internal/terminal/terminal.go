package terminal

import (
	"fmt"
	"io"
	"os"

	"github.com/fatih/color"
)

type Terminal struct {
	out     io.Writer
	verbose io.Writer
	err     io.Writer

	Green  func(format string, a ...interface{}) string
	Yellow func(format string, a ...interface{}) string
	Red    func(format string, a ...interface{}) string
}

func (t *Terminal) Init(verbose bool) {
	var out io.Writer
	if verbose {
		out = os.Stdout
	} else {
		out = silentWriter{}
	}

	t.out = out
	t.verbose = os.Stdout
	t.err = os.Stderr

	t.Green = color.New(color.FgGreen).SprintfFunc()
	t.Yellow = color.New(color.FgYellow).SprintfFunc()
	t.Red = color.New(color.FgRed).SprintfFunc()
}

func (t *Terminal) Print(a string) {
	fmt.Fprintln(t.out, a)
}

func (t *Terminal) Printf(format string, a ...interface{}) {
	fmt.Fprintf(t.out, format, a)
}

func (t *Terminal) Vprint(a string) {
	fmt.Fprintln(t.verbose, a)
}

func (t *Terminal) Vprintf(format string, a ...interface{}) {
	fmt.Fprintf(t.verbose, format, a)
}

func (t *Terminal) Eprint(a string) {
	fmt.Fprintln(t.err, a)
}

func (t *Terminal) Errprint(err error, a ...interface{}) {
	if a != nil {
		t.Vprint(t.Red("Error: " + fmt.Sprint(a)))
	}
	t.Eprint(err.Error())
}

func (t *Terminal) Errprintf(err error, format string, a ...interface{}) {
	if a != nil {
		t.Vprintf(t.Red("Error: " + fmt.Sprintf(format, a)))
	}
	t.Eprint(err.Error())
}

type silentWriter struct{}

func (w silentWriter) Write(p []byte) (n int, err error) {
	return 0, nil
}
