package dash

import (
	"errors"
	"io"
	"os/exec"
	"path"
	"strconv"
	"strings"
)

// Options represents ways that Ffmpeg may be configured to mux video to DASH.
//
// Ffmpeg will step in and use its own defaults if a value is not provided.
type Options struct {
	Fps          int    // Framerate of the output video
	SegmentType  string // Format of the video segment
	SegmentTime  int    // Segment length target duration in seconds
	PlaylistSize int    // Maximum number of playlist entries
	StorageSize  int    // Maximum number of unreferenced segments to keep on disk before removal
}

// Muxer represents the DASH muxer.
type Muxer struct {
	Directory string
	Options   Options
	cmd       *exec.Cmd
}

var execCommand = exec.Command

// Mux begins muxing the video stream to the DASH format.
func (muxer *Muxer) Mux(video io.ReadCloser) error {
	args := []string{
		"-codec", "copy",
		"-f", "dash",
		"-re",
		"-an",
		"-init_seg_name", "init.$ext$",
		"-media_seg_name", "$Time$-$Number$.$ext$",
	}

	segmentType := strings.ToLower(muxer.Options.SegmentType)
	if segmentType == "" || segmentType == "auto" {
		args = append(args, "-dash_segment_type", "auto")
	} else if segmentType == "mp4" {
		args = append(args, "-dash_segment_type", "mp4")
	} else if segmentType == "webm" {
		args = append(args, "-dash_segment_type", "webm")
	} else {
		return errors.New("ffmpeg dash: invalid segment type")
	}

	if muxer.Options.Fps != 0 {
		args = append(args, "-r", strconv.Itoa(muxer.Options.Fps))
	}

	if muxer.Options.SegmentTime != 0 {
		args = append(args, "-seg_duration", strconv.Itoa(muxer.Options.SegmentTime))
	}

	if muxer.Options.PlaylistSize != 0 {
		args = append(args, "-window_size", strconv.Itoa(muxer.Options.PlaylistSize))
	}

	if muxer.Options.StorageSize != 0 {
		args = append(args, "-extra_window_size", strconv.Itoa(muxer.Options.StorageSize))
	}

	args = append(args, path.Join(muxer.Directory, "livestream.mpd"))

	muxer.cmd = execCommand("ffmpeg", args...)
	muxer.cmd.Stdin = video

	return muxer.cmd.Start()
}

// Wait waits for the video stream to finish processing.
//
// The mux operation must have been started by Start.
func (muxer *Muxer) Wait() error {
	if muxer.cmd == nil {
		return errors.New("ffmpeg dash: not started")
	}

	return muxer.cmd.Wait()
}
