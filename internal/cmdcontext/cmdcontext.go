// Package cmdcontext provides the "Context" struct which encapsulates all input values and
// related objects relevant to every execution of the "runner" command.
package cmdcontext

import (
	"io"
	"os"
)

// Context holds all structs which make up the runtime context of a single execution of the
// "runner" command. The most commonly-used structs here are the writers:
//   Out
//     The default writer subcommands should use when writing to standard out. This writer
//     is unlike any other as it will either ignore requests to write to standard out (by
//     default) or will honor them (when the 'verbose' flag is set).
//   VerboseOut
//     The writer subcommands should use when they want to force writing to standard out.
//     The 'version' command, for example, would be cumbersome if the 'verbose' flag were
//     necessary to generate output. Other commands will at times need to generate a minimum
//     level of output.
//   Err
//     The writer to use when generating output related to an error. Note that this writer
//     use the operating system's standard error output, instead of standard out.
type Context struct {
	Out        io.Writer
	VerboseOut io.Writer
	Err        io.Writer
}

// Init will instantiate a new Context object initialized with the state of the 'verbose'
// flag.
func (c *Context) Init(verbose bool) {
	if verbose {
		c.Out = os.Stdout
	} else {
		c.Out = NoopWriter{}
	}

	c.VerboseOut = os.Stdout
	c.Err = os.Stderr
}
