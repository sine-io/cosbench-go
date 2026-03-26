package main

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func TestListCompareLocalFixturesRejectsMalformedManifestGracefully(t *testing.T) {
	pythonBin := mustLookPath(t, "python3")
	manifestDir := t.TempDir()
	manifestPath := filepath.Join(manifestDir, "compare-local-fixtures.txt")
	if err := os.WriteFile(manifestPath, []byte("ok testdata/workloads/mock-stage-aware.xml\ninvalid-line\n"), 0o644); err != nil {
		t.Fatalf("write manifest: %v", err)
	}

	scriptPath, err := filepath.Abs(filepath.Clean("../../scripts/list_compare_local_fixtures.py"))
	if err != nil {
		t.Fatalf("abs script path: %v", err)
	}
	cmd := exec.Command(pythonBin, scriptPath, manifestPath)
	cmd.Dir = repoRootDir()
	output := string(runCommandFailure(t, cmd))

	if strings.Contains(output, "Traceback") {
		t.Fatalf("unexpected traceback: %s", output)
	}
	if !strings.Contains(output, "invalid compare-local manifest line 2") {
		t.Fatalf("unexpected output: %s", output)
	}
	if !strings.Contains(output, "invalid-line") {
		t.Fatalf("unexpected output: %s", output)
	}
}

func TestBuildCompareLocalIndexRejectsMissingFixtureSummaryGracefully(t *testing.T) {
	pythonBin := mustLookPath(t, "python3")
	manifestDir := t.TempDir()
	manifestPath := filepath.Join(manifestDir, "compare-local-fixtures.txt")
	outputDir := filepath.Join(manifestDir, "out")
	if err := os.WriteFile(manifestPath, []byte("mock-stage-aware testdata/workloads/mock-stage-aware.xml\n"), 0o644); err != nil {
		t.Fatalf("write manifest: %v", err)
	}

	scriptPath, err := filepath.Abs(filepath.Clean("../../scripts/build_compare_local_index.py"))
	if err != nil {
		t.Fatalf("abs script path: %v", err)
	}
	cmd := exec.Command(pythonBin, scriptPath, manifestPath, outputDir)
	cmd.Dir = repoRootDir()
	output := string(runCommandFailure(t, cmd))

	if strings.Contains(output, "Traceback") {
		t.Fatalf("unexpected traceback: %s", output)
	}
	if !strings.Contains(output, "missing compare-local summary for fixture mock-stage-aware") {
		t.Fatalf("unexpected output: %s", output)
	}
	if !strings.Contains(output, filepath.Join(outputDir, "mock-stage-aware.json")) {
		t.Fatalf("unexpected output: %s", output)
	}
}

func TestListCompareLocalFixturesRejectsDuplicateFixtureNamesGracefully(t *testing.T) {
	pythonBin := mustLookPath(t, "python3")
	manifestDir := t.TempDir()
	manifestPath := filepath.Join(manifestDir, "compare-local-fixtures.txt")
	data := "" +
		"mock-stage-aware testdata/workloads/mock-stage-aware.xml\n" +
		"mock-stage-aware testdata/workloads/xml-splitrw-subset.xml\n"
	if err := os.WriteFile(manifestPath, []byte(data), 0o644); err != nil {
		t.Fatalf("write manifest: %v", err)
	}

	scriptPath, err := filepath.Abs(filepath.Clean("../../scripts/list_compare_local_fixtures.py"))
	if err != nil {
		t.Fatalf("abs script path: %v", err)
	}
	cmd := exec.Command(pythonBin, scriptPath, manifestPath)
	cmd.Dir = repoRootDir()
	output := string(runCommandFailure(t, cmd))

	if strings.Contains(output, "Traceback") {
		t.Fatalf("unexpected traceback: %s", output)
	}
	if !strings.Contains(output, "duplicate compare-local fixture name 'mock-stage-aware'") {
		t.Fatalf("unexpected output: %s", output)
	}
	if !strings.Contains(output, "line 2") || !strings.Contains(output, "line 1") {
		t.Fatalf("unexpected output: %s", output)
	}
}

func TestBuildCompareLocalIndexCreatesMissingOutputDirForEmptyManifest(t *testing.T) {
	pythonBin := mustLookPath(t, "python3")
	manifestDir := t.TempDir()
	manifestPath := filepath.Join(manifestDir, "compare-local-fixtures.txt")
	outputDir := filepath.Join(manifestDir, "nested", "out")
	if err := os.WriteFile(manifestPath, []byte("# comment only\n"), 0o644); err != nil {
		t.Fatalf("write manifest: %v", err)
	}

	scriptPath, err := filepath.Abs(filepath.Clean("../../scripts/build_compare_local_index.py"))
	if err != nil {
		t.Fatalf("abs script path: %v", err)
	}
	cmd := exec.Command(pythonBin, scriptPath, manifestPath, outputDir)
	cmd.Dir = repoRootDir()
	output := string(runCommandSuccess(t, cmd))

	if strings.TrimSpace(output) != "" {
		t.Fatalf("unexpected stdout: %s", output)
	}
	indexData := mustReadFile(t, filepath.Join(outputDir, "index.json"))
	summaryData := mustReadFile(t, filepath.Join(outputDir, "summary.md"))

	var payload struct {
		Meta struct {
			Filter       string `json:"filter"`
			FixtureCount int    `json:"fixture_count"`
		} `json:"meta"`
		Fixtures []any `json:"fixtures"`
	}
	mustUnmarshalJSON(t, indexData, &payload)
	if payload.Meta.Filter != "all" || payload.Meta.FixtureCount != 0 || len(payload.Fixtures) != 0 {
		t.Fatalf("unexpected index payload: %#v", payload)
	}
	if !strings.Contains(string(summaryData), "Artifact directory: `"+outputDir+"`") {
		t.Fatalf("unexpected summary: %s", summaryData)
	}
	if !strings.Contains(string(summaryData), "Fixture count: 0") {
		t.Fatalf("unexpected summary: %s", summaryData)
	}
}

func TestBuildCompareLocalIndexRejectsMalformedFixtureSummaryGracefully(t *testing.T) {
	pythonBin := mustLookPath(t, "python3")
	manifestDir := t.TempDir()
	manifestPath := filepath.Join(manifestDir, "compare-local-fixtures.txt")
	outputDir := filepath.Join(manifestDir, "out")
	if err := os.WriteFile(manifestPath, []byte("mock-stage-aware testdata/workloads/mock-stage-aware.xml\n"), 0o644); err != nil {
		t.Fatalf("write manifest: %v", err)
	}
	if err := os.MkdirAll(outputDir, 0o755); err != nil {
		t.Fatalf("mkdir output dir: %v", err)
	}
	if err := os.WriteFile(filepath.Join(outputDir, "mock-stage-aware.json"), []byte("{\"stages\":1,\"works\":1,\"samples\":1}\n"), 0o644); err != nil {
		t.Fatalf("write malformed summary: %v", err)
	}

	scriptPath, err := filepath.Abs(filepath.Clean("../../scripts/build_compare_local_index.py"))
	if err != nil {
		t.Fatalf("abs script path: %v", err)
	}
	cmd := exec.Command(pythonBin, scriptPath, manifestPath, outputDir)
	cmd.Dir = repoRootDir()
	output := string(runCommandFailure(t, cmd))

	if strings.Contains(output, "Traceback") {
		t.Fatalf("unexpected traceback: %s", output)
	}
	if !strings.Contains(output, "invalid compare-local summary for fixture mock-stage-aware") {
		t.Fatalf("unexpected output: %s", output)
	}
	if !strings.Contains(output, "missing required field errors") {
		t.Fatalf("unexpected output: %s", output)
	}
	if !strings.Contains(output, filepath.Join(outputDir, "mock-stage-aware.json")) {
		t.Fatalf("unexpected output: %s", output)
	}
}
