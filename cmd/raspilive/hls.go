package main

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/jaredpetersen/raspilive/internal/ffmpeg/hls"
	"github.com/jaredpetersen/raspilive/internal/raspivid"
	"github.com/jaredpetersen/raspilive/internal/server"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

// HlsCfg represents the HLS configuration options
type HlsCfg struct {
	Video        VideoCfg
	Port         int
	Directory    string
	TLSCert      string
	TLSKey       string
	SegmentType  string // Format of the video segment
	SegmentTime  int    // Segment length target duration in seconds
	PlaylistSize int    // Maximum number of playlist entries
	StorageSize  int    // Maximum number of unreferenced segments to keep on disk before removal
}

func newHlsCmd(video VideoCfg) *cobra.Command {
	cfg := HlsCfg{
		Video: video,
	}

	cmd := &cobra.Command{
		Use:   "hls",
		Short: "Stream video using HLS",
		Long:  "Stream video using HLS",
	}

	cmd.Flags().IntVar(&cfg.Port, "port", 0, "static file server port")
	cmd.MarkFlagRequired("port")

	cmd.Flags().StringVar(&cfg.Directory, "directory", "", "static file server directory")

	cmd.Flags().StringVar(&cfg.TLSCert, "tls-cert", "", "static file server TLS certificate")

	cmd.Flags().StringVar(&cfg.TLSKey, "tls-key", "", "static file server TLS key")

	cmd.Flags().StringVar(&cfg.SegmentType, "segment-type", "", "format of the video segments (valid [\"mpegts\", \"fmp4\"], default \"mpegts\")")

	cmd.Flags().IntVar(&cfg.SegmentTime, "segment-time", 2, "target segment duration in seconds")

	cmd.Flags().IntVar(&cfg.PlaylistSize, "playlist-size", 10, "maximum number of playlist entries")

	cmd.Flags().IntVar(&cfg.StorageSize, "storage-size", 1, "maximum number of unreferenced segments to keep on disk before removal")

	cmd.Flags().SortFlags = false

	cmd.Run = func(cmd *cobra.Command, args []string) {
		streamHls(cfg)
	}

	cmd.PreRun = func(cmd *cobra.Command, args []string) {
		isValidCfg := isValidHlsCfg(cfg)
		if !isValidCfg {
			cmd.Usage()
			os.Exit(1)
		}
	}

	return cmd
}

func isValidHlsCfg(cfg HlsCfg) bool {
	isValidCfg := true

	segmentType := strings.ToLower(cfg.SegmentType)
	validSegmentType := segmentType == "" || segmentType == "mpegts" || segmentType == "fmp4"

	if !validSegmentType {
		fmt.Printf("Error: invalid value \"%s\" for flag \"segment-type\"\n", cfg.SegmentType)
		isValidCfg = false
	}

	return isValidCfg
}

func streamHls(cfg HlsCfg) {
	// Set up raspivid stream
	raspiOptions := raspivid.Options{
		Width:          cfg.Video.Width,
		Height:         cfg.Video.Height,
		Fps:            cfg.Video.Fps,
		HorizontalFlip: cfg.Video.HorizontalFlip,
		VerticalFlip:   cfg.Video.VerticalFlip,
	}
	raspiStream, err := raspivid.NewStream(raspiOptions)
	if err != nil {
		log.Fatal().Msg("Encountered an error streaming video from the Raspberry Pi Camera Module")
	}

	// Set up HLS muxer
	muxer := hls.Muxer{
		Directory: cfg.Directory,
		Options: hls.Options{
			Fps:          cfg.Video.Fps,
			SegmentTime:  cfg.SegmentTime,
			PlaylistSize: cfg.PlaylistSize,
			StorageSize:  cfg.StorageSize,
		},
	}

	// Set up static file server
	srv := server.Static{
		Port:      cfg.Port,
		Directory: cfg.Directory,
		Cert:      cfg.TLSCert,
		Key:       cfg.TLSKey,
	}

	// Set up a channel for exiting
	stop := make(chan struct{})
	osStopper(stop)

	// Serve files generated by the video stream
	go func() {
		err := srv.ListenAndServe()
		if errors.Is(err, server.ErrInvalidDirectory) {
			log.Fatal().Msg("Directory does not exist")
		}
		if err != nil {
			log.Debug().Err(err).Msg("Encountered an error serving video")
			log.Fatal().Msg("Encountered an error serving video")
		}
		stop <- struct{}{}
	}()

	// Stream video
	go func() {
		if err := muxHls(raspiStream, &muxer); err != nil {
			log.Fatal().Msg("Encountered an error streaming/muxing video")
		}
		stop <- struct{}{}
	}()

	// Wait for a stop signal
	<-stop

	log.Info().Msg("Shutting down")

	raspiStream.Video.Close()
	srv.Shutdown(serverShutdownDeadline)
}

func muxHls(raspiStream *raspivid.Stream, muxer *hls.Muxer) error {
	if err := muxer.Mux(raspiStream.Video); err != nil {
		log.Debug().Err(err).Msg("Encountered an error starting video mux")
		return err
	}
	log.Debug().Str("cmd", muxer.String()).Msg("Started ffmpeg muxer")

	if err := raspiStream.Start(); err != nil {
		log.Debug().Err(err).Msg("Encountered an error starting video stream")
		return err
	}
	log.Debug().Str("cmd", raspiStream.String()).Msg("Started raspivid")

	if err := muxer.Wait(); err != nil {
		log.Debug().Err(err).Msg("Encountered an error waiting for video mux")
		return err
	}

	if err := raspiStream.Wait(); err != nil {
		log.Debug().Err(err).Msg("Encountered an error waiting for video stream")
		return err
	}

	return nil
}
