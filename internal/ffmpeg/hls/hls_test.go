package hls

import (
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"strings"
	"testing"
	"time"
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

func TestHlsDefault(t *testing.T) {
	execCommand = mockExecCommand
	defer func() { execCommand = exec.Command }()

	videoStream := ioutil.NopCloser(strings.NewReader("totallyfakevideostream"))
	directory := "./camera"
	options := Options{}

	hlsMuxer := Hls(videoStream, directory, options)

	hlsArgs := hlsMuxer.cmd.Args[1:]
	expectedArgs := []string{
		"ffmpeg",
		"-codec", "copy",
		"-f", "hls",
		"-re",
		"-an",
		"-strftime", "1",
		"-hls_segment_filename", "%s-%%d.m4s",
		"-hls_segment_type", "fmp4",
		"-hls_flags", "second_level_segment_index",
		path.Join(directory, "livestream.m3u8"),
	}

	if !equal(hlsArgs, expectedArgs) {
		t.Error("Ffmpeg args do not match, got", hlsArgs)
	}

	if hlsMuxer.cmd.Stdin != videoStream {
		t.Error("Ffmpeg input stream does not match")
	}
}

func TestHlsFps(t *testing.T) {
	execCommand = mockExecCommand
	defer func() { execCommand = exec.Command }()

	videoStream := ioutil.NopCloser(strings.NewReader("totallyfakevideostream"))
	directory := "./camera"
	options := Options{Fps: 30}

	hlsMuxer := Hls(videoStream, directory, options)

	hlsArgs := hlsMuxer.cmd.Args[1:]
	expectedArgs := []string{
		"ffmpeg",
		"-codec", "copy",
		"-f", "hls",
		"-re",
		"-an",
		"-strftime", "1",
		"-hls_segment_filename", "%s-%%d.m4s",
		"-hls_segment_type", "fmp4",
		"-r", "30",
		"-hls_flags", "second_level_segment_index",
		path.Join(directory, "livestream.m3u8"),
	}

	if !equal(hlsArgs, expectedArgs) {
		t.Error("Ffmpeg args do not match, got", hlsArgs)
	}

	if hlsMuxer.cmd.Stdin != videoStream {
		t.Error("Ffmpeg input stream does not match")
	}
}

func TestHlsTime(t *testing.T) {
	execCommand = mockExecCommand
	defer func() { execCommand = exec.Command }()

	videoStream := ioutil.NopCloser(strings.NewReader("totallyfakevideostream"))
	directory := "./camera"
	options := Options{SegmentTime: 4}

	hlsMuxer := Hls(videoStream, directory, options)

	hlsArgs := hlsMuxer.cmd.Args[1:]
	expectedArgs := []string{
		"ffmpeg",
		"-codec", "copy",
		"-f", "hls",
		"-re",
		"-an",
		"-strftime", "1",
		"-hls_segment_filename", "%s-%%d.m4s",
		"-hls_segment_type", "fmp4",
		"-hls_time", "4",
		"-hls_flags", "second_level_segment_index+split_by_time",
		path.Join(directory, "livestream.m3u8"),
	}

	if !equal(hlsArgs, expectedArgs) {
		t.Error("Ffmpeg args do not match, got", hlsArgs)
	}

	if hlsMuxer.cmd.Stdin != videoStream {
		t.Error("Ffmpeg input stream does not match")
	}
}

func TestHlsListSize(t *testing.T) {
	execCommand = mockExecCommand
	defer func() { execCommand = exec.Command }()

	videoStream := ioutil.NopCloser(strings.NewReader("totallyfakevideostream"))
	directory := "./camera"
	options := Options{PlaylistSize: 60}

	hlsMuxer := Hls(videoStream, directory, options)

	hlsArgs := hlsMuxer.cmd.Args[1:]
	expectedArgs := []string{
		"ffmpeg",
		"-codec", "copy",
		"-f", "hls",
		"-re",
		"-an",
		"-strftime", "1",
		"-hls_segment_filename", "%s-%%d.m4s",
		"-hls_segment_type", "fmp4",
		"-hls_list_size", "60",
		"-hls_flags",
		"second_level_segment_index",
		path.Join(directory, "livestream.m3u8"),
	}

	if !equal(hlsArgs, expectedArgs) {
		t.Error("Ffmpeg args do not match, got", hlsArgs)
	}

	if hlsMuxer.cmd.Stdin != videoStream {
		t.Error("Ffmpeg input stream does not match")
	}
}

func TestHlsStorageSize(t *testing.T) {
	execCommand = mockExecCommand
	defer func() { execCommand = exec.Command }()

	videoStream := ioutil.NopCloser(strings.NewReader("totallyfakevideostream"))
	directory := "./camera"
	options := Options{StorageSize: 240}

	hlsMuxer := Hls(videoStream, directory, options)

	hlsArgs := hlsMuxer.cmd.Args[1:]
	expectedArgs := []string{
		"ffmpeg",
		"-codec", "copy",
		"-f", "hls",
		"-re",
		"-an",
		"-strftime", "1",
		"-hls_segment_filename", "%s-%%d.m4s",
		"-hls_segment_type", "fmp4",
		"-hls_delete_threshold", "240",
		"-hls_flags", "second_level_segment_index+delete_segments",
		path.Join(directory, "livestream.m3u8"),
	}

	if !equal(hlsArgs, expectedArgs) {
		t.Error("Ffmpeg args do not match, got", hlsArgs)
	}

	if hlsMuxer.cmd.Stdin != videoStream {
		t.Error("Ffmpeg input stream does not match")
	}
}

func TestHlsAll(t *testing.T) {
	execCommand = mockExecCommand
	defer func() { execCommand = exec.Command }()

	videoStream := ioutil.NopCloser(strings.NewReader("totallyfakevideostream"))
	directory := "./hls"
	options := Options{
		Fps:          60,
		SegmentTime:  1,
		PlaylistSize: 40,
		StorageSize:  100,
	}

	hlsMuxer := Hls(videoStream, directory, options)

	hlsArgs := hlsMuxer.cmd.Args[1:]
	expectedArgs := []string{
		"ffmpeg",
		"-codec", "copy",
		"-f", "hls",
		"-re",
		"-an",
		"-strftime", "1",
		"-hls_segment_filename", "%s-%%d.m4s",
		"-hls_segment_type", "fmp4",
		"-r",
		"60",
		"-hls_time",
		"1",
		"-hls_list_size",
		"40",
		"-hls_delete_threshold",
		"100",
		"-hls_flags",
		"second_level_segment_index+split_by_time+delete_segments",
		path.Join(directory, "livestream.m3u8"),
	}

	if !equal(hlsArgs, expectedArgs) {
		t.Error("Ffmpeg args do not match, got", hlsArgs)
	}

	if hlsMuxer.cmd.Stdin != videoStream {
		t.Error("Ffmpeg input stream does not match")
	}
}

func TestStartAndWait(t *testing.T) {
	execCommand = mockExecCommand
	defer func() { execCommand = exec.Command }()

	videoStream := ioutil.NopCloser(strings.NewReader("totallyfakevideostream"))
	directory := "./camera"
	options := Options{}

	hlsMuxer := Hls(videoStream, directory, options)

	err := hlsMuxer.Start()

	if err != nil {
		t.Error("Start encountered error", err)
	}
	if hlsMuxer.cmd.Process == nil {
		t.Fatal("Start has not started a new process")
	}

	err = hlsMuxer.Wait()

	if err != nil {
		t.Error("Start encountered error", err)
	}
	if !hlsMuxer.cmd.ProcessState.Exited() {
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
