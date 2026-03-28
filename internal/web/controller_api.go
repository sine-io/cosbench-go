package web

import (
	"encoding/csv"
	"encoding/json"
	"net/http"
	"sort"
	"strings"

	"github.com/sine-io/cosbench-go/internal/domain"
)

func (h *Handler) controllerJobsAPI(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	writeJSON(w, http.StatusOK, h.manager.ListJobMatrix())
}

func (h *Handler) controllerJobAPIRoute(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	path := strings.TrimPrefix(r.URL.Path, "/api/controller/jobs/")
	path = strings.Trim(path, "/")
	if path == "" {
		http.NotFound(w, r)
		return
	}
	parts := strings.Split(path, "/")
	jobID := parts[0]
	job, ok := h.manager.GetJob(jobID)
	if !ok {
		http.NotFound(w, r)
		return
	}
	result, _ := h.manager.GetJobResult(jobID)
	events := h.manager.GetJobEvents(jobID)
	timeline, _ := h.manager.GetJobTimeline(jobID)

	if len(parts) == 1 {
		writeJSON(w, http.StatusOK, struct {
			Job    domain.Job        `json:"job"`
			Result domain.JobResult  `json:"result"`
			Events []domain.JobEvent `json:"events"`
		}{Job: job, Result: result, Events: events})
		return
	}

	switch parts[1] {
	case "config":
		if len(parts) == 2 {
			writeJSON(w, http.StatusOK, struct {
				Job    domain.Job `json:"job"`
				RawXML string     `json:"raw_xml"`
			}{Job: job, RawXML: job.RawXML})
			return
		}
		if len(parts) == 3 && parts[2] == "advanced" {
			writeJSON(w, http.StatusOK, struct {
				Job                domain.Job      `json:"job"`
				NormalizedWorkload domain.Workload `json:"normalized_workload"`
			}{Job: job, NormalizedWorkload: job.Workload})
			return
		}
	case "stages":
		if len(parts) == 3 {
			for _, stage := range result.StageTotals {
				if stage.Name == parts[2] {
					writeJSON(w, http.StatusOK, struct {
						JobID string            `json:"job_id"`
						Stage domain.StageState `json:"stage"`
					}{JobID: jobID, Stage: stage})
					return
				}
			}
			for _, stage := range job.Stages {
				if stage.Name == parts[2] {
					writeJSON(w, http.StatusOK, struct {
						JobID string            `json:"job_id"`
						Stage domain.StageState `json:"stage"`
					}{JobID: jobID, Stage: stage})
					return
				}
			}
			http.NotFound(w, r)
			return
		}
	case "timeline":
		if len(parts) == 2 {
			writeJSON(w, http.StatusOK, timeline)
			return
		}
	case "artifacts":
		if len(parts) == 3 {
			h.exportControllerArtifact(w, r, jobID, parts[2])
			return
		}
	}

	if len(parts) == 2 && parts[1] == "timeline.csv" {
		h.exportControllerTimelineCSV(w, timeline)
		return
	}

	http.NotFound(w, r)
}

func (h *Handler) exportControllerTimelineCSV(w http.ResponseWriter, timeline domain.JobTimeline) {
	w.Header().Set("Content-Type", "text/csv; charset=utf-8")
	writer := csv.NewWriter(w)
	_ = writer.Write([]string{"scope", "name", "timestamp", "operation_count", "byte_count", "error_count", "avg_latency_ms"})
	for _, point := range timeline.Job {
		_ = writer.Write([]string{"job", "summary", point.Timestamp.Format(timeLayout), itoa(point.OperationCount), itoa(point.ByteCount), itoa(point.ErrorCount), formatFloat(point.AvgLatencyMs)})
	}
	stageNames := make([]string, 0, len(timeline.Stages))
	for name := range timeline.Stages {
		stageNames = append(stageNames, name)
	}
	sort.Strings(stageNames)
	for _, name := range stageNames {
		for _, point := range timeline.Stages[name] {
			_ = writer.Write([]string{"stage", name, point.Timestamp.Format(timeLayout), itoa(point.OperationCount), itoa(point.ByteCount), itoa(point.ErrorCount), formatFloat(point.AvgLatencyMs)})
		}
	}
	writer.Flush()
}

func writeJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	_ = enc.Encode(payload)
}

const timeLayout = "2006-01-02T15:04:05.000Z07:00"
