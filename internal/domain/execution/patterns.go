package execution

import (
	"fmt"
	"math/rand"
	"strconv"
	"strings"
	"sync/atomic"
)

type IntGenerator interface {
	Next(r *rand.Rand, idx, all int) int
}

type boundedIntGenerator interface {
	Bounds() (int, int)
}

type constantGen struct{ value int }

func (g constantGen) Next(_ *rand.Rand, _, _ int) int { return g.value }
func (g constantGen) Bounds() (int, int)              { return g.value, g.value }

type sequentialGen struct {
	lower  int
	rangeN int
	cursor atomic.Int64
}

func (g *sequentialGen) Next(_ *rand.Rand, _, _ int) int {
	return g.lower + int(g.cursor.Add(1)-1)%g.rangeN
}
func (g *sequentialGen) Bounds() (int, int) { return g.lower, g.lower + g.rangeN - 1 }

type uniformGen struct{ lower, upper int }

func (g uniformGen) Next(r *rand.Rand, idx, all int) int {
	if all <= 0 {
		all = 1
	}
	if idx <= 0 {
		idx = 1
	}
	rangeN := g.upper - g.lower + 1
	base := rangeN / all
	extra := rangeN % all
	offset := base*(idx-1) + min(extra, idx-1)
	segment := base
	if extra >= idx {
		segment++
	}
	if segment <= 0 {
		segment = 1
	}
	return g.lower + offset + r.Intn(segment)
}
func (g uniformGen) Bounds() (int, int) { return g.lower, g.upper }

type rangeGen struct {
	lower, upper int
	cursors      []atomic.Int64
}

func newRangeGen(lower, upper int) *rangeGen { return &rangeGen{lower: lower, upper: upper} }

func (g *rangeGen) Next(_ *rand.Rand, idx, all int) int {
	if all <= 0 {
		all = 1
	}
	if idx <= 0 {
		idx = 1
	}
	if len(g.cursors) == 0 {
		g.cursors = make([]atomic.Int64, all)
	}
	rangeN := g.upper - g.lower + 1
	base := rangeN / all
	extra := rangeN % all
	offset := base*(idx-1) + min(extra, idx-1)
	segment := base
	if extra >= idx {
		segment++
	}
	if segment <= 0 {
		segment = 1
	}
	cur := g.cursors[idx-1].Add(1) - 1
	return g.lower + offset + int(cur%int64(segment))
}
func (g *rangeGen) Bounds() (int, int) { return g.lower, g.upper }

func ParseIntGenerator(pattern string) (IntGenerator, error) {
	pattern = strings.TrimSpace(pattern)
	switch {
	case strings.HasPrefix(pattern, "c("):
		a, b, err := parseBounds(pattern)
		if err != nil {
			return nil, fmt.Errorf("illegal constant distribution pattern: %s", pattern)
		}
		_ = b
		return constantGen{value: a}, nil
	case strings.HasPrefix(pattern, "u("):
		a, b, err := parseBounds(pattern)
		if err != nil {
			return nil, fmt.Errorf("illegal uniform distribution pattern: %s", pattern)
		}
		return uniformGen{lower: a, upper: b}, nil
	case strings.HasPrefix(pattern, "s("):
		a, b, err := parseBounds(pattern)
		if err != nil {
			return nil, fmt.Errorf("illegal sequential distribution pattern: %s", pattern)
		}
		return &sequentialGen{lower: a, rangeN: b - a + 1}, nil
	case strings.HasPrefix(pattern, "r(") || strings.Contains(pattern, "-"):
		a, b, err := parseRange(pattern)
		if err != nil {
			return nil, fmt.Errorf("illegal range distribution pattern: %s", pattern)
		}
		return newRangeGen(a, b), nil
	default:
		return nil, fmt.Errorf("unrecognized distribution: %s", pattern)
	}
}

func parseRange(pattern string) (int, int, error) {
	if strings.HasPrefix(pattern, "r(") {
		return parseBounds(pattern)
	}
	parts := strings.Split(pattern, "-")
	if len(parts) != 2 {
		return 0, 0, fmt.Errorf("invalid")
	}
	a, err := strconv.Atoi(strings.TrimSpace(parts[0]))
	if err != nil {
		return 0, 0, err
	}
	b, err := strconv.Atoi(strings.TrimSpace(parts[1]))
	if err != nil {
		return 0, 0, err
	}
	if a < 0 || b < 0 || a > b {
		return 0, 0, fmt.Errorf("invalid bounds")
	}
	return a, b, nil
}

func parseBounds(pattern string) (int, int, error) {
	inside := strings.TrimSuffix(strings.TrimPrefix(pattern, pattern[:2]), ")")
	parts := strings.Split(inside, ",")
	if len(parts) == 0 || len(parts) > 2 {
		return 0, 0, fmt.Errorf("invalid")
	}
	a, err := strconv.Atoi(strings.TrimSpace(parts[0]))
	if err != nil {
		return 0, 0, err
	}
	b := a
	if len(parts) == 2 {
		b, err = strconv.Atoi(strings.TrimSpace(parts[1]))
		if err != nil {
			return 0, 0, err
		}
	}
	if a < 0 || b < 0 || a > b {
		return 0, 0, fmt.Errorf("invalid bounds")
	}
	return a, b, nil
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
