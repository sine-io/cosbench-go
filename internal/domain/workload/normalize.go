package workload

import (
	"errors"
	"fmt"
	"strings"
)

func NormalizeAndValidate(w Workload) (Workload, error) {
	if strings.TrimSpace(w.Name) == "" {
		return w, errors.New("workload name cannot be empty")
	}
	if len(w.Workflow.Stages) == 0 {
		return w, errors.New("workflow must have stages")
	}

	for si := range w.Workflow.Stages {
		stage := &w.Workflow.Stages[si]
		if strings.TrimSpace(stage.Name) == "" {
			return w, errors.New("stage name cannot be empty")
		}
		stage.Config = inheritConfig(stage.Config, w.Workflow.Config)
		if stage.Storage == nil {
			stage.Storage = cloneStorage(w.Storage)
		}
		if len(stage.Works) == 0 {
			return w, errors.New("stage must have works")
		}

		for wi := range stage.Works {
			work := &stage.Works[wi]
			applyWorkDefaults(work)
			work.Config = inheritConfig(work.Config, stage.Config)
			if work.Storage == nil {
				work.Storage = cloneStorage(stage.Storage)
			}
			if err := normalizeSpecialWork(work); err != nil {
				return w, err
			}
			if err := validateWork(work); err != nil {
				return w, fmt.Errorf("stage %q work %q: %w", stage.Name, work.Name, err)
			}
		}
	}

	return w, nil
}

func applyWorkDefaults(w *Work) {
	if w.Type == "" {
		w.Type = "normal"
	}
	if w.Interval == 0 {
		w.Interval = 5
	}
	if w.Division == "" {
		w.Division = "none"
	}
	if w.AFR == 0 && w.Type == "normal" {
		// keep explicit zero as explicit zero, actual defaulting is handled below only when missing
	}
}

func normalizeSpecialWork(w *Work) error {
	typeName := strings.TrimSpace(w.Type)
	switch typeName {
	case "prepare":
		if w.Name == "" {
			w.Name = "prepare"
		}
		w.Division = "object"
		w.Runtime = 0
		if w.AFR < 0 {
			w.AFR = 0
		}
		w.TotalBytes = 0
		w.TotalOps = w.Workers
		cfg := w.Config
		if !strings.Contains(cfg, "createContainer=") {
			cfg = joinConfig("createContainer=false", cfg)
		}
		w.Operations = []Operation{{Type: "prepare", Ratio: 100, Config: cfg, Division: w.Division}}
	case "mprepare":
		if w.Storage == nil || !isSIOStorage(w.Storage.Type) {
			return errors.New("mprepare requires sio-compatible storage")
		}
		if w.Name == "" {
			w.Name = "mprepare"
		}
		w.Division = "object"
		w.Runtime = 0
		if w.AFR < 0 {
			w.AFR = 0
		}
		w.TotalBytes = 0
		w.TotalOps = w.Workers
		cfg := w.Config
		if !strings.Contains(cfg, "createContainer=") {
			cfg = joinConfig("createContainer=false", cfg)
		}
		w.Operations = []Operation{{Type: "mprepare", Ratio: 100, Config: cfg, Division: w.Division}}
	case "cleanup":
		if w.Name == "" {
			w.Name = "cleanup"
		}
		w.Division = "object"
		w.Runtime = 0
		if w.AFR < 0 {
			w.AFR = 0
		}
		w.TotalBytes = 0
		w.TotalOps = w.Workers
		cfg := w.Config
		if !strings.Contains(cfg, "deleteContainer=") {
			cfg = joinConfig("deleteContainer=false", cfg)
		}
		w.Operations = []Operation{{Type: "cleanup", Ratio: 100, Config: cfg, Division: w.Division}}
	case "init":
		if w.Name == "" {
			w.Name = "init"
		}
		w.Division = "container"
		w.Runtime = 0
		if w.AFR < 0 {
			w.AFR = 0
		}
		w.TotalBytes = 0
		w.TotalOps = w.Workers
		w.Operations = []Operation{{Type: "init", Ratio: 100, Config: joinConfig("objects=r(0,0);sizes=c(0)B", w.Config), Division: w.Division}}
	case "dispose":
		if w.Name == "" {
			w.Name = "dispose"
		}
		w.Division = "container"
		w.Runtime = 0
		if w.AFR < 0 {
			w.AFR = 0
		}
		w.TotalBytes = 0
		w.TotalOps = w.Workers
		w.Operations = []Operation{{Type: "dispose", Ratio: 100, Config: joinConfig("objects=r(0,0);sizes=c(0)B", w.Config), Division: w.Division}}
	case "delay":
		if w.Name == "" {
			w.Name = "delay"
		}
		w.Division = "none"
		w.Runtime = 0
		if w.AFR < 0 {
			w.AFR = 0
		}
		w.TotalBytes = 0
		w.Workers = 1
		w.TotalOps = 1
		w.Operations = []Operation{{Type: "delay", Ratio: 100, Config: "", Division: w.Division}}
	default:
		if w.AFR < 0 {
			w.AFR = 200000
		}
	}
	return nil
}

func validateWork(w *Work) error {
	if strings.TrimSpace(w.Name) == "" {
		return errors.New("work name cannot be empty")
	}
	if w.Workers <= 0 {
		return errors.New("workers must be > 0")
	}
	if w.Runtime == 0 && w.TotalOps == 0 && w.TotalBytes == 0 {
		return errors.New("no work limits detected, either runtime, totalOps or totalBytes")
	}
	if w.TotalOps > 0 && w.Workers > w.TotalOps {
		return errors.New("if totalOps is set, workers should be <= totalOps")
	}
	if w.Storage == nil || strings.TrimSpace(w.Storage.Type) == "" {
		return errors.New("work must have storage")
	}
	if len(w.Operations) == 0 {
		return errors.New("a work must have operations")
	}

	ops := make([]Operation, 0, len(w.Operations))
	sum := 0
	for _, op := range w.Operations {
		op.Config = inheritConfig(op.Config, w.Config)
		if op.Division == "" {
			op.Division = w.Division
		}
		if strings.TrimSpace(op.Type) == "" {
			return errors.New("operation type cannot be empty")
		}
		if strings.TrimSpace(op.Division) == "" {
			return errors.New("operation must have division")
		}
		if op.Ratio < 0 || op.Ratio > 100 {
			return fmt.Errorf("illegal operation ratio: %d", op.Ratio)
		}
		if op.Ratio == 0 {
			continue
		}
		if requiresSIO(op.Type) && !isSIOStorage(w.Storage.Type) {
			return fmt.Errorf("operation %q requires sio-compatible storage", op.Type)
		}
		sum += op.Ratio
		ops = append(ops, op)
	}
	if sum != 100 {
		return fmt.Errorf("op ratio should sum to 100, got %d", sum)
	}
	w.Operations = ops
	return nil
}

func inheritConfig(child, parent string) string {
	child = strings.TrimSpace(child)
	parent = strings.TrimSpace(parent)
	if child == "" {
		return parent
	}
	if parent == "" {
		return child
	}
	return parent + ";" + child
}

func joinConfig(parts ...string) string {
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p != "" {
			out = append(out, p)
		}
	}
	return strings.Join(out, ";")
}

func cloneStorage(s *StorageSpec) *StorageSpec {
	if s == nil {
		return nil
	}
	cp := *s
	return &cp
}

func isSIOStorage(t string) bool {
	switch strings.ToLower(strings.TrimSpace(t)) {
	case "sio", "siov1", "gdas":
		return true
	default:
		return false
	}
}

func requiresSIO(op string) bool {
	switch strings.ToLower(strings.TrimSpace(op)) {
	case "mwrite", "head", "restore", "mfilewrite", "localwrite", "mprepare":
		return true
	default:
		return false
	}
}
