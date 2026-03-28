package execution

import (
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	cfg "github.com/sine-io/cosbench-go/internal/infrastructure/config"
)

type ParsedOpConfig struct {
	Raw           map[string]string
	ContainerGen  IntGenerator
	ObjectGen     IntGenerator
	SizeGen       SizeGenerator
	ContainerPref string
	ObjectPref    string
	FilesDir      string
	FileSelection string
	PartSize      int64
	RestoreDays   int
	Delay         time.Duration
	IsPrefetch    bool
	IsRangeRequest bool
	FileLength    int64
	ChunkLength   int64
}

type SizeGenerator interface {
	Next(r *rand.Rand, idx, all int) int64
}

type scaledIntGen struct {
	base  int64
	inner IntGenerator
}

func (g scaledIntGen) Next(r *rand.Rand, idx, all int) int64 {
	return int64(g.inner.Next(r, idx, all)) * g.base
}

func ParseOpConfig(raw string) (*ParsedOpConfig, error) {
	m := cfg.ParseKVConfig(raw)
	return parseOpConfigMap(m)
}

func ParseOpConfigWithStorage(storageRaw, opRaw string) (*ParsedOpConfig, error) {
	merged := cfg.ParseKVConfig(storageRaw)
	for key, value := range cfg.ParseKVConfig(opRaw) {
		merged[key] = value
	}
	return parseOpConfigMap(merged)
}

func parseOpConfigMap(m map[string]string) (*ParsedOpConfig, error) {
	p := &ParsedOpConfig{
		Raw:           m,
		ContainerPref: valueOr(m["cprefix"], "c"),
		ObjectPref:    valueOr(m["oprefix"], "obj"),
		FilesDir:      strings.TrimSpace(m["files"]),
		FileSelection: strings.TrimSpace(m["fileselection"]),
		PartSize:      parseInt64Or(m["part_size"], 5*1024*1024),
		RestoreDays:   parseIntOr(m["restore_days"], 1),
		Delay:         parseDurationOr(m["duration"], m["delay"]),
		IsPrefetch:    parseBoolOr(m["is_prefetch"]),
		IsRangeRequest: parseBoolOr(m["is_range_request"]),
		FileLength:    parseInt64Or(m["file_length"], 0),
		ChunkLength:   parseInt64Or(m["chunk_length"], 0),
	}
	if s := strings.TrimSpace(m["containers"]); s != "" {
		g, err := ParseIntGenerator(s)
		if err != nil {
			return nil, fmt.Errorf("containers: %w", err)
		}
		p.ContainerGen = g
	}
	if s := strings.TrimSpace(m["objects"]); s != "" {
		g, err := ParseIntGenerator(s)
		if err != nil {
			return nil, fmt.Errorf("objects: %w", err)
		}
		p.ObjectGen = g
	}
	if s := strings.TrimSpace(m["sizes"]); s != "" {
		units := detectSizeUnit(s)
		base := int64(sizeUnitBase(units))
		pattern := strings.TrimSuffix(s, units)
		g, err := ParseIntGenerator(pattern)
		if err != nil {
			return nil, fmt.Errorf("sizes: %w", err)
		}
		p.SizeGen = scaledIntGen{base: base, inner: g}
	}
	return p, nil
}

type OpTarget struct {
	Bucket string
	Key    string
	Size   int64
	File   string
}

func (p *ParsedOpConfig) NextTarget(r *rand.Rand, idx, all int) OpTarget {
	var t OpTarget
	if p.ContainerGen != nil {
		t.Bucket = fmt.Sprintf("%s%d", p.ContainerPref, p.ContainerGen.Next(r, idx, all))
	}
	if p.ObjectGen != nil {
		t.Key = fmt.Sprintf("%s%d", p.ObjectPref, p.ObjectGen.Next(r, idx, all))
	}
	if p.SizeGen != nil {
		t.Size = p.SizeGen.Next(r, idx, all)
	}
	if p.FilesDir != "" {
		t.File = p.resolveFilePath(r, idx, all)
	}
	return t
}

func (p *ParsedOpConfig) ScanTargets() []OpTarget {
	containerValues := expandIntGenerator(p.ContainerGen)
	objectValues := expandIntGenerator(p.ObjectGen)
	if len(containerValues) == 0 {
		containerValues = []int{0}
	}
	if len(objectValues) == 0 {
		objectValues = []int{0}
	}

	targets := make([]OpTarget, 0, len(containerValues)*len(objectValues))
	for _, container := range containerValues {
		for _, object := range objectValues {
			var t OpTarget
			if p.ContainerGen != nil {
				t.Bucket = fmt.Sprintf("%s%d", p.ContainerPref, container)
			}
			if p.ObjectGen != nil {
				t.Key = fmt.Sprintf("%s%d", p.ObjectPref, object)
			}
			targets = append(targets, t)
		}
	}
	return targets
}

func (p *ParsedOpConfig) resolveFilePath(r *rand.Rand, idx, all int) string {
	base := strings.TrimSpace(p.FilesDir)
	if base == "" {
		return ""
	}
	info, err := os.Stat(base)
	if err == nil && !info.IsDir() {
		return base
	}
	if sel := strings.TrimSpace(p.FileSelection); sel != "" {
		return filepath.Join(base, sel)
	}
	if p.ObjectGen != nil {
		return filepath.Join(base, fmt.Sprintf("%d", p.ObjectGen.Next(r, idx, all)))
	}
	return base
}

func expandIntGenerator(gen IntGenerator) []int {
	if gen == nil {
		return nil
	}
	bounded, ok := gen.(boundedIntGenerator)
	if !ok {
		return nil
	}
	lower, upper := bounded.Bounds()
	values := make([]int, 0, upper-lower+1)
	for i := lower; i <= upper; i++ {
		values = append(values, i)
	}
	return values
}

func detectSizeUnit(s string) string {
	for _, u := range []string{"GiB", "MiB", "KiB", "GB", "MB", "KB", "B"} {
		if strings.HasSuffix(s, u) {
			return u
		}
	}
	return "B"
}

func sizeUnitBase(u string) int {
	switch u {
	case "GB":
		return 1000 * 1000 * 1000
	case "GiB":
		return 1024 * 1024 * 1024
	case "MB":
		return 1000 * 1000
	case "MiB":
		return 1024 * 1024
	case "KB":
		return 1000
	case "KiB":
		return 1024
	default:
		return 1
	}
}

func valueOr(v, def string) string {
	if strings.TrimSpace(v) == "" {
		return def
	}
	return strings.TrimSpace(v)
}
func parseIntOr(v string, def int) int {
	n, err := strconv.Atoi(strings.TrimSpace(v))
	if err != nil {
		return def
	}
	return n
}
func parseInt64Or(v string, def int64) int64 {
	n, err := strconv.ParseInt(strings.TrimSpace(v), 10, 64)
	if err != nil {
		return def
	}
	return n
}

func parseBoolOr(v string) bool {
	switch strings.ToLower(strings.TrimSpace(v)) {
	case "1", "true", "yes", "on":
		return true
	default:
		return false
	}
}

func parseDurationOr(v string, fallback string) time.Duration {
	for _, item := range []string{v, fallback} {
		item = strings.TrimSpace(item)
		if item == "" {
			continue
		}
		if d, err := time.ParseDuration(item); err == nil {
			return d
		}
		if n, err := strconv.Atoi(item); err == nil {
			return time.Duration(n) * time.Millisecond
		}
	}
	return 0
}
