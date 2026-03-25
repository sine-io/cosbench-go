package main

import (
	"context"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/sine-io/cosbench-go/internal/domain/execution"
	storagefactory "github.com/sine-io/cosbench-go/internal/infrastructure/storage"
	xmlparser "github.com/sine-io/cosbench-go/internal/infrastructure/xml"
)

type cliSummary struct {
	Workload string `json:"workload"`
	Stages   int    `json:"stages"`
	Works    int    `json:"works"`
	Samples  int    `json:"samples"`
	Errors   int64  `json:"errors"`
}

func main() {
	workloadPath := flag.String("workload", "", "path to workload xml")
	shortWorkloadPath := flag.String("f", "", "path to workload xml")
	backend := flag.String("backend", "", "override storage backend: mock|s3|sio")
	jsonOut := flag.Bool("json", false, "print JSON summary")
	flag.Parse()
	path, err := resolveWorkloadPath(*workloadPath, *shortWorkloadPath, flag.Args())
	if err != nil {
		fmt.Fprintln(os.Stderr, "usage: cosbench-go [-workload <path> | -f <path> | <path>] [-backend mock|s3|sio] [-json]")
		os.Exit(2)
	}
	if err := runCLI(path, *backend, *jsonOut, os.Stdout, os.Stderr); err != nil {
		os.Exit(1)
	}
}

func runCLI(workloadPath, backend string, jsonOut bool, stdout, stderr io.Writer) error {
	wl, err := xmlparser.ParseWorkloadFile(workloadPath)
	if err != nil {
		fmt.Fprintf(stderr, "parse workload: %v\n", err)
		return err
	}

	ctx := context.Background()
	totalWorks := 0
	totalSamples := 0
	var totalErrors int64
	adapters := storagefactory.NewRunAdapters()
	defer adapters.Close()
	progressOut := stdout
	if jsonOut {
		progressOut = stderr
	}

	fmt.Fprintf(progressOut, "workload=%s stages=%d\n", wl.Name, len(wl.Workflow.Stages))
	for _, stage := range wl.Workflow.Stages {
		fmt.Fprintf(progressOut, "stage=%s works=%d\n", stage.Name, len(stage.Works))
		for _, work := range stage.Works {
			totalWorks++
			storageType := ""
			rawConfig := ""
			if work.Storage != nil {
				storageType = work.Storage.Type
				rawConfig = work.Storage.Config
			}
			if backend != "" {
				storageType = backend
			}
			adapter, shared, err := adapters.Acquire(storageType, rawConfig)
			if err != nil {
				fmt.Fprintf(stderr, "storage adapter: %v\n", err)
				return err
			}
			engine := &execution.Engine{Work: work, Storage: adapter}
			res := engine.Run(ctx)
			if !shared {
				_ = adapter.Dispose()
			}
			if res.Err != nil {
				fmt.Fprintf(stderr, "run work %s: %v\n", work.Name, res.Err)
				return res.Err
			}
			var errs int64
			for _, s := range res.Samples {
				errs += s.ErrorCount
			}
			totalSamples += len(res.Samples)
			totalErrors += errs
			fmt.Fprintf(progressOut, "  work=%s type=%s workers=%d ops=%d runtime=%d totalOps=%d samples=%d errors=%d\n", work.Name, work.Type, work.Workers, len(work.Operations), work.Runtime, work.TotalOps, len(res.Samples), errs)
		}
	}

	summary := cliSummary{Workload: wl.Name, Stages: len(wl.Workflow.Stages), Works: totalWorks, Samples: totalSamples, Errors: totalErrors}
	if jsonOut {
		enc := json.NewEncoder(stdout)
		enc.SetIndent("", "  ")
		_ = enc.Encode(summary)
		return nil
	}
	fmt.Fprintf(stdout, "summary: workload=%s works=%d samples=%d errors=%d\n", summary.Workload, summary.Works, summary.Samples, summary.Errors)
	return nil
}

func resolveWorkloadPath(longFlag, shortFlag string, args []string) (string, error) {
	switch {
	case longFlag != "":
		return longFlag, nil
	case shortFlag != "":
		return shortFlag, nil
	case len(args) > 0 && args[0] != "":
		return args[0], nil
	default:
		return "", errors.New("workload path is required")
	}
}
