package main

import (
	"bufio"
	"os"
	"path/filepath"
	"strings"
	"testing"

	xmlparser "github.com/sine-io/cosbench-go/internal/infrastructure/xml"
)

type compareLocalFixture struct {
	Name string
	Path string
}

func TestCompareLocalManifestFixturesParse(t *testing.T) {
	fixtures := readCompareLocalFixtures(t)
	if len(fixtures) == 0 {
		t.Fatal("expected at least one compare-local fixture")
	}

	for _, fixture := range fixtures {
		workloadPath := filepath.Clean(filepath.Join("../..", fixture.Path))
		if _, err := os.Stat(workloadPath); err != nil {
			t.Fatalf("stat %s: %v", workloadPath, err)
		}
		if _, err := xmlparser.ParseWorkloadFile(workloadPath); err != nil {
			t.Fatalf("parse %s: %v", workloadPath, err)
		}
	}
}

func readCompareLocalFixtures(t *testing.T) []compareLocalFixture {
	t.Helper()

	manifestPath := filepath.Clean("../../testdata/workloads/compare-local-fixtures.txt")
	file, err := os.Open(manifestPath)
	if err != nil {
		t.Fatalf("open manifest: %v", err)
	}
	defer file.Close()

	var fixtures []compareLocalFixture
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		fields := strings.Fields(line)
		if len(fields) != 2 {
			t.Fatalf("invalid manifest line %q", line)
		}
		fixtures = append(fixtures, compareLocalFixture{
			Name: fields[0],
			Path: fields[1],
		})
	}
	if err := scanner.Err(); err != nil {
		t.Fatalf("scan manifest: %v", err)
	}
	return fixtures
}
