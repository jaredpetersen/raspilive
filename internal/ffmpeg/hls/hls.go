package hls

import (
	"errors"
	"io"
	"log"
	"os/exec"
	"path"
	"strconv"
	"strings"
)

// Options represents ways that Ffmpeg may be configured to mux video to HLS.
//
// Ffmpeg will step in and use its own defaults if a value is not provided.
type Options struct {
	Fps          int    // Framerate of the output video
	SegmentType  string // Format of the video segment
	SegmentTime  int    // Segment length target duration in seconds
	PlaylistSize int    // Maximum number of playlist entries
	StorageSize  int    // Maximum number of unreferenced segments to keep on disk before removal
}

// Muxer represents the HLS muxer.
type Muxer struct {
	Directory string
	Options   Options
	cmd       *exec.Cmd
}

var execCommand = exec.Command

// Mux begins muxing the video stream to the HLS format.
func (muxer *Muxer) Mux(video io.ReadCloser) error {
	args := []string{
		"-re",
		"-i", "pipe:0",
		"-codec", "copy",
		"-f", "hls",
		"-an",
		"-strftime", "1",
	}
	hlsFlags := []string{"second_level_segment_index"}

	segmentType := strings.ToLower(muxer.Options.SegmentType)
	if segmentType == "" || segmentType == "mpegts" {
		args = append(
			args,
			"-hls_segment_type", "mpegts",
			"-hls_segment_filename", path.Join(muxer.Directory, "%s-%%d.ts"))
	} else if segmentType == "fmp4" {
		args = append(
			args,
			"-hls_segment_type", "fmp4",
			"-hls_segment_filename", path.Join(muxer.Directory, "%s-%%d.m4s"))
	} else {
		return errors.New("ffmpeg dash: invalid segment type")
	}

	if muxer.Options.Fps != 0 {
		args = append(args, "-r", strconv.Itoa(muxer.Options.Fps))
	}

	if muxer.Options.SegmentTime != 0 {
		args = append(args, "-hls_time", strconv.Itoa(muxer.Options.SegmentTime))
		hlsFlags = append(hlsFlags, "split_by_time")
	}

	if muxer.Options.PlaylistSize != 0 {
		args = append(args, "-hls_list_size", strconv.Itoa(muxer.Options.PlaylistSize))
	}

	if muxer.Options.StorageSize != 0 {
		args = append(args, "-hls_delete_threshold", strconv.Itoa(muxer.Options.StorageSize))
		hlsFlags = append(hlsFlags, "delete_segments")
	}

	args = append(args, "-hls_flags", strings.Join(hlsFlags, "+"), path.Join(muxer.Directory, "livestream.m3u8"))

	muxer.cmd = execCommand("ffmpeg", args...)
	muxer.cmd.Stdin = video

	log.Println("ffmpeg", muxer.cmd.String())

	return muxer.cmd.Start()
}

// Wait blocks until the video stream is finished processing by Mux.
func (muxer *Muxer) Wait() error {
	if muxer.cmd == nil {
		return errors.New("ffmpeg hls: not started")
	}

	return muxer.cmd.Wait()
}
