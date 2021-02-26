package hls

import (
	"errors"
	"io"
	"os/exec"
	"path"
	"strconv"
	"strings"
)

// Muxer represents the HLS muxer.
//
// Ffmpeg will step in and use its own defaults if a value is not provided.
type Muxer struct {
	Directory    string
	Fps          int    // Framerate of the output video
	SegmentType  string // Format of the video segment
	SegmentTime  int    // Segment length target duration in seconds
	PlaylistSize int    // Maximum number of playlist entries
	StorageSize  int    // Maximum number of unreferenced segments to keep on disk before removal
	cmd          *exec.Cmd
}

var execCommand = exec.Command

// Start begins muxing the video stream to the HLS format.
func (muxer *Muxer) Start(video io.ReadCloser) error {
	args := []string{
		"-codec", "copy",
		"-f", "hls",
		"-re",
		"-an",
		"-strftime", "1",
	}
	hlsFlags := []string{"second_level_segment_index"}

	segmentType := strings.ToLower(muxer.SegmentType)
	if segmentType == "" || segmentType == "mpegts" {
		args = append(
			args,
			"-hls_segment_type", "mpegts",
			"-hls_segment_filename", "%s-%%d.ts")
	} else if segmentType == "fmp4" {
		args = append(
			args,
			"-hls_segment_type", "fmp4",
			"-hls_segment_filename", "%s-%%d.m4s")
	} else {
		return errors.New("ffmpeg dash: invalid segment type")
	}

	if muxer.Fps != 0 {
		args = append(args, "-r", strconv.Itoa(muxer.Fps))
	}

	if muxer.SegmentTime != 0 {
		args = append(args, "-hls_time", strconv.Itoa(muxer.SegmentTime))
		hlsFlags = append(hlsFlags, "split_by_time")
	}

	if muxer.PlaylistSize != 0 {
		args = append(args, "-hls_list_size", strconv.Itoa(muxer.PlaylistSize))
	}

	if muxer.StorageSize != 0 {
		args = append(args, "-hls_delete_threshold", strconv.Itoa(muxer.StorageSize))
		hlsFlags = append(hlsFlags, "delete_segments")
	}

	args = append(args, "-hls_flags", strings.Join(hlsFlags, "+"), path.Join(muxer.Directory, "livestream.m3u8"))

	muxer.cmd = execCommand("ffmpeg", args...)
	muxer.cmd.Stdin = video

	return muxer.cmd.Start()
}

// Wait blocks until the video stream is finished processing.
//
// The mux operation must have been started by Start.
func (muxer *Muxer) Wait() error {
	if muxer.cmd == nil {
		return errors.New("ffmpeg hls: not started")
	}

	return muxer.cmd.Wait()
}
