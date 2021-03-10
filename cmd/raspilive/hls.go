package main

import (
	"log"

	"github.com/jaredpetersen/raspilive/internal/ffmpeg/hls"
	"github.com/jaredpetersen/raspilive/internal/raspivid"
	"github.com/jaredpetersen/raspilive/internal/server"
	"github.com/jaredpetersen/raspilive/internal/video"
	"github.com/spf13/cobra"
)

// HlsCfg represents the configuration for DASH.
type HlsCfg struct {
	Width          int
	Height         int
	Fps            int
	HorizontalFlip bool
	VerticalFlip   bool
	Port           int
	Directory      string
	SegmentTime    int // Segment length target duration in seconds
	PlaylistSize   int // Maximum number of playlist entries
	StorageSize    int // Maximum number of unreferenced segments to keep on disk before removal
}

// HlsCmd is a HLS command for Cobra
var HlsCmd = &cobra.Command{
	Use:   "hls",
	Short: "Stream video using HLS",
	Long:  "Stream video using HLS",
}

func init() {
	hlsCfg := DashCfg{}

	HlsCmd.Flags().IntVar(&hlsCfg.Width, "width", 1920, "Video width")

	HlsCmd.Flags().IntVar(&hlsCfg.Height, "height", 1080, "Video height")

	HlsCmd.Flags().IntVar(&hlsCfg.Fps, "fps", 30, "Video framerate")

	HlsCmd.Flags().BoolVar(&hlsCfg.HorizontalFlip, "horizontal-flip", false, "Horizontally flip video")

	HlsCmd.Flags().BoolVar(&hlsCfg.VerticalFlip, "vertical-flip", false, "Vertically flip video")

	HlsCmd.Flags().IntVar(&hlsCfg.Port, "port", 0, "Static file server port (required)")
	HlsCmd.MarkFlagRequired("port")

	HlsCmd.Flags().StringVar(&hlsCfg.Directory, "directory", "", "Static file server directory (required)")
	HlsCmd.MarkFlagRequired("directory")

	HlsCmd.Flags().IntVar(&hlsCfg.SegmentTime, "segment-time", 0, "Segment length target duration in seconds")

	HlsCmd.Flags().IntVar(&hlsCfg.PlaylistSize, "playlist-size", 0, "Maximum number of playlist entries")

	HlsCmd.Flags().IntVar(&hlsCfg.StorageSize, "storage-size", 0, "Maximum number of unreferenced segments to keep on disk before removal")

	HlsCmd.Flags().SortFlags = false

	HlsCmd.Run = func(cmd *cobra.Command, args []string) {
		raspiOptions := raspivid.Options{
			Width:          hlsCfg.Width,
			Height:         hlsCfg.Height,
			Fps:            hlsCfg.Fps,
			HorizontalFlip: hlsCfg.HorizontalFlip,
			VerticalFlip:   hlsCfg.VerticalFlip,
		}
		raspiStream, err := raspivid.NewStream(raspiOptions)
		if err != nil {
			log.Fatal(err)
		}

		muxer := hls.Muxer{
			Directory: hlsCfg.Directory,
			Options: hls.Options{
				Fps:          hlsCfg.Fps,
				SegmentTime:  hlsCfg.SegmentTime,
				PlaylistSize: hlsCfg.PlaylistSize,
				StorageSize:  hlsCfg.StorageSize,
			},
		}
		server, err := server.NewStatic(hlsCfg.Port, hlsCfg.Directory)
		if err != nil {
			log.Fatal(err)
		}

		video.MuxAndServe(*raspiStream, &muxer, server)
	}
}
