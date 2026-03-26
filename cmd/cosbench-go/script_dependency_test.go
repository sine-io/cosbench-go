package main

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
)

var scriptTestDependencies = []string{
	"Makefile",
	"scripts/build_compare_local_index.py",
	"scripts/compare_local_manifest.py",
	"scripts/list_compare_local_fixtures.py",
	"scripts/validate_compare_local_filter.py",
	"scripts/worktree_audit.py",
	"scripts/worktree_cleanup_report.py",
	"scripts/worktree_output.py",
	"scripts/worktree_prune_plan.py",
	"testdata/workloads/compare-local-fixtures.txt",
}

func TestMain(m *testing.M) {
	rootDir := filepath.Clean("../..")
	for _, rel := range scriptTestDependencies {
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
