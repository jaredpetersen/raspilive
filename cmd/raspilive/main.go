package main

import (
	"os"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

const serverShutdownDeadline = 10 * time.Second

// VideoCfg represents the video configuration options
type VideoCfg struct {
	Width          int
	Height         int
	Fps            int
	HorizontalFlip bool
	VerticalFlip   bool
}

func main() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.RFC3339})

	cobra.EnableCommandSorting = false

	var debug bool
	var video VideoCfg

	rootCmd := &cobra.Command{
		Use:   "raspilive",
		Short: "raspilive streams video from the Raspberry Pi Camera Module to the web",
		Long: "raspilive streams video from the Raspberry Pi Camera Module to the web\n\n" +
			"For more information visit https://github.com/jaredpetersen/raspilive",
		Version: "1.0.0",
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			setLogLevel(debug)
		},
	}

	rootCmd.AddCommand(newHlsCmd(&video))
	rootCmd.AddCommand(newDashCmd(&video))

	rootCmd.PersistentFlags().IntVar(&video.Width, "width", 1280, "video width")
	rootCmd.PersistentFlags().IntVar(&video.Height, "height", 720, "video height")
	rootCmd.PersistentFlags().IntVar(&video.Fps, "fps", 30, "video framerate")
	rootCmd.PersistentFlags().BoolVar(&video.HorizontalFlip, "horizontal-flip", false, "horizontally flip video")
	rootCmd.PersistentFlags().BoolVar(&video.VerticalFlip, "vertical-flip", false, "vertically flip video")
	rootCmd.PersistentFlags().BoolVar(&debug, "debug", false, "enable debug logging")

	rootCmd.Execute()
}

func setLogLevel(debug bool) {
	var logLevel zerolog.Level
	if debug {
		logLevel = zerolog.DebugLevel
	} else {
		logLevel = zerolog.InfoLevel
	}

	zerolog.SetGlobalLevel(logLevel)
}
