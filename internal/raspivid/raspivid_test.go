package raspivid

import (
	"os"
	"os/exec"
	"testing"
	"time"
)

const raspividSleep = 3 * time.Second

func TestMain(m *testing.M) {
	switch os.Getenv("GO_TEST_MODE") {
	case "":
		os.Exit(m.Run())
	case "raspivid":
		time.Sleep(raspividSleep)
		os.Exit(0)
	}
}

func TestRaspividDefault(t *testing.T) {
	execCommand = mockExecCommand
	defer func() { execCommand = exec.Command }()

	options := Options{}

	raspivid := Stream(options)

	raspividArgs := raspivid.cmd.Args[1:]
	expectedArgs := []string{
		"raspivid",
		"-o", "-",
		"-t", "0",
	}

	if !equal(raspividArgs, expectedArgs) {
		t.Error("Raspivid args do not match, got", raspividArgs)
	}

	if raspivid.cmd.Stdout == nil {
		t.Error("Raspivid output stream does not match")
	}
}

func TestStartAndWait(t *testing.T) {
	execCommand = mockExecCommand
	defer func() { execCommand = exec.Command }()

	options := Options{}

	raspividStream := Stream(options)

	err := raspividStream.Start()

	if err != nil {
		t.Error("Start encountered error", err)
	}
	if raspividStream.cmd.Process == nil {
		t.Fatal("Start has not started a new process")
	}

	err = raspividStream.Wait()

	if err != nil {
		t.Error("Start encountered error", err)
	}
	if !raspividStream.cmd.ProcessState.Exited() {
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
