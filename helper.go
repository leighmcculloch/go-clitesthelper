package testcli

import (
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
	"testing"
)

func New(tb testing.TB) Helper {
	return Helper{
		TB: tb,
	}
}

type Helper struct {
	TB testing.TB
}

type MainFunc func(args []string, stdin io.Reader, stdout, stderr io.Writer) int

func (h Helper) Main(args []string, stdin io.Reader, main MainFunc) (exitCode int, stdout, stderr string) {
	h.TB.Helper()
	combinedBuilder := strings.Builder{}
	stdoutBuilder := strings.Builder{}
	stderrBuilder := strings.Builder{}
	exitCode = main(
		args,
		stdin,
		io.MultiWriter(&stdoutBuilder, &combinedBuilder),
		io.MultiWriter(&stderrBuilder, &combinedBuilder),
	)
	h.TB.Log(strings.Join(args, "") + "\n" + strings.TrimSuffix(string(combinedBuilder.String()), "\n"))
	stdout = stdoutBuilder.String()
	stderr = stderrBuilder.String()
	return

}

func (h Helper) Exec(command string) (exitCode int, stdout, stderr string) {
	h.TB.Helper()
	combinedBuilder := strings.Builder{}
	stdoutBuilder := strings.Builder{}
	stderrBuilder := strings.Builder{}
	cmd := exec.Command("bash", "-c", command)
	cmd.Stdout = io.MultiWriter(&stdoutBuilder, &combinedBuilder)
	cmd.Stderr = io.MultiWriter(&stderrBuilder, &combinedBuilder)
	err := cmd.Run()
	if err != nil {
		h.TB.Fatalf("executing %q: %v", command, err)
	}
	h.TB.Log(command + "\n" + strings.TrimSuffix(string(combinedBuilder.String()), "\n"))
	exitCode = cmd.ProcessState.ExitCode()
	stdout = stdoutBuilder.String()
	stderr = stderrBuilder.String()
	return
}

func (h Helper) Chdir(dir string) {
	h.TB.Helper()

	err := os.Chdir(dir)
	if err != nil {
		h.TB.Fatalf("chdir to %q: %v", dir, err)
	}
	h.TB.Log("cd " + dir)
}

func (h Helper) MkdirTemp() string {
	h.TB.Helper()

	dir, err := os.MkdirTemp("", "")
	if err != nil {
		h.TB.Fatalf("mkdirtemp: %v", err)
	}
	h.TB.Log("mktemp -d\n" + dir)
	return dir
}

func (h Helper) WriteFile(filename string, data []byte) {
	h.TB.Helper()

	err := ioutil.WriteFile(filename, []byte{}, 0644)
	if err != nil {
		h.TB.Fatalf("writing file to %q: %v", filename, err)
	}
}
