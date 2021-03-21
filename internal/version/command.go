package version

import (
	"fmt"
	"time"

	"github.com/enescakir/emoji"
	"github.com/schollz/progressbar/v3"
	"github.com/spf13/cobra"

	"github.com/brevdev/brev-go-cli/internal/cmdcontext"
)

func NewCmdVersion(context *cmdcontext.Context) *cobra.Command {
	cmd := &cobra.Command{
		Use: "version",
		RunE: func(cmd *cobra.Command, args []string) error {
			version, err := buildVersionString(context)
			if err != nil {
				context.PrintErr("Failed to determine version", err)
				return err
			}
			fmt.Fprintln(context.VerboseOut, version)

			bar1 := newProgressBar("Doing things!", "1", "3", func() {
				fmt.Println("\n\nCompleted step 1!")
			})

			for i := 0; i < 1000; i++ {
				bar1.Add(1)
				time.Sleep(5 * time.Millisecond)
			}
			//fmt.Println("\n\nCompleted step 1!")
			fmt.Println("\nStep 3 will replace step 2 on its completion...\n")
			bar2 := newProgressBar("Even more things!", "2", "3", func() {})
			bar3 := newProgressBar("Naderrrrrrrr", "3", "3", func() {
				fmt.Printf("\n\nNow %v that %v\n", emoji.Ship, emoji.PileOfPoo)
			})
			for i := 0; i < 1000; i++ {
				bar2.Add(1)
				time.Sleep(3 * time.Millisecond)
			}
			for i := 0; i < 1000; i++ {
				bar3.Add(1)
				time.Sleep(1 * time.Millisecond)
			}

			return nil
		},
	}
	return cmd
}

func newProgressBar(description string, step string, stepTotal string, onComplete func()) *progressbar.ProgressBar {
	bar := progressbar.NewOptions(1000,
		progressbar.OptionOnCompletion(onComplete),
		progressbar.OptionEnableColorCodes(true),
		progressbar.OptionShowBytes(true),
		progressbar.OptionSetWidth(15),
		progressbar.OptionSetDescription(fmt.Sprintf("[cyan][%s/%s][reset] %s", step, stepTotal, description)),
		progressbar.OptionSetTheme(progressbar.Theme{
			Saucer:        "[green]=[reset]",
			SaucerHead:    "[green]>[reset]",
			SaucerPadding: " ",
			BarStart:      "[",
			BarEnd:        "]",
		}))
	return bar
}
