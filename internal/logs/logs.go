package logs

import (
	"os"
	"os/signal"

	"github.com/brevdev/brev-go-cli/internal/brev_ctx"
	"github.com/brevdev/brev-go-cli/internal/terminal"
)

func LogTask(task string, t *terminal.Terminal) error {

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	brevCtx, err := brev_ctx.New()
	if err != nil {
		return err
	}

	proj, err := brevCtx.Local.GetProject()
	if err != nil {
		return err
	}

	logChan := brevCtx.Remote.TailLogs(proj.InstanceId, task)
	for {
		select {
		case log := <-logChan:
			t.Print(log)
		case <-interrupt:
			// TODO: cleanly close ws client
			return nil
		}

	}
}
