package snapshot

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/sine-io/cosbench-go/internal/domain"
)

type Store struct {
	root string
}

func New(root string) (*Store, error) {
	paths := []string{
		filepath.Join(root, "drivers"),
		filepath.Join(root, "jobs"),
		filepath.Join(root, "missions"),
		filepath.Join(root, "results"),
		filepath.Join(root, "events"),
		filepath.Join(root, "endpoints"),
		filepath.Join(root, "timelines"),
	}
	for _, path := range paths {
		if err := os.MkdirAll(path, 0o755); err != nil {
			return nil, err
		}
	}
	return &Store{root: root}, nil
}

func (s *Store) SaveJob(job domain.Job) error {
	return writeJSON(filepath.Join(s.root, "jobs", job.ID+".json"), job)
}

func (s *Store) SaveDriverNode(driver domain.DriverNode) error {
	return writeJSON(filepath.Join(s.root, "drivers", driver.ID+".json"), driver)
}

func (s *Store) SaveMission(mission domain.Mission) error {
	return writeJSON(filepath.Join(s.root, "missions", mission.ID+".json"), mission)
}

func (s *Store) SaveResult(result domain.JobResult) error {
	return writeJSON(filepath.Join(s.root, "results", result.JobID+".json"), result)
}

func (s *Store) SaveEndpoint(endpoint domain.EndpointConfig) error {
	return writeJSON(filepath.Join(s.root, "endpoints", endpoint.ID+".json"), endpoint)
}

func (s *Store) SaveEvents(jobID string, events []domain.JobEvent) error {
	return writeJSON(filepath.Join(s.root, "events", jobID+".json"), events)
}

func (s *Store) SaveTimeline(timeline domain.JobTimeline) error {
	return writeJSON(filepath.Join(s.root, "timelines", timeline.JobID+".json"), timeline)
}

func (s *Store) LoadJobs() ([]domain.Job, error) {
	var jobs []domain.Job
	if err := readAllJSON(filepath.Join(s.root, "jobs"), &jobs); err != nil {
		return nil, err
	}
	sort.Slice(jobs, func(i, j int) bool { return jobs[i].CreatedAt.After(jobs[j].CreatedAt) })
	return jobs, nil
}

func (s *Store) LoadDriverNodes() ([]domain.DriverNode, error) {
	var drivers []domain.DriverNode
	if err := readAllJSON(filepath.Join(s.root, "drivers"), &drivers); err != nil {
		return nil, err
	}
	sort.Slice(drivers, func(i, j int) bool { return drivers[i].RegisteredAt.After(drivers[j].RegisteredAt) })
	return drivers, nil
}

func (s *Store) LoadMissions() ([]domain.Mission, error) {
	var missions []domain.Mission
	if err := readAllJSON(filepath.Join(s.root, "missions"), &missions); err != nil {
		return nil, err
	}
	sort.Slice(missions, func(i, j int) bool { return missions[i].UpdatedAt.After(missions[j].UpdatedAt) })
	return missions, nil
}

func (s *Store) LoadEndpoints() ([]domain.EndpointConfig, error) {
	var endpoints []domain.EndpointConfig
	if err := readAllJSON(filepath.Join(s.root, "endpoints"), &endpoints); err != nil {
		return nil, err
	}
	sort.Slice(endpoints, func(i, j int) bool { return endpoints[i].UpdatedAt.After(endpoints[j].UpdatedAt) })
	return endpoints, nil
}

func (s *Store) LoadResult(jobID string) (domain.JobResult, error) {
	var result domain.JobResult
	err := readJSON(filepath.Join(s.root, "results", jobID+".json"), &result)
	return result, err
}

func (s *Store) LoadEvents(jobID string) ([]domain.JobEvent, error) {
	var events []domain.JobEvent
	err := readJSON(filepath.Join(s.root, "events", jobID+".json"), &events)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil, nil
		}
		return nil, err
	}
	return events, nil
}

func (s *Store) LoadTimeline(jobID string) (domain.JobTimeline, error) {
	var timeline domain.JobTimeline
	err := readJSON(filepath.Join(s.root, "timelines", jobID+".json"), &timeline)
	return timeline, err
}

func writeJSON(path string, value any) error {
	data, err := json.MarshalIndent(value, "", "  ")
	if err != nil {
		return err
	}
	tmp := path + ".tmp"
	if err := os.WriteFile(tmp, data, 0o644); err != nil {
		return err
	}
	return os.Rename(tmp, path)
}

func readJSON(path string, dest any) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, dest)
}

func readAllJSON[T any](dir string, out *[]T) error {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return fmt.Errorf("read dir %s: %w", dir, err)
	}
	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".json") {
			continue
		}
		var item T
		if err := readJSON(filepath.Join(dir, entry.Name()), &item); err != nil {
			return err
		}
		*out = append(*out, item)
	}
	return nil
}
