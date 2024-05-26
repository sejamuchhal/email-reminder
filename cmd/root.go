package cmd

import (
	"os"

	"github.com/spf13/cobra"
)


var rootCmd = &cobra.Command{
	Use:   "email-reminder",
	Short: "Email reminder application",
	Long: `Email Reminder is a command line application that allows you to schedule 
email reminders. It includes an API server and a background worker. 

Use the 'api-server' command to start the API server and the 'background' 
command to start the background worker.`,
}


func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
