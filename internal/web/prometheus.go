package web

import (
	"fmt"
	"net/http"
	"strings"
)

func (h *Handler) controllerPrometheus(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	w.Header().Set("Content-Type", "text/plain; version=0.0.4; charset=utf-8")
	var lines []string
	lines = append(lines, "# HELP cosbench_job_operation_count Total operations recorded for a job.")
	lines = append(lines, "# TYPE cosbench_job_operation_count gauge")
	for _, row := range h.manager.ListJobMatrix() {
		lines = append(lines, fmt.Sprintf(
			`cosbench_job_operation_count{job_id=%q,job_name=%q,status=%q} %d`,
			row.JobID, row.Name, string(row.Status), row.OperationCount,
		))
	}
	lines = append(lines, "# HELP cosbench_job_error_count Total errors recorded for a job.")
	lines = append(lines, "# TYPE cosbench_job_error_count gauge")
	for _, row := range h.manager.ListJobMatrix() {
		lines = append(lines, fmt.Sprintf(
			`cosbench_job_error_count{job_id=%q,job_name=%q,status=%q} %d`,
			row.JobID, row.Name, string(row.Status), row.ErrorCount,
		))
	}
	_, _ = w.Write([]byte(strings.Join(lines, "\n") + "\n"))
}
