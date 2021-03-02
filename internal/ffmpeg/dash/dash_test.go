package dash

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
				"-codec", "copy",
				"-f", "dash",
				"-re",
				"-an",
				"-init_seg_name", "init.$ext$",
				"-media_seg_name", "$Time$-$Number$.$ext$",
				"-dash_segment_type", "auto",
				"livestream.mpd",
			},
		},
		{
			Muxer{Directory: "camera"},
			[]string{
				"ffmpeg",
				"-codec", "copy",
				"-f", "dash",
				"-re",
				"-an",
				"-init_seg_name", "init.$ext$",
				"-media_seg_name", "$Time$-$Number$.$ext$",
				"-dash_segment_type", "auto",
				path.Join("camera", "livestream.mpd"),
			},
		},
		{
			Muxer{Options: Options{Fps: 60}},
			[]string{
				"ffmpeg",
				"-codec", "copy",
				"-f", "dash",
				"-re",
				"-an",
				"-init_seg_name", "init.$ext$",
				"-media_seg_name", "$Time$-$Number$.$ext$",
				"-dash_segment_type", "auto",
				"-r", "60",
				"livestream.mpd",
			},
		},
		{
			Muxer{Options: Options{SegmentType: "aUtO"}},
			[]string{
				"ffmpeg",
				"-codec", "copy",
				"-f", "dash",
				"-re",
				"-an",
				"-init_seg_name", "init.$ext$",
				"-media_seg_name", "$Time$-$Number$.$ext$",
				"-dash_segment_type", "auto",
				"livestream.mpd",
			},
		},
		{
			Muxer{Options: Options{SegmentType: "mP4"}},
			[]string{
				"ffmpeg",
				"-codec", "copy",
				"-f", "dash",
				"-re",
				"-an",
				"-init_seg_name", "init.$ext$",
				"-media_seg_name", "$Time$-$Number$.$ext$",
				"-dash_segment_type", "mp4",
				"livestream.mpd",
			},
		},
		{
			Muxer{Options: Options{SegmentType: "WeBm"}},
			[]string{
				"ffmpeg",
				"-codec", "copy",
				"-f", "dash",
				"-re",
				"-an",
				"-init_seg_name", "init.$ext$",
				"-media_seg_name", "$Time$-$Number$.$ext$",
				"-dash_segment_type", "webm",
				"livestream.mpd",
			},
		},
		{
			Muxer{Options: Options{SegmentTime: 2}},
			[]string{
				"ffmpeg",
				"-codec", "copy",
				"-f", "dash",
				"-re",
				"-an",
				"-init_seg_name", "init.$ext$",
				"-media_seg_name", "$Time$-$Number$.$ext$",
				"-dash_segment_type", "auto",
				"-seg_duration", "2",
				"livestream.mpd",
			},
		},
		{
			Muxer{Options: Options{PlaylistSize: 50}},
			[]string{
				"ffmpeg",
				"-codec", "copy",
				"-f", "dash",
				"-re",
				"-an",
				"-init_seg_name", "init.$ext$",
				"-media_seg_name", "$Time$-$Number$.$ext$",
				"-dash_segment_type", "auto",
				"-window_size", "50",
				"livestream.mpd",
			},
		},
		{
			Muxer{Options: Options{StorageSize: 100}},
			[]string{
				"ffmpeg",
				"-codec", "copy",
				"-f", "dash",
				"-re",
				"-an",
				"-init_seg_name", "init.$ext$",
				"-media_seg_name", "$Time$-$Number$.$ext$",
				"-dash_segment_type", "auto",
				"-extra_window_size", "100",
				"livestream.mpd",
			},
		},
		{
			Muxer{Directory: "dash", Options: Options{Fps: 30, SegmentTime: 5, PlaylistSize: 25, StorageSize: 50}},
			[]string{
				"ffmpeg",
				"-codec", "copy",
				"-f", "dash",
				"-re",
				"-an",
				"-init_seg_name", "init.$ext$",
				"-media_seg_name", "$Time$-$Number$.$ext$",
				"-dash_segment_type", "auto",
				"-r", "30",
				"-seg_duration", "5",
				"-window_size", "25",
				"-extra_window_size", "50",
				path.Join("dash", "livestream.mpd"),
			},
		},
	}

	for _, tc := range testCases {
		t.Run(fmt.Sprintf("%v", tc.muxer), func(t *testing.T) {
			execCommand = mockExecCommand
			defer func() { execCommand = exec.Command }()

			videoStream := ioutil.NopCloser(strings.NewReader("totallyfakevideostream"))

			mpegdashMuxer := tc.muxer
			err := mpegdashMuxer.Mux(videoStream)

			if err != nil {
				t.Error("Start produced an err", err)
			}

			ffmpegArgs := mpegdashMuxer.cmd.Args[1:]

			if !equal(ffmpegArgs, tc.expectedArgs) {
				t.Error("Command args do not match, got", ffmpegArgs)
			}
		})
	}
}

func TestStartInvalidSegmentTypeReturnsError(t *testing.T) {
	execCommand = mockExecCommand
	defer func() { execCommand = exec.Command }()

	videoStream := ioutil.NopCloser(strings.NewReader("totallyfakevideostream"))

	mpegdashMuxer := Muxer{Options: Options{SegmentType: "badtype"}}
	err := mpegdashMuxer.Mux(videoStream)

	if err.Error() != "ffmpeg dash: invalid segment type" {
		t.Error("Start failed to return an error for inavlid segment type")
	}
}

func TestStartReturnsFfmpegError(t *testing.T) {
	execCommand = mockFailedExecCommand
	defer func() { execCommand = exec.Command }()

	videoStream := ioutil.NopCloser(strings.NewReader("totallyfakevideostream"))

	mpegdashMuxer := Muxer{}
	err := mpegdashMuxer.Mux(videoStream)

	if err == nil {
		t.Error("Start failed to return an error")
	}
}

func TestWait(t *testing.T) {
	execCommand = mockExecCommand
	defer func() { execCommand = exec.Command }()

	videoStream := ioutil.NopCloser(strings.NewReader("totallyfakevideostream"))

	mpegdashMuxer := Muxer{}
	mpegdashMuxer.Mux(videoStream)
	err := mpegdashMuxer.Wait()

	if err != nil {
		t.Error("Wait returned an error", err)
	}
}

func TestWaitWithoutStartReturnsError(t *testing.T) {
	execCommand = mockExecCommand
	defer func() { execCommand = exec.Command }()

	mpegdashMuxer := Muxer{}
	err := mpegdashMuxer.Wait()

	if err == nil || err.Error() != "ffmpeg dash: not started" {
		t.Error("Wait failed to return correct error when run without Start", err)
	}
}

func TestWaitAgainReturnsError(t *testing.T) {
	execCommand = mockExecCommand
	defer func() { execCommand = exec.Command }()

	videoStream := ioutil.NopCloser(strings.NewReader("totallyfakevideostream"))

	mpegdashMuxer := Muxer{}
	mpegdashMuxer.Mux(videoStream)
	mpegdashMuxer.Wait()
	err := mpegdashMuxer.Wait()

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
