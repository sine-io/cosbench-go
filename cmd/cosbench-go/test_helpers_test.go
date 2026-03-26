package main

import (
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func mustLookPath(t *testing.T, name string) string {
	t.Helper()
	path, err := exec.LookPath(name)
	if err != nil {
		t.Fatalf("look path %s: %v", name, err)
	}
	return path
}

func repoRootDir() string {
	return filepath.Clean("../..")
}

func repoRootAbs(t *testing.T) string {
	t.Helper()
	path, err := filepath.Abs(repoRootDir())
	if err != nil {
		t.Fatalf("abs root dir: %v", err)
	}
	return path
}

func runCommandSuccess(t *testing.T, cmd *exec.Cmd) []byte {
	t.Helper()
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("%s failed: %v\n%s", strings.Join(cmd.Args, " "), err, output)
	}
	return output
}

func runCommandFailure(t *testing.T, cmd *exec.Cmd) []byte {
	t.Helper()
	output, err := cmd.CombinedOutput()
	if err == nil {
		t.Fatalf("expected %s to fail\n%s", strings.Join(cmd.Args, " "), output)
	}
	return output
}

func mustReadFile(t *testing.T, path string) []byte {
	t.Helper()
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read %s: %v", path, err)
	}
	return data
}

func mustUnmarshalJSON(t *testing.T, data []byte, target any) {
	t.Helper()
	if err := json.Unmarshal(data, target); err != nil {
		t.Fatalf("unmarshal output: %v\n%s", err, data)
	}
}
