package main

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

func TestCompareLocalPrunesStaleOutputs(t *testing.T) {
	makeBin, err := exec.LookPath("make")
	if err != nil {
		t.Fatalf("look path make: %v", err)
	}
	goBin, err := exec.LookPath("go")
	if err != nil {
		t.Fatalf("look path go: %v", err)
	}

	rootDir := filepath.Clean("../..")
	outputDir := filepath.Join(t.TempDir(), "compare-local")
	if err := os.MkdirAll(outputDir, 0o755); err != nil {
		t.Fatalf("mkdir output dir: %v", err)
	}
	staleFile := filepath.Join(outputDir, "stale.json")
	if err := os.WriteFile(staleFile, []byte("{}\n"), 0o644); err != nil {
		t.Fatalf("seed stale file: %v", err)
	}

	cmd := exec.Command(makeBin, "compare-local", "GO="+goBin, "COMPARE_LOCAL_OUTPUT_DIR="+outputDir)
	cmd.Dir = rootDir
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("make compare-local failed: %v\n%s", err, output)
	}

	if _, err := os.Stat(staleFile); err == nil {
		t.Fatalf("expected stale output %s to be removed", staleFile)
	}

	for _, name := range []string{
		"s3-active-subset.json",
		"mock-stage-aware.json",
		"mock-reusedata-subset.json",
		"xml-splitrw-subset.json",
	} {
		if _, err := os.Stat(filepath.Join(outputDir, name)); err != nil {
			t.Fatalf("expected fresh output %s: %v", name, err)
		}
	}
}

func TestCompareLocalRejectsUnsafeOutputDir(t *testing.T) {
	makeBin, err := exec.LookPath("make")
	if err != nil {
		t.Fatalf("look path make: %v", err)
	}
	goBin, err := exec.LookPath("go")
	if err != nil {
		t.Fatalf("look path go: %v", err)
	}

	rootDir := filepath.Clean("../..")
	unsafeDir := t.TempDir()
	cmd := exec.Command(makeBin, "compare-local", "GO="+goBin, "COMPARE_LOCAL_OUTPUT_DIR="+unsafeDir)
	cmd.Dir = rootDir
	output, err := cmd.CombinedOutput()
	if err == nil {
		t.Fatalf("expected compare-local to reject unsafe output dir %s\n%s", unsafeDir, output)
	}
}
