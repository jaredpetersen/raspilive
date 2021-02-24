package hls

import (
	"io"
	"os/exec"
	"path"
	"strconv"
	"strings"
)

// Muxer does the thing
//
// Ffmpeg will step in and use its own defaults if a value is not provided.
type Muxer struct {
	Directory    string
	Fps          int // Framerate of the output video
	SegmentTime  int // Segment length target duration in seconds
	PlaylistSize int // Maximum number of playlist entries
	StorageSize  int // Maximum number of unreferenced segments to keep on disk before removal
	cmd          *exec.Cmd
}

var execCommand = exec.Command

func (mx *Muxer) Start(video io.ReadCloser) error {
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

	if mx.Fps != 0 {
		args = append(args, "-r", strconv.Itoa(mx.Fps))
	}

	if mx.SegmentTime != 0 {
		args = append(args, "-hls_time", strconv.Itoa(mx.SegmentTime))
		hlsFlags = append(hlsFlags, "split_by_time")
	}

	if mx.PlaylistSize != 0 {
		args = append(args, "-hls_list_size", strconv.Itoa(mx.PlaylistSize))
	}

	if mx.StorageSize != 0 {
		args = append(args, "-hls_delete_threshold", strconv.Itoa(mx.StorageSize))
		hlsFlags = append(hlsFlags, "delete_segments")
	}

	args = append(args, "-hls_flags", strings.Join(hlsFlags, "+"), path.Join(mx.Directory, "livestream.m3u8"))

	mx.cmd = execCommand("ffmpeg", args...)
	mx.cmd.Stdin = video

	return mx.cmd.Start()
}

// Wait waits for the video stream to finish processing.
//
// The mux operation must have been started by Start.
func (mx *Muxer) Wait() error {
	return mx.cmd.Wait()
}
