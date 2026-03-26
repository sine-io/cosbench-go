package main

import (
	"bytes"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func mockStageAwareFixture() string {
	return filepath.Clean("../../testdata/workloads/mock-stage-aware.xml")
}

func runCLITest(t *testing.T, jsonOut, quiet bool, summaryFile string) (string, string) {
	t.Helper()

	var stdout bytes.Buffer
	var stderr bytes.Buffer
	err := runCLI(mockStageAwareFixture(), "mock", jsonOut, quiet, summaryFile, &stdout, &stderr)
	if err != nil {
		t.Fatalf("runCLI(): %v stderr=%s", err, stderr.String())
	}
	return stdout.String(), stderr.String()
}

func TestRunCLIStageAwareMockFixture(t *testing.T) {
	stdout, _ := runCLITest(t, true, false, "")
	if !strings.Contains(stdout, `"workload": "mock-stage-aware"`) {
		t.Fatalf("unexpected stdout: %s", stdout)
	}
	if strings.Contains(stdout, `"errors": 0`) {
		return
	}
	t.Fatalf("expected zero errors in summary: %s", stdout)
}

func TestRunCLIWithFWorkloadAlias(t *testing.T) {
	stdout, _ := runCLITest(t, true, false, "")
	if !strings.HasPrefix(strings.TrimSpace(stdout), "{") {
		t.Fatalf("expected pure json output: %s", stdout)
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
	workload, backend, jsonOut, quiet, summaryFile, err := parseCLIArgs([]string{"testdata/workloads/mock-stage-aware.xml", "-backend", "mock", "-json", "-quiet"})
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
	if !quiet {
		t.Fatal("expected quiet mode")
	}
	if summaryFile != "" {
		t.Fatalf("summaryFile = %q", summaryFile)
	}
}

func TestRunCLIQuietSuppressesProgressOutput(t *testing.T) {
	stdout, stderr := runCLITest(t, true, true, "")
	if strings.TrimSpace(stderr) != "" {
		t.Fatalf("expected no progress output, got: %s", stderr)
	}
	if !strings.HasPrefix(strings.TrimSpace(stdout), "{") {
		t.Fatalf("expected pure json output: %s", stdout)
	}
}

func TestParseCLIArgsSupportsSummaryFile(t *testing.T) {
	workload, backend, jsonOut, quiet, summaryFile, err := parseCLIArgs([]string{"testdata/workloads/mock-stage-aware.xml", "-backend", "mock", "-quiet", "-summary-file", "out/summary.json"})
	if err != nil {
		t.Fatal(err)
	}
	if workload != "testdata/workloads/mock-stage-aware.xml" {
		t.Fatalf("workload = %q", workload)
	}
	if backend != "mock" {
		t.Fatalf("backend = %q", backend)
	}
	if jsonOut {
		t.Fatal("did not expect json mode")
	}
	if !quiet {
		t.Fatal("expected quiet mode")
	}
	if summaryFile != "out/summary.json" {
		t.Fatalf("summaryFile = %q", summaryFile)
	}
}

func TestCLIWritesSummaryFile(t *testing.T) {
	goBin, err := exec.LookPath("go")
	if err != nil {
		t.Fatalf("look path go: %v", err)
	}

	rootDir := filepath.Clean("../..")
	summaryFile := filepath.Join(t.TempDir(), "summary.json")
	cmd := exec.Command(goBin, "run", "./cmd/cosbench-go", "testdata/workloads/mock-stage-aware.xml", "-backend", "mock", "-quiet", "-summary-file", summaryFile)
	cmd.Dir = rootDir
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("go run failed: %v\n%s", err, output)
	}

	data, err := os.ReadFile(summaryFile)
	if err != nil {
		t.Fatalf("read summary file: %v", err)
	}
	if !strings.Contains(string(data), `"workload": "mock-stage-aware"`) {
		t.Fatalf("unexpected summary file: %s", data)
	}
}
