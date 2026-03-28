package web

import (
	"net/http"
	"strings"

	"github.com/sine-io/cosbench-go/internal/domain"
)

func (h *Handler) driverDashboardPage(w http.ResponseWriter, r *http.Request) {
	driver, ok := h.defaultDriver()
	if !ok {
		http.NotFound(w, r)
		return
	}
	overview, _ := h.manager.GetDriverOverview(driver.ID)
	h.render(w, "driver_dashboard.html", pageData{
		Title:          "Driver Dashboard",
		DriverOverview: overview,
		RequestPath:    r.URL.Path,
	})
}

func (h *Handler) driverMissionsPage(w http.ResponseWriter, r *http.Request) {
	driver, ok := h.defaultDriver()
	if !ok {
		http.NotFound(w, r)
		return
	}
	missions := h.manager.ListDriverMissions(driver.ID)
	overview, _ := h.manager.GetDriverOverview(driver.ID)
	h.render(w, "driver_missions.html", pageData{
		Title:          "Driver Missions",
		DriverOverview: overview,
		DriverMissions: missions,
		RequestPath:    r.URL.Path,
	})
}

func (h *Handler) driverMissionDetailPage(w http.ResponseWriter, r *http.Request) {
	driver, ok := h.defaultDriver()
	if !ok {
		http.NotFound(w, r)
		return
	}
	path := strings.TrimPrefix(r.URL.Path, "/driver/missions/")
	missionID := strings.Trim(path, "/")
	mission, ok := h.manager.GetMission(missionID)
	if !ok {
		http.NotFound(w, r)
		return
	}
	overview, _ := h.manager.GetDriverOverview(driver.ID)
	h.render(w, "driver_mission_detail.html", pageData{
		Title:          "Driver Mission Detail",
		DriverOverview: overview,
		DriverMission:  mission,
		RequestPath:    r.URL.Path,
	})
}

func (h *Handler) driverWorkersPage(w http.ResponseWriter, r *http.Request) {
	driver, ok := h.defaultDriver()
	if !ok {
		http.NotFound(w, r)
		return
	}
	overview, _ := h.manager.GetDriverOverview(driver.ID)
	state, _ := h.manager.GetDriverWorkerState(driver.ID)
	h.render(w, "driver_workers.html", pageData{
		Title:             "Driver Workers",
		DriverOverview:    overview,
		DriverWorkerState: state,
		RequestPath:       r.URL.Path,
	})
}

func (h *Handler) defaultDriver() (domain.DriverNode, bool) {
	drivers := h.manager.ListDriverNodes()
	if len(drivers) == 0 {
		return domain.DriverNode{}, false
	}
	return drivers[0], true
}
