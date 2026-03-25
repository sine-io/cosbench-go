package domain

import (
	legacyworkload "github.com/sine-io/cosbench-go/internal/domain/workload"
)

type Workload struct {
	Name        string        `json:"name"`
	Description string        `json:"description,omitempty"`
	Trigger     string        `json:"trigger,omitempty"`
	Config      string        `json:"config,omitempty"`
	Storage     *StorageSpec  `json:"storage,omitempty"`
	Workflow    Workflow      `json:"workflow"`
}

type Workflow struct {
	Config string  `json:"config,omitempty"`
	Stages []Stage `json:"stages"`
}

type Stage struct {
	Name         string        `json:"name"`
	ClosureDelay int           `json:"closure_delay,omitempty"`
	Trigger      string        `json:"trigger,omitempty"`
	Config       string        `json:"config,omitempty"`
	Storage      *StorageSpec  `json:"storage,omitempty"`
	Works        []Work        `json:"works"`
}

type Work struct {
	Name       string       `json:"name"`
	Type       string       `json:"type"`
	Workers    int          `json:"workers"`
	Interval   int          `json:"interval,omitempty"`
	Division   string       `json:"division,omitempty"`
	Runtime    int          `json:"runtime,omitempty"`
	RampUp     int          `json:"ramp_up,omitempty"`
	RampDown   int          `json:"ramp_down,omitempty"`
	AFR        int          `json:"afr,omitempty"`
	TotalOps   int          `json:"total_ops,omitempty"`
	TotalBytes int64        `json:"total_bytes,omitempty"`
	Driver     string       `json:"driver,omitempty"`
	Config     string       `json:"config,omitempty"`
	Storage    *StorageSpec `json:"storage,omitempty"`
	Operations []Operation  `json:"operations"`
}

type Operation struct {
	Type     string `json:"type"`
	Ratio    int    `json:"ratio"`
	Division string `json:"division,omitempty"`
	Config   string `json:"config,omitempty"`
	ID       string `json:"id,omitempty"`
}

type StorageSpec struct {
	Type   string `json:"type"`
	Config string `json:"config,omitempty"`
}

func WorkloadFromLegacy(src legacyworkload.Workload) Workload {
	dst := Workload{
		Name:        src.Name,
		Description: src.Description,
		Trigger:     src.Trigger,
		Config:      src.Config,
		Storage:     storageFromLegacy(src.Storage),
		Workflow: Workflow{
			Config: src.Workflow.Config,
			Stages: make([]Stage, 0, len(src.Workflow.Stages)),
		},
	}
	for _, stage := range src.Workflow.Stages {
		dStage := Stage{
			Name:         stage.Name,
			ClosureDelay: stage.ClosureDelay,
			Trigger:      stage.Trigger,
			Config:       stage.Config,
			Storage:      storageFromLegacy(stage.Storage),
			Works:        make([]Work, 0, len(stage.Works)),
		}
		for _, work := range stage.Works {
			dWork := Work{
				Name:       work.Name,
				Type:       work.Type,
				Workers:    work.Workers,
				Interval:   work.Interval,
				Division:   work.Division,
				Runtime:    work.Runtime,
				RampUp:     work.RampUp,
				RampDown:   work.RampDown,
				AFR:        work.AFR,
				TotalOps:   work.TotalOps,
				TotalBytes: work.TotalBytes,
				Driver:     work.Driver,
				Config:     work.Config,
				Storage:    storageFromLegacy(work.Storage),
				Operations: make([]Operation, 0, len(work.Operations)),
			}
			for _, op := range work.Operations {
				dWork.Operations = append(dWork.Operations, Operation{
					Type:     op.Type,
					Ratio:    op.Ratio,
					Division: op.Division,
					Config:   op.Config,
					ID:       op.ID,
				})
			}
			dStage.Works = append(dStage.Works, dWork)
		}
		dst.Workflow.Stages = append(dst.Workflow.Stages, dStage)
	}
	return dst
}

func (w Workload) ToLegacy() legacyworkload.Workload {
	dst := legacyworkload.Workload{
		Name:        w.Name,
		Description: w.Description,
		Trigger:     w.Trigger,
		Config:      w.Config,
		Storage:     storageToLegacy(w.Storage),
		Workflow: legacyworkload.Workflow{
			Config: w.Workflow.Config,
			Stages: make([]legacyworkload.Stage, 0, len(w.Workflow.Stages)),
		},
	}
	for _, stage := range w.Workflow.Stages {
		lStage := legacyworkload.Stage{
			Name:         stage.Name,
			ClosureDelay: stage.ClosureDelay,
			Trigger:      stage.Trigger,
			Config:       stage.Config,
			Storage:      storageToLegacy(stage.Storage),
			Works:        make([]legacyworkload.Work, 0, len(stage.Works)),
		}
		for _, work := range stage.Works {
			lWork := legacyworkload.Work{
				Name:       work.Name,
				Type:       work.Type,
				Workers:    work.Workers,
				Interval:   work.Interval,
				Division:   work.Division,
				Runtime:    work.Runtime,
				RampUp:     work.RampUp,
				RampDown:   work.RampDown,
				AFR:        work.AFR,
				TotalOps:   work.TotalOps,
				TotalBytes: work.TotalBytes,
				Driver:     work.Driver,
				Config:     work.Config,
				Storage:    storageToLegacy(work.Storage),
				Operations: make([]legacyworkload.Operation, 0, len(work.Operations)),
			}
			for _, op := range work.Operations {
				lWork.Operations = append(lWork.Operations, legacyworkload.Operation{
					Type:     op.Type,
					Ratio:    op.Ratio,
					Division: op.Division,
					Config:   op.Config,
					ID:       op.ID,
				})
			}
			lStage.Works = append(lStage.Works, lWork)
		}
		dst.Workflow.Stages = append(dst.Workflow.Stages, lStage)
	}
	return dst
}

func (w Work) ToLegacy() legacyworkload.Work {
	legacy := legacyworkload.Work{
		Name:       w.Name,
		Type:       w.Type,
		Workers:    w.Workers,
		Interval:   w.Interval,
		Division:   w.Division,
		Runtime:    w.Runtime,
		RampUp:     w.RampUp,
		RampDown:   w.RampDown,
		AFR:        w.AFR,
		TotalOps:   w.TotalOps,
		TotalBytes: w.TotalBytes,
		Driver:     w.Driver,
		Config:     w.Config,
		Storage:    storageToLegacy(w.Storage),
		Operations: make([]legacyworkload.Operation, 0, len(w.Operations)),
	}
	for _, op := range w.Operations {
		legacy.Operations = append(legacy.Operations, legacyworkload.Operation{
			Type:     op.Type,
			Ratio:    op.Ratio,
			Division: op.Division,
			Config:   op.Config,
			ID:       op.ID,
		})
	}
	return legacy
}

func storageFromLegacy(src *legacyworkload.StorageSpec) *StorageSpec {
	if src == nil {
		return nil
	}
	return &StorageSpec{Type: src.Type, Config: src.Config}
}

func storageToLegacy(src *StorageSpec) *legacyworkload.StorageSpec {
	if src == nil {
		return nil
	}
	return &legacyworkload.StorageSpec{Type: src.Type, Config: src.Config}
}
