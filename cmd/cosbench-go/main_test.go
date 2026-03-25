package main

import (
	"bytes"
	"path/filepath"
	"strings"
	"testing"
)

func TestRunCLIStageAwareMockFixture(t *testing.T) {
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	err := runCLI(filepath.Clean("../../testdata/workloads/mock-stage-aware.xml"), "mock", true, &stdout, &stderr)
	if err != nil {
		t.Fatalf("runCLI(): %v stderr=%s", err, stderr.String())
	}
	if !strings.Contains(stdout.String(), `"workload": "mock-stage-aware"`) {
		t.Fatalf("unexpected stdout: %s", stdout.String())
	}
	if strings.Contains(stdout.String(), `"errors": 0`) {
		return
	}
	t.Fatalf("expected zero errors in summary: %s", stdout.String())
}

func TestRunCLIWithFWorkloadAlias(t *testing.T) {
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	err := runCLI(filepath.Clean("../../testdata/workloads/mock-stage-aware.xml"), "mock", true, &stdout, &stderr)
	if err != nil {
		t.Fatalf("runCLI(): %v stderr=%s", err, stderr.String())
	}
	if !strings.HasPrefix(strings.TrimSpace(stdout.String()), "{") {
		t.Fatalf("expected pure json output: %s", stdout.String())
	}
}

func TestResolveWorkloadPathPriority(t *testing.T) {
	path, err := resolveWorkloadPath("flag.xml", "short.xml", []string{"positional.xml"})
	if err != nil {
		t.Fatal(err)
	}
	if path != "flag.xml" {
		t.Fatalf("path = %q", path)
	}

	path, err = resolveWorkloadPath("", "short.xml", []string{"positional.xml"})
	if err != nil {
		t.Fatal(err)
	}
	if path != "short.xml" {
		t.Fatalf("path = %q", path)
	}

	path, err = resolveWorkloadPath("", "", []string{"positional.xml"})
	if err != nil {
		t.Fatal(err)
	}
	if path != "positional.xml" {
		t.Fatalf("path = %q", path)
	}
}

func TestResolveWorkloadPathRequiresInput(t *testing.T) {
	if _, err := resolveWorkloadPath("", "", nil); err == nil {
		t.Fatal("expected error")
	}
}

func TestParseCLIArgsSupportsPositionalPathWithTrailingFlags(t *testing.T) {
	workload, backend, jsonOut, err := parseCLIArgs([]string{"testdata/workloads/mock-stage-aware.xml", "-backend", "mock", "-json"})
	if err != nil {
		t.Fatal(err)
	}
	if workload != "testdata/workloads/mock-stage-aware.xml" {
		t.Fatalf("workload = %q", workload)
	}
	if backend != "mock" {
		t.Fatalf("backend = %q", backend)
	}
	if !jsonOut {
		t.Fatal("expected json mode")
	}
}
