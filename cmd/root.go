package cmd

import (
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "raspilive",
	Short: "raspilive is a livestreaming tool for the Raspberry Pi Camera Module",
	Long: "raspilive is a livestreaming tool for the Raspberry Pi Camera Module\n\n" +
		"For more information visit https://github.com/jaredpetersen/raspilive",
	Version: "1.0.0",
}

// Execute executes the root command.
func Execute() error {
	return rootCmd.Execute()
}
