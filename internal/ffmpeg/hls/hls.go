package hls

import (
	"io"
	"os/exec"
	"path"
	"strconv"
	"strings"
)

// Muxer represents a video transformation operation being prepared or run.
//
// A Muxer cannot be reused after calling its Start method.
type Muxer struct {
	cmd *exec.Cmd
}

// Options represents video muxing options for HLS.
//
// Ffmpeg will step in and use its own defaults if a value is not provided.
type Options struct {
	Fps         *int // Framerate of the output video
	Time        *int // Segment length target duration in seconds
	ListSize    *int // Maximum number of playlist entries
	StorageSize *int // Maximum number of unreferenced segments to keep on disk before removal
}

var execCommand = exec.Command

// Hls prepares to mux a video stream into HLS.
func Hls(inputStream io.ReadCloser, directory string, options Options) *Muxer {
	args := []string{
		"-codec", "copy",
		"-f", "hls",
		"-re",
		"-an",
		"-strftime", "1",
		"-hls_segment_filename", "%s-%%d.m4s",
		"-hls_segment_type", "fmp4",
	}
	hlsFlags := []string{"second_level_segment_index"}

	if options.Fps != nil {
		args = append(args, "-r", strconv.Itoa(*options.Fps))
	}

	if options.Time != nil {
		args = append(args, "-hls_time", strconv.Itoa(*options.Time))
		hlsFlags = append(hlsFlags, "split_by_time")
	}

	if options.ListSize != nil {
		args = append(args, "-hls_list_size", strconv.Itoa(*options.ListSize))
	}

	if options.StorageSize != nil {
		args = append(args, "-hls_delete_threshold", strconv.Itoa(*options.StorageSize))
		hlsFlags = append(hlsFlags, "delete_segments")
	}

	args = append(args, "-hls_flags", strings.Join(hlsFlags, "+"), path.Join(directory, "livestream.m3u8"))

	ffmpegCommand := execCommand("ffmpeg", args...)
	ffmpegCommand.Stdin = inputStream

	return &Muxer{
		cmd: ffmpegCommand,
	}
}

// Start muxes the prepared video stream into HLS.
func (muxer *Muxer) Start() error {
	return muxer.cmd.Start()
}

// Wait waits for the video stream to finish processing.
//
// The mux operation must have been started by Start.
func (muxer *Muxer) Wait() error {
	return muxer.cmd.Wait()
}
