package cmd

import (
	"github.com/sejamuchhal/email-reminder/background"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(newBackgroundCmd())
}

func newBackgroundCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "background",
		Short: "Starts the background worker",
		Run:   backgroundCmd,
	}
	return cmd
}

func backgroundCmd(cmd *cobra.Command, args []string) {
	background.Run()
}
