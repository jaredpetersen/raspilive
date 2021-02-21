package main

import (
	"log"
	"os"
	"strings"

	"github.com/jaredpetersen/raspilive/internal/video/hls"
	"github.com/kelseyhightower/envconfig"
)

// TODO create server directory if it does not exist
// TODO research application logging

// VideoConfig represents the configuration for the video source.
type VideoConfig struct {
	Width          int
	Height         int
	Fps            int
	HorizontalFlip bool
	VerticalFlip   bool
}

// HlsConfig represents the configuraiton for HLS.
type HlsConfig struct {
	Port         int    `required:"true"`
	Directory    string `default:"./camera"`
	SegmentTime  int    // Segment length target duration in seconds
	PlaylistSize int    // Maximum number of playlist entries
	StorageSize  int    // Maximum number of unreferenced segments to keep on disk before removal
}

// Config represents the configuration for raspilive.
type Config struct {
	Mode  string `required:"true"`
	Video VideoConfig
	Hls   HlsConfig
}

func main() {
	var config Config
	err := envconfig.Process("raspilive", &config)
	if err != nil {
		log.Fatal(err.Error())
	}

	switch strings.ToUpper(config.Mode) {
	case "HLS":
		log.Println("Using HLS mode...")
		hls.ServeHls(config.Hls.Port, config.Hls.Directory, toHlsOptions(config.Video, config.Hls))
	default:
		os.Exit(1)
	}
}

func toHlsOptions(videoConfig VideoConfig, hlsConfig HlsConfig) hls.Options {
	return hls.Options{
		Width:          videoConfig.Width,
		Height:         videoConfig.Height,
		Fps:            videoConfig.Fps,
		HorizontalFlip: videoConfig.HorizontalFlip,
		VerticalFlip:   videoConfig.VerticalFlip,
		SegmentTime:    hlsConfig.SegmentTime,
		PlaylistSize:   hlsConfig.PlaylistSize,
		StorageSize:    hlsConfig.StorageSize,
	}
}
