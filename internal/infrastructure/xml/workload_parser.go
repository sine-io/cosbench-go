package xml

import (
	"encoding/xml"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/sine-io/cosbench-go/internal/domain/workload"
)

type xmlWorkload struct {
	XMLName     xml.Name    `xml:"workload"`
	Name        string      `xml:"name,attr"`
	Description string      `xml:"description,attr"`
	Trigger     string      `xml:"trigger,attr"`
	Config      string      `xml:"config,attr"`
	Storage     *xmlStorage `xml:"storage"`
	Workflow    xmlWorkflow `xml:"workflow"`
}

type xmlWorkflow struct {
	Config string     `xml:"config,attr"`
	Stages []xmlStage `xml:"workstage"`
}

type xmlStage struct {
	Name         string      `xml:"name,attr"`
	ClosureDelay int         `xml:"closuredelay,attr"`
	Trigger      string      `xml:"trigger,attr"`
	Config       string      `xml:"config,attr"`
	Storage      *xmlStorage `xml:"storage"`
	Works        []xmlWork   `xml:"work"`
}

type xmlWork struct {
	Name       string         `xml:"name,attr"`
	Type       string         `xml:"type,attr"`
	Workers    int            `xml:"workers,attr"`
	Interval   int            `xml:"interval,attr"`
	Division   string         `xml:"division,attr"`
	Runtime    int            `xml:"runtime,attr"`
	RampUp     int            `xml:"rampup,attr"`
	RampDown   int            `xml:"rampdown,attr"`
	AFR        int            `xml:"afr,attr"`
	TotalOps   int            `xml:"totalOps,attr"`
	TotalBytes int64          `xml:"totalBytes,attr"`
	Driver     string         `xml:"driver,attr"`
	Config     string         `xml:"config,attr"`
	Storage    *xmlStorage    `xml:"storage"`
	Operations []xmlOperation `xml:"operation"`
}

type xmlOperation struct {
	Type     string `xml:"type,attr"`
	Ratio    string `xml:"ratio,attr"`
	Division string `xml:"division,attr"`
	Config   string `xml:"config,attr"`
	ID       string `xml:"id,attr"`
}

type xmlStorage struct {
	Type   string `xml:"type,attr"`
	Config string `xml:"config,attr"`
}

func ParseWorkloadFile(path string) (workload.Workload, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return workload.Workload{}, err
	}
	return ParseWorkload(data)
}

func ParseWorkload(data []byte) (workload.Workload, error) {
	var xw xmlWorkload
	if err := xml.Unmarshal(data, &xw); err != nil {
		return workload.Workload{}, fmt.Errorf("parse workload xml: %w", err)
	}
	w := workload.Workload{
		Name:        xw.Name,
		Description: xw.Description,
		Trigger:     xw.Trigger,
		Config:      xw.Config,
		Storage:     toStorage(xw.Storage),
		Workflow: workload.Workflow{
			Config: xw.Workflow.Config,
			Stages: make([]workload.Stage, 0, len(xw.Workflow.Stages)),
		},
	}
	for _, xs := range xw.Workflow.Stages {
		stage := workload.Stage{
			Name:         xs.Name,
			ClosureDelay: xs.ClosureDelay,
			Trigger:      xs.Trigger,
			Config:       xs.Config,
			Storage:      toStorage(xs.Storage),
			Works:        make([]workload.Work, 0, len(xs.Works)),
		}
		for _, xwork := range xs.Works {
			wk := workload.Work{
				Name:       xwork.Name,
				Type:       xwork.Type,
				Workers:    xwork.Workers,
				Interval:   xwork.Interval,
				Division:   xwork.Division,
				Runtime:    xwork.Runtime,
				RampUp:     xwork.RampUp,
				RampDown:   xwork.RampDown,
				AFR:        normalizeAFR(xwork.AFR),
				TotalOps:   xwork.TotalOps,
				TotalBytes: xwork.TotalBytes,
				Driver:     xwork.Driver,
				Config:     xwork.Config,
				Storage:    toStorage(xwork.Storage),
				Operations: make([]workload.Operation, 0, len(xwork.Operations)),
			}
			for _, xo := range xwork.Operations {
				ratio, err := parseOperationRatio(xo.Ratio, xo.Type)
				if err != nil {
					return workload.Workload{}, fmt.Errorf("operation ratio: %w", err)
				}
				wk.Operations = append(wk.Operations, workload.Operation{
					Type:     xo.Type,
					Ratio:    ratio,
					Division: xo.Division,
					Config:   xo.Config,
					ID:       xo.ID,
				})
			}
			stage.Works = append(stage.Works, wk)
		}
		w.Workflow.Stages = append(w.Workflow.Stages, stage)
	}
	return workload.NormalizeAndValidate(w)
}

func toStorage(s *xmlStorage) *workload.StorageSpec {
	if s == nil {
		return nil
	}
	return &workload.StorageSpec{Type: s.Type, Config: s.Config}
}

func normalizeAFR(v int) int {
	if v == 0 {
		return -1
	}
	return v
}

func parseOperationRatio(raw string, opType string) (int, error) {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		if strings.TrimSpace(opType) != "" {
			return 100, nil
		}
		return 0, nil
	}
	ratio, err := strconv.Atoi(raw)
	if err != nil {
		return 0, fmt.Errorf("invalid ratio %q", raw)
	}
	return ratio, nil
}
