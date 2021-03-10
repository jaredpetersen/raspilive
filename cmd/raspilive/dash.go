package main

import (
	"log"

	"github.com/jaredpetersen/raspilive/internal/ffmpeg/dash"
	"github.com/jaredpetersen/raspilive/internal/raspivid"
	"github.com/jaredpetersen/raspilive/internal/server"
	"github.com/jaredpetersen/raspilive/internal/video"
	"github.com/spf13/cobra"
)

// DashCfg represents the configuration for DASH.
type DashCfg struct {
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

// DashCmd is a DASH command for Cobra
var DashCmd = &cobra.Command{
	Use:   "dash",
	Short: "Stream video using DASH",
	Long:  "Stream video using DASH",
}

func init() {
	dashCfg := DashCfg{}

	DashCmd.Flags().IntVar(&dashCfg.Width, "width", 1920, "Video width")

	DashCmd.Flags().IntVar(&dashCfg.Height, "height", 1080, "Video height")

	DashCmd.Flags().IntVar(&dashCfg.Fps, "fps", 30, "Video framerate")

	DashCmd.Flags().BoolVar(&dashCfg.HorizontalFlip, "horizontal-flip", false, "Horizontally flip video")

	DashCmd.Flags().BoolVar(&dashCfg.VerticalFlip, "vertical-flip", false, "Vertically flip video")

	DashCmd.Flags().IntVar(&dashCfg.Port, "port", 0, "Static file server port (required)")
	DashCmd.MarkFlagRequired("port")

	DashCmd.Flags().StringVar(&dashCfg.Directory, "directory", "", "Static file server directory (required)")
	DashCmd.MarkFlagRequired("directory")

	DashCmd.Flags().IntVar(&dashCfg.SegmentTime, "segment-time", 0, "Segment length target duration in seconds")

	DashCmd.Flags().IntVar(&dashCfg.PlaylistSize, "playlist-size", 0, "Maximum number of playlist entries")

	DashCmd.Flags().IntVar(&dashCfg.StorageSize, "storage-size", 0, "Maximum number of unreferenced segments to keep on disk before removal")

	DashCmd.Flags().SortFlags = false

	DashCmd.Run = func(cmd *cobra.Command, args []string) {
		raspiOptions := raspivid.Options{
			Width:          dashCfg.Width,
			Height:         dashCfg.Height,
			Fps:            dashCfg.Fps,
			HorizontalFlip: dashCfg.HorizontalFlip,
			VerticalFlip:   dashCfg.VerticalFlip,
		}
		raspiStream, err := raspivid.NewStream(raspiOptions)
		if err != nil {
			log.Fatal(err)
		}

		muxer := dash.Muxer{
			Directory: dashCfg.Directory,
			Options: dash.Options{
				Fps:          dashCfg.Fps,
				SegmentTime:  dashCfg.SegmentTime,
				PlaylistSize: dashCfg.PlaylistSize,
				StorageSize:  dashCfg.StorageSize,
			},
		}
		server, err := server.NewStatic(dashCfg.Port, dashCfg.Directory)
		if err != nil {
			log.Fatal(err)
		}

		video.MuxAndServe(*raspiStream, &muxer, server)
	}
}
