package testcli

import (
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
	"testing"
)

type MainFunc func(args []string, stdin io.Reader, stdout, stderr io.Writer) int

func Main(t testing.TB, args []string, stdin io.Reader, main MainFunc) (exitCode int, stdout, stderr string) {
	t.Helper()
	combinedBuilder := strings.Builder{}
	stdoutBuilder := strings.Builder{}
	stderrBuilder := strings.Builder{}
	exitCode = main(
		args,
		stdin,
		io.MultiWriter(&stdoutBuilder, &combinedBuilder),
		io.MultiWriter(&stderrBuilder, &combinedBuilder),
	)
	t.Log(strings.Join(args, " ") + "\n" + strings.TrimSuffix(string(combinedBuilder.String()), "\n"))
	stdout = stdoutBuilder.String()
	stderr = stderrBuilder.String()
	return
}

func Exec(t testing.TB, command string) (exitCode int, stdout, stderr string) {
	t.Helper()
	combinedBuilder := strings.Builder{}
	stdoutBuilder := strings.Builder{}
	stderrBuilder := strings.Builder{}
	cmd := exec.Command("bash", "-c", command)
	cmd.Stdout = io.MultiWriter(&stdoutBuilder, &combinedBuilder)
	cmd.Stderr = io.MultiWriter(&stderrBuilder, &combinedBuilder)
	err := cmd.Run()
	if err != nil {
		if err, ok := err.(*exec.ExitError); ok {
		} else {
			t.Fatalf("executing %q: %v", command, err)
		}
	}
	t.Log(command + "\n" + strings.TrimSuffix(string(combinedBuilder.String()), "\n"))
	exitCode = cmd.ProcessState.ExitCode()
	stdout = stdoutBuilder.String()
	stderr = stderrBuilder.String()
	return
}

func Chdir(t testing.TB, dir string) {
	t.Helper()

	err := os.Chdir(dir)
	if err != nil {
		t.Fatalf("chdir to %q: %v", dir, err)
	}
	t.Log("cd " + dir)
}

func MkdirTemp(t testing.TB) string {
	t.Helper()

	dir, err := os.MkdirTemp("", "")
	if err != nil {
		t.Fatalf("mkdirtemp: %v", err)
	}
	t.Log("mktemp -d\n" + dir)
	return dir
}

func Mkdir(t testing.TB, path string) {
	t.Helper()

	err := os.MkdirAll(path, 0755)
	if err != nil {
		t.Fatalf("mkdir %q: %v", path, err)
	}
}

func WriteFile(t testing.TB, filename string, data []byte) {
	t.Helper()

	err := ioutil.WriteFile(filename, []byte{}, 0644)
	if err != nil {
		t.Fatalf("writing file to %q: %v", filename, err)
	}
}
