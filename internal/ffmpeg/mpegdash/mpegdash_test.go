package mpegdash

import (
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"strings"
	"testing"
	"time"

	"github.com/jaredpetersen/raspilive/internal/utils/pointer"
)

const ffmpegSleep = 3 * time.Second

func TestMain(m *testing.M) {
	switch os.Getenv("GO_TEST_MODE") {
	case "":
		os.Exit(m.Run())
	case "ffmpeg":
		time.Sleep(ffmpegSleep)
		os.Exit(0)
	}
}

func TestMpegDashDefault(t *testing.T) {
	execCommand = mockExecCommand
	defer func() { execCommand = exec.Command }()

	videoStream := ioutil.NopCloser(strings.NewReader("totallyfakevideostream"))
	directory := "./camera"
	options := Options{}

	dashMuxer := MpegDash(videoStream, directory, options)

	dashArgs := dashMuxer.cmd.Args[1:]
	expectedArgs := []string{
		"ffmpeg",
		"-codec", "copy",
		"-f", "dash",
		"-re",
		"-an",
		"-init_seg_name", "init.m4s",
		"-media_seg_name", "$Time$-$Number$.m4s",
		path.Join(directory, "livestream.mpd"),
	}

	if !equal(dashArgs, expectedArgs) {
		t.Error("Ffmpeg args do not match, got", dashArgs)
	}

	if dashMuxer.cmd.Stdin != videoStream {
		t.Error("Ffmpeg input stream does not match")
	}
}

func TestMpegDashFps(t *testing.T) {
	execCommand = mockExecCommand
	defer func() { execCommand = exec.Command }()

	videoStream := ioutil.NopCloser(strings.NewReader("totallyfakevideostream"))
	directory := "./camera"
	options := Options{Fps: pointer.ToInt(30)}

	dashMuxer := MpegDash(videoStream, directory, options)

	dashArgs := dashMuxer.cmd.Args[1:]
	expectedArgs := []string{
		"ffmpeg",
		"-codec", "copy",
		"-f", "dash",
		"-re",
		"-an",
		"-init_seg_name", "init.m4s",
		"-media_seg_name", "$Time$-$Number$.m4s",
		"-r", "30",
		path.Join(directory, "livestream.mpd"),
	}

	if !equal(dashArgs, expectedArgs) {
		t.Error("Ffmpeg args do not match, got", dashArgs)
	}

	if dashMuxer.cmd.Stdin != videoStream {
		t.Error("Ffmpeg input stream does not match")
	}
}

func TestMpegDashTime(t *testing.T) {
	execCommand = mockExecCommand
	defer func() { execCommand = exec.Command }()

	videoStream := ioutil.NopCloser(strings.NewReader("totallyfakevideostream"))
	directory := "./camera"
	options := Options{Time: pointer.ToInt(4)}

	dashMuxer := MpegDash(videoStream, directory, options)

	dashArgs := dashMuxer.cmd.Args[1:]
	expectedArgs := []string{
		"ffmpeg",
		"-codec", "copy",
		"-f", "dash",
		"-re",
		"-an",
		"-init_seg_name", "init.m4s",
		"-media_seg_name", "$Time$-$Number$.m4s",
		"-seg_duration", "4",
		path.Join(directory, "livestream.mpd"),
	}

	if !equal(dashArgs, expectedArgs) {
		t.Error("Ffmpeg args do not match, got", dashArgs)
	}

	if dashMuxer.cmd.Stdin != videoStream {
		t.Error("Ffmpeg input stream does not match")
	}
}

func TestMpegDashListSize(t *testing.T) {
	execCommand = mockExecCommand
	defer func() { execCommand = exec.Command }()

	videoStream := ioutil.NopCloser(strings.NewReader("totallyfakevideostream"))
	directory := "./camera"
	options := Options{ListSize: pointer.ToInt(60)}

	dashMuxer := MpegDash(videoStream, directory, options)

	dashArgs := dashMuxer.cmd.Args[1:]
	expectedArgs := []string{
		"ffmpeg",
		"-codec", "copy",
		"-f", "dash",
		"-re",
		"-an",
		"-init_seg_name", "init.m4s",
		"-media_seg_name", "$Time$-$Number$.m4s",
		"-window_size", "60",
		path.Join(directory, "livestream.mpd"),
	}

	if !equal(dashArgs, expectedArgs) {
		t.Error("Ffmpeg args do not match, got", dashArgs)
	}

	if dashMuxer.cmd.Stdin != videoStream {
		t.Error("Ffmpeg input stream does not match")
	}
}

func TestMpegDashStorageSize(t *testing.T) {
	execCommand = mockExecCommand
	defer func() { execCommand = exec.Command }()

	videoStream := ioutil.NopCloser(strings.NewReader("totallyfakevideostream"))
	directory := "./camera"
	options := Options{StorageSize: pointer.ToInt(240)}

	dashMuxer := MpegDash(videoStream, directory, options)

	dashArgs := dashMuxer.cmd.Args[1:]
	expectedArgs := []string{
		"ffmpeg",
		"-codec", "copy",
		"-f", "dash",
		"-re",
		"-an",
		"-init_seg_name", "init.m4s",
		"-media_seg_name", "$Time$-$Number$.m4s",
		"-extra_window_size", "240",
		path.Join(directory, "livestream.mpd"),
	}

	if !equal(dashArgs, expectedArgs) {
		t.Error("Ffmpeg args do not match, got", dashArgs)
	}

	if dashMuxer.cmd.Stdin != videoStream {
		t.Error("Ffmpeg input stream does not match")
	}
}

func TestMpegDashAll(t *testing.T) {
	execCommand = mockExecCommand
	defer func() { execCommand = exec.Command }()

	videoStream := ioutil.NopCloser(strings.NewReader("totallyfakevideostream"))
	directory := "./dash"
	options := Options{
		Fps:         pointer.ToInt(60),
		Time:        pointer.ToInt(1),
		ListSize:    pointer.ToInt(40),
		StorageSize: pointer.ToInt(100),
	}

	dashMuxer := MpegDash(videoStream, directory, options)

	dashArgs := dashMuxer.cmd.Args[1:]
	expectedArgs := []string{
		"ffmpeg",
		"-codec", "copy",
		"-f", "dash",
		"-re",
		"-an",
		"-init_seg_name", "init.m4s",
		"-media_seg_name", "$Time$-$Number$.m4s",
		"-r", "60",
		"-seg_duration", "1",
		"-window_size", "40",
		"-extra_window_size", "100",
		path.Join(directory, "livestream.mpd"),
	}

	if !equal(dashArgs, expectedArgs) {
		t.Error("Ffmpeg args do not match, got", dashArgs)
	}

	if dashMuxer.cmd.Stdin != videoStream {
		t.Error("Ffmpeg input stream does not match")
	}
}

func TestStartAndWait(t *testing.T) {
	execCommand = mockExecCommand
	defer func() { execCommand = exec.Command }()

	videoStream := ioutil.NopCloser(strings.NewReader("totallyfakevideostream"))
	directory := "./camera"
	options := Options{}

	dashMuxer := MpegDash(videoStream, directory, options)

	startError := dashMuxer.Start()

	if startError != nil {
		t.Error("Start encountered error", startError)
	}
	if dashMuxer.cmd.Process == nil {
		t.Fatal("Start has not started a new process")
	}

	waitError := dashMuxer.Wait()

	if waitError != nil {
		t.Error("Start encountered error", startError)
	}
	if !dashMuxer.cmd.ProcessState.Exited() {
		t.Error("Start execution is not complete")
	}
}

func mockExecCommand(command string, args ...string) *exec.Cmd {
	cs := append([]string{command}, args...)
	cmd := exec.Command(os.Args[0], cs...)
	cmd.Env = append(os.Environ(), "GO_TEST_MODE=ffmpeg")
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
