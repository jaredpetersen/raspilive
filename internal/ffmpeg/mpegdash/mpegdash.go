package mpegdash

import (
	"io"
	"os/exec"
	"path"
	"strconv"
)

// Muxer represents a video transformation operation being prepared or run.
//
// A Muxer cannot be reused after calling its Start method.
type Muxer struct {
	cmd *exec.Cmd
}

// Options represents video muxing options for MPEG-DASH.
//
// Ffmpeg will step in and use its own defaults if a value is not provided.
type Options struct {
	Fps          int // Framerate of the output video
	SegmentTime  int // Segment length target duration in seconds
	PlaylistSize int // Maximum number of playlist entries
	StorageSize  int // Maximum number of unreferenced segments to keep on disk before removal
}

var execCommand = exec.Command

// MpegDash prepares to mux a video stream into MPEG-DASH.
func MpegDash(inputStream io.ReadCloser, directory string, options Options) *Muxer {
	args := []string{
		"-codec", "copy",
		"-f", "dash",
		"-re",
		"-an",
		"-init_seg_name", "init.m4s",
		"-media_seg_name", "$Time$-$Number$.m4s",
	}

	if options.Fps != 0 {
		args = append(args, "-r", strconv.Itoa(options.Fps))
	}

	if options.SegmentTime != 0 {
		args = append(args, "-seg_duration", strconv.Itoa(options.SegmentTime))
	}

	if options.PlaylistSize != 0 {
		args = append(args, "-window_size", strconv.Itoa(options.PlaylistSize))
	}

	if options.StorageSize != 0 {
		args = append(args, "-extra_window_size", strconv.Itoa(options.StorageSize))
	}

	args = append(args, path.Join(directory, "livestream.mpd"))

	ffmpegCommand := execCommand("ffmpeg", args...)
	ffmpegCommand.Stdin = inputStream

	return &Muxer{
		cmd: ffmpegCommand,
	}
}

// Start muxes the prepared video stream into MPEG-DASH.
func (muxer *Muxer) Start() error {
	return muxer.cmd.Start()
}

// Wait waits for the video stream to finish processing.
//
// The mux operation must have been started by Start.
func (muxer *Muxer) Wait() error {
	return muxer.cmd.Wait()
}
