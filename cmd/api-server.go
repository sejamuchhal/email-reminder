package cmd

import (
	"github.com/spf13/cobra"
	"github.com/sejamuchhal/email-reminder/api"

)

func init() {
	rootCmd.AddCommand(newApiServerCmd())
}

func newApiServerCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "api-server",
		Short: "Starts the API server",
		Run: apiServerCmd,
	}
	return cmd
}

func apiServerCmd(cmd *cobra.Command, args []string) {
	api.StartServer()
}