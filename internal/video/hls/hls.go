package hls

import (
	"sync"

	"github.com/jaredpetersen/raspilive/internal/ffmpeg/hls"
	"github.com/jaredpetersen/raspilive/internal/raspivid"
	"github.com/jaredpetersen/raspilive/internal/video/server"
)

// Options represents streaming configuration for HLS.
type Options struct {
	Width          int
	Height         int
	Fps            int
	HorizontalFlip bool
	VerticalFlip   bool
	SegmentTime    int // Segment length target duration in seconds
	PlaylistSize   int // Maximum number of playlist entries
	StorageSize    int // Maximum number of unreferenced segments to keep on disk before removal
}

// ServeHls starts a static file server and stream video from the Raspberry Pi camera module using the HLS format.
//
// This is a blocking operation.
func ServeHls(port int, directory string, options Options) {
	var wg sync.WaitGroup
	wg.Add(2)

	// Serve files generated by the video stream
	go func() {
		server.ServeFiles(port, directory)
		wg.Done()
	}()

	// Stream video
	go func() {
		streamHls(directory, options)
		wg.Done()
	}()

	wg.Wait()
}

// StreamHls streams video from the Raspberry Pi camera module and muxes it to HLS.
//
// This is a blocking operation that will not complete.
func streamHls(directory string, options Options) {
	// Pipe video stream from raspivid into ffmpeg
	raspivid := raspivid.Stream(toRaspividOptions(options))
	ffmpeg := hls.Hls(raspivid.Video, directory, toHlsOptions(options))

	// Start ffmpeg first so that it's ready to accept the stream
	ffmpeg.Start()
	raspivid.Start()
	raspivid.Wait()
	ffmpeg.Wait()
}

func toRaspividOptions(options Options) raspivid.Options {
	return raspivid.Options{
		Width:          options.Width,
		Height:         options.Height,
		Fps:            options.Fps,
		HorizontalFlip: options.HorizontalFlip,
		VerticalFlip:   options.VerticalFlip,
	}
}

func toHlsOptions(options Options) hls.Options {
	return hls.Options{
		Fps:          options.Fps,
		SegmentTime:  options.SegmentTime,
		PlaylistSize: options.PlaylistSize,
		StorageSize:  options.StorageSize,
	}
}
