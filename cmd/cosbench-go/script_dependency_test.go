package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

var scriptTestDependencies = []string{
	"Makefile",
	"scripts/build_compare_local_index.py",
	"scripts/compare_local_manifest.py",
	"scripts/list_compare_local_fixtures.py",
	"scripts/smoke_local.py",
	"scripts/smoke_ready.py",
	"scripts/validate_compare_local_filter.py",
	"scripts/worktree_audit.py",
	"scripts/worktree_cleanup_report.py",
	"scripts/worktree_output.py",
	"scripts/worktree_prune_plan.py",
	"testdata/workloads/compare-local-fixtures.txt",
	"testdata/workloads/mock-stage-aware.xml",
}

func TestMain(m *testing.M) {
	rootDir := filepath.Clean("../..")
	dependencies := append([]string{}, scriptTestDependencies...)
	manifestPath := filepath.Join(rootDir, "testdata/workloads/compare-local-fixtures.txt")
	manifestData, err := os.ReadFile(manifestPath)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "read compare-local manifest %s: %v\n", manifestPath, err)
		os.Exit(1)
	}
	for _, rawLine := range strings.Split(string(manifestData), "\n") {
		line := strings.TrimSpace(rawLine)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		fields := strings.Fields(line)
		if len(fields) != 2 {
			_, _ = fmt.Fprintf(os.Stderr, "unexpected compare-local manifest line %q\n", line)
			os.Exit(1)
		}
		dependencies = append(dependencies, fields[1])
	}

	for _, rel := range dependencies {
		path := filepath.Join(rootDir, rel)
		data, err := os.ReadFile(path)
		if err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "read script test dependency %s: %v\n", path, err)
			os.Exit(1)
		}
		if len(data) == 0 {
			_, _ = fmt.Fprintf(os.Stderr, "script test dependency is empty: %s\n", path)
			os.Exit(1)
		}
	}
	os.Exit(m.Run())
}
