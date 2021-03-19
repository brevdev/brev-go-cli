// Package cmdcontext provides the "Context" struct which encapsulates all input values and
// related objects relevant to every execution of the "runner" command.
package cmdcontext

import (
	"fmt"
	"io"
	"os"

	"github.com/spf13/cobra"
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

func (c *Context) PrintErr(message string, err error) {
	fmt.Fprintln(c.VerboseOut, message)
	fmt.Fprintln(c.Err, err.Error())
}

// InvokeParentPersistentPreRun executes the immediate parent command's
// PersistentPreRunE and PersistentPreRun functions, in that order. If
// an error is returned from PersistentPreRunE, it is immediately returned.
//
// TODO: reverse walk up command tree? would need to ensure no one parent is invoked multiple times
func InvokeParentPersistentPreRun(cmd *cobra.Command, args []string) error {
	parentCmd := cmd.Parent()
	if parentCmd == nil {
		return nil
	}

	var err error

	// Invoke PersistentPreRunE, returning an error if one occurs
	// If no error is returned, proceed with PersistentPreRun
	parentPersistentPreRunE := parentCmd.PersistentPreRunE
	if parentPersistentPreRunE != nil {
		err = parentPersistentPreRunE(parentCmd, args)
	}
	if err != nil {
		return err
	}

	// Invoke PersistentPreRun
	parentPersistentPreRun := parentCmd.PersistentPreRun
	if parentPersistentPreRun != nil {
		parentPersistentPreRun(parentCmd, args)
	}

	return nil
}
