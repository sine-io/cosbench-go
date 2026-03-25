package web

import (
	"context"
	"encoding/csv"
	"encoding/json"
	"html/template"
	"io"
	"net/http"
	"path/filepath"
	"net/url"
	"sort"
	"strconv"
	"strings"

	"github.com/sine-io/cosbench-go/internal/controlplane"
	"github.com/sine-io/cosbench-go/internal/domain"
)

type Handler struct {
	manager   *controlplane.Manager
	templates map[string]*template.Template
	mux       *http.ServeMux
}

type pageData struct {
	Title         string
	Jobs          []domain.Job
	Job           domain.Job
	JobResult     domain.JobResult
	JobEvents     []domain.JobEvent
	RecentErrors  []domain.JobEvent
	Endpoints     []domain.EndpointConfig
	Error         string
	Success       string
	RequestPath   string
	SelectedJobID string
	RawXML        string
}

func NewHandler(manager *controlplane.Manager, viewDir string) (*Handler, error) {
	h := &Handler{manager: manager, mux: http.NewServeMux()}
	if err := h.loadTemplates(viewDir); err != nil {
		return nil, err
	}
	h.routes()
	return h, nil
}

func (h *Handler) loadTemplates(viewDir string) error {
	baseFiles := []string{filepath.Join(viewDir, "layout.html")}
	pages := []string{"dashboard.html", "workload_upload.html", "endpoints.html", "job_detail.html", "history.html"}
	h.templates = map[string]*template.Template{}
	for _, page := range pages {
		pageTemplate, err := template.New(page).ParseFiles(append(baseFiles, filepath.Join(viewDir, page))...)
		if err != nil {
			return err
		}
		h.templates[page] = pageTemplate
	}
	return nil
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.mux.ServeHTTP(w, r)
}

func (h *Handler) routes() {
	h.mux.HandleFunc("/", h.dashboard)
	h.mux.HandleFunc("/workloads/new", h.workloadForm)
	h.mux.HandleFunc("/workloads", h.createWorkload)
	h.mux.HandleFunc("/endpoints", h.endpoints)
	h.mux.HandleFunc("/jobs/", h.jobRoute)
	h.mux.HandleFunc("/exports/jobs/", h.exportRoute)
	h.mux.HandleFunc("/history", h.history)
	h.mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir(filepath.Join("web", "static")))))
}

func (h *Handler) dashboard(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	jobs := h.manager.ListJobs()
	h.render(w, "dashboard.html", pageData{Title: "Dashboard", Jobs: jobs[:min(5, len(jobs))], Endpoints: h.manager.ListEndpoints(), RecentErrors: h.collectRecentErrors(jobs), RequestPath: r.URL.Path})
}

func (h *Handler) workloadForm(w http.ResponseWriter, r *http.Request) {
	h.render(w, "workload_upload.html", pageData{Title: "Upload Workload", Endpoints: h.manager.ListEndpoints(), RequestPath: r.URL.Path, Success: r.URL.Query().Get("success"), Error: r.URL.Query().Get("error")})
}

func (h *Handler) createWorkload(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	if err := r.ParseMultipartForm(4 << 20); err != nil {
		h.redirectError(w, r, "/workloads/new", err)
		return
	}
	endpointID := strings.TrimSpace(r.FormValue("endpoint_id"))
	file, _, err := r.FormFile("workload")
	if err != nil {
		h.redirectError(w, r, "/workloads/new", err)
		return
	}
	defer file.Close()
	raw, err := io.ReadAll(file)
	if err != nil {
		h.redirectError(w, r, "/workloads/new", err)
		return
	}
	job, err := h.manager.CreateJobFromXML(raw, endpointID)
	if err != nil {
		h.redirectError(w, r, "/workloads/new", err)
		return
	}
	http.Redirect(w, r, "/jobs/"+job.ID, http.StatusSeeOther)
}

func (h *Handler) endpoints(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		endpoint := domain.EndpointConfig{
			Name:        strings.TrimSpace(r.FormValue("name")),
			Type:        domain.EndpointType(strings.TrimSpace(r.FormValue("type"))),
			Endpoint:    strings.TrimSpace(r.FormValue("endpoint")),
			Region:      strings.TrimSpace(r.FormValue("region")),
			AccessKey:   strings.TrimSpace(r.FormValue("access_key")),
			SecretKey:   strings.TrimSpace(r.FormValue("secret_key")),
			PathStyle:   r.FormValue("path_style") == "on",
			ExtraConfig: strings.TrimSpace(r.FormValue("extra_config")),
		}
		if _, err := h.manager.CreateEndpoint(endpoint); err != nil {
			h.render(w, "endpoints.html", pageData{Title: "Endpoints", Endpoints: h.manager.ListEndpoints(), Error: err.Error(), RequestPath: r.URL.Path})
			return
		}
		http.Redirect(w, r, "/endpoints?success=1", http.StatusSeeOther)
		return
	}
	h.render(w, "endpoints.html", pageData{Title: "Endpoints", Endpoints: h.manager.ListEndpoints(), RequestPath: r.URL.Path, Success: r.URL.Query().Get("success")})
}

func (h *Handler) jobRoute(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, "/jobs/")
	if strings.HasSuffix(path, "/start") {
		jobID := strings.TrimSuffix(path, "/start")
		jobID = strings.TrimSuffix(jobID, "/")
		if err := h.manager.StartJob(context.Background(), jobID); err != nil {
			h.redirectError(w, r, "/jobs/"+jobID, err)
			return
		}
		http.Redirect(w, r, "/jobs/"+jobID, http.StatusSeeOther)
		return
	}
	if strings.HasSuffix(path, "/cancel") {
		jobID := strings.TrimSuffix(path, "/cancel")
		jobID = strings.TrimSuffix(jobID, "/")
		if err := h.manager.CancelJob(jobID); err != nil {
			h.redirectError(w, r, "/jobs/"+jobID, err)
			return
		}
		http.Redirect(w, r, "/jobs/"+jobID, http.StatusSeeOther)
		return
	}
	jobID := strings.Trim(path, "/")
	job, ok := h.manager.GetJob(jobID)
	if !ok {
		http.NotFound(w, r)
		return
	}
	result, _ := h.manager.GetJobResult(jobID)
	h.render(w, "job_detail.html", pageData{Title: job.Name, Job: job, JobResult: result, JobEvents: h.manager.GetJobEvents(jobID), RequestPath: r.URL.Path, Error: r.URL.Query().Get("error"), RawXML: job.RawXML})
}

func (h *Handler) history(w http.ResponseWriter, r *http.Request) {
	h.render(w, "history.html", pageData{Title: "History", Jobs: h.manager.ListJobs(), RequestPath: r.URL.Path})
}

func (h *Handler) exportRoute(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, "/exports/jobs/")
	if strings.HasSuffix(path, "/result.json") {
		jobID := strings.TrimSuffix(path, "/result.json")
		h.exportJobResultJSON(w, r, strings.Trim(jobID, "/"))
		return
	}
	if strings.HasSuffix(path, "/result.csv") {
		jobID := strings.TrimSuffix(path, "/result.csv")
		h.exportJobResultCSV(w, r, strings.Trim(jobID, "/"))
		return
	}
	http.NotFound(w, r)
}

func (h *Handler) render(w http.ResponseWriter, name string, data pageData) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	tmpl, ok := h.templates[name]
	if !ok {
		http.Error(w, "template not found", http.StatusInternalServerError)
		return
	}
	if err := tmpl.ExecuteTemplate(w, name, data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (h *Handler) collectRecentErrors(jobs []domain.Job) []domain.JobEvent {
	events := make([]domain.JobEvent, 0)
	for _, job := range jobs {
		for _, event := range h.manager.GetJobEvents(job.ID) {
			if event.Level == domain.EventLevelError {
				events = append(events, event)
			}
		}
	}
	sort.Slice(events, func(i, j int) bool { return events[i].OccurredAt.After(events[j].OccurredAt) })
	if len(events) > 5 {
		return events[:5]
	}
	return events
}

func (h *Handler) exportJobResultJSON(w http.ResponseWriter, r *http.Request, jobID string) {
	job, ok := h.manager.GetJob(jobID)
	if !ok {
		http.NotFound(w, r)
		return
	}
	result, _ := h.manager.GetJobResult(jobID)
	payload := struct {
		Job    domain.Job       `json:"job"`
		Result domain.JobResult `json:"result"`
		Events []domain.JobEvent `json:"events"`
	}{Job: job, Result: result, Events: h.manager.GetJobEvents(jobID)}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("Content-Disposition", "attachment; filename=job-"+jobID+"-result.json")
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	_ = enc.Encode(payload)
}

func (h *Handler) exportJobResultCSV(w http.ResponseWriter, r *http.Request, jobID string) {
	job, ok := h.manager.GetJob(jobID)
	if !ok {
		http.NotFound(w, r)
		return
	}
	result, _ := h.manager.GetJobResult(jobID)
	w.Header().Set("Content-Type", "text/csv; charset=utf-8")
	w.Header().Set("Content-Disposition", "attachment; filename=job-"+jobID+"-result.csv")
	writer := csv.NewWriter(w)
	_ = writer.Write([]string{"job_id", "job_name", "status", "scope", "name", "operation_count", "byte_count", "error_count", "avg_latency_ms", "p50_latency_ms", "p95_latency_ms", "p99_latency_ms", "ops_per_second"})
	_ = writer.Write([]string{job.ID, job.Name, string(job.Status), "job", "summary", itoa(job.Metrics.OperationCount), itoa(job.Metrics.ByteCount), itoa(job.Metrics.ErrorCount), formatFloat(job.Metrics.AvgLatencyMs), formatFloat(job.Metrics.P50LatencyMs), formatFloat(job.Metrics.P95LatencyMs), formatFloat(job.Metrics.P99LatencyMs), formatFloat(job.Metrics.OpsPerSecond)})
	for _, stage := range result.StageTotals {
		_ = writer.Write([]string{job.ID, job.Name, string(stage.Status), "stage", stage.Name, itoa(stage.Metrics.OperationCount), itoa(stage.Metrics.ByteCount), itoa(stage.Metrics.ErrorCount), formatFloat(stage.Metrics.AvgLatencyMs), formatFloat(stage.Metrics.P50LatencyMs), formatFloat(stage.Metrics.P95LatencyMs), formatFloat(stage.Metrics.P99LatencyMs), formatFloat(stage.Metrics.OpsPerSecond)})
		for _, work := range stage.WorkResults {
			_ = writer.Write([]string{job.ID, job.Name, string(stage.Status), "work", work.Name, itoa(work.Metrics.OperationCount), itoa(work.Metrics.ByteCount), itoa(work.Metrics.ErrorCount), formatFloat(work.Metrics.AvgLatencyMs), formatFloat(work.Metrics.P50LatencyMs), formatFloat(work.Metrics.P95LatencyMs), formatFloat(work.Metrics.P99LatencyMs), formatFloat(work.Metrics.OpsPerSecond)})
		}
	}
	for _, op := range result.Metrics.ByOperation {
		_ = writer.Write([]string{job.ID, job.Name, string(job.Status), "operation", op.Operation, itoa(op.OperationCount), itoa(op.ByteCount), itoa(op.ErrorCount), formatFloat(op.AvgLatencyMs), formatFloat(op.P50LatencyMs), formatFloat(op.P95LatencyMs), formatFloat(op.P99LatencyMs), ""})
	}
	writer.Flush()
}

func formatFloat(v float64) string {
	return strconv.FormatFloat(v, 'f', 2, 64)
}

func itoa[T ~int64 | ~int](v T) string {
	return strconv.FormatInt(int64(v), 10)
}

func (h *Handler) redirectError(w http.ResponseWriter, r *http.Request, path string, err error) {
	http.Redirect(w, r, path+"?error="+url.QueryEscape(err.Error()), http.StatusSeeOther)
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
