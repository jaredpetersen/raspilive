package dash

import (
	"errors"
	"io"
	"os/exec"
	"path"
	"strconv"
	"strings"
)

// Muxer represents a video transformation operation being prepared or run.
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
		"-f", "dash",
		"-re",
		"-an",
		"-init_seg_name", "init.$ext$",
		"-media_seg_name", "$Time$-$Number$.$ext$",
	}

	segmentType := strings.ToLower(muxer.SegmentType)
	if segmentType == "" || segmentType == "auto" {
		args = append(args, "-dash_segment_type", "auto")
	} else if segmentType == "mp4" {
		args = append(args, "-dash_segment_type", "mp4")
	} else if segmentType == "webm" {
		args = append(args, "-dash_segment_type", "webm")
	} else {
		return errors.New("ffmpeg dash: invalid segment type")
	}

	if muxer.Fps != 0 {
		args = append(args, "-r", strconv.Itoa(muxer.Fps))
	}

	if muxer.SegmentTime != 0 {
		args = append(args, "-seg_duration", strconv.Itoa(muxer.SegmentTime))
	}

	if muxer.PlaylistSize != 0 {
		args = append(args, "-window_size", strconv.Itoa(muxer.PlaylistSize))
	}

	if muxer.StorageSize != 0 {
		args = append(args, "-extra_window_size", strconv.Itoa(muxer.StorageSize))
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
