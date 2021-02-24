package raspivid

import (
	"errors"
	"io"
	"os/exec"
	"strconv"
)

var execCommand = exec.Command

// Stream represents a Raspberry Pi camera video streamer.
type Stream struct {
	Width          int
	Height         int
	Fps            int
	HorizontalFlip bool
	VerticalFlip   bool
	cmd            *exec.Cmd
}

// Start begins the video stream.
func (strm *Stream) Start() (io.ReadCloser, error) {
	args := []string{"-o", "-", "-t", "0"}

	if strm.Width != 0 {
		args = append(args, "--width", strconv.Itoa(strm.Width))
	}

	if strm.Height != 0 {
		args = append(args, "--height", strconv.Itoa(strm.Height))
	}

	if strm.Fps != 0 {
		args = append(args, "--framerate", strconv.Itoa(strm.Fps))
	}

	if strm.HorizontalFlip {
		args = append(args, "--hflip")
	}

	if strm.VerticalFlip {
		args = append(args, "--vflip")
	}

	strm.cmd = execCommand("raspivid", args...)
	video, err := strm.cmd.StdoutPipe()

	if err != nil {
		return nil, err
	}

	err = strm.cmd.Start()

	if err != nil {
		return nil, err
	}

	return video, err
}

// Wait waits for the video stream to complete.
//
// The stream operation must have been started by Start.
func (strm *Stream) Wait() error {
	if strm.cmd == nil {
		return errors.New("raspivid: not started")
	}

	return strm.cmd.Wait()
}
