package env

import (
	"errors"
	"os"
	"syscall"

	"golang.org/x/term"

	"github.com/brevdev/brev-go-cli/internal/auth"
	"github.com/brevdev/brev-go-cli/internal/brev_api"
	"github.com/brevdev/brev-go-cli/internal/brev_ctx"
	"github.com/brevdev/brev-go-cli/internal/terminal"
)

func addVariable(name string, t *terminal.Terminal) error {

	t.Vprintf("Enter value for %s: ", name)

	bytepw, err := term.ReadPassword(int(syscall.Stdin))
	if err != nil {
		os.Exit(1)
	}
	value := string(bytepw)

	bar := t.NewProgressBar("Adding Variable "+t.Yellow(name), func() {})

	brevCtx, err := brev_ctx.New()
	if err != nil {
		return err
	}

	bar.AdvanceTo(40, t)
	project, err := brevCtx.Local.GetProject()
	if err != nil {
		return err
	}

	brevCtx.Remote.SetVariable(*project, name, value)

	finalStr := t.Green("\nVariable ") + t.Yellow("%s", name) + t.Green(" added to your project 🥞")
	bar.AdvanceTo(100, t)
	t.Vprint(finalStr)

	return nil
}

func removeVariable(name string, t *terminal.Terminal) error {

	bar1 := t.NewProgressBar("Removing Variable "+t.Yellow(name), func() {})

	brevCtx, err := brev_ctx.New()
	if err != nil {
		return err
	}

	bar1.AdvanceTo(40, t)
	project, err := brevCtx.Local.GetProject()
	if err != nil {
		return err
	}

	projVars, err := brevCtx.Remote.GetVariables(*project, &brev_ctx.GetVariablesOptions{
		Name: name,
	})
	if err != nil {
		return errors.New(t.Red("There isn't a variable in your project named %s.", name))
	}

	token, err := auth.GetToken()
	if err != nil {
		t.Errprint(err, "Failed to retrieve auth token")
		return err
	}
	brevAgent := brev_api.Agent{
		Key: token,
	}

	// Remove variable by ID
	_, err = brevAgent.RemoveVariable(projVars[0].Id)
	if err != nil {
		t.Errprint(err, "Couldn't remove the variable.")
		return err
	}

	finalStr := t.Green("\nVariable ") + t.Yellow("%s", name) + t.Green(" removed from your project 🥞")
	bar1.AdvanceTo(100, t)
	t.Vprint(finalStr)

	return nil
}
