package hls

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"strings"
	"testing"
)

const fakeVideoStreamContent = "fakevideostream"

func TestMain(m *testing.M) {
	switch os.Getenv("GO_TEST_MODE") {
	case "":
		os.Exit(m.Run())
	case "ffmpeg":
		os.Stdout.WriteString(fakeVideoStreamContent)
		os.Exit(0)
	}
}

func TestStart(t *testing.T) {
	testCases := []struct {
		muxer        Muxer
		expectedArgs []string
	}{
		{
			Muxer{},
			[]string{
				"ffmpeg",
				"-i", "pipe:0",
				"-codec", "copy",
				"-f", "hls",
				"-an",
				"-hls_segment_type", "mpegts",
				"-hls_segment_filename", "raspilive-%03d.ts",
				"livestream.m3u8",
			},
		},
		{
			Muxer{Directory: "camera"},
			[]string{
				"ffmpeg",
				"-i", "pipe:0",
				"-codec", "copy",
				"-f", "hls",
				"-an",
				"-hls_segment_type", "mpegts",
				"-hls_segment_filename", "camera/raspilive-%03d.ts",
				path.Join("camera", "livestream.m3u8"),
			},
		},
		{
			Muxer{Options: Options{Fps: 60}},
			[]string{
				"ffmpeg",
				"-i", "pipe:0",
				"-codec", "copy",
				"-f", "hls",
				"-an",
				"-hls_segment_type", "mpegts",
				"-hls_segment_filename", "raspilive-%03d.ts",
				"-r", "60",
				"livestream.m3u8",
			},
		},
		{
			Muxer{Options: Options{SegmentType: "MpEgTs"}},
			[]string{
				"ffmpeg",
				"-i", "pipe:0",
				"-codec", "copy",
				"-f", "hls",
				"-an",
				"-hls_segment_type", "mpegts",
				"-hls_segment_filename", "raspilive-%03d.ts",
				"livestream.m3u8",
			},
		},
		{
			Muxer{Options: Options{SegmentType: "FmP4"}},
			[]string{
				"ffmpeg",
				"-i", "pipe:0",
				"-codec", "copy",
				"-f", "hls",
				"-an",
				"-hls_segment_type", "fmp4",
				"-hls_segment_filename", "raspilive-%d.m4s",
				"livestream.m3u8",
			},
		},
		{
			Muxer{Options: Options{SegmentTime: 2}},
			[]string{
				"ffmpeg",
				"-i", "pipe:0",
				"-codec", "copy",
				"-f", "hls",
				"-an",
				"-hls_segment_type", "mpegts",
				"-hls_segment_filename", "raspilive-%03d.ts",
				"-hls_time", "2",
				"-hls_flags", "split_by_time",
				"livestream.m3u8",
			},
		},
		{
			Muxer{Options: Options{PlaylistSize: 50}},
			[]string{
				"ffmpeg",
				"-i", "pipe:0",
				"-codec", "copy",
				"-f", "hls",
				"-an",
				"-hls_segment_type", "mpegts",
				"-hls_segment_filename", "raspilive-%03d.ts",
				"-hls_list_size", "50",
				"livestream.m3u8",
			},
		},
		{
			Muxer{Options: Options{StorageSize: 100}},
			[]string{
				"ffmpeg",
				"-i", "pipe:0",
				"-codec", "copy",
				"-f", "hls",
				"-an",
				"-hls_segment_type", "mpegts",
				"-hls_segment_filename", "raspilive-%03d.ts",
				"-hls_delete_threshold", "100",
				"-hls_flags", "delete_segments",
				"livestream.m3u8",
			},
		},
		{
			Muxer{
				Directory: "hls",
				Options:   Options{Fps: 30, SegmentType: "fmp4", SegmentTime: 5, PlaylistSize: 25, StorageSize: 50},
			},
			[]string{
				"ffmpeg",
				"-i", "pipe:0",
				"-codec", "copy",
				"-f", "hls",
				"-an",
				"-hls_segment_type", "fmp4",
				"-hls_segment_filename", "hls/raspilive-%d.m4s",
				"-r", "30",
				"-hls_time", "5",
				"-hls_list_size", "25",
				"-hls_delete_threshold", "50",
				"-hls_flags", "split_by_time+delete_segments",
				path.Join("hls", "livestream.m3u8"),
			},
		},
	}

	for _, tc := range testCases {
		t.Run(fmt.Sprintf("%v", tc.muxer), func(t *testing.T) {
			execCommand = mockExecCommand
			defer func() { execCommand = exec.Command }()

			videoStream := ioutil.NopCloser(strings.NewReader("totallyfakevideostream"))

			hlsMuxer := tc.muxer
			err := hlsMuxer.Mux(videoStream)

			if err != nil {
				t.Error("Start produced an err", err)
			}

			ffmpegArgs := hlsMuxer.cmd.Args[1:]

			if !equal(ffmpegArgs, tc.expectedArgs) {
				t.Error("Command args do not match, got", ffmpegArgs, "but wanted", tc.expectedArgs)
			}
		})
	}
}

func TestStartInvalidSegmentTypeReturnsError(t *testing.T) {
	execCommand = mockExecCommand
	defer func() { execCommand = exec.Command }()

	videoStream := ioutil.NopCloser(strings.NewReader("totallyfakevideostream"))

	hlsMuxer := Muxer{Options: Options{SegmentType: "badtype"}}
	err := hlsMuxer.Mux(videoStream)

	if err.Error() != "ffmpeg dash: invalid segment type" {
		t.Error("Start failed to return an error for inavlid segment type")
	}
}

func TestStartReturnsFfmpegError(t *testing.T) {
	execCommand = mockFailedExecCommand
	defer func() { execCommand = exec.Command }()

	videoStream := ioutil.NopCloser(strings.NewReader("totallyfakevideostream"))

	hlsMuxer := Muxer{}
	err := hlsMuxer.Mux(videoStream)

	if err == nil {
		t.Error("Start failed to return an error")
	}
}

func TestWait(t *testing.T) {
	execCommand = mockExecCommand
	defer func() { execCommand = exec.Command }()

	videoStream := ioutil.NopCloser(strings.NewReader("totallyfakevideostream"))

	hlsMuxer := Muxer{}
	hlsMuxer.Mux(videoStream)
	err := hlsMuxer.Wait()

	if err != nil {
		t.Error("Wait returned an error", err)
	}
}

func TestWaitWithoutStartReturnsError(t *testing.T) {
	execCommand = mockExecCommand
	defer func() { execCommand = exec.Command }()

	hlsMuxer := Muxer{}
	err := hlsMuxer.Wait()

	if err == nil || err.Error() != "ffmpeg hls: not started" {
		t.Error("Wait failed to return correct error when run without Start", err)
	}
}

func TestWaitAgainReturnsError(t *testing.T) {
	execCommand = mockExecCommand
	defer func() { execCommand = exec.Command }()

	videoStream := ioutil.NopCloser(strings.NewReader("totallyfakevideostream"))

	hlsMuxer := Muxer{}
	hlsMuxer.Mux(videoStream)
	hlsMuxer.Wait()
	err := hlsMuxer.Wait()

	if err == nil {
		t.Error("Wait failed to return an error")
	}
}

func mockExecCommand(command string, args ...string) *exec.Cmd {
	cs := append([]string{command}, args...)
	cmd := exec.Command(os.Args[0], cs...)
	cmd.Env = append(os.Environ(), "GO_TEST_MODE=ffmpeg")
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
