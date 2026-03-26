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

func TestListCompareLocalFixturesRejectsMissingManifestGracefully(t *testing.T) {
	pythonBin := mustLookPath(t, "python3")
	manifestPath := filepath.Join(t.TempDir(), "missing-compare-local-fixtures.txt")

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
	if !strings.Contains(output, "compare-local manifest not found") {
		t.Fatalf("unexpected output: %s", output)
	}
	if !strings.Contains(output, manifestPath) {
		t.Fatalf("unexpected output: %s", output)
	}
}

func TestListCompareLocalFixturesRendersNonASCIIManifestPathGracefully(t *testing.T) {
	pythonBin := mustLookPath(t, "python3")
	manifestDir := t.TempDir()
	manifestPath := filepath.Join(manifestDir, "清单.txt")

	scriptPath, err := filepath.Abs(filepath.Clean("../../scripts/list_compare_local_fixtures.py"))
	if err != nil {
		t.Fatalf("abs script path: %v", err)
	}
	cmd := exec.Command(pythonBin, scriptPath, manifestPath)
	cmd.Dir = repoRootDir()
	cmd.Env = append(
		os.Environ(),
		"LC_ALL=C",
		"LANG=C",
		"PYTHONCOERCECLOCALE=0",
		"PYTHONUTF8=0",
	)
	output := string(runCommandFailure(t, cmd))

	if strings.Contains(output, "Traceback") {
		t.Fatalf("unexpected traceback: %s", output)
	}
	if strings.Contains(output, "\\udc") {
		t.Fatalf("unexpected surrogate escapes: %s", output)
	}
	if !strings.Contains(output, "compare-local manifest not found") {
		t.Fatalf("unexpected output: %s", output)
	}
	if !strings.Contains(output, "清单.txt") {
		t.Fatalf("unexpected output: %s", output)
	}
}

func TestListCompareLocalFixturesRejectsUnreadableManifestGracefully(t *testing.T) {
	pythonBin := mustLookPath(t, "python3")
	manifestPath := t.TempDir()

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
	if !strings.Contains(output, "unable to read compare-local manifest") {
		t.Fatalf("unexpected output: %s", output)
	}
	if !strings.Contains(output, manifestPath) {
		t.Fatalf("unexpected output: %s", output)
	}
}

func TestListCompareLocalFixturesRendersUnreadableNonASCIIManifestPathGracefully(t *testing.T) {
	pythonBin := mustLookPath(t, "python3")
	manifestDir := t.TempDir()
	manifestPath := filepath.Join(manifestDir, "清单.txt")
	if err := os.MkdirAll(manifestPath, 0o755); err != nil {
		t.Fatalf("mkdir manifest dir: %v", err)
	}

	scriptPath, err := filepath.Abs(filepath.Clean("../../scripts/list_compare_local_fixtures.py"))
	if err != nil {
		t.Fatalf("abs script path: %v", err)
	}
	cmd := exec.Command(pythonBin, scriptPath, manifestPath)
	cmd.Dir = repoRootDir()
	cmd.Env = append(
		os.Environ(),
		"LC_ALL=C",
		"LANG=C",
		"PYTHONCOERCECLOCALE=0",
		"PYTHONUTF8=0",
	)
	output := string(runCommandFailure(t, cmd))

	if strings.Contains(output, "Traceback") {
		t.Fatalf("unexpected traceback: %s", output)
	}
	if strings.Contains(output, "\\udc") {
		t.Fatalf("unexpected surrogate escapes: %s", output)
	}
	if !strings.Contains(output, "unable to read compare-local manifest") {
		t.Fatalf("unexpected output: %s", output)
	}
	if !strings.Contains(output, "清单.txt") {
		t.Fatalf("unexpected output: %s", output)
	}
}

func TestListCompareLocalFixturesRejectsInvalidManifestEncodingGracefully(t *testing.T) {
	pythonBin := mustLookPath(t, "python3")
	manifestDir := t.TempDir()
	manifestPath := filepath.Join(manifestDir, "compare-local-fixtures.txt")
	if err := os.WriteFile(manifestPath, []byte{0xff, 0xfe, 0x00, 'b', 'a', 'd', '\n'}, 0o644); err != nil {
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
	if !strings.Contains(output, "unable to decode compare-local manifest") {
		t.Fatalf("unexpected output: %s", output)
	}
	if !strings.Contains(output, manifestPath) {
		t.Fatalf("unexpected output: %s", output)
	}
}

func TestListCompareLocalFixturesAcceptsUTF8BOMManifest(t *testing.T) {
	pythonBin := mustLookPath(t, "python3")
	manifestDir := t.TempDir()
	manifestPath := filepath.Join(manifestDir, "compare-local-fixtures.txt")
	if err := os.WriteFile(manifestPath, []byte{0xef, 0xbb, 0xbf, 'm', 'o', 'c', 'k', '-', 's', 't', 'a', 'g', 'e', '-', 'a', 'w', 'a', 'r', 'e', ' ', 't', 'e', 's', 't', 'd', 'a', 't', 'a', '/', 'w', 'o', 'r', 'k', 'l', 'o', 'a', 'd', 's', '/', 'm', 'o', 'c', 'k', '-', 's', 't', 'a', 'g', 'e', '-', 'a', 'w', 'a', 'r', 'e', '.', 'x', 'm', 'l', '\n'}, 0o644); err != nil {
		t.Fatalf("write manifest: %v", err)
	}

	scriptPath, err := filepath.Abs(filepath.Clean("../../scripts/list_compare_local_fixtures.py"))
	if err != nil {
		t.Fatalf("abs script path: %v", err)
	}
	cmd := exec.Command(pythonBin, scriptPath, manifestPath, "--names")
	cmd.Dir = repoRootDir()
	output := string(runCommandSuccess(t, cmd))

	lines := strings.Fields(strings.TrimSpace(output))
	if len(lines) != 1 || lines[0] != "mock-stage-aware" {
		t.Fatalf("unexpected output: %q", output)
	}
}

func TestListCompareLocalFixturesWritesNamesAndPairsWithExplicitUTF8Stdout(t *testing.T) {
	for _, mode := range []string{"--names", "--pairs"} {
		t.Run(mode, func(t *testing.T) {
			pythonBin := mustLookPath(t, "python3")
			manifestDir := t.TempDir()
			manifestPath := filepath.Join(manifestDir, "compare-local-fixtures.txt")
			if err := os.WriteFile(manifestPath, []byte("样例 testdata/workloads/测试.xml\n"), 0o644); err != nil {
				t.Fatalf("write manifest: %v", err)
			}

			scriptPath, err := filepath.Abs(filepath.Clean("../../scripts/list_compare_local_fixtures.py"))
			if err != nil {
				t.Fatalf("abs script path: %v", err)
			}
			cmd := exec.Command(pythonBin, scriptPath, manifestPath, mode)
			cmd.Dir = repoRootDir()
			cmd.Env = append(
				os.Environ(),
				"LC_ALL=C",
				"LANG=C",
				"PYTHONCOERCECLOCALE=0",
				"PYTHONUTF8=0",
			)
			output := string(runCommandSuccess(t, cmd))

			if !strings.Contains(output, "样例") {
				t.Fatalf("unexpected output: %q", output)
			}
			if mode == "--pairs" && !strings.Contains(output, "测试.xml") {
				t.Fatalf("unexpected output: %q", output)
			}
		})
	}
}

func TestListCompareLocalFixturesRejectsUnknownOptionGracefully(t *testing.T) {
	pythonBin := mustLookPath(t, "python3")
	manifestDir := t.TempDir()
	manifestPath := filepath.Join(manifestDir, "compare-local-fixtures.txt")
	if err := os.WriteFile(manifestPath, []byte("# comment only\n"), 0o644); err != nil {
		t.Fatalf("write manifest: %v", err)
	}

	scriptPath, err := filepath.Abs(filepath.Clean("../../scripts/list_compare_local_fixtures.py"))
	if err != nil {
		t.Fatalf("abs script path: %v", err)
	}
	cmd := exec.Command(pythonBin, scriptPath, manifestPath, "--bogus")
	cmd.Dir = repoRootDir()
	output := string(runCommandFailure(t, cmd))

	if strings.Contains(output, "Traceback") {
		t.Fatalf("unexpected traceback: %s", output)
	}
	if !strings.Contains(output, "unknown option: --bogus") {
		t.Fatalf("unexpected output: %s", output)
	}
}

func TestListCompareLocalFixturesRejectsDuplicateOutputModeOptionGracefully(t *testing.T) {
	for _, option := range []string{"--names", "--pairs"} {
		t.Run(option, func(t *testing.T) {
			pythonBin := mustLookPath(t, "python3")
			manifestDir := t.TempDir()
			manifestPath := filepath.Join(manifestDir, "compare-local-fixtures.txt")
			if err := os.WriteFile(manifestPath, []byte("# comment only\n"), 0o644); err != nil {
				t.Fatalf("write manifest: %v", err)
			}

			scriptPath, err := filepath.Abs(filepath.Clean("../../scripts/list_compare_local_fixtures.py"))
			if err != nil {
				t.Fatalf("abs script path: %v", err)
			}
			cmd := exec.Command(pythonBin, scriptPath, manifestPath, option, option)
			cmd.Dir = repoRootDir()
			output := string(runCommandFailure(t, cmd))

			if strings.Contains(output, "Traceback") {
				t.Fatalf("unexpected traceback: %s", output)
			}
			if !strings.Contains(output, "duplicate option: "+option) {
				t.Fatalf("unexpected output: %s", output)
			}
		})
	}
}

func TestListCompareLocalFixturesRejectsSeparatorOnlyFilterGracefully(t *testing.T) {
	pythonBin := mustLookPath(t, "python3")

	scriptPath, err := filepath.Abs(filepath.Clean("../../scripts/list_compare_local_fixtures.py"))
	if err != nil {
		t.Fatalf("abs script path: %v", err)
	}
	cmd := exec.Command(pythonBin, scriptPath, "testdata/workloads/compare-local-fixtures.txt", ",")
	cmd.Dir = repoRootDir()
	output := string(runCommandFailure(t, cmd))

	if strings.Contains(output, "Traceback") {
		t.Fatalf("unexpected traceback: %s", output)
	}
	if !strings.Contains(output, "invalid compare-local filter") {
		t.Fatalf("unexpected output: %s", output)
	}
	if !strings.Contains(output, "filter did not include any fixture names") {
		t.Fatalf("unexpected output: %s", output)
	}
}

func TestListCompareLocalFixturesRejectsExtraFilterArgsGracefully(t *testing.T) {
	pythonBin := mustLookPath(t, "python3")
	manifestDir := t.TempDir()
	manifestPath := filepath.Join(manifestDir, "compare-local-fixtures.txt")
	if err := os.WriteFile(manifestPath, []byte("# comment only\n"), 0o644); err != nil {
		t.Fatalf("write manifest: %v", err)
	}

	scriptPath, err := filepath.Abs(filepath.Clean("../../scripts/list_compare_local_fixtures.py"))
	if err != nil {
		t.Fatalf("abs script path: %v", err)
	}
	cmd := exec.Command(pythonBin, scriptPath, manifestPath, "foo", "bar")
	cmd.Dir = repoRootDir()
	output := string(runCommandFailure(t, cmd))

	if strings.Contains(output, "Traceback") {
		t.Fatalf("unexpected traceback: %s", output)
	}
	if !strings.Contains(output, "expected at most one filter argument") {
		t.Fatalf("unexpected output: %s", output)
	}
	if !strings.Contains(output, "foo") || !strings.Contains(output, "bar") {
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

func TestBuildCompareLocalIndexRejectsNonASCIISummaryPathGracefully(t *testing.T) {
	pythonBin := mustLookPath(t, "python3")
	manifestDir := t.TempDir()
	manifestPath := filepath.Join(manifestDir, "compare-local-fixtures.txt")
	outputDir := filepath.Join(manifestDir, "out")
	if err := os.WriteFile(manifestPath, []byte("样例 testdata/workloads/mock-stage-aware.xml\n"), 0o644); err != nil {
		t.Fatalf("write manifest: %v", err)
	}
	if err := os.MkdirAll(outputDir, 0o755); err != nil {
		t.Fatalf("mkdir output dir: %v", err)
	}

	scriptPath, err := filepath.Abs(filepath.Clean("../../scripts/build_compare_local_index.py"))
	if err != nil {
		t.Fatalf("abs script path: %v", err)
	}
	cmd := exec.Command(pythonBin, scriptPath, manifestPath, outputDir)
	cmd.Dir = repoRootDir()
	cmd.Env = append(
		os.Environ(),
		"LC_ALL=C",
		"LANG=C",
		"PYTHONCOERCECLOCALE=0",
		"PYTHONUTF8=0",
	)
	output := string(runCommandFailure(t, cmd))

	if strings.Contains(output, "Traceback") {
		t.Fatalf("unexpected traceback: %s", output)
	}
	if !strings.Contains(output, "unable to access compare-local summary path for fixture") {
		t.Fatalf("unexpected output: %s", output)
	}
	if !strings.Contains(output, ".json") {
		t.Fatalf("unexpected output: %s", output)
	}
}

func TestBuildCompareLocalIndexRejectsUnreadableFixtureSummaryGracefully(t *testing.T) {
	pythonBin := mustLookPath(t, "python3")
	manifestDir := t.TempDir()
	manifestPath := filepath.Join(manifestDir, "compare-local-fixtures.txt")
	outputDir := filepath.Join(manifestDir, "out")
	if err := os.WriteFile(manifestPath, []byte("mock-stage-aware testdata/workloads/mock-stage-aware.xml\n"), 0o644); err != nil {
		t.Fatalf("write manifest: %v", err)
	}
	if err := os.MkdirAll(filepath.Join(outputDir, "mock-stage-aware.json"), 0o755); err != nil {
		t.Fatalf("mkdir summary dir: %v", err)
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
	if !strings.Contains(output, "unable to read compare-local summary for fixture mock-stage-aware") {
		t.Fatalf("unexpected output: %s", output)
	}
	if !strings.Contains(output, filepath.Join(outputDir, "mock-stage-aware.json")) {
		t.Fatalf("unexpected output: %s", output)
	}
}

func TestBuildCompareLocalIndexRejectsUnreadableFixtureSummaryInNonASCIIOutputDirGracefully(t *testing.T) {
	pythonBin := mustLookPath(t, "python3")
	manifestDir := t.TempDir()
	manifestPath := filepath.Join(manifestDir, "compare-local-fixtures.txt")
	outputDir := filepath.Join(manifestDir, "输出")
	if err := os.WriteFile(manifestPath, []byte("fixture testdata/workloads/mock-stage-aware.xml\n"), 0o644); err != nil {
		t.Fatalf("write manifest: %v", err)
	}
	if err := os.MkdirAll(filepath.Join(outputDir, "fixture.json"), 0o755); err != nil {
		t.Fatalf("mkdir summary dir: %v", err)
	}

	scriptPath, err := filepath.Abs(filepath.Clean("../../scripts/build_compare_local_index.py"))
	if err != nil {
		t.Fatalf("abs script path: %v", err)
	}
	cmd := exec.Command(pythonBin, scriptPath, manifestPath, outputDir)
	cmd.Dir = repoRootDir()
	cmd.Env = append(
		os.Environ(),
		"LC_ALL=C",
		"LANG=C",
		"PYTHONCOERCECLOCALE=0",
		"PYTHONUTF8=0",
	)
	output := string(runCommandFailure(t, cmd))

	if strings.Contains(output, "Traceback") {
		t.Fatalf("unexpected traceback: %s", output)
	}
	if strings.Contains(output, "\\udc") {
		t.Fatalf("unexpected surrogate escapes: %s", output)
	}
	if !strings.Contains(output, "unable to read compare-local summary for fixture fixture") {
		t.Fatalf("unexpected output: %s", output)
	}
	if !strings.Contains(output, "输出") || !strings.Contains(output, "fixture.json") {
		t.Fatalf("unexpected output: %s", output)
	}
}

func TestBuildCompareLocalIndexRejectsInvalidFixtureSummaryEncodingGracefully(t *testing.T) {
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
	if err := os.WriteFile(filepath.Join(outputDir, "mock-stage-aware.json"), []byte{0xff, 0xfe, 0x00, 'b', 'a', 'd', '\n'}, 0o644); err != nil {
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
	if !strings.Contains(output, "unable to decode compare-local summary for fixture mock-stage-aware") {
		t.Fatalf("unexpected output: %s", output)
	}
	if !strings.Contains(output, filepath.Join(outputDir, "mock-stage-aware.json")) {
		t.Fatalf("unexpected output: %s", output)
	}
}

func TestBuildCompareLocalIndexAcceptsUTF8BOMFixtureSummary(t *testing.T) {
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
	summaryPath := filepath.Join(outputDir, "mock-stage-aware.json")
	summaryPayload := append(
		[]byte{0xef, 0xbb, 0xbf},
		[]byte("{\"stages\":1,\"works\":1,\"samples\":1,\"errors\":0}\n")...,
	)
	if err := os.WriteFile(summaryPath, summaryPayload, 0o644); err != nil {
		t.Fatalf("write summary: %v", err)
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
	var payload struct {
		Meta struct {
			FixtureCount int `json:"fixture_count"`
		} `json:"meta"`
		Fixtures []struct {
			Name   string `json:"name"`
			Errors int    `json:"errors"`
		} `json:"fixtures"`
	}
	mustUnmarshalJSON(t, indexData, &payload)
	if payload.Meta.FixtureCount != 1 || len(payload.Fixtures) != 1 {
		t.Fatalf("unexpected payload: %#v", payload)
	}
	if payload.Fixtures[0].Name != "mock-stage-aware" || payload.Fixtures[0].Errors != 0 {
		t.Fatalf("unexpected payload: %#v", payload)
	}
}

func TestBuildCompareLocalIndexWritesArtifactsWithExplicitUTF8Encoding(t *testing.T) {
	pythonBin := mustLookPath(t, "python3")
	manifestDir := t.TempDir()
	manifestPath := filepath.Join(manifestDir, "compare-local-fixtures.txt")
	outputDir := filepath.Join(manifestDir, "out")
	if err := os.WriteFile(manifestPath, []byte("fixture testdata/workloads/测试.xml\n"), 0o644); err != nil {
		t.Fatalf("write manifest: %v", err)
	}
	if err := os.MkdirAll(outputDir, 0o755); err != nil {
		t.Fatalf("mkdir output dir: %v", err)
	}
	if err := os.WriteFile(filepath.Join(outputDir, "fixture.json"), []byte("{\"stages\":1,\"works\":1,\"samples\":1,\"errors\":0}\n"), 0o644); err != nil {
		t.Fatalf("write summary: %v", err)
	}

	scriptPath, err := filepath.Abs(filepath.Clean("../../scripts/build_compare_local_index.py"))
	if err != nil {
		t.Fatalf("abs script path: %v", err)
	}
	cmd := exec.Command(pythonBin, scriptPath, manifestPath, outputDir)
	cmd.Dir = repoRootDir()
	cmd.Env = append(
		os.Environ(),
		"LC_ALL=C",
		"LANG=C",
		"PYTHONCOERCECLOCALE=0",
		"PYTHONUTF8=0",
	)
	output := string(runCommandSuccess(t, cmd))

	if strings.TrimSpace(output) != "" {
		t.Fatalf("unexpected stdout: %s", output)
	}

	summaryData := mustReadFile(t, filepath.Join(outputDir, "summary.md"))
	if !strings.Contains(string(summaryData), "测试.xml") {
		t.Fatalf("unexpected summary: %s", summaryData)
	}
}

func TestBuildCompareLocalIndexRejectsUnknownFilterGracefully(t *testing.T) {
	pythonBin := mustLookPath(t, "python3")
	manifestDir := t.TempDir()
	manifestPath := filepath.Join(manifestDir, "compare-local-fixtures.txt")
	outputDir := filepath.Join(manifestDir, "out")
	if err := os.WriteFile(manifestPath, []byte("# comment only\n"), 0o644); err != nil {
		t.Fatalf("write manifest: %v", err)
	}

	scriptPath, err := filepath.Abs(filepath.Clean("../../scripts/build_compare_local_index.py"))
	if err != nil {
		t.Fatalf("abs script path: %v", err)
	}
	cmd := exec.Command(pythonBin, scriptPath, manifestPath, outputDir, "does-not-exist")
	cmd.Dir = repoRootDir()
	output := string(runCommandFailure(t, cmd))

	if strings.Contains(output, "Traceback") {
		t.Fatalf("unexpected traceback: %s", output)
	}
	if !strings.Contains(output, "unknown compare-local fixture: does-not-exist") {
		t.Fatalf("unexpected output: %s", output)
	}
	if !strings.Contains(output, "known fixtures:") {
		t.Fatalf("unexpected output: %s", output)
	}
}

func TestListCompareLocalFixturesRejectsMixedAllFilterGracefully(t *testing.T) {
	pythonBin := mustLookPath(t, "python3")

	scriptPath, err := filepath.Abs(filepath.Clean("../../scripts/list_compare_local_fixtures.py"))
	if err != nil {
		t.Fatalf("abs script path: %v", err)
	}
	cmd := exec.Command(pythonBin, scriptPath, "testdata/workloads/compare-local-fixtures.txt", "all,mock-stage-aware")
	cmd.Dir = repoRootDir()
	output := string(runCommandFailure(t, cmd))

	if strings.Contains(output, "Traceback") {
		t.Fatalf("unexpected traceback: %s", output)
	}
	if !strings.Contains(output, "invalid compare-local filter") {
		t.Fatalf("unexpected output: %s", output)
	}
	if !strings.Contains(output, "'all' cannot be combined with specific fixtures") {
		t.Fatalf("unexpected output: %s", output)
	}
}

func TestBuildCompareLocalIndexRejectsMixedAllFilterGracefully(t *testing.T) {
	pythonBin := mustLookPath(t, "python3")
	outputDir := filepath.Join(t.TempDir(), "out")
	if err := os.MkdirAll(outputDir, 0o755); err != nil {
		t.Fatalf("mkdir output dir: %v", err)
	}

	scriptPath, err := filepath.Abs(filepath.Clean("../../scripts/build_compare_local_index.py"))
	if err != nil {
		t.Fatalf("abs script path: %v", err)
	}
	cmd := exec.Command(pythonBin, scriptPath, "testdata/workloads/compare-local-fixtures.txt", outputDir, "all,mock-stage-aware")
	cmd.Dir = repoRootDir()
	output := string(runCommandFailure(t, cmd))

	if strings.Contains(output, "Traceback") {
		t.Fatalf("unexpected traceback: %s", output)
	}
	if !strings.Contains(output, "invalid compare-local filter") {
		t.Fatalf("unexpected output: %s", output)
	}
	if !strings.Contains(output, "'all' cannot be combined with specific fixtures") {
		t.Fatalf("unexpected output: %s", output)
	}
}

func TestValidateCompareLocalFilterRejectsUnknownOptionGracefully(t *testing.T) {
	pythonBin := mustLookPath(t, "python3")

	scriptPath, err := filepath.Abs(filepath.Clean("../../scripts/validate_compare_local_filter.py"))
	if err != nil {
		t.Fatalf("abs script path: %v", err)
	}
	cmd := exec.Command(pythonBin, scriptPath, "testdata/workloads/compare-local-fixtures.txt", "--bogus")
	cmd.Dir = repoRootDir()
	output := string(runCommandFailure(t, cmd))

	if strings.Contains(output, "Traceback") {
		t.Fatalf("unexpected traceback: %s", output)
	}
	if !strings.Contains(output, "unknown option: --bogus") {
		t.Fatalf("unexpected output: %s", output)
	}
}

func TestValidateCompareLocalFilterRejectsExtraFilterArgsGracefully(t *testing.T) {
	pythonBin := mustLookPath(t, "python3")

	scriptPath, err := filepath.Abs(filepath.Clean("../../scripts/validate_compare_local_filter.py"))
	if err != nil {
		t.Fatalf("abs script path: %v", err)
	}
	cmd := exec.Command(
		pythonBin,
		scriptPath,
		"testdata/workloads/compare-local-fixtures.txt",
		"mock-stage-aware",
		"xml-splitrw-subset",
	)
	cmd.Dir = repoRootDir()
	output := string(runCommandFailure(t, cmd))

	if strings.Contains(output, "Traceback") {
		t.Fatalf("unexpected traceback: %s", output)
	}
	if !strings.Contains(output, "expected exactly one filter argument") {
		t.Fatalf("unexpected output: %s", output)
	}
	if !strings.Contains(output, "mock-stage-aware xml-splitrw-subset") {
		t.Fatalf("unexpected output: %s", output)
	}
}

func TestValidateCompareLocalFilterRejectsSeparatorOnlyFilterGracefully(t *testing.T) {
	pythonBin := mustLookPath(t, "python3")

	scriptPath, err := filepath.Abs(filepath.Clean("../../scripts/validate_compare_local_filter.py"))
	if err != nil {
		t.Fatalf("abs script path: %v", err)
	}
	cmd := exec.Command(pythonBin, scriptPath, "testdata/workloads/compare-local-fixtures.txt", ",")
	cmd.Dir = repoRootDir()
	output := string(runCommandFailure(t, cmd))

	if strings.Contains(output, "Traceback") {
		t.Fatalf("unexpected traceback: %s", output)
	}
	if !strings.Contains(output, "invalid compare-local filter") {
		t.Fatalf("unexpected output: %s", output)
	}
	if !strings.Contains(output, "filter did not include any fixture names") {
		t.Fatalf("unexpected output: %s", output)
	}
}

func TestBuildCompareLocalIndexRejectsUnknownOptionGracefully(t *testing.T) {
	pythonBin := mustLookPath(t, "python3")
	manifestDir := t.TempDir()
	manifestPath := filepath.Join(manifestDir, "compare-local-fixtures.txt")
	outputDir := filepath.Join(manifestDir, "out")
	if err := os.WriteFile(manifestPath, []byte("# comment only\n"), 0o644); err != nil {
		t.Fatalf("write manifest: %v", err)
	}
	if err := os.MkdirAll(outputDir, 0o755); err != nil {
		t.Fatalf("mkdir output dir: %v", err)
	}

	scriptPath, err := filepath.Abs(filepath.Clean("../../scripts/build_compare_local_index.py"))
	if err != nil {
		t.Fatalf("abs script path: %v", err)
	}
	cmd := exec.Command(pythonBin, scriptPath, manifestPath, outputDir, "--bogus")
	cmd.Dir = repoRootDir()
	output := string(runCommandFailure(t, cmd))

	if strings.Contains(output, "Traceback") {
		t.Fatalf("unexpected traceback: %s", output)
	}
	if !strings.Contains(output, "unknown option: --bogus") {
		t.Fatalf("unexpected output: %s", output)
	}
}

func TestBuildCompareLocalIndexRejectsExtraFilterArgsGracefully(t *testing.T) {
	pythonBin := mustLookPath(t, "python3")
	manifestDir := t.TempDir()
	manifestPath := filepath.Join(manifestDir, "compare-local-fixtures.txt")
	outputDir := filepath.Join(manifestDir, "out")
	if err := os.WriteFile(manifestPath, []byte("# comment only\n"), 0o644); err != nil {
		t.Fatalf("write manifest: %v", err)
	}
	if err := os.MkdirAll(outputDir, 0o755); err != nil {
		t.Fatalf("mkdir output dir: %v", err)
	}

	scriptPath, err := filepath.Abs(filepath.Clean("../../scripts/build_compare_local_index.py"))
	if err != nil {
		t.Fatalf("abs script path: %v", err)
	}
	cmd := exec.Command(
		pythonBin,
		scriptPath,
		manifestPath,
		outputDir,
		"mock-stage-aware",
		"xml-splitrw-subset",
	)
	cmd.Dir = repoRootDir()
	output := string(runCommandFailure(t, cmd))

	if strings.Contains(output, "Traceback") {
		t.Fatalf("unexpected traceback: %s", output)
	}
	if !strings.Contains(output, "expected at most one filter argument") {
		t.Fatalf("unexpected output: %s", output)
	}
	if !strings.Contains(output, "mock-stage-aware xml-splitrw-subset") {
		t.Fatalf("unexpected output: %s", output)
	}
}

func TestBuildCompareLocalIndexRejectsSeparatorOnlyFilterGracefully(t *testing.T) {
	pythonBin := mustLookPath(t, "python3")
	manifestDir := t.TempDir()
	manifestPath := filepath.Join(manifestDir, "compare-local-fixtures.txt")
	outputDir := filepath.Join(manifestDir, "out")
	if err := os.WriteFile(manifestPath, []byte("# comment only\n"), 0o644); err != nil {
		t.Fatalf("write manifest: %v", err)
	}
	if err := os.MkdirAll(outputDir, 0o755); err != nil {
		t.Fatalf("mkdir output dir: %v", err)
	}

	scriptPath, err := filepath.Abs(filepath.Clean("../../scripts/build_compare_local_index.py"))
	if err != nil {
		t.Fatalf("abs script path: %v", err)
	}
	cmd := exec.Command(pythonBin, scriptPath, manifestPath, outputDir, ",")
	cmd.Dir = repoRootDir()
	output := string(runCommandFailure(t, cmd))

	if strings.Contains(output, "Traceback") {
		t.Fatalf("unexpected traceback: %s", output)
	}
	if !strings.Contains(output, "invalid compare-local filter") {
		t.Fatalf("unexpected output: %s", output)
	}
	if !strings.Contains(output, "filter did not include any fixture names") {
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

func TestListCompareLocalFixturesRejectsCaseFoldDuplicateFixtureNamesGracefully(t *testing.T) {
	pythonBin := mustLookPath(t, "python3")
	manifestDir := t.TempDir()
	manifestPath := filepath.Join(manifestDir, "compare-local-fixtures.txt")
	data := "" +
		"Mock-Stage-Aware testdata/workloads/mock-stage-aware.xml\n" +
		"mock-stage-aware testdata/workloads/mock-stage-aware.xml\n"
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
	if !strings.Contains(output, "case-insensitive duplicate compare-local fixture name 'mock-stage-aware'") {
		t.Fatalf("unexpected output: %s", output)
	}
	if !strings.Contains(output, "'Mock-Stage-Aware'") {
		t.Fatalf("unexpected output: %s", output)
	}
}

func TestListCompareLocalFixturesRejectsPathLikeFixtureNamesGracefully(t *testing.T) {
	pythonBin := mustLookPath(t, "python3")
	manifestDir := t.TempDir()
	manifestPath := filepath.Join(manifestDir, "compare-local-fixtures.txt")
	data := "../escape testdata/workloads/mock-stage-aware.xml\n"
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
	if !strings.Contains(output, "invalid compare-local fixture name '../escape'") {
		t.Fatalf("unexpected output: %s", output)
	}
	if !strings.Contains(output, "must not contain path separators or dot-path segments") {
		t.Fatalf("unexpected output: %s", output)
	}
}

func TestListCompareLocalFixturesRejectsReservedFixtureNameAllGracefully(t *testing.T) {
	pythonBin := mustLookPath(t, "python3")
	manifestDir := t.TempDir()
	manifestPath := filepath.Join(manifestDir, "compare-local-fixtures.txt")
	data := "all testdata/workloads/mock-stage-aware.xml\n"
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
	if !strings.Contains(output, "invalid compare-local fixture name 'all'") {
		t.Fatalf("unexpected output: %s", output)
	}
	if !strings.Contains(output, "reserved for the all-fixtures selector") {
		t.Fatalf("unexpected output: %s", output)
	}
}

func TestListCompareLocalFixturesRejectsOptionLikeFixtureNamesGracefully(t *testing.T) {
	pythonBin := mustLookPath(t, "python3")
	manifestDir := t.TempDir()
	manifestPath := filepath.Join(manifestDir, "compare-local-fixtures.txt")
	data := "--bad testdata/workloads/mock-stage-aware.xml\n"
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
	if !strings.Contains(output, "invalid compare-local fixture name '--bad'") {
		t.Fatalf("unexpected output: %s", output)
	}
	if !strings.Contains(output, "must not start with --") {
		t.Fatalf("unexpected output: %s", output)
	}
}

func TestListCompareLocalFixturesRejectsCommaSeparatedFixtureNamesGracefully(t *testing.T) {
	pythonBin := mustLookPath(t, "python3")
	manifestDir := t.TempDir()
	manifestPath := filepath.Join(manifestDir, "compare-local-fixtures.txt")
	data := "foo,bar testdata/workloads/mock-stage-aware.xml\n"
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
	if !strings.Contains(output, "invalid compare-local fixture name 'foo,bar'") {
		t.Fatalf("unexpected output: %s", output)
	}
	if !strings.Contains(output, "must not contain commas") {
		t.Fatalf("unexpected output: %s", output)
	}
}

func TestListCompareLocalFixturesRejectsFilesystemSpecialFixtureNamesGracefully(t *testing.T) {
	pythonBin := mustLookPath(t, "python3")
	manifestDir := t.TempDir()
	manifestPath := filepath.Join(manifestDir, "compare-local-fixtures.txt")
	data := "foo:bar testdata/workloads/mock-stage-aware.xml\n"
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
	if !strings.Contains(output, "invalid compare-local fixture name 'foo:bar'") {
		t.Fatalf("unexpected output: %s", output)
	}
	if !strings.Contains(output, "must not contain filesystem-special characters") {
		t.Fatalf("unexpected output: %s", output)
	}
}

func TestListCompareLocalFixturesRejectsReservedDeviceNamesGracefully(t *testing.T) {
	pythonBin := mustLookPath(t, "python3")
	manifestDir := t.TempDir()
	manifestPath := filepath.Join(manifestDir, "compare-local-fixtures.txt")
	data := "con testdata/workloads/mock-stage-aware.xml\n"
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
	if !strings.Contains(output, "invalid compare-local fixture name 'con'") {
		t.Fatalf("unexpected output: %s", output)
	}
	if !strings.Contains(output, "reserved device name") {
		t.Fatalf("unexpected output: %s", output)
	}
}

func TestListCompareLocalFixturesRejectsReservedDeviceNamesWithExtensionsGracefully(t *testing.T) {
	pythonBin := mustLookPath(t, "python3")
	manifestDir := t.TempDir()
	manifestPath := filepath.Join(manifestDir, "compare-local-fixtures.txt")
	data := "con.txt testdata/workloads/mock-stage-aware.xml\n"
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
	if !strings.Contains(output, "invalid compare-local fixture name 'con.txt'") {
		t.Fatalf("unexpected output: %s", output)
	}
	if !strings.Contains(output, "reserved device name") {
		t.Fatalf("unexpected output: %s", output)
	}
}

func TestListCompareLocalFixturesRejectsUnsafeWorkloadPathsGracefully(t *testing.T) {
	tests := []struct {
		name        string
		workload    string
		expectError string
	}{
		{
			name:        "parent traversal",
			workload:    "../outside.xml",
			expectError: "must be repo-relative without '..' segments",
		},
		{
			name:        "absolute path",
			workload:    "/tmp/outside.xml",
			expectError: "must not be absolute",
		},
		{
			name:        "windows drive absolute",
			workload:    "C:/outside.xml",
			expectError: "must not be absolute",
		},
		{
			name:        "windows drive relative",
			workload:    "C:outside.xml",
			expectError: "must not include a Windows drive prefix",
		},
		{
			name:        "windows backslash traversal",
			workload:    "..\\outside.xml",
			expectError: "must use forward slashes instead of backslashes",
		},
		{
			name:        "dash prefixed",
			workload:    "-bad.xml",
			expectError: "must not start with -",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			pythonBin := mustLookPath(t, "python3")
			manifestDir := t.TempDir()
			manifestPath := filepath.Join(manifestDir, "compare-local-fixtures.txt")
			data := "fixture " + tc.workload + "\n"
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
			if !strings.Contains(output, "invalid compare-local workload path") {
				t.Fatalf("unexpected output: %s", output)
			}
			if !strings.Contains(output, tc.expectError) {
				t.Fatalf("unexpected output: %s", output)
			}
		})
	}
}

func TestListCompareLocalFixturesRejectsNonXMLWorkloadPathsGracefully(t *testing.T) {
	pythonBin := mustLookPath(t, "python3")
	manifestDir := t.TempDir()
	manifestPath := filepath.Join(manifestDir, "compare-local-fixtures.txt")
	data := "fixture README.md\n"
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
	if !strings.Contains(output, "invalid compare-local workload path 'README.md'") {
		t.Fatalf("unexpected output: %s", output)
	}
	if !strings.Contains(output, "must end with .xml") {
		t.Fatalf("unexpected output: %s", output)
	}
}

func TestListCompareLocalFixturesAcceptsUppercaseXMLWorkloadPaths(t *testing.T) {
	pythonBin := mustLookPath(t, "python3")
	manifestDir := t.TempDir()
	manifestPath := filepath.Join(manifestDir, "compare-local-fixtures.txt")
	data := "fixture fixture/WORKLOAD.XML\n"
	if err := os.WriteFile(manifestPath, []byte(data), 0o644); err != nil {
		t.Fatalf("write manifest: %v", err)
	}

	scriptPath, err := filepath.Abs(filepath.Clean("../../scripts/list_compare_local_fixtures.py"))
	if err != nil {
		t.Fatalf("abs script path: %v", err)
	}
	cmd := exec.Command(pythonBin, scriptPath, manifestPath)
	cmd.Dir = repoRootDir()
	output := string(runCommandSuccess(t, cmd))

	if !strings.Contains(output, "\"workload\": \"fixture/WORKLOAD.XML\"") {
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

func TestBuildCompareLocalIndexCreatesNonASCIIOutputDirGracefully(t *testing.T) {
	pythonBin := mustLookPath(t, "python3")
	manifestDir := t.TempDir()
	manifestPath := filepath.Join(manifestDir, "compare-local-fixtures.txt")
	outputDir := filepath.Join(manifestDir, "输出")
	if err := os.WriteFile(manifestPath, []byte("# comment only\n"), 0o644); err != nil {
		t.Fatalf("write manifest: %v", err)
	}

	scriptPath, err := filepath.Abs(filepath.Clean("../../scripts/build_compare_local_index.py"))
	if err != nil {
		t.Fatalf("abs script path: %v", err)
	}
	cmd := exec.Command(pythonBin, scriptPath, manifestPath, outputDir)
	cmd.Dir = repoRootDir()
	cmd.Env = append(
		os.Environ(),
		"LC_ALL=C",
		"LANG=C",
		"PYTHONCOERCECLOCALE=0",
		"PYTHONUTF8=0",
	)
	output := string(runCommandSuccess(t, cmd))

	if strings.TrimSpace(output) != "" {
		t.Fatalf("unexpected stdout: %s", output)
	}
	summaryData := mustReadFile(t, filepath.Join(outputDir, "summary.md"))
	if !strings.Contains(string(summaryData), "Artifact directory: `"+outputDir+"`") {
		t.Fatalf("unexpected summary: %s", summaryData)
	}
}

func TestBuildCompareLocalIndexRejectsFileOutputDirGracefully(t *testing.T) {
	pythonBin := mustLookPath(t, "python3")
	manifestDir := t.TempDir()
	manifestPath := filepath.Join(manifestDir, "compare-local-fixtures.txt")
	outputPath := filepath.Join(manifestDir, "out")
	if err := os.WriteFile(manifestPath, []byte("# comment only\n"), 0o644); err != nil {
		t.Fatalf("write manifest: %v", err)
	}
	if err := os.WriteFile(outputPath, []byte("file\n"), 0o644); err != nil {
		t.Fatalf("write output file: %v", err)
	}

	scriptPath, err := filepath.Abs(filepath.Clean("../../scripts/build_compare_local_index.py"))
	if err != nil {
		t.Fatalf("abs script path: %v", err)
	}
	cmd := exec.Command(pythonBin, scriptPath, manifestPath, outputPath)
	cmd.Dir = repoRootDir()
	output := string(runCommandFailure(t, cmd))

	if strings.Contains(output, "Traceback") {
		t.Fatalf("unexpected traceback: %s", output)
	}
	if !strings.Contains(output, "unable to prepare compare-local output dir") {
		t.Fatalf("unexpected output: %s", output)
	}
	if !strings.Contains(output, outputPath) {
		t.Fatalf("unexpected output: %s", output)
	}
}

func TestBuildCompareLocalIndexRejectsNonASCIIFileOutputDirGracefully(t *testing.T) {
	pythonBin := mustLookPath(t, "python3")
	manifestDir := t.TempDir()
	manifestPath := filepath.Join(manifestDir, "compare-local-fixtures.txt")
	outputPath := filepath.Join(manifestDir, "输出")
	if err := os.WriteFile(manifestPath, []byte("# comment only\n"), 0o644); err != nil {
		t.Fatalf("write manifest: %v", err)
	}
	if err := os.WriteFile(outputPath, []byte("file\n"), 0o644); err != nil {
		t.Fatalf("write output file: %v", err)
	}

	scriptPath, err := filepath.Abs(filepath.Clean("../../scripts/build_compare_local_index.py"))
	if err != nil {
		t.Fatalf("abs script path: %v", err)
	}
	cmd := exec.Command(pythonBin, scriptPath, manifestPath, outputPath)
	cmd.Dir = repoRootDir()
	cmd.Env = append(
		os.Environ(),
		"LC_ALL=C",
		"LANG=C",
		"PYTHONCOERCECLOCALE=0",
		"PYTHONUTF8=0",
	)
	output := string(runCommandFailure(t, cmd))

	if strings.Contains(output, "Traceback") {
		t.Fatalf("unexpected traceback: %s", output)
	}
	if strings.Contains(output, "\\udc") {
		t.Fatalf("unexpected surrogate escapes: %s", output)
	}
	if !strings.Contains(output, "unable to prepare compare-local output dir") {
		t.Fatalf("unexpected output: %s", output)
	}
	if !strings.Contains(output, "输出") {
		t.Fatalf("unexpected output: %s", output)
	}
}

func TestBuildCompareLocalIndexRejectsUnwritableIndexOutputGracefully(t *testing.T) {
	pythonBin := mustLookPath(t, "python3")
	manifestDir := t.TempDir()
	manifestPath := filepath.Join(manifestDir, "compare-local-fixtures.txt")
	outputDir := filepath.Join(manifestDir, "out")
	if err := os.WriteFile(manifestPath, []byte("# comment only\n"), 0o644); err != nil {
		t.Fatalf("write manifest: %v", err)
	}
	if err := os.MkdirAll(filepath.Join(outputDir, "index.json"), 0o755); err != nil {
		t.Fatalf("mkdir index dir: %v", err)
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
	if !strings.Contains(output, "unable to write compare-local artifact") {
		t.Fatalf("unexpected output: %s", output)
	}
	if !strings.Contains(output, filepath.Join(outputDir, "index.json")) {
		t.Fatalf("unexpected output: %s", output)
	}
}

func TestBuildCompareLocalIndexRejectsNonASCIIUnwritableArtifactPathGracefully(t *testing.T) {
	pythonBin := mustLookPath(t, "python3")
	manifestDir := t.TempDir()
	manifestPath := filepath.Join(manifestDir, "compare-local-fixtures.txt")
	outputDir := filepath.Join(manifestDir, "输出")
	if err := os.WriteFile(manifestPath, []byte("# comment only\n"), 0o644); err != nil {
		t.Fatalf("write manifest: %v", err)
	}
	if err := os.MkdirAll(filepath.Join(outputDir, "index.json"), 0o755); err != nil {
		t.Fatalf("mkdir index dir: %v", err)
	}

	scriptPath, err := filepath.Abs(filepath.Clean("../../scripts/build_compare_local_index.py"))
	if err != nil {
		t.Fatalf("abs script path: %v", err)
	}
	cmd := exec.Command(pythonBin, scriptPath, manifestPath, outputDir)
	cmd.Dir = repoRootDir()
	cmd.Env = append(
		os.Environ(),
		"LC_ALL=C",
		"LANG=C",
		"PYTHONCOERCECLOCALE=0",
		"PYTHONUTF8=0",
	)
	output := string(runCommandFailure(t, cmd))

	if strings.Contains(output, "Traceback") {
		t.Fatalf("unexpected traceback: %s", output)
	}
	if strings.Contains(output, "\\udc") {
		t.Fatalf("unexpected surrogate escapes: %s", output)
	}
	if !strings.Contains(output, "unable to write compare-local artifact") {
		t.Fatalf("unexpected output: %s", output)
	}
	if !strings.Contains(output, "输出") {
		t.Fatalf("unexpected output: %s", output)
	}
	if !strings.Contains(output, "index.json") {
		t.Fatalf("unexpected output: %s", output)
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

func TestBuildCompareLocalIndexRejectsNonObjectFixtureSummaryGracefully(t *testing.T) {
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
	if err := os.WriteFile(filepath.Join(outputDir, "mock-stage-aware.json"), []byte("123\n"), 0o644); err != nil {
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
	if !strings.Contains(output, "summary payload must be a JSON object") {
		t.Fatalf("unexpected output: %s", output)
	}
	if !strings.Contains(output, filepath.Join(outputDir, "mock-stage-aware.json")) {
		t.Fatalf("unexpected output: %s", output)
	}
}

func TestBuildCompareLocalIndexRejectsWrongTypedFixtureSummaryFieldsGracefully(t *testing.T) {
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
	if err := os.WriteFile(filepath.Join(outputDir, "mock-stage-aware.json"), []byte("{\"stages\":\"1\",\"works\":1,\"samples\":1,\"errors\":0}\n"), 0o644); err != nil {
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
	if !strings.Contains(output, "field stages must be an integer") {
		t.Fatalf("unexpected output: %s", output)
	}
	if !strings.Contains(output, filepath.Join(outputDir, "mock-stage-aware.json")) {
		t.Fatalf("unexpected output: %s", output)
	}
}

func TestBuildCompareLocalIndexRejectsNegativeFixtureSummaryFieldsGracefully(t *testing.T) {
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
	if err := os.WriteFile(filepath.Join(outputDir, "mock-stage-aware.json"), []byte("{\"stages\":1,\"works\":1,\"samples\":1,\"errors\":-1}\n"), 0o644); err != nil {
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
	if !strings.Contains(output, "field errors must be a non-negative integer") {
		t.Fatalf("unexpected output: %s", output)
	}
	if !strings.Contains(output, filepath.Join(outputDir, "mock-stage-aware.json")) {
		t.Fatalf("unexpected output: %s", output)
	}
}
