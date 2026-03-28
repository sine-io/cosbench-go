package web

import (
	"fmt"
	"net/http"
)

func (h *Handler) exportControllerArtifact(w http.ResponseWriter, r *http.Request, jobID string, artifact string) {
	job, ok := h.manager.GetJob(jobID)
	if !ok {
		http.NotFound(w, r)
		return
	}
	switch artifact {
	case "config":
		w.Header().Set("Content-Type", "application/xml; charset=utf-8")
		w.Header().Set("Content-Disposition", "attachment; filename=job-"+jobID+"-config.xml")
		_, _ = w.Write([]byte(job.RawXML))
	case "log":
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.Header().Set("Content-Disposition", "attachment; filename=job-"+jobID+"-events.log")
		for _, event := range h.manager.GetJobEvents(jobID) {
			_, _ = fmt.Fprintf(w, "%s [%s] %s\n", event.OccurredAt.Format(timeLayout), event.Level, event.Message)
		}
	default:
		http.NotFound(w, r)
	}
}
