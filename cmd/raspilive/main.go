package main

import (
	"os"
	"time"

	"github.com/jaredpetersen/raspilive/cmd/raspilive/dash"
	"github.com/jaredpetersen/raspilive/cmd/raspilive/hls"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

// TODO research application logging

func main() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.RFC3339})
	cobra.EnableCommandSorting = false

	var debug bool

	rootCmd := &cobra.Command{
		Use:   "raspilive",
		Short: "raspilive is a livestreaming tool for the Raspberry Pi Camera Module",
		Long: "raspilive is a livestreaming tool for the Raspberry Pi Camera Module\n\n" +
			"For more information visit https://github.com/jaredpetersen/raspilive",
		Version: "1.0.0",
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			setLogLevel(debug)
		},
	}

	rootCmd.AddCommand(hls.Cmd)
	rootCmd.AddCommand(dash.Cmd)

	rootCmd.PersistentFlags().BoolVar(&debug, "debug", false, "enable debug logging")

	rootCmd.Execute()
}

func setLogLevel(debug bool) {
	if debug {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	}
}
