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
