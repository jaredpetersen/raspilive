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
		Short: "raspilive is a livestreaming tool for the Raspberry Pi Camera Module",
		Long: "raspilive is a livestreaming tool for the Raspberry Pi Camera Module\n\n" +
			"For more information visit https://github.com/jaredpetersen/raspilive",
		Version: "1.0.0",
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			setLogLevel(debug)
		},
	}

	rootCmd.AddCommand(newHlsCmd(video))
	rootCmd.AddCommand(newDashCmd(video))

	rootCmd.PersistentFlags().IntVar(&video.Width, "width", 1920, "video width")
	rootCmd.PersistentFlags().IntVar(&video.Height, "height", 1080, "video height")
	rootCmd.PersistentFlags().IntVar(&video.Fps, "fps", 30, "video framerate")
	rootCmd.PersistentFlags().BoolVar(&video.HorizontalFlip, "horizontal-flip", false, "horizontally flip video")
	rootCmd.PersistentFlags().BoolVar(&video.VerticalFlip, "vertical-flip", false, "vertically flip video")
	rootCmd.PersistentFlags().BoolVar(&debug, "debug", false, "enable debug logging")

	rootCmd.Execute()
}

func setLogLevel(debug bool) {
	if debug {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	}
}
