package hls

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"

	"github.com/jaredpetersen/raspilive/internal/ffmpeg/hls"
	"github.com/jaredpetersen/raspilive/internal/raspivid"
	"github.com/jaredpetersen/raspilive/internal/server"
	"github.com/spf13/cobra"
)

// Cfg represents the configuration for DASH.
type Cfg struct {
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

// Cmd is a HLS command for Cobra.
var Cmd = &cobra.Command{
	Use:   "hls",
	Short: "Stream video using HLS",
	Long:  "Stream video using HLS",
}

func init() {
	cfg := Cfg{}

	Cmd.Flags().IntVar(&cfg.Width, "width", 1920, "video width")

	Cmd.Flags().IntVar(&cfg.Height, "height", 1080, "video height")

	Cmd.Flags().IntVar(&cfg.Fps, "fps", 30, "video framerate")

	Cmd.Flags().BoolVar(&cfg.HorizontalFlip, "horizontal-flip", false, "horizontally flip video")

	Cmd.Flags().BoolVar(&cfg.VerticalFlip, "vertical-flip", false, "vertically flip video")

	Cmd.Flags().IntVar(&cfg.Port, "port", 0, "static file server port (required)")
	Cmd.MarkFlagRequired("port")

	Cmd.Flags().StringVar(&cfg.Directory, "directory", "", "static file server directory (required)")
	Cmd.MarkFlagRequired("directory")

	Cmd.Flags().IntVar(&cfg.SegmentTime, "segment-time", 0, "segment length target duration in seconds")

	Cmd.Flags().IntVar(&cfg.PlaylistSize, "playlist-size", 0, "maximum number of playlist entries")

	Cmd.Flags().IntVar(&cfg.StorageSize, "storage-size", 0, "maximum number of unreferenced segments to keep on disk before removal")

	Cmd.Flags().SortFlags = false

	Cmd.Run = func(cmd *cobra.Command, args []string) {
		streamHls(cfg)
	}
}

func streamHls(cfg Cfg) {
	raspiStream := setupRaspiStream(cfg)
	muxer := setupMuxer(cfg)
	fileServer := setupServer(cfg)

	// Set up a channel for exiting
	stop := make(chan struct{})
	setupOsStopper(stop)

	// Serve files generated by the video stream
	go func() {
		serve(fileServer)
		stop <- struct{}{}
	}()

	// Stream video
	go func() {
		mux(raspiStream, muxer)
		stop <- struct{}{}
	}()

	// Wait for a stop signal
	<-stop

	log.Println("Shutting down")

	raspiStream.Video.Close()
	fileServer.Shutdown(context.Background())
}

func setupRaspiStream(cfg Cfg) *raspivid.Stream {
	raspiOptions := raspivid.Options{
		Width:          cfg.Width,
		Height:         cfg.Height,
		Fps:            cfg.Fps,
		HorizontalFlip: cfg.HorizontalFlip,
		VerticalFlip:   cfg.VerticalFlip,
	}

	raspiStream, err := raspivid.NewStream(raspiOptions)
	if err != nil {
		log.Fatal(err)
	}

	return raspiStream
}

func setupMuxer(cfg Cfg) *hls.Muxer {
	return &hls.Muxer{
		Directory: cfg.Directory,
		Options: hls.Options{
			Fps:          cfg.Fps,
			SegmentTime:  cfg.SegmentTime,
			PlaylistSize: cfg.PlaylistSize,
			StorageSize:  cfg.StorageSize,
		},
	}
}

func setupServer(cfg Cfg) *http.Server {
	server, err := server.NewStatic(cfg.Port, cfg.Directory)
	if err != nil {
		log.Fatal(err)
	}

	return server
}

func setupOsStopper(stop chan struct{}) {
	// Set up a channel for OS signals so that we can quit gracefully if the user terminates the program
	// Once we get this signal, sent a message to the stop channel
	osStop := make(chan os.Signal, 1)
	signal.Notify(osStop, os.Interrupt, os.Kill)

	go func() {
		<-osStop
		stop <- struct{}{}
	}()
}

func serve(server *http.Server) {
	log.Println("Server started", server.Addr)

	if err := server.ListenAndServe(); err != nil {
		if err != http.ErrServerClosed {
			log.Println("Encountered an error while serving vide", err)
		}
	}
}

func mux(raspiStream *raspivid.Stream, muxer *hls.Muxer) {
	if err := muxer.Mux(raspiStream.Video); err != nil {
		log.Fatal(err)
	}
	if err := raspiStream.Start(); err != nil {
		log.Fatal(err)
	}
	if err := muxer.Wait(); err != nil {
		log.Fatal(err)
	}
	if err := raspiStream.Wait(); err != nil {
		log.Fatal(err)
	}
}
