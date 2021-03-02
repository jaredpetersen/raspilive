package raspivid

import (
	"errors"
	"io"
	"os/exec"
	"strconv"
)

var execCommand = exec.Command

// Options represents a Raspberry Pi camera video streamer.
type Options struct {
	Width          int
	Height         int
	Fps            int
	HorizontalFlip bool
	VerticalFlip   bool
}

// Stream represents a Raspberry Pi camera video streamer.
type Stream struct {
	Video io.ReadCloser
	cmd   *exec.Cmd
}

// NewStream creates a new video stream out of the Raspberry Pi Camera Module.
func NewStream(options Options) (*Stream, error) {
	args := []string{"-o", "-", "-t", "0"}

	if options.Width != 0 {
		args = append(args, "--width", strconv.Itoa(options.Width))
	}

	if options.Height != 0 {
		args = append(args, "--height", strconv.Itoa(options.Height))
	}

	if options.Fps != 0 {
		args = append(args, "--framerate", strconv.Itoa(options.Fps))
	}

	if options.HorizontalFlip {
		args = append(args, "--hflip")
	}

	if options.VerticalFlip {
		args = append(args, "--vflip")
	}

	cmd := execCommand("raspivid", args...)
	video, err := cmd.StdoutPipe()

	if err != nil {
		return nil, err
	}

	return &Stream{Video: video, cmd: cmd}, nil
}

// Start begins the video stream.
func (strm *Stream) Start() error {
	if strm.cmd == nil {
		return errors.New("raspivid: not created")
	}

	return strm.cmd.Start()
}

// Wait waits for the video stream to complete.
//
// The stream operation must have been started by Start.
func (strm *Stream) Wait() error {
	if strm.cmd == nil {
		return errors.New("raspivid: not created")
	}
	if strm.cmd.Process == nil {
		return errors.New("raspivid: not started")
	}

	return strm.cmd.Wait()
}
