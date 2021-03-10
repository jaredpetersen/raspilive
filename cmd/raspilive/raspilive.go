package main

import (
	"github.com/spf13/cobra"
)

// TODO research application logging

func main() {
	cobra.EnableCommandSorting = false

	rootCmd := &cobra.Command{
		Use:   "raspilive",
		Short: "raspilive is a livestreaming tool for the Raspberry Pi Camera Module",
		Long: "raspilive is a livestreaming tool for the Raspberry Pi Camera Module\n\n" +
			"For more information visit https://github.com/jaredpetersen/raspilive",
		Version: "1.0.0",
	}

	rootCmd.AddCommand(HlsCmd)
	rootCmd.AddCommand(DashCmd)

	rootCmd.Execute()
}
