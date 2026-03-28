package web

import (
	"net/http"
	"strings"

	"github.com/sine-io/cosbench-go/internal/domain"
)

func (h *Handler) controllerMatrixPage(w http.ResponseWriter, r *http.Request) {
	h.render(w, "controller_matrix.html", pageData{
		Title:      "Controller Matrix",
		MatrixRows: h.manager.ListJobMatrix(),
		RequestPath: r.URL.Path,
	})
}

func (h *Handler) controllerJobPageRoute(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, "/controller/jobs/")
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
	timeline, _ := h.manager.GetJobTimeline(jobID)

	if len(parts) == 2 && parts[1] == "config" {
		h.render(w, "controller_job_config.html", pageData{
			Title:       "Job Config",
			Job:         job,
			RawXML:      job.RawXML,
			RequestPath: r.URL.Path,
		})
		return
	}
	if len(parts) == 3 && parts[1] == "config" && parts[2] == "advanced" {
		h.render(w, "controller_advanced_config.html", pageData{
			Title:              "Advanced Config",
			Job:                job,
			NormalizedWorkload: job.Workload,
			RequestPath:        r.URL.Path,
		})
		return
	}
	if len(parts) == 3 && parts[1] == "stages" {
		stageName := parts[2]
		stage := findStageState(job, result, stageName)
		if stage.Name == "" {
			http.NotFound(w, r)
			return
		}
		h.render(w, "controller_stage.html", pageData{
			Title:       "Stage Detail",
			Job:         job,
			Stage:       stage,
			RequestPath: r.URL.Path,
		})
		return
	}
	if len(parts) == 2 && parts[1] == "timeline" {
		h.render(w, "controller_timeline.html", pageData{
			Title:       "Job Timeline",
			Job:         job,
			Timeline:    timeline,
			RequestPath: r.URL.Path,
		})
		return
	}

	http.NotFound(w, r)
}

func findStageState(job domain.Job, result domain.JobResult, name string) domain.StageState {
	for _, stage := range result.StageTotals {
		if stage.Name == name {
			return stage
		}
	}
	for _, stage := range job.Stages {
		if stage.Name == name {
			return stage
		}
	}
	return domain.StageState{}
}
