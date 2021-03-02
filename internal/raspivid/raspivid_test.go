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
	// Facilitate the "mocking" of os/exec by running a faked CLI program
	switch os.Getenv("GO_TEST_MODE") {
	case "":
		os.Exit(m.Run())
	case "raspivid":
		os.Stdout.WriteString(fakeVideoStreamContent)
		os.Exit(0)
	}
}

func TestNewStream(t *testing.T) {
	testCases := []struct {
		options      Options
		expectedArgs []string
	}{
		{
			Options{},
			[]string{
				"raspivid",
				"-o", "-",
				"-t", "0",
			},
		},
		{
			Options{Width: 1920},
			[]string{
				"raspivid",
				"-o", "-",
				"-t", "0",
				"--width", "1920",
			},
		},
		{
			Options{Height: 1080},
			[]string{
				"raspivid",
				"-o", "-",
				"-t", "0",
				"--height", "1080",
			},
		},
		{
			Options{Fps: 60},
			[]string{
				"raspivid",
				"-o", "-",
				"-t", "0",
				"--framerate", "60",
			},
		},
		{
			Options{HorizontalFlip: true},
			[]string{
				"raspivid",
				"-o", "-",
				"-t", "0",
				"--hflip",
			},
		},
		{
			Options{VerticalFlip: true},
			[]string{
				"raspivid",
				"-o", "-",
				"-t", "0",
				"--vflip",
			},
		},
		{
			Options{Width: 1280, Height: 720, Fps: 30, HorizontalFlip: true, VerticalFlip: true},
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
		t.Run(fmt.Sprintf("%v", tc.options), func(t *testing.T) {
			execCommand = mockExecCommand
			defer func() { execCommand = exec.Command }()

			raspiStream, err := NewStream(tc.options)

			if err != nil {
				t.Error("NewStream produced an err", err)
			}

			raspividArgs := raspiStream.cmd.Args[1:]

			if !equal(raspividArgs, tc.expectedArgs) {
				t.Error("Command args do not match, got", raspividArgs)
			}

			if raspiStream.Video == nil {
				t.Error("NewStream produced a Stream without video output")
			}

			if raspiStream.cmd.Process != nil {
				t.Error("NewStream started the stream prematurely")
			}
		})
	}
}

func TestStart(t *testing.T) {
	execCommand = mockExecCommand
	defer func() { execCommand = exec.Command }()

	raspiStream, _ := NewStream(Options{})
	err := raspiStream.Start()

	if err != nil {
		t.Error("Start produced an err", err)
	}

	if raspiStream.cmd.Process == nil {
		t.Error("Start failed to start raspivid")
	}

	// Unwrap video stream
	buf := new(strings.Builder)
	io.Copy(buf, raspiStream.Video)
	videoText := buf.String()

	if videoText != fakeVideoStreamContent {
		t.Error("Video output is invalid", videoText)
	}
}

func TestStartReturnsError(t *testing.T) {
	execCommand = mockFailedExecCommand
	defer func() { execCommand = exec.Command }()

	raspiStream, _ := NewStream(Options{})
	err := raspiStream.Start()

	if err == nil {
		t.Error("Start failed to return an error")
	}
}

func TestStartBadStreamReturnsError(t *testing.T) {
	execCommand = mockFailedExecCommand
	defer func() { execCommand = exec.Command }()

	raspiStream := Stream{Video: io.NopCloser(strings.NewReader(fakeVideoStreamContent))}
	err := raspiStream.Start()

	if err == nil || err.Error() != "raspivid: not created" {
		t.Error("Start failed to return correct error", err)
	}
}

func TestWait(t *testing.T) {
	execCommand = mockExecCommand
	defer func() { execCommand = exec.Command }()

	raspiStream, _ := NewStream(Options{})
	raspiStream.Start()
	err := raspiStream.Wait()

	if err != nil {
		t.Error("Wait returned an error", err)
	}
}

func TestWaitBadStreamReturnsError(t *testing.T) {
	execCommand = mockExecCommand
	defer func() { execCommand = exec.Command }()

	raspiStream := Stream{Video: io.NopCloser(strings.NewReader(fakeVideoStreamContent))}
	err := raspiStream.Wait()

	if err == nil || err.Error() != "raspivid: not created" {
		t.Error("Wait failed to return correct error", err)
	}
}

func TestWaitWithoutStartReturnsError(t *testing.T) {
	execCommand = mockExecCommand
	defer func() { execCommand = exec.Command }()

	raspiStream, _ := NewStream(Options{})
	err := raspiStream.Wait()

	if err == nil || err.Error() != "raspivid: not started" {
		t.Error("Wait failed to return correct error", err)
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

// mockExecCommaned sets up a mocked exec.Command using TestMain
func mockExecCommand(command string, args ...string) *exec.Cmd {
	cs := append([]string{command}, args...)
	cmd := exec.Command(os.Args[0], cs...)
	cmd.Env = append(os.Environ(), "GO_TEST_MODE=raspivid")
	return cmd
}

// mockFailedExecCommaned sets up a exec.Command that will fail
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
