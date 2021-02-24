package raspivid

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
	"testing"
)

const fakeVideoStreamContent = "fakevideostream"

func TestMain(m *testing.M) {
	switch os.Getenv("GO_TEST_MODE") {
	case "":
		os.Exit(m.Run())
	case "raspivid":
		os.Stdout.WriteString(fakeVideoStreamContent)
		os.Exit(0)
	}
}

func TestStart(t *testing.T) {
	testCases := []struct {
		stream       Stream
		expectedArgs []string
	}{
		{
			Stream{},
			[]string{
				"raspivid",
				"-o", "-",
				"-t", "0",
			},
		},
		{
			Stream{Width: 1920},
			[]string{
				"raspivid",
				"-o", "-",
				"-t", "0",
				"--width", "1920",
			},
		},
		{
			Stream{Height: 1080},
			[]string{
				"raspivid",
				"-o", "-",
				"-t", "0",
				"--height", "1080",
			},
		},
		{
			Stream{Fps: 60},
			[]string{
				"raspivid",
				"-o", "-",
				"-t", "0",
				"--framerate", "60",
			},
		},
		{
			Stream{HorizontalFlip: true},
			[]string{
				"raspivid",
				"-o", "-",
				"-t", "0",
				"--hflip",
			},
		},
		{
			Stream{VerticalFlip: true},
			[]string{
				"raspivid",
				"-o", "-",
				"-t", "0",
				"--vflip",
			},
		},
		{
			Stream{Width: 1280, Height: 720, Fps: 30, HorizontalFlip: true, VerticalFlip: true},
			[]string{
				"raspivid",
				"-o", "-",
				"-t", "0",
				"--width", "1280",
				"--height", "720",
				"--framerate", "30",
				"--hflip", "--vflip",
			},
		},
	}

	for _, tc := range testCases {
		t.Run(fmt.Sprintf("%v", tc.stream), func(t *testing.T) {
			execCommand = mockExecCommand
			defer func() { execCommand = exec.Command }()

			raspividStream := tc.stream
			videoStream, err := raspividStream.Start()

			if err != nil {
				t.Error("Start produced an err", err)
			}

			if videoStream == nil {
				t.Error("Start produced a nil video stream")
			}

			buf := new(strings.Builder)
			io.Copy(buf, videoStream)
			videoStreamText := buf.String()

			if videoStreamText != fakeVideoStreamContent {
				t.Error("Start video stream is incorrect, got", videoStream)
			}

			raspividArgs := raspividStream.cmd.Args[1:]

			if !equal(raspividArgs, tc.expectedArgs) {
				t.Error("Command args do not match, got", raspividArgs)
			}
		})
	}
}

func TestStartReturnsError(t *testing.T) {
	execCommand = mockFailedExecCommand
	defer func() { execCommand = exec.Command }()

	raspividStream := Stream{}
	_, err := raspividStream.Start()

	if err == nil {
		t.Error("Started failed to return an error")
	}
}

func TestWait(t *testing.T) {
	execCommand = mockExecCommand
	defer func() { execCommand = exec.Command }()

	raspividStream := Stream{}
	raspividStream.Start()
	err := raspividStream.Wait()

	if err != nil {
		t.Error("Wait returned an error", err)
	}
}

func TestWaitWithoutStartReturnsError(t *testing.T) {
	execCommand = mockExecCommand
	defer func() { execCommand = exec.Command }()

	raspividStream := Stream{}
	err := raspividStream.Wait()

	if err == nil || err.Error() != "raspivid: not started" {
		t.Error("Wait failed to return correct error when run without Start", err)
	}
}

func TestWaitAgainReturnsError(t *testing.T) {
	execCommand = mockExecCommand
	defer func() { execCommand = exec.Command }()

	raspividStream := Stream{}
	raspividStream.Start()
	raspividStream.Wait()
	err := raspividStream.Wait()

	if err == nil {
		t.Error("Wait failed to return an error")
	}
}

func mockExecCommand(command string, args ...string) *exec.Cmd {
	cs := append([]string{command}, args...)
	cmd := exec.Command(os.Args[0], cs...)
	cmd.Env = append(os.Environ(), "GO_TEST_MODE=raspivid")
	return cmd
}

func mockFailedExecCommand(command string, args ...string) *exec.Cmd {
	cmd := exec.Command("totallyfakecommandthatdoesnotexist")
	return cmd
}

func equal(a, b []string) bool {
	// If one is nil, the other must also be nil.
	if (a == nil) != (b == nil) {
		return false
	}

	if len(a) != len(b) {
		return false
	}

	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}

	return true
}
